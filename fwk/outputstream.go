// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"reflect"
)

// OutputStream implements a task writing data to an OutputStreamer.
//
// OutputStream is concurrent-safe.
//
// OutputStream declares a property 'Ports', a []fwk.Port, which will
// be used to declare the input ports the task will access to,
// writing out data via the underlying OutputStreamer.
//
// OutputStream declares a property 'Streamer', a fwk.OutputStreamer,
// which will be used to actually write data to.
type OutputStream struct {
	TaskBase

	streamer OutputStreamer
	ctrl     StreamControl
}

// Configure declares the input ports defined by the 'Ports' property.
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

// StartTask starts the OutputStreamer task
func (tsk *OutputStream) StartTask(ctx Context) error {
	var err error

	return err
}

// StopTask stops the OutputStreamer task
func (tsk *OutputStream) StopTask(ctx Context) error {
	var err error

	err = tsk.disconnect()
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

// Process gets data from the store and
// writes it out via the underlying OutputStreamer
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
