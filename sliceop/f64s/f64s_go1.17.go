// Copyright ©2020 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.18

// Package f64s provides common operations on float64 slices.
package f64s

import (
	"fmt"
)

var (
	errLength           = fmt.Errorf("sliceop: length mismatch")
	errSortedIndices    = fmt.Errorf("sliceop: indices not sorted")
	errDuplicateIndices = fmt.Errorf("sliceop: duplicate indices")
)

// Filter creates a slice with all the elements x_i of src for which f(x_i) is true.
// Filter uses dst as work buffer, storing elements at the start of the slice.
// Filter clears dst if a slice is passed, and allocates a new slice if dst is nil.
func Filter(dst, src []float64, f func(v float64) bool) []float64 {

	if dst == nil {
		dst = make([]float64, 0, len(src))
	}

	dst = dst[:0]
	for _, x := range src {
		if f(x) {
			dst = append(dst, x)
		}
	}

	return dst
}

// Map creates a slice with all the elements f(x_i) where x_i are elements from src.
// Map uses dst as work buffer, storing elements at the start of the slice.
// Map allocates a new slice if dst is nil.
// Map will panic if the lengths of src and dst differ.
func Map(dst, src []float64, f func(v float64) float64) []float64 {

	if dst == nil {
		dst = make([]float64, len(src))
	}

	if len(src) != len(dst) {
		panic(errLength)
	}

	for i, x := range src {
		dst[i] = f(x)
	}
	return dst
}

// Find creates a slice with all indices corresponding to elements for which f(x) is true.
// Find uses dst as work buffer, storing indices at the start of the slice.
// Find clears dst if a slice is passed, and allocates a new slice if dst is nil.
func Find(dst []int, src []float64, f func(v float64) bool) []int {

	if dst == nil {
		dst = make([]int, 0, len(src))
	}

	dst = dst[:0]
	for i, x := range src {
		if f(x) {
			dst = append(dst, i)
		}
	}

	return dst
}

// Take creates a sub-slice of src with all elements indiced by the provided indices.
// Take uses dst as work buffer, storing elements at the start of the slice.
// Take clears dst if a slice is passed, and allocates a new slice if dst is nil.
// Take will panic if indices is not sorted or has duplicates.
// Take will panic if length of indices is larger than length of src.
// Take will panic if length of indices is different from length of dst.
func Take(dst, src []float64, indices []int) []float64 {

	if len(indices) > len(src) {
		panic(errLength)
	}

	if dst == nil {
		dst = make([]float64, len(indices))
	}

	if len(dst) != len(indices) {
		panic(errLength)
	}

	if len(indices) == 0 {
		return dst
	}

	dst[0] = src[indices[0]]
	for i := 1; i < len(indices); i++ {
		v0 := indices[i-1]
		v1 := indices[i]
		switch {
		case v0 < v1:
			// ok.
		case v0 == v1:
			panic(errDuplicateIndices)
		case v0 > v1:
			panic(errSortedIndices)
		}
		dst[i] = src[v1]
	}

	return dst
}
