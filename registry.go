package fwk

import (
	"reflect"
	"sort"
)

type factoryDb map[string]reflect.Type

var g_factory factoryDb = make(factoryDb)

func Register(t reflect.Type) {
	comp := t.PkgPath() + "." + t.Name()
	g_factory[comp] = t
	//fmt.Printf("### factories ###\n%v\n", g_factory)
}

func Registry() []string {
	comps := make([]string, 0, len(g_factory))
	for k, _ := range g_factory {
		comps = append(comps, k)
	}
	sort.Strings(comps)
	return comps
}

func New(t, n string) (Component, Error) {
	var err Error
	rt, ok := g_factory[t]
	if !ok {
		return nil, Errorf("no component with type [%s] registered\n", t)
	}
	c := reflect.New(rt).Interface().(Component)
	c.SetName(n)
	if g_mgr != nil {
		err = g_mgr.addComponent(c)
	}
	return c, err
}

// EOF
