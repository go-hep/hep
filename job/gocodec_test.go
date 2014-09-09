package job

import (
	"reflect"
	"testing"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/testdata"
)

func TestGoEncode(t *testing.T) {
	appcfg := C{
		Name: "app",
		Type: "github.com/go-hep/fwk.appmgr",
		Props: P{
			"EvtMax": int64(10),
			"NProcs": 42,
		},
	}

	cfg0 := C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t0",
		Props: P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	}

	cfg1 := C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	}

	cfg2 := C{
		Type: "github.com/go-hep/fwk/testdata.svc1",
		Name: "svc1",
		Props: P{
			"MyInt": testdata.MyInt(12),
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
			Type: StmtCreate,
			Data: cfg2,
		},
	}

	stmts := job.Stmts()

	if !reflect.DeepEqual(exp, stmts) {
		t.Fatalf("unexpected statments:\nexp=%#v\ngot=%#v\n", exp, stmts)
	}

	enc := NewGoEncoder(nil)
	err := enc.Encode(stmts)
	if err != nil {
		t.Fatalf("error go-encoding: %v\n", err)
	}
}
