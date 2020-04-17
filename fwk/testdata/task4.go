// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

// task4 is like task2, except it works on float64s
type task4 struct {
	fwk.TaskBase

	input  string
	output string
	fct    func(f float64) float64
}

func (tsk *task4) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("configure...\n")

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf(float64(1)))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(float64(1)))
	if err != nil {
		return err
	}

	msg.Infof("configure... [done]\n")
	return err
}

func (tsk *task4) StartTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("start...\n")
	return nil
}

func (tsk *task4) StopTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("stop...\n")
	return nil
}

func (tsk *task4) Process(ctx fwk.Context) error {
	store := ctx.Store()
	msg := ctx.Msg()
	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}
	i := v.(float64)
	o := tsk.fct(i)
	err = store.Put(tsk.output, o)
	if err != nil {
		return err
	}

	msg.Infof("proc... (id=%d|%d) => [%d -> %d]\n", ctx.ID(), ctx.Slot(), i, o)
	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task4{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &task4{
				TaskBase: fwk.NewTask(typ, name, mgr),
				input:    "floats1",
				output:   "massaged_floats1",
			}
			tsk.fct = func(f float64) float64 {
				return f * f
			}

			err = tsk.DeclProp("Input", &tsk.input)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Output", &tsk.output)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Fct", &tsk.fct)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
