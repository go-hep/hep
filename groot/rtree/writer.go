// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
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
	bufsize  int32 // buffer size for branches
	splitlvl int32 // maximum split-level for branches
	compress int   // compression algorithm name and compression level
}

// WithLZ4 configures a ROOT tree to use LZ4 as a compression mechanism.
func WithLZ4(level int) WriteOption {
	return func(opt *wopt) error {
		opt.compress = 100*int(rcompress.LZ4) + level
		return nil
	}
}

// WithLZMA configures a ROOT tree to use LZMA as a compression mechanism.
func WithLZMA(level int) WriteOption {
	return func(opt *wopt) error {
		opt.compress = 100*int(rcompress.LZMA) + level
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
		opt.compress = 100*int(rcompress.ZLIB) + level
		return nil
	}
}

type wtree struct {
	ttree
}

// WriteVar describes a variable to be written out to a tree.
type WriteVar struct {
	Name  string      // name of the variable
	Value interface{} // pointer to the value to write
	Count string      // name of the branch holding the count-leaf value for slices
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
		//typ: typ,
	}

	cfg := wopt{
		bufsize:  defaultBasketSize,
		splitlvl: defaultSplitLevel,
		compress: int(w.ttree.f.Compression()),
	}

	for _, opt := range opts {
		err := opt(&cfg)
		if err != nil {
			return nil, xerrors.Errorf("rtree: could not configure tree writer: %w", err)
		}
	}

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
	const max = 1 << 32
	for i := 0; i < max; i++ {
		p := d.Parent()
		if p == nil {
			return d.(*riofs.File)
		}
		d = p
	}
	panic("impossible")
}

var (
	_ Tree   = (*wtree)(nil)
	_ Writer = (*wtree)(nil)
)
