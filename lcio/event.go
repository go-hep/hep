// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
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

type Event struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Collections map[string]interface{}
	Names       []string
	Params      Params
}

type McParticles struct {
	Flags     Flags
	Params    Params
	Particles []McParticle
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
	Momentum  [3]float64 // Momentum at production vertex
	Mass      float64
	Charge    float32
	EndPoint  [3]float64
	PEndPoint [3]float64
	Spin      [3]float32
	ColorFlow [2]int32
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
	mom := [3]float32{float32(mc.Momentum[0]), float32(mc.Momentum[1]), float32(mc.Momentum[2])}
	enc.Encode(&mom)
	mass := float32(mc.Mass)
	enc.Encode(&mass)
	enc.Encode(&mc.Charge)
	if mc.SimStatus&uint32(1<<31) != 0 {
		enc.Encode(&mc.EndPoint)
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
	mc.Momentum[0] = float64(mom[0])
	mc.Momentum[1] = float64(mom[1])
	mc.Momentum[2] = float64(mom[2])

	var mass float32
	dec.Decode(&mass)
	mc.Mass = float64(mass)

	dec.Decode(&mc.Charge)

	if mc.SimStatus&uint32(1<<31) != 0 {
		dec.Decode(&mc.EndPoint)
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
		dec.Decode(&hit.CellID1)
		dec.Decode(&hit.Energy)
		if r.VersionSio() > 1009 && hits.Flags.Test(RChBitEnergyError) {
			dec.Decode(&hit.EnergyErr)
		}
		if r.VersionSio() > 1002 && hits.Flags.Test(RChBitTime) {
			dec.Decode(&hit.Time)
		}
		if hits.Flags.Test(RChBitBarrel) {
			dec.Decode(&hit.Pos)
		}
		if r.VersionSio() > 1002 {
			dec.Decode(&hit.Type)
			dec.Pointer(&hit.Raw)
		}
		if r.VersionSio() > 1002 {
			// the logic of the pointer bit has been inverted in v1.3
			if hits.Flags.Test(RChBitNoPtr) {
				dec.Tag(hit)
			}
		} else {
			if !hits.Flags.Test(RChBitNoPtr) {
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

type GenericObjectData struct {
	I32s []int32
	F32s []float32
	F64s []float64
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
