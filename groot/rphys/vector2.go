// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys

import (
	"fmt"
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

func NewVector2(x, y float64) *Vector2 {
	return &Vector2{
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

	hdr := w.WriteHeader(vec.Class(), vec.RVersion())
	w.WriteObject(&vec.obj)
	w.WriteF64(vec.x)
	w.WriteF64(vec.y)

	return w.SetHeader(hdr)
}

func (vec *Vector2) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(vec.Class())
	if hdr.Vers > rvers.Vector2 {
		panic(fmt.Errorf(
			"rphys: invalid %s version=%d > %d",
			vec.Class(), hdr.Vers, vec.RVersion(),
		))
	}

	if hdr.Vers != 2 {
		r.ReadObject(&vec.obj)
	}

	vec.x = r.ReadF64()
	vec.y = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (vec *Vector2) String() string {
	return fmt.Sprintf("TVector2{%v, %v}", vec.x, vec.y)
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
