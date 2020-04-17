// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"go-hep.org/x/exp/vgshiny"
	"go-hep.org/x/hep/hplot"
	"golang.org/x/exp/shiny/screen"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	xmax = vg.Length(400)
	ymax = vg.Length(400 / math.Phi)
)

type winMgr struct {
	scr screen.Screen
}

func newWinMgr(scr screen.Screen) *winMgr {
	return &winMgr{
		scr: scr,
	}
}

func (wmgr *winMgr) newPlot(p *hplot.Plot) error {
	cnv, err := vgshiny.New(wmgr.scr, xmax, ymax)
	if err != nil {
		return err
	}
	p.Draw(draw.New(cnv))
	cnv.Paint()
	go func() {
		cnv.Run(nil)
		cnv.Release()
	}()

	return err
}
