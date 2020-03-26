// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"fmt"
	"sync/atomic"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"go-hep.org/x/hep/groot/rtree"
)

// Record is an in-memory Arrow Record backed by a ROOT Tree.
type Record struct {
	refs int64

	mem  memory.Allocator
	tree rtree.Tree

	schema *arrow.Schema
	nrows  int64
	ncols  int64
	offset int64 // entries offset

	cols []array.Interface
}

// NewRecord creates a new in-memory Arrow Record from the provided ROOT Tree.
func NewRecord(t rtree.Tree, opts ...Option) *Record {
	cfg := newConfig(opts)

	if cfg.end < 0 {
		cfg.end = t.Entries()
	}

	if cfg.beg <= 0 {
		cfg.beg = 0
	}

	if cfg.beg > cfg.end {
		panic("rarrow: invalid entry slice")
	}

	rec := &Record{
		mem:    cfg.mem,
		tree:   t,
		refs:   1,
		schema: SchemaFrom(t),
		offset: cfg.beg,
		nrows:  cfg.end - cfg.beg,
		ncols:  int64(len(t.Branches())),
		cols:   make([]array.Interface, len(t.Branches())),
	}

	rec.load(cfg.beg, cfg.end)

	return rec
}

func (rec *Record) load(beg, end int64) {
	vars := rtree.NewReadVars(rec.tree)
	sc, err := rtree.NewScannerVars(rec.tree, vars...)
	if err != nil {
		panic(fmt.Errorf("could not create scanner from read-vars %#v: %w", vars, err))
	}
	defer sc.Close()

	blds := make([]array.Builder, rec.ncols)
	for i, field := range rec.schema.Fields() {
		blds[i] = builderFrom(rec.mem, field.Type, rec.nrows)
		defer blds[i].Release()
	}

	err = sc.SeekEntry(beg)
	if err != nil {
		panic(fmt.Errorf("could not seek to entry: %w", err))
	}

	n := beg
	for sc.Next() {
		err := sc.Scan()
		if err != nil {
			panic(fmt.Errorf("could not scan entry %d: %w", sc.Entry(), err))
		}

		for i, field := range rec.schema.Fields() {
			appendData(blds[i], vars[i], field.Type)
		}

		n++
		if n >= end {
			break
		}
	}

	for i, bldr := range blds {
		rec.cols[i] = bldr.NewArray()
	}
}

// Retain increases the reference count by 1.
// Retain may be called simultaneously from multiple goroutines.
func (rec *Record) Retain() {
	atomic.AddInt64(&rec.refs, 1)
}

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
// Release may be called simultaneously from multiple goroutines.
func (rec *Record) Release() {
	if atomic.LoadInt64(&rec.refs) <= 0 {
		panic("groot/rarrow: too many releases")
	}

	if atomic.AddInt64(&rec.refs, -1) == 0 {
		for i := range rec.cols {
			rec.cols[i].Release()
		}
		rec.cols = nil
	}
}

func (rec *Record) Schema() *arrow.Schema        { return rec.schema }
func (rec *Record) NumRows() int64               { return rec.nrows }
func (rec *Record) NumCols() int64               { return rec.ncols }
func (rec *Record) Columns() []array.Interface   { return rec.cols }
func (rec *Record) Column(i int) array.Interface { return rec.cols[i] }
func (rec *Record) ColumnName(i int) string      { return rec.schema.Field(i).Name }

// NewSlice constructs a zero-copy slice of the record with the indicated
// indices i and j, corresponding to array[i:j].
// The returned record must be Release()'d after use.
//
// NewSlice panics if the slice is outside the valid range of the record array.
// NewSlice panics if j < i.
func (rec *Record) NewSlice(i, j int64) array.Record {
	return NewRecord(rec.tree, WithStart(rec.offset+i), WithEnd(rec.offset+j))
}

// RecordReader is an ARROW RecordReader for ROOT Trees.
//
// RecordReader does not materialize more than one record at a time.
// The number of rows (or entries, in ROOT speak) that record loads can be configured
// at creation time with the WithChunk function.
// The default is one entry per record.
// One can pass -1 to WithChunk to create a record with all entries of the Tree or Chain.
type RecordReader struct {
	refs int64

	mem    memory.Allocator
	schema *arrow.Schema
	tree   rtree.Tree

	beg   int64 // first entry to read
	end   int64 // last entry to read
	cur   int64 // current entry
	chunk int64 // number of entries to read for each record

	rec *Record
}

// NewRecordReader creates a new ARROW RecordReader from the provided ROOT Tree.
func NewRecordReader(tree rtree.Tree, opts ...Option) *RecordReader {
	cfg := newConfig(opts)

	r := &RecordReader{
		refs:   1,
		mem:    cfg.mem,
		schema: SchemaFrom(tree),
		tree:   tree,
		beg:    cfg.beg,
		end:    cfg.end,
		chunk:  cfg.chunks,
	}

	if r.beg <= 0 {
		r.beg = 0
	}

	if r.end <= 0 {
		r.end = tree.Entries()
	}

	switch {
	case r.chunk == 0:
		r.chunk = 1
	case r.chunk < 0:
		r.chunk = tree.Entries()
	}
	r.cur = r.beg

	return r
}

// Retain increases the reference count by 1.
// Retain may be called simultaneously from multiple goroutines.
func (r *RecordReader) Retain() {
	atomic.AddInt64(&r.refs, 1)
}

// Release decreases the reference count by 1.
// When the reference count goes to zero, the memory is freed.
// Release may be called simultaneously from multiple goroutines.
func (r *RecordReader) Release() {
	if atomic.LoadInt64(&r.refs) <= 0 {
		panic("groot/rarrow: too many releases")
	}

	if atomic.AddInt64(&r.refs, -1) == 0 {
		if r.rec != nil {
			r.rec.Release()
		}
	}
}

func (r *RecordReader) Schema() *arrow.Schema { return r.schema }
func (r *RecordReader) Record() array.Record  { return r.rec }

func (r *RecordReader) Next() bool {
	if r.cur >= r.end {
		return false
	}

	if r.rec != nil {
		r.rec.Release()
	}

	end := minI64(r.cur+r.chunk, r.end)
	r.load(r.cur, end)
	r.cur += r.chunk
	return true
}

func (r *RecordReader) load(beg, end int64) {
	r.rec = NewRecord(r.tree, WithStart(beg), WithEnd(end), WithAllocator(r.mem))
}

var (
	_ array.Record       = (*Record)(nil)
	_ array.RecordReader = (*RecordReader)(nil)
)

func minI64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
