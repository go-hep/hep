// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom_test

import (
	"testing"

	"go-hep.org/x/hep/fmom"
)

func TestVec3(t *testing.T) {
	px := 20.0
	py := 21.0
	pz := 22.0
	p := fmom.Vec3{px, py, pz}

	if got, want := p.X(), px; got != want {
		t.Fatalf("px differ. got=%v, want=%v", got, want)
	}

	if got, want := p.Y(), py; got != want {
		t.Fatalf("py differ. got=%v, want=%v", got, want)
	}

	if got, want := p.Z(), pz; got != want {
		t.Fatalf("pz differ. got=%v, want=%v", got, want)
	}
}
