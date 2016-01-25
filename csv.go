// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Open(fname string) (*Table, error) {
	r, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	table := &Table{
		Reader: csv.NewReader(bufio.NewReader(r)),
		f:      r,
	}
	return table, err
}

type Table struct {
	Reader *csv.Reader

	//TODO(sbinet) add a writer
	//Writer *csv.Writer

	f      *os.File
	closed bool
	err    error
}

func (tbl *Table) Close() error {
	if tbl.closed {
		return tbl.err
	}

	if tbl.f != nil {
		tbl.err = tbl.f.Close()
		tbl.f = nil
		tbl.closed = true
	}
	return tbl.err
}

func (tbl *Table) ReadRows(beg, end int64) (*Rows, error) {
	inc := int64(1)
	rows := &Rows{
		tbl: tbl,
		i:   0,
		n:   end - beg,
		inc: inc,
		cur: beg - inc,
	}
	if end == -1 {
		rows.n = math.MaxInt64
	}
	return rows, nil
}

type Rows struct {
	tbl    *Table
	i      int64    // number of rows iterated over
	n      int64    // number of rows this iterator iters over
	inc    int64    // number of rows to increment by at each iteration
	cur    int64    // current row index
	record []string // last read record
	closed bool
	err    error // last error
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (rows *Rows) Err() error {
	return rows.err
}

// Close closes the Rows, preventing further enumeration.
// Close is idempotent and does not affect the result of Err.
func (rows *Rows) Close() error {
	if rows.closed {
		return nil
	}
	rows.closed = true
	rows.tbl = nil
	return nil
}

// Scan copies the columns in the current row into the values pointed at by
// dest.
func (rows *Rows) Scan(dest ...interface{}) error {
	var err error
	defer func() {
		rows.err = err
	}()

	rows.record, err = rows.tbl.Reader.Read()
	if err != nil {
		return err
	}

	switch len(dest) {
	case 0:
		err = fmt.Errorf("csv: Rows.Scan needs at least one argument")
		return err

	case 1:
		// maybe special case: struct?
		rt := reflect.TypeOf(dest[0]).Elem()
		switch rt.Kind() {
		case reflect.Struct:
			err = rows.scanStruct(dest[0])
			return err
		}
	}

	err = rows.scan(dest...)
	return err
}

func (rows *Rows) scan(args ...interface{}) error {
	var err error
	n := min(len(rows.record), len(args))
	for i := 0; i < n; i++ {
		rec := rows.record[i]
		rv := reflect.ValueOf(args[i]).Elem()
		rt := reflect.TypeOf(args[i]).Elem()
		switch rt.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(rec)
			if err != nil {
				return err
			}
			rv.SetBool(v)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(rec, 10, rt.Bits())
			if err != nil {
				return err
			}
			rv.SetInt(v)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(rec, 10, rt.Bits())
			if err != nil {
				return err
			}
			rv.SetUint(v)

		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(rec, rt.Bits())
			if err != nil {
				return err
			}
			rv.SetFloat(v)

		case reflect.String:
			rv.SetString(rec)

		default:
			return fmt.Errorf("csv: invalid type (%T) %q", rv.Interface(), rec)
		}
	}

	return err
}

func (rows *Rows) scanStruct(ptr interface{}) error {
	rt := reflect.TypeOf(ptr).Elem()
	rv := reflect.ValueOf(ptr).Elem()
	args := make([]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		args[i] = rv.Field(i).Addr().Interface()
	}
	return rows.scan(args...)
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (rows *Rows) Next() bool {
	if rows.closed {
		return false
	}
	if rows.err != nil {
		return false
	}
	next := rows.i < rows.n
	rows.cur += rows.inc
	rows.i += rows.inc
	if !next {
		rows.err = rows.Close()
	}
	return next
}
