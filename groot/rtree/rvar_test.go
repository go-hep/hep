// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func TestNewReadVars(t *testing.T) {
	f, err := riofs.Open("../testdata/leaves.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(Tree)

	vars := NewReadVars(tree)
	want := []ReadVar{
		{Name: "B", Leaf: "B", Value: new(bool)},
		{Name: "Str", Leaf: "Str", Value: new(string)},
		{Name: "I8", Leaf: "I8", Value: new(int8)},
		{Name: "I16", Leaf: "I16", Value: new(int16)},
		{Name: "I32", Leaf: "I32", Value: new(int32)},
		{Name: "I64", Leaf: "I64", Value: new(int64)},
		{Name: "G64", Leaf: "G64", Value: new(int64)},
		{Name: "U8", Leaf: "U8", Value: new(uint8)},
		{Name: "U16", Leaf: "U16", Value: new(uint16)},
		{Name: "U32", Leaf: "U32", Value: new(uint32)},
		{Name: "U64", Leaf: "U64", Value: new(uint64)},
		{Name: "UGG", Leaf: "UGG", Value: new(uint64)},
		{Name: "F32", Leaf: "F32", Value: new(float32)},
		{Name: "F64", Leaf: "F64", Value: new(float64)},
		{Name: "D16", Leaf: "D16", Value: new(root.Float16)},
		{Name: "D32", Leaf: "D32", Value: new(root.Double32)},
		// arrays
		{Name: "ArrBs", Leaf: "ArrBs", Value: new([10]bool)},
		{Name: "ArrI8", Leaf: "ArrI8", Value: new([10]int8)},
		{Name: "ArrI16", Leaf: "ArrI16", Value: new([10]int16)},
		{Name: "ArrI32", Leaf: "ArrI32", Value: new([10]int32)},
		{Name: "ArrI64", Leaf: "ArrI64", Value: new([10]int64)},
		{Name: "ArrG64", Leaf: "ArrG64", Value: new([10]int64)},
		{Name: "ArrU8", Leaf: "ArrU8", Value: new([10]uint8)},
		{Name: "ArrU16", Leaf: "ArrU16", Value: new([10]uint16)},
		{Name: "ArrU32", Leaf: "ArrU32", Value: new([10]uint32)},
		{Name: "ArrU64", Leaf: "ArrU64", Value: new([10]uint64)},
		{Name: "ArrUGG", Leaf: "ArrUGG", Value: new([10]uint64)},
		{Name: "ArrF32", Leaf: "ArrF32", Value: new([10]float32)},
		{Name: "ArrF64", Leaf: "ArrF64", Value: new([10]float64)},
		{Name: "ArrD16", Leaf: "ArrD16", Value: new([10]root.Float16)},
		{Name: "ArrD32", Leaf: "ArrD32", Value: new([10]root.Double32)},
		// slices
		{Name: "N", Leaf: "N", Value: new(int32)},
		{Name: "SliBs", Leaf: "SliBs", Value: new([]bool)},
		{Name: "SliI8", Leaf: "SliI8", Value: new([]int8)},
		{Name: "SliI16", Leaf: "SliI16", Value: new([]int16)},
		{Name: "SliI32", Leaf: "SliI32", Value: new([]int32)},
		{Name: "SliI64", Leaf: "SliI64", Value: new([]int64)},
		{Name: "SliG64", Leaf: "SliG64", Value: new([]int64)},
		{Name: "SliU8", Leaf: "SliU8", Value: new([]uint8)},
		{Name: "SliU16", Leaf: "SliU16", Value: new([]uint16)},
		{Name: "SliU32", Leaf: "SliU32", Value: new([]uint32)},
		{Name: "SliU64", Leaf: "SliU64", Value: new([]uint64)},
		{Name: "SliUGG", Leaf: "SliUGG", Value: new([]uint64)},
		{Name: "SliF32", Leaf: "SliF32", Value: new([]float32)},
		{Name: "SliF64", Leaf: "SliF64", Value: new([]float64)},
		{Name: "SliD16", Leaf: "SliD16", Value: new([]root.Float16)},
		{Name: "SliD32", Leaf: "SliD32", Value: new([]root.Double32)},
	}

	n := len(want)
	if len(vars) < n {
		n = len(vars)
	}

	for i := 0; i < n; i++ {
		got := vars[i]
		if got.Name != want[i].Name {
			t.Fatalf("invalid read-var name[%d]: got=%q, want=%q", i, got.Name, want[i].Name)
		}
		if got.Leaf != want[i].Leaf {
			t.Fatalf("invalid read-var (name=%q) leaf-name[%d]: got=%q, want=%q", got.Name, i, got.Leaf, want[i].Leaf)
		}
		if got, want := reflect.TypeOf(got.Value), reflect.TypeOf(want[i].Value); got != want {
			t.Fatalf("invalid read-var (name=%q) type[%d]: got=%v, want=%v", vars[i].Name, i, got, want)
		}
	}

	if len(want) != len(vars) {
		t.Fatalf("invalid lengths. got=%d, want=%d", len(vars), len(want))
	}
}

