// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vgshiny provides a vg.Canvas implementation backed by a shiny/screen.Window
package vgshiny

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/vgimg"
	"golang.org/x/exp/shiny/screen"
)

// Canvas implements the vg.Canvas interface,
// drawing to a shiny.screen/Window buffer.
type Canvas struct {
	*vgimg.Canvas

	win screen.Window
	buf screen.Buffer

	img draw.Image
}

func New(s screen.Screen, w, h vg.Length) (*Canvas, error) {
	ww := w / vg.Inch * vg.Length(vgimg.DefaultDPI)
	hh := h / vg.Inch * vg.Length(vgimg.DefaultDPI)
	img := draw.Image(image.NewRGBA(image.Rect(0, 0, int(ww+0.5), int(hh+0.5))))
	cc := vgimg.NewWith(vgimg.UseImage(img))

	size := img.Bounds().Size()
	win, err := s.NewWindow(&screen.NewWindowOptions{
		Width:  size.X,
		Height: size.Y,
	})
	if err != nil {
		return nil, err
	}

	buf, err := s.NewBuffer(size)
	if err != nil {
		return nil, err
	}

	return &Canvas{
		win:    win,
		buf:    buf,
		Canvas: cc,
		img:    img,
	}, nil
}

func (c *Canvas) Paint() screen.PublishResult {
	w, h := c.Size()
	rect := image.Rect(0, 0, int(w), int(h))
	sr := c.img.Bounds()

	c.win.Fill(rect, color.Black, draw.Src)
	draw.Draw(c.buf.RGBA(), c.buf.Bounds(), c.img, image.Point{}, draw.Src)
	if !sr.In(rect) {
		sr = rect
	}
	c.win.Upload(image.Point{}, c.buf, sr)

	return c.win.Publish()
}

func (c *Canvas) Release() {
	c.buf.Release()
	c.win.Release()
	c.buf = nil
	c.win = nil
}
