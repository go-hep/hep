// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type task2 struct {
	fwk.TaskBase

	fct func(f float64) float64
}

func (tsk *task2) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("configure...\n")

	msg.Infof("configure... [done]\n")
	return err
}

func (tsk *task2) StartTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("start...\n")
	return nil
}

func (tsk *task2) StopTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("stop...\n")
	return nil
}

func (tsk *task2) Process(ctx fwk.Context) error {
	store := ctx.Store()
	msg := ctx.Msg()
	msg.Infof("proc...\n")
	v, err := store.Get("floats1")
	if err != nil {
		return err
	}
	v = tsk.fct(v.(float64))
	err = store.Put("massaged_floats1", v)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task2{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &task2{
				TaskBase: fwk.NewTask(typ, name, mgr),
			}
			tsk.fct = func(f float64) float64 {
				return f * f
			}

			err = tsk.DeclProp("Fct", &tsk.fct)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclInPort("floats1", reflect.TypeOf(float64(1.0)))
			if err != nil {
				return nil, err
			}

			err = tsk.DeclOutPort("massaged_floats1", reflect.TypeOf(float64(1.0)))
			if err != nil {
				return nil, err
			}
			return tsk, err
		},
	)
}
