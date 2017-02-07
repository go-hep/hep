// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit

import (
	"github.com/go-hep/hbook"
	"github.com/gonum/optimize"
)

// H1D returns the fit of histogram h with function f and optimization method m.
//
// Only bins with at least an entry are considered for the fit.
// In case settings is nil, the optimize.DefaultSettings is used.
// In case m is nil, the same default optimization method than for Curve1D is used.
func H1D(h *hbook.H1D, f Func1D, settings *optimize.Settings, m optimize.Method) (*optimize.Result, error) {
	var (
		n     = h.Len()
		xdata = make([]float64, 0, n)
		ydata = make([]float64, 0, n)
		yerrs = make([]float64, 0, n)
		bins  = h.Binning().Bins()
	)

	for _, bin := range bins {
		if bin.Entries() <= 0 {
			continue
		}
		xdata = append(xdata, bin.XMid())
		ydata = append(ydata, bin.SumW())
		yerrs = append(yerrs, bin.ErrW())
	}

	f.X = xdata
	f.Y = ydata
	f.Err = yerrs

	return Curve1D(f, settings, m)
}
