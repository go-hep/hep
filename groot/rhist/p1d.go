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

// Profile1D is a 1-dim profile histogram.
type Profile1D struct {
	h1d        H1D          // base class
	binEntries rcont.ArrayD // number of entries per bin
	errMode    int32        // Option to compute errors
	ymin       float64      // Lower limit in Y (if set)
	ymax       float64      // Upper limit in Y (if set)
	sumwy      float64      // Total Sum of weight*Y
	sumwy2     float64      // Total Sum of weight*Y*Y
	binSumw2   rcont.ArrayD // Array of sum of squares of weights per bin
}

func newProfile1D() *Profile1D {
	return &Profile1D{
		h1d: *newH1D(),
	}
}

func (*Profile1D) Class() string {
	return "TProfile"
}

func (*Profile1D) RVersion() int16 {
	return rvers.Profile
}

// MarshalROOT implements rbytes.Marshaler
func (p *Profile1D) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(p.RVersion())

	if n, err := p.h1d.MarshalROOT(w); err != nil {
		return n, err
	}

	if n, err := p.binEntries.MarshalROOT(w); err != nil {
		return n, err
	}
	w.WriteI32(p.errMode)
	w.WriteF64(p.ymin)
	w.WriteF64(p.ymax)
	w.WriteF64(p.sumwy)
	w.WriteF64(p.sumwy2)
	if n, err := p.binSumw2.MarshalROOT(w); err != nil {
		return n, err
	}

	return w.SetByteCount(pos, p.Class())
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (p *Profile1D) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(p.Class())
	if vers > rvers.Profile {
		panic(fmt.Errorf("rhist: invalid TProfile version=%d > %d", vers, rvers.Profile))
	}
	if vers < 7 {
		// tested with v7.
		panic(fmt.Errorf("rhist: too old TProfile version=%d < 7", vers))
	}

	if err := p.h1d.UnmarshalROOT(r); err != nil {
		return err
	}
	if err := p.binEntries.UnmarshalROOT(r); err != nil {
		return err
	}
	p.errMode = r.ReadI32()
	p.ymin = r.ReadF64()
	p.ymax = r.ReadF64()
	p.sumwy = r.ReadF64()
	p.sumwy2 = r.ReadF64()
	if err := p.binSumw2.UnmarshalROOT(r); err != nil {
		return err
	}

	r.CheckByteCount(pos, bcnt, start, p.Class())
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		p1d := newProfile1D()
		return reflect.ValueOf(p1d)
	}
	rtypes.Factory.Add("TProfile", f)
}

var (
	_ root.Object        = (*Profile1D)(nil)
	_ rbytes.RVersioner  = (*Profile1D)(nil)
	_ rbytes.Marshaler   = (*Profile1D)(nil)
	_ rbytes.Unmarshaler = (*Profile1D)(nil)
)
