// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	"go-hep.org/x/hep/sio"
)

const (
	VersionMajor uint32 = 2
	VersionMinor uint32 = 8
	Version      uint32 = (VersionMajor << 16) + VersionMinor
)

type RandomAccess struct {
	RunMin         int32
	EventMin       int32
	RunMax         int32
	EventMax       int32
	RunHeaders     int32
	Events         int32
	RecordsInOrder int32
	IndexLoc       int64
	PrevLoc        int64
	NextLoc        int64
	FirstRecordLoc int64
	RecordSize     int32
}

type Index struct {
	// Bit 0 = single run.
	// Bit 1 = int64 offset required
	// Bit 2 = Params included (not yet implemented)
	ControlWord uint32
	RunMin      int32
	BaseOffset  int64
	Offsets     []Offset
}

func (idx *Index) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (idx *Index) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&idx.ControlWord)
	dec.Decode(&idx.RunMin)
	dec.Decode(&idx.BaseOffset)
	var n int32
	dec.Decode(&n)
	idx.Offsets = make([]Offset, int(n))
	for i := range idx.Offsets {
		v := &idx.Offsets[i]
		if idx.ControlWord&1 == 0 {
			dec.Decode(&v.RunOffset)
		}

		dec.Decode(&v.EventNumber)
		switch {
		case idx.ControlWord&2 == 1:
			dec.Decode(&v.Location)
		default:
			var loc int32
			dec.Decode(&loc)
			v.Location = int64(loc)
		}
		if idx.ControlWord&4 == 1 {
			dec.Decode(&v.Ints)
			dec.Decode(&v.Floats)
			dec.Decode(&v.Strings)
		}
	}
	return dec.Err()
}

type Offset struct {
	RunOffset   int32 // run offset relative to Index.RunMin
	EventNumber int32 // event number or -1 for run header records
	Location    int64 // location offset relative to Index.BaseOffset
	Ints        []int32
	Floats      []float32
	Strings     []string
}

type RunHeader struct {
	RunNbr       int32
	Detector     string
	Descr        string
	SubDetectors []string
	Params       Params
}

func (*RunHeader) VersionSio() uint32 {
	return Version
}

type EventHeader struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Blocks      []Block
	Params      Params
}

func (*EventHeader) VersionSio() uint32 {
	return Version
}

type Block struct {
	Name string
	Type string
}

type Params struct {
	Ints    map[string][]int32
	Floats  map[string][]float32
	Strings map[string][]string
}

func (p Params) String() string {
	o := new(bytes.Buffer)
	for k, vec := range p.Ints {
		fmt.Fprintf(o, " parameter %s [int]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	for k, vec := range p.Floats {
		fmt.Fprintf(o, " parameter %s [float]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	for k, vec := range p.Strings {
		fmt.Fprintf(o, " parameter %s [string]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	return string(o.Bytes())
}

type Event struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Collections map[string]interface{}
	Names       []string
	Params      Params
}

func (evt *Event) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, "        Event  : %d - run:   %d - timestamp %v - weight %v\n",
		evt.EventNumber, evt.RunNumber, evt.TimeStamp, evt.Weight(),
	)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, " date       %v\n", time.Unix(0, evt.TimeStamp).UTC().Format("02.01.2006 03:04:05.999999999"))
	fmt.Fprintf(o, " detector : %s\n", evt.Detector)
	fmt.Fprintf(o, " event parameters:\n%v\n", evt.Params)

	for _, name := range evt.Names {
		coll := evt.Collections[name]
		fmt.Fprintf(o, " collection name : %s\n parameters: \n%v\n", name, coll)
	}
	return string(o.Bytes())
}

func (evt *Event) Weight() float64 {
	if v, ok := evt.Params.Floats["_weight"]; ok {
		return float64(v[0])
	}
	return 1.0
}

type McParticles struct {
	Flags     Flags
	Params    Params
	Particles []McParticle
}

func (mcs McParticles) String() string {
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
		"[   id   ]"+
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
			0,
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

func (*McParticles) VersionSio() uint32 {
	return Version
}

func (mc *McParticles) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&mc.Flags)
	enc.Encode(&mc.Params)
	enc.Encode(&mc.Particles)
	return enc.Err()
}

