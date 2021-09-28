// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntcsv_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/hbook/ntup/ntcsv"
)

func ExampleOpen() {
	// Open a new n-tuple pointing at a CSV file "testdata/simple.csv"
	// whose field separator is ';'.
	// We rename the columns v1, v2 and v3.
	nt, err := ntcsv.Open(
		"testdata/simple.csv",
		ntcsv.Comma(';'),
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

	err = nt.Scan("v1, v2, v3", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}

func ExampleOpen_fromRemote() {
	// Open a new n-tuple pointing at a remote CSV file
	// "https://github.com/go-hep/hep/raw/main/hbook/ntup/ntcsv/testdata/simple.csv"
	// whose field separator is ';'.
	// We rename the columns v1, v2 and v3.
	nt, err := ntcsv.Open(
		"https://github.com/go-hep/hep/raw/main/hbook/ntup/ntcsv/testdata/simple.csv",
		ntcsv.Comma(';'),
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

	err = nt.Scan("v1, v2, v3", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}

func ExampleOpen_withDefaultVarNames() {
	// Open a new n-tuple pointing at a CSV file "testdata/simple.csv"
	// whose field separator is ';'.
	// We use the default column names: var1, var2, var3, ...
	nt, err := ntcsv.Open(
		"testdata/simple.csv",
		ntcsv.Comma(';'),
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

	err = nt.Scan("var1, var2, var3", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}

func ExampleOpen_withHeader() {
	// Open a new n-tuple pointing at a CSV file "testdata/simple.csv"
	// whose field separator is ';'.
	// We rename the columns v1, v2 and v3.
	// We tell the CSV driver to handle the CSV header.
	nt, err := ntcsv.Open(
		"testdata/simple-with-header.csv",
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

	err = nt.Scan("v1, v2, v3", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}

func ExampleOpen_withHeaderAndImplicitColumns() {
	// Open a new n-tuple pointing at a CSV file "testdata/simple.csv"
	// whose field separator is ';'.
	// We tell the CSV driver to handle the CSV header.
	// And we implicitly use the column names for the queries.
	nt, err := ntcsv.Open(
		"testdata/simple-with-header.csv",
		ntcsv.Comma(';'),
		ntcsv.Header(),
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

	err = nt.Scan("i, f, str", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}

func ExampleOpen_withHeaderAndExlicitColumns() {
	// Open a new n-tuple pointing at a CSV file "testdata/simple.csv"
	// whose field separator is ';'.
	// We tell the CSV driver to handle the CSV header.
	// And we explicitly use our column names for the queries.
	nt, err := ntcsv.Open(
		"testdata/simple-with-header.csv",
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

	err = nt.Scan("v1, v2, v3", func(i int64, f float64, s string) error {
		fmt.Printf("%d %f %q\n", i, f, s)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0 0.000000 "str-0"
	// 1 1.000000 "str-1"
	// 2 2.000000 "str-2"
	// 3 3.000000 "str-3"
	// 4 4.000000 "str-4"
	// 5 5.000000 "str-5"
	// 6 6.000000 "str-6"
	// 7 7.000000 "str-7"
	// 8 8.000000 "str-8"
	// 9 9.000000 "str-9"
}
