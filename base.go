package fwk

import (
	"reflect"
)

type Base struct {
	name string
}

func (c Base) Name() string {
	return c.name
}

func (c *Base) SetName(n string) {
	c.name = n
}

type TaskBase struct {
	name string
}

func (tsk TaskBase) Name() string {
	return tsk.name
}

func (tsk *TaskBase) SetName(n string) {
	tsk.name = n
}

func (tsk TaskBase) DeclInPort(name string) Error {
	return g_mgr.dflow.addInNode(tsk.name, name)
}

func (tsk TaskBase) DeclOutPort(name string) Error {
	return g_mgr.dflow.addOutNode(tsk.name, name)
}

func (tsk TaskBase) DeclProp(name string, ptr interface{}) Error {
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return Errorf("fwk.DeclProp: no Task [%s] known to TaskMgr", name)
	}

	_, ok := g_mgr.props[c]
	if !ok {
		g_mgr.props[c] = make(map[string]interface{})
	}
	switch reflect.TypeOf(ptr).Kind() {
	case reflect.Ptr:
		// ok
	default:
		return Errorf(
			"fwk.DeclProp: component [%s] didn't pass a pointer for the property [%s] (type=%T)",
			c.Name(),
			name,
			ptr,
		)
	}
	g_mgr.props[c][name] = ptr
	return nil
}

func (tsk TaskBase) SetProp(name string, value interface{}) Error {
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return Errorf("fwk.SetProp: no Task [%s] known to TaskMgr", name)
	}

	m, ok := g_mgr.props[c]
	if !ok {
		return Errorf(
			"fwk.SetProp: component [%s] didn't declare any property",
			c.Name(),
		)
	}
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	ptr := reflect.ValueOf(m[name])
	dst := ptr.Elem().Type()
	if !rt.AssignableTo(dst) {
		return Errorf(
			"fwk.SetProp: component [%s] has property [%s] with type [%s]. got value=%v (type=%s)",
			c.Name(),
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
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return nil, Errorf("fwk.GetProp: no Task [%s] known to TaskMgr", name)
	}

	m, ok := g_mgr.props[c]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property",
			c.Name(),
		)
	}

	ptr, ok := m[name]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property with name [%s]",
			c.Name(),
			name,
		)
	}

	v := reflect.Indirect(reflect.ValueOf(ptr)).Interface()
	return v, nil
}

type SvcBase struct {
	name string
}

func (svc SvcBase) Name() string {
	return svc.name
}

func (svc *SvcBase) SetName(n string) {
	svc.name = n
}

func (svc SvcBase) DeclProp(name string, ptr interface{}) Error {
	c := g_mgr.GetSvc(svc.name)
	if c == nil {
		return Errorf("fwk.DeclProp: no Svc [%s] known to SvcMgr", name)
	}

	return g_mgr.DeclProp(c, name, ptr)
}

func (svc SvcBase) SetProp(name string, value interface{}) Error {
	c := g_mgr.GetSvc(svc.name)
	if c == nil {
		return Errorf("fwk.SetProp: no Svc [%s] known to SvcMgr", name)
	}

	return g_mgr.SetProp(c, name, value)
}

func (svc SvcBase) GetProp(name string) (interface{}, Error) {
	c := g_mgr.GetSvc(svc.name)
	if c == nil {
		return nil, Errorf("fwk.GetProp: no Svc [%s] known to SvcMgr", name)
	}

	return g_mgr.GetProp(c, name)
}

// EOF