func (mc *McParticles) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&mc.Flags)
	dec.Decode(&mc.Params)
	dec.Decode(&mc.Particles)
	return dec.Err()
}

func (mc *McParticles) LinkSio(vers uint32) error {
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

type SimTrackerHits struct {
	Flags  Flags
	Params Params
	Hits   []SimTrackerHit
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
		if r.VersionSio() > 1051 && hits.Flags.Test(ThBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Pos)
		dec.Decode(&hit.EDep)
		dec.Decode(&hit.Time)
		dec.Pointer(&hit.Mc)
		if hits.Flags.Test(ThBitMomentum) {
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

type SimCalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []SimCalorimeterHit
}

func (hits SimCalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of SimCalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "  -> LCIO::CHBIT_LONG   : %v\n", hits.Flags.Test(ChBitLong))
	fmt.Fprintf(o, "     LCIO::CHBIT_BARREL : %v\n", hits.Flags.Test(ChBitBarrel))
	fmt.Fprintf(o, "     LCIO::CHBIT_ID1    : %v\n", hits.Flags.Test(ChBitID1))
	fmt.Fprintf(o, "     LCIO::CHBIT_STEP   : %v\n", hits.Flags.Test(ChBitStep))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |cellId0 |cellId1 |  energy  |        position (x,y,z)          | nMCParticles \n" +
		"           -> MC contribution: prim. PDG |  energy  |   time   | sec. PDG | stepPosition (x,y,z) \n"
	tail := "------------|--------|--------|----------|----------------------------------|--------------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d|%08d|%+.3e|", 0, hit.CellID0, hit.CellID1, hit.Energy)
		if hits.Flags.Test(ChBitLong) {
			fmt.Fprintf(o, "+%.3e, %+.3e, %+.3e", hit.Pos[0], hit.Pos[1], hit.Pos[2])
		} else {
			fmt.Fprintf(o, "    no position available         ")
		}
		fmt.Fprintf(o, "|%+12d\n", len(hit.Contributions))
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "        id-fields: --- unknown/default ----   ")
		for _, c := range hit.Contributions {
			var pdg int32
			if c.Mc != nil {
				pdg = c.Mc.PDG
			}
			fmt.Fprintf(o, "\n           ->                  %+10d|%+1.3e|%+1.3e|", pdg, c.Energy, c.Time)
			if hits.Flags.Test(ChBitStep) {
				fmt.Fprintf(o, "%+d| (%+1.3e, %+1.3e, %+1.3e)", c.PDG, c.StepPos[0], c.StepPos[1], c.StepPos[2])
			} else {
				fmt.Fprintf(o, " no PDG")
			}
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (*SimCalorimeterHits) VersionSio() uint32 {
	return Version
}

func (hits *SimCalorimeterHits) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hits.Flags)
	enc.Encode(&hits.Params)
	enc.Encode(int32(len(hits.Hits)))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		enc.Encode(&hit.CellID0)
		if hits.Flags.Test(ChBitID1) {
			enc.Encode(&hit.CellID1)
		}
		enc.Encode(&hit.Energy)
		if hits.Flags.Test(ChBitLong) {
			enc.Encode(&hit.Pos)
		}
		enc.Encode(int32(len(hit.Contributions)))
		for i := range hit.Contributions {
			c := &hit.Contributions[i]
			enc.Pointer(&c.Mc)
			enc.Encode(&c.Energy)
			enc.Encode(&c.Time)
			if hits.Flags.Test(ChBitStep) {
				enc.Encode(&c.PDG)
				enc.Encode(&c.StepPos)
			}
		}
		enc.Tag(hit)
	}
	return enc.Err()
}

func (hits *SimCalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]SimCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() < 9 || hits.Flags.Test(ChBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Energy)
		if hits.Flags.Test(ChBitLong) {
			dec.Decode(&hit.Pos)
		}
		var n int32
		dec.Decode(&n)
		hit.Contributions = make([]Contrib, int(n))
		for i := range hit.Contributions {
			c := &hit.Contributions[i]
			dec.Pointer(&c.Mc)
			dec.Decode(&c.Energy)
			dec.Decode(&c.Time)
			if hits.Flags.Test(ChBitStep) {
				dec.Decode(&c.PDG)
				if r.VersionSio() > 1051 {
					dec.Decode(&c.StepPos)
				}
			}
		}
		if r.VersionSio() > 1000 {
			dec.Tag(hit)
		}
	}
	return dec.Err()
}

