package hbooksvc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/job"
)

const (
	nentries = 1000
	nhists   = 1000
)

func newapp(evtmax int64, nprocs int) *job.Job {
	app := job.NewJob(nil, job.P{
		"EvtMax":   evtmax,
		"NProcs":   nprocs,
		"MsgLevel": job.MsgLevel("ERROR"),
	})
	return app
}

func TestHbookSvcSeq(t *testing.T) {

	app := newapp(nentries, 0)

	for i := 0; i < nhists; i++ {
		app.Create(job.C{
			Type:  "github.com/go-hep/fwk/hbooksvc.testhsvc",
			Name:  fmt.Sprintf("t%03d", i),
			Props: job.P{},
		})
	}

	app.Create(job.C{
		Type:  "github.com/go-hep/fwk/hbooksvc.hsvc",
		Name:  "histsvc",
		Props: job.P{},
	})

	app.Run()
}

func TestHbookSvcConc(t *testing.T) {

	for _, nprocs := range []int{1, 2, 4, 8} {
		app := newapp(nentries, nprocs)
		app.Infof("=== nprocs: %d ===\n", nprocs)

		for i := 0; i < nhists; i++ {
			app.Create(job.C{
				Type:  "github.com/go-hep/fwk/hbooksvc.testhsvc",
				Name:  fmt.Sprintf("t%03d", i),
				Props: job.P{},
			})
		}

		app.Create(job.C{
			Type:  "github.com/go-hep/fwk/hbooksvc.hsvc",
			Name:  "histsvc",
			Props: job.P{},
		})

		app.Run()
	}
}

type testhsvc struct {
	fwk.TaskBase

	hsvc fwk.HistSvc
	h1d  fwk.H1D
}

func (tsk *testhsvc) Configure(ctx fwk.Context) error {
	var err error

	svc, err := ctx.Svc("histsvc")
	if err != nil {
		return err
	}

	tsk.hsvc = svc.(fwk.HistSvc)

	tsk.h1d, err = tsk.hsvc.BookH1D("h1d-"+tsk.Name(), 100, -10, 10)
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
	if h.Entries() != nentries {
		return fwk.Errorf("expected %d entries. got=%d", nentries, h.Entries())
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
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(testhsvc{}), newtesthsvc)
}
