package fmom_test

import (
	"math"
	"testing"

	"github.com/go-hep/fmom"
)

func TestEtEtaPhiM(t *testing.T) {
	{
		var p4 fmom.EtEtaPhiM
		if p4.Et() != 0 {
			t.Fatalf("expected p4.Et=%v (got=%v)", 0, p4.Et())
		}
		if p4.Eta() != 0 {
			t.Fatalf("expected p4.Eta=%v (got=%v)", 0, p4.Eta())
		}
		if p4.Phi() != 0 {
			t.Fatalf("expected p4.Phi=%v (got=%v)", 0, p4.Phi())
		}
		if p4.M() != 0 {
			t.Fatalf("expected p4.M=%v (got=%v)", 0, p4.M())
		}
	}

	{
		p4 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		if p4.Et() != 10 {
			t.Fatalf("expected p4.Et=%v (got=%v)", 10, p4.Et())
		}
		if p4.Eta() != 11 {
			t.Fatalf("expected p4.Eta=%v (got=%v)", 11, p4.Eta())
		}
		if p4.Phi() != 12 {
			t.Fatalf("expected p4.Phi=%v (got=%v)", 12, p4.Phi())
		}
		if p4.M() != 20 {
			t.Fatalf("expected p4.M=%v (got=%v)", 20, p4.M())
		}

		p1 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		if p1 != p4 {
			t.Fatalf("expected p4=%v (got=%v)", p4, p1)
		}

		var p2 fmom.EtEtaPhiM
		p2 = p1
		if p1 != p2 {
			t.Fatalf("expected p4=%v (got=%v)", p2, p1)
		}
	}

	{
		p1 := fmom.NewEtEtaPhiM(10, 11, 12, 20)
		var p2 fmom.EtEtaPhiM
		p2.Set(&p1)
		if p1 != p2 {
			t.Fatalf("expected p4=%v (got=%v)", p1, p2)
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

// EOF
