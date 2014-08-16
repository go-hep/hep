package fwk

import (
	"io"
	"reflect"
)

// InputStreamer reads data from the underlying io.Reader
// and puts it into fwk's Context
type InputStreamer interface {
	Read(ctx Context) error
}

// OutputStreamer gets data from the Context
// and writes it to the underlying io.Writer
type OutputStreamer interface {
	Write(ctx Context) error
}

type indata struct {
	val float64
	err error
}

type InputStream struct {
	SvcBase

	output string
	max    int64

	data chan indata
}

func (svc *InputStream) Configure(ctx Context) error {
	var err error

	err = svc.mgr.DeclOutPort(svc, svc.output, reflect.TypeOf(float64(1)))
	if err != nil {
		return err
	}

	return err
}

func (svc *InputStream) StartSvc(ctx Context) error {
	var err error

	go func() {
		for i := int64(0); i < svc.max; i++ {
			svc.data <- indata{val: float64(i)}
		}
		svc.data <- indata{err: io.EOF}
	}()

	return err
}

func (svc *InputStream) StopSvc(ctx Context) error {
	var err error

	return err
}

func (svc *InputStream) Process(ctx Context) error {
	var err error

	store := ctx.Store()
	data := <-svc.data

	if data.err != nil {
		return data.err
	}

	err = store.Put(svc.output, data.val)
	if err != nil {
		return err
	}

	return err
}

func newInputStream(typ, name string, mgr App) (Component, error) {
	var err error

	svc := &InputStream{
		SvcBase: NewSvc(typ, name, mgr),
		output:  "Output",
		max:     -1,
		data:    make(chan indata),
	}

	err = svc.DeclProp("Max", &svc.max)
	if err != nil {
		return nil, err
	}

	err = svc.DeclProp("Output", &svc.output)
	if err != nil {
		return nil, err
	}

	return svc, err
}

func init() {
	Register(reflect.TypeOf(InputStream{}), newInputStream)
}

// EOF
