package hplot

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
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
var New = plot.New

type Values struct {
	plotter.Values
}

// EOF
