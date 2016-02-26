package main

import (
	"math"

	"github.com/go-hep/hplot"
	"github.com/go-hep/hplot/vgshiny"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"

	"golang.org/x/exp/shiny/screen"
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
