// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package hplt provides helper functions to plot histograms from n-tuples.
package hplt

import (
	"database/sql"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/go-hep/hbook"
)

// Plot executes a query against the given db and runs the function f against that context.
//
// e.g.
//  err = hplt.Plot(db, "select (x,y) from table where z>10", func(x,y float64) error {
//    h1.Fill(x, 1)
//    h2.Fill(y, 1)
//    return nil
//  })
func Plot(db *sql.DB, query string, f interface{}) error {
	rv := reflect.ValueOf(f)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return fmt.Errorf("hplt: expected a func, got %T", f)
	}
	if rt.NumOut() != 1 || rt.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return fmt.Errorf("hplt: expected a func returning an error. got %T", f)
	}
	vargs := make([]reflect.Value, rt.NumIn())
	args := make([]interface{}, rt.NumIn())
	for i := range args {
		ptr := reflect.New(rt.In(i))
		args[i] = ptr.Interface()
		vargs[i] = ptr.Elem()
	}

	rows, err := db.Query(query)
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

// Plot1D executes a query against the given db and fills the histogram with
// the results of the query.
// If h is nil, a (100-bins, xmin, xmax) histogram is created,
// where xmin and xmax are inferred from the content of the underlying database.
func Plot1D(db *sql.DB, query string, h *hbook.H1D) (*hbook.H1D, error) {
	var data []float64
	var (
		xmin = +math.MaxFloat64
		xmax = -math.MaxFloat64
	)
	err := Plot(db, query, func(x float64) error {
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
		h = hbook.NewH1D(100, xmin, xmax)
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
