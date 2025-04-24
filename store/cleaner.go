package store

// explanation: https://medium.com/@trinad536/deadlocks-in-go-f4ae0ecd05f6

import (
	"time"
)

// CleanExpiredKeys removes expired keys from the store.
func (s *Store) CleanExpiredKeys() error {
	s.mu.RLock()
	keysToDelete := make([]string, 0)
	now := time.Now().Unix()

	// Iterate over the keys and check for expiration
	for key, item := range s.data {
		if item.exp > 0 && now > item.exp {
			keysToDelete = append(keysToDelete, key)
		}
	}
	s.mu.RUnlock()

	// delete outside the mutex to save lock time - so performance gain
	count := 0
	for _, key := range keysToDelete {
		if s.DropKey(key) {
			count++
		}
	}

	return nil
}
