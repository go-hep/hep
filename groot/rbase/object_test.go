// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase_test

import (
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
)

func TestObject(t *testing.T) {
	obj := rbase.NewObject()
	if got, want := obj.ID, uint32(0x0); got != want {
		t.Fatalf("invalid ID. got=0x%x, want=0x%x", got, want)
	}
	obj.SetID(0x2)
	if got, want := obj.ID, uint32(0x2); got != want {
		t.Fatalf("invalid ID. got=0x%x, want=0x%x", got, want)
	}

	if got, want := obj.Bits, uint32(0x3000000); got != want {
		t.Fatalf("invalid bits. got=0x%x, want=0x%x", got, want)
	}

	if got, want := obj.TestBits(rbytes.BypassStreamer), false; got != want {
		t.Fatalf("invalid BypassStreamer-bit-test. got=%v, want=%v", got, want)
	}
	obj.SetBit(rbytes.BypassStreamer)
	if got, want := obj.TestBits(rbytes.BypassStreamer), true; got != want {
		t.Fatalf("invalid BypassStreamer-bit-test. got=%v, want=%v", got, want)
	}
	obj.ResetBit(rbytes.BypassStreamer)
	if got, want := obj.TestBits(rbytes.BypassStreamer), false; got != want {
		t.Fatalf("invalid BypassStreamer-bit-test. got=%v, want=%v", got, want)
	}

	obj.SetBits(rbytes.BypassStreamer)
	if got, want := obj.Bits, uint32(rbytes.BypassStreamer); got != want {
		t.Fatalf("invalid bits. got=0x%x, want=0x%x", got, want)
	}
	if got, want := obj.TestBits(rbytes.BypassStreamer), true; got != want {
		t.Fatalf("invalid BypassStreamer-bit-test. got=%v, want=%v", got, want)
	}
}
