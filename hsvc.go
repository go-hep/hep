package fwk

import (
	"github.com/go-hep/hbook"
)

// HID is a histogram identifier
type HID string

// H1D wraps a hbook.H1D for safe concurrent access
type H1D struct {
	ID   HID // unique id
	Hist *hbook.H1D
}

// HistSvc is the interface providing access to histograms
type HistSvc interface {
	Svc

	// BookH1D books a 1D histogram.
	// name should be of the form: "/fwk/streams/<stream-name>/<path>/<histogram-name>"
	BookH1D(name string, nbins int, low, high float64) (H1D, error)

	// FillH1D fills the 1D-histogram id with data x and weight w.
	FillH1D(id HID, x, w float64)
}
