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
	leaf tleaf
	min	int32
	max int32
}

// Name returns the name of the instance
func (leaf *LeafC) Name() string {
	return leaf.leaf.Name()
}

// Title returns the title of the instance
func (leaf *LeafC) Title() string {
	return leaf.leaf.Title()
}

// Class returns the ROOT class name.
func (leaf *LeafC) Class() string {
	return "TLeafC"
}

func (leaf *LeafC) ArrayDim() int {
	return leaf.leaf.ArrayDim()
}

func (leaf *LeafC) SetBranch(b Branch) {
	leaf.leaf.SetBranch(b)
}

func (leaf *LeafC) Branch() Branch {
	return leaf.leaf.Branch()
}

func (leaf *LeafC) HasRange() bool {
	return leaf.leaf.HasRange()
}

func (leaf *LeafC) IsUnsigned() bool {
	return leaf.leaf.IsUnsigned()
}

func (leaf *LeafC) LeafCount() Leaf {
	return leaf.leaf.LeafCount()
}

func (leaf *LeafC) Len() int {
	return leaf.leaf.Len()
}

func (leaf *LeafC) LenType() int {
	return leaf.leaf.LenType()
}

func (leaf *LeafC) MaxIndex() []int {
	return leaf.leaf.MaxIndex()
}

func (leaf *LeafC) Offset() int {
	return leaf.leaf.Offset()
}

