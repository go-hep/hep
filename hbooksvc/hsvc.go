package hbooksvc

import (
	"reflect"
	"sync"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/fsm"
	"github.com/go-hep/hbook"
)

type h1d struct {
	fwk.H1D
	mu sync.RWMutex
}

type hsvc struct {
	fwk.SvcBase

	h1ds map[fwk.HID]*h1d
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

	return err
}

func (svc *hsvc) StopSvc(ctx fwk.Context) error {
	var err error

	return err
}

func (svc *hsvc) BookH1D(name string, nbins int, low, high float64) (fwk.H1D, error) {
	var err error
	h := fwk.H1D{
		ID:   fwk.HID(name),
		Hist: hbook.NewH1D(nbins, low, high),
	}

	if !(svc.FSMState() < fsm.Running) {
		return h, fwk.Errorf("fwk: can not book histograms during FSM-state %v", svc.FSMState())
	}

	hh := &h1d{H1D: h}
	svc.h1ds[h.ID] = hh
	return hh.H1D, err
}

func (svc *hsvc) FillH1D(id fwk.HID, x, w float64) {
	h := svc.h1ds[id]
	h.mu.Lock()
	h.Hist.Fill(x, w)
	h.mu.Unlock()
}

func newhsvc(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	svc := &hsvc{
		SvcBase: fwk.NewSvc(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
		h1ds: make(map[fwk.HID]*h1d),
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
