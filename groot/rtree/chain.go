// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

type chain struct {
	trees []Tree
	offs  []int64 // number of entries before this tree
	tots  []int64 // total number of entries up to this tree

	cur  int   // index of current tree
	tree Tree  // current tree
	off  int64 // current offset
	tot  int64 // current number of entries
}

// Chain returns a Tree that is the concatenation of all the input Trees.
func Chain(trees ...Tree) Tree {
	if len(trees) == 0 {
		return &chain{}
	}
	n := len(trees)
	ch := &chain{
		trees: make([]Tree, n),
		offs:  make([]int64, n),
		tots:  make([]int64, n),
		cur:   -1,
	}
	var (
		sum int64
		off int64
	)
	for i := range trees {
		t := trees[i]
		n := t.Entries()
		sum += n
		ch.trees[i] = t
		ch.offs[i] = off
		ch.tots[i] = sum
		off += n
	}

	ch.loadTree(ch.cur + 1)
	return ch
}

// ChainOf returns a Tree, a close function and an error if any.
// The tree is the logical concatenation of all the name trees
// located in the input named files.
// The close function allows to close all the open named files.
func ChainOf(name string, files ...string) (Tree, func() error, error) {
	var (
		trees = make([]Tree, len(files))
		fs    = make([]*riofs.File, len(files))
	)

	closef := func(fs []*riofs.File) {
		for _, f := range fs {
			if f == nil {
				continue
			}
			f.Close()
		}
	}

	for i, n := range files {
		f, err := riofs.Open(n)
		if err != nil {
			closef(fs)
			return nil, nil, err
		}
		fs[i] = f
		obj, err := f.Get(name)
		if err != nil {
			closef(fs)
			return nil, nil, err
		}
		t, ok := obj.(Tree)
		if !ok {
			closef(fs)
			return nil, nil, fmt.Errorf("rtree: object %q in file %q is not a Tree", name, n)
		}

		trees[i] = t
	}

	ch := Chain(trees...)
	close := func() error {
		var err error
		for _, f := range fs {
			e := f.Close()
			if e != nil && err == nil {
				err = e
			}
		}
		return err
	}

	return ch, close, nil
}

func (ch *chain) loadTree(i int) {
	ch.cur = i
	if ch.cur >= len(ch.trees) {
		ch.tree = nil
		return
	}
	ch.tree = ch.trees[ch.cur]
	ch.off = ch.offs[ch.cur]
	ch.tot = ch.tots[ch.cur]
}

// Class returns the ROOT class of the argument.
func (*chain) Class() string {
	return "TChain"
}

// Name returns the name of the ROOT objet in the argument.
func (t *chain) Name() string {
	if t.tree == nil {
		return ""
	}
	return t.tree.Name()
}

// Title returns the title of the ROOT object in the argument.
func (t *chain) Title() string {
	if t.tree == nil {
		return ""
	}
	return t.tree.Title()
}

// Entries returns the total number of entries.
func (t *chain) Entries() int64 {
	var v int64
	for _, tree := range t.trees {
		v += tree.Entries()
	}
	return v
}

// Branches returns the list of branches.
func (t *chain) Branches() []Branch {
	if t.tree == nil {
		return nil
	}
	return t.tree.Branches()
}

// Branch returns the branch whose name is the argument.
func (t *chain) Branch(name string) Branch {
	if t.tree == nil {
		return nil
	}
	return t.tree.Branch(name)
}

// Leaves returns direct pointers to individual branch leaves.
func (t *chain) Leaves() []Leaf {
	if t.tree == nil {
		return nil
	}
	return t.tree.Leaves()
}

// Leaf returns the leaf whose name is the argument.
func (t *chain) Leaf(name string) Leaf {
	if t.tree == nil {
		return nil
	}
	return t.tree.Leaf(name)
}

var (
	_ root.Object = (*chain)(nil)
	_ root.Named  = (*chain)(nil)
	_ Tree        = (*chain)(nil)
)
