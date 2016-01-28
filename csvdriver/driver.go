// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package csvdriver registers a database/sql/driver.Driver implementation for CSV files.
package csvdriver

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/cznic/ql/driver"
)

var (
	_ driver.Driver  = (*csvDriver)(nil)
	_ driver.Conn    = (*csvConn)(nil)
	_ driver.Execer  = (*csvConn)(nil)
	_ driver.Queryer = (*csvConn)(nil)
	_ driver.Tx      = (*csvConn)(nil)
)

// Conn describes how a connection to the CSV-driver should be established.
type Conn struct {
	File    string      // name of the file to be open
	Mode    int         // r/w mode (default: read-only)
	Perm    os.FileMode // file permissions
	Comma   rune        // field delimiter (default: ',')
	Comment rune        // comment character for start of line (default: '#')
}

func (c *Conn) setDefaults() {
	if c.Mode == 0 {
		c.Mode = os.O_RDONLY
		c.Perm = 0
	}
	if c.Comma == 0 {
		c.Comma = ','
	}
	if c.Comment == 0 {
		c.Comment = '#'
	}
	return
}

func (c Conn) toJSON() (string, error) {
	c.setDefaults()
	buf, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

// Open opens a database connection with the CSV driver.
func (c Conn) Open() (*sql.DB, error) {
	c.setDefaults()
	str, err := c.toJSON()
	if err != nil {
		return nil, err
	}
	return sql.Open("csv", str)
}

// Open is a CSV-driver helper function for sql.Open.
//
// It opens a database connection to csvdriver.
func Open(name string) (*sql.DB, error) {
	c := Conn{File: name, Mode: os.O_RDONLY, Perm: 0}
	return c.Open()
}

// Create is a CSV-driver helper function for sql.Open.
//
// It creates a new CSV file, connected via the csvdriver.
func Create(name string) (*sql.DB, error) {
	c := Conn{
		File: name,
		Mode: os.O_RDWR | os.O_CREATE | os.O_TRUNC,
		Perm: 0666,
	}
	return c.Open()
}

type csvDriver struct{}

// Open returns a new connection to the database.
// The name is a string in a driver-specific format.
//
// Open may return a cached connection (one previously
// closed), but doing so is unnecessary; the sql package
// maintains a pool of idle connections for efficient re-use.
//
// The returned connection is only used by one goroutine at a
// time.
func (*csvDriver) Open(cfg string) (driver.Conn, error) {
	log.Printf(">>> driver.Open(%s)...\n", cfg)
	c := Conn{}
	if strings.HasPrefix(cfg, "{") {
		err := json.Unmarshal([]byte(cfg), &c)
		if err != nil {
			return nil, err
		}
	} else {
		c.File = cfg
		c.setDefaults()
	}

	doImport := false
	_, err := os.Lstat(c.File)
	if err == nil {
		doImport = true
	}

	f, err := os.OpenFile(c.File, c.Mode, c.Perm)
	if err != nil {
		return nil, err
	}
	conn := &csvConn{
		f:   f,
		cfg: c,
	}

	log.Printf(">>> doimport? %v\n", doImport)
	if doImport {
		err = conn.importCSV()
	} else {
		err = conn.initDB()
	}
	log.Printf(">>> doimport? %v [done]\n", doImport)

	if err != nil {
		return nil, err
	}

	log.Printf(">>> driver.Open(%s)... [done]\n", cfg)
	return conn, err
}

type csvConn struct {
	f   *os.File
	cfg Conn

	ql *sql.DB
	tx []driver.Tx
}

func (conn *csvConn) initDB() error {
	ql, err := qlopen(conn.cfg.File)
	if err != nil {
		return err
	}

	conn.ql = ql
	return nil
}

// Prepare returns a prepared statement, bound to this connection.
func (conn *csvConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := conn.ql.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &csvStmt{conn: conn, query: query, stmt: stmt}, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (conn *csvConn) Close() error {
	var err error
	defer conn.f.Close()

	// FIXME(sbinet) write-back to file if needed.
	// err = conn.exportCSV()

	err = conn.ql.Close()
	if err != nil {
		return err
	}

	err = conn.f.Close()
	if err != nil {
		return err
	}

	return err
}

// Begin starts and returns a new transaction.
func (conn *csvConn) Begin() (driver.Tx, error) {
	log.Printf(">>> conn.Begin()...\n")
	tx, err := conn.ql.Begin()
	if err != nil {
		log.Fatalf("conn-begin: %v\n", err)
		return nil, err
	}
	conn.tx = append(conn.tx, tx)
	return conn, err
}

func (conn *csvConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	log.Printf(">>> conn.Exec(%s)...\n", query)
	return conn.ql.Exec(query, params(args)...)
}

func (conn *csvConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	rows, err := conn.ql.Query(query, params(args)...)
	if err != nil {
		return nil, err
	}
	return &csvRows{rows}, err
}

func (conn *csvConn) Commit() error {
	ntx := len(conn.tx)
	if conn.tx == nil || ntx == 0 {
		return fmt.Errorf("csvdriver: commit while not in transaction")
	}
	tx := conn.tx[ntx-1]
	err := tx.Commit()
	conn.tx = conn.tx[:ntx-1]
	return err
}

func (conn *csvConn) Rollback() error {
	ntx := len(conn.tx)
	if conn.tx == nil || ntx == 0 {
		return fmt.Errorf("csvdriver: commit while not in transaction")
	}
	tx := conn.tx[ntx-1]
	err := tx.Rollback()
	conn.tx = conn.tx[:ntx-1]
	return err
}

func qlopen(name string) (*sql.DB, error) {
	db, err := sql.Open("ql", "memory://"+name)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	sql.Register("csv", &csvDriver{})
}
