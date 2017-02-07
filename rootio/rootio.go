// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootio provides a pure-go read-access to ROOT files.
// rootio might, with time, provide write-access too.
//
// A typical usage is as follow:
//
//   f, err := rootio.Open("ntup.root")
//   obj, err := f.Get("tree")
//   tree := obj.(*rootio.Tree)
//   fmt.Printf("entries= %v\n", t.Entries())
package rootio // import "go-hep.org/x/hep/rootio"

import (
	"bytes"
	"reflect"
)

//go:generate go run ./gen-code.go

// Class represents a ROOT class.
// Class instances are created by a ClassFactory.
type Class interface {
	// GetCheckSum gets the check sum for this ROOT class
	CheckSum() int

	// Members returns the list of members for this ROOT class
	Members() []Member

	// Version returns the version number for this ROOT class
	Version() int

	// ClassName returns the ROOT class name for this ROOT class
	ClassName() string
}

// Member represents a single member of a ROOT class
type Member interface {
	// GetArrayDim returns the dimension of the array (if any)
	ArrayDim() int

	// GetComment returns the comment associated with this member
	Comment() string

	// Name returns the name of this member
	Name() string

	// Type returns the class of this member
	Type() Class

	// GetValue returns the value of this member
	Value(o Object) reflect.Value
}

// Object represents a ROOT object
type Object interface {
	// Class returns the ROOT class of this object
	Class() string
}

// Named represents a ROOT TNamed object
type Named interface {
	Object

	// Name returns the name of this ROOT object
	Name() string

	// Title returns the title of this ROOT object
	Title() string
}

// ClassFactory creates ROOT classes
type ClassFactory interface {
	Create(name string) Class
}

// Collection is a collection of ROOT Objects.
type Collection interface {
	Object

	// Name returns the name of the collection.
	Name() string

	// Last returns the last element index
	Last() int

	// At returns the element at index i
	At(i int) Object

	// Len returns the number of elements in the collection
	Len() int
}

// SeqCollection is a sequential collection of ROOT Objects.
type SeqCollection interface {
	Collection
}

// List is a list of ROOT Objects.
type List interface {
	SeqCollection
}

// ObjArray is an array of ROOT Objects.
type ObjArray interface {
	SeqCollection
	LowerBound() int
}

// Directory describes a ROOT directory structure in memory.
type Directory interface {
	// Get returns the object identified by namecycle
	//   namecycle has the format name;cycle
	//   name  = * is illegal, cycle = * is illegal
	//   cycle = "" or cycle = 9999 ==> apply to a memory object
	//
	//   examples:
	//     foo   : get object named foo in memory
	//             if object is not in memory, try with highest cycle from file
	//     foo;1 : get cycle 1 of foo on file
	Get(namecycle string) (Object, bool)
}

// StreamerInfo describes a ROOT Streamer.
type StreamerInfo interface {
	Named
	CheckSum() int
	ClassVersion() int
	Elements() []StreamerElement
}

// StreamerElement describes a ROOT StreamerElement
type StreamerElement interface {
	Named
	ArrayDim() int
	ArrayLen() int
	Type() int
	Offset() uintptr
	Size() uintptr
	TypeName() string
}

// SetFiler is a simple interface to establish File ownership.
type SetFiler interface {
	SetFile(f *File)
}

// Tree is a collection of branches of data.
type Tree interface {
	Named
	Entries() int64
	TotBytes() int64
	ZipBytes() int64
	Branches() []Branch
	Leaves() []Leaf
}

// Branch describes a branch of a ROOT Tree.
type Branch interface {
	Named
	SetTree(Tree)
	Branches() []Branch
	Leaves() []Leaf
}

// Leaf describes branches data types
type Leaf interface {
	Named
	ArrayDim() int
	SetBranch(Branch)
	Branch() Branch
	HasRange() bool
	IsUnsigned() bool
	LeafCount() Leaf // returns the leaf count if is variable length
	Len() int        // Len returns the number of fixed length elements
	LenType() int    // LenType returns the number of bytes for this data type
	MaxIndex() []int
	Offset() int
	Value(int) interface{}
}

// Array describes ROOT abstract array type.
type Array interface {
	Len() int // number of array elements
	Get(i int) interface{}
	Set(i int, v interface{})
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
type ROOTUnmarshaler interface {
	UnmarshalROOT(r *RBuffer) error
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself into a ROOT buffer
type ROOTMarshaler interface {
	MarshalROOT() (data *bytes.Buffer, err error)
}
