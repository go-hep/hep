// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
)

// Reader reads data from a Tree.
type Reader struct {
	r    reader
	beg  int64
	end  int64
	nrab int // number of read-ahead baskets

	tree  Tree
	rvars []ReadVar

	evals []formula
	dirty bool // whether we need to re-create scanner (if formula needed new branches)
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

// WithPrefetchBaskets specifies the number of baskets to read-ahead, per branch.
// The default is 2.
// The number of prefetch baskets is cap'ed by the number of baskets, per branch.
func WithPrefetchBaskets(n int) ReadOption {
	return func(r *Reader) error {
		r.nrab = n
		return nil
	}
}

// NewReader creates a new Tree Reader from the provided ROOT Tree and
// the set of read-variables into which data will be read.
func NewReader(t Tree, rvars []ReadVar, opts ...ReadOption) (*Reader, error) {
	r := Reader{
		beg:  0,
		end:  -1,
		nrab: 2,
		tree: t,
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
		r.end = t.Entries()
	}

	if r.beg < 0 {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (start=%d < 0)",
			r.beg, r.end, r.beg,
		)
	}

	if r.beg > r.end {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (start=%d > end=%d)",
			r.beg, r.end, r.beg, r.end,
		)
	}

	if r.beg > t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (start=%d > tree-entries=%d)",
			r.beg, r.end, r.beg, t.Entries(),
		)
	}

	if r.end > t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (end=%d > tree-entries=%d)",
			r.beg, r.end, r.end, t.Entries(),
		)
	}

	rvars, err := sanitizeRVars(t, rvars)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not create reader: %w", err)
	}

	r.r = newReader(t, rvars, r.nrab, r.beg, r.end)
	r.rvars = r.r.rvars()

	return &r, nil
}

// Close closes the Reader.
func (r *Reader) Close() error {
	if r.r == nil {
		return nil
	}
	err := r.r.Close()
	r.r = nil
	r.evals = nil
	return err
}

// RCtx provides an entry-wise local context to the tree Reader.
type RCtx struct {
	Entry int64 // Current tree entry.
}

// Read will read data from the underlying tree over the whole specified range.
// Read calls the provided user function f for each entry successfully read.
func (r *Reader) Read(f func(ctx RCtx) error) error {
	if r.dirty {
		r.dirty = false
		_ = r.r.Close()
		r.r = newReader(r.tree, r.rvars, r.nrab, r.beg, r.end)
	}
	r.r.reset()

	const eoff = 0 // entry offset
	return r.r.run(eoff, r.beg, r.end, f)
}

type formula interface {
	eval()
}

// FormulaFunc creates a new formula based on the provided function and
// the list of branches as inputs.
func (r *Reader) FormulaFunc(branches []string, fct interface{}) (*FormulaFunc, error) {
	n := len(r.rvars)
	f, err := newFormulaFunc(r, branches, fct)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not create FormulaFunc: %w", err)
	}
	r.evals = append(r.evals, f)

	if n != len(r.rvars) {
		// formula needed to auto-load new branches.
		// mark reader as dirty to re-create its internal scanner
		// before the event-loop.
		r.dirty = true
	}
	return f, nil
}

func sanitizeRVars(t Tree, rvars []ReadVar) ([]ReadVar, error) {
	for i := range rvars {
		rvar := &rvars[i]
		if rvar.Leaf == "" {
			rvar.Leaf = rvar.Name
		}
		if rvar.count != "" {
			continue
		}
		br := t.Branch(rvar.Name)
		if br == nil {
			return nil, fmt.Errorf("rtree: tree %q has no branch named %q", t.Name(), rvar.Name)
		}
		leaf := br.Leaf(rvar.Leaf)
		if leaf == nil {
			continue
		}
		lfc := leaf.LeafCount()
		if lfc != nil {
			rvar.count = lfc.Name()
		}
	}
	return rvars, nil
}

type reader interface {
	Close() error
	rvars() []ReadVar

	run(off, beg, end int64, f func(RCtx) error) error
	reset()
}

// rtree reads a tree.
type rtree struct {
	tree *ttree
	rvs  []ReadVar
	brs  []rbranch
	lvs  []rleaf
}

var (
	_ reader = (*rtree)(nil)
)

func (r *rtree) rvars() []ReadVar { return r.rvs }

func newReader(t Tree, rvars []ReadVar, n int, beg, end int64) reader {
	rvars, err := sanitizeRVars(t, rvars)
	if err != nil {
		panic(err)
	}

	switch t := t.(type) {
	case *ttree:
		return newRTree(t, rvars, n, beg, end)
	case *tchain:
		return newRChain(t, rvars, n, beg, end)
	default:
		panic(fmt.Errorf("rtree: unknown Tree implementation %T", t))
	}
}

