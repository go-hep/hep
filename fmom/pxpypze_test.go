package fmom_test

import (
	"math"
	"testing"

	"github.com/go-hep/fmom"
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
