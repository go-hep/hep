// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"reflect"
	"testing"
)

func TestFreeList(t *testing.T) {
	var list freeList

	checkSize := func(size int) {
		t.Helper()
		if got, want := len(list), size; got != want {
			t.Fatalf("got=%d, want=%d", got, want)
		}
	}
	checkSegment := func(free *freeSegment, want freeSegment) {
		t.Helper()
		if free == nil {
			t.Fatalf("expected a valid free segment")
		}
		if *free != want {
			t.Fatalf("got=%#v, want=%#v", *free, want)
		}
	}

	list.add(0, 1)
	checkSize(1)

	list.add(3, 10)
	checkSize(2)

	free := list.add(13, 20)
	checkSize(3)
	checkSegment(free, freeSegment{13, 20})

	free = list.add(12, 22)
	checkSize(3)
	checkSegment(free, freeSegment{12, 22})

	if got, want := list, (freeList{
		{0, 1},
		{3, 10},
		{12, 22},
	}); !reflect.DeepEqual(got, want) {
		t.Fatalf("error\ngot = %v\nwant= %v", got, want)
	}

	free = list.add(15, 20)
	checkSize(3)
	checkSegment(free, freeSegment{12, 22})

	free = list.add(40, 50)
	checkSize(4)
	checkSegment(free, freeSegment{40, 50})

	free = list.add(39, 40)
	checkSize(4)
	checkSegment(free, freeSegment{39, 50})

	free = list.add(37, 38)
	checkSize(4)
	checkSegment(free, freeSegment{37, 50})

	list.add(55, 60)
	list.add(65, 70)
	free = list.add(56, 66)
	checkSize(5)
	checkSegment(free, freeSegment{55, 70})

	free = list.add(54, 71)
	checkSize(5)
	checkSegment(free, freeSegment{54, 71})

	for _, tc := range []struct {
		list []freeSegment
		want []freeSegment
		free freeList
	}{
		{
			list: nil,
			want: nil,
			free: nil,
		},
		{
			list: []freeSegment{{0, 1}, {1, 2}},
			want: []freeSegment{{0, 1}, {0, 2}},
			free: freeList{{0, 2}},
		},
		{
			list: []freeSegment{{10, 12}, {10, 13}},
			want: []freeSegment{{10, 12}, {10, 13}},
			free: freeList{{10, 13}},
		},
	} {
		t.Run("", func(t *testing.T) {
			var list freeList
			for i, v := range tc.list {
				free := list.add(v.first, v.last)
				if !reflect.DeepEqual(*free, tc.want[i]) {
					t.Fatalf("error:\ngot[%d] = %#v\nwant[%d]= %#v\n",
						i, *free, i, tc.want[i],
					)
				}
			}
			if !reflect.DeepEqual(list, tc.free) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v\n", list, tc.free)
			}
		})
	}
}

func TestFreeListBest(t *testing.T) {
	for _, tc := range []struct {
		name   string
		nbytes int64
		list   freeList
		want   *freeSegment
	}{
		{
			name:   "empty",
			nbytes: 0,
			list:   nil,
			want:   nil,
		},
		{
			name:   "empty-list",
			nbytes: 10,
			list:   nil,
			want:   nil,
		},
		{
			name:   "exact-match",
			nbytes: 10,
			list:   freeList{{0, 1}, {10, 20 - 1}},
			want:   &freeSegment{10, 20 - 1},
		},
		{
			name:   "match",
			nbytes: 1,
			list:   freeList{{0, 10}},
			want:   &freeSegment{0, 10},
		},
		{
			name:   "match",
			nbytes: 10,
			list:   freeList{{0, 1}, {10, 20 + 4 + 1}},
			want:   &freeSegment{10, 20 + 4 + 1},
		},
		{
			name:   "big-file",
			nbytes: 10,
			list:   freeList{{0, 1}},
			want:   &freeSegment{0, 1000000001},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.list.best(tc.nbytes)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error\ngot = %#v\nwant= %#v\n", got, tc.want)
			}
		})
	}
}

func TestFreeListLast(t *testing.T) {
	for _, tc := range []struct {
		list freeList
		want *freeSegment
	}{
		{
			list: nil,
			want: nil,
		},
		{
			list: freeList{},
			want: nil,
		},
		{
			list: freeList{{0, kStartBigFile}},
			want: &freeSegment{0, kStartBigFile},
		},
		{
			list: freeList{{0, 10}, {12, kStartBigFile}},
			want: &freeSegment{12, kStartBigFile},
		},
	} {
		t.Run("", func(t *testing.T) {
			got := tc.list.last()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error\ngot = %#v\nwant= %#v\n", got, tc.want)
			}
		})
	}
}
