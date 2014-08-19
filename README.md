fwk
===

[![Build Status](https://drone.io/github.com/go-hep/fwk/status.png)](https://drone.io/github.com/go-hep/fwk/latest)

`fwk` is a HEP oriented concurrent framework written in `Go`.
`fwk` should be easy to pick up and use for small and fast analyses but should also support reconstruction, simulation, ... use cases.

## Installation

`fwk`, like any pure-Go package, is `go get` able:

```sh
$ go get github.com/go-hep/fwk/...
```

(yes, with the ellipsis after the slash, to install all the "sub-packages")


## Documentation

The documentation is available on `godoc`:

 http://godoc.org/github.com/go-hep/fwk


## Examples

The [examples](https://github.com/go-hep/fwk/blob/master/examples)
directory contains a few simple applications which exercize the `fwk`
toolkit.

There is also a more physics-oriented example/demonstrator: [fads](https://github.com/go-hep/fads)


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
[0000/0005] github.com/go-hep/fwk.InputStream
[0001/0005] github.com/go-hep/fwk.OutputStream
[0002/0005] github.com/go-hep/fwk.appmgr
[0003/0005] github.com/go-hep/fwk.datastore
[0004/0005] github.com/go-hep/fwk.dflowsvc
```
