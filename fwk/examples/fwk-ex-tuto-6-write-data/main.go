package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-hep/fwk"

	// job is the scripting interface to 'fwk'
	"github.com/go-hep/fwk/job"

	// we need to access some tools defined in testdata (the ascii InputStream)
	// so we need to directly import that package
	"github.com/go-hep/fwk/testdata"

	// for persistency
	"github.com/go-hep/fwk/rio"
)

var (
	lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	evtmax = flag.Int64("evtmax", -1, "number of events to process")
	nprocs = flag.Int("nprocs", -1, "number of events to process concurrently")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fwk-ex-tuto-6-write-data [options] <input-ascii> <output-file>

ex:
 $ %[1]s -l=INFO -evtmax=100 ./input.ascii ./output.rio

options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	input := "input.ascii"
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
		Type: "github.com/go-hep/fwk/testdata.task2",
		Name: "t2",
		Props: job.P{
			"Input":  "t1-ints1",
			"Output": "t1-ints1-massaged",
		},
	})

	f, err := os.Open(input)
	if err != nil {
		app.Errorf("could not open file [%s]: %v\n", input, err)
		os.Exit(1)
	}
	defer f.Close()

	// create an input-stream, reading from some io.Reader
	// note we create it after the one that consumes these integers
	// to exercize the automatic data-flow scheduling.
	app.Create(job.C{
		Type: "github.com/go-hep/fwk.InputStream",
		Name: "input",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "t1-ints1",               // location where to publish our data
					Type: reflect.TypeOf(int64(0)), // type of that data
				},
			},
			"Streamer": &testdata.InputStream{
				R: f,
			},
		},
	})

	// output
	app.Create(job.C{
		Type: "github.com/go-hep/fwk.OutputStream",
		Name: "rio-output",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "t1-ints1-massaged",      // location of data to write out
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

$ ::: fwk-ex-tuto-6-write-data...
t2                   INFO configure...
t2                   INFO configure... [done]
t2                   INFO start...
t2                   INFO proc... (id=1|0) => [1 -> 1]
t2                   INFO proc... (id=0|1) => [0 -> 0]
t2                   INFO proc... (id=2|0) => [2 -> 4]
t2                   INFO proc... (id=3|1) => [3 -> 9]
t2                   INFO proc... (id=4|1) => [4 -> 16]
t2                   INFO proc... (id=5|0) => [5 -> 25]
t2                   INFO proc... (id=6|1) => [6 -> 36]
t2                   INFO proc... (id=7|0) => [7 -> 49]
t2                   INFO proc... (id=8|1) => [8 -> 64]
t2                   INFO proc... (id=9|0) => [9 -> 81]
t2                   INFO proc... (id=10|1) => [10 -> 100]
t2                   INFO proc... (id=11|0) => [11 -> 121]
t2                   INFO proc... (id=12|1) => [12 -> 144]
t2                   INFO proc... (id=13|0) => [13 -> 169]
t2                   INFO proc... (id=14|1) => [14 -> 196]
t2                   INFO proc... (id=15|0) => [15 -> 225]
t2                   INFO proc... (id=16|1) => [16 -> 256]
t2                   INFO proc... (id=17|0) => [17 -> 289]
t2                   INFO proc... (id=18|1) => [18 -> 324]
t2                   INFO proc... (id=19|0) => [19 -> 361]
app                  INFO workers done: 1/2
app                  INFO workers done: 2/2
t2                   INFO stop...
app                  INFO cpu: 5.096738ms
app                  INFO mem: alloc:            139 kB
app                  INFO mem: tot-alloc:        170 kB
app                  INFO mem: n-mallocs:       2185
app                  INFO mem: n-frees:          775
app                  INFO mem: gc-pauses:          0 ms
::: fwk-ex-tuto-6-write-data... [done] (cpu=5.347258ms)
*/
