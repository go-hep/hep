package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// job is the scripting interface to 'fwk'
	"github.com/go-hep/fwk/job"

	// side-effect import 'testdata'.
	// merely importing it will register the components defined in this package
	// with the fwk components' factory.
	_ "github.com/go-hep/fwk/testdata"
)

var (
	g_lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	g_evtmax = flag.Int64("evtmax", 10, "number of events to process")
	g_nprocs = flag.Int("nprocs", 0, "number of events to process concurrently")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fwk-ex-tuto1 [options]

ex:
 $ %[1]s -l=INFO -evtmax=-1

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
	app := job.New(nil, job.P{
		"EvtMax":   *g_evtmax,
		"NProcs":   *g_nprocs,
		"MsgLevel": job.MsgLevel(*g_lvl),
	})

	// create a task that reads integers from some location
	// and publish the square of these integers under some other location
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "t1-ints1-massaged",
		},
	})

	// create a task that publish integers to some location(s)
	// note we create it after the one that consumes these integers
	// to exercize the automatic data-flow scheduling.
	app.Create(job.C{
		Type: "github.com/go-hep/fwk/testdata.task1",
		Name: "t1",
		Props: job.P{
			"Ints1": "t1-ints1",
			"Ints2": "t2-ints2",
			"Int1":  int64(10), // value for the Ints1
			"Int2":  int64(20), // value for the Ints2
		},
	})

	// run the application
	app.Run()

	fmt.Printf("::: %s... [done] (cpu=%v)\n", os.Args[0], time.Since(start))

}
