// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist_test

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
)

func ExampleGraph() {
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tg")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rhist.Graph)
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
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tge")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rhist.GraphErrors)
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
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tgae")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rhist.GraphErrors)
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

func ExampleCreate_graph() {
	const fname = "graph_example.root"
	defer os.Remove(fname)

	f, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hg := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})

	fmt.Printf("original graph:\n")
	for i, pt := range hg.Points() {
		fmt.Printf("pt[%d]=%+v\n", i, pt)
	}

	rg := rhist.NewGraphFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rhist.Graph))
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

	f, err := groot.Create(fname)
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

	rg := rhist.NewGraphErrorsFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rhist.GraphErrors))
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

	f, err := groot.Create(fname)
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

	rg := rhist.NewGraphAsymmErrorsFrom(hg)

	err = f.Put("gr", rg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("gr")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.S2D(robj.(rhist.GraphErrors))
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
