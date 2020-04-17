// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys_test

import (
	"testing"

	"go-hep.org/x/hep/groot/rphys"
)

func TestLorentzVector(t *testing.T) {
	p4 := rphys.NewLorentzVector(1, 2, 3, 4)

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "px",
			fct:  p4.Px,
			want: 1,
		},
		{
			name: "py",
			fct:  p4.Py,
			want: 2,
		},
		{
			name: "pz",
			fct:  p4.Pz,
			want: 3,
		},
		{
			name: "e",
			fct:  p4.E,
			want: 4,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fct()
			if got != tc.want {
				t.Fatalf("invalid getter value: got=%v, want=%v", got, tc.want)
			}
		})
	}

	p4.SetPxPyPzE(-1, -2, -3, 44)

	for _, tc := range []struct {
		name string
		fct  func() float64
		want float64
	}{
		{
			name: "px",
			fct:  p4.Px,
			want: -1,
		},
		{
			name: "py",
			fct:  p4.Py,
			want: -2,
		},
		{
			name: "pz",
			fct:  p4.Pz,
			want: -3,
		},
		{
			name: "e",
			fct:  p4.E,
			want: 44,
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
