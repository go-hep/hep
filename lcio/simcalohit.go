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

type SimCalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []SimCalorimeterHit
}

func (hits SimCalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of SimCalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "  -> LCIO::CHBIT_LONG   : %v\n", hits.Flags.Test(BitsChLong))
	fmt.Fprintf(o, "     LCIO::CHBIT_BARREL : %v\n", hits.Flags.Test(BitsChBarrel))
	fmt.Fprintf(o, "     LCIO::CHBIT_ID1    : %v\n", hits.Flags.Test(BitsChID1))
	fmt.Fprintf(o, "     LCIO::CHBIT_STEP   : %v\n", hits.Flags.Test(BitsChStep))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |cellId0 |cellId1 |  energy  |        position (x,y,z)          | nMCParticles \n" +
		"           -> MC contribution: prim. PDG |  energy  |   time   | sec. PDG | stepPosition (x,y,z) \n"
	tail := "------------|--------|--------|----------|----------------------------------|--------------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d|%08d|%+.3e|", 0, hit.CellID0, hit.CellID1, hit.Energy)
		if hits.Flags.Test(BitsChLong) {
			fmt.Fprintf(o, "+%.3e, %+.3e, %+.3e", hit.Pos[0], hit.Pos[1], hit.Pos[2])
		} else {
			fmt.Fprintf(o, "    no position available         ")
		}
		fmt.Fprintf(o, "|%+12d\n", len(hit.Contributions))
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "        id-fields: --- unknown/default ----   ")
		for _, c := range hit.Contributions {
			var pdg int32
			if c.Mc != nil {
				pdg = c.Mc.PDG
			}
			fmt.Fprintf(o, "\n           ->                  %+10d|%+1.3e|%+1.3e|", pdg, c.Energy, c.Time)
			if hits.Flags.Test(BitsChStep) {
				fmt.Fprintf(o, "%+d| (%+1.3e, %+1.3e, %+1.3e)", c.PDG, c.StepPos[0], c.StepPos[1], c.StepPos[2])
			} else {
				fmt.Fprintf(o, " no PDG")
			}
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*SimCalorimeterHits) VersionSio() uint32 {
	return Version
}

func (hits *SimCalorimeterHits) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hits.Flags)
	enc.Encode(&hits.Params)
	enc.Encode(int32(len(hits.Hits)))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		enc.Encode(&hit.CellID0)
		if hits.Flags.Test(BitsChID1) {
			enc.Encode(&hit.CellID1)
		}
		enc.Encode(&hit.Energy)
		if hits.Flags.Test(BitsChLong) {
			enc.Encode(&hit.Pos)
		}
		enc.Encode(int32(len(hit.Contributions)))
		for i := range hit.Contributions {
			c := &hit.Contributions[i]
			enc.Pointer(&c.Mc)
			enc.Encode(&c.Energy)
			enc.Encode(&c.Time)
			if hits.Flags.Test(BitsChStep) {
				enc.Encode(&c.PDG)
				enc.Encode(&c.StepPos)
			}
		}
		enc.Tag(hit)
	}
	return enc.Err()
}

func (hits *SimCalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]SimCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() < 9 || hits.Flags.Test(BitsChID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Energy)
		if hits.Flags.Test(BitsChLong) {
			dec.Decode(&hit.Pos)
		}
		var n int32
		dec.Decode(&n)
		hit.Contributions = make([]Contrib, int(n))
		for i := range hit.Contributions {
			c := &hit.Contributions[i]
			dec.Pointer(&c.Mc)
			dec.Decode(&c.Energy)
			dec.Decode(&c.Time)
			if hits.Flags.Test(BitsChStep) {
				dec.Decode(&c.PDG)
				if r.VersionSio() > 1051 {
					dec.Decode(&c.StepPos)
				}
			}
		}
		if r.VersionSio() > 1000 {
			dec.Tag(hit)
		}
	}
	return dec.Err()
}

type SimCalorimeterHit struct {
	Params        Params
	CellID0       int32
	CellID1       int32
	Energy        float32
	Pos           [3]float32
	Contributions []Contrib
}

type Contrib struct {
	Mc      *McParticle
	Energy  float32
	Time    float32
	PDG     int32
	StepPos [3]float32
}

var _ sio.Codec = (*SimCalorimeterHits)(nil)
