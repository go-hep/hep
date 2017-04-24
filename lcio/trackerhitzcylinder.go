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

// TrackerHitZCylinderContainer is a collection of tracker hit z-cylinders.
type TrackerHitZCylinderContainer struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHitZCylinder
}

type TrackerHitZCylinder struct {
	CellID0 int32
	CellID1 int32
	Type    int32 // type of Track; encoded in parameters TrackerHitTypeName+TrackerHit TypeValue
	Pos     [3]float64
	Center  [2]float32
	DRPhi   float32 // measurement error along RPhi
	DZ      float32 // measurement error along z
	EDep    float32 // energy deposit on the hit
	EDepErr float32 // error measured on EDep
	Time    float32
	Quality int32 // quality flag word
	RawHits []Hit
}

func (hit *TrackerHitZCylinder) GetCellID0() int32 { return hit.CellID0 }
func (hit *TrackerHitZCylinder) GetCellID1() int32 { return hit.CellID1 }

func (hits TrackerHitZCylinderContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of TrackerHitZCylinder collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::THBIT_BARREL   : %v\n", hits.Flags.Test(BitsThBarrel))

	// FIXME(sbinet): quality-bits

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |cellId0 |cellId1 | position (x,y,z)            | time    |[type]|[qual]| EDep    |EDepError|  dRPhi  |    dZ   |    center (x,y)   |\n"
		tail = "------------|--------|--------|-----------------------------|---------|------|------|---------|---------|---------|---------|-------------------|\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		fmt.Fprintf(o,
			"[%09d] |%08d|%08d|%+.2e,%+.2e,%+.2e|%+.2e|[%04d]|[%04d]|%+.2e|%+.2e|%+.2e|%+.2e|%+.2e,%+.2e|\n",
			ID(hit),
			hit.CellID0, hit.CellID1,
			hit.Pos[0], hit.Pos[1], hit.Pos[2],
			hit.Time, hit.Type, hit.Quality,
			hit.EDep, hit.EDepErr,
			hit.DRPhi, hit.DZ,
			hit.Center[0], hit.Center[1],
		)
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*TrackerHitZCylinderContainer) VersionSio() uint32 {
	return Version
}

func (hits *TrackerHitZCylinderContainer) MarshalSio(w sio.Writer) error {
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
		enc.Encode(&hit.Type)
		enc.Encode(&hit.Pos)
		enc.Encode(&hit.Center)
		enc.Encode(&hit.DRPhi)
		enc.Encode(&hit.DZ)
		enc.Encode(&hit.EDep)
		enc.Encode(&hit.EDepErr)
		enc.Encode(&hit.Time)
		enc.Encode(&hit.Quality)

		enc.Encode(int32(len(hit.RawHits)))
		for ii := range hit.RawHits {
			enc.Pointer(&hit.RawHits[ii])
		}
		enc.Tag(hit)
	}
	return enc.Err()
}

func (hits *TrackerHitZCylinderContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]TrackerHitZCylinder, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		if r.VersionSio() > 1051 {
			dec.Decode(&hit.CellID0)
			if hits.Flags.Test(BitsThID1) {
				dec.Decode(&hit.CellID1)
			}
		}
		if r.VersionSio() > 1002 {
			dec.Decode(&hit.Type)
		}
		dec.Decode(&hit.Pos)
		dec.Decode(&hit.Center)
		dec.Decode(&hit.DRPhi)
		dec.Decode(&hit.DZ)
		dec.Decode(&hit.EDep)
		dec.Decode(&hit.EDepErr)
		dec.Decode(&hit.Time)
		if r.VersionSio() > 1011 {
			dec.Decode(&hit.Quality)
		}

		var n int32 = 1
		if r.VersionSio() > 1002 {
			dec.Decode(&n)
		}
		hit.RawHits = make([]Hit, int(n))
		for ii := range hit.RawHits {
			dec.Pointer(&hit.RawHits[ii])
		}
		dec.Tag(hit)
	}

	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackerHitZCylinderContainer)(nil)
	_ sio.Codec     = (*TrackerHitZCylinderContainer)(nil)
	_ Hit           = (*TrackerHitZCylinder)(nil)
)
