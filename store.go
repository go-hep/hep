package fwk

import (
	"reflect"
)

type achan chan interface{}

type datastore struct {
	SvcBase
	store map[string]achan
}

func (ds *datastore) Configure(ctx Context) Error {
	ds.store = make(map[string]achan)
	return nil
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

func init() {
	Register(reflect.TypeOf(datastore{}))
}

// EOF
