// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"go-hep.org/x/hep/sio"
)

type SimTrackerHits struct {
	Flags  Flags
	Params Params
	Hits   []SimTrackerHit
}

func (*SimTrackerHits) VersionSio() uint32 {
	return Version
}

func (hits *SimTrackerHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *SimTrackerHits) UnmarshalSio(r sio.Reader) error {
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

var _ sio.Codec = (*SimTrackerHits)(nil)
