package store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Action string

const (
	SET  Action = "set"
	DROP Action = "drop"
)

type Operation struct {
	action     Action
	key        string
	stringItem string
}

func initOperation(action Action, item Item) Operation {
	dataStringify, _ := json.Marshal(item)

	return Operation{
		action:     action,
		key:        item.Key,
		stringItem: string(dataStringify),
	}
}

type OptionAOF struct {
	IntervalAnalyzeBuffer time.Duration
	IntervalSnapshot      time.Duration
	QuantityBuffer        int
	AofFolderPath         string
	SnapshotFolderPath    string
}

func InitAOF(store *Store, options OptionAOF) {
	indexCurrentOperation := 0
	var operationsList [][]Operation

	for i := 0; i < int(options.QuantityBuffer); i++ {
		var operations []Operation

		operationsList = append(
			operationsList,
			operations,
		)
	}

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item Item) {
			operationsList[indexCurrentOperation] = append(
				operationsList[indexCurrentOperation],
				initOperation(SET, item),
			)
		},
	)

	store.actionHooks.drop = append(
		store.actionHooks.drop,
		func(item Item) {
			operationsList[indexCurrentOperation] = append(
				operationsList[indexCurrentOperation],
				initOperation(DROP, item),
			)
		},
	)

	go (func() {
		for {
			time.Sleep(options.IntervalAnalyzeBuffer)

			lastIndexOperation := indexCurrentOperation
			indexCurrentOperation = (indexCurrentOperation + 1) % options.QuantityBuffer

			operations := operationsList[lastIndexOperation]

			var newOperation []Operation
			operationsList[lastIndexOperation] = newOperation

			go (func() {
				if len(operations) == 0 {
					return
				}

				var content strings.Builder
				for i := len(operations) - 1; i >= 0; i-- {
					operation := operations[i]
					content.WriteString(string(operation.action) + "\\" + operation.key + "\\" + operation.stringItem + "\n")
				}

				pathFile := options.AofFolderPath + "/" + time.Now().String()
				os.WriteFile(
					pathFile,
					[]byte(content.String()),
					0644,
				)
			})()
		}
	})()

	go (func() {
		for {
			time.Sleep(options.IntervalSnapshot)

			snapshotFiles(options.AofFolderPath, options.SnapshotFolderPath)
		}
	})()

}

func snapshotFiles(aofPath, snapshotPath string) {
	files, err := os.ReadDir(aofPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture du répertoire AOF: %v", err)
		return
	}

	fileNames := make([]string, len(files))
	for i, file := range files {
		fileNames[i] = file.Name()
	}

	sort.Slice(fileNames, func(i, j int) bool {
		numI, errI := strconv.ParseInt(fileNames[i], 10, 64)
		numJ, errJ := strconv.ParseInt(fileNames[j], 10, 64)

		if errI != nil || errJ != nil {
			return fileNames[i] < fileNames[j]
		}

		return numI < numJ
	})

	for _, fileName := range fileNames {
		bytesContent, err := os.ReadFile(filepath.Join(aofPath, fileName))
		if err != nil {
			log.Printf("Erreur lors de la lecture du fichier %s: %v", fileName, err)
			continue
		}

		for _, content := range strings.Split(string(bytesContent), "\n") {
			if content == "" {
				continue
			}

			parts := strings.Split(content, "\\")
			if len(parts) != 3 {
				log.Printf("Format invalide dans %s: %s", fileName, content)
				continue
			}

			action, key, data := Action(parts[0]), parts[1], parts[2]
			filePath := filepath.Join(snapshotPath, key)

			switch action {
			case DROP:
				if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
					log.Printf("Erreur lors de la suppression de %s: %v", key, err)
				}
			case SET:
				if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
					log.Printf("Erreur lors de l'écriture de %s: %v", key, err)
				}
			default:
				log.Printf("Action non implémentée: %s", action)
			}
		}

		os.Remove(filepath.Join(aofPath, fileName))
	}
}

func LoadSnapshot(store *Store, snapshotPath string) {
	files, err := os.ReadDir(snapshotPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture du répertoire de snapshot: %v", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(snapshotPath, file.Name())
		bytesContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Erreur lors de la lecture du fichier %s: %v", file.Name(), err)
			continue
		}

		var item Item
		if err := json.Unmarshal(bytesContent, &item); err != nil {
			log.Printf("Erreur lors de la désérialisation de %s: %v", file.Name(), err)
			continue
		}

		store.data[item.Key] = item
	}
	log.Printf("Chargement des snapshots terminé.")
}
