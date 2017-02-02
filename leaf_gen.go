// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
)

// LeafC implements ROOT TLeafC
type LeafC struct {
	tleaf
	min	int32
	max int32
}

// Class returns the ROOT class name.
func (leaf *LeafC) Class() string {
	return "TLeafC"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafC) Minimum() int32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafC) Maximum() int32 {
	return leaf.max
}

func (leaf *LeafC) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("LeafC: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckByteCount(pos, bcnt, start, "TLeafC")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafC{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafC", f)
	Factory.add("*rootio.LeafC", f)
}

var _ Object = (*LeafC)(nil)
var _ Named = (*LeafC)(nil)
var _ Leaf = (*LeafC)(nil)
var _ ROOTUnmarshaler = (*LeafC)(nil)
// LeafI implements ROOT TLeafI
type LeafI struct {
	tleaf
	min	int32
	max int32
}

// Class returns the ROOT class name.
func (leaf *LeafI) Class() string {
	return "TLeafI"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafI) Minimum() int32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafI) Maximum() int32 {
	return leaf.max
}

func (leaf *LeafI) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("LeafI: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckByteCount(pos, bcnt, start, "TLeafI")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafI{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafI", f)
	Factory.add("*rootio.LeafI", f)
}

var _ Object = (*LeafI)(nil)
var _ Named = (*LeafI)(nil)
var _ Leaf = (*LeafI)(nil)
var _ ROOTUnmarshaler = (*LeafI)(nil)
// LeafL implements ROOT TLeafL
type LeafL struct {
	tleaf
	min	int64
	max int64
}

// Class returns the ROOT class name.
func (leaf *LeafL) Class() string {
	return "TLeafL"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafL) Minimum() int64 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafL) Maximum() int64 {
	return leaf.max
}

func (leaf *LeafL) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("LeafL: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = r.ReadI64()
	leaf.max = r.ReadI64()

	r.CheckByteCount(pos, bcnt, start, "TLeafL")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafL{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafL", f)
	Factory.add("*rootio.LeafL", f)
}

var _ Object = (*LeafL)(nil)
var _ Named = (*LeafL)(nil)
var _ Leaf = (*LeafL)(nil)
var _ ROOTUnmarshaler = (*LeafL)(nil)
// LeafF implements ROOT TLeafF
type LeafF struct {
	tleaf
	min	float32
	max float32
}

// Class returns the ROOT class name.
func (leaf *LeafF) Class() string {
	return "TLeafF"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafF) Minimum() float32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafF) Maximum() float32 {
	return leaf.max
}

func (leaf *LeafF) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("LeafF: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = r.ReadF32()
	leaf.max = r.ReadF32()

	r.CheckByteCount(pos, bcnt, start, "TLeafF")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafF{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafF", f)
	Factory.add("*rootio.LeafF", f)
}

var _ Object = (*LeafF)(nil)
var _ Named = (*LeafF)(nil)
var _ Leaf = (*LeafF)(nil)
var _ ROOTUnmarshaler = (*LeafF)(nil)
// LeafD implements ROOT TLeafD
type LeafD struct {
	tleaf
	min	float64
	max float64
}

// Class returns the ROOT class name.
func (leaf *LeafD) Class() string {
	return "TLeafD"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafD) Minimum() float64 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafD) Maximum() float64 {
	return leaf.max
}

func (leaf *LeafD) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("LeafD: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = r.ReadF64()
	leaf.max = r.ReadF64()

	r.CheckByteCount(pos, bcnt, start, "TLeafD")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafD{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafD", f)
	Factory.add("*rootio.LeafD", f)
}

var _ Object = (*LeafD)(nil)
var _ Named = (*LeafD)(nil)
var _ Leaf = (*LeafD)(nil)
var _ ROOTUnmarshaler = (*LeafD)(nil)
