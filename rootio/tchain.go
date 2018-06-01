// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type tchain struct {
	//	trees []Tree
	trees []itree
	cur   int
	tree  *itree
}

type itree struct {
	tree   Tree
	offset int64 // number of entries before the current one
	total  int64 // total number of entries iterated over
}

func (t itree) Entries() int64 {
	return t.tree.Entries()
}

func (t *itree) loadEntry(i int64) error {
	j := i - t.offset
	return t.tree.loadEntry(j)
}

type ibranch struct {
	br     Branch
	offset int64 // number of entries before the current one
	total  int64 // total number of entries iterated over
}

func (b *ibranch) loadEntry(i int64) error {
	j := i - b.offset
	return b.br.loadEntry(j)
}

func (b *ibranch) Branches() []Branch {
	return b.br.Branches()
}

func (b *ibranch) Leaves() []Leaf {
	return b.br.Leaves()
}

func (b *ibranch) Leaf(name string) Leaf {
	return b.br.Leaf(name)
}

func (b *ibranch) setTree(t Tree) {
	b.br.setTree(t)
}

func (b *ibranch) getTree() Tree {
	return b.br.getTree()
}

func (b *ibranch) getReadEntry() int64 {
	return b.br.getReadEntry()
}

func (b *ibranch) getEntry(i int64) {
	b.br.getEntry(i)
}

func (b *ibranch) Branch(name string) Branch {
	return b.br.Branch(name)
}

func (b *ibranch) scan(ptr interface{}) error {
	return b.br.scan(ptr)
}

func (b *ibranch) setAddress(ptr interface{}) error {
	return b.br.setAddress(ptr)
}

func (b *ibranch) setStreamer(s StreamerInfo, ctx StreamerInfoContext) {
	b.br.setStreamer(s, ctx)
}

func (b *ibranch) setStreamerElement(s StreamerElement, ctx StreamerInfoContext) {
	b.br.setStreamerElement(s, ctx)
}

func (b *ibranch) GoType() reflect.Type {
	return b.br.GoType()
}

// Class returns the ROOT class of the argument.
func (tchain) Class() string {
	return "TChain"
}

// Name returns the name of the ROOT objet in the argument.
func (t tchain) Name() string {
	if len(t.trees) == 0 {
		return ""
	}
	return t.trees[0].tree.Name()
}

// Title returns the title of the ROOT object in the argument.
func (t tchain) Title() string {
	if len(t.trees) == 0 {
		return ""
	}
	return t.trees[0].tree.Title()
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) *tchain {

	ch := &tchain{
		trees: make([]itree, len(trees)),
		cur:   -1,
	}
	var sum int64
	var offset int64
	for i := range trees {
		t := trees[i]
		n := t.Entries()
		sum += n
		ch.trees[i] = itree{tree: t, offset: offset, total: sum}
		offset += n
	}
	if len(trees) > 0 {
		ch.cur = 0
		ch.tree = &ch.trees[ch.cur]
	}
	return ch
}

// Entries returns the total number of entries.
func (t tchain) Entries() int64 {
	if len(t.trees) <= 0 {
		return 0
	}
	return t.trees[len(t.trees)-1].total

}

// TotBytes return the total number of bytes before compression.
func (t tchain) TotBytes() int64 {
	return 0
}

// ZipBytes returns the total number of bytes after compression.
func (t tchain) ZipBytes() int64 {
	return 0
}

// Branches returns the list of branches.
func (t tchain) Branches() []Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Branches()
}

// Branch returns the branch whose name is the argument.
func (t tchain) Branch(name string) Branch {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Branch(name)
}

// Leaves returns direct pointers to individual branch leaves.
func (t tchain) Leaves() []Leaf {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.Leaves()
}

// getFile returns the underlying file.
func (t tchain) getFile() *File {
	if len(t.trees) == 0 {
		return nil
	}
	return t.trees[0].tree.getFile()
}

// loadEntry returns an error if there is a problem during the loading.
func (ch tchain) loadEntry(i int64) error {
	if ch.tree == nil {
		return fmt.Errorf("invalid chain")
	}

	if i >= ch.tree.total {
		ch.cur++
		ch.tree = &ch.trees[ch.cur]
	}
	return ch.tree.loadEntry(i)
}

var (
	_ Object = (*tchain)(nil)
	_ Named  = (*tchain)(nil)
	_ Tree   = (*ttree)(nil)
	_ Tree   = (*tchain)(nil)
)
