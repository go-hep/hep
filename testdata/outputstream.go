package testdata

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/go-hep/fwk"
)

type outdata struct {
	val float64
	err error
}

type outputstream struct {
	fwk.TaskBase

	input string
	w     io.Writer

	ctx   chan fwk.Context
	quit  chan struct{}
	errch chan error
}

func (tsk *outputstream) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf(float64(1)))
	if err != nil {
		return err
	}

	return err
}

func (tsk *outputstream) StartTask(ctx fwk.Context) error {
	var err error

	go tsk.write()

	return err
}

func (tsk *outputstream) StopTask(ctx fwk.Context) error {
	var err error

	go func() {
		tsk.quit <- struct{}{}
	}()

	return err
}

func (tsk *outputstream) Process(ctx fwk.Context) error {
	var err error

	tsk.ctx <- ctx
	err = <-tsk.errch
	if err != nil {
		close(tsk.ctx)
		return err
	}

	return err
}

func (tsk *outputstream) write() {

	for {
		select {
		case ctx, ok := <-tsk.ctx:
			if !ok {
				return
			}
			store := ctx.Store()
			v, err := store.Get(tsk.input)
			if err != nil {
				tsk.errch <- err
				return
			}
			data := v.(float64)
			_, err = tsk.w.Write([]byte(fmt.Sprintf("%f\n", data)))
			tsk.errch <- err
			if err != nil {
				return
			}

		case <-tsk.quit:
			return
		}
	}
}

func newOutputstream(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &outputstream{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "Input",
		ctx:      make(chan fwk.Context),
		errch:    make(chan error),
		quit:     make(chan struct{}),
		w:        new(bytes.Buffer),
	}

	err = tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.w)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(outputstream{}), newOutputstream)
}
