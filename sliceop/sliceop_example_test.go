// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sliceop_test

import (
	"fmt"

	"go-hep.org/x/hep/sliceop"
)

// An example of slice filtering
func ExampleFilter() {
	slice := []float64{1, 2, -99, 4, 5, -99, 7}
	condition := func(x float64) bool { return x > 0 }
	fmt.Println(sliceop.Filter(nil, slice, condition))

	// Output:
	// [1 2 4 5 7]
}

// An example of slice mapping
func ExampleMap() {
	slice := []float64{1, 2, -99, 4, 5, -99, 7}
	operation := func(x float64) float64 { return x * x }
	fmt.Println(sliceop.Map(nil, slice, operation))

	// Output:
	// [1 4 9801 16 25 9801 49]
}

// An example of slice finding
func ExampleFind() {
	slice := []float64{1, 2, -99, 4, 5, -99, 7}
	condition := func(x float64) bool { return x == -99 }
	fmt.Println(sliceop.Find(nil, slice, condition))

	// Output:
	// [2 5]
}

// An example of taking a sub-slice defined by indices
func ExampleTake() {
	slice := []float64{1, 2, -99, 4, 5, -99, 7}
	indices := []int{2, 5}
	fmt.Println(sliceop.Take(nil, slice, indices))

	// Output:
	// [-99 -99]
}

// An example of resizing a slice.
func ExampleResize() {
	slice := []int{1, 2, 3, 4}
	fmt.Printf("slice: len=%d, cap=%d, vs=%v\n", len(slice), cap(slice), slice)
	slice = sliceop.Resize(slice, 8)
	fmt.Printf("slice: len=%d, cap=%d, vs=%v\n", len(slice), cap(slice), slice)
	slice = sliceop.Resize(slice, 3)
	fmt.Printf("slice: len=%d, cap=%d, vs=%v\n", len(slice), cap(slice), slice)

	// Output:
	// slice: len=4, cap=4, vs=[1 2 3 4]
	// slice: len=8, cap=8, vs=[1 2 3 4 0 0 0 0]
	// slice: len=3, cap=8, vs=[1 2 3]
}
