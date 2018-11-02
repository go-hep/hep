// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// disable test on windows because of symlinks
// +build !windows

package rootio_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/rootio"
)

func ExampleGraph() {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tg")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rootio.Graph)
	fmt.Printf("name:  %q\n", g.Name())
	fmt.Printf("title: %q\n", g.Title())
	fmt.Printf("#pts:  %d\n", g.Len())
	for i := 0; i < g.Len(); i++ {
		x, y := g.XY(i)
		fmt.Printf("(x,y)[%d] = (%+e, %+e)\n", i, x, y)
	}

	// Output:
	// name:  "tg"
	// title: "graph without errors"
	// #pts:  4
	// (x,y)[0] = (+1.000000e+00, +2.000000e+00)
	// (x,y)[1] = (+2.000000e+00, +4.000000e+00)
	// (x,y)[2] = (+3.000000e+00, +6.000000e+00)
	// (x,y)[3] = (+4.000000e+00, +8.000000e+00)
}

func ExampleGraphErrors() {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tge")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rootio.GraphErrors)
	fmt.Printf("name:  %q\n", g.Name())
	fmt.Printf("title: %q\n", g.Title())
	fmt.Printf("#pts:  %d\n", g.Len())
	for i := 0; i < g.Len(); i++ {
		x, y := g.XY(i)
		xlo, xhi := g.XError(i)
		ylo, yhi := g.YError(i)
		fmt.Printf("(x,y)[%d] = (%+e +/- [%+e, %+e], %+e +/- [%+e, %+e])\n", i, x, xlo, xhi, y, ylo, yhi)
	}

	// Output:
	// name:  "tge"
	// title: "graph with errors"
	// #pts:  4
	// (x,y)[0] = (+1.000000e+00 +/- [+1.000000e-01, +1.000000e-01], +2.000000e+00 +/- [+2.000000e-01, +2.000000e-01])
	// (x,y)[1] = (+2.000000e+00 +/- [+2.000000e-01, +2.000000e-01], +4.000000e+00 +/- [+4.000000e-01, +4.000000e-01])
	// (x,y)[2] = (+3.000000e+00 +/- [+3.000000e-01, +3.000000e-01], +6.000000e+00 +/- [+6.000000e-01, +6.000000e-01])
	// (x,y)[3] = (+4.000000e+00 +/- [+4.000000e-01, +4.000000e-01], +8.000000e+00 +/- [+8.000000e-01, +8.000000e-01])
}

func ExampleGraphErrors_asymmErrors() {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tgae")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rootio.GraphErrors)
	fmt.Printf("name:  %q\n", g.Name())
	fmt.Printf("title: %q\n", g.Title())
	fmt.Printf("#pts:  %d\n", g.Len())
	for i := 0; i < g.Len(); i++ {
		x, y := g.XY(i)
		xlo, xhi := g.XError(i)
		ylo, yhi := g.YError(i)
		fmt.Printf("(x,y)[%d] = (%+e +/- [%+e, %+e], %+e +/- [%+e, %+e])\n", i, x, xlo, xhi, y, ylo, yhi)
	}

	// Output:
	// name:  "tgae"
	// title: "graph with asymmetric errors"
	// #pts:  4
	// (x,y)[0] = (+1.000000e+00 +/- [+1.000000e-01, +2.000000e-01], +2.000000e+00 +/- [+3.000000e-01, +4.000000e-01])
	// (x,y)[1] = (+2.000000e+00 +/- [+2.000000e-01, +4.000000e-01], +4.000000e+00 +/- [+6.000000e-01, +8.000000e-01])
	// (x,y)[2] = (+3.000000e+00 +/- [+3.000000e-01, +6.000000e-01], +6.000000e+00 +/- [+9.000000e-01, +1.200000e+00])
	// (x,y)[3] = (+4.000000e+00 +/- [+4.000000e-01, +8.000000e-01], +8.000000e+00 +/- [+1.200000e+00, +1.600000e+00])
}

