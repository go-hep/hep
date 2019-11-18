// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rsql provides a convenient access to ROOT files/trees as a database.
package rsql // import "go-hep.org/x/hep/groot/rsql"

import (
	"io"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/rsql/rsqldrv"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"golang.org/x/xerrors"
)

// Scan executes a query against the given tree and runs the function f
// within that context.
func Scan(tree rtree.Tree, query string, f interface{}) error {
	if f == nil {
		return xerrors.Errorf("groot/rsql: nil func")
	}
	rv := reflect.ValueOf(f)
	rt := rv.Type()
	if rt.Kind() != reflect.Func {
		return xerrors.Errorf("groot/rsql: expected a func, got %T", f)
	}
	if rt.NumOut() != 1 || rt.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return xerrors.Errorf("groot/rsql: expected a func returning an error. got %T", f)
	}
	vargs := make([]reflect.Value, rt.NumIn())
	args := make([]interface{}, rt.NumIn())
	for i := range args {
		ptr := reflect.New(rt.In(i))
		args[i] = ptr.Interface()
		vargs[i] = ptr.Elem()
	}

	db := rsqldrv.OpenDB(rtree.FileOf(tree))
	defer db.Close()

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

// ScanH1D executes a query against the tree and fills the histogram with
// the results of the query.
// If h is nil, a (100-bins, xmin, xmax+ULP) histogram is created,
// where xmin and xmax are inferred from the content of the underlying database.
func ScanH1D(tree rtree.Tree, query string, h *hbook.H1D) (*hbook.H1D, error) {
	if h == nil {
		var (
			xmin = +math.MaxFloat64
			xmax = -math.MaxFloat64
		)
		// FIXME(sbinet) leverage the underlying db min/max functions,
		// instead of crawling through the whole data set.
		err := Scan(tree, query, func(x float64) error {
			xmin = math.Min(xmin, x)
			xmax = math.Max(xmax, x)
			return nil
		})
		if err != nil {
			return nil, err
		}

		h = hbook.NewH1D(100, xmin, nextULP(xmax))
	}

	err := Scan(tree, query, func(x float64) error {
		h.Fill(x, 1)
		return nil
	})

	return h, err
}

// ScanH2D executes a query against the ntuple and fills the histogram with
// the results of the query.
// If h is nil, a (100-bins, xmin, xmax+ULP) (100-bins, ymin, ymax+ULP) 2d-histogram
// is created,
// where xmin, xmax and ymin,ymax are inferred from the content of the
// underlying database.
func ScanH2D(tree rtree.Tree, query string, h *hbook.H2D) (*hbook.H2D, error) {
	if h == nil {
		var (
			xmin = +math.MaxFloat64
			xmax = -math.MaxFloat64
			ymin = +math.MaxFloat64
			ymax = -math.MaxFloat64
		)
		// FIXME(sbinet) leverage the underlying db min/max functions,
		// instead of crawling through the whole data set.
		err := Scan(tree, query, func(x, y float64) error {
			xmin = math.Min(xmin, x)
			xmax = math.Max(xmax, x)
			ymin = math.Min(ymin, y)
			ymax = math.Max(ymax, y)
			return nil
		})
		if err != nil {
			return nil, err
		}

		h = hbook.NewH2D(100, xmin, nextULP(xmax), 100, ymin, nextULP(ymax))
	}

	err := Scan(tree, query, func(x, y float64) error {
		h.Fill(x, y, 1)
		return nil
	})

	return h, err
}

func nextULP(v float64) float64 {
	return math.Nextafter(v, v+1)
}
