// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntcsv_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook/ntup/ntcsv"
)

func TestOpen(t *testing.T) {
	type dataType struct {
		i int64
		f float64
		s string
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
	for _, test := range []struct {
		name    string
		comma   rune
		comment rune
		query   string
	}{
		{"testdata/simple.csv", ';', '#', `var1, var2, var3`},
		{"testdata/simple-comma.csv", ',', '#', `var1, var2, var3`},
		// FIXME(sbinet): the name of the variables should be taken from the header
		// {"testdata/simple-with-header.csv", ';', '#', `i, f, str`},
		{"testdata/simple-with-header.csv", ';', '#', `var1, var2, var3`},
	} {
		nt, err := ntcsv.Open(test.name, ntcsv.Comma(test.comma))
		if err != nil {
			t.Errorf("%s: error opening n-tuple: %v", test.name, err)
			continue
		}
		defer nt.DB().Close()

		var got []dataType
		err = nt.Scan(
			test.query,
			func(i int64, f float64, s string) error {
				got = append(got, dataType{i, f, s})
				return nil
			},
		)
		if err != nil {
			t.Errorf("%s: error scanning: %v", test.name, err)
			continue
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got=\n%v\nwant=\n%v\n", test.name, got, want)
			continue
		}
	}
}
