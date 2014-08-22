package hbooksvc

import (
	"reflect"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/fsm"
	"github.com/go-hep/hbook"
)

type hsvc struct {
	fwk.SvcBase

	h1ds map[fwk.HID]fwk.H1D
}

func (svc *hsvc) Configure(ctx fwk.Context) error {
	var err error

	// err = svc.DeclInPort(svc.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = svc.DeclOutPort(svc.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	return err
}

func (svc *hsvc) StartSvc(ctx fwk.Context) error {
	var err error

	svc.h1ds = make(map[fwk.HID]fwk.H1D)
	return err
}

func (svc *hsvc) StopSvc(ctx fwk.Context) error {
	var err error

	svc.h1ds = make(map[fwk.HID]fwk.H1D)
	return err
}

func (svc *hsvc) BookH1D(name string, nbins int, low, high float64) (fwk.H1D, error) {
	var err error
	h1d := fwk.H1D{
		ID:   fwk.HID(name),
		Hist: hbook.NewH1D(nbins, low, high),
	}

	if !(svc.FSMState() < fsm.Running) {
		return h1d, fwk.Errorf("fwk: can not book histograms during FSM-state %v", svc.FSMState())
	}

	return h1d, err
}

func (svc *hsvc) FillH1D(id fwk.HID, x, w float64) {
	// FIXME(sbinet) make it concurrency-safe
	svc.h1ds[id].Hist.Fill(x, w)
}

func newhsvc(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	svc := &hsvc{
		SvcBase: fwk.NewSvc(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
		h1ds: make(map[fwk.HID]fwk.H1D),
	}

	// err = svc.DeclProp("Input", &svc.input)
	// if err != nil {
	// 	return nil, err
	// }

	// err = svc.DeclProp("Output", &svc.output)
	// if err != nil {
	//	return nil, err
	// }

	return svc, err
}

func init() {
	fwk.Register(reflect.TypeOf(hsvc{}), newhsvc)
}

var _ fwk.HistSvc = (*hsvc)(nil)
