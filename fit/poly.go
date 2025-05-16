// Copyright ©2025 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit // import "go-hep.org/x/hep/fit"

import (
	"fmt"
	"slices"

	"gonum.org/v1/gonum/mat"
)

// Poly fits a polynomial p of degree `degree` to points (x, y):
//
//	p(x) = p[0] * x^deg + ... + p[deg]
func Poly(xs, ys []float64, degree int) ([]float64, error) {
	var (
		a = vandermonde(xs, degree+1)
		b = mat.NewDense(len(ys), 1, ys)
		o = make([]float64, degree+1)
		c = mat.NewDense(degree+1, 1, o)
	)

	var qr mat.QR
	qr.Factorize(a)

	const trans = false
	err := qr.SolveTo(c, trans, b)
	if err != nil {
		return nil, fmt.Errorf("could not solve QR: %w", err)
	}

	slices.Reverse(o)
	return o, nil
}

func vandermonde(a []float64, d int) *mat.Dense {
	x := mat.NewDense(len(a), d, nil)
	for i := range a {
		for j, p := 0, 1.0; j < d; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}
