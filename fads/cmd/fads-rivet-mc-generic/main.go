// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fads-rivet-mc-generic is a command mirroring the MC_GENERIC analysis example from Rivet.
//
// More informations about Rivet: https://rivet.hepforge.org/
//
// fads-rivet-mc-generic reads HepMC events from some input file and
// runs a single task: mc-generic.
// mc-generic is modeled after Rivet's MC_GENERIC analysis (rivet/src/Analyses/MC_GENERIC.cc)
//
// mc-generic selects final state particles passing some acceptance
// cuts (|eta|<5 && pt>0.5 GeV), and fills performance plots for these final
// state particles (and the sub-sample of the charged final state particles.)
//
// Example:
//
//  $> curl -O -L http://www.hepforge.org/archive/rivet/Z-hadronic-LEP.hepmc
//  $> fads-rivet-mc-generic -nprocs=1 ./Z-hadronic-LEP.hepmc
//  ::: fads-rivet-mc-generic...
//  app                  INFO workers done: 1/1
//  app                  INFO cpu: 6.115538196s
//  app                  INFO mem: alloc:          12997 kB
//  app                  INFO mem: tot-alloc:     751784 kB
//  app                  INFO mem: n-mallocs:    7300674
//  app                  INFO mem: n-frees:      7336526
//  app                  INFO mem: gc-pauses:          2 ms
//  ::: fads-rivet-mc-generic... [done] (time=6.11575694s)
//
//  $> rio2yoda rivet.rio >| rivet.yoda
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/pprof"
	"time"

	"go-hep.org/x/hep/fads"
	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/hbooksvc"
	"go-hep.org/x/hep/fwk/job"
	"go-hep.org/x/hep/hepmc"
)

var (
	profFlag   = flag.String("profile", "", "filename of cpuprofile")
	lvlFlag    = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	evtmaxFlag = flag.Int("evtmax", -1, "number of events to process")
	nprocsFlag = flag.Int("nprocs", -1, "number of concurrent events to process")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fads-rivet-mc-generic [options] <hepmc-input-file>

ex:
 $ fads-rivet-mc-generic -l=INFO -evtmax=-1 ./testdata/hepmc.data

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		os.Exit(2)
	}

	start := time.Now()

	fmt.Printf("::: fads-rivet-mc-generic...\n")
	if *profFlag != "" {
		f, err := os.Create(*profFlag)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("could not start CPU profile: %+v", err)
		}
		defer pprof.StopCPUProfile()
	}

	app := job.New(job.P{
		"EvtMax":   int64(*evtmaxFlag),
		"NProcs":   *nprocsFlag,
		"MsgLevel": job.MsgLevel(*lvlFlag),
	})

	ifname := flag.Arg(0)

	// create histogram service
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/hbooksvc.hsvc",
		Name: "histsvc",
		Props: job.P{
			"Streams": map[string]hbooksvc.Stream{
				"/MC_GENERIC": {
					Name: "rivet.rio",
					Mode: hbooksvc.Write,
				},
			},
		},
	})

	// read HepMC data
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk.InputStream",
		Name: "hepmc-streamer",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "/fads/McEvent",
					Type: reflect.TypeOf(hepmc.Event{}),
				},
			},
			"Streamer": &fads.HepMcStreamer{
				Name: ifname,
			},
		},
	})

	app.Create(job.C{
		Type: "main.McGeneric",
		Name: "mc-generic",
		Props: job.P{
			"Input": "/fads/McEvent",
		},
	})

	app.Run()
	fmt.Printf("::: fads-rivet-mc-generic... [done] (time=%v)\n", time.Since(start))
}
