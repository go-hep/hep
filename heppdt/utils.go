// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

func calcWidthFromLifetime(lt float64) float64 {
	// totalwidth = hbar / lifetime
	const epsilon = 1.0e-20
	const hbar = 6.58211889e-25 // in GeV s
	if lt < epsilon {
		return 0.
	}
	return hbar / lt
}
