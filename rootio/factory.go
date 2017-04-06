// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"log"
	"reflect"
	"sync"
)

// FactoryFct creates new values of a given type.
type FactoryFct func() reflect.Value

type factory struct {
	mu sync.RWMutex
	db map[string]FactoryFct // a registry of all factory functions by type name
}

func (f *factory) Len() int {
	f.mu.RLock()
	n := len(f.db)
	f.mu.RUnlock()
	return n
}

func (f *factory) Keys() []string {
	f.mu.RLock()
	keys := make([]string, 0, len(f.db))
	for k := range f.db {
		keys = append(keys, k)
	}
	f.mu.RUnlock()
	return keys
}

func (f *factory) HasKey(n string) bool {
	f.mu.RLock()
	_, ok := f.db[n]
	f.mu.RUnlock()
	return ok
}

func (f *factory) Get(n string) FactoryFct {
	f.mu.RLock()
	fct, ok := f.db[n]
	f.mu.RUnlock()
	if ok {
		return fct
	}

	fct = func() reflect.Value {
		o := &dobject{class: n}
		return reflect.ValueOf(o)
	}
	Factory.add(n, fct)
	log.Printf("rootio: adding dummy factory for %q\n", n)
	return fct
}

func (f *factory) get(n string) FactoryFct {
	fct := f.Get(n)
	if fct == nil {
		panic(fmt.Errorf("rootio: no factory for type %q", n))
	}
	return fct
}

func (f *factory) add(n string, fct FactoryFct) {
	f.mu.Lock()
	f.db[n] = fct
	f.mu.Unlock()
}

var Factory = factory{
	db: make(map[string]FactoryFct),
}
