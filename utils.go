package fmom

import (
	"math"
)

const (
	twopi = 2 * math.Pi
)

// DeltaPhi returns the delta Phi in range [-pi,pi[ from two P4
func DeltaPhi(p1, p2 P4) float64 {
	// FIXME: do something more efficient when p1&p2 are PxPyPzE
	dphi := -p1.Phi() + p2.Phi()
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
