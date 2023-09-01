// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop // import "go-hep.org/x/hep/hplot/vgop"

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

// op is a vector graphics operation as defined by the vg.Canvas interface.
// Each method of vg.Canvas has a corresponding vgop.op.
type op interface {
	op(ctx fontCtx, c vg.Canvas) error
}

// opSetLineWidth corresponds to the vg.Canvas.SetWidth method.
type opSetLineWidth struct {
	Width vg.Length `json:"width"`
}

func (op opSetLineWidth) op(ctx fontCtx, c vg.Canvas) error {
	c.SetLineWidth(op.Width)
	return nil
}

// opSetLineDash corresponds to the vg.Canvas.SetLineDash method.
type opSetLineDash struct {
	Dashes  []vg.Length `json:"dashes"`
	Offsets vg.Length   `json:"offsets"`
}

func (op opSetLineDash) op(ctx fontCtx, c vg.Canvas) error {
	c.SetLineDash(op.Dashes, op.Offsets)
	return nil
}

// opSetColor corresponds to the vg.Canvas.SetColor method.
type opSetColor struct {
	Color color.Color
}

func (op opSetColor) MarshalJSON() ([]byte, error) {
	var ctx struct {
		C struct {
			R uint32 `json:"r"`
			G uint32 `json:"g"`
			B uint32 `json:"b"`
			A uint32 `json:"a"`
		} `json:"color"`
	}

	ctx.C.R, ctx.C.G, ctx.C.B, ctx.C.A = op.Color.RGBA()
	return json.Marshal(ctx)
}

func (op *opSetColor) UnmarshalJSON(p []byte) error {
	var ctx struct {
		C struct {
			R uint32 `json:"r"`
			G uint32 `json:"g"`
			B uint32 `json:"b"`
			A uint32 `json:"a"`
		} `json:"color"`
	}

	err := json.Unmarshal(p, &ctx)
	if err != nil {
		return err
	}
	op.Color = &color.RGBA{
		R: uint8(ctx.C.R),
		G: uint8(ctx.C.G),
		B: uint8(ctx.C.B),
		A: uint8(ctx.C.A),
	}
	return nil
}

func (op opSetColor) op(ctx fontCtx, c vg.Canvas) error {
	c.SetColor(op.Color)
	return nil
}

// opRotate corresponds to the vg.Canvas.Rotate method.
type opRotate struct {
	Angle float64 `json:"angle"`
}

func (op opRotate) op(ctx fontCtx, c vg.Canvas) error {
	c.Rotate(op.Angle)
	return nil
}

// opTranslate corresponds to the vg.Canvas.Translate method.
type opTranslate struct {
	Point vg.Point
}

func (op opTranslate) MarshalJSON() ([]byte, error) {
	var v struct {
		P jsonPoint `json:"point"`
	}
	v.P = jsonPointFrom(op.Point)

	return json.Marshal(v)
}

func (op *opTranslate) UnmarshalJSON(p []byte) error {
	var v struct {
		P jsonPoint `json:"point"`
	}
	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}
	op.Point = v.P.cnv()

	return nil
}

func (op opTranslate) op(ctx fontCtx, c vg.Canvas) error {
	c.Translate(op.Point)
	return nil
}

// opScale corresponds to the vg.Canvas.Scale method.
type opScale struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (op opScale) op(ctx fontCtx, c vg.Canvas) error {
	c.Scale(op.X, op.Y)
	return nil
}

// opPush corresponds to the vg.Canvas.Push method.
type opPush struct{}

func (op opPush) op(ctx fontCtx, c vg.Canvas) error {
	c.Push()
	return nil
}

// opPop corresponds to the vg.Canvas.Pop method.
type opPop struct{}

func (op opPop) op(ctx fontCtx, c vg.Canvas) error {
	c.Pop()
	return nil
}

// opStroke corresponds to the vg.Canvas.Stroke method.
type opStroke struct {
	Path vg.Path
}

