package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type FastJetFinder struct {
	fwk.TaskBase
}

func (tsk *FastJetFinder) Configure(ctx fwk.Context) error {
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

func (tsk *FastJetFinder) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *FastJetFinder) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *FastJetFinder) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newFastJetFinder(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &FastJetFinder{
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
	fwk.Register(reflect.TypeOf(FastJetFinder{}), newFastJetFinder)
}
