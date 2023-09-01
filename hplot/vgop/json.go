// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop // import "go-hep.org/x/hep/hplot/vgop"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

const (
	typSetLineWidth = "SetLineWidth"
	typSetLineDash  = "SetLineDash"
	typSetColor     = "SetColor"
	typRotate       = "Rotate"
	typTranslate    = "Translate"
	typScale        = "Scale"
	typPush         = "Push"
	typPop          = "Pop"
	typStroke       = "Stroke"
	typFill         = "Fill"
	typFillString   = "FillString"
	typDrawImage    = "DrawImage"
)

var (
	_ vg.Canvas         = (*JSON)(nil)
	_ vg.CanvasSizer    = (*JSON)(nil)
	_ vg.CanvasWriterTo = (*JSON)(nil)

	_ json.Marshaler   = (*JSON)(nil)
	_ json.Unmarshaler = (*JSON)(nil)
)

// JSON implements JSON serialization for vg.Canvas.
type JSON struct {
	*Canvas
}

// NewJSON creates a new vg.Canvas for JSON serialization.
func NewJSON(opts ...Option) *JSON {
	return &JSON{New(opts...)}
}

func (c *JSON) WriteTo(w io.Writer) (int64, error) {
	ww := cwriter{w: w}
	enc := json.NewEncoder(&ww)
	enc.SetIndent("", " ")
	err := enc.Encode(c)
	if err != nil {
		return 0, fmt.Errorf("could not encode canvas to JSON: %w", err)
	}
	return int64(ww.n), nil
}

func (c *JSON) MarshalJSON() ([]byte, error) {
	jc := jsonCanvas{
		Size:  jsonSize{Width: c.w, Height: c.h},
		Ops:   make([]jsonOp, len(c.ops)),
		Fonts: make([]fontID, len(c.ctx.fonts)),
	}

	var err error
	for i, op := range c.ops {
		jc.Ops[i], err = jsonOpFrom(op)
		if err != nil {
			return nil, err
		}
	}

	keys := make([]fontID, 0, len(jc.Fonts))
	for k := range c.ctx.fonts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ii := font.Font(keys[i])
		jj := font.Font(keys[j])
		if ii.Name() != jj.Name() {
			return ii.Name() < jj.Name()
		}
		return ii.Size < jj.Size
	})
	copy(jc.Fonts, keys)

	return json.Marshal(jc)
}

func (c *JSON) UnmarshalJSON(p []byte) error {
	var (
		jc  jsonCanvas
		err error
	)

	err = json.NewDecoder(bytes.NewReader(p)).Decode(&jc)
	if err != nil {
		return err
	}

	c.w = jc.Size.Width
	c.h = jc.Size.Height

	c.ctx.fonts = make(map[fontID]font.Face, len(jc.Fonts))
	for _, fnt := range jc.Fonts {
		c.ctx.fonts[fnt] = font.Face{
			Font: font.Font(fnt),
			Face: nil,
		}
	}

	c.ops = make([]op, 0, len(jc.Ops))
	for _, jop := range jc.Ops {
		var o op
		switch jop.Type {
		case typSetLineWidth:
			var v opSetLineWidth
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typSetLineDash:
			var v opSetLineDash
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typSetColor:
			var v opSetColor
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typRotate:
			var v opRotate
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typTranslate:
			var v opTranslate
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typScale:
			var v opScale
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typPush:
			o = opPush{}
		case typPop:
			o = opPop{}
		case typStroke:
			var v opStroke
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typFill:
			var v opFill
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typFillString:
			var v opFillString
			err = json.Unmarshal(jop.Value, &v)
			o = v
		case typDrawImage:
			var v opDrawImage
			err = json.Unmarshal(jop.Value, &v)
			o = v
		default:
			return fmt.Errorf("invalid vgop type %q", jop.Type)
		}
		if err != nil {
			return err
		}
		c.ops = append(c.ops, o)
	}

	return nil
}

