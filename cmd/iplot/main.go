package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math/rand"

	"github.com/go-hep/hbook"
	"github.com/go-hep/hplot"
	"github.com/gonum/plot/plotter"
	vgdraw "github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgimg"
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
	driver.Main(func(s screen.Screen) {
		w, err := newWindow(s, image.Point{xmax, ymax})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		img := image.NewRGBA(image.Rect(0, 0, xmax, ymax))
		for {
			switch e := w.w.NextEvent().(type) {
			case key.Event:
				repaint := false
				switch e.Code {
				case key.CodeEscape, key.CodeQ:
					return
				case key.CodeR:
					if e.Direction == key.DirPress {
						repaint = true
					}

				case key.CodeN:
					p, err := newPlot()
					if err != nil {
						log.Fatal(err)
					}
					img = image.NewRGBA(image.Rect(0, 0, xmax, ymax))
					c := vgimg.NewWith(vgimg.UseImage(img))
					p.Draw(vgdraw.New(c))
					repaint = true
				}
				if repaint {
					w.w.Send(paint.Event{})
				}

			case paint.Event:
				w.display(img)
			}
		}
	})
}

type window struct {
	s screen.Screen
	w screen.Window
	b screen.Buffer
}

func newWindow(s screen.Screen, size image.Point) (*window, error) {
	w, err := s.NewWindow(
		&screen.NewWindowOptions{
			Width:  size.X,
			Height: size.Y,
		},
	)
	if err != nil {
		return nil, err
	}
	buf, err := s.NewBuffer(size)
	if err != nil {
		return nil, err
	}

	return &window{s: s, w: w, b: buf}, err
}

func (w *window) Release() {
	if w.b != nil {
		w.b.Release()
		w.b = nil
	}
	if w.w != nil {
		w.w.Release()
		w.w = nil
	}
	w.s = nil
}

func (w *window) display(img image.Image) screen.PublishResult {
	rect := image.Rect(0, 0, xmax, ymax)
	sr := img.Bounds()

	w.w.Fill(rect, bkgCol, draw.Src)
	draw.Draw(w.b.RGBA(), w.b.Bounds(), img, image.Point{}, draw.Src)
	if !sr.In(rect) {
		sr = rect
	}
	w.w.Upload(image.Point{}, w.b, sr)

	o := w.w.Publish()
	return o
}
