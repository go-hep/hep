// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"math"
	"reflect"
	"testing"
)

func TestCoerce(t *testing.T) {
	for _, tc := range []struct {
		v1, v2 any
		w1, w2 any
	}{
		{
			v1: int32(0),
			v2: int32(0),
			w1: int32(0),
			w2: int32(0),
		},
		{
			v1: idealFloat(0),
			v2: int32(0),
			w1: idealFloat(0),
			w2: int32(0),
		},
		{
			v1: int32(0),
			v2: idealInt(0),
			w1: int32(0),
			w2: int32(0),
		},
		{
			v1: nil,
			v2: idealInt(0),
			w1: nil,
			w2: idealInt(0),
		},
		{
			v1: int32(0),
			v2: int64(0),
			w1: int32(0),
			w2: int64(0),
		},
	} {
		t.Run("", func(t *testing.T) {
			{
				v1, v2 := coerce(tc.v1, tc.v2)
				rt1 := reflect.TypeOf(v1)
				rt2 := reflect.TypeOf(v2)
				w1 := reflect.TypeOf(tc.w1)
				w2 := reflect.TypeOf(tc.w2)
				switch {
				case w1 != rt1:
					t.Fatalf("invalid type\ngot1=%v\nwant=%v\n", rt1, w1)
				case w2 != rt2:
					t.Fatalf("invalid type\ngot2=%v\nwant=%v\n", rt2, w2)
				}
			}
			{
				v1, v2 := coerce(tc.v2, tc.v1)
				rt1 := reflect.TypeOf(v1)
				rt2 := reflect.TypeOf(v2)
				w1 := reflect.TypeOf(tc.w2)
				w2 := reflect.TypeOf(tc.w1)
				switch {
				case w1 != rt1:
					t.Fatalf("invalid type\ngot1=%v\nwant=%v\n", rt1, w1)
				case w2 != rt2:
					t.Fatalf("invalid type\ngot2=%v\nwant=%v\n", rt2, w2)
				}
			}
		})
	}
}

