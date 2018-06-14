// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

const (
	rootVersion = 61206 // ROOT version the Go-HEP/rootio library implements
)

//go:generate go run ./gen-code.go
//go:generate go run ./gendata/gen-evnt-tree.go -f ./testdata/small-evnt-tree-nosplit.root
//go:generate go run ./gendata/gen-evnt-tree.go -f ./testdata/small-evnt-tree-fullsplit.root -split=99

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
	Get(namecycle string) (Object, error)
	Keys() []Key
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
	Branch(name string) Branch
	Branches() []Branch
	Leaf(name string) Leaf
	Leaves() []Leaf

	getFile() *File
	loadEntry(i int64) error
}

// Branch describes a branch of a ROOT Tree.
type Branch interface {
	Named
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
	setStreamer(s StreamerInfo, ctx StreamerInfoContext)
	setStreamerElement(s StreamerElement, ctx StreamerInfoContext)
	GoType() reflect.Type
}

// Leaf describes branches data types
type Leaf interface {
	Named
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
	readBasket(r *RBuffer) error
	value() interface{}
	scan(r *RBuffer, ptr interface{}) error
}

// leafCount describes leaves that are used for array length count
type leafCount interface {
	Leaf
	ivalue() int // for leaf-count
	imax() int
}

// Array describes ROOT abstract array type.
type Array interface {
	Len() int // number of array elements
	Get(i int) interface{}
	Set(i int, v interface{})
}

// Axis describes a ROOT TAxis.
type Axis interface {
	Named
	XMin() float64
	XMax() float64
	NBins() int
	XBins() []float64
	BinCenter(int) float64
	BinLowEdge(int) float64
	BinWidth(int) float64
}

// H1 is a 1-dim ROOT histogram
type H1 interface {
	Named

	isH1()

	// Entries returns the number of entries for this histogram.
	Entries() float64
	// SumW returns the total sum of weights
	SumW() float64
	// SumW2 returns the total sum of squares of weights
	SumW2() float64
	// SumWX returns the total sum of weights*x
	SumWX() float64
	// SumWX2 returns the total sum of weights*x*x
	SumWX2() float64
	// SumW2s returns the array of sum of squares of weights
	SumW2s() []float64
}

// H2 is a 2-dim ROOT histogram
type H2 interface {
	Named

	isH2()

	// Entries returns the number of entries for this histogram.
	Entries() float64
	// SumW returns the total sum of weights
	SumW() float64
	// SumW2 returns the total sum of squares of weights
	SumW2() float64
	// SumWX returns the total sum of weights*x
	SumWX() float64
	// SumWX2 returns the total sum of weights*x*x
	SumWX2() float64
	// SumW2s returns the array of sum of squares of weights
	SumW2s() []float64
	// SumWY returns the total sum of weights*y
	SumWY() float64
	// SumWY2 returns the total sum of weights*y*y
	SumWY2() float64
	// SumWXY returns the total sum of weights*x*y
	SumWXY() float64
}

// Graph describes a ROOT TGraph
type Graph interface {
	Named
	Len() int
	XY(i int) (float64, float64)
}

// GraphErrors describes a ROOT TGraphErrors
type GraphErrors interface {
	Graph
	// XError returns two error values for X data.
	XError(i int) (float64, float64)
	// YError returns two error values for Y data.
	YError(i int) (float64, float64)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
type ROOTUnmarshaler interface {
	UnmarshalROOT(r *RBuffer) error
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself into a ROOT buffer
type ROOTMarshaler interface {
	MarshalROOT(w *WBuffer) (int, error)
}
