package fwk_test

import (
	"testing"

	"github.com/go-hep/fwk/job"
	_ "github.com/go-hep/fwk/testdata"
)

func newapp(evtmax int64, nprocs int) *job.Job {
	app := job.New(nil, job.P{
		"EvtMax":   evtmax,
		"NProcs":   nprocs,
		"MsgLevel": job.MsgLevel("ERROR"),
	})
	return app
}

func TestSimpleSeqApp(t *testing.T) {

	app := newapp(10, 0)
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t0",
		Props: job.P{
			"Floats1": "t0-floats1",
			"Floats2": "t0-floats2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Floats1": "t1-floats1",
			"Floats2": "t2-floats2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-floats1",
			"Output": "t1-floats1-massaged",
		},
	})

	app.Run()
}

func TestSimpleConcApp(t *testing.T) {

	for _, nprocs := range []int{1, 2, 4, 8} {
		app := newapp(10, nprocs)
		app.Create(job.C{
			Type: "github.com/go-hep/fwk/testdata.task1",
			Name: "t0",
			Props: job.P{
				"Floats1": "t0-floats1",
				"Floats2": "t0-floats2",
			},
		})

		app.Create(job.C{
			Type: "github.com/go-hep/fwk/testdata.task1",
			Name: "t1",
			Props: job.P{
				"Floats1": "t1-floats1",
				"Floats2": "t2-floats2",
			},
		})

		app.Create(job.C{
			Type: "github.com/go-hep/fwk/testdata.task2",
			Name: "t2",
			Props: job.P{
				"Input":  "t1-floats1",
				"Output": "t1-floats1-massaged",
			},
		})
		app.Run()
	}
}

func TestDuplicateProperty(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected a panic")
		}
		exp := `fwk.DeclOutPort: component [t0] already declared out-port with name [t0-floats1 (type=float64)].
fwk.DeclOutPort: component [t1] is trying to add a duplicate out-port [t0-floats1 (type=float64)]`
		if r.(error).Error() != exp {
			t.Fatalf("expected error {%s}.\ngot={%v}\n", exp, r)
		}
	}()

	app := newapp(1, 1)
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t0",
		Props: job.P{
			"Floats1": "t0-floats1",
			"Floats2": "t0-floats2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Floats1": "t0-floats1",
			"Floats2": "t0-floats2",
		},
	})
	app.Run()
}

func TestInputStream(t *testing.T) {
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4} {
			app := newapp(evtmax, nprocs)

			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-floats1",
					"Output": "t1-floats1-massaged",
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.inputstream",
				Name: "input",
				Props: job.P{
					"Output": "t1-floats1",
				},
			})
			app.Run()
		}
	}
}

func TestOutputStream(t *testing.T) {
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4} {
			app := newapp(evtmax, nprocs)

			// put output-stream before 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.outputstream",
				Name: "output",
				Props: job.P{
					"Input": "t1-floats1-massaged",
				},
			})

			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-floats1",
					"Output": "t1-floats1-massaged",
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.inputstream",
				Name: "input",
				Props: job.P{
					"Output": "t1-floats1",
				},
			})
			app.Run()
		}
	}
}
