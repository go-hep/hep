hbook
=====

[![Build Status](https://secure.travis-ci.org/go-hep/hbook.png)](http://travis-ci.org/go-hep/hbook)
[![GoDoc](https://godoc.org/go-hep.org/x/hep/hbook?status.svg)](https://godoc.org/go-hep.org/x/hep/hbook)

`hbook` is a set of data analysis tools for HEP (histograms (1D, 2D, 3D), profiles and ntuples).

`hbook` is a work in progress of a concurrent friendly histogram filling toolkit.
It is loosely based on `AIDA` interfaces and concepts as well as the "simplicity" of `HBOOK` and the previous work of `YODA`.

## Installation

```sh
$ go get go-hep.org/x/hep/hbook
```

## Documentation

Documentation is available on godoc:

 http://godoc.org/go-hep.org/x/hep/hbook

## Example

### H1D

[embedmd]:# (h1d_test.go go /func ExampleH1D/ /\n}/)
```go
func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:     0,
		Sigma:  1,
		Source: rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}

	fmt.Printf("mean:    %v\n", h.XMean())
	fmt.Printf("rms:     %v\n", h.XRMS())
	fmt.Printf("std-dev: %v\n", h.XStdDev())
	fmt.Printf("std-err: %v\n", h.XStdErr())

	// Output:
	// mean:    -0.017341752512581167
	// rms:     0.9913281479386786
	// std-dev: 0.9912260153046148
	// std-err: 0.00991226015304615
}
```

### H2D

[embedmd]:# (h2d_test.go go /func ExampleH2D/ /\n}/)
```go
func ExampleH2D() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		h.Fill(v[0], v[1], 1)
	}
}
```

### S2D

[embedmd]:# (s2d_test.go go /func ExampleS2D/ /\n}/)
```go
func ExampleS2D() {
	s := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		log.Fatal("nil pointer to S2D")
	}

	fmt.Printf("len=%d\n", s.Len())

	s.Fill(hbook.Point2D{X: 10, Y: -10, ErrX: hbook.Range{Min: 5, Max: 5}, ErrY: hbook.Range{Min: 6, Max: 6}})
	fmt.Printf("len=%d\n", s.Len())
	fmt.Printf("pt[%d]=%+v\n", 3, s.Point(3))

	// Output:
	// len=3
	// len=4
	// pt[3]={X:10 Y:-10 ErrX:{Min:5 Max:5} ErrY:{Min:6 Max:6}}
}
```

### Ntuple

#### Open an existing n-tuple

```go
package main

import (
	"database/sql"
	"fmt"

	_ "go-hep.org/x/hep/csvutil/csvdriver"
	"go-hep.org/x/hep/hbook/ntup"
)

func main() {
	db, err := sql.Open("csv", "data.csv")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	nt, err := ntup.Open(db, "csv")
	if err != nil {
		panic(err)
	}

	h1, err := nt.ScanH1D("px where pt>100", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("h1: %v\n", h1)

	h2 := hbook.NewH1D(100, -10, 10)
	h2, err = nt.ScanH1D("px where pt>100 && pt < 1000", h2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("h2: %v\n", h2)

	h11 := hbook.NewH1D(100, -10, 10)
	h22 := hbook.NewH1D(100, -10, 10)
	err = nt.Scan("px, py where pt>100", func(px, py float64) error {
		h11.Fill(px, 1)
		h22.Fill(py, 1)
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("h11: %v\n", h11)
	fmt.Printf("h22: %v\n", h22)
}
```
