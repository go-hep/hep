fwk
===

[![GoDoc](https://godoc.org/go-hep.org/x/hep/fwk?status.svg)](https://godoc.org/go-hep.org/x/hep/fwk)

`fwk` is a HEP oriented concurrent framework written in `Go`.
`fwk` should be easy to pick up and use for small and fast analyses but should also support reconstruction, simulation, ... use cases.

## Installation

`fwk`, like any pure-Go package, is `go get` able:

```sh
$ go get go-hep.org/x/hep/fwk/...
```

(yes, with the ellipsis after the slash, to install all the "sub-packages")


## Documentation

The documentation is available on `godoc`:

 https://godoc.org/go-hep.org/x/hep/fwk


## Examples


### `fwk` tuto examples

The [examples](https://github.com/go-hep/hep/blob/main/fwk/examples)
directory contains a few simple applications which exercize the `fwk`
toolkit.

The examples/tutorials should be readily available as soon as you've
executed:

```sh
$ go get go-hep.org/x/hep/fwk/examples/...
```

*e.g.:*

```sh
$ fwk-ex-tuto-1 -help
Usage: fwk-ex-tuto1 [options]

ex:
 $ fwk-ex-tuto-1 -l=INFO -evtmax=-1

options:
  -evtmax=10: number of events to process
  -l="INFO": message level (DEBUG|INFO|WARN|ERROR)
  -nprocs=0: number of events to process concurrently
```

```sh
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
app                  INFO cpu: 432.039us
app                  INFO mem: alloc:             68 kB
app                  INFO mem: tot-alloc:         79 kB
app                  INFO mem: n-mallocs:        410
app                  INFO mem: n-frees:           60
app                  INFO mem: gc-pauses:          0 ms
::: fwk-ex-tuto-1... [done] (cpu=625.918us)
```

### Physics-oriented demonstrator

There is also a more physics-oriented example/demonstrator: [fads](https://go-hep.org/x/hep/fads)


## Tools

### `fwk-new-comp`

`fwk-new-comp` is a small tool to generate most of the boilerplate
code to bootstrap the creation of new `fwk.Component`s (either
`fwk.Task` or `fwk.Svc`)

```sh
$ fwk-new-comp -help
Usage: fwk-new-comp [options] <component-name>

ex:
 $ fwk-new-comp -c=task -p=mypackage mytask
 $ fwk-new-comp -c=task -p mypackage mytask >| mytask.go
 $ fwk-new-comp -c=svc  -p mypackage mysvc  >| mysvc.go

options:
  -c="task": type of component to generate (task|svc)
  -p="": name of the package holding the component
```


### `fwk-list-components`

`fwk-list-components` lists all the currently available components.

```sh
$ fwk-list-components
::: components... (5)
[0000/0005] go-hep.org/x/hep/fwk.InputStream
[0001/0005] go-hep.org/x/hep/fwk.OutputStream
[0002/0005] go-hep.org/x/hep/fwk.appmgr
[0003/0005] go-hep.org/x/hep/fwk.datastore
[0004/0005] go-hep.org/x/hep/fwk.dflowsvc
```
