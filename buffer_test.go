// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestRBuffer(t *testing.T) {
	data := make([]byte, 32)
	r := NewRBuffer(data, nil, 0)

	if got, want := r.Len(), int64(32); got != want {
		t.Fatalf("got len=%v. want=%v", got, want)
	}
	start := r.Pos()
	if start != 0 {
		t.Fatalf("got start=%v. want=%v", start, 0)
	}

	_ = r.ReadI16()
	if r.Err() != nil {
		t.Fatalf("error reading int16: %v", r.Err())
	}

	pos := r.Pos()
	if pos != 2 {
		t.Fatalf("got pos=%v. want=%v", pos, 16)
	}

	pos = 0
	data = make([]byte, 2*(2+4+8))
	r = NewRBuffer(data, nil, 0)
	for _, n := range []int{2, 4, 8} {
		beg := r.Pos()
		if beg != pos {
			t.Errorf("pos[%d] error: got=%d, want=%d\n", n, beg, pos)
			t.Errorf("len: %d size: %d\n", r.r.Len(), r.r.Size())
		}
		switch n {
		case 2:
			_ = r.ReadI16()
			_ = r.ReadU16()
		case 4:
			_ = r.ReadI32()
			_ = r.ReadU32()
		case 8:
			_ = r.ReadI64()
			_ = r.ReadU64()
		}
		end := r.Pos()
		pos += int64(2 * n)

		if got, want := end-beg, int64(2*n); got != want {
			t.Errorf("%d-bytes: got=%d. want=%d\n", n, got, want)
		}
	}
}

func TestReadRBuffer(t *testing.T) {
	for _, test := range []struct {
		name string
		file string
		want Object
	}{
		{
			name: "TNamed",
			want: &named{name: "my-name", title: "my-title"},
		},
		{
			name: "TList",
			want: &tlist{
				name: "list-name",
				objs: []Object{
					&named{name: "n0", title: "t0"},
					&named{name: "n1", title: "t1"},
				},
			},
		},
		{
			name: "TObjArray",
			want: &objarray{
				named: named{name: "my-objs"},
				arr: []Object{
					&named{name: "n0", title: "t0"},
					&named{name: "n1", title: "t1"},
					&named{name: "n2", title: "t2"},
				},
				last: 2,
			},
		},
		/*
			{
				name: "TList",
				file: "testdata/tlist-tsi.dat",
				want: &tlist{
					name: "",
					objs: []Object{},
				},
			},
		*/
		/*
			{
				name: "TStreamerInfo",
				want: &tstreamerInfo{},
			},
		*/
	} {
		test := test
		file := test.file
		if file == "" {
			file = "testdata/" + strings.ToLower(test.name) + ".dat"
		}
		t.Run("read-buffer="+test.name, func(t *testing.T) {
			testReadRBuffer(t, test.name, file, test.want)
		})
	}
}

func testReadRBuffer(t *testing.T, name, file string, want Object) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	r := NewRBuffer(data, nil, 0)
	obj := Factory.get(want.Class())().Interface().(Object)
	err = obj.(ROOTUnmarshaler).UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error:\ngot= %+v\nwant=%+v\n", obj, want)
	}
}
