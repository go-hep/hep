// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

func TestIsTObject(t *testing.T) {
	for _, tc := range []struct {
		typ  reflect.Type
		want bool
	}{
		{
			typ:  reflect.TypeOf(rbase.Object{}),
			want: false,
		},
		{
			typ:  reflect.TypeOf(&rbase.Object{}),
			want: true,
		},
		{
			typ:  reflect.TypeOf((*rbase.ObjString)(nil)),
			want: true,
		},
		{
			typ:  reflect.TypeOf(int(0)),
			want: false,
		},
		{
			typ:  reflect.TypeOf((*rcont.List)(nil)),
			want: true,
		},
		{
			typ:  reflect.TypeOf(tobject{}),
			want: true,
		},
		{
			typ:  reflect.TypeOf(&tobject{}),
			want: true,
		},
	} {
		t.Run(fmt.Sprintf("type=%v", tc.typ), func(t *testing.T) {
			got := isTObject(tc.typ)
			if got != tc.want {
				t.Fatalf("invalid is-TObject: got=%v, want=%v", got, tc.want)
			}
		})
	}
}

func TestHasCount(t *testing.T) {
	for _, tc := range []struct {
		tag  reflect.StructTag
		ok   bool
		want string
	}{
		{
			tag: ``,
		},
		{
			tag: `groot:"Name"`,
		},
		{
			tag:  `groot:"Name[N]"`,
			ok:   true,
			want: "N",
		},
		{
			tag:  `groot:"Name[n1_2_3]"`,
			ok:   true,
			want: "n1_2_3",
		},
		{
			tag:  `groot:"Name[ N ]"`,
			ok:   true,
			want: "N",
		},
		{
			tag: `groot:"Name[]"`,
		},
		{
			tag: `groot:"Name[ ]"`,
		},
		{
			tag: `groot:"Name[1]"`,
		},
		{
			tag: `groot:"Name[N;1]"`,
		},
		{
			tag: `groot:"Name[N-1]"`,
		},
		{
			tag: `groot:"Name[N'1]"`,
		},
		{
			tag: `groot:"Name[1N]"`,
		},
		{
			tag: `groot:"Name[1,2,3]"`,
		},
	} {
		t.Run(string(tc.tag), func(t *testing.T) {
			got, ok := hasCount(reflect.StructField{Tag: tc.tag})
			if ok != tc.ok {
				t.Fatalf("invalid ok: got=%v, want=%v", ok, tc.ok)
			}
			if got != tc.want {
				t.Fatalf("invalid count: got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestTypenameOf(t *testing.T) {
	for _, tc := range []struct {
		name   string
		typ    reflect.Type
		want   string
		panics string
	}{
		{
			name: "bool",
			typ:  reflect.TypeOf(false),
			want: "bool",
		},
		{
			name: "*bool",
			typ:  reflect.TypeOf((*bool)(nil)),
			want: "bool*",
		},
		{
			name: "int8",
			typ:  reflect.TypeOf(int8(0)),
			want: "int8_t",
		},
		{
			name: "int16",
			typ:  reflect.TypeOf(int16(0)),
			want: "int16_t",
		},
		{
			name: "int32",
			typ:  reflect.TypeOf(int32(0)),
			want: "int32_t",
		},
		{
			name: "int64",
			typ:  reflect.TypeOf(int64(0)),
			want: "int64_t",
		},
		{
			name: "uint8",
			typ:  reflect.TypeOf(uint8(0)),
			want: "uint8_t",
		},
		{
			name: "uint16",
			typ:  reflect.TypeOf(uint16(0)),
			want: "uint16_t",
		},
		{
			name: "uint32",
			typ:  reflect.TypeOf(uint32(0)),
			want: "uint32_t",
		},
		{
			name: "uint64",
			typ:  reflect.TypeOf(uint64(0)),
			want: "uint64_t",
		},
		{
			name: "float32",
			typ:  reflect.TypeOf(float32(0)),
			want: "float",
		},
		{
			name: "float64",
			typ:  reflect.TypeOf(float64(0)),
			want: "double",
		},
		{
			name: "Float16_t",
			typ:  reflect.TypeOf(root.Float16(0)),
			want: "Float16_t",
		},
		{
			name: "Double32_t",
			typ:  reflect.TypeOf(root.Double32(0)),
			want: "Double32_t",
		},
		{
			name: "string",
			typ:  reflect.TypeOf(""),
			want: "string",
		},
		{
			name: "[2]float32",
			typ:  reflect.TypeOf([2]float32{}),
			want: "float[2]",
		},
		{
			name: "[2]string",
			typ:  reflect.TypeOf([2]string{}),
			want: "string[2]",
		},
		{
			name: "[2][3]string",
			typ:  reflect.TypeOf([2][3]string{}),
			want: "string[2][3]",
		},
		{
			name: "vector<string>",
			typ:  reflect.TypeOf([]string{}),
			want: "vector<string>",
		},
		{
			name: "vector<int32_t>",
			typ:  reflect.TypeOf([]int32{}),
			want: "vector<int32_t>",
		},
		{
			name: "vector<vector<int32_t> >",
			typ:  reflect.TypeOf([][]int32{}),
			want: "vector<vector<int32_t> >",
		},
		{
			name: "vector<vector<vector<int32_t> > >",
			typ:  reflect.TypeOf([][][]int32{}),
			want: "vector<vector<vector<int32_t> > >",
		},
		{
			name: "vector<Float16_t>",
			typ:  reflect.TypeOf([]root.Float16{}),
			want: "vector<Float16_t>",
		},
		{
			name: "vector<Double32_t>",
			typ:  reflect.TypeOf([]root.Double32{}),
			want: "vector<Double32_t>",
		},
		{
			name:   "empty-struct",
			typ:    reflect.TypeOf(struct{}{}),
			panics: "rdict: invalid reflect type struct {}",
		},
		{
			name: "struct-event",
			typ: reflect.TypeOf(func() interface{} {
				type Event struct{}
				return Event{}
			}()),
			want: "Event",
		},
		{
			name: "TObjString*",
			typ:  reflect.TypeOf((*rbase.ObjString)(nil)),
			want: "TObjString*",
		},
		{
			name: "TObjString",
			typ:  reflect.TypeOf((*rbase.ObjString)(nil)).Elem(),
			want: "TObjString",
		},
		{
			name: "vector<TObjString>",
			typ:  reflect.TypeOf([]rbase.ObjString{}),
			want: "vector<TObjString>",
		},
		{
			name: "vector<TObjString*>",
			typ:  reflect.TypeOf([]*rbase.ObjString{}),
			want: "vector<TObjString*>",
		},
		{
			name: "struct0*",
			typ:  reflect.TypeOf((*struct0)(nil)),
			want: "struct0*",
		},
		{
			name: "struct0",
			typ:  reflect.TypeOf((*struct0)(nil)).Elem(),
			want: "struct0",
		},
		{
			name: "vector<struct0>",
			typ:  reflect.TypeOf([]struct0{}),
			want: "vector<struct0>",
		},
		{
			name: "vector<struct0*>",
			typ:  reflect.TypeOf([]*struct0{}),
			want: "vector<struct0*>",
		},
		{
			name: "tobject",
			typ:  reflect.TypeOf(tobject{}),
			want: "tobject",
		},
		{
			name: "*tobject",
			typ:  reflect.TypeOf(&tobject{}),
			want: "tobject*",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic (%q)", tc.panics)
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message:\ngot= %q\nwant=%q",
							got, want,
						)
					}
				}()
			}

			got := typenameOf(tc.typ)
			if got, want := got, tc.want; got != want {
				t.Fatalf("invalid typename:\ngot= %q\nwant=%q", got, want)
			}
		})
	}
}

