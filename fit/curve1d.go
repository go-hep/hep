// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"gonum.org/v1/gonum/optimize"
)

// Curve1D returns the result of a non-linear least squares to fit
// a function f to the underlying data with method m.
func Curve1D(f Func1D, settings *optimize.Settings, m optimize.GlobalMethod) (*optimize.Result, error) {
	f.init()

	p := optimize.Problem{
		Func: f.fct,
		Grad: f.grad,
		Hess: f.hess,
	}

	if m == nil {
		m = &optimize.NelderMead{}
	}

	p0 := make([]float64, len(f.Ps))
	copy(p0, f.Ps)
	if settings == nil {
		settings = optimize.DefaultSettings()
	}
	settings.InitX = p0
	return optimize.Global(p, len(p0), settings, m)
}
