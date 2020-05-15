// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import "testing"

func TestBkReaderFindBaskets(t *testing.T) {
	for _, tc := range []struct {
		spans []rspan
		span  [2]int64
		want  [2]int
	}{
		{
			spans: []rspan{{beg: 0, end: 1}},
			span:  [2]int64{0, 1},
			want:  [2]int{0, 1},
		},
		{
			spans: []rspan{{beg: 0, end: 2}},
			span:  [2]int64{0, 1},
			want:  [2]int{0, 1},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}},
			span:  [2]int64{0, 5},
			want:  [2]int{0, 2},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}},
			span:  [2]int64{0, 2},
			want:  [2]int{0, 1},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}},
			span:  [2]int64{2, 3},
			want:  [2]int{1, 2},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}, {beg: 5, end: 10}},
			span:  [2]int64{2, 3},
			want:  [2]int{1, 2},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}, {beg: 5, end: 10}},
			span:  [2]int64{2, 5},
			want:  [2]int{1, 2},
		},
		{
			spans: []rspan{{beg: 0, end: 2}, {beg: 2, end: 5}, {beg: 5, end: 10}},
			span:  [2]int64{2, 6},
			want:  [2]int{1, 3},
		},
	} {
		t.Run("", func(t *testing.T) {
			bkr := &bkreader{spans: tc.spans}
			ibeg, iend := bkr.findBaskets(tc.span[0], tc.span[1])
			got := [2]int{ibeg, iend}
			if got, want := got, tc.want; got != want {
				t.Fatalf("invalid range for span %#v, got=%v, want=%v",
					tc.span, got, want,
				)
			}
		})
	}
}
