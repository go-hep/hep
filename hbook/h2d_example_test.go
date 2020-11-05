// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"log"

	"go-hep.org/x/hep/hbook"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
)

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
