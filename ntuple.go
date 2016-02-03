// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"math"
	"reflect"
	"strings"
)

var (
	ErrNotExist      = errors.New("hbook: ntuple does not exist")
	ErrMissingColDef = errors.New("hbook: expected at least one column definition")
)

// Ntuple provides read/write access to row-wise data.
type Ntuple struct {
	db     *sql.DB
	name   string
	schema []Descriptor
}

// OpenNtuple inspects the given database handle and tries to return
// an Ntuple connected to a table with the given name.
// OpenNtuple returns ErrNotExist if no such table exists.
// If name is "", OpenNtuple will connect to the one-and-only table in the db.
//
// e.g.:
//  db, err := sql.Open("csv", "file.csv")
//  nt, err := hbook.OpenNtuple(db, "ntup")
func OpenNtuple(db *sql.DB, name string) (*Ntuple, error) {
	nt := &Ntuple{
		db:   db,
		name: name,
	}
	// FIXME(sbinet) test whether the table 'name' actually exists
	// FIXME(sbinet) retrieve underlying schema from db
	return nt, nil
}

// CreateNtuple creates a new ntuple with the given name inside the given database handle.
// The n-tuple schema is inferred from the cols argument. cols can be:
//  - a single struct value (columns are inferred from the names+types of the exported fields)
//  - a list of builtin values (the columns names are varX where X=[1-len(cols)])
//  - a list of hbook.Descriptors
//
// e.g.:
//  nt, err := hbook.CreateNtuple(db, "nt", struct{X float64 `hbook:"x"`}{})
//  nt, err := hbook.CreateNtuple(db, "nt", int64(0), float64(0))
func CreateNtuple(db *sql.DB, name string, cols ...interface{}) (*Ntuple, error) {
	var err error
	nt := &Ntuple{
		db:   db,
		name: name,
	}
	var schema []Descriptor
	switch len(cols) {
	case 0:
		return nil, ErrMissingColDef
	case 1:
		rv := reflect.Indirect(reflect.ValueOf(cols[0]))
		rt := rv.Type()
		switch rt.Kind() {
		case reflect.Struct:
			schema, err = schemaFromStruct(rt)
		default:
			schema, err = schemaFrom(cols...)
		}
	default:
		schema, err = schemaFrom(cols...)
	}
	if err != nil {
		return nil, err
	}
	nt.schema = schema
	return nt, err
}

// Name returns the name of this n-tuple.
func (nt *Ntuple) Name() string {
	return nt.name
}

// Cols returns the columns' descriptors of this n-tuple.
// Modifying it directly leads to undefined behaviour.
func (nt *Ntuple) Cols() []Descriptor {
	return nt.schema
}

// Descriptor describes a column
type Descriptor interface {
	Name() string       // the column name
	Type() reflect.Type // the column type
}

type columnDescr struct {
	name string
	typ  reflect.Type
}

func (col *columnDescr) Name() string {
	return col.name
}

func (col *columnDescr) Type() reflect.Type {
	return col.typ
}

func schemaFromStruct(rt reflect.Type) ([]Descriptor, error) {
	var schema []Descriptor
	var err error
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !ast.IsExported(f.Name) {
			continue
		}
		ft := f.Type
		switch ft.Kind() {
		case reflect.Chan:
			return nil, fmt.Errorf("hbook: chans not supported")
		case reflect.Interface:
			return nil, fmt.Errorf("hbook: interfaces not supported")
		case reflect.Map:
			return nil, fmt.Errorf("hbook: maps not supported")
		case reflect.Slice:
			return nil, fmt.Errorf("hbook: nested slices not supported")
		case reflect.Struct:
			return nil, fmt.Errorf("hbook: nested structs not supported")
		}
		fname := getTag(f.Tag, "hbook", "rio", "db")
		if fname == "" {
			fname = f.Name
		}
		schema = append(schema, &columnDescr{fname, ft})
	}
	return schema, err
}

func schemaFrom(src ...interface{}) ([]Descriptor, error) {
	var schema []Descriptor
	var err error
	for i, col := range src {
		rt := reflect.TypeOf(col)
		switch rt.Kind() {
		case reflect.Chan:
			return nil, fmt.Errorf("hbook: chans not supported")
		case reflect.Interface:
			return nil, fmt.Errorf("hbook: interfaces not supported")
		case reflect.Map:
			return nil, fmt.Errorf("hbook: maps not supported")
		case reflect.Slice:
			return nil, fmt.Errorf("hbook: slices not supported")
		case reflect.Struct:
			return nil, fmt.Errorf("hbook: structs not supported")
		}
		schema = append(schema, &columnDescr{fmt.Sprintf("var%d", i+1), rt})
	}
	return schema, err
}

func getTag(tag reflect.StructTag, keys ...string) string {
	for _, k := range keys {
		v := tag.Get(k)
		if v != "" && v != "-" {
			return v
		}
	}
	return ""
}

// Scan executes a query against the ntuple and runs the function f against that context.
//
// e.g.
//  err = nt.Scan("x,y where z>10", func(x,y float64) error {
//    h1.Fill(x, 1)
//    h2.Fill(y, 1)
//    return nil
//  })
func (nt *Ntuple) Scan(query string, f interface{}) error {
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
func (nt *Ntuple) ScanH1D(query string, h *H1D) (*H1D, error) {
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

func (nt *Ntuple) massageQuery(q string) (string, error) {
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
