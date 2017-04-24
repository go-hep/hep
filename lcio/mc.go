// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"go-hep.org/x/hep/sio"
)

// McParticleContainer is a collection of monte-carlo particles.
type McParticleContainer struct {
	Flags     Flags
	Params    Params
	Particles []McParticle
}

func (mcs McParticleContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of MCParticle collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", mcs.Flags, mcs.Params)

	fmt.Fprintf(o,
		"simulator status bits: [sbvtcls] "+
			"s: created in simulation "+
			"b: backscatter "+
			"v: vertex is not endpoint of parent "+
			"t: decayed in tracker "+
			"c: decayed in calorimeter "+
			"l: has left detector "+
			"s: stopped o: overlay\n",
	)

	p2i := make(map[*McParticle]int, len(mcs.Particles))
	for i := range mcs.Particles {
		p := &mcs.Particles[i]
		p2i[p] = i
	}

	fmt.Fprintf(o,
		"[   id    ]"+
			"index|      PDG |    px,     py,        pz    | "+
			"px_ep,   py_ep , pz_ep      | "+
			"energy  |gen|[simstat ]| "+
			"vertex x,     y   ,   z     | "+
			"endpoint x,    y  ,   z     |    mass |  charge |            spin             | "+
			"colorflow | [parents] - [daughters]\n\n",
	)
	pfmt := "[%8.8d]%5d|%10d|% 1.2e,% 1.2e,% 1.2e|" +
		"% 1.2e,% 1.2e,% 1.2e|" +
		"% 1.2e|" +
		" %1d |" +
		"%s|" +
		"% 1.2e,% 1.2e,% 1.2e|" +
		"% 1.2e,% 1.2e,% 1.2e|" +
		"% 1.2e|" +
		"% 1.2e|" +
		"% 1.2e,% 1.2e,% 1.2e|" +
		"  (%d, %d)   |" +
		" ["

	for i := range mcs.Particles {
		p := &mcs.Particles[i]
		ep := p.EndPoint()
		fmt.Fprintf(o, pfmt,
			ID(p),
			i, p.PDG,
			p.P[0], p.P[1], p.P[2],
			p.PEndPoint[0], p.PEndPoint[1], p.PEndPoint[2],
			p.Energy(),
			p.GenStatus,
			p.SimStatusString(),
			p.Vertex[0], p.Vertex[1], p.Vertex[2],
			ep[0], ep[1], ep[2],
			p.Mass,
			p.Charge,
			p.Spin[0], p.Spin[1], p.Spin[2],
			p.ColorFlow[0], p.ColorFlow[1],
		)
		for ii := range p.Parents {
			if ii > 0 {
				fmt.Fprintf(o, ",")
			}
			fmt.Fprintf(o, "%d", p2i[p.Parents[ii]])
		}
		fmt.Fprintf(o, "] - [")
		for ii := range p.Children {
			if ii > 0 {
				fmt.Fprintf(o, ",")
			}
			fmt.Fprintf(o, "%d", p2i[p.Children[ii]])
		}
		fmt.Fprintf(o, "]\n")
	}

	fmt.Fprintf(o, "\n-------------------------------------------------------------------------------- \n")
	return string(o.Bytes())
}

func (*McParticleContainer) VersionSio() uint32 {
	return Version
}

func (mc *McParticleContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&mc.Flags)
	enc.Encode(&mc.Params)
	enc.Encode(&mc.Particles)
	return enc.Err()
}

func (mc *McParticleContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&mc.Flags)
	dec.Decode(&mc.Params)
	dec.Decode(&mc.Particles)
	return dec.Err()
}

func (mc *McParticleContainer) LinkSio(vers uint32) error {
	var err error
	switch {
	case vers <= 8:
		for i := range mc.Particles {
			p := &mc.Particles[i]
			for _, c := range p.Children {
				if c != nil {
					c.Parents = append(c.Parents, p)
				}
			}
		}

	default:
		for i := range mc.Particles {
			p := &mc.Particles[i]
			for i := range p.Parents {
				mom := p.Parents[i]
				if mom != nil {
					mom.Children = append(mom.Children, p)
				}
			}
		}
	}

	return err
}

type McParticle struct {
	Parents   []*McParticle
	Children  []*McParticle
	PDG       int32
	GenStatus int32
	SimStatus uint32
	Vertex    [3]float64
	Time      float32    // creation time of the particle in ns
	P         [3]float64 // Momentum at production vertex
	Mass      float64
	Charge    float32
	endPoint  [3]float64
	PEndPoint [3]float64 // momentum at end-point
	Spin      [3]float32
	ColorFlow [2]int32
}

func (mc *McParticle) Energy() float64 {
	px := mc.P[0]
	py := mc.P[1]
	pz := mc.P[2]
	return math.Sqrt(px*px + py*py + pz*pz + mc.Mass*mc.Mass)
}

