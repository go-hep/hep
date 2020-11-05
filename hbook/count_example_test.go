// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"

	"go-hep.org/x/hep/hbook"
)

func ExampleCount() {
	c := hbook.Count{}
	c.XRange = hbook.Range{Min: 0, Max: 1}
	c.Val = 10
	c.Err.Low = 2
	c.Err.High = 3

	fmt.Printf("[%v, %v] -> %v +%v -%v",
		c.XRange.Min, c.XRange.Max,
		c.Val, c.Err.High, c.Err.Low)

	// Output:
	// [0, 1] -> 10 +3 -2
}

func ExampleCount_withH1D() {

	h := hbook.NewH1D(6, 0, 6)
	h.Fill(-0.5, 1)
	h.Fill(0, 1.5)
	h.Fill(0.5, 1)
	h.Fill(1.2, 1)
	h.Fill(2.1, 2)
	h.Fill(4.2, 1)
	h.Fill(5.9, 1)
	h.Fill(6, 0.5)

	for _, c := range h.Counts() {
		fmt.Printf("[%v, %v] -> %v +%.1f -%.1f\n",
			c.XRange.Min, c.XRange.Max,
			c.Val, c.Err.High, c.Err.Low)
	}

	// Output:
	// [0, 1] -> 2.5 +0.9 -0.9
	// [1, 2] -> 1 +0.5 -0.5
	// [2, 3] -> 2 +1.0 -1.0
	// [3, 4] -> 0 +0.0 -0.0
	// [4, 5] -> 1 +0.5 -0.5
	// [5, 6] -> 1 +0.5 -0.5
}
