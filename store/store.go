package store

import (
	"sync"
)

type Store struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

var instance *Store
var once sync.Once

func Get() *Store {
	once.Do(func() {
		instance = &Store{
			data: make(map[string]interface{}),
		}
	})
	return instance
}

func (s *Store) SetValue(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) GetValue(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}
