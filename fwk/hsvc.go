package fwk

import (
	"go-hep.org/x/hep/hbook"
)

// Hist is a histogram, scatter or profile object that can
// be saved or loaded by the HistSvc.
type Hist interface {
	Name() string
	Value() interface{}
}

// HID is a histogram, scatter or profile identifier
type HID string

// H1D wraps a hbook.H1D for safe concurrent access
type H1D struct {
	ID   HID // unique id
	Hist *hbook.H1D
}

func (h H1D) Name() string {
	return string(h.ID)
}

func (h H1D) Value() interface{} {
	return h.Hist
}

// H2D wraps a hbook.H2D for safe concurrent access
type H2D struct {
	ID   HID // unique id
	Hist *hbook.H2D
}

func (h H2D) Name() string {
	return string(h.ID)
}

func (h H2D) Value() interface{} {
	return h.Hist
}

// P1D wraps a hbook.P1D for safe concurrent access
type P1D struct {
	ID      HID // unique id
	Profile *hbook.P1D
}

func (p P1D) Name() string {
	return string(p.ID)
}

func (p P1D) Value() interface{} {
	return p.Profile
}

// S2D wraps a hbook.S2D for safe concurrent access
type S2D struct {
	ID      HID // unique id
	Scatter *hbook.S2D
}

func (s S2D) Name() string {
	return string(s.ID)
}

func (s S2D) Value() interface{} {
	return s.Scatter
}

// HistSvc is the interface providing access to histograms
type HistSvc interface {
	Svc

	// BookH1D books a 1D histogram.
	// name should be of the form: "/fwk/streams/<stream-name>/<path>/<histogram-name>"
	BookH1D(name string, nbins int, xmin, xmax float64) (H1D, error)

	// BookH2D books a 2D histogram.
	// name should be of the form: "/fwk/streams/<stream-name>/<path>/<histogram-name>"
	BookH2D(name string, nx int, xmin, xmax float64, ny int, ymin, ymax float64) (H2D, error)

	// BookP1D books a 1D profile.
	// name should be of the form: "/fwk/streams/<stream-name>/<path>/<profile-name>"
	BookP1D(name string, nbins int, xmin, xmax float64) (P1D, error)

	// BookS2D books a 2D scatter.
	// name should be of the form: "/fwk/streams/<stream-name>/<path>/<scatter-name>"
	BookS2D(name string) (S2D, error)

	// FillH1D fills the 1D-histogram id with data x and weight w.
	FillH1D(id HID, x, w float64)

	// FillH2D fills the 2D-histogram id with data (x,y) and weight w.
	FillH2D(id HID, x, y, w float64)

	// FillP1D fills the 1D-profile id with data (x,y) and weight w.
	FillP1D(id HID, x, y, w float64)

	// FillS2D fills the 2D-scatter id with data (x,y).
	FillS2D(id HID, x, y float64)
}

var _ Hist = (*H1D)(nil)
var _ Hist = (*H2D)(nil)
var _ Hist = (*P1D)(nil)
var _ Hist = (*S2D)(nil)