type jsonCanvas struct {
	Size  jsonSize `json:"size"`
	Ops   []jsonOp `json:"ops"`
	Fonts []fontID `json:"fonts,omitempty"`
}

type jsonSize struct {
	Width  vg.Length `json:"width"`
	Height vg.Length `json:"height"`
}

type jsonOp struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func jsonOpFrom(op op) (jsonOp, error) {

	name := "???"
	switch op := op.(type) {
	case opSetLineWidth:
		name = typSetLineWidth
	case opSetLineDash:
		name = typSetLineDash
	case opSetColor:
		name = typSetColor
	case opRotate:
		name = typRotate
	case opTranslate:
		name = typTranslate
	case opScale:
		name = typScale
	case opPush:
		name = typPush
	case opPop:
		name = typPop
	case opStroke:
		name = typStroke
	case opFill:
		name = typFill
	case opFillString:
		name = typFillString
	case opDrawImage:
		name = typDrawImage
	default:
		return jsonOp{}, fmt.Errorf("invalid vgop.op %T", op)
	}

	v, err := json.Marshal(op)
	if err != nil {
		return jsonOp{}, fmt.Errorf("could not encode %T: %w", op, err)
	}
	return jsonOp{Type: name, Value: json.RawMessage(v)}, nil
}

type jsonPoint struct {
	X json.Number `json:"x"`
	Y json.Number `json:"y"`
}

func jsonPointFrom(p vg.Point) jsonPoint {
	return jsonPoint{
		X: json.Number(strconv.FormatFloat(p.X.Points(), 'g', -1, 64)),
		Y: json.Number(strconv.FormatFloat(p.Y.Points(), 'g', -1, 64)),
	}
}

func (p jsonPoint) cnv() vg.Point {
	x, err := p.X.Float64()
	if err != nil {
		panic(err)
	}
	y, err := p.Y.Float64()
	if err != nil {
		panic(err)
	}
	return vg.Point{X: vg.Length(x), Y: vg.Length(y)}
}

type jsonPath []jsonPathComp

func jsonPathFrom(p vg.Path) jsonPath {
	o := make(jsonPath, len(p))
	for i := range p {
		o[i] = jsonPathCompFrom(p[i])
	}
	return o
}

func (jp jsonPath) cnv() vg.Path {
	o := make(vg.Path, len(jp))
	for i := range jp {
		o[i] = jp[i].cnv()
	}

	return o
}

type jsonPathComp struct {
	Type   int         `json:"type"`
	Pos    jsonPoint   `json:"pos,omitempty"`
	Radius vg.Length   `json:"radius,omitempty"`
	Start  float64     `json:"start,omitempty"`
	Angle  float64     `json:"angle,omitempty"`
	Ctl    []jsonPoint `json:"ctl,omitempty"`
}

func jsonPathCompFrom(p vg.PathComp) jsonPathComp {
	o := jsonPathComp{
		Type:   p.Type,
		Pos:    jsonPointFrom(p.Pos),
		Radius: p.Radius,
		Start:  p.Start,
		Angle:  p.Angle,
	}
	if len(p.Control) > 0 {
		o.Ctl = make([]jsonPoint, len(p.Control))
		for i, v := range p.Control {
			o.Ctl[i] = jsonPointFrom(v)
		}
	}

	return o
}

func (jp jsonPathComp) cnv() vg.PathComp {
	o := vg.PathComp{
		Type:   jp.Type,
		Pos:    jp.Pos.cnv(),
		Radius: jp.Radius,
		Start:  jp.Start,
		Angle:  jp.Angle,
	}
	if len(jp.Ctl) > 0 {
		o.Control = make([]vg.Point, len(jp.Ctl))
		for i, v := range jp.Ctl {
			o.Control[i] = v.cnv()
		}
	}
	return o
}

type cwriter struct {
	w io.Writer
	n int
}

func (w *cwriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.n += n
	return n, err
}

var _ io.Writer = (*cwriter)(nil)
