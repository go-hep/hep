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
)

var (
	g_lvl    = flag.String("l", "INFO", "message level (DEBUG|INFO|WARN|ERROR)")
	g_evtmax = flag.Int64("evtmax", -1, "number of events to process")
	g_nprocs = flag.Int("nprocs", 0, "number of events to process concurrently")
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

	ifname := "input.ascii"
	if flag.NArg() > 0 {
		ifname = flag.Arg(0)
	}

	ofname := "output.ascii"
	if flag.NArg() > 1 {
		ofname = flag.Arg(1)
	}

	start := time.Now()
	fmt.Printf("::: %s...\n", os.Args[0])

	// create a default fwk application, with some properties
	// extracted from the CLI
	app := job.New(job.P{
		"EvtMax":   *g_evtmax,
		"NProcs":   *g_nprocs,
		"MsgLevel": job.MsgLevel(*g_lvl),
	})

	r, err := os.Open(ifname)
	if err != nil {
		app.Errorf("could not open input file [%s]: %v\n", ifname, err)
		os.Exit(1)
	}
	defer r.Close()

	w, err := os.Create(ofname)
	if err != nil {
		app.Errorf("could not create output file [%s]: %v\n", ofname, err)
		os.Exit(1)
	}
	defer w.Close()

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
				R: r,
			},
		},
	})

	// create an output-stream
	app.Create(job.C{
		Type: "github.com/go-hep/fwk.OutputStream",
		Name: "output",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "t1-ints1-massaged",      // location of data to write out
					Type: reflect.TypeOf(int64(0)), // type of that data
				},
			},
			"Streamer": &testdata.OutputStream{
				W: w,
			},
		},
	})

	// run the application
	app.Run()

	fmt.Printf("::: %s... [done] (cpu=%v)\n", os.Args[0], time.Since(start))
}

/*
output:

$ ::: fwk-ex-tuto-3...
t2                   INFO configure...
t2                   INFO configure... [done]
t2                   INFO start...
app                  INFO >>> running evt=0...
t2                   INFO proc... (id=0|0) => [0 -> 0]
app                  INFO >>> running evt=1...
t2                   INFO proc... (id=1|0) => [1 -> 1]
app                  INFO >>> running evt=2...
t2                   INFO proc... (id=2|0) => [2 -> 4]
app                  INFO >>> running evt=3...
t2                   INFO proc... (id=3|0) => [3 -> 9]
app                  INFO >>> running evt=4...
t2                   INFO proc... (id=4|0) => [4 -> 16]
app                  INFO >>> running evt=5...
t2                   INFO proc... (id=5|0) => [5 -> 25]
app                  INFO >>> running evt=6...
t2                   INFO proc... (id=6|0) => [6 -> 36]
app                  INFO >>> running evt=7...
t2                   INFO proc... (id=7|0) => [7 -> 49]
app                  INFO >>> running evt=8...
t2                   INFO proc... (id=8|0) => [8 -> 64]
app                  INFO >>> running evt=9...
t2                   INFO proc... (id=9|0) => [9 -> 81]
app                  INFO >>> running evt=10...
t2                   INFO proc... (id=10|0) => [10 -> 100]
app                  INFO >>> running evt=11...
t2                   INFO proc... (id=11|0) => [11 -> 121]
app                  INFO >>> running evt=12...
t2                   INFO proc... (id=12|0) => [12 -> 144]
app                  INFO >>> running evt=13...
t2                   INFO proc... (id=13|0) => [13 -> 169]
app                  INFO >>> running evt=14...
t2                   INFO proc... (id=14|0) => [14 -> 196]
app                  INFO >>> running evt=15...
t2                   INFO proc... (id=15|0) => [15 -> 225]
app                  INFO >>> running evt=16...
t2                   INFO proc... (id=16|0) => [16 -> 256]
app                  INFO >>> running evt=17...
t2                   INFO proc... (id=17|0) => [17 -> 289]
app                  INFO >>> running evt=18...
t2                   INFO proc... (id=18|0) => [18 -> 324]
app                  INFO >>> running evt=19...
t2                   INFO proc... (id=19|0) => [19 -> 361]
app                  INFO >>> running evt=20...
t2                   INFO stop...
::: fwk-ex-tuto-3... [done] (cpu=2.127406ms)

$ cat output.ascii
0
1
4
9
16
25
36
49
64
81
100
121
144
169
196
225
256
289
324
361
*/
