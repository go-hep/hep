// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"reflect"

	"go-hep.org/x/hep/fwk/fsm"
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
// e.g. "go-hep.org/x/hep/fwk/testdata.task1"
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

// FSMState returns the current state of the FSM
func (tsk *TaskBase) FSMState() fsm.State {
	return tsk.mgr.FSMState()
}
