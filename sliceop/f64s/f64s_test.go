// Copyright ©2020 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64s

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func panics(t *testing.T, want error) func() {
	t.Helper()
	return func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected a panic")
		}
		if got, want := err.(error).Error(), want.Error(); got != want {
			t.Fatalf("invalid panic message.\ngot= %v\nwant=%v",
				got, want,
			)
		}
	}
}

func TestMap(t *testing.T) {
	defer panics(t, errLength)()

	_ = Map(make([]float64, 3), make([]float64, 5), nil)
}

func TestTake(t *testing.T) {
	for _, tc := range []struct {
		dst, src []float64
		inds     []int
		want     []float64
		panics   error
	}{
		{
			dst:  nil,
			src:  []float64{1, 2, 3},
			inds: []int{1},
			want: []float64{2},
		},
		{
			dst:  make([]float64, 1),
			src:  []float64{1, 2, 3},
			inds: []int{1},
			want: []float64{2},
		},
		{
			dst:  make([]float64, 0),
			src:  []float64{1, 2, 3},
			inds: []int{},
			want: []float64{},
		},
		{
			dst:  []float64{},
			src:  []float64{1, 2, 3},
			inds: nil,
			want: []float64{},
		},
		{
			dst:  make([]float64, 2),
			src:  []float64{1, 2, 3},
			inds: []int{1, 2},
			want: []float64{2, 3},
		},
		{
			dst:    nil,
			src:    []float64{1, 2, 3},
			inds:   []int{1, 0},
			panics: errSortedIndices,
		},
		{
			dst:    nil,
			src:    []float64{1, 2, 3},
			inds:   []int{0, 1, 2, 3},
			panics: errLength,
		},
		{
			dst:    make([]float64, 1),
			src:    []float64{1, 2, 3},
			inds:   []int{0, 1},
			panics: errLength,
		},
		{
			dst:    make([]float64, 4),
			src:    []float64{1, 2, 3},
			inds:   []int{0, 1},
			panics: errLength,
		},
		{
			dst:    nil,
			src:    []float64{1, 2, 3},
			inds:   []int{1, 1},
			panics: errDuplicateIndices,
		},
	} {
		t.Run("", func(t *testing.T) {
			if tc.panics != nil {
				defer panics(t, tc.panics)()
			}
			got := Take(tc.dst, tc.src, tc.inds)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got= %v\nwant=%v", got, tc.want)
			}
		})
	}
}

var takeSink []float64

func BenchmarkTake(b *testing.B) {
	for _, size := range []int{2, 4, 8, 128, 1024, 1024 * 1024} {
		b.Run(fmt.Sprintf("Len=%d", size), func(b *testing.B) {
			src := make([]float64, size)
			ind := make([]int, 0, len(src))
			rnd := rand.New(rand.NewSource(0))
			for i := range src {
				src[i] = rnd.Float64()
				if rnd.Float64() > 0.5 {
					ind = append(ind, i)
				}
			}
			dst := make([]float64, len(ind))
			b.ReportAllocs()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				takeSink = Take(dst, src, ind)
			}
		})
	}
}
