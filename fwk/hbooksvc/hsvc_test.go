// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbooksvc

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/job"
)

const (
	nentries = 100
	nhists   = 10
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
			Type:  "go-hep.org/x/hep/fwk/hbooksvc.testhsvc",
			Name:  fmt.Sprintf("t%03d", i),
			Props: job.P{},
		})
	}

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/hbooksvc.hsvc",
		Name: "histsvc",
		Props: job.P{
			"Streams": map[string]Stream{
				"/my-hist": {
					Name: "hist-seq.rio",
					Mode: Write,
				},
			},
		},
	})

	app.Run()
	os.Remove("hist-seq.rio")
}

func TestHbookSvcConc(t *testing.T) {

	for _, nprocs := range []int{1, 2, 4, 8} {
		app := newapp(nentries, nprocs)
		app.Infof("=== nprocs: %d ===\n", nprocs)

		for i := 0; i < nhists; i++ {
			app.Create(job.C{
				Type:  "go-hep.org/x/hep/fwk/hbooksvc.testhsvc",
				Name:  fmt.Sprintf("t%03d", i),
				Props: job.P{},
			})
		}

		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/hbooksvc.hsvc",
			Name: "histsvc",
			Props: job.P{
				"Streams": map[string]Stream{
					"/my-hist": {
						Name: fmt.Sprintf("hist-conc-%d.rio", nprocs),
						Mode: Write,
					},
				},
			},
		})

		app.Run()
		os.Remove(fmt.Sprintf("hist-conc-%d.rio", nprocs))
	}
}

func TestHbookStreamName(t *testing.T) {
	var svc hsvc
	for _, test := range []struct {
		name string
		want []string
	}{
		{
			name: "histo",
			want: []string{"", "histo"},
		},
		{
			name: "/histo",
			want: []string{"", "histo"},
		},
		{
			name: "/histo/",
			want: []string{"", "histo"},
		},
		{
			name: "/my-stream/histo",
			want: []string{"my-stream", "histo"},
		},
		{
			name: "my-stream/histo",
			want: []string{"my-stream", "histo"},
		},
		{
			name: "my-stream/histo/",
			want: []string{"my-stream", "histo"},
		},
		{
			name: "/my-stream/histo/",
			want: []string{"my-stream", "histo"},
		},
		{
			name: "/my-stream/hdir/histo",
			want: []string{"my-stream", "hdir/histo"},
		},
		{
			name: "/my-stream/hdir/histo/",
			want: []string{"my-stream", "hdir/histo"},
		},
		{
			name: "my-stream/hdir/histo",
			want: []string{"my-stream", "hdir/histo"},
		},
		{
			name: "my-stream/hdir/histo/",
			want: []string{"my-stream", "hdir/histo"},
		},
	} {
		stream, name := svc.split(test.name)
		got := []string{stream, name}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("test.split(%q): got=%v. want=%v\n", test.name, got, test.want)
		}
	}
}

type testhsvc struct {
	fwk.TaskBase

	hsvc   fwk.HistSvc
	h1d    fwk.H1D
	stream string
}

func (tsk *testhsvc) Configure(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *testhsvc) StartTask(ctx fwk.Context) error {
	var err error

	svc, err := ctx.Svc("histsvc")
	if err != nil {
		return err
	}

	tsk.hsvc = svc.(fwk.HistSvc)

	if !strings.HasPrefix(tsk.stream, "/") {
		tsk.stream = "/" + tsk.stream
	}
	if strings.HasSuffix(tsk.stream, "/") {
		tsk.stream = tsk.stream[:len(tsk.stream)-1]
	}

	tsk.h1d, err = tsk.hsvc.BookH1D(tsk.stream+"/h1d-"+tsk.Name(), 100, -10, 10)
	if err != nil {
		return err
	}

	return err
}

func (tsk *testhsvc) StopTask(ctx fwk.Context) error {
	var err error

	h := tsk.h1d.Hist
	if got := h.Entries(); got != nentries {
		return fwk.Errorf("got %d entries. want=%d", got, nentries)
	}
	if got, want := h.XMean(), 49.5; got != want {
		return fwk.Errorf("got mean=%v. want=%v", got, want)
	}

	if got, want := h.XRMS(), 57.301832431432764; got != want {
		return fwk.Errorf("got RMS=%v. want=%v", got, want)
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
		stream:   "",
	}

	err = tsk.DeclProp("Stream", &tsk.stream)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(testhsvc{}), newtesthsvc)
}
