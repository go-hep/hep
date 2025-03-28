// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"fmt"
	"math"
)

type PtEtaPhiM struct {
	P4 Vec4
}

func NewPtEtaPhiM(pt, eta, phi, m float64) PtEtaPhiM {
	return PtEtaPhiM{P4: Vec4{X: pt, Y: eta, Z: phi, T: m}}
}

func (p4 PtEtaPhiM) String() string {
	return fmt.Sprintf(
		"fmom.P4{Pt:%v, Eta:%v, Phi:%v, M:%v}",
		p4.Pt(), p4.Eta(), p4.Phi(), p4.M(),
	)
}

func (p4 *PtEtaPhiM) Clone() P4 {
	pp := *p4
	return &pp
}

func (p4 *PtEtaPhiM) Pt() float64 {
	return p4.P4.X
}

func (p4 *PtEtaPhiM) Eta() float64 {
	return p4.P4.Y
}

func (p4 *PtEtaPhiM) Phi() float64 {
	return p4.P4.Z
}

func (p4 *PtEtaPhiM) M() float64 {
	return p4.P4.T
}

func (p4 *PtEtaPhiM) E() float64 {
	m := p4.M()
	pt := p4.Pt()
	pz := p4.Pz()

	sign := +1.0
	if pt < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(pt*pt+pz*pz+m*m)
}

func (p4 *PtEtaPhiM) P() float64 {
	pt := p4.Pt()
	pz := p4.Pz()

	sign := +1.0
	if pt < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(pt*pt+pz*pz)
}

func (p4 *PtEtaPhiM) P2() float64 {
	pt := p4.Pt()
	pz := p4.Pz()

	return pt*pt + pz*pz
}

func (p4 *PtEtaPhiM) M2() float64 {
	m := p4.M()
	return m * m
}

func (p4 *PtEtaPhiM) CosPhi() float64 {
	phi := p4.Phi()
	return math.Cos(phi)
}

func (p4 *PtEtaPhiM) SinPhi() float64 {
	phi := p4.Phi()
	return math.Sin(phi)
}

func (p4 *PtEtaPhiM) CotTh() float64 {
	eta := p4.Eta()
	return math.Sinh(eta)
}

func (p4 *PtEtaPhiM) CosTh() float64 {
	eta := p4.Eta()
	return math.Tanh(eta)
}

func (p4 *PtEtaPhiM) SinTh() float64 {
	eta := p4.Eta()
	// avoid numeric overflow if very large eta
	abseta := min(math.Abs(eta), 710)
	return 1 / math.Cosh(abseta)
}

func (p4 *PtEtaPhiM) TanTh() float64 {
	eta := p4.Eta()
	abseta := math.Abs(eta)
	// avoid numeric overflow if very large eta
	if abseta > 710 {
		if eta > 0 {
			eta = +710
		} else {
			eta = -710
		}
	}
	return 1 / math.Sinh(eta)
}

func (p4 *PtEtaPhiM) Et() float64 {
	e := p4.E()
	sinth := p4.SinTh()
	return e * sinth
}

func (p4 *PtEtaPhiM) IPt() float64 {
	pt := p4.Pt()
	return 1 / pt
}

func (p4 *PtEtaPhiM) Rapidity() float64 {
	e := p4.E()
	pz := p4.Pz()
	return 0.5 * math.Log((e+pz)/(e-pz))
}

func (p4 *PtEtaPhiM) Px() float64 {
	pt := p4.Pt()
	cosphi := p4.CosPhi()
	return pt * cosphi
}

func (p4 *PtEtaPhiM) Py() float64 {
	pt := p4.Pt()
	sinphi := p4.SinPhi()
	return pt * sinphi
}

func (p4 *PtEtaPhiM) Pz() float64 {
	pt := p4.Pt()
	cotth := p4.CotTh()
	return pt * cotth
}

func (p4 *PtEtaPhiM) Set(p P4) {
	p4.P4.X = p.Pt()
	p4.P4.Y = p.Eta()
	p4.P4.Z = p.Phi()
	p4.P4.T = p.M()
}

var (
	_ P4           = (*PtEtaPhiM)(nil)
	_ fmt.Stringer = (*PtEtaPhiM)(nil)
)
