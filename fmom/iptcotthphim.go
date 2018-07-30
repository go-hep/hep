// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"
)

type IPtCotThPhiM [4]float64

func NewIPtCotThPhiM(pt, eta, phi, m float64) IPtCotThPhiM {
	return IPtCotThPhiM([4]float64{pt, eta, phi, m})
}

func (p4 *IPtCotThPhiM) Clone() P4 {
	pp := *p4
	return &pp
}

func (p4 *IPtCotThPhiM) IPt() float64 {
	return p4[0]
}

func (p4 *IPtCotThPhiM) CotTh() float64 {
	return p4[1]
}

func (p4 *IPtCotThPhiM) Phi() float64 {
	return p4[2]
}

func (p4 *IPtCotThPhiM) M() float64 {
	return p4[3]
}

func (p4 *IPtCotThPhiM) Pt() float64 {
	ipt := p4.IPt()
	return 1 / ipt
}

func (p4 *IPtCotThPhiM) P() float64 {
	cotth := p4.CotTh()
	ipt := p4.IPt()
	return math.Sqrt(1+cotth*cotth) / ipt
}

func (p4 *IPtCotThPhiM) P2() float64 {
	cotth := p4.CotTh()
	ipt := p4.IPt()
	return (1 + cotth*cotth) / (ipt * ipt)
}

func (p4 *IPtCotThPhiM) M2() float64 {
	m := p4.M()
	return m * m
}

func (p4 *IPtCotThPhiM) TanTh() float64 {
	cotth := p4.CotTh()
	return 1 / cotth
}

func (p4 *IPtCotThPhiM) SinTh() float64 {
	cotth := p4.CotTh()
	return 1 / math.Sqrt(1+cotth*cotth)
}

func (p4 *IPtCotThPhiM) CosTh() float64 {
	cotth := p4.CotTh()
	cotth2 := cotth * cotth
	costh := math.Sqrt(cotth2 / (1 + cotth2))
	sign := 1.0
	if cotth < 0 {
		sign = -1.0
	}
	return sign * costh
}

func (p4 *IPtCotThPhiM) E() float64 {
	m := p4.M()
	p := p4.P()
	if m == 0 {
		return p
	}
	return math.Sqrt(p*p + m*m)
}

func (p4 *IPtCotThPhiM) Et() float64 {
	cotth := p4.CotTh()
	e := p4.E()
	return e / math.Sqrt(1+cotth*cotth)
}

func (p4 *IPtCotThPhiM) Eta() float64 {
	cotth := p4.CotTh()
	aux := math.Sqrt(1 + cotth*cotth)
	return -0.5 * math.Log((aux-cotth)/(aux+cotth))
}

func (p4 *IPtCotThPhiM) Rapidity() float64 {
	e := p4.E()
	pz := p4.Pz()
	return 0.5 * math.Log((e+pz)/(e-pz))
}

func (p4 *IPtCotThPhiM) Px() float64 {
	cosphi := p4.CosPhi()
	ipt := p4.IPt()
	pt := 1 / ipt
	return pt * cosphi
}

func (p4 *IPtCotThPhiM) Py() float64 {
	sinphi := p4.SinPhi()
	ipt := p4.IPt()
	pt := 1 / ipt
	return pt * sinphi
}

func (p4 *IPtCotThPhiM) Pz() float64 {
	cotth := p4.CotTh()
	ipt := p4.IPt()
	pt := 1 / ipt
	return pt * cotth
}

func (p4 *IPtCotThPhiM) CosPhi() float64 {
	phi := p4.Phi()
	return math.Cos(phi)
}

func (p4 *IPtCotThPhiM) SinPhi() float64 {
	phi := p4.Phi()
	return math.Sin(phi)
}

func (p4 *IPtCotThPhiM) Set(p P4) {
	p4[0] = p.IPt()
	p4[1] = p.CotTh()
	p4[2] = p.Phi()
	p4[3] = p.M()
}
