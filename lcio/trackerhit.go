// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"go-hep.org/x/hep/sio"
)

type TrackerHits struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHit
}

func (*TrackerHits) VersionSio() uint32 {
	return Version
}

func (hits *TrackerHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *TrackerHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]TrackerHit, int(n))
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
		var cov [6]float32
		dec.Decode(&cov)
		hit.Cov[0] = float64(cov[0])
		hit.Cov[1] = float64(cov[1])
		hit.Cov[2] = float64(cov[2])
		hit.Cov[3] = float64(cov[3])
		hit.Cov[4] = float64(cov[4])
		hit.Cov[5] = float64(cov[5])

		dec.Decode(&hit.EDep)
		if r.VersionSio() > 1012 {
			dec.Decode(&hit.EDepErr)
		}
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

type TrackerHit struct {
	CellID0 int32
	CellID1 int32
	Type    int32 // type of Track; encoded in parameters TrackerHitTypeName+TrackerHit TypeValue
	Pos     [3]float64
	Cov     [6]float64 // covariance matrix of position (x,y,z)
	EDep    float32    // energy deposit on the hit
	EDepErr float32    // error measured on EDep
	Time    float32
	Quality int32 // quality flag word
	RawHits []*RawCalorimeterHit
}

var _ sio.Codec = (*TrackerHits)(nil)
