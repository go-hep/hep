// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom_test

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/fmom"
	"gonum.org/v1/gonum/spatial/r3"
)

func Example() {
	p1 := fmom.NewPxPyPzE(10, 20, 30, 40)
	p2 := fmom.NewPtEtaPhiM(10, 2, math.Pi/2, 40)

	fmt.Printf("p1 = %v (m=%g)\n", p1, p1.M())
	fmt.Printf("p2 = %v\n", p2)

	p3 := fmom.Add(&p1, &p2)
	fmt.Printf("p3 = p1+p2 = %v\n", p3)

	p4 := fmom.Boost(&p1, r3.Vec{X: 0, Y: 0, Z: 0.99})
	fmt.Printf("p4 = boost(p1, (0,0,0.99)) = %v\n", p4)

	p5 := fmom.Boost(&p1, fmom.BoostOf(&p1).Scale(-1))
	fmt.Printf("p5 = rest-frame(p1) = %v\n", p5)

	// Output:
	// p1 = fmom.P4{Px:10, Py:20, Pz:30, E:40} (m=14.142135623730951)
	// p2 = fmom.P4{Pt:10, Eta:2, Phi:1.5707963267948966, M:40}
	// p3 = p1+p2 = fmom.P4{Px:10, Py:30, Pz:66.26860407847019, E:94.91276392425375}
	// p4 = boost(p1, (0,0,0.99)) = fmom.P4{Px:10, Py:20, Pz:493.3813186858014, E:494.0901998908098}
	// p5 = rest-frame(p1) = fmom.P4{Px:0, Py:0, Pz:0, E:14.14213562373095}
}