func (mc *McParticle) EndPoint() [3]float64 {
	if mc.SimStatus&uint32(1<<31) == 0 {
		if len(mc.Children) == 0 {
			return mc.endPoint
		}
		for _, child := range mc.Children {
			if !child.VertexIsNotEnpointOfParent() {
				return child.Vertex
			}
		}
	}
	return mc.endPoint
}

func (mc *McParticle) SimStatusString() string {
	status := []byte("[    0   ]")
	if mc.SimStatus == 0 {
		return string(status)
	}
	if mc.IsCreatedInSimulation() {
		status[1] = 's'
	}
	if mc.IsBackScatter() {
		status[2] = 'b'
	}
	if mc.VertexIsNotEnpointOfParent() {
		status[3] = 'v'
	}
	if mc.IsDecayedInTracker() {
		status[4] = 't'
	}
	if mc.IsDecayedInCalorimeter() {
		status[5] = 'c'
	} else {
		status[5] = ' '
	}
	if mc.HasLeftDetector() {
		status[6] = 'l'
	}
	if mc.IsStopped() {
		status[7] = 's'
	}
	if mc.IsOverlay() {
		status[8] = 'o'
	}
	return string(status)
}

func (mc *McParticle) IsCreatedInSimulation() bool {
	return mc.SimStatus&uint32(1<<30) != 0
}

func (mc *McParticle) IsBackScatter() bool {
	return mc.SimStatus&uint32(1<<29) != 0
}

func (mc *McParticle) VertexIsNotEnpointOfParent() bool {
	return mc.SimStatus&uint32(1<<28) != 0
}

func (mc *McParticle) IsDecayedInTracker() bool {
	return mc.SimStatus&uint32(1<<27) != 0
}

func (mc *McParticle) IsDecayedInCalorimeter() bool {
	return mc.SimStatus&uint32(1<<26) != 0
}

func (mc *McParticle) HasLeftDetector() bool {
	return mc.SimStatus&uint32(1<<25) != 0
}

func (mc *McParticle) IsStopped() bool {
	return mc.SimStatus&uint32(1<<24) != 0
}

func (mc *McParticle) IsOverlay() bool {
	return mc.SimStatus&uint32(1<<23) != 0
}

func (mc *McParticle) VersionSio() uint32 {
	return Version
}

func (mc *McParticle) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Tag(mc)
	enc.Encode(int32(len(mc.Parents)))
	for ii := range mc.Parents {
		enc.Pointer(&mc.Parents[ii])
	}
	enc.Encode(&mc.PDG)
	enc.Encode(&mc.GenStatus)
	enc.Encode(&mc.SimStatus)
	enc.Encode(&mc.Vertex)
	enc.Encode(&mc.Time)
	mom := [3]float32{float32(mc.P[0]), float32(mc.P[1]), float32(mc.P[2])}
	enc.Encode(&mom)
	mass := float32(mc.Mass)
	enc.Encode(&mass)
	enc.Encode(&mc.Charge)
	if mc.SimStatus&uint32(1<<31) != 0 {
		enc.Encode(&mc.endPoint)
		pend := [3]float32{float32(mc.PEndPoint[0]), float32(mc.PEndPoint[1]), float32(mc.PEndPoint[2])}
		enc.Encode(&pend)
	}
	enc.Encode(&mc.Spin)
	enc.Encode(&mc.ColorFlow)
	return enc.Err()
}

func (mc *McParticle) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Tag(mc)

	var n int32
	dec.Decode(&n)
	if n > 0 {
		mc.Parents = make([]*McParticle, int(n))
		for ii := range mc.Parents {
			dec.Pointer(&mc.Parents[ii])
		}
	}

	dec.Decode(&mc.PDG)
	dec.Decode(&mc.GenStatus)
	dec.Decode(&mc.SimStatus)
	dec.Decode(&mc.Vertex)
	if r.VersionSio() > 1002 {
		dec.Decode(&mc.Time)
	}

	var mom [3]float32
	dec.Decode(&mom)
	mc.P[0] = float64(mom[0])
	mc.P[1] = float64(mom[1])
	mc.P[2] = float64(mom[2])

	var mass float32
	dec.Decode(&mass)
	mc.Mass = float64(mass)

	dec.Decode(&mc.Charge)

	if mc.SimStatus&uint32(1<<31) != 0 {
		dec.Decode(&mc.endPoint)
		if r.VersionSio() > 2006 {
			var mom [3]float32
			dec.Decode(&mom)
			mc.PEndPoint[0] = float64(mom[0])
			mc.PEndPoint[1] = float64(mom[1])
			mc.PEndPoint[2] = float64(mom[2])
		}
	}

	if r.VersionSio() > 1051 {
		dec.Decode(&mc.Spin)
		dec.Decode(&mc.ColorFlow)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*McParticle)(nil)
	_ sio.Codec     = (*McParticle)(nil)
	_ sio.Versioner = (*McParticleContainer)(nil)
	_ sio.Codec     = (*McParticleContainer)(nil)
)
