// Copyright ©2025 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux && ci

package groot_test

import (
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
)

func TestHasCxxROOT(t *testing.T) {
	// test we do have a ROOT/C++ installation on CI.
	if !rtests.HasROOT {
		t.Fatalf("ROOT/C++ must be installed and available on CI")
	}
}
