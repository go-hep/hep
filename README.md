fads
====

[![Build Status](https://drone.io/github.com/go-hep/fads/status.png)](https://drone.io/github.com/go-hep/fads/latest)

`fads`, a FAst Detector Simulation, is a Go-based simulation including a tracking system embedded into a magnetic field, calorimeters and a muon system.

## Installation

Is done via `go get`:

```sh
$ go get github.com/go-hep/fads/...
```

## Documentation

Is available on `godoc`: http://godoc.org/github.com/go-hep/fads

## Example

A test application is available over there:

https://github.com/go-hep/fads/blob/master/cmd/fads-app/main.go

- Install it like so:

```sh
$ go get github.com/go-hep/fads/cmd/fads-app
```

- Run it like so:

```sh
$ fads-app ./go-hep/fads/testdata/hepmc.data
```

- help:

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
```
