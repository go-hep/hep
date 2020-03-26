// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import "fmt"

// Reader reads data from a Tree.
type Reader struct {
	t     Tree
	rvars []ReadVar
	scan  *Scanner
	beg   int64
	end   int64
}

// ReadOption configures how a ROOT tree should be traversed.
type ReadOption func(r *Reader) error

// WithRange specifies the half-open interval [beg, end) of entries
// a Tree reader will read through.
func WithRange(beg, end int64) ReadOption {
	return func(r *Reader) error {
		r.beg = beg
		r.end = end
		return nil
	}
}

// NewReader creates a new Tree Reader from the provided ROOT Tree and
// the set of read-variables into which data will be read.
func NewReader(t Tree, rvars []ReadVar, opts ...ReadOption) (*Reader, error) {
	sc, err := NewScannerVars(t, rvars...)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not create scanner: %w", err)
	}

	r := Reader{
		t:     t,
		rvars: rvars,
		scan:  sc,
		beg:   0,
		end:   -1,
	}

	for i, opt := range opts {
		err := opt(&r)
		if err != nil {
			return nil, fmt.Errorf(
				"rtree: could not set reader option %d: %w",
				i, err,
			)
		}
	}

	if r.end < 0 {
		r.end = r.t.Entries()
	}

	if r.beg < 0 || r.beg > r.end {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d)",
			r.beg, r.end,
		)
	}

	if r.beg > r.t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (start=%d > tree-entries (%d))",
			r.beg, r.end, r.beg, r.t.Entries(),
		)
	}

	if r.end > r.t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (end=%d > tree-entries (%d))",
			r.beg, r.end, r.end, r.t.Entries(),
		)
	}

	return &r, nil
}

// Close closes the Reader.
func (r *Reader) Close() error {
	if r.scan == nil {
		return nil
	}
	err := r.scan.Close()
	r.scan = nil
	return err
}

// RCtx provides an entry-wise local context to the tree Reader.
type RCtx struct {
	Entry int64 // Current tree entry.
}

// Read will read data from the underlying tree over the whole specified range.
// Read calls the provided user function f for each entry successfully read.
func (r *Reader) Read(f func(ctx RCtx) error) error {
	err := r.scan.SeekEntry(r.beg)
	if err != nil {
		return fmt.Errorf("rtree: could not seek to entry %d: %w", r.beg, err)
	}

	for r.scan.Next() && r.scan.Entry() < r.end {
		iev := r.scan.Entry()
		err := r.scan.Scan()
		if err != nil {
			return fmt.Errorf("rtree: could not read entry %d: %w", iev, err)
		}

		err = f(RCtx{Entry: iev})
		if err != nil {
			return fmt.Errorf("rtree: could not process entry %d: %w", iev, err)
		}
	}

	err = r.scan.Err()
	if err != nil {
		return fmt.Errorf("rtree: could not traverse tree: %w", err)
	}

	return nil
}
