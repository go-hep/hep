// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntup_test

import (
	"fmt"
	"log"
	"math"

	"go-hep.org/x/hep/hbook/ntup/ntcsv"
)

func ExampleNtuple_scanH2D() {
	nt, err := ntcsv.Open(
		"ntcsv/testdata/simple-with-header.csv",
		ntcsv.Comma(';'),
		ntcsv.Header(),
		ntcsv.Columns("v1", "v2", "v3"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nt.DB().Close()

	h, err := nt.ScanH2D("v1, v2", nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("XMean:      %f\n", h.XMean())
	fmt.Printf("YMean:      %f\n", h.YMean())
	fmt.Printf("XRMS:       %f\n", h.XRMS())
	fmt.Printf("YRMS:       %f\n", h.YRMS())
	fmt.Printf("XStdDev:    %f\n", h.XStdDev())
	fmt.Printf("YStdDev:    %f\n", h.YStdDev())
	fmt.Printf("XStdErr:    %f\n", h.XStdErr())
	fmt.Printf("YStdErr:    %f\n", h.YStdErr())

	// Output:
	// XMean:      4.500000
	// YMean:      4.500000
	// XRMS:       5.338539
	// YRMS:       5.338539
	// XStdDev:    3.027650
	// YStdDev:    3.027650
	// XStdErr:    0.957427
	// YStdErr:    0.957427
}

func ExampleNtuple_scanH() {
	nt, err := ntcsv.Open(
		"ntcsv/testdata/simple-with-header.csv",
		ntcsv.Comma(';'),
		ntcsv.Header(),
		ntcsv.Columns("v1", "v2", "v3"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer nt.DB().Close()
	var (
		xmin = +math.MaxFloat64
		xmax = -math.MaxFloat64
		ymin = +math.MaxFloat64
		ymax = -math.MaxFloat64
	)
	query := "v1, v2"
	error_ := nt.Scan(query, func(x, y float64) error {
		xmin = math.Min(xmin, x)
		xmax = math.Max(xmax, x)
		ymin = math.Min(ymin, y)
		ymax = math.Max(ymax, y)
		return nil
	})
	if error_ != nil {
		log.Fatal(error_)
	}

	fmt.Printf("Result  %v", error_)

	//Output:
	//Result  <nil>
}
