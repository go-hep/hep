// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"

	"github.com/go-hep/rio"
	"github.com/gonuts/binary"
)

// evenBinning is an evenly-binned 1D axis
type evenBinning struct {
	nbins int     // number of bins for this binning
	low   float64 // low-edge of this binning
	high  float64 // high-edge of this binning
	size  float64 // bin size
}

// newEvenBinning returns a new binning with n bins between xmax and xmax.
// It panics if n is <= 0 or if xmin >= xmax.
func newEvenBinning(n int, xmin, xmax float64) *evenBinning {
	if n <= 0 {
		panic("hbook: X-axis with zero bins")
	}
	if xmin >= xmax {
		panic("hbook: invalid X-axis limits")
	}
	bng := &evenBinning{
		nbins: n,
		low:   xmin,
		high:  xmax,
		size:  (xmax - xmin) / float64(n),
	}
	return bng
}

// Kind returns the binning kind (Fixed,Variable)
func (bng *evenBinning) Kind() BinningKind {
	return FixedBinning
}

// LowerEdge returns the lower edge of the binning.
func (bng *evenBinning) LowerEdge() float64 {
	return bng.low
}

// UpperEdge returns the upper edge of the binning.
func (bng *evenBinning) UpperEdge() float64 {
	return bng.high
}

// Bins returns the number of bins in the binning.
func (bng *evenBinning) Bins() int {
	return bng.nbins
}

// BinLowerEdge returns the lower edge of the bin at index i.
// It panics if i is outside the binning range.
func (bng *evenBinning) BinLowerEdge(i int) float64 {
	if i >= 0 && i <= bng.nbins {
		return bng.low + float64(i)*bng.size
	}
	if i == UnderflowBin {
		return bng.low
	}
	panic(fmt.Errorf("hbook: out of bound index (%d)", i))
}

// BinUpperEdge returns the upper edge of the bin at index i.
// It panics if i is outside the binning range.
func (bng *evenBinning) BinUpperEdge(i int) float64 {
	if i >= 0 && i < bng.nbins {
		return bng.low + float64(i+1)*bng.size
	}
	if i == OverflowBin {
		return bng.high
	}
	panic(fmt.Errorf("hbook: out of bound index (%d)", i))
}

// BinWidth returns the width of the bin at index i.
func (bng *evenBinning) BinWidth(i int) float64 {
	return bng.size
}

// CoordToIndex returns the bin index corresponding to the coordinate x.
func (bng *evenBinning) CoordToIndex(x float64) int {
	switch {
	case x < bng.low:
		return UnderflowBin
	case x >= bng.high:
		return OverflowBin
	default:
		return int(math.Floor((x - bng.low) / float64(bng.size)))
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (bng *evenBinning) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := bng.RioMarshal(buf)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (bng *evenBinning) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	err := bng.RioUnmarshal(buf)
	return err
}

// RioVersion implements rio.RioStreamer
func (bng *evenBinning) RioVersion() rio.Version {
	return 0
}

// RioMarshal implements rio.RioMarshaler
func (bng *evenBinning) RioMarshal(w io.Writer) error {
	var err error

	enc := binary.NewEncoder(w)
	err = enc.Encode(bng.nbins)
	if err != nil {
		return err
	}

	err = enc.Encode(bng.low)
	if err != nil {
		return err
	}

	err = enc.Encode(bng.high)
	if err != nil {
		return err
	}

	err = enc.Encode(bng.size)
	if err != nil {
		return err
	}

	return err
}

// RioUnmarshal implements rio.RioUnmarshaler
func (bng *evenBinning) RioUnmarshal(r io.Reader) error {
	var err error

	dec := binary.NewDecoder(r)
	err = dec.Decode(&bng.nbins)
	if err != nil {
		return err
	}

	err = dec.Decode(&bng.low)
	if err != nil {
		return err
	}

	err = dec.Decode(&bng.high)
	if err != nil {
		return err
	}

	err = dec.Decode(&bng.size)
	if err != nil {
		return err
	}

	return err
}

// check evenBinning satisfies Binning interface
var _ Binning = (*evenBinning)(nil)

// serialization interfaces
var _ rio.Marshaler = (*evenBinning)(nil)
var _ rio.Unmarshaler = (*evenBinning)(nil)
var _ rio.Streamer = (*evenBinning)(nil)

func init() {
	gob.Register((*evenBinning)(nil))
}
