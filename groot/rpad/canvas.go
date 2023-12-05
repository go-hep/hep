// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpad

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type Canvas struct {
	pad             Pad
	fCatt           AttCanvas // Canvas attributes
	fDISPLAY        string    // Name of destination screen
	fXsizeUser      float32   // User specified size of canvas along X in CM
	fYsizeUser      float32   // User specified size of canvas along Y in CM
	fXsizeReal      float32   // Current size of canvas along X in CM
	fYsizeReal      float32   // Current size of canvas along Y in CM
	fHighLightColor int16     // Highlight color of active pad
	fDoubleBuffer   int32     // Double buffer flag (0=off, 1=on)
	fWindowTopX     int32     // Top X position of window (in pixels)
	fWindowTopY     int32     // Top Y position of window (in pixels)
	fWindowWidth    uint32    // Width of window (including borders, etc.)
	fWindowHeight   uint32    // Height of window (including menubar, borders, etc.)
	fCw             uint32    // Width of the canvas along X (pixels)
	fCh             uint32    // Height of the canvas along Y (pixels)
	fRetained       bool      // Retain structure flag
}

func (*Canvas) RVersion() int16 {
	return rvers.Canvas
}

func (*Canvas) Class() string {
	return "TCanvas"
}

func (c *Canvas) Name() string {
	return c.pad.Name()
}

func (c *Canvas) Title() string {
	return c.pad.Title()
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (c *Canvas) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(c.Class())
	if hdr.Vers > rvers.Canvas {
		panic(fmt.Errorf(
			"rpad: invalid %s version=%d > %d",
			c.Class(), hdr.Vers, c.RVersion(),
		))
	}

	r.ReadObject(&c.pad)
	c.fDISPLAY = r.ReadString()
	c.fDoubleBuffer = r.ReadI32()
	c.fRetained = r.ReadBool()

	c.fXsizeUser = r.ReadF32()
	c.fYsizeUser = r.ReadF32()
	c.fXsizeReal = r.ReadF32()
	c.fYsizeReal = r.ReadF32()
	c.fWindowTopX = r.ReadI32()
	c.fWindowTopY = r.ReadI32()
	if hdr.Vers > 2 {
		c.fWindowWidth = r.ReadU32()
		c.fWindowHeight = r.ReadU32()
	}
	c.fCw = r.ReadU32()
	c.fCh = r.ReadU32()
	if hdr.Vers <= 2 {
		c.fWindowWidth = c.fCw
		c.fWindowHeight = c.fCh
	}

	r.ReadObject(&c.fCatt)
	if r.Err() != nil {
		panic(r.Err())
	}

	_ = r.ReadBool() // kMoveOpaque
	_ = r.ReadBool() // kResizeOpaque

	c.fHighLightColor = r.ReadI16()
	_ = r.ReadBool() // fBatch
	if hdr.Vers < 2 {
		r.CheckHeader(hdr)
		return r.Err()
	}
	_ = r.ReadBool() // kShowEventStatus
	if hdr.Vers > 3 {
		_ = r.ReadBool() // kAutoExec
	}
	_ = r.ReadBool() // kMenuBar

	r.CheckHeader(hdr)
	return r.Err()
}

// Keys implements the ObjectFinder interface.
func (c *Canvas) Keys() []string {
	return c.pad.Keys()
}

// Get implements the ObjectFinder interface.
func (c *Canvas) Get(name string) (root.Object, error) {
	return c.pad.Get(name)
}

func init() {
	f := func() reflect.Value {
		var c Canvas
		return reflect.ValueOf(&c)
	}
	rtypes.Factory.Add("TCanvas", f)
}

var (
	_ root.Object        = (*Canvas)(nil)
	_ root.Named         = (*Canvas)(nil)
	_ root.ObjectFinder  = (*Canvas)(nil)
	_ rbytes.Unmarshaler = (*Canvas)(nil)
)
