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

// SimTrackerHitContainer is a collection of simulated tracker hits.
type SimTrackerHitContainer struct {
	Flags  Flags
	Params Params
	Hits   []SimTrackerHit
}

type SimTrackerHit struct {
	CellID0    int32
	CellID1    int32 // second word for cell ID
	Pos        [3]float64
	EDep       float32 // energy deposited on the hit
	Time       float32
	Mc         *McParticle
	Momentum   [3]float32
	PathLength float32
	Quality    int32
}

func (hits SimTrackerHitContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of SimTrackerHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::THBIT_BARREL   : %v\n", hits.Flags.Test(BitsThBarrel))
	fmt.Fprintf(o, "     LCIO::THBIT_MOMENTUM : %v\n", hits.Flags.Test(BitsThMomentum))

	// FIXME(sbinet): quality-bits

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |cellId0 |cellId1 |  position (x,y,z)               |   EDep   |   time   |  PDG  |        (px,  py, pz)          | path-len | Quality \n"
		tail = "------------|--------|--------|---------------------------------|----------|----------|-------|-------------------------------|----------|---------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		var pdg int32 = 0
		if hit.Mc != nil {
			pdg = hit.Mc.PDG
		}
		fmt.Fprintf(o,
			" [%08d] |%08d|%08d|(%+.2e, %+.2e, %+.2e)| %.2e | %.2e | %05d |(%+.2e,%+.2e,%+.2e)|%+.3e|%5d\n",
			0, //id
			hit.CellID0, hit.CellID1,
			hit.Pos[0], hit.Pos[1], hit.Pos[2],
			hit.EDep, hit.Time,
			pdg,
			hit.Momentum[0], hit.Momentum[1], hit.Momentum[2],
			hit.PathLength,
			hit.Quality,
		)
		// FIXME(sbinet): CellIDDecoder
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*SimTrackerHitContainer) VersionSio() uint32 {
	return Version
}

func (hits *SimTrackerHitContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hits.Flags)
	enc.Encode(&hits.Params)
	enc.Encode(int32(len(hits.Hits)))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		enc.Encode(&hit.CellID0)
		if hits.Flags.Test(BitsThID1) {
			enc.Encode(&hit.CellID1)
		}
		enc.Encode(&hit.Pos)
		enc.Encode(&hit.EDep)
		enc.Encode(&hit.Time)
		enc.Pointer(&hit.Mc)
		if hits.Flags.Test(BitsThMomentum) {
			enc.Encode(&hit.Momentum)
			enc.Encode(&hit.PathLength)
		}
		enc.Encode(&hit.Quality)
		enc.Tag(hit)
	}
	return enc.Err()
}

func (hits *SimTrackerHitContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]SimTrackerHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() > 1051 && hits.Flags.Test(BitsThID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Pos)
		dec.Decode(&hit.EDep)
		dec.Decode(&hit.Time)
		dec.Pointer(&hit.Mc)
		if hits.Flags.Test(BitsThMomentum) {
			dec.Decode(&hit.Momentum)
			if r.VersionSio() > 1006 {
				dec.Decode(&hit.PathLength)
			}
		}
		if r.VersionSio() > 2007 {
			dec.Decode(&hit.Quality)
		}
		if r.VersionSio() > 1000 {
			dec.Tag(hit)
		}
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*SimTrackerHitContainer)(nil)
	_ sio.Codec     = (*SimTrackerHitContainer)(nil)
)
