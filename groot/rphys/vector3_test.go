// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys_test

import (
	"testing"

	"go-hep.org/x/hep/groot/rphys"
)

func TestVector3(t *testing.T) {
	p3 := rphys.NewVector3(1, 2, 3)

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "x",
			fct:  p3.X,
			want: 1,
		},
		{
			name: "y",
			fct:  p3.Y,
			want: 2,
		},
		{
			name: "z",
			fct:  p3.Z,
			want: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fct()
			if got != tc.want {
				t.Fatalf("invalid getter value: got=%v, want=%v", got, tc.want)
			}
		})
	}

	p3.SetX(-1)
	p3.SetY(-2)
	p3.SetZ(-3)

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "x",
			fct:  p3.X,
			want: -1,
		},
		{
			name: "y",
			fct:  p3.Y,
			want: -2,
		},
		{
			name: "z",
			fct:  p3.Z,
			want: -3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fct()
			if got != tc.want {
				t.Fatalf("invalid getter value: got=%v, want=%v", got, tc.want)
			}
		})
	}
}
