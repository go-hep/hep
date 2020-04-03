// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

func TestStreamerOf(t *testing.T) {
	for _, tc := range []struct {
		typ  reflect.Type
		want rbytes.StreamerInfo
	}{
		{
			typ: reflect.TypeOf((*struct1)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct1", "struct1"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{StreamerElement{
						named:  *rbase.NewNamed("Name", ""),
						etype:  rmeta.TString,
						esize:  16,
						offset: 0,
						ename:  "golang::string",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Bool", ""),
						etype:  rmeta.Bool,
						esize:  1,
						offset: 0,
						ename:  "golang::bool",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I8", ""),
						etype:  rmeta.Char,
						esize:  1,
						offset: 0,
						ename:  "golang::int8",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I16", ""),
						etype:  rmeta.Short,
						esize:  2,
						offset: 0,
						ename:  "golang::int16",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I32", ""),
						etype:  rmeta.Int,
						esize:  4,
						offset: 0,
						ename:  "golang::int32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I64", ""),
						etype:  rmeta.Long,
						esize:  8,
						offset: 0,
						ename:  "golang::int64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U8", ""),
						etype:  rmeta.UChar,
						esize:  1,
						offset: 0,
						ename:  "golang::uint8",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U16", ""),
						etype:  rmeta.UShort,
						esize:  2,
						offset: 0,
						ename:  "golang::uint16",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U32", ""),
						etype:  rmeta.UInt,
						esize:  4,
						offset: 0,
						ename:  "golang::uint32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U64", ""),
						etype:  rmeta.ULong,
						esize:  8,
						offset: 0,
						ename:  "golang::uint64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F32", ""),
						etype:  rmeta.Float,
						esize:  4,
						offset: 0,
						ename:  "golang::float32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F64", ""),
						etype:  rmeta.Double,
						esize:  8,
						offset: 0,
						ename:  "golang::float64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Cxx::MyFloat64", ""),
						etype:  rmeta.Double,
						esize:  8,
						offset: 0,
						ename:  "golang::float64",
					}},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct2)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct2", "struct2"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("V1", ""),
						etype:  rmeta.Any,
						esize:  72,
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
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Names", ""),
						etype:  rmeta.OffsetL + rmeta.TString,
						esize:  160,
						offset: 0,
						ename:  "golang::string",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("Bools", ""),
						etype:  rmeta.OffsetL + rmeta.Bool,
						esize:  10,
						offset: 0,
						ename:  "golang::bool",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I8s", ""),
						etype:  rmeta.OffsetL + rmeta.Char,
						esize:  10,
						offset: 0,
						ename:  "golang::int8",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I16s", ""),
						etype:  rmeta.OffsetL + rmeta.Short,
						esize:  20,
						offset: 0,
						ename:  "golang::int16",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I32s", ""),
						etype:  rmeta.OffsetL + rmeta.Int,
						esize:  40,
						offset: 0,
						ename:  "golang::int32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("I64s", ""),
						etype:  rmeta.OffsetL + rmeta.Long,
						esize:  80,
						offset: 0,
						ename:  "golang::int64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U8s", ""),
						etype:  rmeta.OffsetL + rmeta.UChar,
						esize:  10,
						offset: 0,
						ename:  "golang::uint8",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U16s", ""),
						etype:  rmeta.OffsetL + rmeta.UShort,
						esize:  20,
						offset: 0,
						ename:  "golang::uint16",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U32s", ""),
						etype:  rmeta.OffsetL + rmeta.UInt,
						esize:  40,
						offset: 0,
						ename:  "golang::uint32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("U64s", ""),
						etype:  rmeta.OffsetL + rmeta.ULong,
						esize:  80,
						offset: 0,
						ename:  "golang::uint64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F32s", ""),
						etype:  rmeta.OffsetL + rmeta.Float,
						esize:  40,
						offset: 0,
						ename:  "golang::float32",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("F64s", ""),
						etype:  rmeta.OffsetL + rmeta.Double,
						esize:  80,
						offset: 0,
						ename:  "golang::float64",
					}},
					&StreamerBasicType{StreamerElement{
						named:  *rbase.NewNamed("S1s", ""),
						etype:  rmeta.OffsetL + rmeta.Any,
						esize:  720,
						offset: 0,
						ename:  "struct1",
					}},
				},
			},
		},
		{
			typ: reflect.TypeOf((*struct4)(nil)).Elem(),
			want: &StreamerInfo{
				named:  *rbase.NewNamed("struct4", "struct4"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("Names", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::string>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("Bools", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::bool>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("I8s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::int8>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("I16s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::int16>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("I32s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::int32>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("I64s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::int64>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("U8s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::uint8>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("U16s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::uint16>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("U32s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::uint32>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("U64s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::uint64>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("F32s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::float32>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("F64s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<golang::float64>",
					}},
					&StreamerObjectAny{StreamerElement{
						named:  *rbase.NewNamed("S1s", ""),
						etype:  rmeta.Any,
						esize:  3 * int32(ptrSize),
						offset: 0,
						ename:  "golang::slice<struct1>",
					}},
				},
			},
		},
	} {
		t.Run(tc.want.Name(), func(t *testing.T) {
			got := StreamerOf(ctx, tc.typ)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid streamer info.\ngot:\n%v\nwant:\n%v", got, tc.want)
			}
		})
	}
}

var (
	ctx = newStreamerStore(nil)
)

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
	Float64 float64 `groot:"Cxx::MyFloat64"`
}

type struct2 struct {
	V1 struct1
}

type struct3 struct {
	Names [10]string
	Bools [10]bool
	I8s   [10]int8
	I16s  [10]int16
	I32s  [10]int32
	I64s  [10]int64
	U8s   [10]uint8
	U16s  [10]uint16
	U32s  [10]uint32
	U64s  [10]uint64
	F32s  [10]float32
	F64s  [10]float64
	S1s   [10]struct1
}

type struct4 struct {
	Names []string
	Bools []bool
	I8s   []int8
	I16s  []int16
	I32s  []int32
	I64s  []int64
	U8s   []uint8
	U16s  []uint16
	U32s  []uint32
	U64s  []uint64
	F32s  []float32
	F64s  []float64
	S1s   []struct1
}

func TestNameOfStructField(t *testing.T) {
	for _, tc := range []struct {
		field reflect.StructField
		want  string
	}{
		{
			field: reflect.StructField{},
			want:  "",
		},
		{
			field: reflect.StructField{Name: "f1"},
			want:  "f1",
		},
		{
			field: reflect.StructField{Name: "F1"},
			want:  "F1",
		},
		{
			field: reflect.StructField{Name: "f1", Tag: `groot:"F1"`},
			want:  "F1",
		},
		{
			field: reflect.StructField{Name: "F1", Tag: `groot:"f1"`},
			want:  "f1",
		},
		{
			field: reflect.StructField{Name: "F1", Tag: `groot:"f1[N]"`},
			want:  "f1",
		},
		{
			field: reflect.StructField{Name: "F1", Tag: `groot:"f1[10]"`},
			want:  "f1",
		},
	} {
		t.Run(fmt.Sprintf("%v", tc.field), func(t *testing.T) {
			got := nameOf(tc.field)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestBuildStreamerInfo(t *testing.T) {
	type builtinsT1 struct {
		F0  bool
		F1  uint8
		F2  uint16
		F3  uint32
		F4  uint64
		F5  int8
		F6  int16
		F7  int32
		F8  int64
		F9  float32
		F10 float64
		F11 root.Float16
		F12 root.Double32
		F13 string
	}

	for _, tc := range []struct {
		name string
		v    interface{}
		wfct func(w *rbytes.WBuffer, ptr interface{}) // FIXME
	}{
		{
			name: "builtins",
			v: builtinsT1{
				F0: true,
				F1: 1, F2: 2, F3: 3, F4: 4,
				F5: -5, F6: -6, F7: -7, F8: -8,
				F9: 9.9, F10: 10.10,
				F11: -11, F12: -12,
				F13: "hello\nworld",
			},
			wfct: func(w *rbytes.WBuffer, ptr interface{}) {
				v := ptr.(*builtinsT1)
				w.WriteBool(v.F0)
				w.WriteU8(v.F1)
				w.WriteU16(v.F2)
				w.WriteU32(v.F3)
				w.WriteU64(v.F4)
				w.WriteI8(v.F5)
				w.WriteI16(v.F6)
				w.WriteI32(v.F7)
				w.WriteI64(v.F8)
				w.WriteF32(v.F9)
				w.WriteF64(v.F10)
				w.WriteF16(v.F11, nil)
				w.WriteD32(v.F12, nil)
				w.WriteString(v.F13)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				rt  = reflect.TypeOf(tc.v)
				ctx = newStreamerStore(nil)
				si  = StreamerOf(ctx, rt).(*StreamerInfo)
			)

			err := si.BuildStreamers()
			if err != nil {
				t.Fatalf("could not build streamers: %+v", err)
			}

			ptr := reflect.New(rt)
			ptr.Elem().Set(reflect.ValueOf(tc.v))

			wbuf := rbytes.NewWBuffer(nil, nil, 0, ctx)
			//	enc, err := si.NewEncoder(rbytes.ObjectWise, wbuf)
			//	if err != nil {
			//		t.Fatalf("could not create streamer encoder: %+v", err)
			//	}
			//
			//	err = enc.EncodeROOT(ptr.Interface())
			//	if err != nil {
			//		t.Fatalf("could not encode value %T: %+v", tc.v, err)
			//	}

			tc.wfct(wbuf, ptr.Interface()) // FIXME(sbinet): use encoder

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, ctx)

			dec, err := si.NewDecoder(rbytes.ObjectWise, rbuf)
			if err != nil {
				t.Fatalf("could not create streamer decoder: %+v", err)
			}

			got := reflect.New(reflect.TypeOf(tc.v))
			err = dec.DecodeROOT(got.Interface())
			if err != nil {
				t.Fatalf("could not decode value %T: %+v", tc.v, err)
			}

			if got, want := got.Elem().Interface(), tc.v; !reflect.DeepEqual(got, want) {
				t.Fatalf(
					"invalid enc/dec round-trip:\ngot= %#v\nwant=%#v",
					got, want,
				)
			}
		})
	}
}
