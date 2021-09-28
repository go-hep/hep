// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntup_test

import (
	"database/sql"
	"fmt"
	"log"
	"math"

	"go-hep.org/x/hep/hbook/ntup"
	"go-hep.org/x/hep/hbook/ntup/ntcsv"
)

func ExampleNtuple_open() {
	db, err := sql.Open("csv", "ntcsv/testdata/simple-with-header.csv")
	if err != nil {
		log.Fatalf("could not open csv-db file: %+v", err)
	}
	defer db.Close()

	nt, err := ntup.Open(db, "ntup")
	if err != nil {
		log.Fatalf("could not open ntup: %+v", err)
	}
	fmt.Printf("name=%q\n", nt.Name())

	// Output:
	// name="ntup"
}

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
	defer func() {
		err = nt.DB().Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

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

func ExampleNtuple_scan() {
	nt, err := ntcsv.Open(
		"ntcsv/testdata/simple-with-header.csv",
		ntcsv.Comma(';'),
		ntcsv.Header(),
		ntcsv.Columns("v1", "v2", "v3"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = nt.DB().Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var (
		v1min = +math.MaxFloat64
		v1max = -math.MaxFloat64
		v2min = +math.MaxFloat64
		v2max = -math.MaxFloat64
	)
	err = nt.Scan("v1, v2", func(v1, v2 float64) error {
		v1min = math.Min(v1min, v1)
		v1max = math.Max(v1max, v1)
		v2min = math.Min(v2min, v2)
		v2max = math.Max(v2max, v2)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("V1Min  %v\n", v1min)
	fmt.Printf("V1Max  %v\n", v1max)
	fmt.Printf("V2Min  %v\n", v2min)
	fmt.Printf("V2Max  %v\n", v2max)

	//Output:
	// V1Min  0
	// V1Max  9
	// V2Min  0
	// V2Max  9

}

func ExampleNtuple_scanH1D() {
	nt, err := ntcsv.Open(
		"ntcsv/testdata/simple-with-header.csv",
		ntcsv.Comma(';'),
		ntcsv.Header(),
		ntcsv.Columns("v1", "v2", "v3"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = nt.DB().Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	h, err := nt.ScanH1D("v1", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("V1Mean:      %f\n", h.XMean())
	fmt.Printf("V1RMS:       %f\n", h.XRMS())
	fmt.Printf("V1StdDev:    %f\n", h.XStdDev())
	fmt.Printf("V1StdErr:    %f\n", h.XStdErr())

	// Output:
	// V1Mean:      4.500000
	// V1RMS:       5.338539
	// V1StdDev:    3.027650
	// V1StdErr:    0.957427
}
