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

type RawCalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []RawCalorimeterHit
}

func (hits RawCalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of RawCalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::RCHBIT_ID1    : %v\n", hits.Flags.Test(RChBitID1))
	fmt.Fprintf(o, "     LCIO::RCHBIT_TIME   : %v\n", hits.Flags.Test(RChBitTime))
	fmt.Fprintf(o, "     LCIO::RCHBIT_NO_PTR : %v\n", hits.Flags.Test(RChBitNoPtr))

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

func (*RawCalorimeterHits) VersionSio() uint32 {
	return Version
}

func (hits *RawCalorimeterHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *RawCalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]RawCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() == 8 || hits.Flags.Test(RChBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Amplitude)
		if hits.Flags.Test(RChBitTime) {
			dec.Decode(&hit.TimeStamp)
		}
		if !hits.Flags.Test(RChBitNoPtr) {
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

var _ sio.Codec = (*RawCalorimeterHits)(nil)
