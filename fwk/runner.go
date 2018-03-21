// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"go-hep.org/x/hep/fwk/fsm"
)

// irunner wraps an appmgr to implement fwk.Scripter
type irunner struct {
	app *appmgr
}

var _ Scripter = (*irunner)(nil)

func (ui irunner) lvl() Level {
	return ui.app.msg.lvl
}

func (ui irunner) state() fsm.State {
	return ui.app.state
}

func (ui *irunner) Configure() error {
	var err error
	ctx := ctxType{
		id:    0,
		slot:  0,
		store: nil,
		msg:   newMsgStream("<root>", ui.lvl(), nil),
	}

	err = ui.app.configure(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Start() error {
	var err error
	ctx := ctxType{
		id:    0,
		slot:  0,
		store: nil,
		msg:   newMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsm.Configured {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsm.Configured)
	}

	err = ui.app.start(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Run(evtmax int64) error {
	var err error
	ctx := ctxType{
		id:    0,
		slot:  0,
		store: nil,
		msg:   newMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsm.Started {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsm.Started)
	}

	err = ui.app.run(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Stop() error {
	var err error
	ctx := ctxType{
		id:    0,
		slot:  0,
		store: nil,
		msg:   newMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsm.Running {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsm.Running)
	}

	err = ui.app.stop(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Shutdown() error {
	var err error
	ctx := ctxType{
		id:    0,
		slot:  0,
		store: nil,
		msg:   newMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsm.Stopped {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsm.Stopped)
	}

	err = ui.app.start(ctx)
	if err != nil {
		return err
	}

	return err
}
