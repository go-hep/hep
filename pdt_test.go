package heppdt

import (
	"math"
	"testing"
)

const epsilon = 1e-20

func TestDefaultTable(t *testing.T) {
	if defaultTable.Len() != 534 {
		t.Fatalf("expected default table to hold [534] particles. got=%v", defaultTable.Len())
	}

	p := ParticleByID(1)
	if p == nil {
		t.Fatalf("could not retrieve info about pid=1")
	}
	if p.ID != 1 {
		t.Fatalf("expected pid=1. got=%d", p.ID)
	}
	if p.Name != "d" {
		t.Fatalf("expected name=d. got=%q", p.Name)
	}
	if p.Mass != 0.33 {
		t.Fatalf("expected mass=0.33. got=%v", p.Mass)
	}
	if math.Abs(-1./3.-p.Charge) > epsilon {
		t.Fatalf("expected e-charge=1/3. got=%v", p.Charge)
	}
}

func TestLocation(t *testing.T) {
	if nj != 1 {
		t.Fatalf("expected nj==1. got=%d", nj)
	}

	if nq3 != 2 {
		t.Fatalf("expected nq3==2. got=%d", nq3)
	}

	if nq2 != 3 {
		t.Fatalf("expected nq2==3. got=%d", nq2)
	}

	if nq1 != 4 {
		t.Fatalf("expected nq1==4. got=%d", nq1)
	}

	if nl != 5 {
		t.Fatalf("expected nl==5. got=%d", nl)
	}

	if nr != 6 {
		t.Fatalf("expected nr==6. got=%d", nr)
	}

	if n != 7 {
		t.Fatalf("expected n==7. got=%d", n)
	}

	if n8 != 8 {
		t.Fatalf("expected n8==8. got=%d", n8)
	}

	if n9 != 9 {
		t.Fatalf("expected n9==9. got=%d", n9)
	}

	if n10 != 10 {
		t.Fatalf("expected n10==10. got=%d", n10)
	}

}
