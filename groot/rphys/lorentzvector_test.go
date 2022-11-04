// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys_test

import (
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/rphys"
	"go-hep.org/x/hep/internal/diff"
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

func TestLorentzVectorRStreamer(t *testing.T) {
	const (
		deep = true
		want = `key[000]: tlv;1 "A four vector with (-,-,-,+) metric" (TLorentzVector) => "TLorentzVector{P: {10, 20, 30}, E: 40}"
key[001]: tree;1 "my tree title" (TTree)
[000][p4]: {{0 50331648} {{0 50331648} 0 1 2} 3}
[001][p4]: {{0 50331648} {{0 50331648} 1 2 3} 4}
[002][p4]: {{0 50331648} {{0 50331648} 2 3 4} 5}
[003][p4]: {{0 50331648} {{0 50331648} 3 4 5} 6}
[004][p4]: {{0 50331648} {{0 50331648} 4 5 6} 7}
[005][p4]: {{0 50331648} {{0 50331648} 5 6 7} 8}
[006][p4]: {{0 50331648} {{0 50331648} 6 7 8} 9}
[007][p4]: {{0 50331648} {{0 50331648} 7 8 9} 10}
[008][p4]: {{0 50331648} {{0 50331648} 8 9 10} 11}
[009][p4]: {{0 50331648} {{0 50331648} 9 10 11} 12}
`
	)

	for _, fname := range []string{
		"../testdata/tlv-split00.root", // exercizes T{Branch,Leaf}Object
		"../testdata/tlv-split01.root",
		"../testdata/tlv-split99.root",
	} {
		t.Run(fname, func(t *testing.T) {
			got := new(strings.Builder)
			err := rcmd.Dump(got, fname, deep, nil)
			if err != nil {
				t.Fatalf("could not run root-dump: %+v", err)
			}

			if got, want := got.String(), want; got != want {
				t.Fatalf("invalid root-dump output:\n%s", diff.Format(got, want))
			}
		})
	}
}
