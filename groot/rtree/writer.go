// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
)

// Writer is the interface that wraps the Write method for Trees.
type Writer interface {
	Tree

	// Write writes the event data to ROOT storage and returns the number
	// of bytes (before compression, if any) written.
	Write() (int, error)

	// Flush commits the current contents of the tree to stable storage.
	Flush() error

	// Close writes metadata and closes the tree.
	Close() error
}

// WriteOption configures how a ROOT tree (and its branches) should be created.
type WriteOption func(opt *wopt) error

type wopt struct {
	title    string // title of the writer tree
	bufsize  int32  // buffer size for branches
	splitlvl int32  // maximum split-level for branches
	compress int32  // compression algorithm name and compression level
}

// WithLZ4 configures a ROOT tree to use LZ4 as a compression mechanism.
func WithLZ4(level int) WriteOption {
	return func(opt *wopt) error {
		opt.compress = rcompress.Settings{Alg: rcompress.LZ4, Lvl: level}.Compression()
		return nil
	}
}

// WithLZMA configures a ROOT tree to use LZMA as a compression mechanism.
func WithLZMA(level int) WriteOption {
	return func(opt *wopt) error {
		opt.compress = rcompress.Settings{Alg: rcompress.LZMA, Lvl: level}.Compression()
		return nil
	}
}

// WithoutCompression configures a ROOT tree to not use any compression mechanism.
func WithoutCompression() WriteOption {
	return func(opt *wopt) error {
		opt.compress = 0
		return nil
	}
}

// WithZlib configures a ROOT tree to use zlib as a compression mechanism.
func WithZlib(level int) WriteOption {
	return func(opt *wopt) error {
		opt.compress = rcompress.Settings{Alg: rcompress.ZLIB, Lvl: level}.Compression()
		return nil
	}
}

// WithBasketSize configures a ROOT tree to use 'size' (in bytes) as a basket buffer size.
// if size is <= 0, the default buffer size is used (DefaultBasketSize).
func WithBasketSize(size int) WriteOption {
	return func(opt *wopt) error {
		if size <= 0 {
			size = defaultBasketSize
		}
		opt.bufsize = int32(size)
		return nil
	}
}

// WithTitle sets the title of the tree writer.
func WithTitle(title string) WriteOption {
	return func(opt *wopt) error {
		opt.title = title
		return nil
	}
}

type wtree struct {
	ttree
	wvars []WriteVar
}

// WriteVar describes a variable to be written out to a tree.
type WriteVar struct {
	Name  string      // name of the variable
	Value interface{} // pointer to the value to write
	Count string      // name of the branch holding the count-leaf value for slices
}

// WriteVarsFromStruct creates a slice of WriteVars from the ptr value.
// WriteVarsFromStruct panics if ptr is not a pointer to a struct value.
// WriteVarsFromStruct ignores fields that are not exported.
func WriteVarsFromStruct(ptr interface{}) []WriteVar {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic(xerrors.Errorf("rtree: expect a pointer value, got %T", ptr))
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		panic(xerrors.Errorf("rtree: expect a pointer to struct value, got %T", ptr))
	}
	var (
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

	var (
		rt    = rv.Type()
		wvars = make([]WriteVar, 0, rt.NumField())
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			ft = rt.Field(i)
			fv = rv.Field(i)
		)
		if ft.Name != strings.Title(ft.Name) {
			// not exported. ignore.
			continue
		}
		wvar := WriteVar{
			Name:  ft.Tag.Get("groot"),
			Value: fv.Addr().Interface(),
		}
		if wvar.Name == "" {
			wvar.Name = ft.Name
		}

		if strings.Contains(wvar.Name, "[") {
			switch ft.Type.Kind() {
			case reflect.Slice:
				sli, dims := split(wvar.Name)
				if len(dims) > 1 {
					panic(xerrors.Errorf("rtree: invalid number of slice-dimensions for field %q: %q", ft.Name, wvar.Name))
				}
				wvar.Name = sli
				wvar.Count = dims[0]

			case reflect.Array:
				arr, dims := split(wvar.Name)
				if len(dims) > 3 {
					panic(xerrors.Errorf("rtree: invalid number of array-dimension for field %q: %q", ft.Name, wvar.Name))
				}
				wvar.Name = arr
			default:
				panic(xerrors.Errorf("rtree: invalid field type for %q, or invalid struct-tag %q: %T", ft.Name, wvar.Name, fv.Interface()))
			}
		}
		switch ft.Type.Kind() {
		case reflect.Int, reflect.Uint, reflect.UnsafePointer, reflect.Uintptr, reflect.Chan, reflect.Interface:
			panic(xerrors.Errorf("rtree: invalid field type for %q: %T", ft.Name, fv.Interface()))
		case reflect.Map:
			panic(xerrors.Errorf("rtree: invalid field type for %q: %T (not yet supported)", ft.Name, fv.Interface()))
		}

		wvars = append(wvars, wvar)
	}

	return wvars
}

