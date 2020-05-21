// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type dumper struct {
	fwk.TaskBase

	input string
}

func (tsk *dumper) Configure(ctx fwk.Context) error {
	err := tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}
	return nil
}

func (tsk *dumper) StartTask(ctx fwk.Context) error {
	return nil
}

func (tsk *dumper) StopTask(ctx fwk.Context) error {
	return nil
}

func (tsk *dumper) Process(ctx fwk.Context) error {
	store := ctx.Store()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	//msg.Debugf(">>> particles: %v\n", len(parts))
	fmt.Printf("%s: %d\n", tsk.input, len(input))
	return nil
}

func newDumper(typ, name string, mgr fwk.App) (fwk.Component, error) {
	tsk := &dumper{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "Input",
	}

	err := tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	return tsk, nil
}

func init() {
	fwk.Register(reflect.TypeOf(dumper{}), newDumper)
}
