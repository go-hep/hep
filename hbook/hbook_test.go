// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "testing"

func TestRangeClone(t *testing.T) {
	r1 := Range{1, 2}
	r2 := r1.clone()

	if r1 != r2 {
		t.Fatalf("range clone failed: got=%#v, want=%v", r2, r1)
	}

	r1.Min = -1
	if r1 == r2 {
		t.Fatalf("range clone did not deep-copy")
	}
}
