package store

import (
	"sync"
	"time"
)

type Item struct {
	value interface{}
	iat   int64 // Issued At Time (timestamp de création)
	exp   int64 // Expiration Time (timestamp d'expiration)
}

type Store struct {
	data map[string]Item
	mu   sync.RWMutex
}

var instance *Store
var once sync.Once

// Get returns the singleton instance of Store.
func Get() *Store {
	once.Do(func() {
		instance = &Store{
			data: make(map[string]Item),
		}
	})
	return instance
}

// SetValue sets the value associated with the key and an expiration duration.
func (s *Store) SetValue(key string, value interface{}, duration int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	exp := now + duration
	item := Item{
		value: value,
		iat:   now,
		exp:   exp,
	}

	s.data[key] = item
}

// GetItem retrieves the Item associated with the key.
func (s *Store) GetItem(key string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.data[key]
	if !ok {
		return Item{}, false
	}

	if item.exp > 0 && time.Now().Unix() > item.exp {
		return Item{}, false
	}

	return item, true
}

// GetValue retrieves the value associated with the key.
func (s *Store) GetValue(key string) (interface{}, bool) {
	item, ok := s.GetItem(key)
	if !ok {
		return nil, false
	}

	return item.value, true
}
