package fwk_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/go-hep/fwk"
	"github.com/go-hep/fwk/job"
	"github.com/go-hep/fwk/testdata"
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
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t2-ints2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "t1-ints1-massaged",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.svc1",
		Name: "svc1",
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
				"Ints1": "t0-ints1",
				"Ints2": "t0-ints2",
			},
		})

		app.Create(job.C{
			Type: "github.com/go-hep/fwk/testdata.task1",
			Name: "t1",
			Props: job.P{
				"Ints1": "t1-ints1",
				"Ints2": "t2-ints2",
			},
		})

		app.Create(job.C{
			Type: "github.com/go-hep/fwk/testdata.task2",
			Name: "t2",
			Props: job.P{
				"Input":  "t1-ints1",
				"Output": "t1-ints1-massaged",
			},
		})
		app.Run()
	}
}

func TestDuplicateOutputPort(t *testing.T) {
	app := newapp(1, 1)
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})
	err := app.App().Run()
	if err == nil {
		t.Fatalf("expected an error\n")
	}
	exp := fwk.Errorf(`fwk.DeclOutPort: component [t0] already declared out-port with name [t0-ints1 (type=int64)].
fwk.DeclOutPort: component [t1] is trying to add a duplicate out-port [t0-ints1 (type=int64)]`)
	if !reflect.DeepEqual(err, exp) {
		t.Fatalf("invalid error.\nexp=%v (type=%[1]T)\ngot=%v (type=%[2]T)\n", exp, err)
	}
}

func TestMissingInputPort(t *testing.T) {
	app := newapp(1, 1)
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1--NOT-THERE",
			"Output": "t2-ints2",
		},
	})

	err := app.App().Run()
	if err == nil {
		t.Fatalf("expected an error\n")
	}
	exp := fwk.Errorf("dataflow: component [%s] declared port [t1-ints1--NOT-THERE] as input but NO KNOWN producer", "t2")
	if !reflect.DeepEqual(err, exp) {
		t.Fatalf("invalid error.\nexp=%v (type=%[1]T)\ngot=%v (type=%[2]T)\n", exp, err)
	}
}

func getsum(n int64) int64 {
	sum := int64(0)
	for i := int64(0); i < n; i++ {
		sum += i
	}
	return sum
}

func getsumsq(n int64) int64 {
	sum := int64(0)
	for i := int64(0); i < n; i++ {
		sum += i * i
	}
	return sum
}

func newTestReader(max int) io.Reader {
	buf := new(bytes.Buffer)
	for i := 0; i < max; i++ {
		fmt.Fprintf(buf, "%d\n", int64(i))
	}
	return buf
}

func TestInputStream(t *testing.T) {
	const max = 1000
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4, 8} {
			nmax := evtmax
			if nmax < 0 {
				nmax = max
			}

			app := newapp(evtmax, nprocs)

			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-ints1",
					"Output": "t1-ints1-massaged",
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk.InputStream",
				Name: "input",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &testdata.InputStream{
						R: newTestReader(max),
					},
				},
			})

			// check we read the expected amount values
			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.reducer",
				Name: "reducer",
				Props: job.P{
					"Input": "t1-ints1",
					"Sum":   getsum(nmax),
				},
			})

			app.Run()

		}
	}
}

func TestOutputStream(t *testing.T) {
	const max = 1000
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4, 8} {
			nmax := evtmax
			if nmax < 0 {
				nmax = max
			}

			app := newapp(evtmax, nprocs)

			fname := fmt.Sprintf("test-output-stream_%d_%d.txt", evtmax, nprocs)
			w, err := os.Create(fname)
			if err != nil {
				t.Fatalf("could not create output file [%s]: %v\n", fname, err)
			}
			defer w.Close()

			// put output-stream before 'reducer', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk.OutputStream",
				Name: "output",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1-massaged",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &testdata.OutputStream{
						W: w,
					},
				},
			})

			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-ints1",
					"Output": "t1-ints1-massaged",
				},
			})

			// check we read the expected amount values
			app.Create(job.C{
				Type: "github.com/go-hep/fwk/testdata.reducer",
				Name: "reducer",
				Props: job.P{
					"Input": "t1-ints1",
					"Sum":   getsum(nmax),
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "github.com/go-hep/fwk.InputStream",
				Name: "input",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &testdata.InputStream{
						R: newTestReader(max),
					},
				},
			})
			app.Run()

			err = w.Close()
			if err != nil {
				t.Fatalf("could not close file [%s]: %v\n", fname, err)
			}
			w, err = os.Open(fname)
			if err != nil {
				t.Fatalf("could not open file [%s]: %v\n", fname, err)
			}
			defer w.Close()
			exp := getsumsq(nmax)
			sum := int64(0)
			for {
				var val int64
				_, err = fmt.Fscanf(w, "%d\n", &val)
				if err != nil {
					break
				}
				sum += val
			}
			if err == io.EOF {
				err = nil
			}
			if err != nil {
				t.Fatalf("problem scanning output file [%s]: %v\n", fname, err)
			}
			if sum != exp {
				t.Fatalf("problem validating file [%s]: expected sum=%d. got=%d\n",
					fname, exp, sum,
				)
			}
			os.Remove(fname)
		}
	}
}
