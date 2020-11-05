// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/hbook"
)

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

func ExampleS2D_newS2DFrom() {
	s := hbook.NewS2DFrom([]float64{1, 2, -1}, []float64{1, 1.5, 2})
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

func ExampleS2D_newS2DFromH1D() {
	h := hbook.NewH1D(20, -4, +4)
	h.Fill(1, 2)
	h.Fill(2, 3)
	h.Fill(3, 1)
	h.Fill(1, 1)
	h.Fill(-2, 1)
	h.Fill(-3, 1)

	s := hbook.NewS2DFromH1D(h)
	s.Sort()
	for _, pt := range s.Points() {
		fmt.Printf("point=(%+3.2f +/- (%+3.2f,%+3.2f), %+3.2f +/- (%+3.2f, %+3.2f))\n", pt.X, pt.ErrX.Min, pt.ErrX.Max, pt.Y, pt.ErrY.Min, pt.ErrY.Max)
	}

	// Output:
	// point=(-3.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-3.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-3.00 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(-2.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-2.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-1.80 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(-1.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-1.00 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-0.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-0.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+0.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+0.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+1.00 +/- (+0.20,+0.20), +7.50 +/- (+5.59, +5.59))
	// point=(+1.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+1.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+2.20 +/- (+0.20,+0.20), +7.50 +/- (+7.50, +7.50))
	// point=(+2.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+3.00 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(+3.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+3.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
}
