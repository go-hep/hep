// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import "testing"

func TestWriteVarsFromStruct(t *testing.T) {
	for _, tc := range []struct {
		name   string
		ptr    interface{}
		wopts  []WriteOption
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
			name: "slices-no-count",
			ptr: &struct {
				F1 int32
				X2 []float32 `groot:"F2[F1]"`
				X3 []float64 `groot:"F3"`
				F4 []float64
			}{},
			want: []WriteVar{
				{Name: "F1"},
				{Name: "F2", Count: "F1"},
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
			want: []WriteVar{
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
			want: []WriteVar{
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
		{
			name: "no-split-struct",
			ptr: &struct {
				N   int32
				F32 float32
			}{},
			wopts: []WriteOption{WithTitle("evt"), WithSplitLevel(0)},
			want: []WriteVar{
				{Name: "evt"},
			},
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
			got := WriteVarsFromStruct(tc.ptr, tc.wopts...)
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
