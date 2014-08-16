package fwk_test

import (
	"testing"

	"github.com/go-hep/fwk/job"
	_ "github.com/go-hep/fwk/testdata"
)

func TestSimpleApp(t *testing.T) {

	app := job.New(nil, job.P{
		"EvtMax":   int64(10),
		"MsgLevel": job.MsgLevel("INFO"),
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "task1",
	})

	app.Run()
}
