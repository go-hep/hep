// Copyright ©2025 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/fit"
)

func ExamplePoly() {
	var (
		xs      = []float64{0.0, 1.0, 2.0, 3.0, +4.0, +5.0}
		ys      = []float64{0.0, 0.8, 0.9, 0.1, -0.8, -1.0}
		degree  = 3
		zs, err = fit.Poly(xs, ys, degree)
	)

	if err != nil {
		log.Fatalf("could not fit polynomial: %v", err)
	}

	fmt.Printf("z = %+.5f\n", zs)
	// Output:
	// z = [+0.08704 -0.81349 +1.69312 -0.03968]
}
