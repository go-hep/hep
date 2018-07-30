// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"reflect"
)

// InputStream implements a task reading data from an InputStreamer.
//
// InputStream is concurrent-safe.
//
// InputStream declares a property 'Ports', a []fwk.Port, which will
// be used to declare the output ports the streamer will publish,
// loading in data from the underlying InputStreamer.
//
// InputStream declares a property 'Streamer', a fwk.InputStreamer,
// which will be used to actually read data from.
type InputStream struct {
	TaskBase

	streamer InputStreamer
	ctrl     StreamControl
}

// Configure declares the output ports defined by the 'Ports' property.
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

// StartTask starts the InputStreamer task
func (tsk *InputStream) StartTask(ctx Context) error {
	var err error

	return err
}

// StopTask stops the InputStreamer task
func (tsk *InputStream) StopTask(ctx Context) error {
	var err error

	err = tsk.disconnect()
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

// Process loads data from the underlying InputStreamer
// and puts it in the event store.
func (tsk *InputStream) Process(ctx Context) error {
	var err error

	tsk.ctrl.Ctx <- ctx
	err = <-tsk.ctrl.Err

	if err != nil {
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

// inputStream provides fake input events when no input file is present
type inputStream struct {
	TaskBase
}

func (tsk *inputStream) StartTask(ctx Context) error {
	return nil
}

func (tsk *inputStream) StopTask(ctx Context) error {
	return nil
}

func (tsk *inputStream) Process(ctx Context) error {
	return nil
}

func init() {
	Register(reflect.TypeOf(InputStream{}), newInputStream)
}
