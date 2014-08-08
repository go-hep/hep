package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type tautagging struct {
	fwk.TaskBase
}

func (tsk *tautagging) Configure(ctx fwk.Context) error {
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

func (tsk *tautagging) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *tautagging) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *tautagging) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newTauTagging(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &tautagging{
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
	fwk.Register(reflect.TypeOf(tautagging{}), newTauTagging)
}
