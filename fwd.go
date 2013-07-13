package hplot

import (
	"github.com/go-hep/hplot/plotinum/plotter"
	"github.com/go-hep/hplot/plotinum/plot"
)

// NewFunction returns a Function that plots F using
// the default line style with 50 samples.
var NewFunction = plotter.NewFunction

// NewScatter returns a Scatter that uses the
// default glyph style.
var NewScatter = plotter.NewScatter

// NewGrid returns a new grid with both vertical and
// horizontal lines using the default grid line style.
var NewGrid = plotter.NewGrid

// New returns a new plot with some reasonable
// default settings.
func NewPlot() (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.X.Padding = 0
	p.Y.Padding = 0
	p.Style = plot.GnuplotStyle
	return p, err
}


type Values struct {
	plotter.Values
}

// EOF
