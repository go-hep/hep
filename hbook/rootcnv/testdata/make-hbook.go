// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/hbook"
)

func main() {
	{
		f, err := os.Open("gauss-1d-data.dat")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		h1 := hbook.NewH1D(10, -4, 4)
		h1.Annotation()["name"] = "h1"
		h1.Annotation()["title"] = "h1"

		const n = 10004
		for i := 0; i < n; i++ {
			var x, w float64
			_, err = fmt.Fscanf(f, "%g %g\n", &x, &w)
			if err != nil {
				log.Fatal(err)
			}
			h1.Fill(x, w)
		}

		raw1, err := h1.MarshalYODA()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("=== YODA ===\n%v\n", string(raw1))
	}

	{
		f, err := os.Open("gauss-2d-data.dat")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		h2 := hbook.NewH2D(3, 0, 3, 3, 0, 3)
		h2.Annotation()["name"] = "h2d"
		h2.Annotation()["title"] = "h2d"

		for i := 0; i < 10000; i++ {
			var x, y, w float64
			_, err = fmt.Fscanf(f, "%g %g %g\n", &x, &y, &w)
			if err != nil {
				log.Fatal(err)
			}
			h2.Fill(x, y, w)
		}
		h2.Fill(+5, +5, 101) // NE
		h2.Fill(+0, +5, 102) // N
		h2.Fill(-5, +5, 103) // NW
		h2.Fill(-5, +0, 104) // W
		h2.Fill(-5, -5, 105) // SW
		h2.Fill(+0, -5, 106) // S
		h2.Fill(+5, -5, 107) // SE
		h2.Fill(+5, +0, 108) // E

		raw2, err := h2.MarshalYODA()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("=== YODA ===\n%v\n", string(raw2))
	}
}
