// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	// job is the scripting interface to 'fwk'
	"go-hep.org/x/hep/fwk/job"

	// side-effect import 'testdata'.
	// merely importing it will register the components defined in this package
	// with the fwk components' factory.
	_ "go-hep.org/x/hep/fwk/testdata"
)

var (
	lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	evtmax = flag.Int64("evtmax", 10, "number of events to process")
	nprocs = flag.Int("nprocs", -1, "number of events to process concurrently")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options]

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
	app := job.New(job.P{
		"EvtMax":   *evtmax,
		"NProcs":   *nprocs,
		"MsgLevel": job.MsgLevel(*lvl),
	})

	// create a task that reads integers from some location
	// and publish the square of these integers under some other location
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk/testdata.task2",
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
		Type: "go-hep.org/x/hep/fwk/testdata.task1",
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

/*
output:

$ fwk-ex-tuto-1
::: fwk-ex-tuto-1...
t2                   INFO configure...
t2                   INFO configure... [done]
t1                   INFO configure ...
t1                   INFO configure ... [done]
t2                   INFO start...
t1                   INFO start...
app                  INFO >>> running evt=0...
t1                   INFO proc... (id=0|0) => [10, 20]
t2                   INFO proc... (id=0|0) => [10 -> 100]
app                  INFO >>> running evt=1...
t1                   INFO proc... (id=1|0) => [10, 20]
t2                   INFO proc... (id=1|0) => [10 -> 100]
app                  INFO >>> running evt=2...
t1                   INFO proc... (id=2|0) => [10, 20]
t2                   INFO proc... (id=2|0) => [10 -> 100]
app                  INFO >>> running evt=3...
t1                   INFO proc... (id=3|0) => [10, 20]
t2                   INFO proc... (id=3|0) => [10 -> 100]
app                  INFO >>> running evt=4...
t1                   INFO proc... (id=4|0) => [10, 20]
t2                   INFO proc... (id=4|0) => [10 -> 100]
app                  INFO >>> running evt=5...
t1                   INFO proc... (id=5|0) => [10, 20]
t2                   INFO proc... (id=5|0) => [10 -> 100]
app                  INFO >>> running evt=6...
t1                   INFO proc... (id=6|0) => [10, 20]
t2                   INFO proc... (id=6|0) => [10 -> 100]
app                  INFO >>> running evt=7...
t1                   INFO proc... (id=7|0) => [10, 20]
t2                   INFO proc... (id=7|0) => [10 -> 100]
app                  INFO >>> running evt=8...
t1                   INFO proc... (id=8|0) => [10, 20]
t2                   INFO proc... (id=8|0) => [10 -> 100]
app                  INFO >>> running evt=9...
t1                   INFO proc... (id=9|0) => [10, 20]
t2                   INFO proc... (id=9|0) => [10 -> 100]
t2                   INFO stop...
t1                   INFO stop...
::: fwk-ex-tuto-1... [done] (cpu=5.482751ms)

*/
