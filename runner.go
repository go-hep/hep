package fwk

type irunner struct {
	app *appmgr
}

func (ui irunner) lvl() Level {
	return ui.app.msg.lvl
}

func (ui irunner) state() fsm {
	return ui.app.state
}

func (ui *irunner) Configure() error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", ui.lvl(), nil),
	}

	err = ui.app.configure(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Start() error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsmCONFIGURED {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsmCONFIGURED)
	}

	err = ui.app.start(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Run(evtmax int64) error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsmSTARTED {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsmSTARTED)
	}

	err = ui.app.run(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Stop() error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsmRUNNING {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsmRUNNING)
	}

	err = ui.app.stop(ctx)
	if err != nil {
		return err
	}

	return err
}

func (ui *irunner) Shutdown() error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", ui.lvl(), nil),
	}

	if ui.state() < fsmSTOPPED {
		return Errorf("fwk: invalid app state (%v). need at least %s", ui.state(), fsmSTOPPED)
	}

	err = ui.app.start(ctx)
	if err != nil {
		return err
	}

	return err
}
