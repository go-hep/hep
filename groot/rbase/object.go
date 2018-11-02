// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Object struct {
	ID   uint32
	Bits uint32
}

func NewObject() *Object {
	return &Object{ID: 0x0, Bits: 0x3000000}
}

func (*Object) Class() string {
	return "TObject"
}

func (*Object) RVersion() int {
	return rvers.Object
}

func (obj *Object) SetID(id uint32)     { obj.ID = id }
func (obj *Object) SetBits(bits uint32) { obj.Bits = bits }

func (obj *Object) UnmarshalROOT(r *rbytes.RBuffer) error {
	r.SkipVersion("")
	obj.ID = r.ReadU32()
	obj.Bits = r.ReadU32()
	obj.Bits |= kIsOnHeap
	if obj.Bits&kIsReferenced != 0 {
		_ = r.ReadU16()
	}
	return r.Err()
}

func (obj *Object) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	n := w.Pos()
	w.WriteU16(uint16(obj.RVersion()))
	w.WriteU32(obj.ID)
	w.WriteU32(obj.Bits)
	if obj.Bits&kIsReferenced != 0 {
		w.WriteU16(0) // FIXME(sbinet): implement referenced objects.
		panic(errors.Errorf("rbase: writing referenced objects are not supported"))
	}

	return int(w.Pos() - n), w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &Object{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TObject", f)
}

var (
	_ root.Object        = (*Object)(nil)
	_ rbytes.Marshaler   = (*Object)(nil)
	_ rbytes.Unmarshaler = (*Object)(nil)
)
