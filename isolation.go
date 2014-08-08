package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type isolation struct {
	fwk.TaskBase
}

func (tsk *isolation) Configure(ctx fwk.Context) error {
	var err error

	// err = tsk.DeclInPort(tsk.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	return err
}

func (tsk *isolation) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *isolation) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *isolation) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newIsolation(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &isolation{
		TaskBase: fwk.NewTask(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
	}

	// err = tsk.DeclProp("Input", &tsk.input)
	// if err != nil {
	// 	return nil, err
	// }

	// err = tsk.DeclProp("Output", &tsk.output)
	// if err != nil {
	//	return nil, err
	// }

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(isolation{}), newIsolation)
}
