// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
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
	t     Tree
	rvars []ReadVar
	scan  *Scanner
	beg   int64
	end   int64

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

	if r.beg > r.t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (start=%d > tree-entries=%d)",
			r.beg, r.end, r.beg, r.t.Entries(),
		)
	}

	if r.end > r.t.Entries() {
		return nil, fmt.Errorf("rtree: invalid event reader range [%d, %d) (end=%d > tree-entries=%d)",
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
		_ = r.scan.Close()
		sc, err := NewScannerVars(r.t, r.rvars...)
		if err != nil {
			return fmt.Errorf("rtree: could not re-create scanner: %w", err)
		}
		r.scan = sc
	}

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
