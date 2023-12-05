// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpad

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Pad struct {
	vpad rbase.VirtualPad
	//	bbox2d rbase.Object

	fX1               float64 // X of lower X coordinate
	fY1               float64 // Y of lower Y coordinate
	fX2               float64 // X of upper X coordinate
	fY2               float64 // Y of upper Y coordinate
	fXtoAbsPixelk     float64 // Conversion coefficient for X World to absolute pixel
	fXtoPixelk        float64 // Conversion coefficient for X World to pixel
	fXtoPixel         float64 // xpixel = fXtoPixelk + fXtoPixel*xworld
	fYtoAbsPixelk     float64 // Conversion coefficient for Y World to absolute pixel
	fYtoPixelk        float64 // Conversion coefficient for Y World to pixel
	fYtoPixel         float64 // ypixel = fYtoPixelk + fYtoPixel*yworld
	fUtoAbsPixelk     float64 // Conversion coefficient for U NDC to absolute pixel
	fUtoPixelk        float64 // Conversion coefficient for U NDC to pixel
	fUtoPixel         float64 // xpixel = fUtoPixelk + fUtoPixel*undc
	fVtoAbsPixelk     float64 // Conversion coefficient for V NDC to absolute pixel
	fVtoPixelk        float64 // Conversion coefficient for V NDC to pixel
	fVtoPixel         float64 // ypixel = fVtoPixelk + fVtoPixel*vndc
	fAbsPixeltoXk     float64 // Conversion coefficient for absolute pixel to X World
	fPixeltoXk        float64 // Conversion coefficient for pixel to X World
	fPixeltoX         float64 // xworld = fPixeltoXk + fPixeltoX*xpixel
	fAbsPixeltoYk     float64 // Conversion coefficient for absolute pixel to Y World
	fPixeltoYk        float64 // Conversion coefficient for pixel to Y World
	fPixeltoY         float64 // yworld = fPixeltoYk + fPixeltoY*ypixel
	fXlowNDC          float64 // X bottom left corner of pad in NDC [0,1]
	fYlowNDC          float64 // Y bottom left corner of pad in NDC [0,1]
	fXUpNDC           float64
	fYUpNDC           float64
	fWNDC             float64     // Width of pad along X in Normalized Coordinates (NDC)
	fHNDC             float64     // Height of pad along Y in Normalized Coordinates (NDC)
	fAbsXlowNDC       float64     // Absolute X top left corner of pad in NDC [0,1]
	fAbsYlowNDC       float64     // Absolute Y top left corner of pad in NDC [0,1]
	fAbsWNDC          float64     // Absolute Width of pad along X in NDC
	fAbsHNDC          float64     // Absolute Height of pad along Y in NDC
	fUxmin            float64     // Minimum value on the X axis
	fUymin            float64     // Minimum value on the Y axis
	fUxmax            float64     // Maximum value on the X axis
	fUymax            float64     // Maximum value on the Y axis
	fTheta            float64     // theta angle to view as lego/surface
	fPhi              float64     // phi angle   to view as lego/surface
	fAspectRatio      float64     // ratio of w/h in case of fixed ratio
	fNumber           int32       // pad number identifier
	fTickx            int32       // Set to 1 if tick marks along X
	fTicky            int32       // Set to 1 if tick marks along Y
	fLogx             int32       // (=0 if X linear scale, =1 if log scale)
	fLogy             int32       // (=0 if Y linear scale, =1 if log scale)
	fLogz             int32       // (=0 if Z linear scale, =1 if log scale)
	fPadPaint         int32       // Set to 1 while painting the pad
	fCrosshair        int32       // Crosshair type (0 if no crosshair requested)
	fCrosshairPos     int32       // Position of crosshair
	fBorderSize       int16       // pad bordersize in pixels
	fBorderMode       int16       // Bordermode (-1=down, 0 = no border, 1=up)
	fModified         bool        // Set to true when pad is modified
	fGridx            bool        // Set to true if grid along X
	fGridy            bool        // Set to true if grid along Y
	fAbsCoord         bool        // Use absolute coordinates
	fEditable         bool        // True if canvas is editable
	fFixedAspectRatio bool        // True if fixed aspect ratio
	fPrimitives       *rcont.List // ->List of primitives (subpads)
	fExecs            *rcont.List // List of commands to be executed when a pad event occurs
	fName             string      // Pad name
	fTitle            string      // Pad title
	fNumPaletteColor  int32       // Number of objects with an automatic color
	fNextPaletteColor int32       // Next automatic color
}

func (*Pad) RVersion() int16 {
	return rvers.Pad
}

func (*Pad) Class() string {
	return "TPad"
}

func (p *Pad) Name() string {
	return p.fName
}

