// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvdriver

import (
	"database/sql"
	"database/sql/driver"
)

var (
	_ driver.Stmt = (*csvStmt)(nil)
)

type csvStmt struct {
	conn  *csvConn
	query string
	stmt  *sql.Stmt
}

// Close closes the statement.
//
// As of Go 1.1, a Stmt will not be closed if it's in use
// by any queries.
func (stmt *csvStmt) Close() error {
	return stmt.stmt.Close()
}

// NumInput returns the number of placeholder parameters.
//
// If NumInput returns >= 0, the sql package will sanity check
// argument counts from callers and return errors to the caller
// before the statement's Exec or Query methods are called.
//
// NumInput may also return -1, if the driver doesn't know
// its number of placeholders. In that case, the sql package
// will not sanity check Exec or Query argument counts.
func (stmt *csvStmt) NumInput() int {
	return -1
}

// Exec executes a query that doesn't return rows, such
// as an INSERT or UPDATE.
func (stmt *csvStmt) Exec(args []driver.Value) (driver.Result, error) {
	return stmt.stmt.Exec(params(args)...)
}

// Query executes a query that may return rows, such as a
// SELECT.
func (stmt *csvStmt) Query(args []driver.Value) (driver.Rows, error) {
	rows, err := stmt.stmt.Query(params(args)...)
	if err != nil {
		return nil, err
	}
	return &csvRows{rows}, err
}
