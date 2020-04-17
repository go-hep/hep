// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slha implements encoding and decoding of SUSY Les Houches Accords (SLHA) data format.
package slha // import "go-hep.org/x/hep/slha"

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SLHA holds informations about a SUSY Les Houches Accords file.
type SLHA struct {
	Blocks    Blocks
	Particles Particles
}

// Value represents a value (string,int,float64) + comment in a SLHA line.
type Value struct {
	v reflect.Value
	c string // comment attached to value
}

// Int returns the value as an int64.
// Int panics if the underlying value isn't an int64.
func (v *Value) Int() int64 {
	return v.v.Int()
}

// Float returns the value as a float64.
// Float panics if the underlying value isn't a float64.
func (v *Value) Float() float64 {
	return v.v.Float()
}

// Interface returns the value as an interface{}.
func (v *Value) Interface() interface{} {
	return v.v.Interface()
}

// Kind returns the kind of the value (reflect.String, reflect.Float64, reflect.Int)
func (v *Value) Kind() reflect.Kind {
	return v.v.Kind()
}

// Comment returns the comment string attached to this value
func (v *Value) Comment() string {
	return v.c
}

// Block represents a block in a SLHA file.
type Block struct {
	Name    string
	Comment string
	Q       float64
	Data    DataArray
}

// Get returns the Value at index args.
// Note that args are 1-based indices.
func (b *Block) Get(args ...int) (Value, error) {
	var err error
	var val Value

	idx := NewIndex(args...)
	val, ok := b.Data.Get(idx)
	if !ok {
		return val, fmt.Errorf("slha: no index (%s) in block %q", strings.Join(strindex(args...), ", "), b.Name)
	}
	return val, err
}

// Set sets the Value at index args with v.
// Set creates a new empty Value if none exists at args.
// Note that args are 1-based indices.
func (b *Block) Set(v interface{}, args ...int) error {
	var err error
	val, _ := b.Get(args...)
	val.v = reflect.ValueOf(v)
	idx := NewIndex(args...)
	pos := b.Data.pos(idx)
	if pos < 0 {
		pos = len(b.Data)
		b.Data = append(b.Data, DataItem{
			Index: idx,
		})
	}
	b.Data[pos].Value = val
	return err
}

// DataArray is an ordered list of DataItems.
type DataArray []DataItem

// DataItem is a pair of (Index,Value).
// Index is a n-dim index (1-based indices)
type DataItem struct {
	Index Index
	Value Value
}

// Get returns the value at the n-dim index idx.
func (d DataArray) Get(idx Index) (Value, bool) {
	var val Value
	for _, v := range d {
		if v.Index == idx {
			return v.Value, true
		}
	}
	return val, false
}

func (d DataArray) pos(idx Index) int {
	for i, v := range d {
		if v.Index == idx {
			return i
		}
	}
	return -1
}

// Index is an n-dimensional index.
// Note that the indices are 1-based.
type Index struct {
	rank   int
	coords string
}

// NewIndex creates a new n-dim index from args.
// Note that args are 1-based indices.
func NewIndex(args ...int) Index {
	sargs := strindex(args...)
	return Index{
		rank:   len(args),
		coords: strings.Join(sargs, "#"),
	}
}

// Index returns the n-dim indices.
// Note that the indices are 1-based.
func (idx Index) Index() []int {
	sargs := strings.Split(idx.coords, "#")
	args := make([]int, len(sargs))
	for i, v := range sargs {
		if v == "" {
			continue
		}
		var err error
		args[i], err = strconv.Atoi(v)
		if err != nil {
			panic(fmt.Errorf("slha.index: %v", err))
		}
	}
	return args
}

func strindex(args ...int) []string {
	sargs := make([]string, len(args))
	for i, v := range args {
		sargs[i] = strconv.Itoa(v)
	}
	return sargs
}

// Blocks is a list of Blocks.
type Blocks []Block

// Keys returns the names of the contained blocks.
func (b Blocks) Keys() []string {
	keys := make([]string, len(b))
	for i := range b {
		keys[i] = b[i].Name
	}
	return keys
}

// Get returns the block named name or nil.
func (b Blocks) Get(name string) *Block {
	for i := range b {
		blk := &b[i]
		if blk.Name == name {
			return blk
		}
	}
	return nil
}

// Decay is a decay line in an SLHA file.
type Decay struct {
	Br      float64 // Branching Ratio
	IDs     []int   // list of PDG IDs to which the decay occur
	Comment string  // comment attached to this decay line - if any
}

// Decays is a list of decays in a Decay block.
type Decays []Decay

// Particle is the representation of a single, specific particle, decay block from a SLHA file.
type Particle struct {
	PdgID   int     // PDG-ID code
	Width   float64 // total width of that particle
	Mass    float64 // mass of that particle
	Comment string
	Decays  Decays
}

// Particles is a block of particle's decays in an SLHA file.
type Particles []Particle

func (p Particles) Len() int           { return len(p) }
func (p Particles) Less(i, j int) bool { return p[i].PdgID < p[j].PdgID }
func (p Particles) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Get returns the Particle with matching pdgid or nil.
func (p Particles) Get(pdgid int) *Particle {
	for i := range p {
		part := &p[i]
		if part.PdgID == pdgid {
			return part
		}
	}
	return nil
}
