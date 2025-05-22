package store

type Item struct {
	Value interface{}
	Key   string
}

type ActionHooks struct {
	get  []func(item Item)
	set  []func(item Item)
	drop []func(item Item)
}

type Store struct {
	data        map[string]Item
	actionHooks ActionHooks
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Item),
		actionHooks: ActionHooks{
			get:  []func(item Item){},
			set:  []func(item Item){},
			drop: []func(item Item){},
		},
	}
}

func launchHook(item Item, hooks []func(item Item)) {
	for _, hook := range hooks {
		go hook(item)
	}
}

func (store *Store) Get(key string) (Item, bool) {
	item, err := store.data[key]

	go launchHook(item, store.actionHooks.get)

	return item, err
}

func (store *Store) Set(item Item) Item {
	store.data[item.Key] = item

	go launchHook(item, store.actionHooks.set)

	return item
}

func (store *Store) Drop(key string) (Item, bool) {
	item, exist := store.data[key]

	if exist {
		delete(store.data, key)

		go launchHook(item, store.actionHooks.drop)

		return item, true
	}
	return item, false
}
