package main

import (
	"fmt"
	"reflect"

	"github.com/go-hep/fads"
	"github.com/go-hep/fwk"
)

type task3 struct {
	fwk.TaskBase

	parts string
}

func (tsk *task3) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	fmt.Printf(">>> configure [%v]...\n", tsk.Name())

	tsk.parts = "/fads/StableParticles"
	err = tsk.DeclProp("Output", &tsk.parts)
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.parts, reflect.TypeOf([]fads.Candidate{}))
	if err != nil {
		return err
	}

	fmt.Printf(">>> configure [%v]... [done]\n", tsk.Name())
	return err
}

func (tsk *task3) StartTask(ctx fwk.Context) fwk.Error {
	fmt.Printf(">>> start [%v]...\n", tsk.Name())
	return nil
}

func (tsk *task3) StopTask(ctx fwk.Context) fwk.Error {
	fmt.Printf(">>> stop [%v]...\n", tsk.Name())
	return nil
}

func (tsk *task3) Process(ctx fwk.Context) fwk.Error {
	fmt.Printf(">>> proc [%v]...\n", tsk.Name())
	store := ctx.Store()

	parts := make([]fads.Candidate, 0)
	err := store.Put(tsk.parts, parts)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	fwk.Register(reflect.TypeOf(task3{}))
}
