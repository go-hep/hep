// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

var (
	ErrNotExist = errors.New("ntuple does not exist")
)

// NTuple provides read/write access to row-wise data.
type NTuple struct {
	db     *sql.DB
	name   string
	schema []columnDescr
}

// OpenNTuple inspects the given database handle and tries to return
// an NTuple connected to a table with the given name.
// OpenNTuple returns ErrNotExist if no such table exists.
//
// e.g.:
//  db, err := sql.Open("csv", "file.csv")
//  nt, err := hbook.OpenNTuple(db, "ntup")
func OpenNTuple(db *sql.DB, name string) (*NTuple, error) {
	nt := &NTuple{
		db:   db,
		name: name,
	}
	// FIXME(sbinet) test whether the table 'name' actually exists
	// FIXME(sbinet) retrieve underlying schema from db
	return nt, nil
}

// CreateNTuple creates a new ntuple with the given name inside the given database handle.
// The n-tuple schema is inferred from the cols argument. cols can be:
//  - a single struct value (columns are inferred from the names+types of the exported fields)
//  - a list of builtin values (the columns names are varX where X=[1-len(cols)])
//  - a list of hbook.Descriptors
//
// e.g.:
//  nt, err := hbook.CreateNTuple(db, "nt", struct{X float64 `hbook:"x"`}{})
//  nt, err := hbook.CreateNTuple(db, "nt", int64(0), float64(0))
func CreateNTuple(db *sql.DB, name string, cols ...interface{}) (*NTuple, error) {
	nt := &NTuple{
		db:   db,
		name: name,
	}
	return nt, nil
}

type columnDescr struct {
	Name string
	Type reflect.Type
}

// Scan executes a query against the ntuple and runs the function f against that context.
//
// e.g.
//  err = nt.Scan("x,y where z>10", func(x,y float64) error {
//    h1.Fill(x, 1)
//    h2.Fill(y, 1)
//    return nil
//  })
func (nt *NTuple) Scan(query string, f interface{}) error {
	rv := reflect.ValueOf(f)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return fmt.Errorf("hbook: expected a func, got %T", f)
	}
	if rt.NumOut() != 1 || rt.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("hbook: expected a func returning an error. got %T", f)
	}
	vargs := make([]reflect.Value, rt.NumIn())
	args := make([]interface{}, rt.NumIn())
	for i := range args {
		ptr := reflect.New(rt.In(i))
		args[i] = ptr.Interface()
		vargs[i] = ptr.Elem()
	}

	query, err := nt.massageQuery(query)
	if err != nil {
		return err
	}

	rows, err := nt.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(args...)
		if err != nil {
			return err
		}

		out := rv.Call(vargs)[0].Interface()
		if out != nil {
			return out.(error)
		}
	}

	err = rows.Err()
	if err == io.EOF {
		err = nil
	}
	return err
}

// ScanH1D executes a query against the ntuple and fills the histogram with
// the results of the query.
// If h is nil, a (100-bins, xmin, xmax) histogram is created,
// where xmin and xmax are inferred from the content of the underlying database.
func (nt *NTuple) ScanH1D(query string, h *H1D) (*H1D, error) {
	var data []float64
	var (
		xmin = +math.MaxFloat64
		xmax = -math.MaxFloat64
	)
	err := nt.Scan(query, func(x float64) error {
		data = append(data, x)
		if xmin > x {
			xmin = x
		}
		if xmax < x {
			xmax = x
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if h == nil {
		h = NewH1D(100, xmin, xmax)
	}

	for _, x := range data {
		h.Fill(x, 1)
	}
	return h, err
}

/*
// Scatter1D
func Scatter1D(db *sql.DB, query string) error {
	var err error
	return err
}
*/

func (nt *NTuple) massageQuery(q string) (string, error) {
	const (
		tokWHERE = " WHERE "
		tokWhere = " where "
	)
	vars := q
	where := ""
	switch {
	case strings.Contains(q, tokWHERE):
		toks := strings.Split(q, tokWHERE)
		vars = toks[0]
		where = " where " + toks[1]
	case strings.Contains(q, tokWhere):
		toks := strings.Split(q, tokWhere)
		vars = toks[0]
		where = " where " + toks[1]
	}

	// FIXME(sbinet) this is vulnerable to SQL injections...
	return "select " + vars + " from " + nt.name + where, nil
}
