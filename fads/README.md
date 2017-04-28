fads
====

[![GoDoc](https://godoc.org/go-hep.org/x/hep/fads?status.svg)](https://godoc.org/go-hep.org/x/hep/fads)

`fads`, a FAst Detector Simulation, is a Go-based detector simulation including a tracking system embedded into a magnetic field, calorimeters and a muon system.

## Installation

Is done via `go get`:

```sh
$ go get go-hep.org/x/hep/fads/...
```

## Documentation

Is available on `godoc`: https://godoc.org/go-hep.org/x/hep/fads

## Example

A test application is available over there:

https://go-hep.org/x/hep/fads/blob/master/cmd/fads-app/main.go

A more in-depth tutorial is available at [go-hep/tutos](https://github.com/go-hep/tutos) but, in a nutshell:

```sh
$ fads-app -help
Usage: fads-app [options] <hepmc-input-file>

ex:
 $ fads-app -l=INFO -evtmax=-1 ./testdata/hepmc.data

options:
  -cpu-prof=false: enable CPU profiling
  -evtmax=-1: number of events to process
  -l="INFO": log level (DEBUG|INFO|WARN|ERROR)
  -nprocs=0: number of concurrent events to process

$ fads-app $GOPATH/src/go-hep.org/x/hep/fads/testdata/hepmc.data
::: fads-app...
app                  INFO >>> running evt=0...
app                  INFO >>> running evt=1...
app                  INFO >>> running evt=2...
app                  INFO >>> running evt=3...
app                  INFO >>> running evt=4...
app                  INFO >>> running evt=5...
app                  INFO cpu: 1.212611252s
app                  INFO mem: alloc:           3219 kB
app                  INFO mem: tot-alloc:      26804 kB
app                  INFO mem: n-mallocs:      53058
app                  INFO mem: n-frees:        52419
app                  INFO mem: gc-pauses:         36 ms
::: fads-app... [done] (time=1.216341021s)
```

## Using the execution tracer

It is now possible to run `fads-app` with the `runtime/trace` execution tracer
enabled:

```sh
$ fads-app -trace=foo.out ./testdata/hepmc.data
```

The traces are then analyzed and displayed via `trace-viewer` from the
`catapult` project:

```sh
$ go tool trace `which fads-app` ./foo.out
```

See: [catapult-tracing](https://github.com/catapult-project/catapult/tree/master/tracing#readme)
