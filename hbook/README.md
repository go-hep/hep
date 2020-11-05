hbook
=====

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

 https://godoc.org/go-hep.org/x/hep/hbook

## Example

### H1D

[embedmd]:# (h1d_example_test.go go /func ExampleH1D/ /\n}/)
```go
func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	// fill h with a slice of values and their weights
	h.FillN([]float64{1, 2, 3}, []float64{1, 1, 1})
	h.FillN([]float64{1, 2, 3}, nil) // all weights are 1.

	fmt.Printf("mean:    %v\n", h.XMean())
	fmt.Printf("rms:     %v\n", h.XRMS())
	fmt.Printf("std-dev: %v\n", h.XStdDev())
	fmt.Printf("std-err: %v\n", h.XStdErr())

	// Output:
	// mean:    0.005589967511734562
	// rms:     1.0062596231244403
	// std-dev: 1.0062943821322063
	// std-err: 0.010059926295994191
}
```

[embedmd]:# (h1d_example_test.go go /func ExampleAddH1D/ /\n}/)
```go
func ExampleAddH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsum := hbook.AddH1D(h1, h2)
	fmt.Printf("Under: %.1f +/- %.1f\n", hsum.Binning.Outflows[0].SumW(), math.Sqrt(hsum.Binning.Outflows[0].SumW2()))
	for i := 0; i < hsum.Len(); i++ {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, hsum.Binning.Bins[i].SumW(), math.Sqrt(hsum.Binning.Bins[i].SumW2()))
	}
	fmt.Printf("Over : %.1f +/- %.1f\n", hsum.Binning.Outflows[1].SumW(), math.Sqrt(hsum.Binning.Outflows[1].SumW2()))

	// Output:
	// Under: 1.7 +/- 1.2
	// Bin 0: 4.7 +/- 2.4
	// Bin 1: 1.8 +/- 1.3
	// Bin 2: 2.7 +/- 2.1
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: 2.3 +/- 1.6
	// Bin 5: 3.0 +/- 2.2
	// Over : 1.5 +/- 1.1
}
```
[embedmd]:# (h1d_example_test.go go /func ExampleAddScaledH1D/ /\n}/)
```go
func ExampleAddScaledH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsum := hbook.AddScaledH1D(h1, 10, h2)
	fmt.Printf("Under: %.1f +/- %.1f\n", hsum.Binning.Outflows[0].SumW(), math.Sqrt(hsum.Binning.Outflows[0].SumW2()))
	for i := 0; i < hsum.Len(); i++ {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, hsum.Binning.Bins[i].SumW(), math.Sqrt(hsum.Binning.Bins[i].SumW2()))
	}
	fmt.Printf("Over : %.1f +/- %.1f\n", hsum.Binning.Outflows[1].SumW(), math.Sqrt(hsum.Binning.Outflows[1].SumW2()))

	// Output:
	// Under: 8.0 +/- 7.1
	// Bin 0: 24.5 +/- 15.7
	// Bin 1: 9.0 +/- 8.1
	// Bin 2: 9.0 +/- 7.3
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: 14.0 +/- 13.0
	// Bin 5: 21.0 +/- 20.0
	// Over : 10.5 +/- 10.0
}
```
[embedmd]:# (h1d_example_test.go go /func ExampleSubH1D/ /\n}/)
```go
func ExampleSubH1D() {

	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := hbook.NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	hsub := hbook.SubH1D(h1, h2)
	under := hsub.Binning.Outflows[0]
	fmt.Printf("Under: %.1f +/- %.1f\n", under.SumW(), math.Sqrt(under.SumW2()))
	for i, bin := range hsub.Binning.Bins {
		fmt.Printf("Bin %v: %.1f +/- %.1f\n", i, bin.SumW(), math.Sqrt(bin.SumW2()))
	}
	over := hsub.Binning.Outflows[1]
	fmt.Printf("Over : %.1f +/- %.1f\n", over.SumW(), math.Sqrt(over.SumW2()))

	// Output:
	// Under: 0.3 +/- 1.2
	// Bin 0: 0.3 +/- 2.4
	// Bin 1: 0.2 +/- 1.3
	// Bin 2: 1.3 +/- 2.1
	// Bin 3: 0.0 +/- 0.0
	// Bin 4: -0.3 +/- 1.6
	// Bin 5: -1.0 +/- 2.2
	// Over : -0.5 +/- 1.1
}
```

### H2D

[embedmd]:# (h2d_example_test.go go /func ExampleH2D/ /\n}/)
```go
func ExampleH2D() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
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

	// fill h with slices of values and their weights
	h.FillN(
		[]float64{1, 2, 3}, // xs
		[]float64{1, 2, 3}, // ys
		[]float64{1, 1, 1}, // ws
	)

	// fill h with slices of values. all weights are 1.
	h.FillN(
		[]float64{1, 2, 3}, // xs
		[]float64{1, 2, 3}, // ys
		nil,                // ws
	)
}
```

### S2D

[embedmd]:# (s2d_example_test.go go /func ExampleS2D/ /\n}/)
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