type SimCalorimeterHit struct {
	Params        Params
	CellID0       int32
	CellID1       int32
	Energy        float32
	Pos           [3]float32
	Contributions []Contrib
}

type Contrib struct {
	Mc      *McParticle
	Energy  float32
	Time    float32
	PDG     int32
	StepPos [3]float32
}

type FloatVec struct {
	Flags    Flags
	Params   Params
	Elements [][]float32
}

func (*FloatVec) VersionSio() uint32 {
	return Version
}

func (vec *FloatVec) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&vec.Flags)
	enc.Encode(&vec.Params)
	enc.Encode(vec.Elements)
	enc.Encode(int32(len(vec.Elements)))
	for i := range vec.Elements {
		enc.Encode(int32(len(vec.Elements[i])))
		for _, v := range vec.Elements[i] {
			enc.Encode(v)
		}
		if w.VersionSio() > 1002 {
			enc.Tag(&vec.Elements[i])
		}
	}
	return enc.Err()
}

func (vec *FloatVec) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&vec.Flags)
	dec.Decode(&vec.Params)
	var nvecs int32
	dec.Decode(&nvecs)
	vec.Elements = make([][]float32, int(nvecs))
	for i := range vec.Elements {
		var n int32
		dec.Decode(&n)
		vec.Elements[i] = make([]float32, int(n))
		for j := range vec.Elements[i] {
			dec.Decode(&vec.Elements[i][j])
		}
		if r.VersionSio() > 1002 {
			dec.Tag(&vec.Elements[i])
		}
	}
	return dec.Err()
}

type IntVec struct {
	Flags    Flags
	Params   Params
	Elements [][]int32
}

func (*IntVec) VersionSio() uint32 {
	return Version
}

func (vec *IntVec) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&vec.Flags)
	enc.Encode(&vec.Params)
	enc.Encode(vec.Elements)
	enc.Encode(int32(len(vec.Elements)))
	for i := range vec.Elements {
		enc.Encode(int32(len(vec.Elements[i])))
		for _, v := range vec.Elements[i] {
			enc.Encode(v)
		}
		if w.VersionSio() > 1002 {
			enc.Tag(&vec.Elements[i])
		}
	}
	return enc.Err()
}

func (vec *IntVec) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&vec.Flags)
	dec.Decode(&vec.Params)
	var nvecs int32
	dec.Decode(&nvecs)
	vec.Elements = make([][]int32, int(nvecs))
	for i := range vec.Elements {
		var n int32
		dec.Decode(&n)
		vec.Elements[i] = make([]int32, int(n))
		for j := range vec.Elements[i] {
			dec.Decode(&vec.Elements[i][j])
		}
		if r.VersionSio() > 1002 {
			dec.Tag(&vec.Elements[i])
		}
	}
	return dec.Err()
}

type StrVec struct {
	Flags    Flags
	Params   Params
	Elements [][]string
}

func (*StrVec) VersionSio() uint32 {
	return Version
}

func (vec *StrVec) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&vec.Flags)
	enc.Encode(&vec.Params)
	enc.Encode(vec.Elements)
	enc.Encode(int32(len(vec.Elements)))
	for i := range vec.Elements {
		enc.Encode(int32(len(vec.Elements[i])))
		for _, v := range vec.Elements[i] {
			enc.Encode(v)
		}
		if w.VersionSio() > 1002 {
			enc.Tag(&vec.Elements[i])
		}
	}
	return enc.Err()
}

func (vec *StrVec) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&vec.Flags)
	dec.Decode(&vec.Params)
	var nvecs int32
	dec.Decode(&nvecs)
	vec.Elements = make([][]string, int(nvecs))
	for i := range vec.Elements {
		var n int32
		dec.Decode(&n)
		vec.Elements[i] = make([]string, int(n))
		for j := range vec.Elements[i] {
			dec.Decode(&vec.Elements[i][j])
		}
		if r.VersionSio() > 1002 {
			dec.Tag(&vec.Elements[i])
		}
	}
	return dec.Err()
}

type RawCalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []RawCalorimeterHit
}

