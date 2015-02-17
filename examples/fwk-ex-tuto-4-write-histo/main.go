package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// job is the scripting interface to 'fwk'
	"github.com/go-hep/fwk/job"

	// for hsbooksvc.Stream
	"github.com/go-hep/fwk/hbooksvc"
)

var (
	lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	evtmax = flag.Int64("evtmax", 100, "number of events to process")
	nprocs = flag.Int("nprocs", -1, "number of events to process concurrently")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options] <input-file> <output-file>

ex:
 $ %[1]s -l=INFO -evtmax=-1 ./input.ascii ./output.ascii

options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	start := time.Now()
	fmt.Printf("::: %s...\n", os.Args[0])

	// create a default fwk application, with some properties
	// extracted from the CLI
	app := job.NewJob(nil, job.P{
		"EvtMax":   *evtmax,
		"NProcs":   *nprocs,
		"MsgLevel": job.MsgLevel(*lvl),
	})

	app.Create(job.C{
		Type: "main.testhsvc",
		Name: "t-01",
		Props: job.P{
			"Stream": "/my-hist",
		},
	})

	app.Create(job.C{
		Type: "main.testhsvc",
		Name: "t-02",
		Props: job.P{
			"Stream": "/my-hist",
		},
	})

	app.Create(job.C{
		Type: "main.testhsvc",
		Name: "t-03",
		Props: job.P{
			"Stream": "", // in-memory temporary hist.
		},
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fwk/hbooksvc.hsvc",
		Name: "histsvc",
		Props: job.P{
			"Streams": map[string]hbooksvc.Stream{
				"/my-hist": {
					Name: "hist.rio",
					Mode: hbooksvc.Write,
				},
			},
		},
	})

	app.Run()
	fmt.Printf("::: %s... [done] (cpu=%v)\n", os.Args[0], time.Since(start))
}

/*
output:

$ fwk-ex-tuto-4-write-histo
::: fwk-ex-tuto-4-write-histo...
app                  INFO workers done: 1/2
app                  INFO workers done: 2/2
t-01                 INFO histo[h1d-t-01]: entries=100 mean=4.5 RMS=2.8722813232690143
t-02                 INFO histo[h1d-t-02]: entries=100 mean=4.5 RMS=2.8722813232690143
t-03                 INFO histo[h1d-t-03]: entries=100 mean=4.5 RMS=2.8722813232690143
app                  INFO cpu: 6.827126ms
app                  INFO mem: alloc:            237 kB
app                  INFO mem: tot-alloc:        555 kB
app                  INFO mem: n-mallocs:       6540
app                  INFO mem: n-frees:         5597
app                  INFO mem: gc-pauses:          1 ms
::: fwk-ex-tuto-4-write-histo... [done] (cpu=7.094478ms)

*/