func TestGraph(t *testing.T) {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tg")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rootio.Graph)
	if !ok {
		t.Fatalf("'tg' not a rootio.Graph: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
	}
}

func TestGraphErrors(t *testing.T) {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tge")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rootio.GraphErrors)
	if !ok {
		t.Fatalf("'tge' not a rootio.GraphErrors: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
		xlo, xhi := g.XError(i)
		if want := 0.1 * v; want != xlo {
			t.Errorf("xerr[%d].low=%v want=%v", i, xlo, want)
		}
		if want := 0.1 * v; want != xhi {
			t.Errorf("xerr[%d].high=%v want=%v", i, xhi, want)
		}
		ylo, yhi := g.YError(i)
		if want := 0.2 * v; want != ylo {
			t.Errorf("yerr[%d].low=%v want=%v", i, ylo, want)
		}
		if want := 0.2 * v; want != yhi {
			t.Errorf("yerr[%d].high=%v want=%v", i, yhi, want)
		}
	}
}

func TestGraphAsymmErrors(t *testing.T) {
	f, err := rootio.Open("testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tgae")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rootio.GraphErrors)
	if !ok {
		t.Fatalf("'tgae' not a rootio.GraphErrors: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
		xlo, xhi := g.XError(i)
		if want := 0.1 * v; want != xlo {
			t.Errorf("xerr[%d].low=%v want=%v", i, xlo, want)
		}
		if want := 0.2 * v; want != xhi {
			t.Errorf("xerr[%d].high=%v want=%v", i, xhi, want)
		}
		ylo, yhi := g.YError(i)
		if want := 0.3 * v; want != ylo {
			t.Errorf("yerr[%d].low=%v want=%v", i, ylo, want)
		}
		if want := 0.4 * v; want != yhi {
			t.Errorf("yerr[%d].high=%v want=%v", i, yhi, want)
		}
	}
}

func ExampleCreate_graph() {
	const fname = "graph_example.root"
	defer os.Remove(fname)

	f, err := rootio.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hg := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})

	fmt.Printf("original graph:\n")
	for i, pt := range hg.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	rg := rootio.NewGraphFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := rootio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rootio.Graph))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ngraph read back:\n")
	for i, pt := range hr.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	// Output:
	// original graph:
	// pt[0]={X:1 Y:1 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
	//
	// graph read back:
	// pt[0]={X:1 Y:1 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:0 Max:0} ErrY:{Min:0 Max:0}}
}

func ExampleCreate_graphErrors() {
	const fname = "grapherr_example.root"
	defer os.Remove(fname)

	f, err := rootio.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hg := hbook.NewS2D(
		hbook.Point2D{X: 1, Y: 1, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
		hbook.Point2D{X: 2, Y: 1.5, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
		hbook.Point2D{X: -1, Y: +2, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
	)

	fmt.Printf("original graph:\n")
	for i, pt := range hg.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	rg := rootio.NewGraphErrorsFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := rootio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rootio.GraphErrors))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ngraph read back:\n")
	for i, pt := range hr.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	// Output:
	// original graph:
	// pt[0]={X:1 Y:1 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
	//
	// graph read back:
	// pt[0]={X:1 Y:1 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:2 Max:2} ErrY:{Min:3 Max:3}}
}

func ExampleCreate_graphAsymmErrors() {
	const fname = "graphasymmerr_example.root"
	defer os.Remove(fname)

	f, err := rootio.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hg := hbook.NewS2D(
		hbook.Point2D{X: 1, Y: 1, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
		hbook.Point2D{X: 2, Y: 1.5, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
		hbook.Point2D{X: -1, Y: +2, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
	)

	fmt.Printf("original graph:\n")
	for i, pt := range hg.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	rg := rootio.NewGraphAsymmErrorsFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := rootio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rootio.GraphErrors))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\ngraph read back:\n")
	for i, pt := range hr.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	// Output:
	// original graph:
	// pt[0]={X:1 Y:1 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
	//
	// graph read back:
	// pt[0]={X:1 Y:1 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
	// pt[1]={X:2 Y:1.5 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
	// pt[2]={X:-1 Y:2 ErrX:{Min:1 Max:2} ErrY:{Min:3 Max:4}}
}
