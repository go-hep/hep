// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"reflect"
	"testing"
)

func TestSpanSplit(t *testing.T) {
	mk := func(beg, end int64) span {
		return span{
			off: beg,
			len: end - beg,
		}
	}
	from := func(vs ...int64) []span {
		if len(vs) == 0 {
			return nil
		}
		o := make([]span, 0, len(vs)/2)
		for i := 0; i < len(vs); i += 2 {
			o = append(o, mk(vs[i], vs[i+1]))
		}
		return o
	}
	for _, tc := range []struct {
		spans []span
		sp    span
		want  []span
	}{
		{
			spans: from(),
			sp:    span{0, 2},
			want:  from(0, 2),
		},
		// 1-1 intersects
		{
			spans: from(2, 4),
			sp:    mk(0, 3),
			want:  from(0, 2),
		},
		{
			spans: from(2, 4),
			sp:    mk(3, 5),
			want:  from(4, 5),
		},
		{
			spans: from(2, 5),
			sp:    mk(3, 4),
			want:  from(),
		},
		{
			spans: from(2, 4),
			sp:    mk(0, 5),
			want:  from(0, 2, 4, 5),
		},
		{
			spans: from(2, 4),
			sp:    mk(0, 4),
			want:  from(0, 2),
		},
		{
			spans: from(2, 4),
			sp:    mk(2, 5),
			want:  from(4, 5),
		},
		{
			spans: from(2, 4),
			sp:    mk(0, 2),
			want:  from(0, 2),
		},
		{
			spans: from(2, 4),
			sp:    mk(0, 1),
			want:  from(0, 1),
		},
		{
			spans: from(2, 4),
			sp:    mk(4, 6),
			want:  from(4, 6),
		},
		{
			spans: from(2, 4),
			sp:    mk(5, 6),
			want:  from(5, 6),
		},
		//
		{
			spans: from(0, 4, 5, 7),
			sp:    mk(0, 7),
			want:  from(4, 5),
		},
		{
			spans: from(0, 4, 5, 7),
			sp:    mk(0, 6),
			want:  from(4, 5),
		},
		{
			spans: from(0, 4, 5, 7),
			sp:    mk(7, 9),
			want:  from(7, 9),
		},
		// 2-1 intersects
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(3, 7),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(3, 10),
			want:  from(4, 6, 8, 10),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(3, 5),
			want:  from(4, 5),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 5),
			want:  from(4, 5),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 7),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 10),
			want:  from(4, 6, 8, 10),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(0, 5),
			want:  from(0, 1, 4, 5),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(0, 7),
			want:  from(0, 1, 4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(0, 10),
			want:  from(0, 1, 4, 6, 8, 10),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(4, 5),
			want:  from(4, 5),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(4, 7),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(4, 10),
			want:  from(4, 6, 8, 10),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 8),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 4),
			want:  from(),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(1, 6),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(4, 6),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(4, 8),
			want:  from(4, 6),
		},
		{
			spans: from(1, 4, 6, 8),
			sp:    mk(6, 8),
			want:  from(),
		},
	} {
		t.Run("", func(t *testing.T) {
			got := split(tc.sp, tc.spans)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid split:\ngot= %+v\nwant=%+v", got, tc.want)
			}
		})
	}
}

func TestSpanAdd(t *testing.T) {
	mk := func(beg, end int64) span {
		return span{
			off: beg,
			len: end - beg,
		}
	}
	from := func(vs ...int64) []span {
		if len(vs) == 0 {
			return nil
		}
		o := make([]span, 0, len(vs)/2)
		for i := 0; i < len(vs); i += 2 {
			o = append(o, mk(vs[i], vs[i+1]))
		}
		return o
	}
	for _, tc := range []struct {
		spans []span
		sp    span
		want  []span
	}{
		{
			spans: from(),
			sp:    mk(0, 10),
			want:  from(0, 10),
		},
		{
			spans: from(),
			sp:    mk(9, 10),
			want:  from(9, 10),
		},
		{
			spans: from(1, 3, 4, 6),
			sp:    mk(3, 4),
			want:  from(1, 6),
		},
		{
			spans: from(1, 3, 4, 6),
			sp:    mk(0, 1),
			want:  from(0, 3, 4, 6),
		},
		{
			spans: from(1, 3, 4, 6),
			sp:    mk(6, 10),
			want:  from(1, 3, 4, 10),
		},
		{
			spans: from(1, 3, 4, 6),
			sp:    mk(7, 10),
			want:  from(1, 3, 4, 6, 7, 10),
		},
		{
			spans: from(2, 3, 4, 6),
			sp:    mk(0, 1),
			want:  from(0, 1, 2, 3, 4, 6),
		},
	} {
		t.Run("", func(t *testing.T) {
			got := make(spans, len(tc.spans))
			copy(got, tc.spans)
			got.add(tc.sp)
			if !reflect.DeepEqual(got, spans(tc.want)) {
				t.Fatalf("invalid span-add:\ngot= %+v\nwant=%+v", got, tc.want)
			}
		})
	}
}
