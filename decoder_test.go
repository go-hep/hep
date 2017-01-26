// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"testing"
)

func TestDecoder(t *testing.T) {
	data := make([]byte, 32)
	dec := newDecoder(bytes.NewBuffer(data))

	if dec.Len() != 32 {
		t.Fatalf("expected len=%v. got %v", len(data), dec.Len())
	}
	start := dec.Pos()
	if start != 0 {
		t.Fatalf("expected start=%v. got %v", 0, start)
	}

	var x int16
	dec.readBin(&x)
	if dec.err != nil {
		t.Fatalf("error reading int16: %v", dec.err)
	}

	pos := dec.Pos()
	if pos != 2 {
		t.Fatalf("expected pos=%v. got %v", 16, pos)
	}
}
