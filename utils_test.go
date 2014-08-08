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
		if math.Abs(dR-table.exp) > epsilon_test {
			t.Fatalf("DeltaR differ\np1=%#v\np2=%#v\nexp=%+e\ngot=%+e\n",
				table.p1,
				table.p2,
				table.exp,
				dR,
			)
		}
	}

}