func (leaf *LeafC) Value(i int) interface{} {
	return leaf.leaf.Value(i)
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

	if err := leaf.leaf.UnmarshalROOT(r); err != nil {
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
	leaf tleaf
	min	int32
	max int32
}

// Name returns the name of the instance
func (leaf *LeafI) Name() string {
	return leaf.leaf.Name()
}

// Title returns the title of the instance
func (leaf *LeafI) Title() string {
	return leaf.leaf.Title()
}

// Class returns the ROOT class name.
func (leaf *LeafI) Class() string {
	return "TLeafI"
}

func (leaf *LeafI) ArrayDim() int {
	return leaf.leaf.ArrayDim()
}

func (leaf *LeafI) SetBranch(b Branch) {
	leaf.leaf.SetBranch(b)
}

func (leaf *LeafI) Branch() Branch {
	return leaf.leaf.Branch()
}

func (leaf *LeafI) HasRange() bool {
	return leaf.leaf.HasRange()
}

func (leaf *LeafI) IsUnsigned() bool {
	return leaf.leaf.IsUnsigned()
}

func (leaf *LeafI) LeafCount() Leaf {
	return leaf.leaf.LeafCount()
}

func (leaf *LeafI) Len() int {
	return leaf.leaf.Len()
}

func (leaf *LeafI) LenType() int {
	return leaf.leaf.LenType()
}

func (leaf *LeafI) MaxIndex() []int {
	return leaf.leaf.MaxIndex()
}

func (leaf *LeafI) Offset() int {
	return leaf.leaf.Offset()
}

func (leaf *LeafI) Value(i int) interface{} {
	return leaf.leaf.Value(i)
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

	if err := leaf.leaf.UnmarshalROOT(r); err != nil {
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
	leaf tleaf
	min	int64
	max int64
}

// Name returns the name of the instance
func (leaf *LeafL) Name() string {
	return leaf.leaf.Name()
}

// Title returns the title of the instance
func (leaf *LeafL) Title() string {
	return leaf.leaf.Title()
}

// Class returns the ROOT class name.
func (leaf *LeafL) Class() string {
	return "TLeafL"
}

func (leaf *LeafL) ArrayDim() int {
	return leaf.leaf.ArrayDim()
}

func (leaf *LeafL) SetBranch(b Branch) {
	leaf.leaf.SetBranch(b)
}

func (leaf *LeafL) Branch() Branch {
	return leaf.leaf.Branch()
}

func (leaf *LeafL) HasRange() bool {
	return leaf.leaf.HasRange()
}

func (leaf *LeafL) IsUnsigned() bool {
	return leaf.leaf.IsUnsigned()
}

func (leaf *LeafL) LeafCount() Leaf {
	return leaf.leaf.LeafCount()
}

func (leaf *LeafL) Len() int {
	return leaf.leaf.Len()
}

func (leaf *LeafL) LenType() int {
	return leaf.leaf.LenType()
}

func (leaf *LeafL) MaxIndex() []int {
	return leaf.leaf.MaxIndex()
}

func (leaf *LeafL) Offset() int {
	return leaf.leaf.Offset()
}

func (leaf *LeafL) Value(i int) interface{} {
	return leaf.leaf.Value(i)
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

	if err := leaf.leaf.UnmarshalROOT(r); err != nil {
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
	leaf tleaf
	min	float32
	max float32
}

// Name returns the name of the instance
func (leaf *LeafF) Name() string {
	return leaf.leaf.Name()
}

// Title returns the title of the instance
func (leaf *LeafF) Title() string {
	return leaf.leaf.Title()
}

// Class returns the ROOT class name.
func (leaf *LeafF) Class() string {
	return "TLeafF"
}

func (leaf *LeafF) ArrayDim() int {
	return leaf.leaf.ArrayDim()
}

func (leaf *LeafF) SetBranch(b Branch) {
	leaf.leaf.SetBranch(b)
}

func (leaf *LeafF) Branch() Branch {
	return leaf.leaf.Branch()
}

func (leaf *LeafF) HasRange() bool {
	return leaf.leaf.HasRange()
}

func (leaf *LeafF) IsUnsigned() bool {
	return leaf.leaf.IsUnsigned()
}

func (leaf *LeafF) LeafCount() Leaf {
	return leaf.leaf.LeafCount()
}

func (leaf *LeafF) Len() int {
	return leaf.leaf.Len()
}

func (leaf *LeafF) LenType() int {
	return leaf.leaf.LenType()
}

func (leaf *LeafF) MaxIndex() []int {
	return leaf.leaf.MaxIndex()
}

func (leaf *LeafF) Offset() int {
	return leaf.leaf.Offset()
}

func (leaf *LeafF) Value(i int) interface{} {
	return leaf.leaf.Value(i)
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

	if err := leaf.leaf.UnmarshalROOT(r); err != nil {
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
	leaf tleaf
	min	float64
	max float64
}

// Name returns the name of the instance
func (leaf *LeafD) Name() string {
	return leaf.leaf.Name()
}

// Title returns the title of the instance
func (leaf *LeafD) Title() string {
	return leaf.leaf.Title()
}

// Class returns the ROOT class name.
func (leaf *LeafD) Class() string {
	return "TLeafD"
}

func (leaf *LeafD) ArrayDim() int {
	return leaf.leaf.ArrayDim()
}

func (leaf *LeafD) SetBranch(b Branch) {
	leaf.leaf.SetBranch(b)
}

func (leaf *LeafD) Branch() Branch {
	return leaf.leaf.Branch()
}

func (leaf *LeafD) HasRange() bool {
	return leaf.leaf.HasRange()
}

func (leaf *LeafD) IsUnsigned() bool {
	return leaf.leaf.IsUnsigned()
}

func (leaf *LeafD) LeafCount() Leaf {
	return leaf.leaf.LeafCount()
}

func (leaf *LeafD) Len() int {
	return leaf.leaf.Len()
}

func (leaf *LeafD) LenType() int {
	return leaf.leaf.LenType()
}

func (leaf *LeafD) MaxIndex() []int {
	return leaf.leaf.MaxIndex()
}

func (leaf *LeafD) Offset() int {
	return leaf.leaf.Offset()
}

func (leaf *LeafD) Value(i int) interface{} {
	return leaf.leaf.Value(i)
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

	if err := leaf.leaf.UnmarshalROOT(r); err != nil {
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
