// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type attline struct {
	rvers int16
	color int16
	style int16
	width int16
}

func (a *attline) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(a.rvers)
	w.WriteI16(a.color)
	w.WriteI16(a.style)
	w.WriteI16(a.width)
	return w.SetByteCount(pos, "TAttLine")

}

func (a *attline) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	a.rvers = vers

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

var (
	_ ROOTMarshaler   = (*attline)(nil)
	_ ROOTUnmarshaler = (*attline)(nil)
)
