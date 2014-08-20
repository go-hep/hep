package fwk

import (
	"reflect"
)

// TaskBase provides a base implementation for fwk.Task
type TaskBase struct {
	t   string
	n   string
	mgr App
}

// NewTask creates a new TaskBase of type typ and name name,
// managed by the fwk.App mgr.
func NewTask(typ, name string, mgr App) TaskBase {
	return TaskBase{
		t:   typ,
		n:   name,
		mgr: mgr,
	}
}

// Type returns the fully qualified type of the underlying task.
// e.g. "github.com/go-hep/fwk/testdata.task1"
func (tsk *TaskBase) Type() string {
	return tsk.t
}

// Name returns the name of the underlying task.
// e.g. "my-task"
func (tsk *TaskBase) Name() string {
	return tsk.n
}

// DeclInPort declares this task has an input Port with name n and type t.
func (tsk *TaskBase) DeclInPort(n string, t reflect.Type) error {
	return tsk.mgr.DeclInPort(tsk, n, t)
}

// DeclOutPort declares this task has an output Port with name n and type t.
func (tsk *TaskBase) DeclOutPort(n string, t reflect.Type) error {
	return tsk.mgr.DeclOutPort(tsk, n, t)
}

// DeclProp declares this task has a property named n,
// and takes a pointer to the associated value.
func (tsk *TaskBase) DeclProp(n string, ptr interface{}) error {
	return tsk.mgr.DeclProp(tsk, n, ptr)
}

// SetProp sets the property name n with the value v.
func (tsk *TaskBase) SetProp(n string, v interface{}) error {
	return tsk.mgr.SetProp(tsk, n, v)
}

// GetProp returns the value of the property named n.
func (tsk *TaskBase) GetProp(n string) (interface{}, error) {
	return tsk.mgr.GetProp(tsk, n)
}

// SvcBase provides a base implementation for fwk.Svc
type SvcBase struct {
	t   string
	n   string
	mgr App
}

// NewSvc creates a new SvcBase of type typ and name name,
// managed by the fwk.App mgr.
func NewSvc(typ, name string, mgr App) SvcBase {
	return SvcBase{
		t:   typ,
		n:   name,
		mgr: mgr,
	}
}

// Type returns the fully qualified type of the underlying service.
// e.g. "github.com/go-hep/fwk/testdata.svc1"
func (svc *SvcBase) Type() string {
	return svc.t
}

// Name returns the name of the underlying service.
// e.g. "my-service"
func (svc *SvcBase) Name() string {
	return svc.n
}

// DeclProp declares this service has a property named n,
// and takes a pointer to the associated value.
func (svc *SvcBase) DeclProp(n string, ptr interface{}) error {
	return svc.mgr.DeclProp(svc, n, ptr)
}

// SetProp sets the property name n with the value v.
func (svc *SvcBase) SetProp(name string, value interface{}) error {
	return svc.mgr.SetProp(svc, name, value)
}

// GetProp returns the value of the property named n.
func (svc *SvcBase) GetProp(name string) (interface{}, error) {
	return svc.mgr.GetProp(svc, name)
}

// EOF
