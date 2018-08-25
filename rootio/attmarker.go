// Copyright 2017 The go-hep Authors. All rights reserved.
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

func newAttMarker() *attmarker {
	return &attmarker{
		rvers: 2, // FIXME(sbinet): harmonize versions
		color: 1,
		style: 1,
		width: 1,
	}
}

func (a *attmarker) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(a.rvers)
	w.WriteI16(a.color)
	w.WriteI16(a.style)
	w.WriteF32(a.width)
	return w.SetByteCount(pos, "TAttMarker")
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
		o := newAttMarker()
		return reflect.ValueOf(o)
	}
	Factory.add("TAttMarker", f)
	Factory.add("*rootio.attmarker", f)
}

var (
	_ ROOTMarshaler   = (*attmarker)(nil)
	_ ROOTUnmarshaler = (*attmarker)(nil)
)
