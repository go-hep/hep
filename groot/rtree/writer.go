// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rvers"
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

	closed bool
}

// NewWriter creates a new Tree with the given name and under the given
// directory dir, ready to be filled with data.
func NewWriter(dir riofs.Directory, name string, vars []WriteVar, opts ...WriteOption) (Writer, error) {
	if dir == nil {
		return nil, fmt.Errorf("rtree: missing parent directory")
	}

	w := &wtree{
		ttree: ttree{
			f:         fileOf(dir),
			dir:       dir,
			rvers:     rvers.Tree,
			named:     *rbase.NewNamed(name, ""),
			attline:   *rbase.NewAttLine(),
			attfill:   *rbase.NewAttFill(),
			attmarker: *rbase.NewAttMarker(),
			weight:    1,
			scanField: 25,

			defaultEntryOffsetLen: 1000,
			maxEntries:            1000000000000,
			maxEntryLoop:          1000000000000,
			autoSave:              -300000000,
			autoFlush:             -30000000,
			estimate:              1000000,
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
			return nil, fmt.Errorf("rtree: could not configure tree writer: %w", err)
		}
	}

	w.ttree.named.SetTitle(cfg.title)

	for _, v := range vars {
		b, err := newBranchFromWVar(w, v.Name, v, nil, 0, cfg)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not create branch for write-var %#v: %w", v, err)
		}
		w.ttree.branches = append(w.ttree.branches, b)
	}

	return w, nil
}

func (w *wtree) SetTitle(title string) { w.ttree.named.SetTitle(title) }

func (w *wtree) ROOTMerge(src root.Object) error {
	switch src := src.(type) {
	case Tree:
		r, err := NewReader(src, nil)
		if err != nil {
			return fmt.Errorf("rtree: could not create tree reader: %w", err)
		}
		defer r.Close()

		_, err = Copy(w, r)
		if err != nil {
			return fmt.Errorf("rtree: could not merge tree: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("rtree: can not merge src=%T into dst=%T", src, w)
	}
}

// Write writes the event data to ROOT storage and returns the number
// of bytes (before compression, if any) written.
func (w *wtree) Write() (int, error) {
	var (
		tot int
		zip int
	)
	for _, b := range w.ttree.branches {
		nbytes, err := b.write()
		if err != nil {
			return tot, fmt.Errorf("rtree: could not write branch %q: %w", b.Name(), err)
		}
		tot += nbytes
	}
	w.ttree.entries++
	w.ttree.totBytes += int64(tot)
	w.ttree.zipBytes += int64(zip)
	// FIXME(sbinet): autoflush

	return tot, nil
}

// Flush commits the current contents of the tree to stable storage.
func (w *wtree) Flush() error {
	for _, b := range w.ttree.branches {
		err := b.flush()
		if err != nil {
			return fmt.Errorf("rtree: could not flush branch %q: %w", b.Name(), err)
		}
	}
	return nil
}

// Close writes metadata and closes the tree.
func (w *wtree) Close() error {
	if w.closed {
		return nil
	}
	defer func() {
		w.closed = true
	}()

	if err := w.Flush(); err != nil {
		return fmt.Errorf("rtree: could not flush tree %q: %w", w.Name(), err)
	}

	if err := w.ttree.dir.Put(w.Name(), w); err != nil {
		return fmt.Errorf("rtree: could not save tree %q: %w", w.Name(), err)
	}

	return nil
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
	_ Tree        = (*wtree)(nil)
	_ Writer      = (*wtree)(nil)
	_ root.Merger = (*wtree)(nil)
)
