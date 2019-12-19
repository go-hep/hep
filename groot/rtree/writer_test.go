// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
)

func TestWriteVarsFromStruct(t *testing.T) {
	for _, tc := range []struct {
		name   string
		ptr    interface{}
		want   []WriteVar
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
			want: []WriteVar{{Name: "I32"}, {Name: "F32"}, {Name: "Str"}},
		},
		{
			name: "simple-with-unexported",
			ptr: &struct {
				I32 int32
				F32 float32
				val float32
				Str string
			}{},
			want: []WriteVar{{Name: "I32"}, {Name: "F32"}, {Name: "Str"}},
		},
		{
			name: "slices",
			ptr: &struct {
				N      int32
				NN     int64
				SliF32 []float32 `groot:"F32s[N]"`
				SliF64 []float64 `groot:"F64s[NN]"`
			}{},
			want: []WriteVar{
				{Name: "N"},
				{Name: "NN"},
				{Name: "F32s", Count: "N"},
				{Name: "F64s", Count: "NN"},
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
			want: []WriteVar{
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
			got := WriteVarsFromStruct(tc.ptr)
			if got, want := len(got), len(tc.want); got != want {
				t.Fatalf("invalid number of wvars: got=%d, want=%d", got, want)
			}
			for i := range got {
				if got, want := got[i].Name, tc.want[i].Name; got != want {
					t.Fatalf("invalid name for wvar[%d]: got=%q, want=%q", i, got, want)
				}
				if got, want := got[i].Count, tc.want[i].Count; got != want {
					t.Fatalf("invalid count for wvar[%d]: got=%q, want=%q", i, got, want)
				}
			}
		})
	}
}

func TestFlattenArrayType(t *testing.T) {
	for _, tc := range []struct {
		typ   interface{}
		want  interface{}
		shape []int
	}{
		{
			typ:   int32(0),
			want:  int32(0),
			shape: nil,
		},
		{
			typ:   [2]int32{},
			want:  int32(0),
			shape: []int{2},
		},
		{
			typ:   [2][3]int32{},
			want:  int32(0),
			shape: []int{2, 3},
		},
		{
			typ:   [2][3][4]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4},
		},
		{
			typ:   [2][3][4][0]int32{},
			want:  int32(0),
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   [2][3][4][0]struct{}{},
			want:  struct{}{},
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   [2][3][4][0][]string{},
			want:  []string{},
			shape: []int{2, 3, 4, 0},
		},
		{
			typ:   []string{},
			want:  []string{},
			shape: nil,
		},
	} {
		t.Run(fmt.Sprintf("%T", tc.typ), func(t *testing.T) {
			rt, shape := flattenArrayType(reflect.TypeOf(tc.typ))
			if got, want := rt, reflect.TypeOf(tc.want); got != want {
				t.Fatalf("invalid array element type: got=%v, want=%v", got, want)
			}

			if got, want := shape, tc.shape; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid shape: got=%v, want=%v", got, want)
			}
		})
	}
}

func TestInvalidTreeMerger(t *testing.T) {
	var (
		w   wtree
		src = rbase.NewObjString("foo")
	)

	err := w.ROOTMerge(src)
	if err == nil {
		t.Fatalf("expected an error")
	}

	const want = "rtree: can not merge src=*rbase.ObjString into dst=*rtree.wtree"
	if got, want := err.Error(), want; got != want {
		t.Fatalf("invalid ROOTMerge error. got=%q, want=%q", got, want)
	}
}
