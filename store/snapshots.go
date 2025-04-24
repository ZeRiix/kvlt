package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

func MakeSnapshot() error {
	if instance == nil || len(instance.data) == 0 {
		log.Println("Attention: Store vide ou non initialisé lors de la création du snapshot")
	}

	type SnapshotData struct {
		Data      map[string]interface{} `json:"data"`
		Timestamp string                 `json:"timestamp"`
	}

	snapshotData := SnapshotData{
		Data:      make(map[string]interface{}),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if instance != nil {
		for key, value := range instance.data {
			snapshotData.Data[key] = value
		}
	}

	snapshotDir := "snapshots"
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		return fmt.Errorf("erreur lors de la création du dossier snapshots: %w", err)
	}

	data, err := json.MarshalIndent(snapshotData, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation des données du store: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(snapshotDir, fmt.Sprintf("store_snapshot_%s.json", timestamp))

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier de snapshot: %w", err)
	}

	log.Printf("Snapshot du store enregistré dans %s avec %d entrées", filename, len(snapshotData.Data))
	return nil
}

func LoadSnapshot(filename string) error {
	snapshotDir := "snapshots"
	filePath := filepath.Join(snapshotDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("le fichier de snapshot '%s' n'existe pas", filename)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier snapshot: %w", err)
	}

	type SnapshotData struct {
		Data      map[string]interface{} `json:"data"`
		Timestamp string                 `json:"timestamp"`
	}

	var snapshotData SnapshotData
	if err := json.Unmarshal(data, &snapshotData); err != nil {
		return fmt.Errorf("erreur lors de la désérialisation du snapshot: %w", err)
	}

	if len(snapshotData.Data) == 0 {
		log.Println("Attention: Le snapshot chargé ne contient aucune donnée")
	}

	if instance == nil {
		instance = &Store{
			data: make(map[string]interface{}),
			mu:   sync.RWMutex{},
		}
	}

	instance.mu.Lock()
	defer instance.mu.Unlock()

	instance.data = make(map[string]interface{})
	for key, value := range snapshotData.Data {
		instance.data[key] = value
	}

	log.Printf("Snapshot '%s' chargé avec succès (%d entrées, timestamp: %s)",
		filename, len(snapshotData.Data), snapshotData.Timestamp)

	return nil
}

func ListSnapshots() ([]string, error) {
	snapshotDir := "snapshots"

	if _, err := os.Stat(snapshotDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(snapshotDir)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du dossier snapshots: %w", err)
	}

	var snapshots []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			snapshots = append(snapshots, entry.Name())
		}
	}

	sort.Slice(snapshots, func(i, j int) bool {
		fileI, _ := os.Stat(filepath.Join(snapshotDir, snapshots[i]))
		fileJ, _ := os.Stat(filepath.Join(snapshotDir, snapshots[j]))
		return fileI.ModTime().After(fileJ.ModTime())
	})

	return snapshots, nil
}

func GetLatestSnapshot() (string, error) {
	snapshots, err := ListSnapshots()
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération des snapshots: %w", err)
	}

	if len(snapshots) == 0 {
		return "", nil
	}

	return snapshots[0], nil
}