func (hits RawCalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of RawCalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "     LCIO::RCHBIT_ID1    : %v\n", hits.Flags.Test(RChBitID1))
	fmt.Fprintf(o, "     LCIO::RCHBIT_TIME   : %v\n", hits.Flags.Test(RChBitTime))
	fmt.Fprintf(o, "     LCIO::RCHBIT_NO_PTR : %v\n", hits.Flags.Test(RChBitNoPtr))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |  cellId0 ( M, S, I, J, K) |cellId1 | amplitude |  time  \n"
	tail := "------------|---------------------------|--------|-----------|---------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d%19s|%08d|%10d |%8d", 0, hit.CellID0, "", hit.CellID1, hit.Amplitude, hit.TimeStamp)
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "\n        id-fields: --- unknown/default ----   ")
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (hits *RawCalorimeterHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *RawCalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]RawCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() == 8 || hits.Flags.Test(RChBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Amplitude)
		if hits.Flags.Test(RChBitTime) {
			dec.Decode(&hit.TimeStamp)
		}
		if !hits.Flags.Test(RChBitNoPtr) {
			dec.Tag(hit)
		}
	}
	return dec.Err()
}

type RawCalorimeterHit struct {
	CellID0   int32
	CellID1   int32
	Amplitude int32
	TimeStamp int32
}

type CalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []CalorimeterHit
}

func (hits CalorimeterHits) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of CalorimeterHit collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", hits.Flags, hits.Params)
	fmt.Fprintf(o, "  -> LCIO::RCHBIT_LONG   : %v\n", hits.Flags.Test(RChBitLong))
	fmt.Fprintf(o, "     LCIO::RCHBIT_BARREL : %v\n", hits.Flags.Test(RChBitBarrel))
	fmt.Fprintf(o, "     LCIO::RCHBIT_ID1    : %v\n", hits.Flags.Test(RChBitID1))
	fmt.Fprintf(o, "     LCIO::RCHBIT_TIME   : %v\n", hits.Flags.Test(RChBitTime))
	fmt.Fprintf(o, "     LCIO::RCHBIT_NO_PTR : %v\n", hits.Flags.Test(RChBitNoPtr))
	fmt.Fprintf(o, "     LCIO::RCHBIT_ENERGY_ERROR  : %v\n", hits.Flags.Test(RChBitEnergyError))

	// FIXME(sbinet): CellIDDecoder

	fmt.Fprintf(o, "\n")

	head := " [   id   ] |cellId0 |cellId1 |  energy  |energyerr |        position (x,y,z)           \n"
	tail := "------------|--------|--------|----------|----------|-----------------------------------\n"
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, hit := range hits.Hits {
		fmt.Fprintf(o, " [%08d] |%08d|%08d|%+.3e|%+.3e|", 0, hit.CellID0, hit.CellID1, hit.Energy, hit.EnergyErr)
		if hits.Flags.Test(ChBitLong) {
			fmt.Fprintf(o, "+%.3e, %+.3e, %+.3e", hit.Pos[0], hit.Pos[1], hit.Pos[2])
		} else {
			fmt.Fprintf(o, "    no position available         ")
		}
		// FIXME(sbinet): CellIDDecoder
		fmt.Fprintf(o, "\n        id-fields: --- unknown/default ----   ")
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
}

func (hits *CalorimeterHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *CalorimeterHits) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hits.Flags)
	dec.Decode(&hits.Params)
	var n int32
	dec.Decode(&n)
	hits.Hits = make([]CalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		dec.Decode(&hit.CellID0)
		if r.VersionSio() == 8 || hits.Flags.Test(RChBitID1) {
			dec.Decode(&hit.CellID1)
		}
		dec.Decode(&hit.Energy)
		if r.VersionSio() > 1009 && hits.Flags.Test(RChBitEnergyError) {
			dec.Decode(&hit.EnergyErr)
		}
		if r.VersionSio() > 1002 && hits.Flags.Test(RChBitTime) {
			dec.Decode(&hit.Time)
		}
		if hits.Flags.Test(RChBitLong) {
			dec.Decode(&hit.Pos)
		}
		if r.VersionSio() > 1002 {
			dec.Decode(&hit.Type)
			dec.Pointer(&hit.Raw)
		}
		if r.VersionSio() > 1002 {
			// the logic of the pointer bit has been inverted in v1.3
			if !hits.Flags.Test(RChBitNoPtr) {
				dec.Tag(hit)
			}
		} else {
			if hits.Flags.Test(RChBitNoPtr) {
				dec.Tag(hit)
			}
		}
	}
	return dec.Err()
}

