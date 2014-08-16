package testdata

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type task2 struct {
	fwk.TaskBase

	input  string
	output string
	fct    func(f float64) float64
}

func (tsk *task2) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("configure...\n")

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf(float64(1.0)))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(float64(1.0)))
	if err != nil {
		return err
	}

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
	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}
	v = tsk.fct(v.(float64))
	err = store.Put(tsk.output, v)
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
