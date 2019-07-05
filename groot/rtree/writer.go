// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/riofs"
)

// Writer is the interface that wraps the Write method for Trees.
type Writer interface {
	Tree

	// Write writes the event data to ROOT storage.
	Write() error

	// Sync commits the current contents of the tree to stable storage.
	Sync() error

	// Close writes metadata and closes the tree.
	Close() error
}

type wtree struct {
	ttree

	dir riofs.Directory // directory holding this tree
	typ reflect.Type
}

// WriteVar describes a variable to be written out to a tree.
type WriteVar struct {
	Name  string      // name of the variable
	Value interface{} // pointer to the value to write
}

// NewWriter creates a new Tree with the given name and under the given
// directory dir, ready to be filled with data.
func NewWriter(dir riofs.Directory, name string, vars []WriteVar) (Writer, error) {
	if dir == nil {
		return nil, errors.Errorf("rtree: missing parent directory")
	}

	w := &wtree{
		dir: dir,
		//typ: typ,
	}
	w.ttree.f = fileOf(dir)
	w.ttree.named.SetName(name)

	const compress = 1 // FIXME: make it func-opt
	for _, v := range vars {
		b, err := newBranchFromWVars(w, v.Name, []WriteVar{v}, nil, compress)
		if err != nil {
			return nil, errors.Wrapf(err, "rtree: could not create branch for write-var %#v", v)
		}
		w.ttree.branches = append(w.ttree.branches, b)
	}

	return w, nil
}

func (w *wtree) SetTitle(title string) { w.ttree.named.SetTitle(title) }

// Write writes the event data to ROOT storage.
func (w *wtree) Write() error {
	for _, b := range w.ttree.branches {
		err := b.write()
		if err != nil {
			return errors.Wrapf(err, "rtree: could not write branch %q", b.Name())
		}
	}
	w.ttree.entries++
	// FIXME(sbinet): autoflush
	return nil
}

// Sync commits the current contents of the tree to stable storage.
func (w *wtree) Sync() error {
	panic("not implemented") // FIXME(sbinet): implement
}

// Close writes metadata and closes the tree.
func (w *wtree) Close() error {
	if err := w.Sync(); err != nil {
		return err
	}
	panic("not implemented") // FIXME(sbinet): implement
}

func (w *wtree) loadEntry(i int64) error {
	return errors.Errorf("rtree: Tree writer can not be read from")
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
