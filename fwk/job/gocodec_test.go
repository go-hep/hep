// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"bytes"
	"reflect"
	"testing"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/testdata"
)

func TestGoEncode(t *testing.T) {
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

	cfg2 := C{
		Type: "go-hep.org/x/hep/fwk/testdata.svc1",
		Name: "svc1",
		Props: P{
			"Int":    testdata.MyInt(12),
			"Struct": testdata.MyStruct{I: 12},
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

	comp1 := job.Create(cfg1)
	prop11 := P{
		"Ints1": "t1-ints1-modified",
	}
	job.SetProp(comp1, "Ints1", prop11["Ints1"])

	job.Create(cfg2)

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
		{
			Type: StmtSetProp,
			Data: C{
				Type:  comp1.Type(),
				Name:  comp1.Name(),
				Props: prop11,
			},
		},
		{
			Type: StmtCreate,
			Data: cfg2,
		},
	}

	stmts := job.Stmts()

	if !reflect.DeepEqual(exp, stmts) {
		t.Fatalf("unexpected statments:\nexp=%#v\ngot=%#v\n", exp, stmts)
	}

	buf := new(bytes.Buffer)
	enc := NewGoEncoder(buf)
	err := enc.Encode(stmts)
	if err != nil {
		t.Fatalf("error go-encoding: %v\n", err)
	}
}
