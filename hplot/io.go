// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"image/color"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgeps"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgpdf"
	"gonum.org/v1/plot/vg/vgsvg"
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
	w, h = Dims(w, h)

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	format := strings.ToLower(filepath.Ext(file))
	if len(format) != 0 {
		format = format[1:]
	}

	dc, err := WriterTo(p, w, h, format)
	if err != nil {
		return err
	}

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

// WriterTo returns an io.WriterTo that will write the plots as
// the specified image format.
//
// Supported formats are the same ones than hplot.Plot.WriterTo
//
// If w or h are <= 0, the value is chosen such that it follows the Golden Ratio.
// If w and h are <= 0, the values are chosen such that they follow the Golden Ratio
// (the width is defaulted to vgimg.DefaultWidth).
func WriterTo(p Drawer, w, h vg.Length, format string) (io.WriterTo, error) {
	w, h = Dims(w, h)

	dpi := float64(vgimg.DefaultDPI)
	if p, ok := p.(*P); ok {
		dpi = p.DPI
	}

	c, err := newFormattedCanvas(w, h, format, dpi)
	if err != nil {
		return nil, fmt.Errorf("hplot: could not create canvas: %w", err)
	}
	p.Draw(draw.New(c))

	return c, nil
}

// newFormattedCanvas creates a new vg.CanvasWriterTo with the specified
// image format.
//
// Supported formats are:
//
//  eps, jpg|jpeg, pdf, png, svg, tex and tif|tiff.
func newFormattedCanvas(w, h vg.Length, format string, dpi float64) (vg.CanvasWriterTo, error) {
	var c vg.CanvasWriterTo
	switch format {
	case "eps":
		c = vgeps.New(w, h)

	case "jpg", "jpeg":
		c = vgimg.JpegCanvas{Canvas: vgimg.NewWith(
			vgimg.UseDPI(int(dpi)),
			vgimg.UseWH(w, h),
		)}

	case "pdf":
		c = vgpdf.New(w, h)

	case "png":
		c = vgimg.PngCanvas{Canvas: vgimg.NewWith(
			vgimg.UseDPI(int(dpi)),
			vgimg.UseWH(w, h),
		)}

	case "svg":
		c = vgsvg.New(w, h)

	case "tex":
		c = vgtex.NewDocument(w, h)

	case "tif", "tiff":
		c = vgimg.TiffCanvas{Canvas: vgimg.NewWith(
			vgimg.UseDPI(int(dpi)),
			vgimg.UseWH(w, h),
		)}

	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
	return c, nil
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

func Dims(width, height vg.Length) (w, h vg.Length) {
	w = width
	h = height
	switch {
	case w <= 0 && h <= 0:
		w = vgimg.DefaultWidth
		h = vgimg.DefaultWidth / math.Phi
	case w <= 0:
		w = h * math.Phi
	case h <= 0:
		h = w / math.Phi
	}
	return w, h
}
