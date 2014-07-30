package main

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type task1 struct {
	fwk.TaskBase

	f1 float64
	f2 float64
}

func (tsk *task1) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	msg := ctx.Msg()

	msg.Infof("configure ...\n")

	msg.Infof("configure ... [done]\n")
	return err
}

func (tsk *task1) StartTask(ctx fwk.Context) fwk.Error {
	msg := ctx.Msg()
	msg.Infof("start...\n")
	return nil
}

func (tsk *task1) StopTask(ctx fwk.Context) fwk.Error {
	msg := ctx.Msg()
	msg.Infof("stop...\n")
	return nil
}

func (tsk *task1) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	msg := ctx.Msg()
	msg.Infof("proc...\n")
	store := ctx.Store()

	err = store.Put("floats1", tsk.f1)
	if err != nil {
		return err
	}

	err = store.Put("floats2", tsk.f2)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task1{}),
		func(name string, mgr fwk.App) (fwk.Component, fwk.Error) {
			var err fwk.Error
			tsk := &task1{
				TaskBase: fwk.NewTask(name, mgr),
				f1:       -1,
				f2:       +2,
			}

			err = tsk.DeclProp("Float1", &tsk.f1)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Float2", &tsk.f2)
			if err != nil {
				return nil, err
			}

			err = tsk.SetProp("Float1", 1.)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclOutPort("floats1", reflect.TypeOf(float64(1.0)))
			if err != nil {
				return nil, err
			}

			err = tsk.DeclOutPort("floats2", reflect.TypeOf(float64(1.0)))
			if err != nil {
				return nil, err
			}
			return tsk, err
		},
	)
}
