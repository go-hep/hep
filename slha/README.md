slha
====

[![Build Status](https://secure.travis-ci.org/go-hep/slha.png)](http://travis-ci.org/go-hep/slha)
[![GoDoc](https://godoc.org/go-hep.org/x/hep/slha?status.svg)](https://godoc.org/go-hep.org/x/hep/slha)

Package `slha` implements encoding and decoding of SUSY Les Houches
Accords (SLHA) data format.

## Installation

```sh
$ go get go-hep.org/x/hep/slha
```

## Example

```go
package main

import (
	"fmt"
	"os"

	"go-hep.org/x/hep/slha"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fname := "testdata/sps1a.spc"
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	f, err := os.Open(fname)
	handle(err)

	defer f.Close()

	data, err := slha.Decode(f)
	handle(err)

	spinfo := data.Blocks.Get("SPINFO")
	value, err := spinfo.Get(1)
	handle(err)
	fmt.Printf("spinfo: %s -- %q\n", value.Interface(), value.Comment())

	modsel := data.Blocks.Get("MODSEL")
	value, err = modsel.Get(1)
	handle(err)
	fmt.Printf("modsel: %d -- %q\n", value.Interface(), value.Comment())

	mass := data.Blocks.Get("MASS")
	value, err = mass.Get(5)
	handle(err)
	fmt.Printf("mass[pdgid=5]: %v -- %q\n", value.Interface(), value.Comment())

	nmix := data.Blocks.Get("NMIX")
	value, err = nmix.Get(1, 2)
	handle(err)
	fmt.Printf("nmix[1,2] = %v -- %q\n", value.Interface(), value.Comment())
}


// Output:
// spinfo: SOFTSUSY -- "spectrum calculator"
// modsel: 1 -- "sugra"
// mass[pdgid=5]: 4.88991651 -- "b-quark pole mass calculated from mb(mb)_Msbar"
// nmix[1,2] = -0.0531103553 -- "N_12"
```

## Documentation

Documentation is available on [godoc](https://godoc.org/go-hep.org/x/hep/slha):

  https://godoc.org/go-hep.org/x/hep/slha

