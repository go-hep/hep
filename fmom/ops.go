// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// Equal returns true if p1==p2
func Equal(p1, p2 P4) bool {
	return p4equal(p1, p2, 1e-14)
}

func p4equal(p1, p2 P4, epsilon float64) bool {
	if cmpeq(p1.E(), p2.E(), epsilon) &&
		cmpeq(p1.Px(), p2.Px(), epsilon) &&
		cmpeq(p1.Py(), p2.Py(), epsilon) &&
		cmpeq(p1.Pz(), p2.Pz(), epsilon) {
		return true
	}
	return false
}

func cmpeq(x, y, epsilon float64) bool {
	if x == y {
		return true
	}

	return math.Abs(x-y) < epsilon
}

// Add returns the sum p1+p2.
func Add(p1, p2 P4) P4 {
	// FIXME(sbinet):
	// dispatch most efficient/less-lossy addition
	// based on type(dst) (and, optionally, type(src))
	var sum P4
	switch p1 := p1.(type) {

	case *PxPyPzE:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		sum = &p

	case *EEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp EEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *EtEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp EtEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *PtEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp PtEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *IPtCotThPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp IPtCotThPhiM
		pp.Set(&p)
		sum = &pp

	default:
		panic(fmt.Errorf("fmom: invalid P4 concrete value: %#v", p1))
	}
	return sum
}

// IAdd adds src into dst, and returns dst
func IAdd(dst, src P4) P4 {
	// FIXME(sbinet):
	// dispatch most efficient/less-lossy addition
	// based on type(dst) (and, optionally, type(src))
	var sum P4
	var p4 *PxPyPzE = nil
	switch p1 := dst.(type) {

	case *PxPyPzE:
		p4 = p1
		sum = dst

	case *EEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *EtEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *PtEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *IPtCotThPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	default:
		panic(fmt.Errorf("fmom: invalid P4 concrete value: %#v", dst))
	}
	p4.P4.X += src.Px()
	p4.P4.Y += src.Py()
	p4.P4.Z += src.Pz()
	p4.P4.T += src.E()
	sum.Set(p4)
	return sum
}

// Scale returns a*p
func Scale(a float64, p P4) P4 {
	// FIXME(sbinet):
	// dispatch most efficient/less-lossy operation
	// based on type(dst) (and, optionally, type(src))
	var out P4
	switch p := p.(type) {

	case *PxPyPzE:
		dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
		out = &dst

	case *EEtaPhiM:
		dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
		var pp EEtaPhiM
		pp.Set(&dst)
		out = &pp

	case *EtEtaPhiM:
		dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
		var pp EtEtaPhiM
		pp.Set(&dst)
		out = &pp

	case *PtEtaPhiM:
		dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
		var pp PtEtaPhiM
		pp.Set(&dst)
		out = &pp

	case *IPtCotThPhiM:
		dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
		var pp IPtCotThPhiM
		pp.Set(&dst)
		out = &pp

	default:
		panic(fmt.Errorf("fmom: invalid P4 concrete value: %#v", p))
	}

	return out
}

// InvMass computes the invariant mass of two incoming 4-vectors p1 and p2.
func InvMass(p1, p2 P4) float64 {
	p := Add(p1, p2)
	return p.M()
}

// BoostOf returns the 3d boost vector of the provided four-vector p.
// It panics if p has zero energy and a non-zero |p|^2.
// It panics if p isn't a timelike four-vector.
func BoostOf(p P4) Vec3 {
	e := p.E()
	if e == 0 {
		if p.P2() == 0 {
			return Vec3{}
		}
		panic("fmom: zero-energy four-vector")
	}
	if p.M2() <= 0 {
		panic("fmom: non-timelike four-vector")
	}

	inv := 1 / e
	return Vec3{inv * p.Px(), inv * p.Py(), inv * p.Pz()}
}

// Boost returns a copy of the provided four-vector
// boosted by the provided three-vector.
func Boost(p4 P4, boost Vec3) P4 {
	o := p4.Clone()
	if boost == (Vec3{}) {
		return o
	}

	var (
		px = p4.Px()
		py = p4.Py()
		pz = p4.Pz()
		ee = p4.E()

		b  = mat.NewVecDense(3, boost[:])
		b2 = mat.Dot(b, b)

		p3 = mat.NewVecDense(3, []float64{px, py, pz})
		bp = mat.Dot(b, p3)

		ga = 1 / math.Sqrt(1-b2)
		gg = (ga - 1) / b2

		alpha = gg*bp + ga*ee
	)

	pp := NewPxPyPzE(
		px+alpha*boost.X(),
		py+alpha*boost.Y(),
		pz+alpha*boost.Z(),
		ga*(ee+bp),
	)

	o.Set(&pp)

	return o
}
