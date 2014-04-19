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

func (tsk TaskBase) DeclInPort(name string, t reflect.Type) Error {
	return g_mgr.dflow.addInNode(tsk.name, name, t)
}

func (tsk TaskBase) DeclOutPort(name string, t reflect.Type) Error {
	return g_mgr.dflow.addOutNode(tsk.name, name, t)
}

func (tsk TaskBase) DeclProp(name string, ptr interface{}) Error {
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return Errorf("fwk.DeclProp: no Task [%s] known to TaskMgr", tsk.name)
	}

	return g_mgr.DeclProp(c, name, ptr)
}

func (tsk TaskBase) SetProp(name string, value interface{}) Error {
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return Errorf("fwk.SetProp: no Task [%s] known to TaskMgr", tsk.name)
	}

	return g_mgr.SetProp(c, name, value)
}

func (tsk TaskBase) GetProp(name string) (interface{}, Error) {
	c := g_mgr.GetTask(tsk.name)
	if c == nil {
		return nil, Errorf("fwk.GetProp: no Task [%s] known to TaskMgr", tsk.name)
	}

	return g_mgr.GetProp(c, name)
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