func TestReadVarsFromStruct(t *testing.T) {
	for _, tc := range []struct {
		name   string
		ptr    interface{}
		want   []ReadVar
		panics string
	}{
		{
			name: "not-ptr",
			ptr: struct {
				I32 int32
			}{},
			panics: "rtree: expect a pointer value, got struct { I32 int32 }",
		},
		{
			name:   "not-ptr-to-struct",
			ptr:    new(int32),
			panics: "rtree: expect a pointer to struct value, got *int32",
		},
		{
			name: "struct-with-int",
			ptr: &struct {
				I32 int
				F32 float32
				Str string
			}{},
			panics: "rtree: invalid field type for \"I32\": int",
		},
		{
			name: "struct-with-map", // FIXME(sbinet)
			ptr: &struct {
				Map map[int32]string
			}{},
			panics: "rtree: invalid field type for \"Map\": map[int32]string (not yet supported)",
		},
		{
			name: "invalid-struct-tag",
			ptr: &struct {
				N int32 `groot:"N[42]"`
			}{},
			panics: "rtree: invalid field type for \"N\", or invalid struct-tag \"N[42]\": int32",
		},
		{
			name: "simple",
			ptr: &struct {
				I32 int32
				F32 float32
				Str string
			}{},
			want: []ReadVar{{Name: "I32"}, {Name: "F32"}, {Name: "Str"}},
		},
		{
			name: "simple-with-unexported",
			ptr: &struct {
				I32 int32
				F32 float32
				val float32
				Str string
			}{},
			want: []ReadVar{{Name: "I32"}, {Name: "F32"}, {Name: "Str"}},
		},
		{
			name: "slices",
			ptr: &struct {
				N      int32
				NN     int64
				SliF32 []float32 `groot:"F32s[N]"`
				SliF64 []float64 `groot:"F64s[NN]"`
			}{},
			want: []ReadVar{
				{Name: "N"},
				{Name: "NN"},
				{Name: "F32s", count: "N"},
				{Name: "F64s", count: "NN"},
			},
		},
		{
			name: "slices-no-count",
			ptr: &struct {
				F1 int32
				F2 []float32 `groot:"F2[F1]"`
				X3 []float64 `groot:"F3"`
				F4 []float64
			}{},
			want: []ReadVar{
				{Name: "F1"},
				{Name: "F2", count: "F1"},
				{Name: "F3"},
				{Name: "F4"},
			},
		},
		{
			name: "arrays",
			ptr: &struct {
				N     int32 `groot:"n"`
				Arr01 [10]float64
				Arr02 [10][10]float64
				Arr03 [10][10][10]float64
				Arr11 [10]float64         `groot:"arr11[10]"`
				Arr12 [10][10]float64     `groot:"arr12[10][10]"`
				Arr13 [10][10][10]float64 `groot:"arr13[10][10][10]"`
				Arr14 [10][10][10]float64 `groot:"arr14"`
			}{},
			want: []ReadVar{
				{Name: "n"},
				{Name: "Arr01"},
				{Name: "Arr02"},
				{Name: "Arr03"},
				{Name: "arr11"},
				{Name: "arr12"},
				{Name: "arr13"},
				{Name: "arr14"},
			},
		},
		{
			name: "struct-with-struct",
			ptr: &struct {
				F1 int64
				F2 struct {
					FF1 int64
					FF2 float64
					FF3 struct {
						FFF1 float64
					}
				}
			}{},
			want: []ReadVar{
				{Name: "F1"},
				{Name: "F2"},
			},
		},
		{
			name: "struct-with-struct+slice",
			ptr: &struct {
				F1 int64
				F2 struct {
					FF1 int64
					FF2 float64
					FF3 []float64
					FF4 []struct {
						FFF1 float64
						FFF2 []float64
					}
				}
			}{},
			want: []ReadVar{
				{Name: "F1"},
				{Name: "F2"},
			},
		},
		{
			name: "invalid-slice-tag",
			ptr: &struct {
				N   int32
				Sli []int32 `groot:"vs[N][N]"`
			}{},
			panics: "rtree: invalid number of slice-dimensions for field \"Sli\": \"vs[N][N]\"",
		},
		{
			name: "invalid-array-tag",
			ptr: &struct {
				N   int32
				Arr [12]int32 `groot:"vs[1][2][3][4]"`
			}{},
			panics: "rtree: invalid number of array-dimension for field \"Arr\": \"vs[1][2][3][4]\"",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message:\ngot= %q\nwant=%q", got, want)
					}
				}()
			}
			got := ReadVarsFromStruct(tc.ptr)
			if got, want := len(got), len(tc.want); got != want {
				t.Fatalf("invalid number of rvars: got=%d, want=%d", got, want)
			}
			for i := range got {
				if got, want := got[i].Name, tc.want[i].Name; got != want {
					t.Fatalf("invalid name for rvar[%d]: got=%q, want=%q", i, got, want)
				}
				if got, want := got[i].count, tc.want[i].count; got != want {
					t.Fatalf("invalid count for rvar[%d]: got=%q, want=%q", i, got, want)
				}
			}
		})
	}
}

