// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type attline struct {
	color int16
	style int16
	width int16
}

func (a *attline) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("attline-vers=%v\n", vers)

	a.color = r.ReadI16()
	a.style = r.ReadI16()
	a.width = r.ReadI16()
	r.CheckByteCount(pos, bcnt, start, "TAttLine")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &attline{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttLine", f)
	Factory.add("*rootio.attline", f)
}

var _ ROOTUnmarshaler = (*attline)(nil)
