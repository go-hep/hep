// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type AttPad struct {
	fLeftMargin      float32 // LeftMargin
	fRightMargin     float32 // RightMargin
	fBottomMargin    float32 // BottomMargin
	fTopMargin       float32 // TopMargin
	fXfile           float32 // X position where to draw the file name
	fYfile           float32 // Y position where to draw the file name
	fAfile           float32 // Alignment for the file name
	fXstat           float32 // X position where to draw the statistics
	fYstat           float32 // Y position where to draw the statistics
	fAstat           float32 // Alignment for the statistics
	fFrameFillColor  int16   // Pad frame fill color
	fFrameLineColor  int16   // Pad frame line color
	fFrameFillStyle  int16   // Pad frame fill style
	fFrameLineStyle  int16   // Pad frame line style
	fFrameLineWidth  int16   // Pad frame line width
	fFrameBorderSize int16   // Pad frame border size
	fFrameBorderMode int32   // Pad frame border mode
}

func (*AttPad) RVersion() int16 {
	return rvers.AttPad
}

func (*AttPad) Class() string {
	return "TAttPad"
}

func (a *AttPad) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(a.Class(), a.RVersion())

	w.WriteF32(a.fLeftMargin)
	w.WriteF32(a.fRightMargin)
	w.WriteF32(a.fBottomMargin)
	w.WriteF32(a.fTopMargin)
	w.WriteF32(a.fXfile)
	w.WriteF32(a.fYfile)
	w.WriteF32(a.fAfile)
	w.WriteF32(a.fXstat)
	w.WriteF32(a.fYstat)
	w.WriteF32(a.fAstat)
	w.WriteI16(a.fFrameFillColor)
	w.WriteI16(a.fFrameLineColor)
	w.WriteI16(a.fFrameFillStyle)
	w.WriteI16(a.fFrameLineStyle)
	w.WriteI16(a.fFrameLineWidth)
	w.WriteI16(a.fFrameBorderSize)
	w.WriteI32(a.fFrameBorderMode)

	return w.SetHeader(hdr)
}

func (a *AttPad) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(a.Class(), a.RVersion())
	a.fLeftMargin = r.ReadF32()
	a.fRightMargin = r.ReadF32()
	a.fBottomMargin = r.ReadF32()
	a.fTopMargin = r.ReadF32()
	a.fXfile = r.ReadF32()
	a.fYfile = r.ReadF32()
	a.fAfile = r.ReadF32()
	a.fXstat = r.ReadF32()
	a.fYstat = r.ReadF32()
	a.fAstat = r.ReadF32()
	if hdr.Vers > 1 {
		a.fFrameFillColor = r.ReadI16()
		a.fFrameLineColor = r.ReadI16()
		a.fFrameFillStyle = r.ReadI16()
		a.fFrameLineStyle = r.ReadI16()
		a.fFrameLineWidth = r.ReadI16()
		a.fFrameBorderSize = r.ReadI16()
		a.fFrameBorderMode = r.ReadI32()
	}

	r.CheckHeader(hdr)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		var v AttPad
		return reflect.ValueOf(&v)
	}
	rtypes.Factory.Add("TAttPad", f)
}

var (
	_ root.Object        = (*AttPad)(nil)
	_ rbytes.Marshaler   = (*AttPad)(nil)
	_ rbytes.Unmarshaler = (*AttPad)(nil)
)
