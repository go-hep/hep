// Copyright 2017 The go-hep Authors. All rights reserved.
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

	for _, table := range []struct {
		p1  P4
		p2  P4
		exp float64
	}{
		// pxpypze
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			exp: 3.4064618746379645,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.5707963267948966,
		},

		// eetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: 3.4064618746379645,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.5707963267948966,
		},

		// etetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: 3.4064618746379645,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.5707963267948966,
		},

		// ptetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: 3.4064618746379645,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.5707963267948966,
		},

		// iptcotthphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: 3.4064618746379645,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.5707963267948966,
		},
	} {
		dR := DeltaR(table.p1, table.p2)
		if dR-table.exp > epsilon_test {
			t.Fatalf("DeltaR error\np1=%#v\np2=%#v\nexp=%+e\ngot=%+e\n",
				table.p1,
				table.p2,
				table.exp,
				dR,
			)
		}
	}
}

func TestDeltaPhi(t *testing.T) {

	for _, table := range []struct {
		p1  P4
		p2  P4
		exp float64
	}{
		// pxpypze
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -math.Pi,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			exp: math.Pi / 2.0,
		},

		// eetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -math.Pi,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: math.Pi / 2.0,
		},

		// etetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -math.Pi,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: math.Pi / 2.0,
		},

		// ptetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -math.Pi,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: math.Pi / 2.0,
		},

		// iptcotthphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 0,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -math.Pi,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: math.Pi / 2.0,
		},
	} {
		dphi := DeltaPhi(table.p1, table.p2)
		if dphi-table.exp > epsilon_test {
			t.Fatalf("DeltaPhi error\np1=%#v\np2=%#v\nexp=%+e\ngot=%+e\n",
				table.p1,
				table.p2,
				table.exp,
				dphi,
			)
		}
	}
}

func TestCosTheta(t *testing.T) {
	for _, table := range []struct {
		p1  P4
		p2  P4
		exp float64
	}{
		// pxpypze
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: 1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPxPyPzE(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.0 / 3,
		},

		// eetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.0 / 3,
		},

		// etetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newEtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.0 / 3,
		},

		// ptetaphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newPtEtaPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.0 / 3,
		},

		// iptcotthphim
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			exp: 1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(-10, -10, -10, +20)),
			exp: -1,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(+10, +10, +10, +20)),
			p2:  newIPtCotThPhiM(NewPxPyPzE(+10, -10, +10, +20)),
			exp: 1.0 / 3,
		},
	} {
		costh := CosTheta(table.p1, table.p2)
		if costh-table.exp > epsilon_test {
			t.Fatalf("CosTheta error\np1=%#v\np2=%#v\nexp=%+e\ngot=%+e\n",
				table.p1,
				table.p2,
				table.exp,
				costh,
			)
		}
	}
}
