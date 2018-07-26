// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
)

// ArrayI implements ROOT TArrayI
type ArrayI struct {
	Data []int32
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

func (arr *ArrayI) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI32(arr.Data)

	return int(w.Pos() - pos), w.err
}

func (arr *ArrayI) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = r.ReadFastArrayI32(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayI{}
		return reflect.ValueOf(o)
	}
	Factory.add("TArrayI", f)
	Factory.add("*rootio.ArrayI", f)
}

var (
	_ Array           = (*ArrayI)(nil)
	_ ROOTMarshaler   = (*ArrayI)(nil)
	_ ROOTUnmarshaler = (*ArrayI)(nil)
)

// ArrayL64 implements ROOT TArrayL64
type ArrayL64 struct {
	Data []int64
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

func (arr *ArrayL64) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayI64(arr.Data)

	return int(w.Pos() - pos), w.err
}

func (arr *ArrayL64) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = r.ReadFastArrayI64(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayL64{}
		return reflect.ValueOf(o)
	}
	Factory.add("TArrayL64", f)
	Factory.add("*rootio.ArrayL64", f)
}

var (
	_ Array           = (*ArrayL64)(nil)
	_ ROOTMarshaler   = (*ArrayL64)(nil)
	_ ROOTUnmarshaler = (*ArrayL64)(nil)
)

// ArrayF implements ROOT TArrayF
type ArrayF struct {
	Data []float32
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

func (arr *ArrayF) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayF32(arr.Data)

	return int(w.Pos() - pos), w.err
}

func (arr *ArrayF) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = r.ReadFastArrayF32(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayF{}
		return reflect.ValueOf(o)
	}
	Factory.add("TArrayF", f)
	Factory.add("*rootio.ArrayF", f)
}

var (
	_ Array           = (*ArrayF)(nil)
	_ ROOTMarshaler   = (*ArrayF)(nil)
	_ ROOTUnmarshaler = (*ArrayF)(nil)
)

// ArrayD implements ROOT TArrayD
type ArrayD struct {
	Data []float64
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

func (arr *ArrayD) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteI32(int32(len(arr.Data)))
	w.WriteFastArrayF64(arr.Data)

	return int(w.Pos() - pos), w.err
}

func (arr *ArrayD) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = r.ReadFastArrayF64(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &ArrayD{}
		return reflect.ValueOf(o)
	}
	Factory.add("TArrayD", f)
	Factory.add("*rootio.ArrayD", f)
}

var (
	_ Array           = (*ArrayD)(nil)
	_ ROOTMarshaler   = (*ArrayD)(nil)
	_ ROOTUnmarshaler = (*ArrayD)(nil)
)
