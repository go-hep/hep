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

// EvenBinAxis is an evenly-binned 1D axis
type EvenBinAxis struct {
	nbins int     // number of bins for this axis
	low   float64 // low-edge of this axis
	high  float64 // high-edge of this axis
	size  float64 // bin size
}

// NewEvenBinAxis returns a new axis with n bins between xmax and xmax.
// It panics if n is <= 0 or if xmin >= xmax.
func NewEvenBinAxis(n int, xmin, xmax float64) *EvenBinAxis {
	if n <= 0 {
		panic("hbook: X-axis with zero bins")
	}
	if xmin >= xmax {
		panic("hbook: invalid X-axis limits")
	}
	axis := &EvenBinAxis{
		nbins: n,
		low:   xmin,
		high:  xmax,
		size:  (xmax - xmin) / float64(n),
	}
	return axis
}

// Kind returns the binning kind (Fixed,Variable) of an axis
func (axis *EvenBinAxis) Kind() AxisKind {
	return FixedBinning
}

// LowerEdge returns the lower edge of the axis.
func (axis *EvenBinAxis) LowerEdge() float64 {
	return axis.low
}

// UpperEdge returns the upper edge of the axis.
func (axis *EvenBinAxis) UpperEdge() float64 {
	return axis.high
}

// Bins returns the number of bins in the axis.
func (axis *EvenBinAxis) Bins() int {
	return axis.nbins
}

// BinLowerEdge returns the lower edge of the bin at index i.
// It panics if i is outside the axis range.
func (axis *EvenBinAxis) BinLowerEdge(i int) float64 {
	if i >= 0 && i <= axis.nbins {
		return axis.low + float64(i)*axis.size
	}
	if i == UnderflowBin {
		return axis.low
	}
	panic(fmt.Errorf("hbook: out of bound index (%d)", i))
}

// BinUpperEdge returns the upper edge of the bin at index i.
// It panics if i is outside the axis range.
func (axis *EvenBinAxis) BinUpperEdge(i int) float64 {
	if i >= 0 && i < axis.nbins {
		return axis.low + float64(i+1)*axis.size
	}
	if i == OverflowBin {
		return axis.high
	}
	panic(fmt.Errorf("hbook: out of bound index (%d)", i))
}

// BinWidth returns the width of the bin at index i.
func (axis *EvenBinAxis) BinWidth(i int) float64 {
	return axis.size
}

// CoordToIndex returns the bin index corresponding to the coordinate x.
func (axis *EvenBinAxis) CoordToIndex(x float64) int {
	switch {
	case x < axis.low:
		return UnderflowBin
	case x >= axis.high:
		return OverflowBin
	default:
		return int(math.Floor((x - axis.low) / float64(axis.size)))
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (axis *EvenBinAxis) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := axis.RioMarshal(buf)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (axis *EvenBinAxis) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	err := axis.RioUnmarshal(buf)
	return err
}

// RioVersion implements rio.RioStreamer
func (axis *EvenBinAxis) RioVersion() rio.Version {
	return 0
}

// RioMarshal implements rio.RioMarshaler
func (axis *EvenBinAxis) RioMarshal(w io.Writer) error {
	var err error

	enc := binary.NewEncoder(w)
	err = enc.Encode(axis.nbins)
	if err != nil {
		return err
	}

	err = enc.Encode(axis.low)
	if err != nil {
		return err
	}

	err = enc.Encode(axis.high)
	if err != nil {
		return err
	}

	err = enc.Encode(axis.size)
	if err != nil {
		return err
	}

	return err
}

// RioUnmarshal implements rio.RioUnmarshaler
func (axis *EvenBinAxis) RioUnmarshal(r io.Reader) error {
	var err error

	dec := binary.NewDecoder(r)
	err = dec.Decode(&axis.nbins)
	if err != nil {
		return err
	}

	err = dec.Decode(&axis.low)
	if err != nil {
		return err
	}

	err = dec.Decode(&axis.high)
	if err != nil {
		return err
	}

	err = dec.Decode(&axis.size)
	if err != nil {
		return err
	}

	return err
}

// check EvenBinAxis satisfies Axis interface
var _ Axis = (*EvenBinAxis)(nil)

// serialization interfaces
var _ rio.Marshaler = (*EvenBinAxis)(nil)
var _ rio.Unmarshaler = (*EvenBinAxis)(nil)
var _ rio.Streamer = (*EvenBinAxis)(nil)

func init() {
	gob.Register((*EvenBinAxis)(nil))
}

// EOF
