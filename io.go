package fwk

import (
	"io"
	"reflect"
)

type StreamControl struct {
	Ports []Port
	Ctx   chan Context
	Err   chan error
	Quit  chan struct{}
}

// InputStreamer reads data from the underlying io.Reader
// and puts it into fwk's Context
type InputStreamer interface {
	Connect(ports []Port) error
	Read(ctx Context) error
	Disconnect() error
}

// OutputStreamer gets data from the Context
// and writes it to the underlying io.Writer
type OutputStreamer interface {
	Open(w io.Writer) error
	Write(ctx Context) error
	Close() error
}

type InputStream struct {
	TaskBase

	streamer InputStreamer
	ctrl     StreamControl
}

func (tsk *InputStream) Configure(ctx Context) error {
	var err error

	for _, port := range tsk.ctrl.Ports {
		err = tsk.DeclOutPort(port.Name, port.Type)
		if err != nil {
			return err
		}
	}

	return err
}

func (tsk *InputStream) StartTask(ctx Context) error {
	var err error

	return err
}

func (tsk *InputStream) StopTask(ctx Context) error {
	var err error

	return err
}

func (tsk *InputStream) connect(ctrl StreamControl) error {
	ctrl.Ports = make([]Port, len(tsk.ctrl.Ports))
	copy(ctrl.Ports, tsk.ctrl.Ports)

	tsk.ctrl = ctrl
	err := tsk.streamer.Connect(ctrl.Ports)
	if err != nil {
		return err
	}

	go tsk.read()

	return err
}

func (tsk *InputStream) disconnect() error {
	return tsk.streamer.Disconnect()
}

func (tsk *InputStream) read() {
	for {
		select {

		case ctx := <-tsk.ctrl.Ctx:
			tsk.ctrl.Err <- tsk.streamer.Read(ctx)

		case <-tsk.ctrl.Quit:
			return
		}
	}
}

func (tsk *InputStream) Process(ctx Context) error {
	var err error

	tsk.ctrl.Ctx <- ctx
	err = <-tsk.ctrl.Err
	if err != nil {
		close(tsk.ctrl.Quit)
		return err
	}

	return err
}

func newInputStream(typ, name string, mgr App) (Component, error) {
	var err error

	tsk := &InputStream{
		TaskBase: NewTask(typ, name, mgr),
		streamer: nil,
		ctrl: StreamControl{
			Ports: make([]Port, 0),
		},
	}

	err = tsk.DeclProp("Ports", &tsk.ctrl.Ports)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Streamer", &tsk.streamer)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	Register(reflect.TypeOf(InputStream{}), newInputStream)
}

// EOF
