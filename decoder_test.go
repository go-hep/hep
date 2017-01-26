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

	if got, want := dec.Len(), int64(32); got != want {
		t.Fatalf("got len=%v. want=%v", got, want)
	}
	start := dec.Pos()
	if start != 0 {
		t.Fatalf("got start=%v. want=%v", start, 0)
	}

	var x int16
	dec.readBin(&x)
	if dec.err != nil {
		t.Fatalf("error reading int16: %v", dec.err)
	}

	pos := dec.Pos()
	if pos != 2 {
		t.Fatalf("got pos=%v. want=%v", pos, 16)
	}
}
