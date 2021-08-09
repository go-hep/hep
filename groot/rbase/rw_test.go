// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestWRBuffer(t *testing.T) {
	for _, tc := range []struct {
		name string
		want rtests.ROOTer
	}{
		{
			name: "TObject",
			want: &Object{ID: 0x0, Bits: 0x3000000},
		},
		{
			name: "TObject",
			want: &Object{ID: 0x1, Bits: 0x3000001},
		},
		{
			name: "TUUID",
			want: &UUID{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
				10, 11, 12, 13, 14, 15,
			},
		},
		{
			name: "TNamed",
			want: &Named{obj: Object{ID: 0x0, Bits: 0x3000000}, name: "my-name", title: "my-title"},
		},
		{
			name: "TNamed",
			want: &Named{
				obj:  Object{ID: 0x0, Bits: 0x3000000},
				name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			want: &Named{
				obj:  Object{ID: 0x0, Bits: 0x3500000},
				name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			want: &Named{
				obj:   Object{ID: 0x0, Bits: 0x3000000},
				name:  strings.Repeat("*", 256),
				title: "my-title",
			},
		},
		{
			name: "TObjString",
			want: &ObjString{
				obj: Object{ID: 0x0, Bits: 0x3000008},
				str: "tobjstring-string",
			},
		},
		{
			name: "TProcessID",
			want: &ProcessID{
				named: Named{obj: Object{ID: 0x0, Bits: 0x3000000}, name: "my-name", title: "my-title"},
			},
		},
		{
			name: "TRef",
			want: &Ref{
				obj: Object{ID: 0x0, Bits: 0x3000000},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			class := tc.want.Class()
			obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
			{
				rbuf.SetErr(io.EOF)
				err = obj.UnmarshalROOT(rbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
				rbuf.SetErr(nil)
			}
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}

func TestWriteWBuffer(t *testing.T) {
	for _, test := range rwBufferCases {
		t.Run("write-buffer="+test.file, func(t *testing.T) {
			testWriteWBuffer(t, test.name, test.file, test.want)
		})
	}
}

func testWriteWBuffer(t *testing.T, name, file string, want interface{}) {
	rdata, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	{
		wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
		wbuf.SetErr(io.EOF)
		_, err := want.(rbytes.Marshaler).MarshalROOT(wbuf)
		if err == nil {
			t.Fatalf("expected an error")
		}
		if err != io.EOF {
			t.Fatalf("got=%v, want=%v", err, io.EOF)
		}
	}

	w := rbytes.NewWBuffer(nil, nil, 0, nil)
	_, err = want.(rbytes.Marshaler).MarshalROOT(w)
	if err != nil {
		t.Fatal(err)
	}
	wdata := w.Bytes()

	r := rbytes.NewRBuffer(wdata, nil, 0, nil)
	obj := rtypes.Factory.Get(name)().Interface().(rbytes.Unmarshaler)
	{
		r.SetErr(io.EOF)
		err = obj.UnmarshalROOT(r)
		if err == nil {
			t.Fatalf("expected an error")
		}
		if err != io.EOF {
			t.Fatalf("got=%v, want=%v", err, io.EOF)
		}
		r.SetErr(nil)
	}
	err = obj.UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(file+".new", wdata, 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error: %q\ngot= %+v\nwant=%+v\ngot= %+v\nwant=%+v", file, wdata, rdata, obj, want)
	}

	os.Remove(file + ".new")
}
func TestReadRBuffer(t *testing.T) {
	for _, test := range rwBufferCases {
		test := test
		file := test.file
		if file == "" {
			file = "../testdata/" + strings.ToLower(test.name) + ".dat"
		}
		t.Run("read-buffer="+file, func(t *testing.T) {
			testReadRBuffer(t, test.name, file, test.want)
		})
	}
}

func testReadRBuffer(t *testing.T, name, file string, want interface{}) {
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	r := rbytes.NewRBuffer(data, nil, 0, nil)
	obj := rtypes.Factory.Get(name)().Interface().(rbytes.Unmarshaler)
	err = obj.UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error: %q\ngot= %#v\nwant=%#v\n", file, obj, want)
	}
}

var rwBufferCases = []struct {
	name string
	file string
	want rbytes.Unmarshaler
}{
	{
		name: "TObject",
		file: "../testdata/tobject.dat",
		want: &Object{ID: 0x0, Bits: 0x3000000},
	},
	{
		name: "TNamed",
		file: "../testdata/tnamed.dat",
		want: &Named{obj: Object{ID: 0x0, Bits: 0x3000000}, name: "my-name", title: "my-title"},
	},
	{
		name: "TNamed",
		file: "../testdata/tnamed-cmssw.dat",
		want: &Named{
			obj:  Object{ID: 0x0, Bits: 0x3000000},
			name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
		},
	},
	{
		name: "TNamed",
		file: "../testdata/tnamed-cmssw-2.dat",
		want: &Named{
			obj:  Object{ID: 0x0, Bits: 0x3500000},
			name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
		},
	},
	{
		name: "TNamed",
		file: "../testdata/tnamed-long-string.dat",
		want: &Named{
			obj:   Object{ID: 0x0, Bits: 0x3000000},
			name:  strings.Repeat("*", 256),
			title: "my-title",
		},
	},
	{
		name: "TObjString",
		file: "../testdata/tobjstring.dat",
		want: &ObjString{
			obj: Object{ID: 0x0, Bits: 0x3000008},
			str: "tobjstring-string",
		},
	},
}

func TestUUIDv1(t *testing.T) {
	uuid := UUID([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
	_, err := wbuf.Write(uuid[:])
	if err != nil {
		t.Fatalf("could not serialize v1: %+v", err)
	}

	var (
		got  UUID
		want = uuid
	)
	rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
	err = got.UnmarshalROOTv1(rbuf)
	if err != nil {
		t.Fatalf("could not unserialize v1: %+v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid UUID-v1 round-trip:\ngot= %v\nwant=%v\n", got, want)
	}
}
