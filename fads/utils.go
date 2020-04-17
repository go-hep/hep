// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"math"

	"go-hep.org/x/hep/fmom"
)

type int64Slice []int64

func (p int64Slice) Len() int {
	return len(p)
}

func (p int64Slice) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p int64Slice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// ByPt sorts candidate by descending Pt
type ByPt []Candidate

func (p ByPt) Len() int {
	return len(p)
}

func (p ByPt) Less(i, j int) bool {
	return p[i].Mom.Pt() > p[j].Mom.Pt()
}

func (p ByPt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func newPtEtaPhiE(pt, eta, phi, ene float64) fmom.PxPyPzE {
	pt = math.Abs(pt)

	px := pt * math.Cos(phi)
	py := pt * math.Sin(phi)
	pz := pt * math.Sinh(eta)

	return fmom.NewPxPyPzE(px, py, pz, ene)
}
