package store

import (
	"log"
)

type Index struct {
	intStore     map[int64]map[string]*Item
	stringStore  map[string]map[string]*Item
	nullStore    map[string]*Item
	booleanStore map[bool]map[string]*Item
}

type Key = string

type Indexes map[Key]Index

func flatten(input interface{}) interface{} {
	if objectInput, ok := input.(map[string]interface{}); ok {
		newObject := make(map[string]interface{})

		for key, value := range objectInput {
			flatted := flatten(value)

			if objectFlatted, ok := flatted.(map[string]interface{}); ok {
				for subKey, value := range objectFlatted {
					newObject[key+"."+subKey] = value
				}
			} else {
				newObject[key] = flatted
			}

		}

		return newObject
	}

	return input
}

func sortIndex(indexes *Indexes, item *Item, key string, value interface{}) {
	if _, exist := (*indexes)[key]; !exist {
		(*indexes)[key] = Index{
			intStore:     make(map[int64]map[string]*Item),
			stringStore:  make(map[string]map[string]*Item),
			nullStore:    make(map[string]*Item),
			booleanStore: make(map[bool]map[string]*Item),
		}
	}

	switch typedValue := value.(type) {
	case string:
		(*indexes)[key].stringStore[typedValue][item.Key] = item
	case int64:
		(*indexes)[key].intStore[typedValue][item.Key] = item
	case bool:
		(*indexes)[key].booleanStore[typedValue][item.Key] = item
	case nil:
		(*indexes)[key].nullStore[item.Key] = item
	default:
		log.Println(typedValue)
		return
	}
}

func InitIndexes(store *Store) {
	indexes := make(Indexes)

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item Item) {
			flattedValue := flatten(item.Value)

			if flattedObject, ok := flattedValue.(map[string]interface{}); ok {
				for key, value := range flattedObject {
					sortIndex(&indexes, &item, key, value)
				}
			} else {
				sortIndex(&indexes, &item, "", flattedValue)
			}
		},
	)

	store.actionHooks.drop = append(
		store.actionHooks.drop,
		func(item Item) {
			flattedValue := flatten(item.Value)

			if flattedObject, ok := flattedValue.(map[string]interface{}); ok {
				for key, value := range flattedObject {
					sortIndex(&indexes, &item, key, value)
				}
			} else {

			}
		},
	)
}
