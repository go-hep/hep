// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil_test

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"testing"

	"go-hep.org/x/hep/csvutil"
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

func TestCSVReaderInvalidScan(t *testing.T) {
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

	if !rows.Next() {
		t.Fatalf("could not get data")
	}

	err = rows.Scan()
	if err == nil {
		t.Errorf("expected an error")
	}

	err = rows.Err()
	if err == nil {
		t.Fatalf("expected a sticky error")
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

func TestCSVReaderScanArgsSubSample(t *testing.T) {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		t.Errorf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(2, 10)
	if err != nil {
		t.Errorf("could read rows [2, 10): %v\n", err)
	}
	defer rows.Close()

	irow := 2
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

	if irow-2 != 8 {
		t.Errorf("error: got %d rows. expected 8\n", irow-2)
	}
}

func TestCSVWriterArgs(t *testing.T) {
	fname := "testdata/out-args.csv"
	tbl, err := csvutil.Create(fname)
	if err != nil {
		t.Errorf("could not create %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	err = tbl.WriteHeader("## a simple set of data: int64;float64;string\n")
	if err != nil {
		t.Errorf("error writing header: %v\n", err)
	}

	for i := 0; i < 10; i++ {
		var (
			f = float64(i)
			s = fmt.Sprintf("str-%d", i)
		)
		err = tbl.WriteRow(i, f, s)
		if err != nil {
			t.Errorf("error writing row %d: %v\n", i, err)
			break
		}
	}

	err = tbl.Close()
	if err != nil {
		t.Errorf("error closing table: %v\n", err)
	}

	err = diff("testdata/simple.csv", fname)
	if err != nil {
		t.Errorf("files differ: %v\n", err)
	}
}

func TestCSVWriterStruct(t *testing.T) {
	fname := "testdata/out-struct.csv"
	tbl, err := csvutil.Create(fname)
	if err != nil {
		t.Errorf("could not create %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	// test WriteHeader w/o a trailing newline
	err = tbl.WriteHeader("## a simple set of data: int64;float64;string")
	if err != nil {
		t.Errorf("error writing header: %v\n", err)
	}

	for i := 0; i < 10; i++ {
		data := struct {
			I int
			F float64
			S string
		}{
			I: i,
			F: float64(i),
			S: fmt.Sprintf("str-%d", i),
		}
		err = tbl.WriteRow(data)
		if err != nil {
			t.Errorf("error writing row %d: %v\n", i, err)
			break
		}
	}

	err = tbl.Close()
	if err != nil {
		t.Errorf("error closing table: %v\n", err)
	}

	err = diff("testdata/simple.csv", fname)
	if err != nil {
		t.Errorf("files differ: %v\n", err)
	}
}

func TestCSVAppend(t *testing.T) {
	fname := "testdata/append-test.csv"
	tbl, err := csvutil.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer tbl.Close()

	tbl.Writer.Comma = ';'

	// test WriteHeader w/o a trailing newline
	err = tbl.WriteHeader("## a simple set of data: int64;float64;string")
	if err != nil {
		t.Errorf("error writing header: %v\n", err)
	}

	for i := 0; i < 10; i++ {
		data := struct {
			I int
			F float64
			S string
		}{
			I: i,
			F: float64(i),
			S: fmt.Sprintf("str-%d", i),
		}
		err = tbl.WriteRow(data)
		if err != nil {
			t.Errorf("error writing row %d: %v\n", i, err)
			break
		}
	}

	err = tbl.Close()
	if err != nil {
		t.Errorf("error closing table: %v\n", err)
	}

	// re-open to append
	tbl, err = csvutil.Append(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer tbl.Close()

	tbl.Writer.Comma = ';'
	for i := 10; i < 20; i++ {
		data := struct {
			I int
			F float64
			S string
		}{
			I: i,
			F: float64(i),
			S: fmt.Sprintf("str-%d", i),
		}
		err = tbl.WriteRow(data)
		if err != nil {
			t.Errorf("error writing row %d: %v\n", i, err)
			break
		}
	}

	err = tbl.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = diff("testdata/append.csv", fname)
	if err != nil {
		t.Errorf("files differ: %v\n", err)
	}
}

func TestCSVReaderTypes(t *testing.T) {
	fname := "testdata/types.csv"
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

	const nfields = 14
	type Data struct {
		Bool   bool
		Int    int
		Int8   int8
		Int16  int16
		Int32  int32
		Int64  int64
		UInt   uint
		UInt8  uint8
		UInt16 uint16
		UInt32 uint32
		UInt64 uint64
		F32    float32
		F64    float64
		Str    string
	}

	wants := []Data{
		{true, +1, -1, -1, -1, -1, +1, +1, +1, +1, +1, 1.1, 1.1, "str-1"},
		{false, -2, -2, -2, -2, -2, +2, +2, +2, +2, +2, 2.2, 2.2, "str-2"},
	}
	irow := 0
	for rows.Next() {
		want := wants[irow]
		{
			var got Data
			err = rows.Scan(&got.Bool, &got.Int, &got.Int8, &got.Int16, &got.Int32, &got.Int64, &got.UInt, &got.UInt8, &got.UInt16, &got.UInt32, &got.UInt64, &got.F32, &got.F64, &got.Str)
			if err != nil {
				t.Errorf("error reading row %d: %v\n", irow, err)
			}
			if want != got {
				t.Errorf("error reading row %d\ngot= %#v\nwant=%#v\n",
					irow, got, want,
				)
			}
			if got, want := rows.NumFields(), nfields; got != want {
				t.Errorf("invalid number of fields. got=%d. want=%d", got, want)

			}
			if got, want := len(rows.Fields()), nfields; got != want {
				t.Errorf("invalid number of fields. got=%d. want=%d", got, want)

			}
		}
		{
			var got Data
			err = rows.Scan(&got)
			if err != nil {
				t.Errorf("error reading row %d: %v\n", irow, err)
			}
			if want != got {
				t.Errorf("error reading row %d\ngot= %#v\nwant=%#v\n",
					irow, got, want,
				)
			}
			if got, want := rows.NumFields(), nfields; got != want {
				t.Errorf("invalid number of fields. got=%d. want=%d", got, want)
			}
			if got, want := len(rows.Fields()), nfields; got != want {
				t.Errorf("invalid number of fields. got=%d. want=%d", got, want)
			}
		}
		irow++
	}

	err = rows.Err()
	if err != nil {
		t.Errorf("error iterating over rows: %v\n", err)
	}
}

func TestCSVWriterTypes(t *testing.T) {
	fname := "testdata/out-types.csv"
	tbl, err := csvutil.Create(fname)
	if err != nil {
		t.Errorf("could not create %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	// test WriteHeader w/o a trailing newline
	err = tbl.WriteHeader("## supported types: bool;int;int8;int16;int32;int64;uint;uint8;uint16;uint32;uint64;float32;float64;string")
	if err != nil {
		t.Errorf("error writing header: %v\n", err)
	}

	type Data struct {
		Bool   bool
		Int    int
		Int8   int8
		Int16  int16
		Int32  int32
		Int64  int64
		UInt   uint
		UInt8  uint8
		UInt16 uint16
		UInt32 uint32
		UInt64 uint64
		F32    float32
		F64    float64
		Str    string
	}

	wants := []Data{
		{true, +1, -1, -1, -1, -1, +1, +1, +1, +1, +1, 1.1, 1.1, "str-1"},
		{false, -2, -2, -2, -2, -2, +2, +2, +2, +2, +2, 2.2, 2.2, "str-2"},
	}
	for i := range wants {
		err = tbl.WriteRow(wants[i])
		if err != nil {
			t.Errorf("error writing row %d: %v\n", i, err)
			break
		}
	}

	err = tbl.Close()
	if err != nil {
		t.Errorf("error closing table: %v\n", err)
	}

	err = diff("testdata/types.csv.ref", fname)
	if err != nil {
		t.Errorf("files differ: %v\n", err)
	}
}

func diff(ref, chk string) error {
	cmd := exec.Command("diff", "-urN", ref, chk)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("diff %v %v failed: %v\n%v\n",
			ref, chk, err,
			buf.String(),
		)
	}
	return nil
}
