// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"log"
	"math/rand/v2"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmv"
)

func ExampleP1D() {
	const npoints = 1000

	p := hbook.NewP1D(100, -10, 10)
	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewPCG(1234, 1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for range npoints {
		v = dist.Rand(v)
		p.Fill(v[0], v[1], 1)
	}

	fmt.Printf("mean:    %v\n", p.XMean())
	fmt.Printf("rms:     %v\n", p.XRMS())
	fmt.Printf("std-dev: %v\n", p.XStdDev())
	fmt.Printf("std-err: %v\n", p.XStdErr())

	// Output:
	// mean:    0.0965063103552738
	// rms:     2.0024868956390707
	// std-dev: 2.0011608991313086
	// std-err: 0.06328226405725404
}
