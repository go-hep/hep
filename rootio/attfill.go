// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type attfill struct {
	color int16
	style int16
}

func (a *attfill) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	a.color = r.ReadI16()
	a.style = r.ReadI16()
	r.CheckByteCount(pos, bcnt, start, "TAttFill")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &attfill{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttFill", f)
	Factory.add("*rootio.attfill", f)
}

var _ ROOTUnmarshaler = (*attfill)(nil)
