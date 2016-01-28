// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvdriver

import (
	"database/sql"
	"database/sql/driver"
	"io"
)

type csvRows struct {
	rows *sql.Rows
}

// Columns returns the names of the columns. The number of
// columns of the result is inferred from the length of the
// slice.  If a particular column name isn't known, an empty
// string should be returned for that entry.
func (rows *csvRows) Columns() []string {
	cols, err := rows.rows.Columns()
	if err != nil {
		return nil
	}
	return cols
}

// Close closes the rows iterator.
func (rows *csvRows) Close() error {
	return rows.rows.Close()
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// The dest slice may be populated only with
// a driver Value type, but excluding string.
// All string values must be converted to []byte.
//
// Next should return io.EOF when there are no more rows.
func (rows *csvRows) Next(dest []driver.Value) error {
	if !rows.rows.Next() {
		return io.EOF
	}
	args := params(dest)
	return rows.rows.Scan(args...)
}
