// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type attmarker struct {
	rvers int16
	color int16
	style int16
	width float32
}

func (a *attmarker) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	a.rvers = vers
	a.color = r.ReadI16()
	a.style = r.ReadI16()
	a.width = r.ReadF32()
	r.CheckByteCount(pos, bcnt, start, "TAttMarker")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &attmarker{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttMarker", f)
	Factory.add("*rootio.attmarker", f)
}

var _ ROOTUnmarshaler = (*attmarker)(nil)
