// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom_test

import (
	"math"
	"testing"

	"go-hep.org/x/hep/fmom"
)

func TestEtEtaPhiM(t *testing.T) {
	{
		var p4 fmom.EtEtaPhiM
		if got, want := p4.Et(), 0.0; got != want {
			t.Fatalf("p4.Et=%v, want=%v", got, want)
		}
		if got, want := p4.Eta(), 0.0; got != want {
			t.Fatalf("p4.Eta=%v, want=%v", got, want)
		}
		if got, want := p4.Phi(), 0.0; got != want {
			t.Fatalf("p4.Phi=%v, want=%v", got, want)
		}
		if got, want := p4.M(), 0.0; got != want {
			t.Fatalf("p4.M=%v, want=%v", got, want)
		}
	}

	{
		p4 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		if got, want := p4.Et(), 10.0; got != want {
			t.Fatalf("p4.Et=%v, want=%v", got, want)
		}
		if got, want := p4.Eta(), 11.0; got != want {
			t.Fatalf("p4.Eta=%v, want=%v", got, want)
		}
		if got, want := p4.Phi(), 12.0; got != want {
			t.Fatalf("p4.Phi=%v, want=%v", got, want)
		}
		if got, want := p4.M(), 20.0; got != want {
			t.Fatalf("p4.M=%v, want=%v", got, want)
		}

		p1 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		if p1 != p4 {
			t.Fatalf("p4=%v, want=%v", p1, p4)
		}

		var p2 fmom.EtEtaPhiM
		p2 = p1
		if p1 != p2 {
			t.Fatalf("p4=%v, want=%v", p1, p2)
		}
	}

	{
		p1 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		var p2 fmom.EtEtaPhiM
		p2.Set(&p1)
		if p1 != p2 {
			t.Fatalf("p4=%v want=%v", p2, p1)
		}
	}

	p := fmom.NewPxPyPzE(10, 11, 12, 20)

	for i, v := range []float64{
		math.Abs(math.Atan2(p.SinPhi(), p.CosPhi()) - p.Phi()),
		math.Abs(p.SinPhi()*p.SinPhi() + p.CosPhi()*p.CosPhi() - 1.0),
		math.Abs(-math.Log(math.Tan(math.Atan(p.TanTh())*0.5)) - p.Eta()),
	} {
		if v > epsilon {
			t.Fatalf("test [%d]: value out of tolerance", i)
		}

	}
}
