// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

// TrackerHitPlaneContainer is a collection of tracker hit planes.
type TrackerHitPlaneContainer struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHitPlane
}

type TrackerHitPlane struct {
	CellID0 int32
	CellID1 int32
	Type    int32 // type of Track; encoded in parameters TrackerHitTypeName+TrackerHit TypeValue
	Pos     [3]float64
	U       [2]float32
	V       [2]float32
	DU      float32 // measurement error along u
	DV      float32 // measurement error along v
	EDep    float32 // energy deposit on the hit
	EDepErr float32 // error measured on EDep
	Time    float32
	Quality int32 // quality flag word
	RawHits []*RawCalorimeterHit
}

func (hit *TrackerHitPlane) GetCellID0() int32 { return hit.CellID0 }
func (hit *TrackerHitPlane) GetCellID1() int32 { return hit.CellID1 }

func (hits TrackerHitPlaneContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of TrackerHitPlane collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::THBIT_BARREL   : %v\n", hits.Flags.Test(BitsThBarrel))

	// FIXME(sbinet): quality-bits

	fmt.Fprintf(o, "\n")

	dec := NewCellIDDecoderFrom(hits.Params)
	const (
		head = " [   id   ] |cellId0 |cellId1 | position (x,y,z)            | time    |[type]|[qual]| EDep    |EDepError|   du    |   dv    |  u (theta, phi)   |  v (theta, phi)\n"
		tail = "------------|--------|--------|-----------------------------|---------|------|------|---------|---------|---------|---------|-------------------|-------------------|\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		fmt.Fprintf(o,
			"[%09d] |%08d|%08d|%+.2e,%+.2e,%+.2e|%+.2e|[%04d]|[%04d]|%+.2e|%+.2e|%+.2e|%+.2e|%+.2e,%+.2e|%+.2e,%+.2e|\n",
			ID(hit),
			hit.CellID0, hit.CellID1,
			hit.Pos[0], hit.Pos[1], hit.Pos[2],
			hit.Time, hit.Type, hit.Quality,
			hit.EDep, hit.EDepErr,
			hit.DU, hit.DV,
			hit.U[0], hit.U[1],
			hit.V[0], hit.V[1],
		)
		if dec != nil {
			fmt.Fprintf(o, "        id-fields: (%s)\n", dec.ValueString(hit))
		}
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*TrackerHitPlaneContainer) VersionSio() uint32 {
	return Version
}

func (hits *TrackerHitPlaneContainer) MarshalSio(w sio.Writer) error {
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
		enc.Encode(&hit.U)
		enc.Encode(&hit.V)
		enc.Encode(&hit.DU)
		enc.Encode(&hit.DV)
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

func (hits *TrackerHitPlaneContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]TrackerHitPlane, int(n))
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
		dec.Decode(&hit.U)
		dec.Decode(&hit.V)
		dec.Decode(&hit.DU)
		dec.Decode(&hit.DV)
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
		hit.RawHits = make([]*RawCalorimeterHit, int(n))
		for ii := range hit.RawHits {
			dec.Pointer(&hit.RawHits[ii])
		}
		dec.Tag(hit)
	}

	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackerHitPlaneContainer)(nil)
	_ sio.Codec     = (*TrackerHitPlaneContainer)(nil)
	_ Hit           = (*TrackerHitPlane)(nil)
)
