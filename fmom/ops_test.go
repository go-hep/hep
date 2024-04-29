// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"reflect"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/spatial/r3"
)

func p3equal(p1, p2 r3.Vec, epsilon float64) bool {
	if cmpeq(p1.X, p2.X, epsilon) &&
		cmpeq(p1.Y, p2.Y, epsilon) &&
		cmpeq(p1.Z, p2.Z, epsilon) {
		return true
	}
	return false
}

func newPxPyPzE(p4 PxPyPzE) P4 {
	return &p4
}

func newEEtaPhiM(p4 PxPyPzE) P4 {
	var pp EEtaPhiM
	pp.Set(&p4)
	return &pp
}

func newEtEtaPhiM(p4 PxPyPzE) P4 {
	var pp EtEtaPhiM
	pp.Set(&p4)
	return &pp
}

func newPtEtaPhiM(p4 PxPyPzE) P4 {
	var pp PtEtaPhiM
	pp.Set(&p4)
	return &pp
}

func newIPtCotThPhiM(p4 PxPyPzE) P4 {
	var pp IPtCotThPhiM
	pp.Set(&p4)
	return &pp
}

func deepEqual(p1, p2 P4) bool {
	return Equal(p1, p2)
}

func TestAdd(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want P4
	}{
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p1 := tc.p1.Clone()
		p2 := tc.p2.Clone()

		got := Add(p1, p2)

		if !deepEqual(got, tc.want) {
			t.Fatalf("got= %#v\nwant=%#v", got, tc.want)
		}
		if !reflect.DeepEqual(p1, tc.p1) {
			t.Fatalf("add modified p1:\ngot= %#v\nwant=%#v", p1, tc.p1)
		}
		if !reflect.DeepEqual(p2, tc.p2) {
			t.Fatalf("add modified p2:\ngot= %#v\nwant=%#v", p2, tc.p2)
		}
	}
}

func TestIAdd(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want P4
	}{
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p1 := tc.p1.Clone()
		p2 := tc.p2.Clone()

		got := IAdd(p1, p2)

		if !deepEqual(got, tc.want) {
			t.Fatalf("got= %#v\nwant=%#v", got, tc.want)
		}

		if !reflect.DeepEqual(got, p1) {
			t.Fatalf("fmom.IAdd did not modify p1 in-place:\ngot= %#v\nwant=%#v", got, p1)
		}
		if !reflect.DeepEqual(p2, tc.p2) {
			t.Fatalf("fmom.IAdd modified p2:\ngot= %#v\nwant=%#v", p2, tc.p2)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want bool
	}{
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: true,
		},
		{
			p1:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: true,
		},
		{
			p1:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: true,
		},
		{
			p1:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: true,
		},
		{
			p1:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: true,
		},

		{
			p1:   newPxPyPzE(NewPxPyPzE(10+1e-14, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: false,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10+1e-14, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: false,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10+1e-14, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: false,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20+1e-14)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: false,
		},
	} {
		{
			got := deepEqual(tc.p1, tc.p2)
			if got != tc.want {
				t.Fatalf("got= %#v\nwant=%#v\np1=%#v\np2=%#v\n", got, tc.want, tc.p1, tc.p2)
			}
		}
		got := Equal(tc.p1, tc.p2)
		if got != tc.want {
			t.Fatalf("got= %#v\nwant=%#v\np1=%#v\np2=%#v\n", got, tc.want, tc.p1, tc.p2)
		}
	}
}

