package fwk

import (
	"reflect"
	"sort"
)

// FactoryFunc creates a Component of type t and name n, managed by the fwk.App mgr.
type FactoryFunc func(t, n string, mgr App) (Component, error)

// factoryDb associates a fully-qualified type-name (pkg-path + type-name) with
// a component factory-function.
type factoryDb map[string]FactoryFunc

var gFactory factoryDb = make(factoryDb)

// Register registers a type t with the FactoryFunc fct.
//
// fwk.ComponentMgr will then be able to create new values of that type t
// using the associated FactoryFunc fct.
// If a type t was already registered, the previous FactoryFunc value will be
// silently overridden with the new FactoryFunc value.
func Register(t reflect.Type, fct FactoryFunc) {
	comp := t.PkgPath() + "." + t.Name()
	gFactory[comp] = fct
	//fmt.Printf("### factories ###\n%v\n", gFactory)
}

// Registry returns the list of all registered and known components.
func Registry() []string {
	comps := make([]string, 0, len(gFactory))
	for k, _ := range gFactory {
		comps = append(comps, k)
	}
	sort.Strings(comps)
	return comps
}

// New creates a new Component value with type t and name n.
func (app *appmgr) New(t, n string) (Component, error) {
	var err error
	fct, ok := gFactory[t]
	if !ok {
		return nil, Errorf("no component with type [%s] registered", t)
	}

	if _, dup := app.props[n]; dup {
		return nil, Errorf("component with name [%s] already created", n)
	}
	app.props[n] = make(map[string]interface{})

	c, err := fct(t, n, app)
	if err != nil {
		return nil, Errorf("error creating [%s:%s] %v", t, n, err)
	}
	if c.Name() == "" {
		return nil, Errorf("factory for [%s] does NOT set the name of the component", t)
	}

	err = app.addComponent(c)
	if err != nil {
		return nil, err
	}

	switch c := c.(type) {
	case Svc:
		err = app.AddSvc(c)
		if err != nil {
			return nil, err
		}
	case Task:
		err = app.AddTask(c)
		if err != nil {
			return nil, err
		}
	}

	return c, err
}

// EOF
