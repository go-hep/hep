package fwk

import (
	"reflect"
)

type OutputStream struct {
	TaskBase

	streamer OutputStreamer
	ctrl     StreamControl
}

func (tsk *OutputStream) Configure(ctx Context) error {
	var err error

	for _, port := range tsk.ctrl.Ports {
		err = tsk.DeclInPort(port.Name, port.Type)
		if err != nil {
			return err
		}
	}

	return err
}

func (tsk *OutputStream) StartTask(ctx Context) error {
	var err error

	return err
}

func (tsk *OutputStream) StopTask(ctx Context) error {
	var err error

	return err
}

func (tsk *OutputStream) connect(ctrl StreamControl) error {
	ctrl.Ports = make([]Port, len(tsk.ctrl.Ports))
	copy(ctrl.Ports, tsk.ctrl.Ports)

	tsk.ctrl = ctrl
	err := tsk.streamer.Connect(ctrl.Ports)
	if err != nil {
		return err
	}

	go tsk.write()

	return err
}

func (tsk *OutputStream) disconnect() error {
	select {
	case tsk.ctrl.Quit <- struct{}{}:
	default:
	}
	return tsk.streamer.Disconnect()
}

func (tsk *OutputStream) write() {
	for {
		select {

		case ctx := <-tsk.ctrl.Ctx:
			tsk.ctrl.Err <- tsk.streamer.Write(ctx)

		case <-tsk.ctrl.Quit:
			return
		}
	}
}

func (tsk *OutputStream) Process(ctx Context) error {
	var err error

	tsk.ctrl.Ctx <- ctx
	err = <-tsk.ctrl.Err
	if err != nil {
		return err
	}

	return err
}

func newOutputStream(typ, name string, mgr App) (Component, error) {
	var err error

	tsk := &OutputStream{
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
	Register(reflect.TypeOf(OutputStream{}), newOutputStream)
}
