// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vgshiny provides a vg.Canvas implementation backed by a shiny/screen.Window
package vgshiny // import "go-hep.org/x/hep/hplot/vgshiny"

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

// Canvas implements the vg.Canvas interface,
// drawing to a shiny.screen/Window buffer.
type Canvas struct {
	*vgimg.Canvas

	win screen.Window
	buf screen.Buffer

	img draw.Image
}

// New creates a new canvas with the given width and height.
func New(s screen.Screen, w, h vg.Length) (*Canvas, error) {
	ww := w / vg.Inch * vg.Length(vgimg.DefaultDPI)
	hh := h / vg.Inch * vg.Length(vgimg.DefaultDPI)
	size := image.Pt(int(ww+0.5), int(hh+0.5))
	img := draw.Image(image.NewRGBA(image.Rect(0, 0, size.X, size.Y)))
	cc := vgimg.NewWith(vgimg.UseImage(img))

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

// Paint paints the canvas' content on the screen.
func (c *Canvas) Paint() screen.PublishResult {
	w, h := c.Size()
	rect := image.Rect(0, 0, int(w), int(h))
	sr := c.img.Bounds()

	c.win.Fill(rect, color.Black, draw.Src)
	draw.Draw(c.buf.RGBA(), c.buf.Bounds(), c.img, image.Point{}, draw.Src)
	c.win.Upload(image.Point{}, c.buf, sr)

	return c.win.Publish()
}

// Release releases shiny/screen resources.
func (c *Canvas) Release() {
	c.buf.Release()
	c.win.Release()
	c.buf = nil
	c.win = nil
}

// Send sends an event to the underlying shiny window.
func (c *Canvas) Send(evt interface{}) {
	c.win.Send(evt)
}

// Run runs the function f for each event on the event queue of the underlying shiny window.
// f is expected to return true to continue processing events and false otherwise.
// If f is nil, a default processing function will be used.
// The default processing functions handles paint.Event events and exits when 'q' or 'ESC' are pressed.
func (c *Canvas) Run(f func(e interface{}) bool) {
	if f == nil {
		f = func(e interface{}) bool {
			switch e := e.(type) {
			case paint.Event:
				c.Paint()
			case key.Event:
				switch e.Code {
				case key.CodeEscape, key.CodeQ:
					if e.Direction == key.DirPress {
						return false
					}
				}
			}
			return true

		}
	}
	for {
		e := c.win.NextEvent()
		if !f(e) {
			return
		}
	}
}
