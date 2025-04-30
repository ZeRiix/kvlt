package store

import (
	"sync"
	"time"
)

type Cleaner struct {
	store     *Store
	interval  time.Duration
	stopChan  chan struct{}
	waitGroup *sync.WaitGroup
}

func CleanExpiredKeys(store *Store) error {
	now := time.Now().Unix()

	for key, item := range store.data {
		if item.Exp > 0 && now > item.Exp {
			store.Drop(key)
		}
	}

	return nil
}

func StartCleaner(store *Store, interval time.Duration) *Cleaner {
	cleaner := &Cleaner{
		store:     store,
		interval:  interval,
		stopChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}

	temp := make(map[int64][]*Item)

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item Item) {
			if temp[item.Exp] == nil {
				temp[item.Exp] = []*Item{}
			}
			temp[item.Exp] = append(temp[item.Exp], &item)
		},
	)

	cleaner.waitGroup.Add(1)
	go func() {
		defer cleaner.waitGroup.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			if items, ok := temp[time.Now().Unix()]; ok {
				for _, item := range items {
					store.Drop(item.Key)
				}
				delete(temp, time.Now().Unix())
			}
		}
	}()

	return cleaner
}

func (c *Cleaner) Stop() {
	close(c.stopChan)
	c.waitGroup.Wait()
}
