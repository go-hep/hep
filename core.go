package fwk

import (
	"fmt"
	"reflect"
)

type statuscode int

func (sc statuscode) Error() string {
	return fmt.Sprintf("fwk: error code [%d]", int(sc))
}

type Context interface {
	Id() int64      // id of this context (e.g. entry number or some kind of event number)
	Slot() int      // slot number in the pool of event sequences
	Store() Store   // data store corresponding to the id+slot
	Msg() MsgStream // messaging for this context (id+slot)
}

type Component interface {
	Type() string // Type of the component (ex: "github.com/go-hep/fads.MomentumSmearing")
	Name() string // Name of the component (ex: "MyPropagator")
}

type ComponentMgr interface {
	Component(n string) Component
	HasComponent(n string) bool
	Components() []Component
	New(t, n string) (Component, error)
}

type Task interface {
	Component

	StartTask(ctx Context) error
	Process(ctx Context) error
	StopTask(ctx Context) error
}

type TaskMgr interface {
	AddTask(tsk Task) error
	DelTask(tsk Task) error
	HasTask(n string) bool
	GetTask(n string) Task
	Tasks() []Task
}

type Configurer interface {
	Component
	Configure(ctx Context) error
}

type Svc interface {
	Component

	StartSvc(ctx Context) error
	StopSvc(ctx Context) error
}

type SvcMgr interface {
	AddSvc(svc Svc) error
	DelSvc(svc Svc) error
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
	PortMgr

	Run() error

	Msg() MsgStream
}

type PropMgr interface {
	DeclProp(c Component, name string, ptr interface{}) error
	SetProp(c Component, name string, value interface{}) error
	GetProp(c Component, name string) (interface{}, error)
	HasProp(c Component, name string) bool
}

type Property interface {
	DeclProp(name string, ptr interface{}) error
	SetProp(name string, value interface{}) error
	GetProp(name string) (interface{}, error)
}

type Store interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Has(key string) bool
}

// DeclPorter is the interface to declare input/output ports for the data flow.
type DeclPorter interface {
	DeclInPort(name string, t reflect.Type) error
	DeclOutPort(name string, t reflect.Type) error
}

// PortMgr is the interface to manage input/output ports for the data flow
type PortMgr interface {
	DeclInPort(c Component, name string, t reflect.Type) error
	DeclOutPort(c Component, name string, t reflect.Type) error
}

type Level int

const (
	//LvlVerbose Level = -20
	LvlDebug   Level = -10
	LvlInfo    Level = 0
	LvlWarning Level = 10
	LvlError   Level = 20
)

func (lvl Level) msgstring() string {
	switch lvl {
	case LvlDebug:
		return "DBG "
	case LvlInfo:
		return "INFO"
	case LvlWarning:
		return "WARN"
	case LvlError:
		return "ERR "
	}
	panic(Errorf("fwk.Level: invalid fwk.Level value [%d]", int(lvl)))
}

func (lvl Level) String() string {
	switch lvl {
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO"
	case LvlWarning:
		return "WARN"
	case LvlError:
		return "ERROR"
	}
	panic(Errorf("fwk.Level: invalid fwk.Level value [%d]", int(lvl)))
}

type MsgStream interface {
	Debugf(format string, a ...interface{}) (int, error)
	Infof(format string, a ...interface{}) (int, error)
	Warnf(format string, a ...interface{}) (int, error)
	Errorf(format string, a ...interface{}) (int, error)

	Msg(lvl Level, format string, a ...interface{}) (int, error)
}

// Deleter prepares values to be GC-reclaimed
type Deleter interface {
	Delete() error
}

// EOF
