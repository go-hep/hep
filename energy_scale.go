package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type EnergyScale struct {
	fwk.TaskBase
}

func (tsk *EnergyScale) Configure(ctx fwk.Context) error {
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

func (tsk *EnergyScale) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *EnergyScale) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *EnergyScale) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newEnergyScale(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &EnergyScale{
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
	fwk.Register(reflect.TypeOf(EnergyScale{}), newEnergyScale)
}
