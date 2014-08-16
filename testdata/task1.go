package testdata

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type task1 struct {
	fwk.TaskBase

	f1prop string
	f2prop string

	f1 float64
	f2 float64
}

func (tsk *task1) Configure(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()

	msg.Infof("configure ...\n")

	err = tsk.DeclOutPort(tsk.f1prop, reflect.TypeOf(float64(1.0)))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.f2prop, reflect.TypeOf(float64(1.0)))
	if err != nil {
		return err
	}

	msg.Infof("configure ... [done]\n")
	return err
}

func (tsk *task1) StartTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("start...\n")
	return nil
}

func (tsk *task1) StopTask(ctx fwk.Context) error {
	msg := ctx.Msg()
	msg.Infof("stop...\n")
	return nil
}

func (tsk *task1) Process(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("proc...\n")
	store := ctx.Store()

	err = store.Put(tsk.f1prop, tsk.f1)
	if err != nil {
		return err
	}

	err = store.Put(tsk.f2prop, tsk.f2)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task1{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &task1{
				TaskBase: fwk.NewTask(typ, name, mgr),
				f1prop:   "floats1",
				f2prop:   "floats2",
				f1:       -1,
				f2:       +2,
			}

			err = tsk.DeclProp("Floats1", &tsk.f1prop)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Floats2", &tsk.f2prop)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Float1", &tsk.f1)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Float2", &tsk.f2)
			if err != nil {
				return nil, err
			}

			err = tsk.SetProp("Float1", 1.0)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
