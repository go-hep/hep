package fwk

type achan chan interface{}

type datastore struct {
	Base
	store map[string]achan
}

func newDataStore(n string) *datastore {
	ds := &datastore{
		Base: Base{
			Name: n,
			Type: "fwk.datastore",
		},
		store: make(map[string]achan),
	}
	return ds
}

func (ds *datastore) Get(k string) (interface{}, Error) {
	//fmt.Printf(">>> get(%v)...\n", k)
	ch, ok := ds.store[k]
	if !ok {
		return nil, Errorf("Store.Get: no such key [%v]", k)
	}
	v := <-ch
	ch <- v
	//fmt.Printf("<<< get(%v, %v)...\n", k, v)
	return v, nil
}

func (ds *datastore) Put(k string, v interface{}) Error {
	//fmt.Printf(">>> put(%v, %v)...\n", k, v)
	ds.store[k] <- v
	//fmt.Printf("<<< put(%v, %v)...\n", k, v)
	return nil
}

func (ds *datastore) StartSvc(ctx Context) Error {
	ds.store = make(map[string]achan)
	return nil
}

func (ds *datastore) StopSvc(ctx Context) Error {
	ds.store = nil
	return nil
}

// EOF
