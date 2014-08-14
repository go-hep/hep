package fmom

import (
	"reflect"
	"testing"
)

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
	for _, table := range []struct {
		p1  P4
		p2  P4
		exp P4
	}{
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p1 := table.p1.Clone()
		p2 := table.p2.Clone()

		sum := Add(p1, p2)

		if !deepEqual(sum, table.exp) {
			t.Fatalf("exp: %#v\ngot: %#v", table.exp, sum)
		}
		if !reflect.DeepEqual(p1, table.p1) {
			t.Fatalf("add modified p1:\np1=%#v (ref)\np1=%#v (new)", table.p1, p1)
		}
		if !reflect.DeepEqual(p2, table.p2) {
			t.Fatalf("add modified p2:\np1=%#v (ref)\np2=%#v (new)", table.p2, p2)
		}

	}
}

func TestIAdd(t *testing.T) {
	for _, table := range []struct {
		p1  P4
		p2  P4
		exp P4
	}{
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p1:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p1 := table.p1.Clone()
		p2 := table.p2.Clone()

		sum := IAdd(p1, p2)

		if !deepEqual(sum, table.exp) {
			t.Fatalf("exp: %#v\ngot: %#v", table.exp, sum)
		}

		if !reflect.DeepEqual(sum, p1) {
			t.Fatalf("fmom.IAdd did not modify p1 in-place:\nexp: %#v\ngot: %#v", sum, p1)
		}
		if !reflect.DeepEqual(p2, table.p2) {
			t.Fatalf("fmom.IAdd modified p2:\np1=%#v (ref)\np2=%#v (new)", table.p2, p2)
		}
	}
}

func TestEqual(t *testing.T) {
	for _, table := range []struct {
		p1  P4
		p2  P4
		exp bool
	}{
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: true,
		},
		{
			p1:  newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: true,
		},
		{
			p1:  newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: true,
		},
		{
			p1:  newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: true,
		},
		{
			p1:  newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: true,
		},

		{
			p1:  newPxPyPzE(NewPxPyPzE(10+1e-14, 10, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: false,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10+1e-14, 10, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: false,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10+1e-14, 20)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: false,
		},
		{
			p1:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20+1e-14)),
			p2:  newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			exp: false,
		},
	} {
		{
			eq := deepEqual(table.p1, table.p2)
			if eq != table.exp {
				t.Fatalf("exp: %#v\ngot: %#v\np1=%#v\np2=%#v\n", table.exp, eq, table.p1, table.p2)
			}
		}
		eq := Equal(table.p1, table.p2)
		if eq != table.exp {
			t.Fatalf("exp: %#v\ngot: %#v\np1=%#v\np2=%#v\n", table.exp, eq, table.p1, table.p2)
		}
	}
}

func TestScale(t *testing.T) {
	for _, table := range []struct {
		p   P4
		a   float64
		exp P4
	}{
		{
			p:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:   1,
			exp: newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
		},

		{
			p:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:   0,
			exp: newPxPyPzE(NewPxPyPzE(0, 0, 0, 0)),
		},

		{
			p:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:   -1,
			exp: newPxPyPzE(NewPxPyPzE(-10, -10, -10, -20)),
		},

		{
			p:   newPxPyPzE(NewPxPyPzE(10, 10, 10, 20)),
			a:   2,
			exp: newPxPyPzE(NewPxPyPzE(20, 20, 20, 40)),
		},

		{
			p:   newEEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:   2,
			exp: newEEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:   newEtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:   2,
			exp: newEtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:   newPtEtaPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:   2,
			exp: newPtEtaPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
		{
			p:   newIPtCotThPhiM(NewPxPyPzE(10, 10, 10, 20)),
			a:   2,
			exp: newIPtCotThPhiM(NewPxPyPzE(20, 20, 20, 40)),
		},
	} {
		p := table.p.Clone()

		o := Scale(table.a, p)

		if !deepEqual(o, table.exp) {
			t.Fatalf("exp: %#v\ngot: %#v", table.exp, o)
		}
		if !reflect.DeepEqual(p, table.p) {
			t.Fatalf("add modified p:\np=%#v (ref)\np=%#v (new)", table.p, p)
		}
	}
}
