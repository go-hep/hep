package hplot

import (
	"github.com/go-hep/hplot/plotinum/plotter"
)

// NewFunction returns a Function that plots F using
// the default line style with 50 samples.
// NewFunction returns a Function that plots F using
// the default line style with 50 samples.
func NewFunction(f func(float64) float64) *plotter.Function {
	return plotter.NewFunction(f)
}

// NewLine returns a Line that uses the default line style and
// does not draw glyphs.
func NewLine(xys plotter.XYer) (*plotter.Line, error) {
	return plotter.NewLine(xys)
}

// NewScatter returns a Scatter that uses the
// default glyph style.
func NewScatter(xys plotter.XYer) (*plotter.Scatter, error) {
	return plotter.NewScatter(xys)
}

// NewGrid returns a new grid with both vertical and
// horizontal lines using the default grid line style.
func NewGrid() *plotter.Grid {
	return plotter.NewGrid()
}

// EOF
