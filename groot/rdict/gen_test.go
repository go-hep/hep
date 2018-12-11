// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import "testing"

func TestGoName2Cxx(t *testing.T) {
	for _, tc := range []struct {
		name, want string
	}{
		{
			name: "go-hep.org/x/hep/hbook.H1D",
			want: "go_hep_org::x::hep::hbook::H1D",
		},
		{
			name: "go-hep.org/x.H1D",
			want: "go_hep_org::x::H1D",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := GoName2Cxx(tc.name)
			if tc.want != got {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestTypename(t *testing.T) {
	for _, tc := range []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;go-hep.org/x/hep/hbook.H1D",
			want:  "go-hep.org/x/hep/hbook.H1D",
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "",
			want:  "go_hep_org::x::hep::hbook::H1D",
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;hbook.H1D",
			want:  "hbook.H1D",
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D",
			want:  "hbook.H1D",
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D ",
			want:  "hbook.H1D",
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Rust; stl::hbook::H1D",
			want:  "stl::hbook::H1D",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := GoName2Cxx(tc.name)
			got, _ := Typename(name, tc.title)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}
