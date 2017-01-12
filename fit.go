// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fit provides functions to fit data.
package fit

import (
	"github.com/gonum/diff/fd"
	"github.com/gonum/optimize"
)

type Func1D struct {
	F  func(x float64, ps []float64) float64 // function to minimize
	N  int                                   // number of parameters
	Ps []float64                             // initial parameters values

	X   []float64
	Y   []float64
	Err []float64

	sig2 []float64

	fct  func(ps []float64) float64
	grad func(grad, ps []float64)
}

func (f *Func1D) init() {

	f.sig2 = make([]float64, len(f.Y))
	switch {
	default:
		for i := range f.Y {
			f.sig2[i] = 1
		}
	case f.Err != nil:
		for i, v := range f.Err {
			f.sig2[i] = 1 / (v * v)
		}
	}

	if f.Ps == nil {
		f.Ps = make([]float64, f.N)
	}

	if len(f.Ps) == 0 {
		panic("fit: invalid number of initial parameters")
	}

	if len(f.X) != len(f.Y) {
		panic("fit: mismatch length")
	}

	if len(f.sig2) != len(f.Y) {
		panic("fit: mismatch length")
	}

	f.fct = func(ps []float64) float64 {
		var chi2 float64
		for i := range f.X {
			res := f.F(f.X[i], ps) - f.Y[i]
			chi2 += res * res * f.sig2[i]
		}
		return 0.5 * chi2
	}

	f.grad = func(grad, ps []float64) {
		fd.Gradient(grad, f.fct, ps, nil)
	}
}

func Curve1D(f Func1D, settings *optimize.Settings, m optimize.Method) (*optimize.Result, error) {
	f.init()

	p := optimize.Problem{
		Func: f.fct,
		Grad: f.grad,
	}

	if m == nil {
		m = &optimize.NelderMead{}
	}

	p0 := make([]float64, len(f.Ps))
	copy(p0, f.Ps)
	return optimize.Local(p, p0, settings, m)
}
