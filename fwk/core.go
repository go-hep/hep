// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/fwk/fsm"
)

// Context is the interface to access context-local data.
type Context interface {
	ID() int64      // id of this context (e.g. entry number or some kind of event number)
	Slot() int      // slot number in the pool of event sequences
	Store() Store   // data store corresponding to the id+slot
	Msg() MsgStream // messaging for this context (id+slot)

	Svc(n string) (Svc, error) // retrieve an already existing Svc by name
}

// Component is the interface satisfied by all values in fwk.
//
// A component can be asked for:
// its Type() (ex: "go-hep.org/x/hep/fads.MomentumSmearing")
// its Name() (ex: "MyPropagator")
type Component interface {
	Type() string // Type of the component (ex: "go-hep.org/x/hep/fads.MomentumSmearing")
	Name() string // Name of the component (ex: "MyPropagator")
}

// ComponentMgr manages components.
// ComponentMgr creates and provides access to all the components in a fwk App.
type ComponentMgr interface {
	Component(n string) Component
	HasComponent(n string) bool
	Components() []Component
	New(t, n string) (Component, error)
}

// Task is a component processing event-level data.
// Task.Process is called for every component and for every input event.
type Task interface {
	Component

	StartTask(ctx Context) error
	Process(ctx Context) error
	StopTask(ctx Context) error
}

// TaskMgr manages tasks.
type TaskMgr interface {
	AddTask(tsk Task) error
	DelTask(tsk Task) error
	HasTask(n string) bool
	GetTask(n string) Task
	Tasks() []Task
}

// Configurer are components which can be configured via properties
// declared or created by the job-options.
type Configurer interface {
	Component
	Configure(ctx Context) error
}

// Svc is a component providing services or helper features.
// Services are started before the main event loop processing and
// stopped just after.
type Svc interface {
	Component

	StartSvc(ctx Context) error
	StopSvc(ctx Context) error
}

// SvcMgr manages services.
type SvcMgr interface {
	AddSvc(svc Svc) error
	DelSvc(svc Svc) error
	HasSvc(n string) bool
	GetSvc(n string) Svc
	Svcs() []Svc
}

// App is the component orchestrating all the other components
// in a coherent application to process physics events.
type App interface {
	Component
	ComponentMgr
	SvcMgr
	TaskMgr
	PropMgr
	PortMgr

	FSMStater

	Runner
	Scripter() Scripter

	Msg() MsgStream
}

// Runner runs a fwk App in a batch fashion:
//   - Configure
//   - Start
//   - Run event loop
//   - Stop
//   - Shutdown
type Runner interface {
	Run() error
}

// Scripter gives finer control to running a fwk App
type Scripter interface {
	Configure() error
	Start() error
	Run(evtmax int64) error
	Stop() error
	Shutdown() error
}

// PropMgr manages properties attached to components.
type PropMgr interface {
	DeclProp(c Component, name string, ptr interface{}) error
	SetProp(c Component, name string, value interface{}) error
	GetProp(c Component, name string) (interface{}, error)
	HasProp(c Component, name string) bool
}

// Property is a pair key/value, associated to a component.
// Properties of a given component can be modified
// by a job-option or by other components.
type Property interface {
	DeclProp(name string, ptr interface{}) error
	SetProp(name string, value interface{}) error
	GetProp(name string) (interface{}, error)
}

// Store provides access to a concurrent-safe map[string]interface{} store.
type Store interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
	Has(key string) bool
}

// Port holds the name and type of a data item in a store
type Port struct {
	Name string
	Type reflect.Type
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

// FSMStater is the interface used to query the current state of the fwk application
type FSMStater interface {
	FSMState() fsm.State
}

// Level regulates the verbosity level of a component.
type Level int

// Default verbosity levels.
const (
	LvlDebug   Level = -10 // LvlDebug defines the DBG verbosity level
	LvlInfo    Level = 0   // LvlInfo defines the INFO verbosity level
	LvlWarning Level = 10  // LvlWarning defines the WARN verbosity level
	LvlError   Level = 20  // LvlError defines the ERR verbosity level
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
	panic(fmt.Errorf("fwk.Level: invalid fwk.Level value [%d]", int(lvl)))
}

// String prints the human-readable representation of a Level value.
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
	panic(fmt.Errorf("fwk.Level: invalid fwk.Level value [%d]", int(lvl)))
}

// MsgStream provides access to verbosity-defined formated messages, a la fmt.Printf.
type MsgStream interface {
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})

	Msg(lvl Level, format string, a ...interface{})
}

// Deleter prepares values to be GC-reclaimed
type Deleter interface {
	Delete() error
}
