// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
)

// ReadVar describes a variable to be read out of a tree.
type ReadVar struct {
	Name  string      // name of the branch to read
	Leaf  string      // name of the leaf to read
	Value interface{} // pointer to the value to fill
	count string      // name of the leaf-count, if any
}

// NewReadVars returns the complete set of ReadVars to read all the data
// contained in the provided Tree.
func NewReadVars(t Tree) []ReadVar {
	var vars []ReadVar
	for _, b := range t.Branches() {
		for _, leaf := range b.Leaves() {
			ptr := newValue(leaf)
			cnt := ""
			if leaf.LeafCount() != nil {
				cnt = leaf.LeafCount().Name()
			}
			vars = append(vars, ReadVar{Name: b.Name(), Leaf: leaf.Name(), Value: ptr, count: cnt})
		}
	}

	return vars
}

// ReadVarsFromStruct returns a list of ReadVars bound to the exported fields
// of the provided pointer to a struct value.
//
// ReadVarsFromStruct panicks if the provided value is not a pointer to
// a struct value.
func ReadVarsFromStruct(ptr interface{}) []ReadVar {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("rtree: expect a pointer value, got %T", ptr))
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("rtree: expect a pointer to struct value, got %T", ptr))
	}

	var (
		rt     = rv.Type()
		rvars  = make([]ReadVar, 0, rt.NumField())
		reDims = regexp.MustCompile(`\w*?\[(\w*)\]+?`)
	)

	split := func(s string) (string, []string) {
		n := s
		if i := strings.Index(s, "["); i > 0 {
			n = s[:i]
		}

		out := reDims.FindAllStringSubmatch(s, -1)
		if len(out) == 0 {
			return n, nil
		}

		dims := make([]string, len(out))
		for i := range out {
			dims[i] = out[i][1]
		}
		return n, dims
	}

	for i := 0; i < rt.NumField(); i++ {
		var (
			ft = rt.Field(i)
			fv = rv.Field(i)
		)
		if ft.Name != strings.Title(ft.Name) {
			// not exported. ignore.
			continue
		}
		rvar := ReadVar{
			Name:  ft.Tag.Get("groot"),
			Value: fv.Addr().Interface(),
		}
		if rvar.Name == "" {
			rvar.Name = ft.Name
		}

		if strings.Contains(rvar.Name, "[") {
			switch ft.Type.Kind() {
			case reflect.Slice:
				sli, dims := split(rvar.Name)
				if len(dims) > 1 {
					panic(fmt.Errorf("rtree: invalid number of slice-dimensions for field %q: %q", ft.Name, rvar.Name))
				}
				rvar.Name = sli
				rvar.count = dims[0]

			case reflect.Array:
				arr, dims := split(rvar.Name)
				if len(dims) > 3 {
					panic(fmt.Errorf("rtree: invalid number of array-dimension for field %q: %q", ft.Name, rvar.Name))
				}
				rvar.Name = arr
			default:
				panic(fmt.Errorf("rtree: invalid field type for %q, or invalid struct-tag %q: %T", ft.Name, rvar.Name, fv.Interface()))
			}
		}
		switch ft.Type.Kind() {
		case reflect.Int, reflect.Uint, reflect.UnsafePointer, reflect.Uintptr, reflect.Chan, reflect.Interface:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T", ft.Name, fv.Interface()))
		case reflect.Map:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T (not yet supported)", ft.Name, fv.Interface()))
		}

		rvar.Leaf = rvar.Name
		rvars = append(rvars, rvar)
	}
	return rvars
}

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
			})
		}
	}
	r.rvs = append(rcounts, r.rvs...)

	r.lvs = make([]rleaf, len(r.rvs))
	for i, rvar := range r.rvs {
		br := t.Branch(rvar.Name)
		if br == nil {
			continue
		}
		leaf := br.Leaf(rvar.Leaf)
		r.lvs[i] = rleafFrom(leaf, rvar, r)
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

func (r *rtree) rcount(name string) func() int {
	for _, leaf := range r.lvs {
		n := leaf.Leaf().Name()
		if n != name {
			continue
		}
		switch leaf := leaf.(type) {
		case *rleafValI8:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValI16:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValI32:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValI64:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValU8:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValU16:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValU32:
			return func() int {
				return int(*leaf.v)
			}
		case *rleafValU64:
			return func() int {
				return int(*leaf.v)
			}
		default:
			panic(fmt.Errorf("rleaf %T not implemented", leaf))
		}
	}
	panic("impossible")
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
			rb.stop()
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
