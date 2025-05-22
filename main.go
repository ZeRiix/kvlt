package main

import (
	"kvlt/store"
	"time"
)

func main() {

	storeInstance := store.NewStore()

	store.InitIndexes(storeInstance)

	store.LoadSnapshot(storeInstance, "./data")

	store.StartCleaner(storeInstance, 15*time.Second)

	options := store.OptionAOF{
		IntervalAnalyzeBuffer: 1 * time.Second,
		IntervalSnapshot:      10 * time.Second,
		QuantityBuffer:        10,
		AofFolderPath:         "./buffer",
		SnapshotFolderPath:    "./data",
	}

	store.InitAOF(storeInstance, options)

	storeInstance.Set(store.Item{
		Key: "test",
		Value: map[string]interface{}{
			"firstname": "john",
			"lastname":  "doe",
			"toto": map[string]interface{}{
				"hihi": 111,
			},
		},
		Exp: time.Now().Unix() + 20,
	})

	storeInstance.Set(store.Item{
		Key:   "test1",
		Value: 1000,
		Exp:   time.Now().Unix() + 1000,
	})

	storeInstance.Set(store.Item{
		Key:   "test2",
		Value: "test",
		Exp:   time.Now().Unix() + 1000,
	})

	storeInstance.Drop("test1")

	select {}
}
