package testdata

import (
	"io"
	"reflect"

	"github.com/go-hep/fwk"
)

type InputStream struct {
	output string
	max    int
	data   chan indata
}

func (stream *InputStream) Connect(ports []fwk.Port) error {
	var err error
	stream.data = make(chan indata)
	stream.output = ports[0].Name
	stream.max = 10000

	go func() {
		for i := 0; i < stream.max; i++ {
			stream.data <- indata{val: float64(i)}
		}
		stream.data <- indata{err: io.EOF}
	}()

	return err
}

func (stream *InputStream) Read(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	data := <-stream.data

	if data.err != nil {
		return data.err
	}

	err = store.Put(stream.output, data.val)
	if err != nil {
		return err
	}

	return err
}

func (stream *InputStream) Disconnect() error {
	var err error
	return err
}

type indata struct {
	val float64
	err error
}

type inputstream struct {
	fwk.TaskBase

	output string
	max    int

	data chan indata
}

func (tsk *inputstream) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(float64(1)))
	if err != nil {
		return err
	}

	return err
}

func (tsk *inputstream) StartTask(ctx fwk.Context) error {
	var err error

	go func() {
		for i := 0; i < tsk.max; i++ {
			tsk.data <- indata{val: float64(i)}
		}
		tsk.data <- indata{err: io.EOF}
	}()

	return err
}

func (tsk *inputstream) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *inputstream) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	data := <-tsk.data

	if data.err != nil {
		return data.err
	}

	err = store.Put(tsk.output, data.val)
	if err != nil {
		return err
	}

	return err
}

func newInputstream(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &inputstream{
		TaskBase: fwk.NewTask(typ, name, mgr),
		output:   "Output",
		max:      10000,
		data:     make(chan indata),
	}

	err = tsk.DeclProp("Max", &tsk.max)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(inputstream{}), newInputstream)
}
