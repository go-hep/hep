// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !cross_compile

package main

import (
	"log"
	"math"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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
	msg  *log.Logger
	quit chan int
	once sync.Once
	wg   sync.WaitGroup
}

func newWinMgr(msg *log.Logger) *winMgr {
	return &winMgr{
		msg:  msg,
		quit: make(chan int),
	}
}

func (wmgr *winMgr) Close() error {
	wmgr.once.Do(wmgr.doClose)
	return nil
}

func (wmgr *winMgr) doClose() {
	close(wmgr.quit)
}

func (wmgr *winMgr) newPlot(p *hplot.Plot) *window {
	wmgr.wg.Add(1)
	win := newWindow(p)
	go win.run(wmgr)
	return win
}

type window struct {
	w     *app.Window
	ready chan int

	mu  sync.Mutex
	plt *hplot.Plot
}

func newWindow(p *hplot.Plot) *window {
	title := p.Plot.Title.Text
	switch title {
	case "":
		title = "PAW-Go"
	default:
		title = "PAW-Go [" + title + "]"
	}

	x := unit.Px(float32(xmax.Dots(dpi)))
	y := unit.Px(float32(ymax.Dots(dpi)))

	win := &window{
		w:     app.NewWindow(app.Title(title), app.Size(x, y)),
		plt:   p,
		ready: make(chan int),
	}
	return win
}

func (w *window) run(wmgr *winMgr) {
	defer wmgr.wg.Done()
	close(w.ready)

	for {
		select {
		case e := <-w.w.Events():
			o := w.handle(e)
			if o == winStop {
				return
			}
		case <-wmgr.quit:
			w.w.Close()
			return
		}
	}
}

type winState byte

const (
	winContinue winState = iota
	winStop
)

func (w *window) handle(e event.Event) winState {
	switch e := e.(type) {
	case system.DestroyEvent:
		return winStop
	case system.FrameEvent:
		cnv := vggio.New(
			layout.NewContext(new(op.Ops), e),
			xmax, ymax,
			vggio.UseDPI(dpi),
		)
		w.mu.Lock()
		w.plt.Draw(draw.New(cnv))
		w.mu.Unlock()
		e.Frame(cnv.Paint())

	case key.Event:
		switch e.Name {
		case "Q", key.NameEscape:
			w.w.Invalidate()
			w.w.Close()
		}
	}
	return winContinue
}
