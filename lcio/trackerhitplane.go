// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"go-hep.org/x/hep/sio"
)

type TrackerHitPlanes struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHitPlane
}

func (*TrackerHitPlanes) VersionSio() uint32 {
	return Version
}

func (hits *TrackerHitPlanes) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *TrackerHitPlanes) UnmarshalSio(r sio.Reader) error {
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
			if hits.Flags.Test(ThBitID1) {
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

var _ sio.Codec = (*TrackerHitPlanes)(nil)
