// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"go-hep.org/x/hep/groot/riofs"
)

// Connector returns a database/sql/driver.Connector from a ROOT file.
//
// Connector can be used to open a database/sql.DB from an already
// open ROOT file.
func Connector(file *riofs.File) driver.Connector {
	c := &rootConnector{file: file}
	return c
}

type rootConnector struct {
	drv  rootDriver
	file *riofs.File
}

// Connect returns a connection to the database.
// Connect may return a cached connection (one previously
// closed), but doing so is unnecessary; the sql package
// maintains a pool of idle connections for efficient re-use.
//
// The provided context.Context is for dialing purposes only
// (see net.DialContext) and should not be stored or used for
// other purposes.
//
// The returned connection is only used by one goroutine at a
// time.
func (c *rootConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return c.drv.connect(c.file), nil
}

// Driver returns the underlying Driver of the Connector,
// mainly to maintain compatibility with the Driver method
// on sql.DB.
func (c *rootConnector) Driver() driver.Driver {
	return &c.drv
}

// OpenDB opens a database/sql.DB from an already open ROOT file.
func OpenDB(file *riofs.File) *sql.DB {
	return sql.OpenDB(Connector(file))
}

var (
	_ driver.Connector = (*rootConnector)(nil)
)
