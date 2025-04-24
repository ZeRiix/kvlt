package store

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type storedItem struct {
	Value interface{} `json:"value"`
	Iat   int64       `json:"iat"`
	Exp   int64       `json:"exp"`
}

type storedDb struct {
	Data map[string]storedItem `json:"data"`
}

// SaveToFile save data in json format to a file
func (s *Store) SaveToFile(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	s.mu.RLock()
	data := storedDb{
		Data: make(map[string]storedItem, len(s.data)),
	}

	for key, item := range s.data {
		data.Data[key] = storedItem{
			Value: item.value,
			Iat:   item.iat,
			Exp:   item.exp,
		}
	}

	s.mu.RUnlock()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}

// LoadFromFile load data from a json file
func (s *Store) LoadFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var storedData storedDb
	if err := json.Unmarshal(data, &storedData); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, item := range storedData.Data {
		s.data[key] = Item{
			value: item.Value,
			iat:   item.Iat,
			exp:   item.Exp,
		}
	}

	return nil
}

// EnableAutoPersistence Activate the auto-persistence feature
func (s *Store) EnableAutoPersistence(filePath string) {
	_ = s.LoadFromFile(filePath)

	originalSetValue := s.SetValue
	s.SetValue = func(key string, value interface{}, duration int64) {
		originalSetValue(key, value, duration)
		go func() {
			if err := s.SaveToFile(filePath); err != nil {
				log.Printf("error saving to file: %v", err)
			}
		}()
	}

	originalDelete := s.DropKey
	s.DropKey = func(key string) bool {
		result := originalDelete(key)
		if result {
			go func() {
				if err := s.SaveToFile(filePath); err != nil {
					log.Printf("error saving to file after deletion: %v", err)
				}
			}()
		}
		return result
	}
}
