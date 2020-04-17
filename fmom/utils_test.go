// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

import (
	"math"
	"testing"
)

const (
	epsilon_test = 1e-6
)

func TestDeltaR(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want float64
	}{
		// pxpypze
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			want: 3.4064618746379645,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.5707963267948966,
		},

		// eetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 3.4064618746379645,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.5707963267948966,
		},

		// etetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 3.4064618746379645,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.5707963267948966,
		},

		// ptetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 3.4064618746379645,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.5707963267948966,
		},

		// iptcotthphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 3.4064618746379645,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.5707963267948966,
		},
	} {
		got := DeltaR(tc.p1, tc.p2)
		if got-tc.want > epsilon_test {
			t.Fatalf("DeltaR error\np1=%#v\np2=%#v\ngot= %+e\nwantt=%+e\n",
				tc.p1,
				tc.p2,
				got,
				tc.want,
			)
		}
	}
}

func TestDeltaPhi(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want float64
	}{
		// pxpypze
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			want: -math.Pi,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			want: math.Pi / 2.0,
		},

		// eetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -math.Pi,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: math.Pi / 2.0,
		},

		// etetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -math.Pi,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: math.Pi / 2.0,
		},

		// ptetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -math.Pi,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: math.Pi / 2.0,
		},

		// iptcotthphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 0,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -math.Pi,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: math.Pi / 2.0,
		},
	} {
		got := DeltaPhi(tc.p1, tc.p2)
		if got-tc.want > epsilon_test {
			t.Fatalf("DeltaPhi error\np1=%#v\np2=%#v\ngot= %+e\nwant=%+e\n",
				tc.p1,
				tc.p2,
				got,
				tc.want,
			)
		}
	}
}

func TestDot(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want float64
	}{
		// pxpypze
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: 100,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			want: 700,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			want: 300,
		},

		// eetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 100,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 700,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 300,
		},

		// etetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 100,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 700,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 300,
		},

		// ptetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 100,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 700,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 300,
		},

		// iptcotthphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 100,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: 700,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 300,
		},
	} {
		got := Dot(tc.p1, tc.p2)
		if got-tc.want > epsilon_test {
			t.Fatalf("Dot error\np1=%#v\np2=%#v\ngot= %+e\nwant=%+e\n",
				tc.p1,
				tc.p2,
				got,
				tc.want,
			)
		}
	}
}
func TestCosTheta(t *testing.T) {
	for _, tc := range []struct {
		p1   P4
		p2   P4
		want float64
	}{
		// pxpypze
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			want: 1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			want: -1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.0 / 3,
		},

		// eetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.0 / 3,
		},

		// etetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.0 / 3,
		},

		// ptetaphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.0 / 3,
		},

		// iptcotthphim
		{
			p1:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			want: 1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			want: -1,
		},
		{
			p1:   newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:   newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			want: 1.0 / 3,
		},
	} {
		got := CosTheta(tc.p1, tc.p2)
		if got-tc.want > epsilon_test {
			t.Fatalf("CosTheta error\np1=%#v\np2=%#v\ngot= %+e\nwant=%+e\n",
				tc.p1,
				tc.p2,
				got,
				tc.want,
			)
		}
	}
}
