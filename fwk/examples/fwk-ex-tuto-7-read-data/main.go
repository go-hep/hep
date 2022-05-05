// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

	"go-hep.org/x/hep/fwk"

	// job is the scripting interface to 'fwk'
	"go-hep.org/x/hep/fwk/job"

	// for persistency
	"go-hep.org/x/hep/fwk/rio"

	// we need to access some tools defined in fwktest
	// so we need to directly import that package
	_ "go-hep.org/x/hep/fwk/internal/fwktest"
)

var (
	lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	evtmax = flag.Int64("evtmax", -1, "number of events to process")
	nprocs = flag.Int("nprocs", -1, "number of events to process concurrently")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options] <input-file> <output-file>

ex:
 $ %[1]s -l=INFO -evtmax=100 ./input.rio ./output.rio

options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	input := "input.rio"
	if flag.NArg() > 0 {
		input = flag.Arg(0)
	}

	output := "output.rio"
	if flag.NArg() > 1 {
		output = flag.Arg(1)
	}

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
		Type: "go-hep.org/x/hep/fwk/internal/fwktest.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1-massaged",
			"Output": "t1-ints1-massaged-new",
		},
	})

	// create an input-stream, reading from some io.Reader
	// note we create it after the one that consumes these integers
	// to exercize the automatic data-flow scheduling.
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk.InputStream",
		Name: "input",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "t1-ints1-massaged",      // location where to publish our data
					Type: reflect.TypeOf(int64(0)), // type of that data
				},
			},
			"Streamer": &rio.InputStreamer{
				Names: []string{input},
			},
		},
	})

	// output
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk.OutputStream",
		Name: "rio-output",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "t1-ints1-massaged-new",  // location of data to write out
					Type: reflect.TypeOf(int64(0)), // type of that data
				},
			},
			"Streamer": &rio.OutputStreamer{
				Name: output,
			},
		},
	})

	// run the application
	app.Run()

	fmt.Printf("::: %s... [done] (cpu=%v)\n", os.Args[0], time.Since(start))
}

/*
output:

$ fwk-ex-tuto-7-read-data
::: fwk-ex-tuto-7-read-data...
t2                   INFO configure...
t2                   INFO configure... [done]
t2                   INFO start...
t2                   INFO proc... (id=0|1) => [1 -> 1]
t2                   INFO proc... (id=1|0) => [0 -> 0]
t2                   INFO proc... (id=2|1) => [4 -> 16]
t2                   INFO proc... (id=3|0) => [9 -> 81]
t2                   INFO proc... (id=4|1) => [16 -> 256]
t2                   INFO proc... (id=5|0) => [25 -> 625]
t2                   INFO proc... (id=6|1) => [36 -> 1296]
t2                   INFO proc... (id=7|0) => [49 -> 2401]
t2                   INFO proc... (id=8|1) => [81 -> 6561]
t2                   INFO proc... (id=9|0) => [64 -> 4096]
t2                   INFO proc... (id=10|1) => [121 -> 14641]
t2                   INFO proc... (id=11|0) => [100 -> 10000]
t2                   INFO proc... (id=12|1) => [169 -> 28561]
t2                   INFO proc... (id=13|0) => [144 -> 20736]
t2                   INFO proc... (id=14|1) => [196 -> 38416]
t2                   INFO proc... (id=15|0) => [225 -> 50625]
t2                   INFO proc... (id=16|1) => [256 -> 65536]
t2                   INFO proc... (id=17|0) => [289 -> 83521]
t2                   INFO proc... (id=18|1) => [324 -> 104976]
t2                   INFO proc... (id=19|0) => [361 -> 130321]
app                  INFO workers done: 1/2
app                  INFO workers done: 2/2
t2                   INFO stop...
app                  INFO cpu: 7.033018ms
app                  INFO mem: alloc:            197 kB
app                  INFO mem: tot-alloc:       1058 kB
app                  INFO mem: n-mallocs:       2953
app                  INFO mem: n-frees:         2332
app                  INFO mem: gc-pauses:          1 ms
::: fwk-ex-tuto-7-read-data... [done] (cpu=7.271177ms)
*/