func (p *Pad) Title() string {
	return p.fTitle
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (p *Pad) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(p.Class())
	if hdr.Vers > rvers.Pad {
		panic(fmt.Errorf(
			"rpad: invalid %s version=%d > %d",
			p.Class(), hdr.Vers, p.RVersion(),
		))
	}

	if hdr.Vers != 13 {
		panic(fmt.Errorf(
			"rpad: invalid %s version=%d > %d",
			p.Class(), hdr.Vers, p.RVersion(),
		))
	}

	r.ReadObject(&p.vpad)

	_ = r.ReadHeader("TAttBBox2D")

	p.fX1 = r.ReadF64()
	p.fY1 = r.ReadF64()
	p.fX2 = r.ReadF64()
	p.fY2 = r.ReadF64()
	p.fXtoAbsPixelk = r.ReadF64()
	p.fXtoPixelk = r.ReadF64()
	p.fXtoPixel = r.ReadF64()
	p.fYtoAbsPixelk = r.ReadF64()
	p.fYtoPixelk = r.ReadF64()
	p.fYtoPixel = r.ReadF64()
	p.fUtoAbsPixelk = r.ReadF64()
	p.fUtoPixelk = r.ReadF64()
	p.fUtoPixel = r.ReadF64()
	p.fVtoAbsPixelk = r.ReadF64()
	p.fVtoPixelk = r.ReadF64()
	p.fVtoPixel = r.ReadF64()
	p.fAbsPixeltoXk = r.ReadF64()
	p.fPixeltoXk = r.ReadF64()
	p.fPixeltoX = r.ReadF64()
	p.fAbsPixeltoYk = r.ReadF64()
	p.fPixeltoYk = r.ReadF64()
	p.fPixeltoY = r.ReadF64()
	p.fXlowNDC = r.ReadF64()
	p.fYlowNDC = r.ReadF64()
	p.fXUpNDC = r.ReadF64()
	p.fYUpNDC = r.ReadF64()
	p.fWNDC = r.ReadF64()
	p.fHNDC = r.ReadF64()
	p.fAbsXlowNDC = r.ReadF64()
	p.fAbsYlowNDC = r.ReadF64()
	p.fAbsWNDC = r.ReadF64()
	p.fAbsHNDC = r.ReadF64()
	p.fUxmin = r.ReadF64()
	p.fUymin = r.ReadF64()
	p.fUxmax = r.ReadF64()
	p.fUymax = r.ReadF64()
	p.fTheta = r.ReadF64()
	p.fPhi = r.ReadF64()
	p.fAspectRatio = r.ReadF64()
	p.fNumber = r.ReadI32()
	p.fTickx = r.ReadI32()
	p.fTicky = r.ReadI32()
	p.fLogx = r.ReadI32()
	p.fLogy = r.ReadI32()
	p.fLogz = r.ReadI32()
	p.fPadPaint = r.ReadI32()
	p.fCrosshair = r.ReadI32()
	p.fCrosshairPos = r.ReadI32()
	p.fBorderSize = r.ReadI16()
	p.fBorderMode = r.ReadI16()
	p.fModified = r.ReadBool()
	p.fGridx = r.ReadBool()
	p.fGridy = r.ReadBool()
	p.fAbsCoord = r.ReadBool()
	p.fEditable = r.ReadBool()
	p.fFixedAspectRatio = r.ReadBool()

	{
		var prims rcont.List
		r.ReadObject(&prims)
		if prims.Len() > 0 {
			p.fPrimitives = &prims
		}
	}

	{
		execs := r.ReadObjectAny()
		if execs != nil {
			p.fExecs = execs.(*rcont.List)
		}
	}

	p.fName = r.ReadString()
	p.fTitle = r.ReadString()

	p.fNumPaletteColor = r.ReadI32()
	p.fNextPaletteColor = r.ReadI32()

	r.CheckHeader(hdr)
	return r.Err()
}

// Keys implements the ObjectFinder interface.
func (pad *Pad) Keys() []string {
	var keys []string
	if pad.fPrimitives != nil && pad.fPrimitives.Len() > 0 {
		for i := 0; i < pad.fPrimitives.Len(); i++ {
			o, ok := pad.fPrimitives.At(i).(root.Named)
			if !ok {
				continue
			}
			keys = append(keys, o.Name())
		}
	}
	return keys
}

// Get implements the ObjectFinder interface.
func (pad *Pad) Get(name string) (root.Object, error) {
	for i := 0; i < pad.fPrimitives.Len(); i++ {
		v := pad.fPrimitives.At(i)
		o, ok := v.(root.Named)
		if !ok {
			continue
		}
		if o.Name() == name {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no object named %q", name)
}

func init() {
	f := func() reflect.Value {
		var p Pad
		return reflect.ValueOf(&p)
	}
	rtypes.Factory.Add("TPad", f)
}

var (
	_ root.Object        = (*Pad)(nil)
	_ root.Named         = (*Pad)(nil)
	_ root.ObjectFinder  = (*Pad)(nil)
	_ rbytes.Unmarshaler = (*Pad)(nil)
)
