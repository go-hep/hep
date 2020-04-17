// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
		t.Fatalf("expected e-charge=-1/3. got=%v", p.Charge)
	}

	p = ParticleByID(1000020040)
	if p == nil {
		t.Fatalf("could not retrieve info about pid=1000020040")
	}
	if p.ID != 1000020040 {
		t.Fatalf("expected pid=1000020040. got=%d", p.ID)
	}
	if p.Name != "Alpha-(He4)" {
		t.Fatalf("expected name=Alpha-(He4). got=%q", p.Name)
	}
	if p.Mass != 3.72742 {
		t.Fatalf("expected mass=3.72742. got=%v", p.Mass)
	}
	if math.Abs(2.-p.Charge) > epsilon {
		t.Fatalf("expected e-charge=2. got=%v", p.Charge)
	}

}

func TestLocation(t *testing.T) {
	if Nj != 1 {
		t.Fatalf("expected Nj==1. got=%d", Nj)
	}

	if Nq3 != 2 {
		t.Fatalf("expected Nq3==2. got=%d", Nq3)
	}

	if Nq2 != 3 {
		t.Fatalf("expected Nq2==3. got=%d", Nq2)
	}

	if Nq1 != 4 {
		t.Fatalf("expected Nq1==4. got=%d", Nq1)
	}

	if Nl != 5 {
		t.Fatalf("expected Nl==5. got=%d", Nl)
	}

	if Nr != 6 {
		t.Fatalf("expected Nr==6. got=%d", Nr)
	}

	if N != 7 {
		t.Fatalf("expected N==7. got=%d", N)
	}

	if N8 != 8 {
		t.Fatalf("expected N8==8. got=%d", N8)
	}

	if N9 != 9 {
		t.Fatalf("expected N9==9. got=%d", N9)
	}

	if N10 != 10 {
		t.Fatalf("expected N10==10. got=%d", N10)
	}

}
