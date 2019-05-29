// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"testing"
)

func TestLeafDims(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want []int
	}{
		{s: "Leaf", want: nil},
		{s: "Leaf[]", want: []int{-1}},
		{s: "Leaf[20]", want: []int{20}},
		{s: "Leaf[2000]", want: []int{2000}},
		{s: "Leaf[1][2]", want: []int{1, 2}},
		{s: "Leaf[2][]", want: []int{2, -1}},
		{s: "Leaf[1][2][3]", want: []int{1, 2, 3}},
	} {
		t.Run(tc.s, func(t *testing.T) {
			dims := leafDims(tc.s)
			if !reflect.DeepEqual(dims, tc.want) {
				t.Fatalf("invalid dims. got=%#v, want=%#v", dims, tc.want)
			}
		})
	}
}
