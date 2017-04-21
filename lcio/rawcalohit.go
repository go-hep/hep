// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

// RawCalorimeterHitContainer is a collection of raw calorimeter hits.
type RawCalorimeterHitContainer struct {
	Flags  Flags
	Params Params
	Hits   []RawCalorimeterHit
}

func (hits RawCalorimeterHitContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of RawCalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::RCHBIT_ID1    : %v\n", hits.Flags.Test(BitsRChID1))
	fmt.Fprintf(o, "     LCIO::RCHBIT_TIME   : %v\n", hits.Flags.Test(BitsRChTime))
	fmt.Fprintf(o, "     LCIO::RCHBIT_NO_PTR : %v\n", hits.Flags.Test(BitsRChNoPtr))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |  cellId0 ( M, S, I, J, K) |cellId1 | amplitude |  time  \n"
	tail := "------------|---------------------------|--------|-----------|---------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d%19s|%08d|%10d |%8d", 0, hit.CellID0, "", hit.CellID1, hit.Amplitude, hit.TimeStamp)
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "\n        id-fields: --- unknown/default ----   ")
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*RawCalorimeterHitContainer) VersionSio() uint32 {
	return Version
}

func (hits *RawCalorimeterHitContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hits.Flags)
	enc.Encode(&hits.Params)
	enc.Encode(int32(len(hits.Hits)))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		enc.Encode(&hit.CellID0)
		if hits.Flags.Test(BitsRChID1) {
			enc.Encode(&hit.CellID1)
		}
		enc.Encode(&hit.Amplitude)
		if hits.Flags.Test(BitsRChTime) {
			enc.Encode(&hit.TimeStamp)
		}
		if !hits.Flags.Test(BitsRChNoPtr) {
			enc.Tag(hit)
		}
	}
	return enc.Err()
}

func (hits *RawCalorimeterHitContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]RawCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() == 8 || hits.Flags.Test(BitsRChID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Amplitude)
		if hits.Flags.Test(BitsRChTime) {
			dec.Decode(&hit.TimeStamp)
		}
		if !hits.Flags.Test(BitsRChNoPtr) {
			dec.Tag(hit)
		}
	}
	return dec.Err()
}

type RawCalorimeterHit struct {
	CellID0   int32
	CellID1   int32
	Amplitude int32
	TimeStamp int32
}

var (
	_ sio.Versioner = (*RawCalorimeterHitContainer)(nil)
	_ sio.Codec     = (*RawCalorimeterHitContainer)(nil)
)
