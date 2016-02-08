package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/go-hep/hplot/vgshiny"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	vgdraw "github.com/gonum/plot/vg/draw"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
)

const (
	NPOINTS = 100000
	xmax    = 400
	ymax    = 400
)

var (
	bkgCol = color.Black
)

func newPlot() (*hplot.Plot, error) {
	// Draw some random values from the standard
	// normal distribution.
	hist1 := hbook.NewH1D(100, -5, +5)
	hist2 := hbook.NewH1D(100, -5, +5)
	for i := 0; i < NPOINTS; i++ {
		v1 := rand.NormFloat64() - 1
		v2 := rand.NormFloat64() + 1
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	// Make a plot and set its title.
	p, err := hplot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1, err := hplot.NewH1D(hist1)
	if err != nil {
		return nil, err
	}

	h2, err := hplot.NewH1D(hist2)
	if err != nil {
		return nil, err
	}

	h1.Infos.Style = hplot.HInfoSummary
	h2.Infos.Style = hplot.HInfoNone

	h1.Color = color.Black
	h1.FillColor = nil
	h2.Color = color.RGBA{255, 0, 0, 255}
	h2.FillColor = nil

	p.Add(h1)
	p.Add(h2)

	p.Add(plotter.NewGrid())
	return p, err
}

func main() {
	driver.Main(func(scr screen.Screen) {
		{
			p, err := newPlot()
			if err != nil {
				log.Fatal(err)
			}
			c, err := p.Show(-1, -1, scr)
			if err != nil {
				log.Fatal(err)
			}
			defer c.Release()
		}
		w, err := newWidget(scr, image.Point{xmax, ymax})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		for {
			switch e := w.canvas.NextEvent().(type) {
			case key.Event:
				repaint := false
				switch e.Code {
				case key.CodeEscape, key.CodeQ:
					return
				case key.CodeR:
					if e.Direction == key.DirPress {
						repaint = true
					}

				case key.CodeN, key.CodeSpacebar:
					p, err := newPlot()
					if err != nil {
						log.Fatal(err)
					}
					p.Draw(vgdraw.New(w.canvas))
					repaint = true
				}
				if repaint {
					w.canvas.Send(paint.Event{})
				}

			case paint.Event:
				w.canvas.Paint()
			}
		}
	})
}

type widget struct {
	s      screen.Screen
	canvas *vgshiny.Canvas
}

func newWidget(s screen.Screen, size image.Point) (*widget, error) {
	c, err := vgshiny.New(s, vg.Length(size.X), vg.Length(size.Y))
	if err != nil {
		return nil, err
	}

	return &widget{s: s, canvas: c}, err
}

func (w *widget) Release() {
	if w.canvas != nil {
		w.canvas.Release()
		w.canvas = nil
	}
	w.s = nil
}
