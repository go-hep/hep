// Copyright ©2021 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.18

// Package sliceop provides operations on slices not available in the stdlib
// slices package.
package sliceop // import "go-hep.org/x/hep/sliceop"

import "errors"

var (
	errLength           = errors.New("sliceop: length mismatch")
	errSortedIndices    = errors.New("sliceop: indices not sorted")
	errDuplicateIndices = errors.New("sliceop: duplicate indices")
)

// Filter creates a slice with all the elements x_i of src for which f(x_i) is true.
// Filter uses dst as work buffer, storing elements at the start of the slice.
// Filter clears dst if a slice is passed, and allocates a new slice if dst is nil.
func Filter[T any](dst, src []T, f func(v T) bool) []T {

	if dst == nil {
		dst = make([]T, 0, len(src))
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
func Map[T, U any](dst []U, src []T, f func(v T) U) []U {

	if dst == nil {
		dst = make([]U, len(src))
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
func Find[S ~[]E, E any](dst []int, src S, f func(v E) bool) []int {

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
func Take[S ~[]E, E any](dst, src S, indices []int) S {

	if len(indices) > len(src) {
		panic(errLength)
	}

	if dst == nil {
		dst = make(S, len(indices))
	}

	if len(dst) != len(indices) {
		panic(errLength)
	}

	if len(indices) == 0 {
		return dst
	}

	dst[0] = src[indices[0]]
	var (
		v0 = indices[0]
		nn = len(indices)
	)
	for i := 1; i < nn; i++ {
		v1 := indices[i]
		switch {
		case v0 < v1:
			// ok.
		case v0 == v1:
			panic(errDuplicateIndices)
		case v0 > v1:
			panic(errSortedIndices)
		}
		dst[i-1] = src[v0]
		v0 = v1
	}
	dst[nn-1] = src[v0]

	return dst
}

// Resize returns a slice of size n, reusing the storage of the provided
// slice, appending new elements if the capacity is not sufficient.
func Resize[S ~[]E, E any](sli S, n int) S {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make(S, n-m)...)
	}
	sli = sli[:n]
	return sli
}
