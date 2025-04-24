package store

import "time"

func CleanExpiredKeys() {
	if instance == nil {
		return
	}

	instance.mu.Lock()
	defer instance.mu.Unlock()

	for key, value := range instance.data {
		if value.exp != 0 && value.exp < time.Now().Unix() {
			delete(instance.data, key)
		}
	}
}
