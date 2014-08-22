hbook
=====

[![Build Status](https://drone.io/github.com/go-hep/hbook/status.png)](https://drone.io/github.com/go-hep/hbook/latest)
[![GoDoc](https://godoc.org/github.com/go-hep/hbook?status.svg)](https://godoc.org/github.com/go-hep/hbook)

`hbook` is a set of data analysis tools for HEP (histograms (1D, 2D, 3D), profiles and ntuples).

`hbook` is a work in progress of a concurrent friendly histogram filling toolkit.
It is loosely based on `AIDA` interfaces and concepts as well as the "simplicity" of `HBOOK` and the previous work of `YODA`.

## Installation

```sh
$ go get github.com/go-hep/hbook
```

## Documentation

Documentation is available on godoc:

 http://godoc.org/github.com/go-hep/hbook

## Example

```go
package main

import (
	   "math/rand"
	   "github.com/go-hep/hbook"
)

func main() {
	 h := hbook.NewH1D(100, 0, 100)
	 for i := 0; i < 100; i++ {
	 	 h.Fill(rand.Float64()*100, 1.)
	 }
}

```
