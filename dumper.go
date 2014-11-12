package fads

import (
	"fmt"
	"reflect"

	"github.com/go-hep/fwk"
)

type dumper struct {
	fwk.TaskBase

	input string
}

func (tsk *dumper) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *dumper) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *dumper) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *dumper) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	//msg.Debugf(">>> particles: %v\n", len(parts))
	fmt.Printf("%s: %d\n", tsk.input, len(input))
	return err
}

func newDumper(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &dumper{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "Input",
	}

	err = tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(dumper{}), newDumper)
}
