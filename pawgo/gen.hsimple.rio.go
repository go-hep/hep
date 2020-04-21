// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/rio"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
	"gonum.org/v1/gonum/stat/distuv"
)

func main() {
	h1 := genH1()
	h2 := genH2()
	s2 := genS2()
	p1 := genP1()

	f, err := os.Create("testdata/hsimple.rio")
	if err != nil {
		log.Fatalf("could not create file: %+v", err)
	}
	defer f.Close()

	w, err := rio.NewWriter(f)
	if err != nil {
		log.Fatalf("could not create rio writer: %+v", err)
	}

	for _, v := range []struct {
		name  string
		value interface{}
	}{
		{"h1", h1},
		{"h2", h2},
		{"s2", s2},
		{"p1", p1},
	} {
		err = w.WriteValue(v.name, v.value)
		if err != nil {
			log.Fatalf("could not write %q: %+v", v.name, err)
		}
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("could not close rio writer: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close file: %+v", err)
	}
}

func genH1() *hbook.H1D {
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

	return h
}

func genH2() *hbook.H2D {
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

	return h
}

func genS2() *hbook.S2D {
	s := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		log.Fatal("nil pointer to S2D")
	}

	s.Fill(hbook.Point2D{X: 10, Y: -10, ErrX: hbook.Range{Min: 5, Max: 5}, ErrY: hbook.Range{Min: 6, Max: 6}})

	return s
}

func genP1() *hbook.P1D {
	const npoints = 1000

	p := hbook.NewP1D(100, -10, 10)
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
		p.Fill(v[0], v[1], 1)
	}

	return p
}
