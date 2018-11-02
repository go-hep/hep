// Copyright 2018 The go-hep Authors. All rights reserved.
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
	arr.Data = r.ReadFastArrayI32(n)

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
	arr.Data = r.ReadFastArrayI64(n)

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
	arr.Data = r.ReadFastArrayF32(n)

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
	arr.Data = r.ReadFastArrayF64(n)

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
