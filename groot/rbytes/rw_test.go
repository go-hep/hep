// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
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
				err := wbuf.WriteObjectAny(v)
				if err != nil {
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
