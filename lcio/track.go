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

// TrackContainer is a collection of tracks.
type TrackContainer struct {
	Flags  Flags
	Params Params
	Tracks []Track
}

type Track struct {
	Type       int32 // type of track (e.g TPC, VTX, SIT)
	States     []TrackState
	Chi2       float32  // chi^2 of fit
	NdF        int32    // ndf of fit
	DEdx       float32  // dEdx
	DEdxErr    float32  // error of dEdx
	Radius     float32  // radius of innermost hit used in track fit
	SubDetHits []int32  // number of hits in particular sub-detectors
	Tracks     []*Track // tracks that have been combined into this track
	Hits       []*TrackerHit
}

func (trk *Track) D0() float64 {
	if len(trk.States) <= 0 {
		return 0
	}
	return float64(trk.States[0].D0)
}

func (trk *Track) Phi() float64 {
	if len(trk.States) <= 0 {
		return 0
	}
	return float64(trk.States[0].Phi)
}

func (trk *Track) Omega() float64 {
	if len(trk.States) <= 0 {
		return 0
	}
	return float64(trk.States[0].Omega)
}

func (trk *Track) Z0() float64 {
	if len(trk.States) <= 0 {
		return 0
	}
	return float64(trk.States[0].Z0)
}

func (trk *Track) TanL() float64 {
	if len(trk.States) <= 0 {
		return 0
	}
	return float64(trk.States[0].TanL)
}

type TrackState struct {
	Loc   int32       // location of the track state
	D0    float32     // impact parameter in r-phi
	Phi   float32     // phi of track in r-phi
	Omega float32     // curvature signed with charge
	Z0    float32     // impact parameter in r-z
	TanL  float32     // tangent of dip angle in r-z
	Cov   [15]float32 // covariance matrix
	Ref   [3]float32  // reference point (x,y,z)
}

func (trks *TrackContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of Track collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", trks.Flags, trks.Params)
	fmt.Fprintf(o, "     LCIO::TRBIT_HITS : %v\n", trks.Flags.Test(BitsClHits))

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |   type   |    d0    |  phi     | omega    |    z0     | tan lambda|   reference point(x,y,z)        |    dEdx  |  dEdxErr |   chi2   |  ndf   \n"
		tail = "------------|----------|----------|----------|----------|-----------|-----------|---------------------------------|----------|----------|-------- \n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range trks.Tracks {
		trk := &trks.Tracks[i]
		var ref [3]float32
		if len(trk.States) > 0 {
			ref = trk.States[0].Ref
		}
		fmt.Fprintf(o,
			"[%09d] | %08d |%+.2e |%+.2e |%+.2e |%+.3e |%+.3e |(%+.2e, %+.2e, %+.2e)|%+.2e |%+.2e |%+.2e |%5d\n",
			ID(trk),
			trk.Type, trk.D0(), trk.Phi(), trk.Omega(), trk.Z0(), trk.TanL(),
			ref[0], ref[1], ref[2],
			trk.DEdx, trk.DEdxErr, trk.Chi2, trk.NdF,
		)
	}

	return string(o.Bytes())
}

func (*TrackContainer) VersionSio() uint32 {
	return Version
}

func (trks *TrackContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&trks.Flags)
	enc.Encode(&trks.Params)
	enc.Encode(int32(len(trks.Tracks)))
	for i := range trks.Tracks {
		trk := &trks.Tracks[i]
		enc.Encode(&trk.Type)
		enc.Encode(int32(len(trk.States)))
		for i := range trk.States {
			state := &trk.States[i]
			enc.Encode(&state.Loc)
			enc.Encode(&state.D0)
			enc.Encode(&state.Phi)
			enc.Encode(&state.Omega)
			enc.Encode(&state.Z0)
			enc.Encode(&state.TanL)
			enc.Encode(&state.Cov)
			enc.Encode(&state.Ref)
		}
		enc.Encode(&trk.Chi2)
		enc.Encode(&trk.NdF)
		enc.Encode(&trk.DEdx)
		enc.Encode(&trk.DEdxErr)
		enc.Encode(&trk.Radius)
		enc.Encode(&trk.SubDetHits)

		enc.Encode(int32(len(trk.Tracks)))
		for i := range trk.Tracks {
			enc.Pointer(&trk.Tracks[i])
		}

		if trks.Flags.Test(BitsTrHits) {
			enc.Encode(int32(len(trk.Hits)))
			for i := range trk.Hits {
				enc.Pointer(&trk.Hits[i])
			}
		}
		enc.Tag(trk)
	}
	return enc.Err()
}

func (trks *TrackContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&trks.Flags)
	dec.Decode(&trks.Params)
	var n int32
	dec.Decode(&n)
	trks.Tracks = make([]Track, int(n))
	for i := range trks.Tracks {
		trk := &trks.Tracks[i]
		dec.Decode(&trk.Type)
		var n int32 = 1 // set to 1 by default for bwd compat
		if r.VersionSio() >= 2000 {
			dec.Decode(&n)
		}
		trk.States = make([]TrackState, int(n))
		for i := range trk.States {
			state := &trk.States[i]
			if r.VersionSio() >= 2000 {
				dec.Decode(&state.Loc)
			}
			dec.Decode(&state.D0)
			dec.Decode(&state.Phi)
			dec.Decode(&state.Omega)
			dec.Decode(&state.Z0)
			dec.Decode(&state.TanL)
			dec.Decode(&state.Cov)
			dec.Decode(&state.Ref)
		}
		dec.Decode(&trk.Chi2)
		dec.Decode(&trk.NdF)
		dec.Decode(&trk.DEdx)
		dec.Decode(&trk.DEdxErr)
		dec.Decode(&trk.Radius)
		dec.Decode(&trk.SubDetHits)

		dec.Decode(&n)
		trk.Tracks = make([]*Track, int(n))
		for i := range trk.Tracks {
			dec.Pointer(&trk.Tracks[i])
		}

		if trks.Flags.Test(BitsTrHits) {
			dec.Decode(&n)
			trk.Hits = make([]*TrackerHit, int(n))
			for i := range trk.Hits {
				dec.Pointer(&trk.Hits[i])
			}
		}
		dec.Tag(trk)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackContainer)(nil)
	_ sio.Codec     = (*TrackContainer)(nil)
)
