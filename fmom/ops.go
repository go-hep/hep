// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/spatial/r3"
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
	out := p.Clone()
	dst := NewPxPyPzE(a*p.Px(), a*p.Py(), a*p.Pz(), a*p.E())
	out.Set(&dst)
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
func BoostOf(p P4) r3.Vec {
	e := p.E()
	if e == 0 {
		if p.P2() == 0 {
			return r3.Vec{}
		}
		panic("fmom: zero-energy four-vector")
	}
	if p.M2() <= 0 {
		panic("fmom: non-timelike four-vector")
	}

	inv := 1 / e
	return r3.Vec{X: inv * p.Px(), Y: inv * p.Py(), Z: inv * p.Pz()}
}

// Boost returns a copy of the provided four-vector
// boosted by the provided three-vector.
func Boost(p P4, vec r3.Vec) P4 {
	o := p.Clone()
	if vec == (r3.Vec{}) {
		return o
	}

	var (
		px = p.Px()
		py = p.Py()
		pz = p.Pz()
		ee = p.E()

		p3 = r3.Vec{X: px, Y: py, Z: pz}
		v2 = vecDot(vec, vec)
		bp = vecDot(vec, p3)

		gamma = 1 / math.Sqrt(1-v2)
		beta  = (gamma - 1) / v2

		alpha = beta*bp + gamma*ee
	)

	pp := NewPxPyPzE(
		px+alpha*vec.X,
		py+alpha*vec.Y,
		pz+alpha*vec.Z,
		gamma*(ee+bp),
	)

	o.Set(&pp)

	return o
}

func vecDot(u, v r3.Vec) float64 {
	return u.Dot(v)
}
