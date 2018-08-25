// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type attfill struct {
	rvers int16
	color int16
	style int16
}

func newAttFill() *attfill {
	return &attfill{
		rvers: 2, // FIXME(sbinet): harmonize versions
		color: 0,
		style: 1001, // FIXME(sbinet)
	}
}

func (a *attfill) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(a.rvers)
	w.WriteI16(a.color)
	w.WriteI16(a.style)
	return w.SetByteCount(pos, "TAttFill")

}

func (a *attfill) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	a.rvers = vers

	a.color = r.ReadI16()
	a.style = r.ReadI16()
	r.CheckByteCount(pos, bcnt, start, "TAttFill")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := newAttFill()
		return reflect.ValueOf(o)
	}
	Factory.add("TAttFill", f)
	Factory.add("*rootio.attfill", f)
}

var (
	_ ROOTMarshaler   = (*attfill)(nil)
	_ ROOTUnmarshaler = (*attfill)(nil)
)
