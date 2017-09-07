// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type task1 struct {
	fwk.TaskBase

	i1prop string
	i2prop string

	i1 int64
	i2 int64
}

func (tsk *task1) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()

	msg.Infof("configure ...\n")

	err = tsk.DeclOutPort(tsk.i1prop, reflect.TypeOf(int64(1.0)))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.i2prop, reflect.TypeOf(int64(1.0)))
	if err != nil {
		return err
	}

	msg.Infof("configure ... [done]\n")
	return err
}

func (tsk *task1) StartTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("start...\n")
	return nil
}

func (tsk *task1) StopTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("stop...\n")
	return nil
}

func (tsk *task1) Process(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("proc... (id=%d|%d) => [%d, %d]\n", ctx.ID(), ctx.Slot(), tsk.i1, tsk.i2)
	store := ctx.Store()

	err = store.Put(tsk.i1prop, tsk.i1)
	if err != nil {
		return err
	}

	err = store.Put(tsk.i2prop, tsk.i2)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task1{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &task1{
				TaskBase: fwk.NewTask(typ, name, mgr),
				i1prop:   "ints1",
				i2prop:   "ints2",
				i1:       -1,
				i2:       +2,
			}

			err = tsk.DeclProp("Ints1", &tsk.i1prop)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Ints2", &tsk.i2prop)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Int1", &tsk.i1)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Int2", &tsk.i2)
			if err != nil {
				return nil, err
			}

			err = tsk.SetProp("Int1", int64(1))
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
