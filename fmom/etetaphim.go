// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"
)

type EtEtaPhiM [4]float64

func NewEtEtaPhiM(et, eta, phi, m float64) EtEtaPhiM {
	return EtEtaPhiM([4]float64{et, eta, phi, m})
}

func (p4 *EtEtaPhiM) Clone() P4 {
	pp := *p4
	return &pp
}

func (p4 *EtEtaPhiM) Et() float64 {
	return p4[0]
}

func (p4 *EtEtaPhiM) Eta() float64 {
	return p4[1]
}

func (p4 *EtEtaPhiM) Phi() float64 {
	return p4[2]
}

func (p4 *EtEtaPhiM) M() float64 {
	return p4[3]
}

func (p4 *EtEtaPhiM) M2() float64 {
	m := p4.M()
	return m * m
}

func (p4 *EtEtaPhiM) P() float64 {
	m := p4.M()
	e := p4.E()
	if m == 0 {
		return e
	}
	sign := 1.0
	if e < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(e*e-m*m)
}
func (p4 *EtEtaPhiM) P2() float64 {
	m := p4.M()
	e := p4.E()
	return e*e - m*m
}

func (p4 *EtEtaPhiM) CosPhi() float64 {
	phi := p4.Phi()
	return math.Cos(phi)
}

func (p4 *EtEtaPhiM) SinPhi() float64 {
	phi := p4.Phi()
	return math.Sin(phi)
}

func (p4 *EtEtaPhiM) TanTh() float64 {
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
	return 1. / math.Sinh(eta)
}

func (p4 *EtEtaPhiM) CotTh() float64 {
	eta := p4.Eta()
	return math.Sinh(eta)
}

func (p4 *EtEtaPhiM) CosTh() float64 {
	eta := p4.Eta()
	return math.Tanh(eta)
}

func (p4 *EtEtaPhiM) SinTh() float64 {
	eta := p4.Eta()
	abseta := math.Abs(eta)
	if abseta > 710 {
		abseta = 710
	}
	return 1 / math.Cosh(abseta)
}

func (p4 *EtEtaPhiM) Pt() float64 {
	p := p4.P()
	sinth := p4.SinTh()
	return p * sinth
}

func (p4 *EtEtaPhiM) E() float64 {
	et := p4.Et()
	sinth := p4.SinTh()
	return et / sinth
}

func (p4 *EtEtaPhiM) IPt() float64 {
	pt := p4.Pt()
	return 1 / pt
}

func (p4 *EtEtaPhiM) Rapidity() float64 {
	e := p4.E()
	pz := p4.Pz()
	return 0.5 * math.Log((e+pz)/(e-pz))
}

func (p4 *EtEtaPhiM) Px() float64 {
	pt := p4.Pt()
	cosphi := p4.CosPhi()
	return pt * cosphi
}

func (p4 *EtEtaPhiM) Py() float64 {
	pt := p4.Pt()
	sinphi := p4.SinPhi()
	return pt * sinphi
}

func (p4 *EtEtaPhiM) Pz() float64 {
	p := p4.P()
	costh := p4.CosTh()
	return p * costh
}

func (p4 *EtEtaPhiM) Set(p P4) {
	p4[0] = p.Et()
	p4[1] = p.Eta()
	p4[2] = p.Phi()
	p4[3] = p.M()
}
