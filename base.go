package fwk

import (
	"reflect"
)

type TaskBase struct {
	t   string
	n   string
	mgr App
}

func NewTask(typ, name string, mgr App) TaskBase {
	return TaskBase{
		t:   typ,
		n:   name,
		mgr: mgr,
	}
}

func (tsk *TaskBase) Type() string {
	return tsk.t
}

func (tsk *TaskBase) Name() string {
	return tsk.n
}

func (tsk *TaskBase) DeclInPort(name string, t reflect.Type) error {
	return tsk.mgr.DeclInPort(tsk, name, t)
}

func (tsk *TaskBase) DeclOutPort(name string, t reflect.Type) error {
	return tsk.mgr.DeclOutPort(tsk, name, t)
}

func (tsk *TaskBase) DeclProp(name string, ptr interface{}) error {
	return tsk.mgr.DeclProp(tsk, name, ptr)
}

func (tsk *TaskBase) SetProp(name string, value interface{}) error {
	return tsk.mgr.SetProp(tsk, name, value)
}

func (tsk *TaskBase) GetProp(name string) (interface{}, error) {
	return tsk.mgr.GetProp(tsk, name)
}

type SvcBase struct {
	t   string
	n   string
	mgr App
}

func NewSvc(typ, name string, mgr App) SvcBase {
	return SvcBase{
		t:   typ,
		n:   name,
		mgr: mgr,
	}
}

func (svc *SvcBase) Type() string {
	return svc.t
}

func (svc *SvcBase) Name() string {
	return svc.n
}

func (svc *SvcBase) DeclProp(name string, ptr interface{}) error {
	return svc.mgr.DeclProp(svc, name, ptr)
}

func (svc *SvcBase) SetProp(name string, value interface{}) error {
	return svc.mgr.SetProp(svc, name, value)
}

func (svc *SvcBase) GetProp(name string) (interface{}, error) {
	return svc.mgr.GetProp(svc, name)
}

// EOF
