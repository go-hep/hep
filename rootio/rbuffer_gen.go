// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"encoding/binary"
	"math"
)

// ReadI8 reads a int8 from a ROOT buffer.
func (r *RBuffer) ReadI8(v *int8) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 1
	*v = int8((r.r.p[beg:r.r.c][0]))
}

// ReadFastArrayI8 reads a slice of int8 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayI8(v []int8) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 1
		v[i] = int8((r.r.p[beg:r.r.c][0]))
	}
}

// ReadI16 reads a int16 from a ROOT buffer.
func (r *RBuffer) ReadI16(v *int16) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 2
	*v = int16(binary.BigEndian.Uint16(r.r.p[beg:r.r.c]))
}

// ReadFastArrayI16 reads a slice of int16 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayI16(v []int16) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 2
		v[i] = int16(binary.BigEndian.Uint16(r.r.p[beg:r.r.c]))
	}
}

// ReadI32 reads a int32 from a ROOT buffer.
func (r *RBuffer) ReadI32(v *int32) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 4
	*v = int32(binary.BigEndian.Uint32(r.r.p[beg:r.r.c]))
}

// ReadFastArrayI32 reads a slice of int32 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayI32(v []int32) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 4
		v[i] = int32(binary.BigEndian.Uint32(r.r.p[beg:r.r.c]))
	}
}

// ReadI64 reads a int64 from a ROOT buffer.
func (r *RBuffer) ReadI64(v *int64) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 8
	*v = int64(binary.BigEndian.Uint64(r.r.p[beg:r.r.c]))
}

// ReadFastArrayI64 reads a slice of int64 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayI64(v []int64) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 8
		v[i] = int64(binary.BigEndian.Uint64(r.r.p[beg:r.r.c]))
	}
}

// ReadU8 reads a uint8 from a ROOT buffer.
func (r *RBuffer) ReadU8(v *uint8) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 1
	*v = uint8(r.r.p[beg:r.r.c][0])
}

// ReadFastArrayU8 reads a slice of uint8 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayU8(v []uint8) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 1
		v[i] = uint8(r.r.p[beg:r.r.c][0])
	}
}

// ReadU16 reads a uint16 from a ROOT buffer.
func (r *RBuffer) ReadU16(v *uint16) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 2
	*v = binary.BigEndian.Uint16(r.r.p[beg:r.r.c])
}

// ReadFastArrayU16 reads a slice of uint16 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayU16(v []uint16) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 2
		v[i] = binary.BigEndian.Uint16(r.r.p[beg:r.r.c])
	}
}

// ReadU32 reads a uint32 from a ROOT buffer.
func (r *RBuffer) ReadU32(v *uint32) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 4
	*v = binary.BigEndian.Uint32(r.r.p[beg:r.r.c])
}

// ReadFastArrayU32 reads a slice of uint32 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayU32(v []uint32) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 4
		v[i] = binary.BigEndian.Uint32(r.r.p[beg:r.r.c])
	}
}

// ReadU64 reads a uint64 from a ROOT buffer.
func (r *RBuffer) ReadU64(v *uint64) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 8
	*v = binary.BigEndian.Uint64(r.r.p[beg:r.r.c])
}

// ReadFastArrayU64 reads a slice of uint64 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayU64(v []uint64) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 8
		v[i] = binary.BigEndian.Uint64(r.r.p[beg:r.r.c])
	}
}

// ReadF32 reads a float32 from a ROOT buffer.
func (r *RBuffer) ReadF32(v *float32) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 4
	*v = math.Float32frombits(binary.BigEndian.Uint32(r.r.p[beg:r.r.c]))
}

// ReadFastArrayF32 reads a slice of float32 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayF32(v []float32) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 4
		v[i] = math.Float32frombits(binary.BigEndian.Uint32(r.r.p[beg:r.r.c]))
	}
}

// ReadF64 reads a float64 from a ROOT buffer.
func (r *RBuffer) ReadF64(v *float64) {
	if r.err != nil {
		return
	}

	beg := r.r.c
	r.r.c += 8
	*v = math.Float64frombits(binary.BigEndian.Uint64(r.r.p[beg:r.r.c]))
}

// ReadFastArrayF64 reads a slice of float64 from a ROOT buffer.
func (r *RBuffer) ReadFastArrayF64(v []float64) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}
	for i := range v {
		beg := r.r.c
		r.r.c += 8
		v[i] = math.Float64frombits(binary.BigEndian.Uint64(r.r.p[beg:r.r.c]))
	}
}
