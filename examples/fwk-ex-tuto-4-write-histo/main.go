package main

import (
	"flag"
	"fmt"
	"os"

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
		fmt.Fprintf(os.Stderr, `Usage: fwk-ex-tuto1 [options] <input-file> <output-file>

ex:
 $ %[1]s -l=INFO -evtmax=-1 ./input.ascii ./output.ascii

options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()

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
}
