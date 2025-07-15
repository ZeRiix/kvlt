package main

import (
	"kvlt/store"
	"log"
	"time"

	"github.com/samber/lo"
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
	finder := store.InitIndexes(storeInstance)

	storeInstance.Set(store.Item{
		Key:   "test2",
		Value: int64(12),
	})

	storeInstance.Drop("test2")

	storeInstance.Set(store.Item{
		Key: "test",
		Value: map[string]interface{}{
			"firstname": "john",
			"lastname":  "doe",
			"toto": map[string]interface{}{
				"hihi": int64(111),
				"test": map[string]interface{}{
					"deep": "ok",
				},
			},
			"expireAt": time.Now().Unix() + 20,
		},
	})

	time.Sleep(1 * time.Millisecond)

	result := lo.Map(
		finder("toto.hihi", int64(111)),
		func(item *store.Item, i int) store.Item {
			val := *item
			return val
		},
	)

	log.Println("findItems: %#v\n", result)

	storeInstance.Drop("test1")

	select {}
}
