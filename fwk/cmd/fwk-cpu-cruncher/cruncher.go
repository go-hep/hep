package main

import (
	"reflect"
	"time"

	"go-hep.org/x/hep/fwk"
)

type CPUCruncher struct {
	fwk.TaskBase

	cpuch chan int64
	quit  chan struct{}

	inputs  []string
	outputs []string
	cpus    []int64
}

func (tsk *CPUCruncher) Configure(ctx fwk.Context) error {
	var err error

	for _, input := range tsk.inputs {
		if input == "" {
			continue
		}
		err = tsk.DeclInPort(input, reflect.TypeOf(int64(0)))
		if err != nil {
			return err
		}
	}

	for _, output := range tsk.outputs {
		if output == "" {
			continue
		}
		err = tsk.DeclOutPort(output, reflect.TypeOf(int64(0)))
		if err != nil {
			return err
		}
	}

	if len(tsk.cpus) <= 0 {
		msg := ctx.Msg()
		msg.Errorf("invalid cpu-timings list: %v\n", tsk.cpus)
		return fwk.Errorf("invalid cpu-timings")
	}

	return err
}

func (tsk *CPUCruncher) StartTask(ctx fwk.Context) error {
	var err error

	go func() {
		i := 0
		n := len(tsk.cpus)
		for {
			select {
			case <-tsk.quit:
				return
			default:
				if i >= n {
					i = 0
				}
				v := tsk.cpus[i]
				i++
				tsk.cpuch <- v
			}
		}
	}()

	return err
}

func (tsk *CPUCruncher) StopTask(ctx fwk.Context) error {
	var err error
	go func() { tsk.quit <- struct{}{} }()
	return err
}

func (tsk *CPUCruncher) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	for _, input := range tsk.inputs {
		_, err = store.Get(input)
		if err != nil {
			return err
		}
	}

	cpu := <-tsk.cpuch
	time.Sleep(time.Duration(cpu) * time.Microsecond)

	for _, output := range tsk.outputs {
		err = store.Put(output, cpu)
		if err != nil {
			return err
		}
	}

	return err
}

func newCPUCruncher(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &CPUCruncher{
		TaskBase: fwk.NewTask(typ, name, mgr),
		cpuch:    make(chan int64),
		quit:     make(chan struct{}),
		inputs:   make([]string, 0),
		outputs:  make([]string, 0),
		cpus:     make([]int64, 0),
	}

	err = tsk.DeclProp("Inputs", &tsk.inputs)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Outputs", &tsk.outputs)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("CPU", &tsk.cpus)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(CPUCruncher{}), newCPUCruncher)
}
