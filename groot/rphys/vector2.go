// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Vector2 struct {
	obj rbase.Object
	x   float64
	y   float64
}

func NewVector2(x, y float64) *Vector3 {
	return &Vector3{
		obj: *rbase.NewObject(),
		x:   x,
		y:   y,
	}
}

func (*Vector2) RVersion() int16 {
	return rvers.Vector2
}

func (*Vector2) Class() string {
	return "TVector2"
}

func (vec *Vector2) X() float64 { return vec.x }
func (vec *Vector2) Y() float64 { return vec.y }

func (vec *Vector2) SetX(x float64) { vec.x = x }
func (vec *Vector2) SetY(y float64) { vec.y = y }

func (vec *Vector2) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(vec.RVersion())
	if _, err := vec.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteF64(vec.x)
	w.WriteF64(vec.y)

	return w.SetByteCount(pos, vec.Class())
}

func (vec *Vector2) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(vec.Class())

	if vers != 2 {
		if err := vec.obj.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	vec.x = r.ReadF64()
	vec.y = r.ReadF64()

	r.CheckByteCount(pos, bcnt, beg, vec.Class())
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := &Vector2{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TVector2", f)
	}
}

var (
	_ root.Object        = (*Vector2)(nil)
	_ rbytes.Marshaler   = (*Vector2)(nil)
	_ rbytes.Unmarshaler = (*Vector2)(nil)
)
