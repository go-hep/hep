// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

import (
	"math"
)

// Particle holds informations on a particle as per the PDG booklet
type Particle struct {
	ID          PID           // particle ID
	Name        string        // particle name
	PDG         int           // PDG code of the particle
	Mass        float64       // particle mass in GeV
	Charge      float64       // electrical charge
	ColorCharge float64       // color charge
	Spin        SpinState     // spin state
	Quarks      []Constituent // constituents
	Resonance   Resonance     // resonance
}

// IsStable returns whether this particle is stable
func (p *Particle) IsStable() bool {
	res := &p.Resonance
	if res.Width.Value == -1. {
		return false
	}
	lt := res.Lifetime()
	if res.Width.Value > 0 || lt.Value > 0 {
		// FIXME(sbinet): res.Width.Value should be == -1.
		// when lifetime.Value == +inf
		if math.IsInf(lt.Value, +1) {
			return true
		}
		return false
	}
	return true
}
