package main

import (
	"fmt"

	"github.com/go-hep/fwk"
)

type task1 struct {
	fwk.TaskBase

	f1 float64
	f2 float64
}

func (tsk *task1) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	fmt.Printf(">>> configure [%v]...\n", tsk.CompName())

	tsk.f1 = -1
	tsk.f2 = 2

	err = tsk.DeclProp("Float1", &tsk.f1)
	if err != nil {
		return err
	}

	err = tsk.DeclProp("Float2", &tsk.f2)
	if err != nil {
		return err
	}

	err = tsk.SetProp("Float1", 1.)
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort("floats1")
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort("floats2")
	if err != nil {
		return err
	}

	fmt.Printf(">>> configure [%v]... [done]\n", tsk.CompName())
	return err
}

func (tsk *task1) StartTask(ctx fwk.Context) fwk.Error {
	fmt.Printf(">>> start [%v]...\n", tsk.CompName())
	return nil
}

func (tsk *task1) StopTask(ctx fwk.Context) fwk.Error {
	fmt.Printf(">>> stop [%v]...\n", tsk.CompName())
	return nil
}

func (tsk *task1) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	fmt.Printf(">>> proc [%v]...\n", tsk.CompName())
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
