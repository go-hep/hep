// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys_test

import (
	"fmt"
	"testing"

	"go-hep.org/x/hep/groot/rphys"
)

func TestVector2(t *testing.T) {
	p2 := rphys.NewVector2(1, 2)

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "x",
			fct:  p2.X,
			want: 1,
		},
		{
			name: "y",
			fct:  p2.Y,
			want: 2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fct()
			if got != tc.want {
				t.Fatalf("invalid getter value: got=%v, want=%v", got, tc.want)
			}
		})
	}
	if got, want := fmt.Sprintf("%v", p2), "TVector2{1, 2}"; got != want {
		t.Fatalf("invalid stringer value:\ngot= %q\nwant=%q", got, want)
	}

	p2.SetX(-1)
	p2.SetY(-2)

	if got, want := fmt.Sprintf("%v", p2), "TVector2{-1, -2}"; got != want {
		t.Fatalf("invalid stringer value:\ngot= %q\nwant=%q", got, want)
	}

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "x",
			fct:  p2.X,
			want: -1,
		},
		{
			name: "y",
			fct:  p2.Y,
			want: -2,
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
