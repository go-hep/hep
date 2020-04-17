// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"fmt"
	"math"
)

// DivideH1D divides 2 1D-histograms and returns a 2D scatter.
// DivideH1D returns an error if the binning of the 1D histograms are not compatible.
func DivideH1D(num, den *H1D) (*S2D, error) {
	var s2d S2D

	bins1 := num.Binning.Bins
	bins2 := den.Binning.Bins

	for i := range bins1 {
		b1 := bins1[i]
		b2 := bins2[i]

		if !fuzzyEq(b1.XMin(), b2.XMin()) || !fuzzyEq(b1.XMax(), b2.XMax()) {
			return nil, fmt.Errorf("hbook: x binnings are not equivalent in %v / %v", num.Name(), den.Name())
		}

		// assemble the x value and error
		// use the midpoint of the "bin" for the new central value
		x := b1.XMid()
		exm := x - b1.XMin()
		exp := b1.XMax() - x

		// assemble the y value and error
		// TODO(sbinet): provide optional alternative behaviours to fill with NaN
		//               or remove the invalid points
		var y, ey float64
		b2h := b2.SumW() / b2.XWidth() // height of the bin
		b1h := b1.SumW() / b1.XWidth() // ditto
		b2herr := math.Sqrt(b2.SumW2()) / b2.XWidth()
		b1herr := math.Sqrt(b1.SumW2()) / b1.XWidth()

		switch {
		case b2h == 0 || (b1h == 0 && b1herr != 0): // TODO(sbinet): is it OK?
			y = math.NaN()
			ey = math.NaN()
		default:
			y = b1h / b2h
			// TODO(sbinet): is this the exact error treatment for all (uncorrelated) cases?
			// What should be the behaviour around 0? +1 and -1 fills?
			relerr1 := 0.0
			if b1herr != 0 {
				relerr1 = math.Sqrt(b1.SumW2()) / b1.SumW() // TODO(sbinet) refactor as bin1d.RelErr() ?
			}
			relerr2 := 0.0
			if b2herr != 0 {
				relerr2 = math.Sqrt(b2.SumW2()) / b2.SumW()
			}
			ey = y * math.Sqrt(relerr1*relerr1+relerr2*relerr2)
		}

		// deal with +/- errors separately, inverted for the denominator contributions:
		// TODO(sbinet): check correctness with different signed numerator and denominator.

		s2d.Fill(Point2D{X: x, Y: y, ErrX: Range{Min: exm, Max: exp}, ErrY: Range{Min: ey, Max: ey}})
	}
	return &s2d, nil
}

// fuzzyEq returns true if a and b are equal with a degree of fuzziness
func fuzzyEq(a, b float64) bool {
	const tol = 1e-5
	aa := math.Abs(a)
	bb := math.Abs(b)
	absavg := 0.5 * (aa + bb)
	absdiff := math.Abs(a - b)
	return (aa < 1e-8 && bb < 1e-8) || absdiff < tol*absavg
}