func TestCoerce1(t *testing.T) {
	for _, tc := range []struct {
		v1, v2 any
		want   any
	}{
		{
			v1:   nil,
			v2:   "",
			want: nil,
		},
		// idealFloat
		{
			v1:   idealFloat(1),
			v2:   idealFloat(2),
			want: idealFloat(1),
		},
		{
			v1:   idealFloat(1),
			v2:   float32(2),
			want: float32(1),
		},
		{
			v1:   idealFloat(1),
			v2:   float64(2),
			want: float64(1),
		},
		// idealInt
		{
			v1:   idealInt(1),
			v2:   idealFloat(2),
			want: idealFloat(1),
		},
		{
			v1:   idealInt(1),
			v2:   idealInt(2),
			want: idealInt(1),
		},
		{
			v1:   idealInt(1),
			v2:   idealUint(2),
			want: idealUint(1),
		},
		{
			v1:   idealInt(-1),
			v2:   idealUint(2),
			want: idealInt(-1),
		},
		{
			v1:   idealInt(1),
			v2:   float32(2),
			want: float32(1),
		},
		{
			v1:   idealInt(1),
			v2:   float64(2),
			want: float64(1),
		},
		{
			v1:   idealInt(1),
			v2:   int8(2),
			want: int8(1),
		},
		{
			v1:   idealInt(math.MaxInt8 + 1),
			v2:   int8(2),
			want: idealInt(math.MaxInt8 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   int16(2),
			want: int16(1),
		},
		{
			v1:   idealInt(math.MaxInt16 + 1),
			v2:   int16(2),
			want: idealInt(math.MaxInt16 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   int32(2),
			want: int32(1),
		},
		{
			v1:   idealInt(math.MaxInt32 + 1),
			v2:   int32(2),
			want: idealInt(math.MaxInt32 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   int64(2),
			want: int64(1),
		},
		{
			v1:   idealInt(1),
			v2:   uint8(2),
			want: uint8(1),
		},
		{
			v1:   idealInt(math.MaxUint8 + 1),
			v2:   uint8(2),
			want: idealInt(math.MaxUint8 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   uint16(2),
			want: uint16(1),
		},
		{
			v1:   idealInt(math.MaxUint16 + 1),
			v2:   uint16(2),
			want: idealInt(math.MaxUint16 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   uint32(2),
			want: uint32(1),
		},
		{
			v1:   idealInt(math.MaxUint32 + 1),
			v2:   uint32(2),
			want: idealInt(math.MaxUint32 + 1),
		},
		{
			v1:   idealInt(1),
			v2:   uint64(2),
			want: uint64(1),
		},
		// idealUint
		{
			v1:   idealUint(1),
			v2:   idealFloat(2),
			want: idealFloat(1),
		},
		{
			v1:   idealUint(1),
			v2:   idealInt(2),
			want: idealInt(1),
		},
		{
			v1:   idealUint(math.MaxInt64 + 1),
			v2:   idealInt(2),
			want: idealUint(math.MaxInt64 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   idealUint(2),
			want: idealUint(1),
		},
		{
			v1:   idealUint(1),
			v2:   float32(2),
			want: float32(1),
		},
		{
			v1:   idealUint(1),
			v2:   float64(2),
			want: float64(1),
		},
		{
			v1:   idealUint(1),
			v2:   int8(2),
			want: int8(1),
		},
		{
			v1:   idealUint(math.MaxInt8 + 1),
			v2:   int8(2),
			want: idealUint(math.MaxInt8 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   int16(2),
			want: int16(1),
		},
		{
			v1:   idealUint(math.MaxInt16 + 1),
			v2:   int16(2),
			want: idealUint(math.MaxInt16 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   int32(2),
			want: int32(1),
		},
		{
			v1:   idealUint(math.MaxInt32 + 1),
			v2:   int32(2),
			want: idealUint(math.MaxInt32 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   int64(2),
			want: int64(1),
		},
		{
			v1:   idealUint(math.MaxInt64 + 1),
			v2:   int64(2),
			want: idealUint(math.MaxInt64 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   uint8(2),
			want: uint8(1),
		},
		{
			v1:   idealUint(math.MaxUint8 + 1),
			v2:   uint8(2),
			want: idealUint(math.MaxUint8 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   uint16(2),
			want: uint16(1),
		},
		{
			v1:   idealUint(math.MaxUint16 + 1),
			v2:   uint16(2),
			want: idealUint(math.MaxUint16 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   uint32(2),
			want: uint32(1),
		},
		{
			v1:   idealUint(math.MaxUint32 + 1),
			v2:   uint32(2),
			want: idealUint(math.MaxUint32 + 1),
		},
		{
			v1:   idealUint(1),
			v2:   uint64(2),
			want: uint64(1),
		},
	} {
		t.Run("", func(t *testing.T) {
			got := coerce1(tc.v1, tc.v2)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got=%#v (%T), want=%#v (%T)", got, got, tc.want, tc.want)
			}
		})
	}
}

func TestColFromDesc(t *testing.T) {
	for _, tc := range []struct {
		name     string
		etyp     reflect.Type
		kind     reflect.Kind
		count    bool
		size     int
		unsigned bool

		want colDescr
	}{
		// int8
		{
			name:     "int8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    false,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "int8", Len: -1, Nullable: false, Type: reflect.TypeOf(int8(0))},
		},
		{
			name:     "[2]int8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    false,
			size:     2,
			unsigned: false,

			want: colDescr{Name: "[2]int8", Len: 2, Nullable: false, Type: reflect.TypeOf(int8(0))},
		},
		{
			name:     "[]int8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    true,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "[]int8", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(int8(0))},
		},
		// uint8
		{
			name:     "uint8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    false,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "uint8", Len: -1, Nullable: false, Type: reflect.TypeOf(uint8(0))},
		},
		{
			name:     "[2]uint8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    false,
			size:     2,
			unsigned: true,

			want: colDescr{Name: "[2]uint8", Len: 2, Nullable: false, Type: reflect.TypeOf(uint8(0))},
		},
		{
			name:     "[]uint8",
			etyp:     reflect.TypeOf(int8(0)),
			kind:     reflect.Int8,
			count:    true,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "[]uint8", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(uint8(0))},
		},
		// int16
		{
			name:     "int16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    false,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "int16", Len: -1, Nullable: false, Type: reflect.TypeOf(int16(0))},
		},
		{
			name:     "[2]int16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    false,
			size:     2,
			unsigned: false,

			want: colDescr{Name: "[2]int16", Len: 2, Nullable: false, Type: reflect.TypeOf(int16(0))},
		},
		{
			name:     "[]int16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    true,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "[]int16", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(int16(0))},
		},
		// uint16
		{
			name:     "uint16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    false,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "uint16", Len: -1, Nullable: false, Type: reflect.TypeOf(uint16(0))},
		},
		{
			name:     "[2]uint16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    false,
			size:     2,
			unsigned: true,

			want: colDescr{Name: "[2]uint16", Len: 2, Nullable: false, Type: reflect.TypeOf(uint16(0))},
		},
		{
			name:     "[]uint16",
			etyp:     reflect.TypeOf(int16(0)),
			kind:     reflect.Int16,
			count:    true,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "[]uint16", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(uint16(0))},
		},
		// int32
		{
			name:     "int32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    false,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "int32", Len: -1, Nullable: false, Type: reflect.TypeOf(int32(0))},
		},
		{
			name:     "[2]int32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    false,
			size:     2,
			unsigned: false,

			want: colDescr{Name: "[2]int32", Len: 2, Nullable: false, Type: reflect.TypeOf(int32(0))},
		},
		{
			name:     "[]int32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    true,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "[]int32", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(int32(0))},
		},
		// uint32
		{
			name:     "uint32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    false,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "uint32", Len: -1, Nullable: false, Type: reflect.TypeOf(uint32(0))},
		},
		{
			name:     "[2]uint32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    false,
			size:     2,
			unsigned: true,

			want: colDescr{Name: "[2]uint32", Len: 2, Nullable: false, Type: reflect.TypeOf(uint32(0))},
		},
		{
			name:     "[]uint32",
			etyp:     reflect.TypeOf(int32(0)),
			kind:     reflect.Int32,
			count:    true,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "[]uint32", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(uint32(0))},
		},
		// int64
		{
			name:     "int64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    false,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "int64", Len: -1, Nullable: false, Type: reflect.TypeOf(int64(0))},
		},
		{
			name:     "[2]int64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    false,
			size:     2,
			unsigned: false,

			want: colDescr{Name: "[2]int64", Len: 2, Nullable: false, Type: reflect.TypeOf(int64(0))},
		},
		{
			name:     "[]int64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    true,
			size:     0,
			unsigned: false,

			want: colDescr{Name: "[]int64", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(int64(0))},
		},
		// uint64
		{
			name:     "uint64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    false,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "uint64", Len: -1, Nullable: false, Type: reflect.TypeOf(uint64(0))},
		},
		{
			name:     "[2]uint64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    false,
			size:     2,
			unsigned: true,

			want: colDescr{Name: "[2]uint64", Len: 2, Nullable: false, Type: reflect.TypeOf(uint64(0))},
		},
		{
			name:     "[]uint64",
			etyp:     reflect.TypeOf(int64(0)),
			kind:     reflect.Int64,
			count:    true,
			size:     0,
			unsigned: true,

			want: colDescr{Name: "[]uint64", Len: math.MaxInt64, Nullable: true, Type: reflect.TypeOf(uint64(0))},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := colDescrFrom(tc.name, tc.etyp, tc.kind, tc.count, tc.size, tc.unsigned)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid column descriptor:\ngot= %#v\nwant=%#v", got, tc.want)
			}
		})
	}
}
