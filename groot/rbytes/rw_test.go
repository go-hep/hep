// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes_test

import (
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

func TestReadWriteObjectAny(t *testing.T) {
	var (
		v1 = rbase.NewNamed("name-1", "title-1")
		v2 = rbase.NewNamed("name-2", "title-2")
		v3 *rbase.Named
	)

	for _, tc := range []struct {
		name string
		vs   []*rbase.Named
	}{
		{"1213", []*rbase.Named{v1, v2, v1, v3}},
		{"1123", []*rbase.Named{v1, v1, v2, v3}},
		{"2113", []*rbase.Named{v2, v1, v1, v3}},
		{"2131", []*rbase.Named{v2, v1, v3, v1}},

		{"12213", []*rbase.Named{v1, v2, v2, v1, v3}},
		{"12123", []*rbase.Named{v1, v2, v1, v2, v3}},
		{"11223", []*rbase.Named{v1, v1, v2, v2, v3}},
		{"21213", []*rbase.Named{v2, v1, v2, v1, v3}},
		{"22131", []*rbase.Named{v2, v2, v1, v3, v1}},
	} {
		t.Run(tc.name, func(t *testing.T) {

			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			for i, v := range tc.vs {
				wbuf.WriteObjectAny(v)
				if err := wbuf.Err(); err != nil {
					t.Fatalf("could not write named[%d]=%v: %+v", i, v, err)
				}
			}
			if err := wbuf.Err(); err != nil {
				t.Fatalf("could not fill wbuffer: %+v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			for i := range tc.vs {
				var v *rbase.Named
				obj := rbuf.ReadObjectAny()
				if err := rbuf.Err(); err != nil {
					t.Fatalf("could not read object[%d]: %+v", i, err)
				}

				if obj != nil {
					v = obj.(*rbase.Named)
				}

				if got, want := v, tc.vs[i]; !reflect.DeepEqual(got, want) {
					t.Fatalf("invalid named[%d] value:\ngot = %v\nwant= %v\n", i, got, want)
				}
			}
		})
	}
}

func TestRWStrings(t *testing.T) {
	want := []string{"", "x", "", "xx", "", "xxx"}
	wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
	for i, str := range want {
		wbuf.WriteString(str)
		if err := wbuf.Err(); err != nil {
			t.Errorf("could not write string #%d: %+v", i, err)
		}
	}
	rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
	for i := range want {
		got := rbuf.ReadString()
		if got != want[i] {
			t.Errorf("invalid string at %d: got=%q, want=%q", i, got, want[i])
		}
	}
}

func TestRWStdVecStrings(t *testing.T) {
	want := [][]string{
		{"", "x", "", "xx", "", "xxx"},
		{"x", "", "xx", "", "xxx"},
		{"", "x", "", "xx", "", "xxx"},
		{"x", "", "xx", "", "xxx", strings.Repeat("1!", 256)},
	}
	wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
	for i, str := range want {
		wbuf.WriteStdVectorStrs(str)
		if err := wbuf.Err(); err != nil {
			t.Errorf("could not write string #%d: %+v", i, err)
		}
	}
	rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
	for i := range want {
		var got []string
		rbuf.ReadStdVectorStrs(&got)
		if got, want := got, want[i]; !reflect.DeepEqual(got, want) {
			t.Errorf("invalid string at %d: got=%q, want=%q", i, got, want)
		}
	}
}

func TestRWFloat16(t *testing.T) {
	makeElm := func(title string) rbytes.StreamerElement {
		elm := rdict.Element{
			Name: *rbase.NewNamed("f16", title),
			Type: rmeta.Float16,
		}.New()
		return &elm
	}

	for _, tc := range []struct {
		name string
		v    root.Float16
		want root.Float16
	}{
		{
			name: "",
			v:    42,
			want: 42,
		},
		{
			name: "[0,42]",
			v:    42,
			want: 42,
		},
		{
			name: "[43,44]",
			v:    42,
			want: 43,
		},
		{
			name: "[0,41]",
			v:    42,
			want: 41,
		},
		{
			name: "[10,10,2]",
			v:    14,
			want: 14,
		},
		{
			name: "[-10,-10,3]",
			v:    -10,
			want: -10,
		},
		{
			name: "[10,10,20]",
			v:    10,
			want: 10,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var elm rbytes.StreamerElement
			if tc.name != "" {
				elm = makeElm(tc.name)
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			wbuf.WriteArrayF16([]root.Float16{tc.v}, elm)
			if err := wbuf.Err(); err != nil {
				t.Fatalf("could not write f16=%v: %+v", tc.v, err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			got := make([]root.Float16, 1)
			rbuf.ReadArrayF16(got, elm)
			if err := rbuf.Err(); err != nil {
				t.Fatalf("could not read f16=%v: %+v", tc.v, err)
			}

			if got, want := got[0], tc.want; got != want {
				t.Fatalf("invalid r/w round-trip: got=%v, want=%v", got, want)
			}
		})
	}
}

func TestRWDouble32(t *testing.T) {
	makeElm := func(title string) rbytes.StreamerElement {
		elm := rdict.Element{
			Name: *rbase.NewNamed("d32", title),
			Type: rmeta.Double32,
		}.New()
		return &elm
	}

	for _, tc := range []struct {
		name string
		v    root.Double32
		want root.Double32
	}{
		{
			name: "",
			v:    42,
			want: 42,
		},
		{
			name: "[0,42]",
			v:    42,
			want: 42,
		},
		{
			name: "[43,44]",
			v:    42,
			want: 43,
		},
		{
			name: "[0,41]",
			v:    42,
			want: 41,
		},
		{
			name: "[10,10,2]",
			v:    14,
			want: 14,
		},
		{
			name: "[-10,-10,3]",
			v:    -10,
			want: -10,
		},
		{
			name: "[10,10,20]",
			v:    10,
			want: 10,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var elm rbytes.StreamerElement
			if tc.name != "" {
				elm = makeElm(tc.name)
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			wbuf.WriteArrayD32([]root.Double32{tc.v}, elm)
			if err := wbuf.Err(); err != nil {
				t.Fatalf("could not write d32=%v: %+v", tc.v, err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			got := make([]root.Double32, 1)
			rbuf.ReadArrayD32(got, elm)
			if err := rbuf.Err(); err != nil {
				t.Fatalf("could not read d32=%v: %+v", tc.v, err)
			}

			if got, want := got[0], tc.want; got != want {
				t.Fatalf("invalid r/w round-trip: got=%v, want=%v", got, want)
			}
		})
	}
}
