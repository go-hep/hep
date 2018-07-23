// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type attaxis struct {
	rvers   int16
	ndivs   int32   // number of divisions (10000*n3 + 100*n2 + n1)
	acolor  int16   // color of the line axis
	lcolor  int16   // color of labels
	lfont   int16   // font for labels
	loffset float32 // offset of labels
	lsize   float32 // size of labels
	ticks   float32 // length of tick marks
	toffset float32 // offset of axis title
	tsize   float32 // size of axis title
	tcolor  int16   // color of axis title
	tfont   int16   // font for axis title
}

func (*attaxis) Class() string {
	return "TAttAxis"
}

func (a *attaxis) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(a.rvers)
	w.WriteI32(a.ndivs)
	w.WriteI16(a.acolor)
	w.WriteI16(a.lcolor)
	w.WriteI16(a.lfont)
	w.WriteF32(a.loffset)
	w.WriteF32(a.lsize)
	w.WriteF32(a.ticks)
	w.WriteF32(a.toffset)
	w.WriteF32(a.tsize)
	w.WriteI16(a.tcolor)
	w.WriteI16(a.tfont)

	return w.SetByteCount(pos, "TAttAxis")
}

func (a *attaxis) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	a.rvers = vers
	if vers < 4 {
		return fmt.Errorf("rootio: TAttAxis version too old (%d < 4)", vers)
	}

	a.ndivs = r.ReadI32()
	a.acolor = r.ReadI16()
	a.lcolor = r.ReadI16()
	a.lfont = r.ReadI16()
	a.loffset = r.ReadF32()
	a.lsize = r.ReadF32()
	a.ticks = r.ReadF32()
	a.toffset = r.ReadF32()
	a.tsize = r.ReadF32()
	a.tcolor = r.ReadI16()
	a.tfont = r.ReadI16()

	r.CheckByteCount(pos, bcnt, beg, "TAttAxis")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &attaxis{}
		return reflect.ValueOf(o)
	}
	Factory.add("TAttAxis", f)
	Factory.add("*rootio.attaxis", f)
}

var (
	_ Object          = (*attaxis)(nil)
	_ ROOTMarshaler   = (*attaxis)(nil)
	_ ROOTUnmarshaler = (*attaxis)(nil)
)
