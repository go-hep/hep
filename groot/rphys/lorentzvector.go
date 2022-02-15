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

type LorentzVector struct {
	obj rbase.Object
	p   Vector3 // 3-vector component
	e   float64 // time or energy
}

func NewLorentzVector(px, py, pz, e float64) *LorentzVector {
	return &LorentzVector{
		obj: *rbase.NewObject(),
		p:   *NewVector3(px, py, pz),
		e:   e,
	}
}

func (*LorentzVector) RVersion() int16 {
	return rvers.LorentzVector
}

func (*LorentzVector) Class() string {
	return "TLorentzVector"
}

func (vec *LorentzVector) Px() float64 { return vec.p.x }
func (vec *LorentzVector) Py() float64 { return vec.p.y }
func (vec *LorentzVector) Pz() float64 { return vec.p.z }
func (vec *LorentzVector) E() float64  { return vec.e }

func (vec *LorentzVector) SetPxPyPzE(px, py, pz, e float64) {
	vec.p.x = px
	vec.p.y = py
	vec.p.z = pz
	vec.e = e
}

func (vec *LorentzVector) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(vec.Class(), vec.RVersion())
	w.WriteObject(&vec.obj)
	w.WriteObject(&vec.p)
	w.WriteF64(vec.e)

	return w.SetHeader(hdr)
}

func (vec *LorentzVector) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(vec.Class())
	if hdr.Vers > rvers.LorentzVector {
		panic(fmt.Errorf(
			"rphys: invalid %s version=%d > %d",
			vec.Class(), hdr.Vers, vec.RVersion(),
		))
	}

	r.ReadObject(&vec.obj)
	r.ReadObject(&vec.p)
	vec.e = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (vec *LorentzVector) String() string {
	return fmt.Sprintf(
		"TLorentzVector{P: {%v, %v, %v}, E: %v}",
		vec.p.x, vec.p.y, vec.p.z,
		vec.e,
	)
}

func init() {
	{
		f := func() reflect.Value {
			o := &LorentzVector{}
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TLorentzVector", f)
	}
}

var (
	_ root.Object        = (*LorentzVector)(nil)
	_ rbytes.Marshaler   = (*LorentzVector)(nil)
	_ rbytes.Unmarshaler = (*LorentzVector)(nil)
)
