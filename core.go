package fwk

import (
	"fmt"
	"reflect"
)

type Error interface {
	error
}

type statuscode int

func (sc statuscode) Error() string {
	return fmt.Sprintf("fwk: error code [%d]", int(sc))
}

type Context interface {
	Id() int64
	Slot() int
	Store() Store
}

type Component interface {
	CompName() string // Name of the component (ex: "MyPropagator")
	CompType() string // Type of the component (ex: "github.com/foo/bar.Propagator")
}

type Task interface {
	Component

	StartTask(ctx Context) Error
	Process(ctx Context) Error
	StopTask(ctx Context) Error
}

type TaskMgr interface {
	AddTask(tsk Task) Error
	DelTask(tsk Task) Error
	HasTask(n string) bool
	GetTask(n string) Task
	Tasks() []Task
}

type Configurer interface {
	Component
	Configure(ctx Context) Error
}

type Svc interface {
	Component

	StartSvc(ctx Context) Error
	StopSvc(ctx Context) Error
}

type SvcMgr interface {
	AddSvc(svc Svc) Error
	DelSvc(svc Svc) Error
	HasSvc(n string) bool
	GetSvc(n string) Svc
	Svcs() []Svc
}

type App interface {
	Component
	Run() Error
}

type Property interface {
	DeclProp(name string, defvalue interface{})
	SetProp(name string, value interface{})
	GetProp(name string) interface{}
}

type Store interface {
	Get(key string) (interface{}, Error)
	Put(key string, value interface{}) Error
}

type Level int

const (
	LvlVerbose Level = -20
	LvlDebug   Level = -10
	LvlInfo    Level = 0
	LvlWarning Level = 10
	LvlEror    Level = 20
)

type MsgStream interface {
	MsgVerbose(format string, a ...interface{}) (int, error)
	MsgDebug(format string, a ...interface{}) (int, error)
	MsgInfo(format string, a ...interface{}) (int, error)
	MsgWarning(format string, a ...interface{}) (int, error)
	MsgError(format string, a ...interface{}) (int, error)

	Msg(lvl Level, format string, a ...interface{}) (int, error)
}

func DeclProp(c Component, name string, ptr interface{}) Error {
	app := g_app.(*appmgr)
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

func GetProp(c Component, name string) (interface{}, Error) {
	app := g_app.(*appmgr)
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

func SetProp(c Component, name string, value interface{}) Error {
	app := g_app.(*appmgr)
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

func DeclInPort(tsk Task, name string, value interface{}) Error {
	app := g_app.(*appmgr)
	return app.dflow.addInNode(tsk, name, value)
}

func DeclOutPort(tsk Task, name string, value interface{}) Error {
	app := g_app.(*appmgr)
	return app.dflow.addOutNode(tsk, name, value)
}

/*
func DeclInOutPort(comp Component, name string, value interface{}) Error {
	return nil
}
*/

// EOF
