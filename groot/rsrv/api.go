// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"encoding/json"
	"image/color"
	"math"
	"strings"

	"go-hep.org/x/hep/groot/rsrv/internal/hexcolor"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

// OpenFileRequest describes a request to open a file located at the
// provided URI.
type OpenFileRequest struct {
	URI string `json:"uri"`
}

type CloseFileRequest struct {
	URI string `json:"uri"`
}

type File struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

type ListResponse struct {
	Files []File `json:"files"`
}

type DirentRequest struct {
	URI       string `json:"uri"`
	Dir       string `json:"dir,omitempty"`
	Recursive bool   `json:"recursive,omitempty"`
}

type DirentResponse struct {
	URI     string   `json:"uri"`
	Content []Dirent `json:"content,omitempty"`
}

type Dirent struct {
	Path  string `json:"path"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
	Cycle int    `json:"cycle"`
}

type Tree struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Entries  int64    `json:"entries"`
	Branches []Branch `json:"branches"`
	Leaves   []Leaf   `json:"leaves"`
}

type Branch struct {
	Type     string   `json:"type"`
	Name     string   `json:"name"`
	Branches []Branch `json:"branches"`
	Leaves   []Leaf   `json:"leaves"`
}

type Leaf struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type TreeRequest struct {
	URI string `json:"uri"`
	Dir string `json:"dir"`
	Obj string `json:"obj"`
}

type TreeResponse struct {
	URI  string `json:"uri"`
	Dir  string `json:"dir"`
	Obj  string `json:"obj"`
	Tree Tree   `json:"tree"`
}

type PlotH1Request struct {
	URI string `json:"uri"`
	Dir string `json:"dir"`
	Obj string `json:"obj"`

	Options PlotOptions `json:"options"`
}

type PlotH2Request struct {
	URI string `json:"uri"`
	Dir string `json:"dir"`
	Obj string `json:"obj"`

	Options PlotOptions `json:"options"`
}

type PlotS2Request struct {
	URI string `json:"uri"`
	Dir string `json:"dir"`
	Obj string `json:"obj"`

	Options PlotOptions `json:"options"`
}

type PlotTreeRequest struct {
	URI  string   `json:"uri"`
	Dir  string   `json:"dir"`
	Obj  string   `json:"obj"`
	Vars []string `json:"vars"`

	Options PlotOptions `json:"options"`
}

type PlotResponse struct {
	URI string `json:"uri"`
	Dir string `json:"dir"`
	Obj string `json:"obj"`

	Data string `json:"data"`
}

type PlotOptions struct {
	Title string `json:"title,omitempty"`
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`

	Type   string    `json:"type"`
	Width  vg.Length `json:"width"`
	Height vg.Length `json:"height"`

	Line      LineStyle   `json:"line,omitempty"`
	FillColor color.Color `json:"fill_color,omitempty"`
}

func (opt *PlotOptions) init() {
	if opt.Type == "" {
		opt.Type = "png"
	}
	opt.Type = strings.ToLower(opt.Type)

	switch {
	case opt.Width <= 0 && opt.Height <= 0:
		opt.Width = vgimg.DefaultWidth
		opt.Height = vgimg.DefaultWidth / math.Phi
	case opt.Width <= 0:
		opt.Width = opt.Height * math.Phi
	case opt.Height <= 0:
		opt.Height = opt.Width / math.Phi
	}

	opt.Line.init()
	if opt.FillColor == nil {
		opt.FillColor = color.Transparent
	}
}

func (opt PlotOptions) MarshalJSON() ([]byte, error) {
	if opt.FillColor == nil {
		opt.FillColor = color.Transparent
	}
	raw := jsonPlotOptions{
		Title:     opt.Title,
		X:         opt.X,
		Y:         opt.Y,
		Type:      opt.Type,
		Width:     opt.Width,
		Height:    opt.Height,
		Line:      opt.Line.toJSON(),
		FillColor: hexcolor.HexModel.Convert(opt.FillColor).(hexcolor.Hex),
	}
	return json.Marshal(raw)
}

func (opt *PlotOptions) UnmarshalJSON(p []byte) error {
	var raw jsonPlotOptions
	err := json.Unmarshal(p, &raw)
	if err != nil {
		return err
	}
	opt.Title = raw.Title
	opt.X = raw.X
	opt.Y = raw.Y
	opt.Type = raw.Type
	opt.Width = raw.Width
	opt.Height = raw.Height
	opt.Line = raw.Line.fromJSON()
	opt.FillColor = raw.FillColor
	return nil
}

type LineStyle struct {
	Color   color.Color `json:"color,omitempty"`
	Width   vg.Length   `json:"width,omitempty"`
	Dashes  []vg.Length `json:"dashes,omitempty"`
	DashOff vg.Length   `json:"dash_offset,omitempty"`
}

func (sty *LineStyle) init() {
	if sty.Color == nil {
		sty.Color = color.RGBA{255, 0, 0, 255}
	}
}

type jsonPlotOptions struct {
	Title string `json:"title,omitempty"`
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`

	Type   string    `json:"type"`
	Width  vg.Length `json:"width"`
	Height vg.Length `json:"height"`

	Line      jsonLineStyle `json:"line,omitempty"`
	FillColor hexcolor.Hex  `json:"fill_color,omitempty"`
}

type jsonLineStyle struct {
	Color   hexcolor.Hex `json:"color,omitempty"`
	Width   vg.Length    `json:"width,omitempty"`
	Dashes  []vg.Length  `json:"dashes,omitempty"`
	DashOff vg.Length    `json:"dash_offset,omitempty"`
}

func (sty jsonLineStyle) fromJSON() LineStyle {
	return LineStyle{
		Color:   sty.Color,
		Width:   sty.Width,
		Dashes:  sty.Dashes,
		DashOff: sty.DashOff,
	}
}

func (sty LineStyle) toJSON() jsonLineStyle {
	if sty.Color == nil {
		sty.Color = defaultLineColor
	}
	o := jsonLineStyle{
		Width:   sty.Width,
		Dashes:  sty.Dashes,
		DashOff: sty.DashOff,
		Color:   hexcolor.HexModel.Convert(sty.Color).(hexcolor.Hex),
	}
	return o
}

var (
	defaultLineColor = color.RGBA{R: 255, A: 255}
)
