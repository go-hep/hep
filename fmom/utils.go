// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"

	"gonum.org/v1/gonum/spatial/r3"
)

const (
	twopi = 2 * math.Pi
)

// DeltaPhi returns the delta Phi in range [-pi,pi[ from two P4
func DeltaPhi(p1, p2 P4) float64 {
	// FIXME: do something more efficient when p1&p2 are PxPyPzE
	dphi := p1.Phi() - p2.Phi()
	return -math.Remainder(dphi, twopi)
}

// DeltaEta returns the delta Eta between two P4
func DeltaEta(p1, p2 P4) float64 {
	// FIXME: do something more efficient when p1&p2 are PxPyPzE
	return p1.Eta() - p2.Eta()
}

// DeltaR returns the delta R between two P4
func DeltaR(p1, p2 P4) float64 {
	deta := DeltaEta(p1, p2)
	dphi := DeltaPhi(p1, p2)
	return math.Sqrt(deta*deta + dphi*dphi)
}

// Dot returns the dot product p1·p2.
func Dot(p1, p2 P4) float64 {
	return p1.E()*p2.E() - p1.Px()*p2.Px() - p1.Py()*p2.Py() - p1.Pz()*p2.Pz()
}

// CosTheta returns the cosine of the angle between the momentum of two 4-vectors.
func CosTheta(p1, p2 P4) float64 {
	mag1 := p1.P()
	mag2 := p2.P()
	dot := vecDot(vecOf(p1), vecOf(p2))
	cosTh := dot / (mag1 * mag2)
	return cosTh
}

// vecOf returns the space-components of a 4-vector.
func vecOf(p P4) r3.Vec {
	return r3.Vec{X: p.Px(), Y: p.Py(), Z: p.Pz()}
}
