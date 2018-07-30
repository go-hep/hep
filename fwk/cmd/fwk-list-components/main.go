// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"go-hep.org/x/hep/fwk"
)

func main() {
	comps := fwk.Registry()
	fmt.Printf("::: components... (%d)\n", len(comps))
	for i, c := range comps {
		fmt.Printf("[%04d/%04d] %s\n", i, len(comps), c)
	}
}
