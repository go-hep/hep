package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type BTagging struct {
    fwk.TaskBase
}

func (tsk *BTagging) Configure(ctx fwk.Context) fwk.Error {
    var err fwk.Error

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

func (tsk *BTagging) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *BTagging) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *BTagging) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func newBTagging(typ, name string, mgr fwk.App) (fwk.Component, fwk.Error) {
	var err fwk.Error

	tsk := &BTagging{
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
	fwk.Register(reflect.TypeOf(BTagging{}), newBTagging)
}