func TestNameOf(t *testing.T) {
	for _, tc := range []struct {
		ptr  interface{}
		want string
	}{
		{
			ptr: &struct {
				F int64 `groot:"field"`
			}{},
			want: "field",
		},
		{
			ptr: &struct {
				F [1]int64 `groot:"field[1]"`
			}{},
			want: "field[1]",
		},
		{
			ptr: &struct {
				F [1]int64 `groot:"field"`
			}{},
			want: "field[1]",
		},
		{
			ptr: &struct {
				F [1][2][3]int64 `groot:"field"`
			}{},
			want: "field[1][2][3]",
		},
		{
			ptr: &struct {
				F [1][2][3]int64 `groot:"field[6]"`
			}{},
			want: "field[6]",
		},
		{
			ptr: &struct {
				F []int64 `groot:"field"`
			}{},
			want: "field",
		},
		{
			ptr: &struct {
				N int32   `groot:"N"`
				F []int64 `groot:"field[N]"`
			}{},
			want: "field[N]",
		},
	} {
		t.Run("", func(t *testing.T) {
			rt := reflect.TypeOf(tc.ptr).Elem()
			field, ok := rt.FieldByName("F")
			if !ok {
				t.Fatalf("could not retrieve field named \"F\" for %T", tc.ptr)
			}
			got := nameOf(field)
			if got, want := got, tc.want; got != want {
				t.Fatalf("invalid name: got=%q, want=%q", got, want)
			}
		})
	}
}
