// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/fwk"
	_ "go-hep.org/x/hep/fwk/testdata"
)

func TestStmt(t *testing.T) {

	appcfg := C{
		Name: "app",
		Type: "go-hep.org/x/hep/fwk.appmgr",
		Props: P{
			"EvtMax": int64(10),
			"NProcs": 42,
		},
	}

	cfg0 := C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t0",
		Props: P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	}

	cfg1 := C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t1",
		Props: P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	}

	job := NewJob(
		fwk.NewApp(),
		appcfg.Props,
	)

	if job == nil {
		t.Fatalf("got nil job.Job")
	}

	job.Create(cfg0)
	job.Create(cfg1)

	exp := []Stmt{
		{
			Type: StmtNewApp,
			Data: appcfg,
		},
		{
			Type: StmtCreate,
			Data: cfg0,
		},
		{
			Type: StmtCreate,
			Data: cfg1,
		},
	}

	stmts := job.Stmts()

	if !reflect.DeepEqual(exp, stmts) {
		t.Fatalf("unexpected statments:\nexp=%#v\ngot=%#v\n", exp, stmts)
	}
}

func TestStmtWithProps(t *testing.T) {

	appcfg := C{
		Name: "app",
		Type: "go-hep.org/x/hep/fwk.appmgr",
		Props: P{
			"EvtMax": int64(10),
			"NProcs": 42,
		},
	}

	cfg0 := C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t0",
		Props: P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	}

	cfg1 := C{
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
		Name: "t1",
		Props: P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	}

	job := NewJob(
		fwk.NewApp(),
		appcfg.Props,
	)

	if job == nil {
		t.Fatalf("got nil job.Job")
	}

	comp0 := job.Create(cfg0)
	prop01 := P{
		"Ints1": "t0-ints1-modified",
	}
	prop02 := P{
		"Ints2": "t0-ints2-modified",
	}

	job.SetProp(comp0, "Ints1", prop01["Ints1"])
	job.SetProp(comp0, "Ints2", prop02["Ints2"])

	comp1 := job.Create(cfg1)

	prop11 := P{
		"Ints1": "t1-ints1-modified",
	}
	prop12 := P{
		"Ints2": "t1-ints2-modified",
	}

	job.SetProp(comp1, "Ints1", prop11["Ints1"])
	job.SetProp(comp1, "Ints2", prop12["Ints2"])

	exp := []Stmt{
		{
			Type: StmtNewApp,
			Data: appcfg,
		},
		{
			Type: StmtCreate,
			Data: cfg0,
		},
		{
			Type: StmtSetProp,
			Data: C{
				Type:  comp0.Type(),
				Name:  comp0.Name(),
				Props: prop01,
			},
		},
		{
			Type: StmtSetProp,
			Data: C{
				Type:  comp0.Type(),
				Name:  comp0.Name(),
				Props: prop02,
			},
		},
		{
			Type: StmtCreate,
			Data: cfg1,
		},
		{
			Type: StmtSetProp,
			Data: C{
				Type:  comp1.Type(),
				Name:  comp1.Name(),
				Props: prop11,
			},
		},
		{
			Type: StmtSetProp,
			Data: C{
				Type:  comp1.Type(),
				Name:  comp1.Name(),
				Props: prop12,
			},
		},
	}

	stmts := job.Stmts()

	if !reflect.DeepEqual(exp, stmts) {
		t.Fatalf("unexpected statments:\nexp=%#v\ngot=%#v\n", exp, stmts)
	}
}
