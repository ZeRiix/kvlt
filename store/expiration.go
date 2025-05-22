package store

import (
	"time"

	"github.com/samber/lo"
)

func InitExpiration(store *Store) {
	expirationStore := make(map[int64][]*Item)

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item Item) {
			if expireAt, ok := findExpireAtProperty(item.Value); ok {
				if expireAt <= time.Now().Unix() {
					store.Drop(item.Key)
					return
				} else if expirationStore[expireAt] == nil {
					expirationStore[expireAt] = []*Item{}
				}

				expirationStore[expireAt] = append(
					expirationStore[expireAt],
					&item,
				)
			}
		},
	)

	store.actionHooks.drop = append(
		store.actionHooks.drop,
		func(item Item) {
			if expireAt, ok := findExpireAtProperty(item.Value); ok {
				if expirationStore[expireAt] == nil {
					return
				}

				expirationStore[expireAt] = lo.Filter(
					expirationStore[expireAt],
					func(storageItem *Item, index int) bool {
						return item.Value != &storageItem.Value
					},
				)

				if len(expirationStore[expireAt]) == 0 {
					delete(expirationStore, expireAt)
				}
			}
		},
	)

	go func() {
		var lastExpirationTime int64

		for {
			oldLastExpirationTime := lastExpirationTime
			lastExpirationTime = time.Now().Unix()
			offset := lastExpirationTime - oldLastExpirationTime

			lo.Times(
				int(offset),
				func(index int) bool {
					expirationTime := int64(index) + oldLastExpirationTime

					if items, ok := expirationStore[expirationTime]; ok {
						lo.ForEach(
							items,
							func(item *Item, index int) {
								store.Drop(item.Key)
							},
						)
					}

					return true
				},
			)

		}
	}()
}

func findExpireAtProperty(data interface{}) (int64, bool) {
	if obj, ok := data.(map[string]interface{}); ok {
		if expireAt, exists := obj["expireAt"].(int); exists {
			return int64(expireAt), true
		} else if expireAt, exists := obj["expireAt"].(int64); exists {
			return expireAt, true
		} else if expireAt, exists := obj["expireAt"].(float32); exists {
			return int64(expireAt), true
		} else if expireAt, exists := obj["expireAt"].(float64); exists {
			return int64(expireAt), true
		}
		return 0, false
	}
	return 0, false
}
