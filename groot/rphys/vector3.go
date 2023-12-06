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

type Vector3 struct {
	obj rbase.Object
	x   float64
	y   float64
	z   float64
}

func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{
		obj: *rbase.NewObject(),
		x:   x,
		y:   y,
		z:   z,
	}
}

func (*Vector3) RVersion() int16 {
	return rvers.Vector3
}

func (*Vector3) Class() string {
	return "TVector3"
}

func (vec *Vector3) X() float64 { return vec.x }
func (vec *Vector3) Y() float64 { return vec.y }
func (vec *Vector3) Z() float64 { return vec.z }

func (vec *Vector3) SetX(x float64) { vec.x = x }
func (vec *Vector3) SetY(y float64) { vec.y = y }
func (vec *Vector3) SetZ(z float64) { vec.z = z }

func (vec *Vector3) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(vec.Class(), vec.RVersion())
	w.WriteObject(&vec.obj)
	w.WriteF64(vec.x)
	w.WriteF64(vec.y)
	w.WriteF64(vec.z)

	return w.SetHeader(hdr)
}

func (vec *Vector3) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(vec.Class(), vec.RVersion())

	if hdr.Vers == 1 || hdr.Vers > 2 {
		r.ReadObject(&vec.obj)
	}

	vec.x = r.ReadF64()
	vec.y = r.ReadF64()
	vec.z = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (vec *Vector3) String() string {
	return fmt.Sprintf(
		"TVector3{%v, %v, %v}",
		vec.x, vec.y, vec.z,
	)
}

func init() {
	{
		f := func() reflect.Value {
			o := &Vector3{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TVector3", f)
	}
}

var (
	_ root.Object        = (*Vector3)(nil)
	_ rbytes.Marshaler   = (*Vector3)(nil)
	_ rbytes.Unmarshaler = (*Vector3)(nil)
)
