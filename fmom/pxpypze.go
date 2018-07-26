// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"
)

type PxPyPzE [4]float64

func NewPxPyPzE(px, py, pz, e float64) PxPyPzE {
	return PxPyPzE([4]float64{px, py, pz, e})
}

func (p4 *PxPyPzE) Clone() P4 {
	pp := *p4
	return &pp
}

func (p4 *PxPyPzE) Px() float64 {
	return p4[0]
}

func (p4 *PxPyPzE) Py() float64 {
	return p4[1]
}

func (p4 *PxPyPzE) Pz() float64 {
	return p4[2]
}

func (p4 *PxPyPzE) E() float64 {
	return p4[3]
}

func (p4 *PxPyPzE) X() float64 {
	return p4[0]
}

func (p4 *PxPyPzE) Y() float64 {
	return p4[1]
}

func (p4 *PxPyPzE) Z() float64 {
	return p4[2]
}

func (p4 *PxPyPzE) T() float64 {
	return p4[3]
}

func (p4 *PxPyPzE) M2() float64 {
	px := p4.Px()
	py := p4.Py()
	pz := p4.Pz()
	e := p4.E()

	m2 := e*e - (px*px + py*py + pz*pz)
	return m2
}

func (p4 *PxPyPzE) M() float64 {
	m2 := p4.M2()
	if m2 < 0.0 {
		return -math.Sqrt(-m2)
	}
	return +math.Sqrt(+m2)
}

func (p4 *PxPyPzE) Eta() float64 {
	px := p4.Px()
	py := p4.Py()
	pz := p4.Pz()
	e := p4.E()

	// FIXME: should we use a more underflow-friendly formula:
	//  sqrt(a**2 + b**2)
	//   => y.sqrt(1+(x/y)**2) where y=max(|a|,|b|) and x=min(|a|,|b|)
	//
	p := math.Sqrt(px*px + py*py + pz*pz)
	switch p {
	case 0.0:
		return 0
	case +pz:
		return math.Inf(+1)
	case -pz:
		return math.Inf(-1)
	}
	// flip if negative e
	sign := 1.0
	if e < 0 {
		sign = -1.0
	}
	return sign * 0.5 * math.Log((p+pz)/(p-pz))
}

func (p4 *PxPyPzE) Phi() float64 {
	e := p4.E()
	// flip if negative e
	sign := 1.0
	if e < 0 {
		sign = -1.0
	}
	px := sign * p4.Px()
	py := sign * p4.Py()
	if px == 0.0 && py == 0.0 {
		return 0
	}
	return math.Atan2(py, px)
}

func (p4 *PxPyPzE) P2() float64 {
	px := p4.Px()
	py := p4.Py()
	pz := p4.Pz()

	return px*px + py*py + pz*pz
}

func (p4 *PxPyPzE) P() float64 {
	e := p4.E()
	// flip if negative e
	sign := 1.0
	if e < 0 {
		sign = -1.0
	}
	p2 := p4.P2()
	return sign * math.Sqrt(p2)
}

func (p4 *PxPyPzE) CosPhi() float64 {
	px := p4.Px()
	ipt := p4.IPt()
	return px * ipt
}

func (p4 *PxPyPzE) SinPhi() float64 {
	py := p4.Py()
	ipt := p4.IPt()
	return py * ipt
}

func (p4 *PxPyPzE) TanTh() float64 {
	pt := p4.Pt()
	pz := p4.Pz()
	return pt / pz
}

func (p4 *PxPyPzE) CotTh() float64 {
	pt := p4.Pt()
	pz := p4.Pz()
	return pz / pt
}

func (p4 *PxPyPzE) CosTh() float64 {
	pz := p4.Pz()
	p := p4.P()
	return pz / p
}

func (p4 *PxPyPzE) SinTh() float64 {
	pt := p4.Pt()
	p := p4.P()
	return pt / p
}

func (p4 *PxPyPzE) Pt() float64 {
	e := p4.E()
	px := p4.Px()
	py := p4.Py()
	// flip if negative e
	sign := 1.0
	if e < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(px*px+py*py)
}

func (p4 *PxPyPzE) Et() float64 {
	// to be improved
	e := p4.E()
	sinth := p4.SinTh()
	return e * sinth
}

func (p4 *PxPyPzE) IPt() float64 {
	pt := p4.Pt()
	return 1.0 / pt
}

func (p4 *PxPyPzE) Rapidity() float64 {
	e := p4.E()
	pz := p4.Pz()
	switch e {
	case 0.0:
		return 0.0
	case +pz:
		return math.Inf(+1)
	case -pz:
		return math.Inf(-1)
	}
	// invariant under flipping of 4-mom with negative energy
	return 0.5 * math.Log((e+pz)/(e-pz))
}

func (p4 *PxPyPzE) Set(p P4) {
	p4[0] = p.Px()
	p4[1] = p.Py()
	p4[2] = p.Pz()
	p4[3] = p.E()
}
