package fwk

import (
	"reflect"
)

type TaskBase struct {
	name string
	mgr  App
}

func NewTask(name string, mgr App) TaskBase {
	return TaskBase{
		name: name,
		mgr:  mgr,
	}
}

func (tsk TaskBase) Name() string {
	return tsk.name
}

func (tsk *TaskBase) SetName(n string) {
	tsk.name = n
}

func (tsk *TaskBase) DeclInPort(name string, t reflect.Type) Error {
	return tsk.mgr.DeclInPort(tsk, name, t)
}

func (tsk *TaskBase) DeclOutPort(name string, t reflect.Type) Error {
	return tsk.mgr.DeclOutPort(tsk, name, t)
}

func (tsk *TaskBase) DeclProp(name string, ptr interface{}) Error {
	return tsk.mgr.DeclProp(tsk, name, ptr)
}

func (tsk *TaskBase) SetProp(name string, value interface{}) Error {
	return tsk.mgr.SetProp(tsk, name, value)
}

func (tsk *TaskBase) GetProp(name string) (interface{}, Error) {
	return tsk.mgr.GetProp(tsk, name)
}

type SvcBase struct {
	name string
	mgr  App
}

func NewSvc(name string, mgr App) SvcBase {
	return SvcBase{
		name: name,
		mgr:  mgr,
	}
}

func (svc SvcBase) Name() string {
	return svc.name
}

func (svc *SvcBase) SetName(n string) {
	svc.name = n
}

func (svc *SvcBase) DeclProp(name string, ptr interface{}) Error {
	return svc.mgr.DeclProp(svc, name, ptr)
}

func (svc *SvcBase) SetProp(name string, value interface{}) Error {
	return svc.mgr.SetProp(svc, name, value)
}

func (svc *SvcBase) GetProp(name string) (interface{}, Error) {
	return svc.mgr.GetProp(svc, name)
}

// EOF
