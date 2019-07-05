// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rtree contains the interfaces and types to decode, read, concatenate
// and iterate over ROOT Trees.
package rtree // import "go-hep.org/x/hep/groot/rtree"

import (
	"reflect" // Tree is a collection of branches of data.

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

// FileOf returns the file hosting the given Tree.
// If the tree is not connected to any ROOT file, nil is returned.
func FileOf(tree Tree) *riofs.File { return tree.getFile() }

type Tree interface {
	root.Named

	Entries() int64
	TotBytes() int64
	ZipBytes() int64
	Branch(name string) Branch
	Branches() []Branch
	Leaf(name string) Leaf
	Leaves() []Leaf

	getFile() *riofs.File
	loadEntry(i int64) error
}

// Branch describes a branch of a ROOT Tree.
type Branch interface {
	root.Named

	Branches() []Branch
	Leaves() []Leaf
	Branch(name string) Branch
	Leaf(name string) Leaf

	setTree(Tree)
	getTree() Tree
	loadEntry(i int64) error
	getReadEntry() int64
	getEntry(i int64)
	scan(ptr interface{}) error
	setAddress(ptr interface{}) error
	setStreamer(s rbytes.StreamerInfo, ctx rbytes.StreamerInfoContext)
	setStreamerElement(s rbytes.StreamerElement, ctx rbytes.StreamerInfoContext)
	GoType() reflect.Type

	// write interface part
	writeToBuffer(w *rbytes.WBuffer) error
	write() error
}

// Leaf describes branches data types
type Leaf interface {
	root.Named

	ArrayDim() int
	Branch() Branch
	HasRange() bool
	IsUnsigned() bool
	LeafCount() Leaf // returns the leaf count if is variable length
	Len() int        // Len returns the number of fixed length elements
	LenType() int    // LenType returns the number of bytes for this data type
	MaxIndex() []int
	Offset() int
	Kind() reflect.Kind
	Type() reflect.Type
	Value(int) interface{}
	TypeName() string

	setBranch(Branch)
	readFromBuffer(r *rbytes.RBuffer) error
	value() interface{}
	scan(r *rbytes.RBuffer, ptr interface{}) error
	setAddress(ptr interface{}) error

	// write interface part
	writeToBuffer(w *rbytes.WBuffer) error

	canGenerateOffsetArray() bool
	computeOffsetArray(base, nevts int) []int32
}

// leafCount describes leaves that are used for array length count
type leafCount interface {
	Leaf
	ivalue() int // for leaf-count
	imax() int
}
