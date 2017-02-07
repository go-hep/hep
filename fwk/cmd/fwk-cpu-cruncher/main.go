package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"go-hep.org/x/hep/fwk/job"
)

var (
	lvl     = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	evtmax  = flag.Int("evtmax", 10, "number of events to process")
	nprocs  = flag.Int("nprocs", 0, "number of concurrent events to process")
	dotfile = flag.String("dotfile", "", "path to dotfile for dumping dataflow graph")
	cpu     = flag.Bool("cpu-prof", false, "enable CPU profiling")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fwk-cpu-cruncher [options] <config-file>

ex:
 $ fwk-cpu-cruncher -l=INFO -evtmax=100 ./testdata/athena.json

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() <= 0 {
		fmt.Fprintf(os.Stderr, "** error: needs an input cpu-cruncher configuration file\n")
		flag.Usage()
		os.Exit(1)
	}
	fname := flag.Arg(0)

	start := time.Now()

	fmt.Printf("::: fwk-cpu-cruncher...\n")
	if *cpu {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app := job.New(job.P{
		"EvtMax":   int64(*evtmax),
		"NProcs":   *nprocs,
		"MsgLevel": job.MsgLevel(*lvl),
	})

	if *dotfile != "" {
		dflow := app.App().GetSvc("dataflow")
		if dflow == nil {
			panic(fmt.Errorf("could not retrieve dataflow service"))
		}

		app.SetProp(dflow, "DotFile", *dotfile)
	}

	loadConfig(fname, app)

	app.Run()
	fmt.Printf("::: fwk-cpu-cruncher... [done] (time=%v)\n", time.Since(start))
}
