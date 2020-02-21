// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package root defines ROOT core interfaces.
package root // import "go-hep.org/x/hep/groot/root"

import (
	"go-hep.org/x/hep/groot/rvers"
)

const (
	Version = rvers.ROOT // ROOT version the groot library implements
)

// Object represents a ROOT object
type Object interface {
	// Class returns the ROOT class of this object
	Class() string
}

// UIDer is the interface for objects that can be referenced.
type UIDer interface {
	// UID returns the unique ID of this object
	UID() uint32
}

// Named represents a ROOT TNamed object
type Named interface {
	Object

	// Name returns the name of this ROOT object
	Name() string

	// Title returns the title of this ROOT object
	Title() string
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

// Array describes ROOT abstract array type.
type Array interface {
	Len() int // number of array elements
	Get(i int) interface{}
	Set(i int, v interface{})
}

// ObjString is a ROOT string that implements ROOT TObject.
type ObjString interface {
	Name() string
	String() string
}

// Merger is a ROOT object that can ingest data from another ROOT object.
type Merger interface {
	ROOTMerge(src Object) error
}
