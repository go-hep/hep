// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/hbook"
)

func main() {
	f, err := os.Open("gauss-1d-data.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := hbook.NewH1D(10, -4, 4)
	h.Annotation()["name"] = "h1"
	h.Annotation()["title"] = "h1"
	const n = 10004
	for i := 0; i < n; i++ {
		var x, w float64
		_, err = fmt.Fscanf(f, "%g %g\n", &x, &w)
		if err != nil {
			log.Fatal(err)
		}
		h.Fill(x, w)
	}

	raw, err := h.MarshalYODA()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("=== YODA ===\n%v\n", string(raw))
}
