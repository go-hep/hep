// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"reflect"
	"sync"

	"go-hep.org/x/hep/fwk"
)

type reducer struct {
	fwk.TaskBase

	input string
	exp   int64
	sum   int64
	nevts int
	mux   sync.RWMutex
}

func (tsk *reducer) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf(int64(0)))
	if err != nil {
		return err
	}

	return err
}

func (tsk *reducer) StartTask(ctx fwk.Context) error {
	var err error
	return err
}

func (tsk *reducer) StopTask(ctx fwk.Context) error {
	var err error

	tsk.mux.RLock()
	sum := tsk.sum
	nevts := tsk.nevts
	tsk.mux.RUnlock()

	msg := ctx.Msg()
	if sum != tsk.exp {
		msg.Errorf("expected sum=%v. got=%v (nevts=%d)\n", tsk.exp, sum, nevts)
		return fwk.Errorf("%s: expected sum=%v. got=%v (nevts=%d)", tsk.Name(), tsk.exp, sum, nevts)
	}
	msg.Debugf("expected sum=%v. got=%v (all GOOD) (nevts=%d)\n", tsk.exp, sum, nevts)

	return err
}

func (tsk *reducer) Process(ctx fwk.Context) error {
	var err error

	tsk.mux.Lock()
	tsk.nevts += 1
	tsk.mux.Unlock()

	store := ctx.Store()
	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	val := v.(int64)
	tsk.mux.Lock()
	tsk.sum += val
	sum := tsk.sum
	tsk.mux.Unlock()

	msg := ctx.Msg()
	msg.Infof("sum=%d (id=%d|%d)\n", sum, ctx.ID(), ctx.Slot())
	return err
}

func newreducer(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &reducer{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "Input",
		sum:      0,
		exp:      0,
		nevts:    0,
	}

	err = tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Sum", &tsk.exp)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(reducer{}), newreducer)
}
