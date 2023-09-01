// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop // import "go-hep.org/x/hep/hplot/vgop"

import (
	"fmt"
	"image"
	"image/color"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

const (
	// DefaultWidth and DefaultHeight are the default canvas
	// dimensions.
	DefaultWidth  = 4 * vg.Inch
	DefaultHeight = 4 * vg.Inch
)

var (
	_ vg.Canvas      = (*Canvas)(nil)
	_ vg.CanvasSizer = (*Canvas)(nil)
)

// Canvas implements vg.Canvas for serialization.
type Canvas struct {
	w, h vg.Length

	ops []op
	ctx fontCtx
}

type Option func(c *Canvas)

func WithSize(w, h vg.Length) Option {
	return func(c *Canvas) {
		c.w = w
		c.h = h
	}
}

// New returns a new canvas.
func New(opts ...Option) *Canvas {
	c := &Canvas{
		w: DefaultWidth,
		h: DefaultHeight,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Canvas) Size() (w, h vg.Length) {
	return c.w, c.h
}

// Reset resets the canvas to the base state.
func (c *Canvas) Reset() {
	c.ops = c.ops[:0]

	// FIXME(sbinet): should we also reset the fonts context ?
}

// ReplayOn replays the set of vector graphics operations onto the
// destination canvas.
func (c *Canvas) ReplayOn(dst vg.Canvas) error {
	if c.ctx.fonts == nil {
		c.ctx.fonts = make(map[fontID]font.Face)
	}

	for _, op := range c.ops {
		err := op.op(c.ctx, dst)
		if err != nil {
			return fmt.Errorf("could not apply op %T: %w", op, err)
		}
	}

	return nil
}

func (c *Canvas) add(face font.Face) {
	if c.ctx.fonts == nil {
		c.ctx.fonts = make(map[fontID]font.Face)
	}

	id := fontID(face.Font)
	c.ctx.fonts[id] = face
}

func (c *Canvas) append(o op) {
	c.ops = append(c.ops, o)
}

// SetLineWidth implements the SetLineWidth method of the vg.Canvas interface.
func (c *Canvas) SetLineWidth(w vg.Length) {
	c.append(opSetLineWidth{Width: w})
}

// SetLineDash implements the SetLineDash method of the vg.Canvas interface.
func (c *Canvas) SetLineDash(dashes []vg.Length, offs vg.Length) {
	c.append(opSetLineDash{
		Dashes:  append([]vg.Length(nil), dashes...),
		Offsets: offs,
	})
}

// SetColor implements the SetColor method of the vg.Canvas interface.
func (c *Canvas) SetColor(col color.Color) {
	c.append(opSetColor{Color: col})
}

// Rotate implements the Rotate method of the vg.Canvas interface.
func (c *Canvas) Rotate(a float64) {
	c.append(opRotate{Angle: a})
}

// Translate implements the Translate method of the vg.Canvas interface.
func (c *Canvas) Translate(pt vg.Point) {
	c.append(opTranslate{Point: pt})
}

// Scale implements the Scale method of the vg.Canvas interface.
func (c *Canvas) Scale(x, y float64) {
	c.append(opScale{X: x, Y: y})
}

// Push implements the Push method of the vg.Canvas interface.
func (c *Canvas) Push() {
	c.append(opPush{})
}

// Pop implements the Pop method of the vg.Canvas interface.
func (c *Canvas) Pop() {
	c.append(opPop{})
}

// Stroke implements the Stroke method of the vg.Canvas interface.
func (c *Canvas) Stroke(path vg.Path) {
	c.append(opStroke{Path: append(vg.Path(nil), path...)})
}

// Fill implements the Fill method of the vg.Canvas interface.
func (c *Canvas) Fill(path vg.Path) {
	c.append(opFill{Path: append(vg.Path(nil), path...)})
}

// FillString implements the FillString method of the vg.Canvas interface.
func (c *Canvas) FillString(font font.Face, pt vg.Point, str string) {
	c.add(font)
	c.append(opFillString{
		Font:   font.Font,
		Point:  pt,
		String: str,
	})
}

// DrawImage implements the DrawImage method of the vg.Canvas interface.
func (c *Canvas) DrawImage(rect vg.Rectangle, img image.Image) {
	c.append(opDrawImage{
		Rect:  rect,
		Image: img,
	})
}
