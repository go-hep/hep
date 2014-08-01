package fmom

import (
	"math"
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
	const epsilon = 1e-6
	if math.Abs(p1.Px()-p2.Px()) > epsilon {
		return false
	}
	if math.Abs(p1.Py()-p2.Py()) > epsilon {
		return false
	}
	if math.Abs(p1.Pz()-p2.Pz()) > epsilon {
		return false
	}
	if math.Abs(p1.E()-p2.E()) > epsilon {
		return false
	}
	return true
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
