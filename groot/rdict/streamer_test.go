// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
)

func TestStreamerOf(t *testing.T) {
	for _, tc := range []struct {
		typ  reflect.Type
		want rbytes.StreamerInfo
	}{
		{
			typ: reflect.TypeOf((*struct1)(nil)).Elem(),
			want: &StreamerInfo{
				named: *rbase.NewNamed("struct1", "struct1"),
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
				named: *rbase.NewNamed("struct2", "struct2"),
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
				named: *rbase.NewNamed("struct3", "struct3"),
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
				named: *rbase.NewNamed("struct4", "struct4"),
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

var (
	ctx = newStreamerStore(nil)
)
