package fads

import (
	"math"

	"github.com/go-hep/fmom"
)

type Int64Slice []int64

func (p Int64Slice) Len() int {
	return len(p)
}

func (p Int64Slice) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p Int64Slice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type ByPt []Candidate

func (p ByPt) Len() int {
	return len(p)
}

func (p ByPt) Less(i, j int) bool {
	return p[i].Mom.Pt() < p[j].Mom.Pt()
}

func (p ByPt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func newPtEtaPhiE(pt, eta, phi, ene float64) fmom.PxPyPzE {
	pt = math.Abs(pt)

	px := pt * math.Cos(phi)
	py := pt * math.Sin(phi)
	pz := pt * math.Sinh(eta)

	return fmom.PxPyPzE{px, py, pz, ene}
}
