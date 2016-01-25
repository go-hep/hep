// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/go-hep/csvutil"
)

func TestCSVReaderScanArgs(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, 10)
	if err != nil {
		t.Errorf("could read rows [0, 10): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		var (
			i int
			f float64
			s string
		)
		err = rows.Scan(&i, &f, &s)
		if err != nil {
			t.Errorf("error reading row %d: %v\n", irow, err)
		}
		exp := fmt.Sprintf("%d;%d;str-%d", irow, irow, irow)
		got := fmt.Sprintf("%v;%v;%v", i, f, s)
		if exp != got {
			t.Errorf("error reading row %d\nexp=%q\ngot=%q\n",
				irow, exp, got,
			)
		}
		irow++
	}

	err = rows.Err()
	if err != nil {
		t.Errorf("error iterating over rows: %v\n", err)
	}
}

func TestCSVReaderScanStruct(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, 10)
	if err != nil {
		t.Errorf("could read rows [0, 10): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		data := struct {
			I int
			F float64
			S string
		}{}
		err = rows.Scan(&data)
		if err != nil {
			t.Errorf("error reading row %d: %v\n", irow, err)
		}
		exp := fmt.Sprintf("%d;%d;str-%d", irow, irow, irow)
		got := fmt.Sprintf("%v;%v;%v", data.I, data.F, data.S)
		if exp != got {
			t.Errorf("error reading row %d\nexp=%q\ngot=%q\n",
				irow, exp, got,
			)
		}
		irow++
	}

	err = rows.Err()
	if err != nil {
		t.Errorf("error iterating over rows: %v\n", err)
	}
}

func TestCSVReaderScanSmallRead(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, 2)
	if err != nil {
		t.Errorf("could read rows [0, 2): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		data := struct {
			I int
			F float64
			S string
		}{}
		err = rows.Scan(&data)
		if err != nil {
			t.Errorf("error reading row %d: %v\n", irow, err)
		}
		exp := fmt.Sprintf("%d;%d;str-%d", irow, irow, irow)
		got := fmt.Sprintf("%v;%v;%v", data.I, data.F, data.S)
		if exp != got {
			t.Errorf("error reading row %d\nexp=%q\ngot=%q\n",
				irow, exp, got,
			)
		}
		irow++
	}

	err = rows.Err()
	if err != nil {
		t.Errorf("error iterating over rows: %v\n", err)
	}
}

func TestCSVReaderScanEOF(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, 12)
	if err != nil {
		t.Errorf("could read rows [0, 12): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		data := struct {
			I int
			F float64
			S string
		}{}
		err = rows.Scan(&data)
		if err != nil {
			if irow == 10 {
				break
			}
			t.Errorf("error reading row %d: %v\n", irow, err)
		}
		exp := fmt.Sprintf("%d;%d;str-%d", irow, irow, irow)
		got := fmt.Sprintf("%v;%v;%v", data.I, data.F, data.S)
		if exp != got {
			t.Errorf("error reading row %d\nexp=%q\ngot=%q\n",
				irow, exp, got,
			)
		}
		irow++
	}

	if irow != 10 {
		t.Errorf("error. expected irow==10. got=%v\n", irow)
	}

	err = rows.Err()
	if err != io.EOF {
		t.Errorf("error: expected io.EOF. got=%v\n", err)
	}
}

func TestCSVReaderScanUntilEOF(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, -1)
	if err != nil {
		t.Errorf("could read rows [0, -1): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		data := struct {
			I int
			F float64
			S string
		}{}
		err = rows.Scan(&data)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("error reading row %d: %v\n", irow, err)
		}
		exp := fmt.Sprintf("%d;%d;str-%d", irow, irow, irow)
		got := fmt.Sprintf("%v;%v;%v", data.I, data.F, data.S)
		if exp != got {
			t.Errorf("error reading row %d\nexp=%q\ngot=%q\n",
				irow, exp, got,
			)
		}
		irow++
	}

	err = rows.Err()
	if err != io.EOF {
		t.Errorf("error: expected io.EOF. got=%v\n", err)
	}
}
