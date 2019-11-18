// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"sync/atomic"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
)

// NewTable creates a new in-memory Arrow Table from the provided ROOT Tree.
func NewTable(t rtree.Tree, opts ...Option) array.Table {
	cfg := newConfig(opts)

	tbl := &rootTable{
		mem:    cfg.mem,
		tree:   t,
		refs:   1,
		schema: SchemaFrom(t),
		nrows:  t.Entries(),
		ncols:  int64(len(t.Branches())),
		cols:   make([]*array.Column, len(t.Branches())),
	}

	tbl.init()

	return tbl
}

type rootTable struct {
	mem  memory.Allocator
	tree rtree.Tree

	refs   int64
	schema *arrow.Schema
	nrows  int64
	ncols  int64

	cols []*array.Column
}

func (tbl *rootTable) Schema() *arrow.Schema      { return tbl.schema }
func (tbl *rootTable) NumRows() int64             { return tbl.nrows }
func (tbl *rootTable) NumCols() int64             { return tbl.ncols }
func (tbl *rootTable) Column(i int) *array.Column { return tbl.cols[i] }

// Retain increases the reference count by 1.
// Retain may be called simultaneously from multiple goroutines.
func (tbl *rootTable) Retain() {
	atomic.AddInt64(&tbl.refs, 1)
}

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
// Release may be called simultaneously from multiple goroutines.
func (tbl *rootTable) Release() {
	if atomic.LoadInt64(&tbl.refs) <= 0 {
		panic("groot/rarrow: too many releases")
	}

	if atomic.AddInt64(&tbl.refs, -1) == 0 {
		for i := range tbl.cols {
			tbl.cols[i].Release()
		}
		tbl.cols = nil
	}
}

func (tbl *rootTable) init() {
	// FIXME(sbinet): infer clusters sizes
	// FIXME(sbinet): lazily populate rootTable

	vars := rtree.NewScanVars(tbl.tree)
	sc, err := rtree.NewScannerVars(tbl.tree, vars...)
	if err != nil {
		panic(xerrors.Errorf("could not create scanner from scan-vars %#v: %w", vars, err))
	}
	defer sc.Close()

	arrs := make([]array.Interface, tbl.ncols)
	blds := make([]array.Builder, tbl.ncols)
	for i, field := range tbl.schema.Fields() {
		blds[i] = builderFrom(tbl.mem, field.Type, tbl.nrows)
		defer blds[i].Release()
	}

	for sc.Next() {
		err := sc.Scan()
		if err != nil {
			panic(xerrors.Errorf("could not scan entry %d: %w", sc.Entry(), err))
		}

		for i, field := range tbl.schema.Fields() {
			appendData(blds[i], vars[i], field.Type)
		}
	}

	for i, bldr := range blds {
		arrs[i] = bldr.NewArray()
		defer arrs[i].Release()
	}

	tbl.cols = make([]*array.Column, tbl.ncols)
	for i, arr := range arrs {
		field := tbl.schema.Field(i)
		if !arrow.TypeEquals(field.Type, arr.DataType()) {
			panic(xerrors.Errorf("field[%d][%s]: type=%v|%v array=%v", i, field.Name, field.Type, arr.DataType(), arr))
		}
		chunked := array.NewChunked(field.Type, []array.Interface{arr})
		defer chunked.Release()
		tbl.cols[i] = array.NewColumn(field, chunked)
	}
}

var (
	_ array.Table = (*rootTable)(nil)
)