func TestScale(t *testing.T) {
	for _, tc := range []struct {
		p    P4
		a    float64
		want P4
	}{
		{
			p:    newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:    1,
			want: newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
		},

		{
			p:    newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:    0,
			want: newPxPyPzE(NewPxPyPzE(0, 0, 0, 0)),
		},

		{
			p:    newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:    -1,
			want: newPxPyPzE(NewPxPyPzE(-10, -10, -10, -20)),
		},

		{
			p:    newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:    2,
			want: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},

		{
			p:    newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:    2,
			want: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:    newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:    2,
			want: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:    newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:    2,
			want: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:    newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:    2,
			want: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p := tc.p.Clone()

		got := Scale(tc.a, p)

		if !deepEqual(got, tc.want) {
			t.Fatalf("got= %#v\nwant=%#v", got, tc.want)
		}
		if !reflect.DeepEqual(p, tc.p) {
			t.Fatalf("add modified p:\np=%#v (ref)\np=%#v (new)", tc.p, p)
		}
	}
}

func TestInvMass(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want float64
	}{
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)).M(),
		},
		{
			p1:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)).M(),
		},
		{
			p1:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)).M(),
		},
		{
			p1:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)).M(),
		},
		{
			p1:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)).M(),
		},
	} {
		p1 := tc.p1.Clone()
		p2 := tc.p2.Clone()

		got := InvMass(p1, p2)

		if !scalar.EqualWithinULP(got, tc.want, 2) {
			t.Fatalf("got= %#v\nwant=%#v", got, tc.want)
		}

		if !reflect.DeepEqual(tc.p1, p1) {
			t.Fatalf("fmom.InvMass modified p1 in-place:\ngot: %#v\nwant:%#v", p1, tc.p1)
		}
		if !reflect.DeepEqual(tc.p2, p2) {
			t.Fatalf("fmom.InvMass modified p2 in-place:\ngot: %#v\nwant:%#v", p2, tc.p2)
		}
	}
}

func TestBoost(t *testing.T) {
	var (
		p1      = NewPxPyPzE(1, 2, 3, 4)
		boost   = BoostOf(&p1)
		p1RF    = Boost(&p1, r3.Vec{X: -boost.X, Y: -boost.Y, Z: -boost.Z})
		boostRF = BoostOf(p1RF)
		zero    r3.Vec
	)

	if !p3equal(boostRF, zero, 1e-14) {
		t.Fatalf("invalid boost: got=%v, want=%v", boostRF, zero)
	}

	if got, want := Boost(&p1, r3.Vec{}), &p1; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid zero-boost: got=%v, want=%v", got, want)
	}
}

func TestBoostOf(t *testing.T) {
	for _, tc := range []struct {
		p      P4
		v      r3.Vec
		panics string
	}{
		{
			p: newPxPyPzE(NewPxPyPzE(1, 2, 4, 10)),
			v: r3.Vec{X: 0.1, Y: 0.2, Z: 0.4},
		},
		{
			p: newPxPyPzE(NewPxPyPzE(1, 2, 4, -10)),
			v: r3.Vec{X: -0.1, Y: -0.2, Z: -0.4},
		},
		{
			p:      newPxPyPzE(NewPxPyPzE(1, 2, 3, 1)),
			v:      r3.Vec{},
			panics: "fmom: non-timelike four-vector",
		},
		{
			p:      newPxPyPzE(NewPxPyPzE(1, 2, 3, 0)),
			v:      r3.Vec{},
			panics: "fmom: zero-energy four-vector",
		},
		{
			p: newPxPyPzE(NewPxPyPzE(0, 0, 0, 0)),
			v: r3.Vec{},
		},
	} {
		t.Run("", func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					e := recover()
					if e == nil {
						t.Fatalf("expected a panic: got=%v, want=%v", e, tc.panics)
					}
					if got, want := e.(string), tc.panics; got != want {
						t.Fatalf("invalid panic message:\ngot= %v\nwant=%v", got, want)
					}
				}()
			}

			got := BoostOf(tc.p)
			if got, want := got, tc.v; got != want {
				t.Fatalf("invalid boost vector:\ngot= %v\nwant=%v", got, want)
			}
		})
	}
}

func TestVecOf(t *testing.T) {
	for _, tc := range []struct {
		p    P4
		want r3.Vec
	}{
		{
			p:    newPxPyPzE(NewPxPyPzE(0, 10, 20, 30)),
			want: r3.Vec{X: 0, Y: 10, Z: 20},
		},

		{
			p:    newPxPyPzE(NewPxPyPzE(10, 0, 20, 30)),
			want: r3.Vec{X: 10, Y: 0, Z: 20},
		},

		{
			p:    newPxPyPzE(NewPxPyPzE(10, 20, 0, 30)),
			want: r3.Vec{X: 10, Y: 20, Z: 0},
		},
	} {
		got := VecOf(tc.p)

		if got != tc.want {
			t.Fatalf("invalid spatial components for %#v: got= %#v\nwant=%#v",
				tc.p, got, tc.want)
		}
	}
}
