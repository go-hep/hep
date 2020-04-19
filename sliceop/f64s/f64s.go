// Copyright 2019 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Common operations on float64 slices
package f64s


// Return a slice with all elements for which f(x) returns true
func Filter(dst, src []float64, f func(v float64) bool) []float64 {
	if dst == nil {
		dst = []float64{}
	}
	for _, x := range src {
		if f(x) {
			dst = append(dst, x)
		}
	}
	return dst	
}

// Return a slice of {f(x_i)} when src = {x_i}
func Map(dst, src []float64, f func(v float64) float64) []float64 {
	if dst == nil {
		dst = []float64{}
	}
	for _, x := range src {
		dst = append(dst, f(x))
	}
	return dst
}

// Return a slice of all indices corresponding to elements for which f(x) is true
func Find(dst []int, src []float64, f func(v float64) bool) []int {
	if dst == nil {
		dst = []int{}
	}
	for i, x := range src {
		if f(x) {
			dst = append(dst, i)
		}
	}
	return dst
}


