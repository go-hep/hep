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
	Id() int64    // id of this context (e.g. entry number or some kind of event number)
	Slot() int    // slot number in the pool of event sequences
	Store() Store // data store corresponding to the id+slot
}

type Component interface {
	Name() string     // Name of the component (ex: "MyPropagator")
	SetName(n string) // rename a component (should only be done by the framework!)
}

type ComponentMgr interface {
	Component(n string) Component
	HasComponent(n string) bool
	Components() []Component
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
	ComponentMgr
	SvcMgr
	TaskMgr
	PropMgr

	Run() Error
}

type PropMgr interface {
	DeclProp(c Component, name string, ptr interface{}) Error
	SetProp(c Component, name string, value interface{}) Error
	GetProp(c Component, name string) (interface{}, Error)
}

type Property interface {
	DeclProp(name string, ptr interface{}) Error
	SetProp(name string, value interface{}) Error
	GetProp(name string) (interface{}, Error)
}

type Store interface {
	Get(key string) (interface{}, Error)
	Put(key string, value interface{}) Error
	Has(key string) bool
}

// DeclPorter is the interface to declare input/output ports for the data flow.
type DeclPorter interface {
	DeclInPort(name string, t reflect.Type) Error
	DeclOutPort(name string, t reflect.Type) Error
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
	Verbosef(format string, a ...interface{}) (int, error)
	Debugf(format string, a ...interface{}) (int, error)
	Infof(format string, a ...interface{}) (int, error)
	Warnf(format string, a ...interface{}) (int, error)
	Errorf(format string, a ...interface{}) (int, error)

	Msg(lvl Level, format string, a ...interface{}) (int, error)
}

// EOF
