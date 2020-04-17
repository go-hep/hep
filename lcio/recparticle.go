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

// RecParticleContainer is a collection of RecParticles.
type RecParticleContainer struct {
	Flags  Flags
	Params Params
	Parts  []RecParticle
}

type RecParticle struct {
	Type          int32
	P             [3]float32  // momentum (Px,PyPz)
	Energy        float32     // energy of particle
	Cov           [10]float32 // covariance matrix for 4-vector (Px,Py,Pz,E)
	Mass          float32     // mass of object used for 4-vector
	Charge        float32     // charge of particle
	Ref           [3]float32  // reference point of 4-vector
	PIDs          []ParticleID
	PIDUsed       *ParticleID
	GoodnessOfPID float32 // overall quality of the particle identification
	Recs          []*RecParticle
	Tracks        []*Track
	Clusters      []*Cluster
	StartVtx      *Vertex // start vertex associated to the particle
}

type ParticleID struct {
	Likelihood float32
	Type       int32
	PDG        int32
	AlgType    int32
	Params     []float32
}

func (recs *RecParticleContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of ReconstructedParticle collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", recs.Flags, recs.Params)

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |com|type|     momentum( px,py,pz)       | energy | mass   | charge  |        position ( x,y,z)      | pidUsed |GoodnessOfPID|\n"
		tail = "------------|---|----|-------------------------------|--------|--------|---------|-------------------------------|---------|-------------|\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range recs.Parts {
		rec := &recs.Parts[i]
		compound := 0
		if len(rec.Recs) > 0 {
			compound = 1
		}
		fmt.Fprintf(o,
			"[%09d] |%3d|%4d|%+.2e, %+.2e, %+.2e|%.2e|%.2e|%+.2e|%+.2e, %+.2e, %+.2e|%09d|%+.2e| \n",
			ID(rec),
			compound, rec.Type,
			rec.P[0], rec.P[1], rec.P[2], rec.Energy, rec.Mass, rec.Charge,
			rec.Ref[0], rec.Ref[1], rec.Ref[2],
			ID(rec.PIDUsed),
			rec.GoodnessOfPID,
		)
	}
	return string(o.Bytes())
}

func (*RecParticleContainer) VersionSio() uint32 {
	return Version
}

func (recs *RecParticleContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&recs.Flags)
	enc.Encode(&recs.Params)
	enc.Encode(int32(len(recs.Parts)))
	for i := range recs.Parts {
		rec := &recs.Parts[i]
		enc.Encode(&rec.Type)
		enc.Encode(&rec.P)
		enc.Encode(&rec.Energy)
		enc.Encode(&rec.Cov)
		enc.Encode(&rec.Mass)
		enc.Encode(&rec.Charge)
		enc.Encode(&rec.Ref)

		enc.Encode(int32(len(rec.PIDs)))
		for i := range rec.PIDs {
			pid := &rec.PIDs[i]
			enc.Encode(&pid.Likelihood)
			enc.Encode(&pid.Type)
			enc.Encode(&pid.PDG)
			enc.Encode(&pid.AlgType)
			enc.Encode(&pid.Params)
			enc.Tag(pid)
		}

		enc.Pointer(&rec.PIDUsed)
		enc.Encode(&rec.GoodnessOfPID)

		enc.Encode(int32(len(rec.Recs)))
		for i := range rec.Recs {
			enc.Pointer(&rec.Recs[i])
		}

		enc.Encode(int32(len(rec.Tracks)))
		for i := range rec.Tracks {
			enc.Pointer(&rec.Tracks[i])
		}

		enc.Encode(int32(len(rec.Clusters)))
		for i := range rec.Clusters {
			enc.Pointer(&rec.Clusters[i])
		}

		enc.Pointer(&rec.StartVtx)
		enc.Tag(rec)
	}
	return enc.Err()
}

func (recs *RecParticleContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&recs.Flags)
	dec.Decode(&recs.Params)
	var n int32
	dec.Decode(&n)
	recs.Parts = make([]RecParticle, int(n))
	if r.VersionSio() <= 1002 {
		return fmt.Errorf("lcio: too old file (%d)", r.VersionSio())
	}

	for i := range recs.Parts {
		rec := &recs.Parts[i]
		dec.Decode(&rec.Type)
		dec.Decode(&rec.P)
		dec.Decode(&rec.Energy)
		dec.Decode(&rec.Cov)
		dec.Decode(&rec.Mass)
		dec.Decode(&rec.Charge)
		dec.Decode(&rec.Ref)

		var n int32
		dec.Decode(&n)
		rec.PIDs = make([]ParticleID, int(n))
		for i := range rec.PIDs {
			pid := &rec.PIDs[i]
			dec.Decode(&pid.Likelihood)
			dec.Decode(&pid.Type)
			dec.Decode(&pid.PDG)
			dec.Decode(&pid.AlgType)
			dec.Decode(&pid.Params)
			dec.Tag(pid)
		}
		dec.Pointer(&rec.PIDUsed)
		dec.Decode(&rec.GoodnessOfPID)

		dec.Decode(&n)
		rec.Recs = make([]*RecParticle, int(n))
		for i := range rec.Recs {
			dec.Pointer(&rec.Recs[i])
		}

		dec.Decode(&n)
		rec.Tracks = make([]*Track, int(n))
		for i := range rec.Tracks {
			dec.Pointer(&rec.Tracks[i])
		}

		dec.Decode(&n)
		rec.Clusters = make([]*Cluster, int(n))
		for i := range rec.Clusters {
			dec.Pointer(&rec.Clusters[i])
		}

		if r.VersionSio() > 1007 {
			dec.Pointer(&rec.StartVtx)
		}

		dec.Tag(rec)
	}

	return dec.Err()
}

var (
	_ sio.Versioner = (*RecParticleContainer)(nil)
	_ sio.Codec     = (*RecParticleContainer)(nil)
)
