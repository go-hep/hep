// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt_test

import (
	"testing"

	"go-hep.org/x/hep/heppdt"
)

func TestDefaultTable(t *testing.T) {
	if got, want := heppdt.Name(), "particle.tbl"; got != want {
		t.Fatalf("invalid table name. got=%q, want=%q", got, want)
	}

	if got, want := heppdt.Len(), 534; got != want {
		t.Fatalf("invalid table length. got=%d, want=%d", got, want)
	}

	particles := heppdt.PDT()
	if got, want := particles[1].Name, "d"; got != want {
		t.Fatalf("invalid particle for pid=1. got=%q, want=%q", got, want)
	}
}
