// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestPxPyPzE(t *testing.T) {
	{
		var p4 PxPyPzE
		if got, want := p4.Px(), 0.0; got != want {
			t.Fatalf("p4.Px=%v, want=%v", got, want)
		}
		if got, want := p4.Py(), 0.0; got != want {
			t.Fatalf("p4.Py=%v, want=%v", got, want)
		}
		if got, want := p4.Pz(), 0.0; got != want {
			t.Fatalf("p4.Pz=%v, want=%v", got, want)
		}
		if got, want := p4.E(), 0.0; got != want {
			t.Fatalf("p4.E=%v, want=%v", got, want)
		}
		if got, want := p4.String(), "fmom.P4{Px:0, Py:0, Pz:0, E:0}"; got != want {
			t.Fatalf("p4=%v, want=%v", got, want)
		}
	}

	{
		p4 := NewPxPyPzE(10, 11, 12, 20)
		if got, want := p4.Px(), 10.0; got != want {
			t.Fatalf("p4.Px=%v, want=%v", got, want)
		}
		if got, want := p4.Py(), 11.0; got != want {
			t.Fatalf("p4.Py=%v, want=%v", got, want)
		}
		if got, want := p4.Pz(), 12.0; got != want {
			t.Fatalf("p4.Pz=%v, want=%v", got, want)
		}
		if got, want := p4.E(), 20.0; got != want {
			t.Fatalf("p4.E=%v, want=%v", got, want)
		}
		if got, want := p4.X(), 10.0; got != want {
			t.Fatalf("p4.X=%v, want=%v", got, want)
		}
		if got, want := p4.Y(), 11.0; got != want {
			t.Fatalf("p4.Y=%v, want=%v", got, want)
		}
		if got, want := p4.Z(), 12.0; got != want {
			t.Fatalf("p4.Z=%v, want=%v", got, want)
		}
		if got, want := p4.T(), 20.0; got != want {
			t.Fatalf("p4.T=%v, want=%v", got, want)
		}
		if got, want := p4.String(), "fmom.P4{Px:10, Py:11, Pz:12, E:20}"; got != want {
			t.Fatalf("p4=%v, want=%v", got, want)
		}

		p1 := NewPxPyPzE(10, 11, 12, 20)
		if p1 != p4 {
			t.Fatalf("p4=%v, want=%v", p1, p4)
		}

		var p2 PxPyPzE = p1
		if p1 != p2 {
			t.Fatalf("p4=%v, want=%v", p1, p2)
		}
	}

	{
		p1 := NewPxPyPzE(10, 11, 12, 20)
		var p2 PxPyPzE
		p2.Set(&p1)
		if p1 != p2 {
			t.Fatalf("p4=%v want=%v", p2, p1)
		}
	}

	p := NewPxPyPzE(10, 11, 12, 20)

	// values obtained with ROOT-6.14.00
	for _, tc := range []struct {
		name string
		got  float64
		want float64
		ulp  uint
	}{
		{
			name: "phi",
			got:  p.Phi(),
			want: 8.329812666744317e-01,
			ulp:  1,
		},
		{
			name: "sin-phi",
			got:  p.SinPhi(),
			want: 7.399400733959437e-01,
			ulp:  1,
		},
		{
			name: "cos-phi",
			got:  p.CosPhi(),
			want: 6.726727939963124e-01,
			ulp:  1,
		},
		{
			name: "eta",
			got:  p.Eta(),
			want: 7.382863647914931e-01,
			ulp:  1,
		},
		{
			name: "tan-th",
			got:  p.TanTh(),
			want: 1.238839062276542e+00,
			ulp:  1,
		},
		{
			name: "cos-th",
			got:  p.CosTh(),
			want: 6.281087071082564e-01,
			ulp:  1,
		},
		{
			name: "rapidity",
			got:  p.Rapidity(),
			want: 6.931471805599453e-01,
			ulp:  1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !scalar.EqualWithinULP(tc.got, tc.want, tc.ulp) {
				t.Fatalf("error (ulp=%d)\ngot = %v\nwant= %v", tc.ulp, tc.got, tc.want)
			}
		})
	}

	const epsilon = 1e-12
	t.Run("set-PtEtaPhiM", func(t *testing.T) {
		p1 := NewPxPyPzE(10, 20, 30, 40)
		p1.SetPtEtaPhiM(100, 1.5, 1/3.*math.Pi, 10)
		want := NewPxPyPzE(
			49.99999999999999,
			86.60254037844388,
			212.9279455094818,
			235.45341360636257,
		)
		if got := p1; !p4equal(&got, &want, epsilon) {
			t.Fatalf("invalid p4:\ngot= %v\nwant=%v", got, want)
		}
	})

	t.Run("set-PtEtaPhiE", func(t *testing.T) {
		p1 := NewPxPyPzE(10, 20, 30, 40)
		p1.SetPtEtaPhiE(100, 1.5, 1/3.*math.Pi, 10)
		want := NewPxPyPzE(
			49.99999999999999,
			86.60254037844388,
			212.9279455094818,
			10,
		)
		if got := p1; !p4equal(&got, &want, epsilon) {
			t.Fatalf("invalid p4:\ngot= %v\nwant=%v", got, want)
		}
	})
}
