// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"go-hep.org/x/hep/fwk/job"
)

func TestDFlowSvcGraph(t *testing.T) {
	app := job.NewJob(nil, job.P{
		"EvtMax":   int64(1),
		"NProcs":   1,
		"MsgLevel": job.MsgLevel("ERROR"),
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t2-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "t1-ints1-massaged",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/testdata.svc1",
		Name: "svc1",
	})

	dflow := app.App().GetSvc("dataflow")
	if dflow == nil {
		t.Fatalf("could not retrieve dataflow svc")
	}

	const (
		dotfile  = "testdata/simple_dflow.dot"
		wantfile = "testdata/simple_dflow.dot.golden"
	)

	app.SetProp(dflow, "DotFile", dotfile)

	app.Run()

	got, err := ioutil.ReadFile(dotfile)
	if err != nil {
		t.Fatalf("could not read %q: %v", dotfile, err)
	}

	want, err := ioutil.ReadFile(wantfile)
	if err != nil {
		t.Fatalf("could not read reference file %q: %v", wantfile, err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("dot files differ.\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	os.Remove(dotfile)
}
