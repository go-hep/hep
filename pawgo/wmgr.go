// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !cross_compile

package main

import (
	"log"
	"math"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vggio"
)

const (
	xmax = vg.Length(400)
	ymax = vg.Length(400 / math.Phi)
	dpi  = vggio.DefaultDPI // FIXME(sbinet): remove?
)

type winMgr struct {
	msg *log.Logger
}

func newWinMgr(msg *log.Logger) *winMgr {
	return &winMgr{
		msg: msg,
	}
}

func (wmgr *winMgr) newPlot(p *hplot.Plot) error {
	win := app.NewWindow(
		app.Title("PAW-Go"),
		app.Size(
			unit.Px(float32(xmax.Dots(dpi))),
			unit.Px(float32(ymax.Dots(dpi))),
		),
	)

	go func() {
		defer win.Close()

		for e := range win.Events() {
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				cnv := vggio.New(e, xmax, ymax, vggio.UseDPI(dpi))
				p.Draw(draw.New(cnv))
				cnv.Paint(e)

			case key.Event:
				switch e.Name {
				case "Q", key.NameEscape:
					return
				}
			}
		}
	}()

	return nil
}
