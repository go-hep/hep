package dao

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
)

type EvenBinAxis struct {
	nbins int     // number of bins for this axis
	low   float64 // low-edge of this axis
	high  float64 // high-edge of this axis
	size  float64 // bin size
}

func NewEvenBinAxis(nbins int, low, high float64) *EvenBinAxis {
	axis := &EvenBinAxis{
		nbins: nbins,
		low:   low,
		high:  high,
		size:  (high - low) / float64(nbins),
	}
	return axis
}

func (axis *EvenBinAxis) Kind() AxisKind {
	return FixedBinning
}

func (axis *EvenBinAxis) LowerEdge() float64 {
	return axis.low
}

func (axis *EvenBinAxis) UpperEdge() float64 {
	return axis.high
}

func (axis *EvenBinAxis) Bins() int {
	return axis.nbins
}

func (axis *EvenBinAxis) BinLowerEdge(idx int) float64 {
	if idx >= 0 && idx <= axis.nbins {
		return axis.low + float64(idx)*axis.size
	}
	if idx == UnderflowBin {
		return axis.low
	}
	panic(fmt.Errorf("hist: out of bound index (%d)", idx))
}

func (axis *EvenBinAxis) BinUpperEdge(idx int) float64 {
	if idx >= 0 && idx < axis.nbins {
		return axis.low + float64(idx+1)*axis.size
	}
	if idx == OverflowBin {
		return axis.high
	}
	panic(fmt.Errorf("hist: out of bound index (%d)", idx))
}

func (axis *EvenBinAxis) BinWidth(idx int) float64 {
	return axis.size
}

func (axis *EvenBinAxis) CoordToIndex(coord float64) int {
	switch {
	case coord < axis.low:
		return UnderflowBin
	case coord >= axis.high:
		return OverflowBin
	default:
		return int(math.Floor((coord - axis.low) / float64(axis.size)))
	}
	panic("unreachable")
}

func (axis *EvenBinAxis) MarshalBinary(buf *bytes.Buffer) error {
	enc := gob.NewEncoder(buf)
	err := axis.gobEncode(enc)
	return err
}

func (axis *EvenBinAxis) UnmarshalBinary(buf *bytes.Buffer) error {
	dec := gob.NewDecoder(buf)
	err := axis.gobDecode(dec)
	return err
}

func (axis *EvenBinAxis) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := axis.gobEncode(enc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (axis *EvenBinAxis) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := axis.gobDecode(dec)
	return err
}

func (axis *EvenBinAxis) gobEncode(enc *gob.Encoder) error {
	var err error

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

func (axis *EvenBinAxis) gobDecode(dec *gob.Decoder) error {
	var err error

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

func init() {
	gob.Register((*EvenBinAxis)(nil))
}

// EOF
