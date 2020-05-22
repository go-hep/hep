// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntcsv_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook/ntup/ntcsv"
)

func TestOpen(t *testing.T) {
	for _, test := range []struct {
		name  string
		query string
		opts  []ntcsv.Option
	}{
		{"testdata/simple.csv", `var1, var2, var3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
			},
		},
		{"testdata/simple.csv", `var1, var2, var3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
				ntcsv.Columns("var1", "var2", "var3"),
			},
		},
		{"testdata/simple.csv", `v1, v2, v3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
				ntcsv.Columns("v1", "v2", "v3"),
			},
		},
		{"testdata/simple-comma.csv", `var1, var2, var3`,
			[]ntcsv.Option{
				ntcsv.Comma(','),
				ntcsv.Comment('#'),
			},
		},
		{"testdata/simple-with-comment.csv", `var1, var2, var3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
			},
		},
		{"testdata/simple-with-comment.csv", `v1, v2, v3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
				ntcsv.Columns("v1", "v2", "v3"),
			},
		},
		{"testdata/simple-with-header.csv", `i, f, str`,
			[]ntcsv.Option{
				ntcsv.Header(),
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
			},
		},
		{"testdata/simple-with-header.csv", `i64, f64, str`,
			[]ntcsv.Option{
				ntcsv.Header(),
				ntcsv.Columns("i64", "f64", "str"),
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
			},
		},
		{"http://github.com/go-hep/hep/raw/master/hbook/ntup/ntcsv/testdata/simple-with-comment.csv", `v1, v2, v3`,
			[]ntcsv.Option{
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
				ntcsv.Columns("v1", "v2", "v3"),
			},
		},
		{"https://github.com/go-hep/hep/raw/master/hbook/ntup/ntcsv/testdata/simple-with-header.csv", `i64, f64, str`,
			[]ntcsv.Option{
				ntcsv.Header(),
				ntcsv.Columns("i64", "f64", "str"),
				ntcsv.Comma(';'),
				ntcsv.Comment('#'),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testCSV(t, test.name, test.query, test.opts...)
		})
	}
}

func testCSV(t *testing.T, name, query string, opts ...ntcsv.Option) {
	nt, err := ntcsv.Open(name, opts...)
	if err != nil {
		t.Fatalf("%s: error opening n-tuple: %v", name, err)
	}
	defer nt.DB().Close()

	type dataType struct {
		i int64
		f float64
		s string
	}

	var got []dataType
	err = nt.Scan(
		query,
		func(i int64, f float64, s string) error {
			got = append(got, dataType{i, f, s})
			return nil
		},
	)
	if err != nil {
		t.Fatalf("%s: error scanning: %v", name, err)
	}

	want := []dataType{
		{0, 0, "str-0"},
		{1, 1, "str-1"},
		{2, 2, "str-2"},
		{3, 3, "str-3"},
		{4, 4, "str-4"},
		{5, 5, "str-5"},
		{6, 6, "str-6"},
		{7, 7, "str-7"},
		{8, 8, "str-8"},
		{9, 9, "str-9"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: got=\n%v\nwant=\n%v\n", name, got, want)
	}
}
