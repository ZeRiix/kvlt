package main

import (
	"kvlt/store"
	"time"
)

func main() {

	optionsAOF := store.OptionsAOF{
		IntervalAnalyzeBuffer: 1 * time.Second,
		IntervalSnapshot:      10 * time.Second,
		QuantityBuffer:        10,
		AOFFolderPath:         "./buffer",
		SnapshotFolderPath:    "./data",
		SplitChar:             "|\\|\\|",
	}

	storeInstance := store.NewStore()

	store.InitExpiration(storeInstance)
	store.InitAOF(storeInstance, optionsAOF)
	// store.InitIndexes(storeInstance)

	storeInstance.Set(store.Item{
		Key: "test",
		Value: map[string]interface{}{
			"firstname": "john",
			"lastname":  "doe",
			"toto": map[string]interface{}{
				"hihi": 111,
			},
			"expireAt": time.Now().Unix() + 20,
		},
	})

	storeInstance.Drop("test1")

	select {}
}
