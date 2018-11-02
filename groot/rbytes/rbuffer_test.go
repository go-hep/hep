// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import (
	"testing"
)

func TestRBuffer(t *testing.T) {
	data := make([]byte, 32)
	r := NewRBuffer(data, nil, 0, nil)

	if got, want := r.Len(), int64(32); got != want {
		t.Fatalf("got len=%v. want=%v", got, want)
	}
	start := r.Pos()
	if start != 0 {
		t.Fatalf("got start=%v. want=%v", start, 0)
	}

	_ = r.ReadI16()
	if r.Err() != nil {
		t.Fatalf("error reading int16: %v", r.Err())
	}

	pos := r.Pos()
	if pos != 2 {
		t.Fatalf("got pos=%v. want=%v", pos, 16)
	}

	pos = 0
	data = make([]byte, 2*(2+4+8))
	r = NewRBuffer(data, nil, 0, nil)
	for _, n := range []int{2, 4, 8} {
		beg := r.Pos()
		if beg != pos {
			t.Errorf("pos[%d] error: got=%d, want=%d\n", n, beg, pos)
		}
		switch n {
		case 2:
			_ = r.ReadI16()
			_ = r.ReadU16()
		case 4:
			_ = r.ReadI32()
			_ = r.ReadU32()
		case 8:
			_ = r.ReadI64()
			_ = r.ReadU64()
		}
		end := r.Pos()
		pos += int64(2 * n)

		if got, want := end-beg, int64(2*n); got != want {
			t.Errorf("%d-bytes: got=%d. want=%d\n", n, got, want)
		}
	}
}