// WriteVarsFromTree creates a slice of WriteVars from the tree value.
func WriteVarsFromTree(t Tree) []WriteVar {
	rvars := NewScanVars(t)
	wvars := make([]WriteVar, len(rvars))
	for i, rvar := range rvars {
		wvars[i] = WriteVar{
			Name:  rvar.Name,
			Value: reflect.New(reflect.TypeOf(rvar.Value).Elem()).Interface(),
			Count: rvar.count,
		}
	}
	return wvars
}

// NewWriter creates a new Tree with the given name and under the given
// directory dir, ready to be filled with data.
func NewWriter(dir riofs.Directory, name string, vars []WriteVar, opts ...WriteOption) (Writer, error) {
	if dir == nil {
		return nil, xerrors.Errorf("rtree: missing parent directory")
	}

	w := &wtree{
		ttree: ttree{
			f:     fileOf(dir),
			dir:   dir,
			rvers: rvers.Tree,
			named: *rbase.NewNamed(name, ""),
		},
		wvars: vars,
	}

	cfg := wopt{
		bufsize:  defaultBasketSize,
		splitlvl: defaultSplitLevel,
		compress: w.ttree.f.Compression(),
	}

	for _, opt := range opts {
		err := opt(&cfg)
		if err != nil {
			return nil, xerrors.Errorf("rtree: could not configure tree writer: %w", err)
		}
	}

	w.ttree.named.SetTitle(cfg.title)

	for _, v := range vars {
		b, err := newBranchFromWVars(w, v.Name, []WriteVar{v}, nil, cfg)
		if err != nil {
			return nil, xerrors.Errorf("rtree: could not create branch for write-var %#v: %w", v, err)
		}
		w.ttree.branches = append(w.ttree.branches, b)
	}

	return w, nil
}

func (w *wtree) SetTitle(title string) { w.ttree.named.SetTitle(title) }

// Write writes the event data to ROOT storage and returns the number
// of bytes (before compression, if any) written.
func (w *wtree) Write() (int, error) {
	var tot int
	for _, b := range w.ttree.branches {
		nbytes, err := b.write()
		if err != nil {
			return tot, xerrors.Errorf("rtree: could not write branch %q: %w", b.Name(), err)
		}
		tot += nbytes
	}
	w.ttree.entries++
	w.ttree.totBytes += int64(tot)
	// FIXME(sbinet): autoflush

	return tot, nil
}

// Flush commits the current contents of the tree to stable storage.
func (w *wtree) Flush() error {
	for _, b := range w.ttree.branches {
		err := b.flush()
		if err != nil {
			return xerrors.Errorf("rtree: could not flush branch %q: %w", b.Name(), err)
		}
	}
	return nil
}

// Close writes metadata and closes the tree.
func (w *wtree) Close() error {
	if err := w.Flush(); err != nil {
		return xerrors.Errorf("rtree: could not flush tree %q: %w", w.Name(), err)
	}

	if err := w.ttree.dir.Put(w.Name(), w); err != nil {
		return xerrors.Errorf("rtree: could not save tree %q: %w", w.Name(), err)
	}

	return nil
}

func (w *wtree) loadEntry(i int64) error {
	return xerrors.Errorf("rtree: Tree writer can not be read from")
}

func fileOf(d riofs.Directory) *riofs.File {
	const max = 1<<31 - 1
	for i := 0; i < max; i++ {
		p := d.Parent()
		if p == nil {
			return d.(*riofs.File)
		}
		d = p
	}
	panic("impossible")
}

func flattenArrayType(rt reflect.Type) (reflect.Type, []int) {
	var (
		shape []int
		kind  = rt.Kind()
	)
	const max = 1<<31 - 1
	for i := 0; i < max; i++ {
		if kind != reflect.Array {
			return rt, shape
		}
		shape = append(shape, rt.Len())
		rt = rt.Elem()
		kind = rt.Kind()
	}
	panic("impossible")
}

var (
	_ Tree   = (*wtree)(nil)
	_ Writer = (*wtree)(nil)
)
