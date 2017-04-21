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

// ClusterContainer is a collection of clusters.
type ClusterContainer struct {
	Flags    Flags
	Params   Params
	Clusters []Cluster
}

type Cluster struct {
	// Type of cluster:
	//  - bits 31-16: ECAL, HCAL, COMBINED, LAT, LCAL
	//  - bits 15-00: NEUTRAL, CHARGED, UNDEFINED
	Type       int32
	Energy     float32    // energy of the cluster
	EnergyErr  float32    // energy error of the cluster
	Pos        [3]float32 // center of cluster (x,y,z)
	PosErr     [6]float32 // covariance matrix of position
	Theta      float32    // intrinsic direction: theta at position
	Phi        float32    // intrinsic direction: phi at position
	DirErr     [3]float32 // covariance matrix of direct
	Shape      []float32  // shape parameters, defined in collection parameter 'ShapeParameterNames'
	PIDs       []ParticleID
	Clusters   []*Cluster        // clusters combined into this cluster
	Hits       []*CalorimeterHit // hits that made this cluster
	Weights    []float32         // energy fraction of the hit that contributed to this cluster
	SubDetEnes []float32         // energy observed in a particular subdetector
}

func (clus ClusterContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of Cluster collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", clus.Flags, clus.Params)
	fmt.Fprintf(o, "     LCIO::CLBIT_HITS : %v\n", clus.Flags.Test(BitsClHits))

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |type|  energy  |energyerr |      position ( x,y,z)           |  itheta  |   iphi   \n"
		tail = "------------|----|----------|----------|----------------------------------|----------|----------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range clus.Clusters {
		clu := &clus.Clusters[i]
		fmt.Fprintf(o,
			" [%08d] |%4d|%+.3e|%+.3e|%+.3e, %+.3e, %+.3e|%+.3e|%+.3e\n",
			0, // id
			clu.Type, clu.Energy, clu.EnergyErr,
			clu.Pos[0], clu.Pos[1], clu.Pos[2],
			clu.Theta,
			clu.Phi,
		)
		fmt.Fprintf(o,
			"            errors (6 pos)/(3 dir): (%+.3e, %+.3e, %+.3e, %+.3e, %+.3e, %+.3e)/(%+.3e, %+.3e, %+.3e)\n",
			clu.PosErr[0], clu.PosErr[1], clu.PosErr[2], clu.PosErr[3], clu.PosErr[4], clu.PosErr[5],
			clu.DirErr[0], clu.DirErr[1], clu.DirErr[2],
		)
		fmt.Fprintf(o, "            clusters(e): ")
		for ii, cc := range clu.Clusters {
			var e float32
			if cc != nil {
				e = cc.Energy
			}
			if ii > 0 {
				fmt.Fprintf(o, ", ")
			}
			fmt.Fprintf(o, "%+.3e", e)
		}
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "            subdetector energies : ")
		for ii, ee := range clu.SubDetEnes {
			if ii > 0 {
				fmt.Fprintf(o, ", ")
			}
			fmt.Fprintf(o, "%+.3e", ee)
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)

	return string(o.Bytes())
}

func (*ClusterContainer) VersionSio() uint32 {
	return Version
}

func (clus *ClusterContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&clus.Flags)
	enc.Encode(&clus.Params)
	enc.Encode(int32(len(clus.Clusters)))
	for i := range clus.Clusters {
		clu := &clus.Clusters[i]
		enc.Encode(&clu.Type)
		enc.Encode(&clu.Energy)
		enc.Encode(&clu.EnergyErr)
		enc.Encode(&clu.Pos)
		enc.Encode(&clu.PosErr)
		enc.Encode(&clu.Theta)
		enc.Encode(&clu.Phi)
		enc.Encode(&clu.DirErr)
		enc.Encode(&clu.Shape)
		enc.Encode(int32(len(clu.PIDs)))
		for i := range clu.PIDs {
			pid := &clu.PIDs[i]
			enc.Encode(&pid.Likelihood)
			enc.Encode(&pid.Type)
			enc.Encode(&pid.PDG)
			enc.Encode(&pid.AlgType)
			enc.Encode(&pid.Params)
		}

		enc.Encode(int32(len(clu.Clusters)))
		for i := range clu.Clusters {
			enc.Pointer(&clu.Clusters[i])
		}

		if clus.Flags.Test(BitsClHits) {
			enc.Encode(int32(len(clu.Hits)))
			for i := range clu.Hits {
				enc.Pointer(&clu.Hits[i])
				enc.Encode(&clu.Weights[i])
			}
		}
		enc.Encode(&clu.SubDetEnes)

		enc.Tag(clu)

	}
	return enc.Err()
}

func (clus *ClusterContainer) UnmarshalSio(r sio.Reader) error {
	const NShapeOld = 6

	dec := sio.NewDecoder(r)
	dec.Decode(&clus.Flags)
	dec.Decode(&clus.Params)
	var n int32
	dec.Decode(&n)
	clus.Clusters = make([]Cluster, int(n))
	for i := range clus.Clusters {
		clu := &clus.Clusters[i]
		dec.Decode(&clu.Type)
		dec.Decode(&clu.Energy)
		if r.VersionSio() > 1051 {
			dec.Decode(&clu.EnergyErr)
		}
		dec.Decode(&clu.Pos)
		dec.Decode(&clu.PosErr)
		dec.Decode(&clu.Theta)
		dec.Decode(&clu.Phi)
		dec.Decode(&clu.DirErr)

		var n int32 = NShapeOld
		if r.VersionSio() > 1002 {
			dec.Decode(&n)
		}
		clu.Shape = make([]float32, int(n))
		for i := range clu.Shape {
			dec.Decode(&clu.Shape[i])
		}

		if r.VersionSio() > 1002 {
			var n int32
			dec.Decode(&n)
			clu.PIDs = make([]ParticleID, int(n))
			for i := range clu.PIDs {
				pid := &clu.PIDs[i]
				dec.Decode(&pid.Likelihood)
				dec.Decode(&pid.Type)
				dec.Decode(&pid.PDG)
				dec.Decode(&pid.AlgType)
				dec.Decode(&pid.Params)
			}
		} else {
			var ptype [3]float32
			dec.Decode(&ptype)
		}

		dec.Decode(&n)
		clu.Clusters = make([]*Cluster, int(n))
		for i := range clu.Clusters {
			dec.Pointer(&clu.Clusters[i])
		}

		if clus.Flags.Test(BitsClHits) {
			dec.Decode(&n)
			clu.Hits = make([]*CalorimeterHit, int(n))
			clu.Weights = make([]float32, int(n))

			for i := range clu.Hits {
				dec.Pointer(&clu.Hits[i])
				dec.Decode(&clu.Weights[i])
			}
		}
		dec.Decode(&clu.SubDetEnes)

		dec.Tag(clu)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*ClusterContainer)(nil)
	_ sio.Codec     = (*ClusterContainer)(nil)
)