type CalorimeterHit struct {
	CellID0   int32
	CellID1   int32
	Energy    float32
	EnergyErr float32
	Time      float32
	Pos       [3]float32
	Type      int32
	Raw       *RawCalorimeterHit
}

type TrackerHits struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHit
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
			if hits.Flags.Test(ThBitID1) {
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

type TrackerHitPlanes struct {
	Flags  Flags
	Params Params
	Hits   []TrackerHitPlane
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

type GenericObject struct {
	Flag   Flags
	Params Params
	Data   []GenericObjectData
}

func (obj GenericObject) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of LCGenericObject collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v\n", obj.Flag, obj.Params)
	fmt.Fprintf(o, " [   id   ] ")
	if obj.Data != nil {
		descr := ""
		if v := obj.Params.Strings["DataDescription"]; len(v) > 0 {
			descr = v[0]
		}
		fmt.Fprintf(o,
			"%s - isFixedSize: %v\n",
			descr,
			obj.Flag.Test(GOBitFixed),
		)
	} else {
		fmt.Fprintf(o, " Data.... \n")
	}

	tail := fmt.Sprintf(" %s", strings.Repeat("-", 55))

	fmt.Fprintf(o, "%s\n", tail)
	for _, iobj := range obj.Data {
		fmt.Fprintf(o, "%v\n", iobj)
		fmt.Fprintf(o, "%s\n", tail)
	}
	return string(o.Bytes())
}

type GenericObjectData struct {
	I32s []int32
	F32s []float32
	F64s []float64
}

func (obj GenericObjectData) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, " [%08d] ", 0)
	for _, v := range obj.I32s {
		fmt.Fprintf(o, "i:%d; ", v)
	}
	for _, v := range obj.F32s {
		fmt.Fprintf(o, "f:%f; ", v)
	}
	for _, v := range obj.F64s {
		fmt.Fprintf(o, "d:%f; ", v)
	}
	return string(o.Bytes())
}

func (obj *GenericObject) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (obj *GenericObject) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&obj.Flag)
	dec.Decode(&obj.Params)

	var (
		ni32  int32
		nf32  int32
		nf64  int32
		nobjs int32
	)

	if obj.Flag.Test(GOBitFixed) {
		dec.Decode(&ni32)
		dec.Decode(&nf32)
		dec.Decode(&nf64)
	}
	dec.Decode(&nobjs)
	obj.Data = make([]GenericObjectData, int(nobjs))
	for iobj := range obj.Data {
		data := &obj.Data[iobj]
		if !obj.Flag.Test(GOBitFixed) {
			dec.Decode(&ni32)
			dec.Decode(&nf32)
			dec.Decode(&nf64)
		}
		data.I32s = make([]int32, int(ni32))
		for i := range data.I32s {
			dec.Decode(&data.I32s[i])
		}
		data.F32s = make([]float32, int(nf32))
		for i := range data.F32s {
			dec.Decode(&data.F32s[i])
		}
		data.F64s = make([]float64, int(nf64))
		for i := range data.F64s {
			dec.Decode(&data.F64s[i])
		}

		dec.Tag(data)
	}

	return dec.Err()
}

var _ sio.Codec = (*Index)(nil)
var _ sio.Codec = (*McParticle)(nil)
var _ sio.Codec = (*McParticles)(nil)
var _ sio.Codec = (*GenericObject)(nil)
var _ sio.Codec = (*SimTrackerHits)(nil)
var _ sio.Codec = (*SimCalorimeterHits)(nil)
var _ sio.Codec = (*FloatVec)(nil)
var _ sio.Codec = (*IntVec)(nil)
var _ sio.Codec = (*StrVec)(nil)
var _ sio.Codec = (*RawCalorimeterHits)(nil)
var _ sio.Codec = (*CalorimeterHits)(nil)
var _ sio.Codec = (*TrackerHits)(nil)
var _ sio.Codec = (*TrackerHitPlanes)(nil)
