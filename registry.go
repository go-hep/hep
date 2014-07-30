package fwk

import (
	"reflect"
	"sort"
)

type factoryDb map[string]func(name string, mgr App) (Component, Error)

var g_factory factoryDb = make(factoryDb)

func Register(t reflect.Type, fct func(name string, mgr App) (Component, Error)) {
	comp := t.PkgPath() + "." + t.Name()
	g_factory[comp] = fct
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

func (app *appmgr) New(t, n string) (Component, Error) {
	var err Error
	fct, ok := g_factory[t]
	if !ok {
		return nil, Errorf("no component with type [%s] registered", t)
	}

	if _, dup := app.props[n]; dup {
		return nil, Errorf("component with name [%s] already created", n)
	}
	app.props[n] = make(map[string]interface{})

	c, err := fct(n, app)
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
