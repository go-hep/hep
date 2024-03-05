// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtypes

import (
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
	if n == "" {
		panic("rtypes: invalid classname")
	}

	f.mu.RLock()
	fct, ok := f.db[n]
	f.mu.RUnlock()
	if ok {
		return fct
	}

	// if we are here, nobody registered a streamer+factory for 'string'.
	// try our streamer-less rbase.String version.
	if n == "string" {
		f.mu.RLock()
		fct, ok := f.db["*rbase.String"]
		f.mu.RUnlock()
		if ok {
			return fct
		}
	}

	f.mu.RLock()
	obj := f.db["*rdict.Object"]
	f.mu.RUnlock()

	fct = func() reflect.Value {
		v := obj()
		v.Interface().(setClasser).SetClass(n)
		return v
	}

	return fct
}

func (f *factory) Add(n string, fct FactoryFct) {
	f.mu.Lock()
	f.db[n] = fct
	f.mu.Unlock()
}

type setClasser interface {
	SetClass(name string)
}

var Factory = &factory{
	db: make(map[string]FactoryFct),
}
