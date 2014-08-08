package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type Isolation struct {
	fwk.TaskBase
}

func (tsk *Isolation) Configure(ctx fwk.Context) error {
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

func (tsk *Isolation) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Isolation) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Isolation) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newIsolation(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &Isolation{
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
	fwk.Register(reflect.TypeOf(Isolation{}), newIsolation)
}
