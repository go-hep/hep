// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgtex"
)

// Drawer is the interface that wraps the Draw method.
type Drawer interface {
	Draw(draw.Canvas)
}

// Save saves the plot to an image file.  The file format is determined
// by the extension.
//
// Supported extensions are:
//
//  .eps, .jpg, .jpeg, .pdf, .png, .svg, .tex, .tif and .tiff.
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func Save(p Drawer, w, h vg.Length, file string) (err error) {
	switch {
	case w <= 0 && h <= 0:
		w = vgimg.DefaultWidth
		h = vgimg.DefaultWidth / math.Phi
	case w <= 0:
		w = h * math.Phi
	case h <= 0:
		h = w / math.Phi
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	format := strings.ToLower(filepath.Ext(file))
	if len(format) != 0 {
		format = format[1:]
	}

	dc, err := draw.NewFormattedCanvas(w, h, format)
	if err != nil {
		return err
	}

	p.Draw(draw.New(dc))

	_, err = dc.WriteTo(f)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	if format == "tex" {
		if p, ok := p.(*wplot); ok {
			err = p.latex.CompileLatex(file)
			if err != nil {
				return fmt.Errorf("hplot: could not generate PDF: %w", err)
			}
		}
	}

	return nil
}

func vgtexBorder(dc draw.Canvas) {
	switch dc.Canvas.(type) {
	case *vgtex.Canvas:
		// prevent pgf/tikz to crop-out the bounding box
		// by filling the whole image with a transparent box.
		dc.FillPolygon(color.Transparent, []vg.Point{
			{X: dc.Min.X, Y: dc.Min.Y},
			{X: dc.Max.X, Y: dc.Min.Y},
			{X: dc.Max.X, Y: dc.Max.Y},
			{X: dc.Min.X, Y: dc.Max.Y},
		})
	}
}
