package rootio

import (
	"reflect"
)

type FactoryFct func() reflect.Value

type factory struct {
	db map[string]FactoryFct // a registry of all factory functions by type name
}

func (f *factory) NumKey() int {
	return len(f.db)
}

func (f *factory) Keys() []string {
	keys := make([]string, 0, len(f.db))
	for k, _ := range f.db {
		keys = append(keys, k)
	}
	return keys
}

func (f *factory) HasKey(n string) bool {
	_, ok := f.db[n]
	return ok
}

func (f *factory) Get(n string) FactoryFct {
	fct, ok := f.db[n]
	if ok {
		return fct
	}
	return nil
}

// the registry of all factory functions, by type name
var Factory = factory{
	db: make(map[string]FactoryFct),
}

// EOF
