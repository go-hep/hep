// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// Profile2D is a 2-dim profile histogram.
type Profile2D struct {
	h2d        H2D          // base class
	binEntries rcont.ArrayD // number of entries per bin
	errMode    int32        // Option to compute errors
	zmin       float64      // Lower limit in Z (if set)
	zmax       float64      // Upper limit in Z (if set)
	sumwz      float64      // Total Sum of weight*Z
	sumwz2     float64      // Total Sum of weight*Z*Z
	binSumw2   rcont.ArrayD // Array of sum of squares of weights per bin
}

func newProfile2D() *Profile2D {
	return &Profile2D{
		h2d: *newH2D(),
	}
}

func (*Profile2D) Class() string {
	return "TProfile2D"
}

func (*Profile2D) RVersion() int16 {
	return rvers.Profile2D
}

// MarshalROOT implements rbytes.Marshaler
func (p2d *Profile2D) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(p2d.RVersion())

	w.WriteObject(&p2d.h2d)
	w.WriteObject(&p2d.binEntries)
	w.WriteI32(p2d.errMode)
	w.WriteF64(p2d.zmin)
	w.WriteF64(p2d.zmax)
	w.WriteF64(p2d.sumwz)
	w.WriteF64(p2d.sumwz2)
	w.WriteObject(&p2d.binSumw2)

	return w.SetByteCount(pos, p2d.Class())
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (p2d *Profile2D) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(p2d.Class())
	if vers > rvers.Profile2D {
		panic(fmt.Errorf("rhist: invalid TProfile2D version=%d > %d", vers, rvers.Profile2D))
	}
	if vers < 8 {
		// tested with v8.
		panic(fmt.Errorf("rhist: too old TProfile2D version=%d < 8", vers))
	}

	r.ReadObject(&p2d.h2d)
	r.ReadObject(&p2d.binEntries)
	p2d.errMode = r.ReadI32()
	p2d.zmin = r.ReadF64()
	p2d.zmax = r.ReadF64()
	p2d.sumwz = r.ReadF64()
	p2d.sumwz2 = r.ReadF64()
	r.ReadObject(&p2d.binSumw2)

	r.CheckByteCount(pos, bcnt, start, p2d.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		p2d := newProfile2D()
		return reflect.ValueOf(p2d)
	}
	rtypes.Factory.Add("TProfile2D", f)
}

var (
	_ root.Object        = (*Profile2D)(nil)
	_ rbytes.RVersioner  = (*Profile2D)(nil)
	_ rbytes.Marshaler   = (*Profile2D)(nil)
	_ rbytes.Unmarshaler = (*Profile2D)(nil)
)
