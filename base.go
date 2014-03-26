package fwk

import (
	"reflect"
)

type Base struct {
	Type string
	Name string
}

func (c Base) CompName() string {
	return c.Name
}

func (c Base) CompType() string {
	return c.Type
}

type TaskBase struct {
	Type string
	Name string
}

func (tsk TaskBase) CompName() string {
	return tsk.Name
}

func (tsk TaskBase) CompType() string {
	return tsk.Type
}

func (tsk TaskBase) DeclInPort(name string) Error {
	app := g_app.(*appmgr)
	return app.dflow.addInNode(tsk.Name, name)
}

func (tsk TaskBase) DeclOutPort(name string) Error {
	app := g_app.(*appmgr)
	return app.dflow.addOutNode(tsk.Name, name)
}

func (tsk TaskBase) DeclProp(name string, ptr interface{}) Error {
	app := g_app.(*appmgr)
	c := app.GetTask(tsk.Name)
	if c == nil {
		return Errorf("fwk.DeclProp: no Task [%s] known to TaskMgr", name)
	}

	_, ok := app.props[c]
	if !ok {
		app.props[c] = make(map[string]interface{})
	}
	switch reflect.TypeOf(ptr).Kind() {
	case reflect.Ptr:
		// ok
	default:
		return Errorf(
			"fwk.DeclProp: component [%s] didn't pass a pointer for the property [%s] (type=%T)",
			c.CompName(),
			name,
			ptr,
		)
	}
	app.props[c][name] = ptr
	return nil
}

func (tsk TaskBase) SetProp(name string, value interface{}) Error {
	app := g_app.(*appmgr)
	c := app.GetTask(tsk.Name)
	if c == nil {
		return Errorf("fwk.SetProp: no Task [%s] known to TaskMgr", name)
	}

	m, ok := app.props[c]
	if !ok {
		return Errorf(
			"fwk.SetProp: component [%s] didn't declare any property",
			c.CompName(),
		)
	}
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	ptr := reflect.ValueOf(m[name])
	dst := ptr.Elem().Type()
	if !rt.AssignableTo(dst) {
		return Errorf(
			"fwk.SetProp: component [%s:%s] has property [%s] with type [%s]. got value=%v (type=%s)",
			c.CompType(),
			c.CompName(),
			name,
			dst.Name(),
			value,
			rt.Name(),
		)
	}
	ptr.Elem().Set(rv)
	return nil

}

func (tsk TaskBase) GetProp(name string) (interface{}, Error) {
	app := g_app.(*appmgr)
	c := app.GetTask(tsk.Name)
	if c == nil {
		return nil, Errorf("fwk.GetProp: no Task [%s] known to TaskMgr", name)
	}

	m, ok := app.props[c]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property",
			c.CompName(),
		)
	}

	ptr, ok := m[name]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property with name [%s]",
			c.CompName(),
			name,
		)
	}

	v := reflect.Indirect(reflect.ValueOf(ptr)).Interface()
	return v, nil
}

type SvcBase struct {
	Type string
	Name string
}

func (svc SvcBase) CompName() string {
	return svc.Name
}

func (svc SvcBase) CompType() string {
	return svc.Type
}

func (svc SvcBase) DeclProp(name string, ptr interface{}) Error {
	app := g_app.(*appmgr)
	c := app.GetSvc(svc.Name)
	if c == nil {
		return Errorf("fwk.DeclProp: no Svc [%s] known to SvcMgr", name)
	}

	_, ok := app.props[c]
	if !ok {
		app.props[c] = make(map[string]interface{})
	}
	switch reflect.TypeOf(ptr).Kind() {
	case reflect.Ptr:
		// ok
	default:
		return Errorf(
			"fwk.DeclProp: component [%s] didn't pass a pointer for the property [%s] (type=%T)",
			c.CompName(),
			name,
			ptr,
		)
	}
	app.props[c][name] = ptr
	return nil
}

func (svc SvcBase) SetProp(name string, value interface{}) Error {
	app := g_app.(*appmgr)
	c := app.GetSvc(svc.Name)
	if c == nil {
		return Errorf("fwk.SetProp: no Svc [%s] known to SvcMgr", name)
	}

	m, ok := app.props[c]
	if !ok {
		return Errorf(
			"fwk.SetProp: component [%s] didn't declare any property",
			c.CompName(),
		)
	}
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	ptr := reflect.ValueOf(m[name])
	dst := ptr.Elem().Type()
	if !rt.AssignableTo(dst) {
		return Errorf(
			"fwk.SetProp: component [%s:%s] has property [%s] with type [%s]. got value=%v (type=%s)",
			c.CompType(),
			c.CompName(),
			name,
			dst.Name(),
			value,
			rt.Name(),
		)
	}
	ptr.Elem().Set(rv)
	return nil

}

func (svc SvcBase) GetProp(name string) (interface{}, Error) {
	app := g_app.(*appmgr)
	c := app.GetSvc(svc.Name)
	if c == nil {
		return nil, Errorf("fwk.GetProp: no Svc [%s] known to SvcMgr", name)
	}

	m, ok := app.props[c]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property",
			c.CompName(),
		)
	}

	ptr, ok := m[name]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property with name [%s]",
			c.CompName(),
			name,
		)
	}

	v := reflect.Indirect(reflect.ValueOf(ptr)).Interface()
	return v, nil
}

// EOF
