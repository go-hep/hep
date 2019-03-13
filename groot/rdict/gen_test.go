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
		ok    bool
	}{
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;go-hep.org/x/hep/hbook.H1D",
			want:  "go-hep.org/x/hep/hbook.H1D",
			ok:    true,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "",
			want:  "go_hep_org::x::hep::hbook::H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;hbook.H1D",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D ",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Rust; stl::hbook::H1D",
			want:  "stl::hbook::H1D",
			ok:    false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := GoName2Cxx(tc.name)
			got, ok := Typename(name, tc.title)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
			if ok != tc.ok {
				t.Fatalf("got=%q, want=%q, ok=%v (want=%v)", got, tc.want, ok, tc.ok)
			}
		})
	}

	if _, ok := Typename("go_hep_org::x::hep::groot::redm::HLV", "Go;go-hep.org/x/hep/groot/redm.Event"); ok {
		t.Fatalf("typename did not fail!")
	}
}

func TestROOTComment(t *testing.T) {
	var g genGoType
	for _, tc := range []struct {
		title string
		meta  string
		doc   string
	}{
		{
			title: "A comment",
			meta:  "",
			doc:   "A comment",
		},
		{
			title: " A comment ",
			meta:  "",
			doc:   "A comment",
		},
		{
			title: "[N]",
			meta:  "[N]",
			doc:   "",
		},
		{
			title: "[N] this is an array. ",
			meta:  "[N]",
			doc:   "this is an array.",
		},
		{
			title: "[-1,1,2]",
			meta:  "[-1,1,2]",
			doc:   "",
		},
		{
			title: "[-1,1,2] a Double32 with min,max,factor",
			meta:  "[-1,1,2]",
			doc:   "a Double32 with min,max,factor",
		},
		{
			title: "[fN][-1,1,2] an array of Double32-s with min,max,factor",
			meta:  "[fN][-1,1,2]",
			doc:   "an array of Double32-s with min,max,factor",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			meta, doc := g.rcomment(tc.title)
			if meta != tc.meta {
				t.Fatalf("meta: got=%q, want=%q", meta, tc.meta)
			}
			if doc != tc.doc {
				t.Fatalf("doc: got=%q, want=%q", doc, tc.doc)
			}
		})
	}
}
