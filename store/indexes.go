package store

import (
	"log"

	"github.com/samber/lo"
)

type Any = interface{}

type RecordKey interface {
	string | int64 | bool
}

type Record[Key RecordKey, Value any] = map[Key]Value

type RecordItem = Record[string, *Item]

type Index struct {
	intStore     Record[int64, RecordItem]
	stringStore  Record[string, RecordItem]
	nullStore    RecordItem
	booleanStore Record[bool, RecordItem]
}

type Indexes Record[string, Index]

func flatten(input Any) Any {
	if objectInput, ok := input.(Record[string, Any]); ok {
		return lo.Reduce(
			lo.Entries(objectInput),
			func(acc Record[string, Any], entry lo.Entry[string, Any], index int) Record[string, Any] {
				flatted := flatten(entry.Value)
				if objectFlatted, ok := flatted.(Record[string, Any]); ok {
					return lo.Reduce(
						lo.Entries(objectFlatted),
						func(acc Record[string, Any], subEntry lo.Entry[string, Any], index int) Record[string, Any] {
							acc[entry.Key+"."+subEntry.Key] = subEntry.Value
							return acc
						},
						acc,
					)
				} else {
					acc[entry.Key] = flatted
					return acc
				}
			},
			make(Record[string, Any]),
		)
	}

	return input
}

func sortIndex(indexes *Indexes, item *Item, key string, value Any) {
	if _, exist := (*indexes)[key]; !exist {
		(*indexes)[key] = Index{
			intStore:     make(Record[int64, RecordItem]),
			stringStore:  make(Record[string, RecordItem]),
			nullStore:    make(RecordItem),
			booleanStore: make(Record[bool, RecordItem]),
		}
	}

	switch typedValue := value.(type) {
	case string:
		if _, exist := (*indexes)[key].stringStore[typedValue]; !exist {
			(*indexes)[key].stringStore[typedValue] = make(RecordItem)
		}

		(*indexes)[key].stringStore[typedValue][item.Key] = item
	case int64:
		if _, exist := (*indexes)[key].intStore[typedValue]; !exist {
			(*indexes)[key].intStore[typedValue] = make(RecordItem)
		}

		(*indexes)[key].intStore[typedValue][item.Key] = item
	case bool:
		if _, exist := (*indexes)[key].booleanStore[typedValue]; !exist {
			(*indexes)[key].booleanStore[typedValue] = make(RecordItem)
		}

		(*indexes)[key].booleanStore[typedValue][item.Key] = item
	case nil:
		(*indexes)[key].nullStore[item.Key] = item
	default:
		log.Println("unsupport value: %#v", typedValue)
		return
	}
}

func deleteIndex(indexes *Indexes, item *Item, key string, value Any) {
	switch typedValue := value.(type) {
	case string:
		delete((*indexes)[key].stringStore[typedValue], item.Key)
	case int64:
		delete((*indexes)[key].intStore[typedValue], item.Key)
	case bool:
		delete((*indexes)[key].booleanStore[typedValue], item.Key)
	case nil:
		delete((*indexes)[key].nullStore, item.Key)
	default:
		log.Println("unsupport value: %#v", typedValue)
		return
	}
}

func InitIndexes(store *Store) func(path string, value Any) []*Item {
	indexes := make(Indexes)

	store.actionHooks.set = append(
		store.actionHooks.set,
		func(item *Item) {
			flattedValue := flatten(item.Value)

			log.Printf("flattedValue: %#v", flattedValue)

			if flattedObject, ok := flattedValue.(Record[string, Any]); ok {
				lo.ForEach(
					lo.Entries(flattedObject),
					func(subEntry lo.Entry[string, Any], index int) {
						sortIndex(&indexes, item, subEntry.Key, subEntry.Value)
					},
				)
			} else {
				sortIndex(&indexes, item, "", flattedValue)
			}
		},
	)

	store.actionHooks.drop = append(
		store.actionHooks.drop,
		func(item *Item) {
			flattedValue := flatten(item.Value)

			if flattedObject, ok := flattedValue.(Record[string, Any]); ok {
				lo.ForEach(
					lo.Entries(flattedObject),
					func(entry lo.Entry[string, Any], index int) {
						deleteIndex(&indexes, item, entry.Key, entry.Value)
					},
				)
			} else {
				deleteIndex(&indexes, item, "", flattedValue)
			}
		},
	)

	return func(key string, value Any) []*Item {
		if _, exist := indexes[key]; !exist {
			return make([]*Item, 0)
		}

		switch typedValue := value.(type) {
		case string:
			if _, exist := indexes[key].stringStore[typedValue]; !exist {
				return make([]*Item, 0)
			}

			return lo.Values(indexes[key].stringStore[typedValue])
		case int64:
			if _, exist := indexes[key].intStore[typedValue]; !exist {
				return make([]*Item, 0)
			}

			return lo.Values(indexes[key].intStore[typedValue])
		case bool:
			if _, exist := indexes[key].booleanStore[typedValue]; !exist {
				return make([]*Item, 0)
			}

			return lo.Values(indexes[key].booleanStore[typedValue])
		case nil:

			return lo.Values(indexes[key].nullStore)
		default:
			log.Println("unsupport value: %#v", typedValue)
			return make([]*Item, 0)
		}
	}
}
