package dao

import (
	"bytes"
	"encoding/gob"
	"math"
)

type Bin1D struct {
	entries int64
	sw      float64 // sum of weights
	swc     float64 // sum of weights times coord
	sw2     float64 // sum of squared weights
}

func (b *Bin1D) Entries() int64 {
	return b.entries
}

func (b *Bin1D) fill(coord, weight float64) {
	b.entries += 1
	b.sw += weight
	b.swc += weight * coord
	b.sw2 += weight * weight
}

func (b *Bin1D) increment(n int64, height, error, centre float64) {
	b.entries += n
	b.sw += height
	b.swc += centre * height
	b.sw2 += error * error
}

func (b *Bin1D) Set(n int64, height, error, centre float64) {
	b.entries = n
	b.sw = height
	b.swc = centre * height
	b.sw2 = error * error
}

func (b *Bin1D) Err() float64 {
	return math.Sqrt(b.sw2)
}

func (b *Bin1D) Err2() float64 {
	return b.sw2
}

func (b *Bin1D) Centre() float64 {
	if b.sw == 0. {
		return 0.
	}
	return b.swc / b.sw
}

func (b *Bin1D) Scale(factor float64) {
	b.sw *= factor
	b.swc *= factor
	b.sw2 *= factor * factor
}

func (b *Bin1D) MarshalBinary(buf *bytes.Buffer) error {
	enc := gob.NewEncoder(buf)
	err := b.gobEncode(enc)
	return err
}

func (b *Bin1D) UnmarshalBinary(buf *bytes.Buffer) error {
	dec := gob.NewDecoder(buf)
	err := b.gobDecode(dec)
	return err
}

func (b *Bin1D) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := b.gobEncode(enc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func (b *Bin1D) gobEncode(enc *gob.Encoder) error {
	var err error

	err = enc.Encode(b.entries)
	if err != nil {
		return err
	}

	err = enc.Encode(b.sw)
	if err != nil {
		return err
	}

	err = enc.Encode(b.swc)
	if err != nil {
		return err
	}

	err = enc.Encode(b.sw2)
	if err != nil {
		return err
	}

	return err
}

func (b *Bin1D) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := b.gobDecode(dec)
	return err
}

func (b *Bin1D) gobDecode(dec *gob.Decoder) error {
	var err error
	err = dec.Decode(&b.entries)
	if err != nil {
		return err
	}

	err = dec.Decode(&b.sw)
	if err != nil {
		return err
	}

	err = dec.Decode(&b.swc)
	if err != nil {
		return err
	}

	err = dec.Decode(&b.sw2)
	if err != nil {
		return err
	}
	return err
}

func init() {
	gob.Register((*Bin1D)(nil))
}

// EOF
