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

// TrackerHitContainer is a collection of tracker hits.
type TrackerHitContainer struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHit
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

func (hits TrackerHitContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of TrackerHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::THBIT_BARREL   : %v\n", hits.Flags.Test(BitsThBarrel))

	// FIXME(sbinet): quality-bits

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |cellId0 |cellId1 | position (x,y,z)            | time    |[type]|[qual]| EDep    |EDepError|  cov(x,x),  cov(y,x),  cov(y,y),  cov(z,x),  cov(z,y),  cov(z,z)\n"
		tail = "------------|--------|--------|-----------------------------|---------|------|------|---------|---------|-----------------------------------------------------------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		fmt.Fprintf(o,
			" [%08d] |%08d|%08d|%+.2e,%+.2e,%+.2e|%+.2e|[%04d]|[%04d]|%+.2e|%+.2e|%+.2e, %+.2e, %+.2e, %+.2e, %+.2e, %+.2e\n",
			0, // id
			hit.CellID0, hit.CellID1,
			hit.Pos[0], hit.Pos[1], hit.Pos[2],
			hit.Time, hit.Type, hit.Quality,
			hit.EDep, hit.EDepErr,
			hit.Cov[0], hit.Cov[1], hit.Cov[2], hit.Cov[3], hit.Cov[4], hit.Cov[5],
		)
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*TrackerHitContainer) VersionSio() uint32 {
	return Version
}

func (hits *TrackerHitContainer) MarshalSio(w sio.Writer) error {
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
		var cov [6]float32
		cov[0] = float32(hit.Cov[0])
		cov[1] = float32(hit.Cov[1])
		cov[2] = float32(hit.Cov[2])
		cov[3] = float32(hit.Cov[3])
		cov[4] = float32(hit.Cov[4])
		cov[5] = float32(hit.Cov[5])
		enc.Encode(&cov)
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

func (hits *TrackerHitContainer) UnmarshalSio(r sio.Reader) error {
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

var (
	_ sio.Versioner = (*TrackerHitContainer)(nil)
	_ sio.Codec     = (*TrackerHitContainer)(nil)
)
