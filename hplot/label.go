package hplot

import (
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type Label struct {
	Text  string         // Text of the label
	X, Y  float64        // Position of the label
	Style draw.TextStyle // Style of the label

	// Compute the position wrt canvas size. This
	// feature would need to be implemented into the plotter.
	// (There is an example in gonum/v1/plot/plotter).
	CanPos bool
}

// NewLabel creates a new value of type Label.
func NewLabel(x, y float64, txt string, opts ...LabelOption) Label {

	// Handle configuration
	cfg := &LabelConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Create the basic label
	label := Label{Text: txt, X: x, Y: y}

	// Apply configuration
	label.Style = cfg.Style
	label.CanPos = cfg.CanPos

	// Return the customized object
	return label
}

// This function could be internal only to convert
// the a slice of labels into a plotter.Labels object?
func PlotterLabels(Ls []Label) *plotter.Labels {

	// Wrap up all the labels into a field slices.
	xys := make([]plotter.XY, len(Ls))
	txt := make([]string, len(Ls))
	stl := make([]draw.TextStyle, len(Ls))
	for i, l := range Ls {
		xys[i] = plotter.XY{X: l.X, Y: l.Y}
		txt[i] = l.Text
		stl[i] = l.Style
	}

	// Create of the YXlabels.
	xyL := plotter.XYLabels{
		XYs:    xys,
		Labels: txt,
	}

	// Create the plotter.Labels
	labels, err := plotter.NewLabels(xyL)
	if err != nil {
		panic("cannot create plotter.Labels")
	}

	// Add the text styles.
	// FIXME[rmadar]: uncomment this line when l.Style will
	//                be set to a default value.
	// labels.TextStyle = stl

	// Return the result
	return labels
}

type labelConfig struct {
	Style  draw.TextStyle
	CanPos bool
}

// Label option
type LabelOption func(cfg *labelConfig)

// withTextStyle specifies the text style of the label.
func withTextStyle(style draw.TextStyle) LabelOption {
	return func(cfg *labelConfig) {
		cfg.Style = style
	}
}

// withTextStyle specifies the text style of the label.
func withCanPos(doIt bool) LabelOption {
	return func(cfg *labelConfig) {
		cfg.CanPos = doIt
	}
}
