// Copyright ©2020 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package f64s provides common operations on float64 slices.
package f64s

import (
	"sort"
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
		panic("f64s: length mismatch")
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

	if !sort.IntsAreSorted(indices) {
		panic("f64s: indices not sorted")
	}

	if len(indices) > len(src) {
		panic("f64s: length mismatch")
	}

	if dst == nil {
		dst = make([]float64, len(indices))
	}

	if len(dst) != len(indices) {
		panic("f64s: length mismatch")
	}

	dst = dst[0:len(indices)]

	m := make(map[int]int, len(indices))
	for i, v := range indices {
		if _, there := m[v]; !there {
			m[v] = v
		} else {
			panic("f64s: duplicate indices")
		}
		dst[i] = src[v]
	}

	return dst
}
