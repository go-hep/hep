// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.12
// +build go1.12

package hep_test

import (
	"fmt"

	"go-hep.org/x/hep"
)

func Example_version() {
	version, sum := hep.Version()
	fmt.Printf("Go-HEP version:  %q\n", version)
	fmt.Printf("       checksum: %q\n", sum)
}