func newRTree(t *ttree, rvars []ReadVar, n int, beg, end int64) *rtree {
	r := &rtree{
		tree: t,
		rvs:  rvars,
	}
	usr := make(map[string]struct{}, len(rvars))
	for _, rvar := range rvars {
		usr[rvar.Name+"."+rvar.Leaf] = struct{}{}
	}

	var rcounts []ReadVar
	for _, rvar := range rvars {
		if rvar.count == "" {
			continue
		}
		leaf := t.Branch(rvar.Name).Leaf(rvar.Leaf).LeafCount()
		name := leaf.Branch().Name() + "." + leaf.Name()
		if _, ok := usr[name]; !ok {
			var ptr interface{}
			switch leaf := leaf.(type) {
			case *LeafB:
				ptr = new(int8)
			case *LeafS:
				ptr = new(int16)
			case *LeafI:
				ptr = new(int32)
			case *LeafL:
				ptr = new(int64)
			default:
				panic(fmt.Errorf("unknown Leaf count type %T", leaf))
			}
			rcounts = append(rcounts, ReadVar{
				Name:  leaf.Branch().Name(),
				Leaf:  leaf.Name(),
				Value: ptr,
				leaf:  leaf,
			})
		}
	}
	r.rvs = append(rcounts, r.rvs...)
	r.rvs = bindRVarsTo(t, r.rvs)

	r.lvs = make([]rleaf, 0, len(r.rvs))
	for i := range r.rvs {
		rv := r.rvs[i]
		r.lvs = append(r.lvs, rleafFrom(rv.leaf, rv, r))
	}

	// regroup leaves by holding branch
	set := make(map[string]int)
	brs := make([][]rleaf, 0, len(r.lvs))
	for _, leaf := range r.lvs {
		br := leaf.Leaf().Branch().Name()
		if _, ok := set[br]; !ok {
			set[br] = len(brs)
			brs = append(brs, []rleaf{})
		}
		id := set[br]
		brs[id] = append(brs[id], leaf)
	}

	r.brs = make([]rbranch, len(brs))
	for i, leaves := range brs {
		branch := leaves[0].Leaf().Branch()
		r.brs[i] = newRBranch(branch, n, beg, end, leaves, r)
	}

	return r
}
func (r *rtree) Close() error {
	for i := range r.brs {
		rb := &r.brs[i]
		rb.rb.close()
	}
	return nil
}

func (r *rtree) reset() {
	for i := range r.brs {
		rb := &r.brs[i]
		rb.reset()
	}
}

func (r *rtree) rcountFunc(name string) func() int {
	for _, leaf := range r.lvs {
		n := leaf.Leaf().Name()
		if n != name {
			continue
		}
		switch leaf := leaf.(type) {
		case *rleafValI8:
			return leaf.ivalue
		case *rleafValI16:
			return leaf.ivalue
		case *rleafValI32:
			return leaf.ivalue
		case *rleafValI64:
			return leaf.ivalue
		case *rleafValU8:
			return leaf.ivalue
		case *rleafValU16:
			return leaf.ivalue
		case *rleafValU32:
			return leaf.ivalue
		case *rleafValU64:
			return leaf.ivalue
		case *rleafElem:
			leaf.bindCount()
			return leaf.ivalue

		default:
			panic(fmt.Errorf("rleaf %T not implemented", leaf))
		}
	}
	panic(fmt.Errorf("impossible: no leaf for %s", name))
}

func (r *rtree) rcountLeaf(name string) leafCount {
	for _, leaf := range r.lvs {
		n := leaf.Leaf().Name()
		if n != name {
			continue
		}
		return &rleafCount{
			Leaf: leaf.Leaf(),
			n:    r.rcountFunc(name),
			leaf: leaf,
		}
	}
	panic(fmt.Errorf("impossible: no leaf for %s", name))
}

func (r *rtree) run(off, beg, end int64, f func(RCtx) error) error {
	var (
		err  error
		rctx RCtx
	)

	defer r.Close()

	for i := range r.brs {
		rb := &r.brs[i]
		err = rb.start()
		if err != nil {
			if err == io.EOF {
				// empty range.
				return nil
			}
			return err
		}
	}
	defer func() {
		for i := range r.brs {
			rb := &r.brs[i]
			_ = rb.stop()
		}
	}()

	for i := beg; i < end; i++ {
		err = r.read(i)
		if err != nil {
			return fmt.Errorf("rtree: could not read entry %d: %w", i, err)
		}
		rctx.Entry = i + off
		err = f(rctx)
		if err != nil {
			return fmt.Errorf("rtree: could not process entry %d: %w", i, err)
		}
	}

	return err
}

func (r *rtree) read(ievt int64) error {
	for i := range r.brs {
		rb := &r.brs[i]
		err := rb.read(ievt)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	_ rleafCtx = (*rtree)(nil)
)