func (op opStroke) MarshalJSON() ([]byte, error) {
	var v struct {
		P jsonPath `json:"path"`
	}
	v.P = jsonPathFrom(op.Path)

	return json.Marshal(v)
}

func (op *opStroke) UnmarshalJSON(p []byte) error {
	var v struct {
		P jsonPath `json:"path"`
	}
	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}
	op.Path = v.P.cnv()

	return nil
}

func (op opStroke) op(ctx fontCtx, c vg.Canvas) error {
	c.Stroke(op.Path)
	return nil
}

// opFill corresponds to the vg.Canvas.Fill method.
type opFill struct {
	Path vg.Path
}

func (op opFill) MarshalJSON() ([]byte, error) {
	var v struct {
		P jsonPath `json:"path"`
	}
	v.P = jsonPathFrom(op.Path)

	return json.Marshal(v)
}

func (op *opFill) UnmarshalJSON(p []byte) error {
	var v struct {
		P jsonPath `json:"path"`
	}
	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}
	op.Path = v.P.cnv()

	return nil
}

func (op opFill) op(ctx fontCtx, c vg.Canvas) error {
	c.Fill(op.Path)
	return nil
}

// opFillString corresponds to the vg.Canvas.FillString method.
type opFillString struct {
	Font   font.Font
	Point  vg.Point
	String string
}

func (op opFillString) MarshalJSON() ([]byte, error) {
	var v struct {
		ID  fontID    `json:"font"`
		Pt  jsonPoint `json:"point"`
		Str string    `json:"string"`
	}
	v.ID = fontID(op.Font)
	v.Pt = jsonPointFrom(op.Point)
	v.Str = op.String

	return json.Marshal(v)
}

func (op *opFillString) UnmarshalJSON(p []byte) error {
	var v struct {
		ID  fontID    `json:"font"`
		Pt  jsonPoint `json:"point"`
		Str string    `json:"string"`
	}
	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}
	op.Font = font.Font(v.ID)
	op.Point = v.Pt.cnv()
	op.String = v.Str

	return nil
}

func (op opFillString) op(ctx fontCtx, c vg.Canvas) error {
	id := fontID(op.Font)
	face, ok := ctx.fonts[id]
	if !ok {
		return fmt.Errorf("unknown font name=%q, size=%v", op.Font.Name(), op.Font.Size)
	}
	c.FillString(face, op.Point, op.String)
	return nil
}

// opDrawImage corresponds to the vg.Canvas.DrawImage method
type opDrawImage struct {
	Rect  vg.Rectangle
	Image image.Image
}

func (op opDrawImage) MarshalJSON() ([]byte, error) {
	var v struct {
		Rect struct {
			Min jsonPoint `json:"min"`
			Max jsonPoint `json:"max"`
		}
		Image []byte `json:"data"`
	}
	v.Rect.Min = jsonPointFrom(op.Rect.Min)
	v.Rect.Max = jsonPointFrom(op.Rect.Max)

	var buf bytes.Buffer
	err := png.Encode(&buf, op.Image)
	if err != nil {
		return nil, fmt.Errorf("could not encode image to PNG: %w", err)
	}
	v.Image = []byte(base64.StdEncoding.EncodeToString(buf.Bytes()))

	return json.Marshal(v)
}

func (op *opDrawImage) UnmarshalJSON(p []byte) error {
	var v struct {
		Rect struct {
			Min jsonPoint `json:"min"`
			Max jsonPoint `json:"max"`
		}
		Image []byte `json:"data"`
	}
	err := json.Unmarshal(p, &v)
	if err != nil {
		return err
	}
	op.Rect.Min = v.Rect.Min.cnv()
	op.Rect.Max = v.Rect.Max.cnv()

	buf, err := base64.StdEncoding.DecodeString(string(v.Image))
	if err != nil {
		return err
	}
	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return err
	}
	op.Image = img

	return nil
}

func (op opDrawImage) op(ctx fontCtx, c vg.Canvas) error {
	c.DrawImage(op.Rect, op.Image)
	return nil
}
