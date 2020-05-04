// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rcont

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// ArrayC implements ROOT TArrayC
type ArrayC struct {
	Data []int8
}

func (*ArrayC) RVersion() int16 {
	return rvers.ArrayC
}

// Class returns the ROOT class name.
func (*ArrayC) Class() string {
	return "TArrayC"
}

func (arr *ArrayC) Len() int {
	return len(arr.Data)
}

func (arr *ArrayC) At(i int) int8 {
	return arr.Data[i]
}

func (arr *ArrayC) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayC) Set(i int, v interface{}) {
	arr.Data[i] = v.(int8)
}

func (arr *ArrayC) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI8(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayC) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeI8(arr.Data, n)
	r.ReadArrayI8(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayC{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayC", f)
}

var (
	_ root.Array         = (*ArrayC)(nil)
	_ rbytes.Marshaler   = (*ArrayC)(nil)
	_ rbytes.Unmarshaler = (*ArrayC)(nil)
)

// ArrayS implements ROOT TArrayS
type ArrayS struct {
	Data []int16
}

func (*ArrayS) RVersion() int16 {
	return rvers.ArrayS
}

// Class returns the ROOT class name.
func (*ArrayS) Class() string {
	return "TArrayS"
}

func (arr *ArrayS) Len() int {
	return len(arr.Data)
}

func (arr *ArrayS) At(i int) int16 {
	return arr.Data[i]
}

func (arr *ArrayS) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayS) Set(i int, v interface{}) {
	arr.Data[i] = v.(int16)
}

func (arr *ArrayS) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI16(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayS) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeI16(arr.Data, n)
	r.ReadArrayI16(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayS{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayS", f)
}

var (
	_ root.Array         = (*ArrayS)(nil)
	_ rbytes.Marshaler   = (*ArrayS)(nil)
	_ rbytes.Unmarshaler = (*ArrayS)(nil)
)

// ArrayI implements ROOT TArrayI
type ArrayI struct {
	Data []int32
}

func (*ArrayI) RVersion() int16 {
	return rvers.ArrayI
}

// Class returns the ROOT class name.
func (*ArrayI) Class() string {
	return "TArrayI"
}

func (arr *ArrayI) Len() int {
	return len(arr.Data)
}

func (arr *ArrayI) At(i int) int32 {
	return arr.Data[i]
}

func (arr *ArrayI) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayI) Set(i int, v interface{}) {
	arr.Data[i] = v.(int32)
}

func (arr *ArrayI) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI32(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayI) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeI32(arr.Data, n)
	r.ReadArrayI32(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayI{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayI", f)
}

var (
	_ root.Array         = (*ArrayI)(nil)
	_ rbytes.Marshaler   = (*ArrayI)(nil)
	_ rbytes.Unmarshaler = (*ArrayI)(nil)
)

// ArrayL implements ROOT TArrayL
type ArrayL struct {
	Data []int64
}

func (*ArrayL) RVersion() int16 {
	return rvers.ArrayL
}

// Class returns the ROOT class name.
func (*ArrayL) Class() string {
	return "TArrayL"
}

func (arr *ArrayL) Len() int {
	return len(arr.Data)
}

func (arr *ArrayL) At(i int) int64 {
	return arr.Data[i]
}

func (arr *ArrayL) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayL) Set(i int, v interface{}) {
	arr.Data[i] = v.(int64)
}

func (arr *ArrayL) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI64(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayL) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeI64(arr.Data, n)
	r.ReadArrayI64(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayL{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayL", f)
}

var (
	_ root.Array         = (*ArrayL)(nil)
	_ rbytes.Marshaler   = (*ArrayL)(nil)
	_ rbytes.Unmarshaler = (*ArrayL)(nil)
)

// ArrayL64 implements ROOT TArrayL64
type ArrayL64 struct {
	Data []int64
}

func (*ArrayL64) RVersion() int16 {
	return rvers.ArrayL64
}

// Class returns the ROOT class name.
func (*ArrayL64) Class() string {
	return "TArrayL64"
}

func (arr *ArrayL64) Len() int {
	return len(arr.Data)
}

func (arr *ArrayL64) At(i int) int64 {
	return arr.Data[i]
}

func (arr *ArrayL64) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayL64) Set(i int, v interface{}) {
	arr.Data[i] = v.(int64)
}

func (arr *ArrayL64) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI64(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayL64) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeI64(arr.Data, n)
	r.ReadArrayI64(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayL64{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayL64", f)
}

var (
	_ root.Array         = (*ArrayL64)(nil)
	_ rbytes.Marshaler   = (*ArrayL64)(nil)
	_ rbytes.Unmarshaler = (*ArrayL64)(nil)
)

// ArrayF implements ROOT TArrayF
type ArrayF struct {
	Data []float32
}

func (*ArrayF) RVersion() int16 {
	return rvers.ArrayF
}

// Class returns the ROOT class name.
func (*ArrayF) Class() string {
	return "TArrayF"
}

func (arr *ArrayF) Len() int {
	return len(arr.Data)
}

func (arr *ArrayF) At(i int) float32 {
	return arr.Data[i]
}

func (arr *ArrayF) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayF) Set(i int, v interface{}) {
	arr.Data[i] = v.(float32)
}

func (arr *ArrayF) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayF32(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayF) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeF32(arr.Data, n)
	r.ReadArrayF32(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayF{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayF", f)
}

var (
	_ root.Array         = (*ArrayF)(nil)
	_ rbytes.Marshaler   = (*ArrayF)(nil)
	_ rbytes.Unmarshaler = (*ArrayF)(nil)
)

// ArrayD implements ROOT TArrayD
type ArrayD struct {
	Data []float64
}

func (*ArrayD) RVersion() int16 {
	return rvers.ArrayD
}

// Class returns the ROOT class name.
func (*ArrayD) Class() string {
	return "TArrayD"
}

func (arr *ArrayD) Len() int {
	return len(arr.Data)
}

func (arr *ArrayD) At(i int) float64 {
	return arr.Data[i]
}

func (arr *ArrayD) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *ArrayD) Set(i int, v interface{}) {
	arr.Data[i] = v.(float64)
}

func (arr *ArrayD) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayF64(arr.Data)

	return int(w.Pos() - pos), w.Err()
}

func (arr *ArrayD) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	n := int(r.ReadI32())
	arr.Data = rbytes.ResizeF64(arr.Data, n)
	r.ReadArrayF64(arr.Data)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayD{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayD", f)
}

var (
	_ root.Array         = (*ArrayD)(nil)
	_ rbytes.Marshaler   = (*ArrayD)(nil)
	_ rbytes.Unmarshaler = (*ArrayD)(nil)
)
