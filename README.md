dao
===

`dao` is a set of data analysis objects for HEP (histograms (1D, 2D, 3D), profiles and ntuples).

`dao` is a work in progress of a concurrent friendly histogram filling toolkit.
It is loosely based on `AIDA` interfaces and concepts as well as the "simplicity" of `HBOOK` and the previous work of `YODA`.

## Installation

```sh
$ go get github.com/go-hep/dao
```

## Documentation

Documentation is available on godoc:

 http://godoc.org/github.com/go-hep/dao

## Example

```go
package main

import (
	   "math/rand"
	   "github.com/go-hep/dao"
)

func main() {
	 h := dao.NewH1D(100, 0, 100)
	 for i := 0; i < 100; i++ {
	 	 h.Fill(rand.Float64()*100, 1.)
	 }
}

```
