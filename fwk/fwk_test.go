// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/internal/fwktest"
	"go-hep.org/x/hep/fwk/job"
)

func newapp(evtmax int64, nprocs int) *job.Job {
	app := job.NewJob(nil, job.P{
		"EvtMax":   evtmax,
		"NProcs":   nprocs,
		"MsgLevel": job.MsgLevel("ERROR"),
	})
	return app
}

func TestSimpleSeqApp(t *testing.T) {

	app := newapp(10, 0)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t2-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "t1-ints1-massaged",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.svc1",
		Name: "svc1",
	})

	app.Run()
}

func TestSimpleConcApp(t *testing.T) {

	for _, nprocs := range []int{1, 2, 4, 8} {
		app := newapp(10, nprocs)
		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
			Name: "t0",
			Props: job.P{
				"Ints1": "t0-ints1",
				"Ints2": "t0-ints2",
			},
		})

		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
			Name: "t1",
			Props: job.P{
				"Ints1": "t1-ints1",
				"Ints2": "t2-ints2",
			},
		})

		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
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
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
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
	want := fmt.Errorf(`fwk.DeclOutPort: component [t0] already declared out-port with name [t0-ints1 (type=int64)].
fwk.DeclOutPort: component [t1] is trying to add a duplicate out-port [t0-ints1 (type=int64)]`)
	if got, want := err.Error(), want.Error(); got != want {
		t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
	}
}

func TestMissingInputPort(t *testing.T) {
	app := newapp(1, 1)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
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
	want := fmt.Errorf("dataflow: component [%s] declared port [t1-ints1--NOT-THERE] as input but NO KNOWN producer", "t2")
	if got, want := err.Error(), want.Error(); got != want {
		t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
	}
}

func TestMismatchPortTypes(t *testing.T) {
	app := newapp(1, 1)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t1-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "data",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task4",
		Name: "t4",
		Props: job.P{
			"Input":  "data",
			"Output": "out-data",
		},
	})

	err := app.App().Run()
	if err == nil {
		t.Fatalf("expected an error\n")
	}
	want := fmt.Errorf(`fwk.DeclInPort: detected type inconsistency for port [data]:
 component=%[1]q port=out type=int64
 component=%[2]q port=in  type=float64
`,
		"t2",
		"t4",
	)
	if got, want := err.Error(), want.Error(); got != want {
		t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
	}
}

func TestPortsCycles(t *testing.T) {
	app := newapp(1, 1)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t1-cycle",
		Props: job.P{
			"Input":  "input",
			"Output": "data-1",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "data-1",
			"Output": "data-2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t3",
		Props: job.P{
			"Input":  "data-2",
			"Output": "input",
		},
	})

	err := app.App().Run()
	if err == nil {
		t.Fatalf("expected an error\n")
	}
	want := fmt.Errorf("dataflow: cycle detected: 1")
	if got, want := err.Error(), want.Error(); got != want {
		t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
	}
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
	for i := range max {
		fmt.Fprintf(buf, "%d\n", int64(i))
	}
	return buf
}

func TestInputStream(t *testing.T) {
	const max = 1000
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4, 8, -1} {
			nmax := evtmax
			if nmax < 0 {
				nmax = max
			}

			app := newapp(evtmax, nprocs)

			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-ints1",
					"Output": "t1-ints1-massaged",
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk.InputStream",
				Name: "input",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &fwktest.InputStream{
						R: newTestReader(max),
					},
				},
			})

			// check we read the expected amount values
			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk/internal/fwktest.reducer",
				Name: "reducer",
				Props: job.P{
					"Input": "t1-ints1-massaged",
					"Sum":   getsumsq(nmax),
				},
			})

			err := app.App().Run()
			if err != nil {
				t.Errorf("error (evtmax=%d nprocs=%d): %v\n", evtmax, nprocs, err)
			}
		}
	}
}

func TestOutputStream(t *testing.T) {
	const max = 1000
	for _, evtmax := range []int64{0, 1, 10, 100, -1} {
		for _, nprocs := range []int{0, 1, 2, 4, 8, -1} {
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
				Type: "go-hep.org/x/hep/fwk.OutputStream",
				Name: "output",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1-massaged",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &fwktest.OutputStream{
						W: w,
					},
				},
			})

			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
				Name: "t2",
				Props: job.P{
					"Input":  "t1-ints1",
					"Output": "t1-ints1-massaged",
				},
			})

			// check we read the expected amount values
			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk/internal/fwktest.reducer",
				Name: "reducer",
				Props: job.P{
					"Input": "t1-ints1-massaged",
					"Sum":   getsumsq(nmax),
				},
			})

			// put input-stream after 't2', to test dataflow re-ordering
			app.Create(job.C{
				Type: "go-hep.org/x/hep/fwk.InputStream",
				Name: "input",
				Props: job.P{
					"Ports": []fwk.Port{
						{
							Name: "t1-ints1",
							Type: reflect.TypeOf(int64(1)),
						},
					},
					"Streamer": &fwktest.InputStream{
						R: newTestReader(max),
					},
				},
			})
			err = app.App().Run()
			if err != nil {
				t.Errorf("error (evtmax=%d nprocs=%d): %v\n", evtmax, nprocs, err)
			}

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

func Benchmark___SeqApp(b *testing.B) {
	app := newapp(100, 0)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t0",
			"Ints2": "t2-ints2",
		},
	})

	input := "t0"
	for i := range 100 {
		name := fmt.Sprintf("tx-%d", i)
		out := fmt.Sprintf("tx-%d", i)
		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
			Name: name,
			Props: job.P{
				"Input":  input,
				"Output": out,
			},
		})
		input = out
	}

	ui := app.App().Scripter()
	err := ui.Configure()
	if err != nil {
		b.Fatalf("error: %v\n", err)
	}

	err = ui.Start()
	if err != nil {
		b.Fatalf("error: %v\n", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ui.Run(-1)
		if err != nil && err != io.EOF {
			b.Fatalf("error: %v\n", err)
		}
	}
}

func Benchmark__ConcApp(b *testing.B) {
	app := newapp(100, 4)
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t0",
		Props: job.P{
			"Ints1": "t0-ints1",
			"Ints2": "t0-ints2",
		},
	})

	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t0",
			"Ints2": "t2-ints2",
		},
	})

	input := "t0"
	for i := range 100 {
		name := fmt.Sprintf("tx-%d", i)
		out := fmt.Sprintf("tx-%d", i)
		app.Create(job.C{
			Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
			Name: name,
			Props: job.P{
				"Input":  input,
				"Output": out,
			},
		})
		input = out
	}

	ui := app.App().Scripter()
	err := ui.Configure()
	if err != nil {
		b.Fatalf("error: %v\n", err)
	}

	err = ui.Start()
	if err != nil {
		b.Fatalf("error: %v\n", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = ui.Run(-1)
		if err != nil && err != io.EOF {
			b.Fatalf("error: %v\n", err)
		}
	}
}
