package store

import (
	"sync"
	"time"
)

type Item struct {
	value interface{}
	iat   int64 // Issued At Time
	exp   int64 // Expiration Time
}

type Store struct {
	data     map[string]Item
	mu       sync.RWMutex
	SetValue func(key string, value interface{}, duration int64)
}

var instance *Store
var once sync.Once

// Get returns the singleton instance of Store.
func Get() *Store {
	once.Do(func() {
		s := &Store{
			data: make(map[string]Item),
		}

		// SetValue is a method to set a value in the store with an expiration duration.
		s.SetValue = func(key string, value interface{}, duration int64) {
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

		instance = s
	})
	return instance
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
