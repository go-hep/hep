package testdata

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type task2 struct {
	fwk.TaskBase

	input  string
	output string
	fct    func(f int64) int64
}

func (tsk *task2) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("configure...\n")

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf(int64(1)))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(int64(1)))
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
	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}
	i := v.(int64)
	o := tsk.fct(i)
	err = store.Put(tsk.output, o)
	if err != nil {
		return err
	}

	msg.Infof("proc... (id=%d|%d) => [%d -> %d]\n", ctx.ID(), ctx.Slot(), i, o)
	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task2{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &task2{
				TaskBase: fwk.NewTask(typ, name, mgr),
				input:    "ints1",
				output:   "massaged_ints1",
			}
			tsk.fct = func(f int64) int64 {
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
