// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom_test

import (
	"testing"

	"go-hep.org/x/hep/fmom"
	"gonum.org/v1/gonum/floats"
)

func TestPxPyPzE(t *testing.T) {
	{
		var p4 fmom.PxPyPzE
		if p4.Px() != 0 {
			t.Fatalf("expected p4.Px=%v (got=%v)", 0, p4.Px())
		}
		if p4.Py() != 0 {
			t.Fatalf("expected p4.Py=%v (got=%v)", 0, p4.Py())
		}
		if p4.Pz() != 0 {
			t.Fatalf("expected p4.Pz=%v (got=%v)", 0, p4.Pz())
		}
		if p4.E() != 0 {
			t.Fatalf("expected p4.E=%v (got=%v)", 0, p4.E())
		}
	}

	{
		p4 := fmom.NewPxPyPzE(10, 11, 12, 20)
		if p4.Px() != 10 {
			t.Fatalf("expected p4.Px=%v (got=%v)", 10, p4.Px())
		}
		if p4.Py() != 11 {
			t.Fatalf("expected p4.Py=%v (got=%v)", 11, p4.Py())
		}
		if p4.Pz() != 12 {
			t.Fatalf("expected p4.Pz=%v (got=%v)", 12, p4.Pz())
		}
		if p4.E() != 20 {
			t.Fatalf("expected p4.E=%v (got=%v)", 20, p4.E())
		}
		if p4.X() != 10 {
			t.Fatalf("expected p4.X=%v (got=%v)", 10, p4.X())
		}
		if p4.Y() != 11 {
			t.Fatalf("expected p4.Y=%v (got=%v)", 11, p4.Y())
		}
		if p4.Z() != 12 {
			t.Fatalf("expected p4.Z=%v (got=%v)", 12, p4.Z())
		}
		if p4.T() != 20 {
			t.Fatalf("expected p4.T=%v (got=%v)", 20, p4.T())
		}

		p1 := fmom.NewPxPyPzE(10, 11, 12, 20)
		if p1 != p4 {
			t.Fatalf("expected p4=%v (got=%v)", p4, p1)
		}

		var p2 fmom.PxPyPzE
		p2 = p1
		if p1 != p2 {
			t.Fatalf("expected p4=%v (got=%v)", p2, p1)
		}
	}

	{
		p1 := fmom.NewPxPyPzE(10, 11, 12, 20)
		var p2 fmom.PxPyPzE
		p2.Set(&p1)
		if p1 != p2 {
			t.Fatalf("expected p4=%v (got=%v)", p1, p2)
		}
	}

	p := fmom.NewPxPyPzE(10, 11, 12, 20)

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
			if !floats.EqualWithinULP(tc.got, tc.want, tc.ulp) {
				t.Fatalf("error (ulp=%d)\ngot = %v\nwant= %v", tc.ulp, tc.got, tc.want)
			}
		})
	}
}
