// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

// Resonance holds mass and width informations for a Breit-Wigner
// distribution about a given mass
type Resonance struct {
	Mass  Measurement // mass measurement
	Width Measurement // total width measurement
	Lower float64     // lower cutoff of allowed width values
	Upper float64     // upper cutoff of allowed width values
}

// Lifetime computes and returns the lifetime from the total width
func (r *Resonance) Lifetime() Measurement {
	// lifetime = hbar / totalwidth
	const hbar = 6.58211889e-25 // in GeV s
	var lt Measurement
	lt.Value = hbar / r.Width.Value
	lt.Sigma = lt.Value * r.Width.Sigma / r.Width.Value
	return lt
}

func (r *Resonance) SetTotalWidthFromLifetime(lifetime Measurement) {
	// totalwidth = hbar / lifetime
	const epsilon = 1.0e-20
	const hbar = 6.58211889e-25 // in GeV s
	var width float64
	var sigma float64

	// make no changes if lifetime is not greater than zero
	if lifetime.Value < epsilon {
		return
	}

	width = hbar / lifetime.Value

	if lifetime.Sigma < epsilon {
		sigma = 0.0
	} else {
		sigma = (lifetime.Sigma / lifetime.Value) * width
	}
	r.Width = Measurement{
		Value: width,
		Sigma: sigma,
	}
}
