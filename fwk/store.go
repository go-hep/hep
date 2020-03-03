// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"fmt"
	"reflect"
)

type achan chan interface{}

// datastore stores (event) data and provides concurrent-safe access to it.
type datastore struct {
	SvcBase
	store map[string]achan
	quit  chan struct{}
}

func (ds *datastore) Configure(ctx Context) error {
	return nil
}

func (ds *datastore) Get(k string) (interface{}, error) {
	ch, ok := ds.store[k]
	if !ok {
		return nil, fmt.Errorf("Store.Get: no such key [%v]", k)
	}
	select {
	case v, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("%s: closed channel for key [%s]", ds.Name(), k)
		}
		ch <- v
		return v, nil
	case <-ds.quit:
		return nil, fmt.Errorf("%s: timeout to get [%s]", ds.Name(), k)
	}
}

func (ds *datastore) Put(k string, v interface{}) error {
	select {
	case ds.store[k] <- v:
		return nil
	case <-ds.quit:
		return fmt.Errorf("%s: timeout to put [%s]", ds.Name(), k)
	}
}

func (ds *datastore) Has(k string) bool {
	_, ok := ds.store[k]
	return ok
}

func (ds *datastore) StartSvc(ctx Context) error {
	ds.store = make(map[string]achan)
	return nil
}

func (ds *datastore) StopSvc(ctx Context) error {
	ds.store = make(map[string]achan)
	return nil
}

// reset deletes the payload and resets the associated channel
func (ds *datastore) reset(keys []string) error {
	var err error
	for _, k := range keys {
		ch, ok := ds.store[k]
		if ok {
			select {
			case vv := <-ch:
				if vv, ok := vv.(Deleter); ok {
					err = vv.Delete()
					if err != nil {
						return err
					}
				}
			default:
			}
		}
		ds.store[k] = make(achan, 1)
	}
	ds.quit = make(chan struct{})
	return err
}

// close notifies components hanging on store.Get or .Put that event has been aborted
func (ds *datastore) close() {
	close(ds.quit)
}

func init() {
	Register(reflect.TypeOf(datastore{}),
		func(typ, name string, mgr App) (Component, error) {
			return &datastore{
				SvcBase: NewSvc(typ, name, mgr),
				store:   make(map[string]achan),
				quit:    make(chan struct{}),
			}, nil
		},
	)
}

// interface tests
var _ Store = (*datastore)(nil)
