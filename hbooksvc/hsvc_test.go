package hbooksvc

import (
	"reflect"
	"testing"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/job"
)

func newapp(evtmax int64, nprocs int) *job.Job {
	app := job.NewJob(nil, job.P{
		"EvtMax":   evtmax,
		"NProcs":   nprocs,
		"MsgLevel": job.MsgLevel("ERROR"),
	})
	return app
}

func TestHbookSvc(t *testing.T) {

	app := newapp(10, 0)
	app.Create(job.C{
		Type:  "github.com/go-hep/fwk/hbooksvc.testhsvc",
		Name:  "t0",
		Props: job.P{},
	})

	app.Create(job.C{
		Type:  "github.com/go-hep/fwk/hbooksvc.hsvc",
		Name:  "histsvc",
		Props: job.P{},
	})

	app.Run()
}

type testhsvc struct {
	fwk.TaskBase

	hsvc fwk.HistSvc
	h1d  fwk.H1D
}

func (tsk *testhsvc) Configure(ctx fwk.Context) error {
	var err error

	// err = tsk.DeclInPort(tsk.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	svc, err := ctx.Svc("histsvc")
	if err != nil {
		return err
	}

	tsk.hsvc = svc.(fwk.HistSvc)

	tsk.h1d, err = tsk.hsvc.BookH1D("h1d", 100, -10, 10)
	if err != nil {
		return err
	}

	return err
}

func (tsk *testhsvc) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *testhsvc) StopTask(ctx fwk.Context) error {
	var err error

	h := tsk.h1d.Hist
	if h.Entries() != 10 {
		return fwk.Errorf("expected 10 entries. got=%d", h.Entries())
	}
	mean := h.Mean()
	if mean != 4.5 {
		return fwk.Errorf("expected mean=%v. got=%v", 4.5, mean)
	}

	rms := h.RMS()
	if rms != 2.8722813232690143 {
		return fwk.Errorf("expected RMS=%v. got=%v", 2.8722813232690143, rms)
	}
	return err
}

func (tsk *testhsvc) Process(ctx fwk.Context) error {
	var err error
	id := ctx.ID()
	tsk.hsvc.FillH1D(tsk.h1d.ID, float64(id), 1)
	return err
}

func newtesthsvc(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &testhsvc{
		TaskBase: fwk.NewTask(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
	}

	// err = tsk.DeclProp("Input", &tsk.input)
	// if err != nil {
	// 	return nil, err
	// }

	// err = tsk.DeclProp("Output", &tsk.output)
	// if err != nil {
	//	return nil, err
	// }

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(testhsvc{}), newtesthsvc)
}
