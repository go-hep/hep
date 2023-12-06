// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpad

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type AttCanvas struct {
	fXBetween     float32 // X distance between pads
	fYBetween     float32 // Y distance between pads
	fTitleFromTop float32 // Y distance of Global Title from top
	fXdate        float32 // X position where to draw the date
	fYdate        float32 // X position where to draw the date
	fAdate        float32 // Alignment for the date
}

func (*AttCanvas) RVersion() int16 {
	return rvers.AttCanvas
}

func (*AttCanvas) Class() string {
	return "TAttCanvas"
}

func init() {
	f := func() reflect.Value {
		var v AttCanvas
		return reflect.ValueOf(&v)
	}
	rtypes.Factory.Add("TAttCanvas", f)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (att *AttCanvas) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(att.Class(), att.RVersion())

	att.fXBetween = r.ReadF32()
	att.fYBetween = r.ReadF32()
	att.fTitleFromTop = r.ReadF32()
	att.fXdate = r.ReadF32()
	att.fYdate = r.ReadF32()
	att.fAdate = r.ReadF32()

	r.CheckHeader(hdr)
	return r.Err()
}

var (
	_ root.Object        = (*AttCanvas)(nil)
	_ rbytes.Unmarshaler = (*AttCanvas)(nil)
)