func TestStreamerOf(t *testing.T) {
	for _, tc := range []struct {
		typ    reflect.Type
		want   rbytes.StreamerInfo
		panics string
	}{
		{
			typ: reflect.TypeOf(&rbase.ObjString{}),
			want: func() rbytes.StreamerInfo {
				si, ok := StreamerInfos.Get("TObjString", -1)
				if !ok {
					t.Fatalf("could not get streamer for TObjString")
				}
				return si
			}(),
		},
		{
			typ: reflect.TypeOf((*struct0)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct0", "struct0"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectPointer{
						StreamerElement: StreamerElement{
							named:  *rbase.NewNamed("ObjPtr", ""),
							etype:  rmeta.ObjectP,
							esize:  diskPtrSize,
							offset: 0,
							ename:  "TObjString*",
						},
					},
					&StreamerObjectAnyPointer{
						StreamerElement: StreamerElement{
							named:  *rbase.NewNamed("UsrPtr", ""),
							etype:  rmeta.AnyP,
							esize:  diskPtrSize,
							offset: 0,
							ename:  "struct1*",
						},
					},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct1)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct1", "struct1"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerSTLstring{
						StreamerSTL: StreamerSTL{
							StreamerElement: StreamerElement{
								named:  *rbase.NewNamed("Name", ""),
								etype:  rmeta.Streamer,
								esize:  sizeOfStdString,
								offset: 0,
								ename:  "string",
							},
						},
					},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Bool", ""),
						etype:  rmeta.Bool,
						esize:  1,
						offset: 0,
						ename:  "bool",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I8", ""),
						etype:  rmeta.Int8,
						esize:  1,
						offset: 0,
						ename:  "int8_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I16", ""),
						etype:  rmeta.Int16,
						esize:  2,
						offset: 0,
						ename:  "int16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I32", ""),
						etype:  rmeta.Int32,
						esize:  4,
						offset: 0,
						ename:  "int32_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I64", ""),
						etype:  rmeta.Int64,
						esize:  8,
						offset: 0,
						ename:  "int64_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U8", ""),
						etype:  rmeta.Uint8,
						esize:  1,
						offset: 0,
						ename:  "uint8_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U16", ""),
						etype:  rmeta.Uint16,
						esize:  2,
						offset: 0,
						ename:  "uint16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U32", ""),
						etype:  rmeta.Uint32,
						esize:  4,
						offset: 0,
						ename:  "uint32_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U64", ""),
						etype:  rmeta.Uint64,
						esize:  8,
						offset: 0,
						ename:  "uint64_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F32", ""),
						etype:  rmeta.Float32,
						esize:  4,
						offset: 0,
						ename:  "float",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F64", ""),
						etype:  rmeta.Float64,
						esize:  8,
						offset: 0,
						ename:  "double",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F16", ""),
						etype:  rmeta.Float16,
						esize:  4,
						offset: 0,
						ename:  "Float16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("D32", ""),
						etype:  rmeta.Double32,
						esize:  8,
						offset: 0,
						ename:  "Double32_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Cxx::MyFloat64", ""),
						etype:  rmeta.Float64,
						esize:  8,
						offset: 0,
						ename:  "double",
					}},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct2)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct2", "struct2"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("V1", ""),
						etype:  rmeta.Any,
						esize:  88,
						offset: 0,
						ename:  "struct1",
					}},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct3)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct3", "struct3"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{StreamerElement{
						named:  *rbase.NewNamed("Names", ""),
						etype:  rmeta.OffsetL + rmeta.TString,
						esize:  10 * sizeOfTString,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "TString",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Bools", ""),
						etype:  rmeta.OffsetL + rmeta.Bool,
						esize:  10 * 1,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "bool",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I8s", ""),
						etype:  rmeta.OffsetL + rmeta.Int8,
						esize:  10 * 1,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "int8_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I16s", ""),
						etype:  rmeta.OffsetL + rmeta.Int16,
						esize:  10 * 2,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "int16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I32s", ""),
						etype:  rmeta.OffsetL + rmeta.Int32,
						esize:  10 * 4,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "int32_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I64s", ""),
						etype:  rmeta.OffsetL + rmeta.Int64,
						esize:  10 * 8,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "int64_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U8s", ""),
						etype:  rmeta.OffsetL + rmeta.Uint8,
						esize:  10 * 1,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "uint8_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U16s", ""),
						etype:  rmeta.OffsetL + rmeta.Uint16,
						esize:  10 * 2,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "uint16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U32s", ""),
						etype:  rmeta.OffsetL + rmeta.Uint32,
						esize:  10 * 4,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "uint32_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U64s", ""),
						etype:  rmeta.OffsetL + rmeta.Uint64,
						esize:  10 * 8,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "uint64_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F32s", ""),
						etype:  rmeta.OffsetL + rmeta.Float32,
						esize:  10 * 4,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "float",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F64s", ""),
						etype:  rmeta.OffsetL + rmeta.Float64,
						esize:  10 * 8,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "double",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F16s", ""),
						etype:  rmeta.OffsetL + rmeta.Float16,
						esize:  10 * 4,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "Float16_t",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("D32s", ""),
						etype:  rmeta.OffsetL + rmeta.Double32,
						esize:  10 * 8,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "Double32_t",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("S1s", ""),
						etype:  rmeta.OffsetL + rmeta.Any,
						esize:  10 * 88,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "struct1",
					}},
					&StreamerObject{StreamerElement{
						named:  *rbase.NewNamed("ObjStrs", ""),
						etype:  rmeta.OffsetL + rmeta.Object,
						esize:  10 * sizeOfTObjString,
						offset: 0,
						arrlen: 10,
						arrdim: 1,
						maxidx: [5]int32{10, 0, 0, 0, 0},
						ename:  "TObjString",
					}},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct4)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct4", "struct4"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("Names", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<string>",
					}, rmeta.STLvector, rmeta.STLstring),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("Bools", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<bool>",
					}, rmeta.STLvector, rmeta.Bool),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I8s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<int8_t>",
					}, rmeta.STLvector, rmeta.Int8),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<int16_t>",
					}, rmeta.STLvector, rmeta.Int16),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<int32_t>",
					}, rmeta.STLvector, rmeta.Int32),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<int64_t>",
					}, rmeta.STLvector, rmeta.Int64),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U8s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<uint8_t>",
					}, rmeta.STLvector, rmeta.Uint8),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<uint16_t>",
					}, rmeta.STLvector, rmeta.Uint16),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<uint32_t>",
					}, rmeta.STLvector, rmeta.Uint32),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<uint64_t>",
					}, rmeta.STLvector, rmeta.Uint64),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<float>",
					}, rmeta.STLvector, rmeta.Float32),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<double>",
					}, rmeta.STLvector, rmeta.Float64),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<Float16_t>",
					}, rmeta.STLvector, rmeta.Float16),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("D32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<Double32_t>",
					}, rmeta.STLvector, rmeta.Double32),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("S1s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<struct1>",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("ObjStrs", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<TObjString>",
					}, rmeta.STLvector, rmeta.Object),
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct5)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct5", "struct5"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("N", ""),
							Type:  rmeta.Int32, // should be rmeta.Counter
							Size:  4,
							EName: "int32_t",
						}.New(),
					},
					NewStreamerLoop(
						Element{
							Name:  *rbase.NewNamed("Names", "[N]"),
							Size:  4,
							EName: "TString*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("Bools", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Bool,
							Size:  1,
							EName: "bool*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("I8s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int8,
							Size:  1,
							EName: "int8_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("I16s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int16,
							Size:  2,
							EName: "int16_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("I32s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int32,
							Size:  4,
							EName: "int32_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("I64s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int64,
							Size:  8,
							EName: "int64_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("U8s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint8,
							Size:  1,
							EName: "uint8_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("U16s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint16,
							Size:  2,
							EName: "uint16_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("U32s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint32,
							Size:  4,
							EName: "uint32_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("U64s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint64,
							Size:  8,
							EName: "uint64_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("F32s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Float32,
							Size:  4,
							EName: "float*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("F64s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Float64,
							Size:  8,
							EName: "double*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("F16s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Float16,
							Size:  4,
							EName: "Float16_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("D32s", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Double32,
							Size:  8,
							EName: "Double32_t*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerLoop(
						Element{
							Name:  *rbase.NewNamed("S1s", "[N]"),
							Size:  8,
							EName: "struct1*",
						}.New(), 1, "N", "struct5",
					),
					NewStreamerLoop(
						Element{
							Name:  *rbase.NewNamed("ObjStrs", "[N]"),
							Size:  8,
							EName: "TObjString*",
						}.New(), 1, "N", "struct5",
					),
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct6)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct6", "struct6"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("Names", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<string> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("Bools", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<bool> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I8s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<int8_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<int16_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<int32_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("I64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<int64_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U8s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<uint8_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<uint16_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<uint32_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("U64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<uint64_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<float> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F64s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<double> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("F16s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<Float16_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("D32s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<Double32_t> >",
					}, rmeta.STLvector, rmeta.Any),
					NewCxxStreamerSTL(StreamerElement{
						named:  *rbase.NewNamed("S1s", ""),
						etype:  rmeta.Streamer,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "vector<vector<struct1> >",
					}, rmeta.STLvector, rmeta.Any),
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct7)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct7", "struct7"),
				clsver: 1,
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("ArrI32", ""),
						etype:  rmeta.OffsetL + rmeta.Int32,
						esize:  1 * 2 * 3 * 4 * 5 * 4,
						offset: 0,
						arrlen: 1,
						arrdim: 5,
						maxidx: [5]int32{1, 2, 3, 4, 5},
						ename:  "int32_t",
					}},
					&StreamerString{StreamerElement{
						named:  *rbase.NewNamed("ArrStr", ""),
						etype:  rmeta.OffsetL + rmeta.TString,
						esize:  1 * 2 * 3 * 4 * 5 * sizeOfTString,
						offset: 0,
						arrlen: 1,
						arrdim: 5,
						maxidx: [5]int32{1, 2, 3, 4, 5},
						ename:  "TString",
					}},
					&StreamerObject{StreamerElement{
						named:  *rbase.NewNamed("ArrObj", ""),
						etype:  rmeta.OffsetL + rmeta.Object,
						esize:  1 * 2 * 3 * 4 * 5 * 40,
						offset: 0,
						arrlen: 1,
						arrdim: 5,
						maxidx: [5]int32{1, 2, 3, 4, 5},
						ename:  "TObjString",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("ArrUsr", ""),
						etype:  rmeta.OffsetL + rmeta.Any,
						esize:  1 * 2 * 3 * 4 * 5 * 88,
						offset: 0,
						arrlen: 1,
						arrdim: 5,
						maxidx: [5]int32{1, 2, 3, 4, 5},
						ename:  "struct1",
					}},
				},
			},
		},
		{
			// FIXME(sbinet): add support for maps.
			typ: reflect.TypeOf(panicFIXMEStruct0{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicFIXMEStruct0", "panicFIXMEStruct0"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Map, type=map[int32]int32, kind=map)`,
		},
		{
			// FIXME(sbinet): add support for interfaces?
			typ: reflect.TypeOf(panicFIXMEStruct1{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicFIXMEStruct1", "panicFIXMEStruct1"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Iface, type=error, kind=interface)`,
		},
		{
			// FIXME(sbinet): add support for complex-64
			typ: reflect.TypeOf(panicFIXMEStruct2{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicFIXMEStruct2", "panicFIXMEStruct2"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=C64, type=complex64, kind=complex64)`,
		},
		{
			// FIXME(sbinet): add support for complex-128
			typ: reflect.TypeOf(panicFIXMEStruct3{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicFIXMEStruct3", "panicFIXMEStruct3"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=C128, type=complex128, kind=complex128)`,
		},
		{
			typ: reflect.TypeOf(panicStruct0{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct0", "panicStruct0"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Chan, type=chan int32, kind=chan)`,
		},
		{
			typ: reflect.TypeOf(panicStruct1{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct1", "panicStruct1"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Int, type=int, kind=int)`,
		},
		{
			typ: reflect.TypeOf(panicStruct2{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct2", "panicStruct2"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Uint, type=uint, kind=uint)`,
		},
		{
			typ: reflect.TypeOf(panicStruct3{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct3", "panicStruct3"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Uintptr, type=uintptr, kind=uintptr)`,
		},
		{
			typ: reflect.TypeOf(panicStruct4{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct4", "panicStruct4"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Unsafe, type=unsafe.Pointer, kind=unsafe.Pointer)`,
		},
		{
			typ: reflect.TypeOf(panicStruct5{}),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("panicStruct5", "panicStruct5"),
				clsver: 1,
			},
			panics: `rdict: invalid struct field (name=Func, type=func(), kind=func)`,
		},
	} {
		t.Run(tc.want.Name(), func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic (%s)", tc.panics)
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf(
							"invalid panic message:\ngot= %s\nwant=%s",
							got, want,
						)
					}
				}()
			}

			ctx := newStreamerStore(StreamerInfos)
			got := StreamerOf(ctx, tc.typ)
			if !reflect.DeepEqual(got, tc.want) {
				egot := got.Elements()
				ewat := tc.want.Elements()
				for i := range egot {
					if !reflect.DeepEqual(egot[i], ewat[i]) {
						t.Logf("i=%d\ngot= %#v\nwant=%#v", i, egot[i], ewat[i])
					}
				}
				t.Fatalf("invalid streamer info.\ngot:\n%v\nwant:\n%v", got, tc.want)
			}
		})
	}
}

type struct0 struct {
	TObjPtr *rbase.ObjString `groot:"ObjPtr"`
	TUsrPtr *struct1         `groot:"UsrPtr"`
}

type struct1 struct {
	Name    string
	Bool    bool
	I8      int8
	I16     int16
	I32     int32
	I64     int64
	U8      uint8
	U16     uint16
	U32     uint32
	U64     uint64
	F32     float32
	F64     float64
	F16     root.Float16
	D32     root.Double32
	Float64 float64 `groot:"Cxx::MyFloat64"`
}

type struct2 struct {
	V1 struct1
}

type struct3 struct {
	Names   [10]string
	Bools   [10]bool
	I8s     [10]int8
	I16s    [10]int16
	I32s    [10]int32
	I64s    [10]int64
	U8s     [10]uint8
	U16s    [10]uint16
	U32s    [10]uint32
	U64s    [10]uint64
	F32s    [10]float32
	F64s    [10]float64
	F16s    [10]root.Float16
	D32s    [10]root.Double32
	S1s     [10]struct1
	ObjStrs [10]rbase.ObjString
}

type struct4 struct {
	Names   []string
	Bools   []bool
	I8s     []int8
	I16s    []int16
	I32s    []int32
	I64s    []int64
	U8s     []uint8
	U16s    []uint16
	U32s    []uint32
	U64s    []uint64
	F32s    []float32
	F64s    []float64
	F16s    []root.Float16
	D32s    []root.Double32
	S1s     []struct1
	ObjStrs []rbase.ObjString
}

type struct5 struct {
	N       int32
	Names   []string          `groot:"Names[N]"`
	Bools   []bool            `groot:"Bools[N]"`
	I8s     []int8            `groot:"I8s[N]"`
	I16s    []int16           `groot:"I16s[N]"`
	I32s    []int32           `groot:"I32s[N]"`
	I64s    []int64           `groot:"I64s[N]"`
	U8s     []uint8           `groot:"U8s[N]"`
	U16s    []uint16          `groot:"U16s[N]"`
	U32s    []uint32          `groot:"U32s[N]"`
	U64s    []uint64          `groot:"U64s[N]"`
	F32s    []float32         `groot:"F32s[N]"`
	F64s    []float64         `groot:"F64s[N]"`
	F16s    []root.Float16    `groot:"F16s[N]"`
	D32s    []root.Double32   `groot:"D32s[N]"`
	S1s     []struct1         `groot:"S1s[N]"`
	ObjStrs []rbase.ObjString `groot:"ObjStrs[N]"`
}

type struct6 struct {
	Names [][]string
	Bools [][]bool
	I8s   [][]int8
	I16s  [][]int16
	I32s  [][]int32
	I64s  [][]int64
	U8s   [][]uint8
	U16s  [][]uint16
	U32s  [][]uint32
	U64s  [][]uint64
	F32s  [][]float32
	F64s  [][]float64
	F16s  [][]root.Float16
	D32s  [][]root.Double32
	S1s   [][]struct1
}

type struct7 struct {
	ArrI32 [1][2][3][4][5]int32           `groot:"ArrI32[1][2][3][4][5]"`
	ArrStr [1][2][3][4][5]string          `groot:"ArrStr[1][2][3][4][5]"`
	ArrObj [1][2][3][4][5]rbase.ObjString `groot:"ArrObj[1][2][3][4][5]"`
	ArrUsr [1][2][3][4][5]struct1         `groot:"ArrUsr[1][2][3][4][5]"`
}

type panicFIXMEStruct0 struct {
	Map map[int32]int32 `groot:"Map"`
}

type panicFIXMEStruct1 struct {
	Iface error
}

type panicFIXMEStruct2 struct {
	C64 complex64
}

type panicFIXMEStruct3 struct {
	C128 complex128
}

type panicStruct0 struct {
	Chan chan int32
}

type panicStruct1 struct {
	Int int
}

type panicStruct2 struct {
	Uint uint
}

type panicStruct3 struct {
	Uintptr uintptr
}

type panicStruct4 struct {
	Unsafe unsafe.Pointer
}

type panicStruct5 struct {
	Func func()
}

type tobject struct{}

func (tobject) Class() string { return "tobject" }

var (
	_ root.Object = (*tobject)(nil)
)
