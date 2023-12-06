// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Limit struct{}

func (*Limit) Class() string {
	return "TLimit"
}

func (*Limit) RVersion() int16 {
	return rvers.Limit
}

// MarshalROOT implements rbytes.Marshaler
func (o *Limit) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(o.Class(), o.RVersion())
	return w.SetHeader(hdr)
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (o *Limit) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(o.Class(), o.RVersion())

	r.CheckHeader(hdr)
	return r.Err()
}

type LimitDataSource struct {
	base     rbase.Object   `groot:"BASE-TObject"`       // base class
	sig      rcont.ObjArray `groot:"fSignal"`            // packed input signal
	bkg      rcont.ObjArray `groot:"fBackground"`        // packed input background
	data     rcont.ObjArray `groot:"fCandidates"`        // packed input candidates (data)
	sigErr   rcont.ObjArray `groot:"fErrorOnSignal"`     // packed error sources for signal
	bkgErr   rcont.ObjArray `groot:"fErrorOnBackground"` // packed error sources for background
	ids      rcont.ObjArray `groot:"fIds"`               // packed IDs for the different error sources
	dummyTA  rcont.ObjArray `groot:"fDummyTA"`           // array of dummy object (used for bookeeping)
	dummyIDs rcont.ObjArray `groot:"fDummyIds"`          // array of dummy object (used for bookeeping)
}

func (*LimitDataSource) Class() string {
	return "TLimitDataSource"
}

func (*LimitDataSource) RVersion() int16 {
	return rvers.LimitDataSource
}

// MarshalROOT implements rbytes.Marshaler
func (o *LimitDataSource) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(o.Class(), o.RVersion())

	w.WriteObject(&o.base)
	w.WriteObject(&o.sig)
	w.WriteObject(&o.bkg)
	w.WriteObject(&o.data)
	w.WriteObject(&o.sigErr)
	w.WriteObject(&o.bkgErr)
	w.WriteObject(&o.ids)
	w.WriteObject(&o.dummyTA)
	w.WriteObject(&o.dummyIDs)

	return w.SetHeader(hdr)
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (o *LimitDataSource) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(o.Class(), o.RVersion())

	r.ReadObject(&o.base)
	r.ReadObject(&o.sig)
	r.ReadObject(&o.bkg)
	r.ReadObject(&o.data)
	r.ReadObject(&o.sigErr)
	r.ReadObject(&o.bkgErr)
	r.ReadObject(&o.ids)
	r.ReadObject(&o.dummyTA)
	r.ReadObject(&o.dummyIDs)

	r.CheckHeader(hdr)
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			var o LimitDataSource
			return reflect.ValueOf(&o)
		}
		rtypes.Factory.Add("TLimitDataSource", f)
	}

	{
		f := func() reflect.Value {
			var o Limit
			return reflect.ValueOf(&o)
		}
		rtypes.Factory.Add("TLimit", f)
	}
}

var (
	_ root.Object        = (*Limit)(nil)
	_ rbytes.RVersioner  = (*Limit)(nil)
	_ rbytes.Marshaler   = (*Limit)(nil)
	_ rbytes.Unmarshaler = (*Limit)(nil)

	_ root.Object        = (*LimitDataSource)(nil)
	_ rbytes.RVersioner  = (*LimitDataSource)(nil)
	_ rbytes.Marshaler   = (*LimitDataSource)(nil)
	_ rbytes.Unmarshaler = (*LimitDataSource)(nil)
)
