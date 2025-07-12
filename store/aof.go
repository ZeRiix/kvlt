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

	"github.com/samber/lo"
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

func initOperation(action Action, item *Item) Operation {
	dataStringify, _ := json.Marshal(item)

	return Operation{
		action:     action,
		key:        item.Key,
		stringItem: string(dataStringify),
	}
}

type OptionsAOF struct {
	IntervalAnalyzeBuffer time.Duration
	IntervalSnapshot      time.Duration
	QuantityBuffer        int
	AOFFolderPath         string
	SnapshotFolderPath    string
	SplitChar             string
}

func InitAOF(store *Store, options OptionsAOF) {
	createFolder(options.AOFFolderPath)
	createFolder(options.SnapshotFolderPath)

	indexCurrentOperation := 0

	operationsList := lo.Times(
		int(options.QuantityBuffer),
		func(i int) []Operation {
			var operations []Operation

			return operations
		},
	)

	lo.ForEach(
		loadSnapshots(options),
		func(item Item, index int) {
			store.Set(item)
		},
	)

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item *Item) {
			operationsList[indexCurrentOperation] = append(
				operationsList[indexCurrentOperation],
				initOperation(SET, item),
			)
		},
	)

	store.actionHooks.drop = append(
		store.actionHooks.drop,
		func(item *Item) {
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

			exportOperations(options, operations)
		}
	})()

	go (func() {
		for {
			time.Sleep(options.IntervalSnapshot)

			applyAllAOF(options)
		}
	})()

}

func exportOperations(options OptionsAOF, operations []Operation) {
	if len(operations) == 0 {
		return
	}

	content := lo.Reduce(
		operations,
		func(acc string, operation Operation, index int) string {
			return acc + string(operation.action) + options.SplitChar + operation.key + options.SplitChar + operation.stringItem + "\n"
		},
		string(""),
	)

	pathFile := filepath.Join(options.AOFFolderPath, strconv.FormatInt(time.Now().Unix(), 10))
	os.WriteFile(
		pathFile,
		[]byte(content),
		0644,
	)
}

func applyAllAOF(options OptionsAOF) {
	files, err := os.ReadDir(options.AOFFolderPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture du répertoire AOF: %v", err)
		return
	}

	intFileNames := lo.Map(
		files,
		func(item os.DirEntry, index int) int {
			num, err := strconv.Atoi(item.Name())

			if err != nil {
				log.Printf("Erreur le nom d'un fichier AOF n'est pas convertible en int: %v", err)
				os.Exit(1)
			}

			return num
		},
	)

	sort.Ints(intFileNames)

	filePaths := lo.Map(
		intFileNames,
		func(intFileName int, index int) string {
			return filepath.Join(
				options.AOFFolderPath,
				strconv.Itoa(intFileName),
			)
		},
	)

	lo.ForEach(
		filePaths,
		func(filePath string, index int) {
			applyAOF(options, filePath)
		},
	)
}

func applyAOF(options OptionsAOF, filePath string) {
	bytesContent, err := os.ReadFile(filePath)

	if err != nil {
		log.Printf("Erreur lors de la lecture du fichier %s: %v", filePath, err)
		os.Exit(1)
	}

	lo.ForEach(
		strings.Split(string(bytesContent), "\n"),
		func(stringOperation string, index int) {
			if stringOperation == "" {
				return
			}

			parts := strings.Split(stringOperation, options.SplitChar)

			if len(parts) != 3 {
				log.Printf("Format invalide dans %s: %s %s", filePath, stringOperation, parts)
				os.Exit(1)
			}

			action, key, data := Action(parts[0]), parts[1], parts[2]

			itemPath := filepath.Join(options.SnapshotFolderPath, key)

			switch action {
			case DROP:
				if err := os.Remove(itemPath); err != nil && !os.IsNotExist(err) {
					log.Printf("Erreur lors de la suppression de %s: %v", key, err)
					os.Exit(1)
				}
			case SET:
				if err := os.WriteFile(itemPath, []byte(data), 0644); err != nil {
					log.Printf("Erreur lors de l'écriture de %s: %v", key, err)
					os.Exit(1)
				}
			default:
				log.Printf("Action non implémentée: %s", action)
				os.Exit(1)
			}
		},
	)

	os.Remove(filePath)
}

func loadSnapshots(options OptionsAOF) []Item {
	applyAllAOF(options)

	files, err := os.ReadDir(options.SnapshotFolderPath)
	if err != nil {
		log.Printf("Erreur lors de la lecture du répertoire de snapshot: %v", err)
		os.Exit(1)
	}

	data := lo.Map(
		files,
		func(file os.DirEntry, index int) Item {
			pathFile := filepath.Join(options.SnapshotFolderPath, file.Name())
			bytesContent, err := os.ReadFile(pathFile)
			if err != nil {
				log.Printf("Erreur lors de la lecture du fichier %s: %v", file.Name(), err)
				os.Exit(1)
			}

			var item Item
			if err := json.Unmarshal(bytesContent, &item); err != nil {
				log.Printf("Erreur lors de la désérialisation de %s: %v", file.Name(), err)
				os.Exit(1)
			}

			return item
		},
	)

	log.Printf("Chargement des snapshots terminé.")

	return data
}

func createFolder(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			log.Printf("Erreur lors de la création : %s", err)
			os.Exit(1)
		}
	}
}
