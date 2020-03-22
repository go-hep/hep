// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntroot_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/hbook/ntup/ntroot"
)

func ExampleOpen() {
	nt, err := ntroot.Open("../../../groot/testdata/simple.root", "tree")
	if err != nil {
		log.Fatalf("could not open n-tuple: %+v", err)
	}
	defer nt.DB().Close()

	err = nt.Scan(
		"(one, two, three)",
		func(i int32, f float32, s string) error {
			fmt.Printf("row=(%v, %v, %q)\n", i, f, s)
			return nil
		},
	)

	if err != nil {
		log.Fatalf("could not scan n-tuple: %+v", err)
	}

	// Output:
	// row=(1, 1.1, "uno")
	// row=(2, 2.2, "dos")
	// row=(3, 3.3, "tres")
	// row=(4, 4.4, "quatro")
}
