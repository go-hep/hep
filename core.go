package fwk

import (
	"fmt"
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

// DeclPorter is the interface to declare input/output ports for the data flow.
type DeclPorter interface {
	DeclInPort(name string) Error
	DeclOutPort(name string) Error
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

// EOF
