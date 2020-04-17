// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package groot

//go:generate go run ./gen.rboot.go
//go:generate go run ./gen.rcont.go
//go:generate go run ./gen.rhist.go
//go:generate go run ./gen.rtree.go

import (
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	_ "go-hep.org/x/hep/groot/ztypes"
)

const (
	Version = root.Version // ROOT version hep/groot implements
)

// Open opens the named ROOT file for reading. If successful, methods on the
// returned file can be used for reading; the associated file descriptor
// has mode os.O_RDONLY.
func Open(path string) (*File, error) {
	return riofs.Open(path)
}

// NewReader creates a new ROOT file reader.
func NewReader(r Reader) (*File, error) {
	return riofs.NewReader(r)
}

// Create creates the named ROOT file for writing.
func Create(name string, opts ...FileOption) (*File, error) {
	return riofs.Create(name, opts...)
}

type (
	File       = riofs.File
	FileOption = riofs.FileOption
	Reader     = riofs.Reader
)

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
