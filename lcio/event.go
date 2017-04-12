// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"encoding/binary"
	"log"

	"go-hep.org/x/hep/sio"
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
	var err error
	err = sio.Unmarshal(r, &idx.ControlWord)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &idx.RunMin)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &idx.BaseOffset)
	if err != nil {
		return err
	}

	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}

	idx.Offsets = make([]Offset, int(n))
	for i := range idx.Offsets {
		v := &idx.Offsets[i]
		if idx.ControlWord&1 == 0 {
			err = sio.Unmarshal(r, &v.RunOffset)
			if err != nil {
				return err
			}
		}

		err = sio.Unmarshal(r, &v.EventNumber)
		if err != nil {
			return err
		}

		switch {
		case idx.ControlWord&2 == 1:
			err = sio.Unmarshal(r, &v.Location)
		default:
			var loc int32
			err = sio.Unmarshal(r, &loc)
			v.Location = int64(loc)
		}
		if err != nil {
			return err
		}

		if idx.ControlWord&4 == 1 {
			err = sio.Unmarshal(r, &v.Ints)
			if err != nil {
				return err
			}

			err = sio.Unmarshal(r, &v.Floats)
			if err != nil {
				return err
			}

			err = sio.Unmarshal(r, &v.Strings)
			if err != nil {
				return err
			}
		}
	}

	return err
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

type EventHeader struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Blocks      []Block
	Params      Params
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

func (mc *McParticles) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (mc *McParticles) UnmarshalSio(r sio.Reader) error {
	var err error

	err = sio.Unmarshal(r, &mc.Flags)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &mc.Params)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &mc.Particles)
	if err != nil {
		return err
	}

	return err
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
			for _, mom := range p.Parents {
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

func (mc *McParticle) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (mc *McParticle) UnmarshalSio(r sio.Reader) error {
	var err error

	err = r.Tag(mc)
	if err != nil {
		return err
	}

	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}
	if n > 0 {
		mc.Parents = make([]*McParticle, int(n))
		for ii := range mc.Parents {
			err = r.Pointer(&mc.Parents[ii])
			if err != nil {
				return err
			}
		}
	}

	err = sio.Unmarshal(r, &mc.PDG)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &mc.GenStatus)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &mc.SimStatus)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &mc.Vertex)
	if err != nil {
		return err
	}

	if r.Version() > 1002 {
		err = sio.Unmarshal(r, &mc.Time)
		if err != nil {
			return err
		}
	}

	var mom [3]float32
	err = sio.Unmarshal(r, &mom)
	if err != nil {
		return err
	}
	mc.Momentum[0] = float64(mom[0])
	mc.Momentum[1] = float64(mom[1])
	mc.Momentum[2] = float64(mom[2])

	var mass float32
	err = sio.Unmarshal(r, &mass)
	if err != nil {
		return err
	}
	mc.Mass = float64(mass)

	err = sio.Unmarshal(r, &mc.Charge)
	if err != nil {
		return err
	}

	if mc.SimStatus&uint32(1<<31) != 0 {
		err = sio.Unmarshal(r, &mc.EndPoint)
		if err != nil {
			return err
		}
		if r.Version() > 2006 {
			var mom [3]float32
			err = sio.Unmarshal(r, &mom)
			if err != nil {
				return err
			}
			mc.PEndPoint[0] = float64(mom[0])
			mc.PEndPoint[1] = float64(mom[1])
			mc.PEndPoint[2] = float64(mom[2])
		}
	}

	if r.Version() > 1051 {
		err = sio.Unmarshal(r, &mc.Spin)
		if err != nil {
			return err
		}

		err = sio.Unmarshal(r, &mc.ColorFlow)
		if err != nil {
			return err
		}
	}
	return err
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
	var err error
	err = sio.Unmarshal(r, &hits.Flags)
	if err != nil {
		return err
	}
	err = sio.Unmarshal(r, &hits.Params)
	if err != nil {
		return err
	}
	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}
	hits.Hits = make([]SimTrackerHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		err = sio.Unmarshal(r, &hit.CellID0)
		if err != nil {
			return err
		}
		if r.Version() > 1051 && hits.Flags.Test(ThBitID1) {
			err = sio.Unmarshal(r, &hit.CellID1)
			if err != nil {
				return err
			}
		}
		err = sio.Unmarshal(r, &hit.Pos)
		if err != nil {
			return err
		}
		err = sio.Unmarshal(r, &hit.EDep)
		if err != nil {
			return err
		}
		err = sio.Unmarshal(r, &hit.Time)
		if err != nil {
			return err
		}
		err = r.Pointer(&hit.Mc)
		if err != nil {
			return err
		}
		if hits.Flags.Test(ThBitMomentum) {
			err = sio.Unmarshal(r, &hit.Momentum)
			if err != nil {
				return err
			}
			err = sio.Unmarshal(r, &hit.PathLength)
			if err != nil {
				return err
			}
		}
		if r.Version() > 1000 {
			err = r.Tag(hit)
			if err != nil {
				return err
			}
		}
	}
	return err
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
}

type SimCalorimeterHits struct {
	Flags  Flags
	Params Params
	Hits   []SimCalorimeterHit
}

func (hits *SimCalorimeterHits) MarshalSio(w sio.Writer) error {
	panic("not implemented")
}

func (hits *SimCalorimeterHits) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &hits.Flags)
	if err != nil {
		return err
	}
	err = sio.Unmarshal(r, &hits.Params)
	if err != nil {
		return err
	}
	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}
	hits.Hits = make([]SimCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		err = sio.Unmarshal(r, &hit.CellID0)
		if err != nil {
			return err
		}
		if r.Version() < 9 || hits.Flags.Test(ChBitID1) {
			err = sio.Unmarshal(r, &hit.CellID1)
			if err != nil {
				return err
			}
		}
		err = sio.Unmarshal(r, &hit.Energy)
		if err != nil {
			return err
		}
		if hits.Flags.Test(ChBitLong) {
			err = sio.Unmarshal(r, &hit.Pos)
			if err != nil {
				return err
			}
		}
		var n int32
		err = sio.Unmarshal(r, &n)
		if err != nil {
			return err
		}
		hit.Contributions = make([]Contrib, int(n))
		for i := range hit.Contributions {
			c := &hit.Contributions[i]
			err = r.Pointer(&c.Mc)
			if err != nil {
				return err
			}
			err = sio.Unmarshal(r, &c.Energy)
			if err != nil {
				return err
			}
			err = sio.Unmarshal(r, &c.Time)
			if err != nil {
				return err
			}
			if hits.Flags.Test(ChBitStep) {
				err = sio.Unmarshal(r, &c.PDG)
				if err != nil {
					return err
				}
				if r.Version() > 1051 {
					err = sio.Unmarshal(r, &c.StepPos)
					if err != nil {
						return err
					}
				}
			}
		}
		if r.Version() > 1000 {
			err = r.Tag(hit)
			if err != nil {
				return err
			}
		}
	}
	return err
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
	Elements []float32
}

type IntVec struct {
	Flags    Flags
	Params   Params
	Elements []int32
}

type StrVec struct {
	Flags    Flags
	Params   Params
	Elements []string
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
	var err error
	err = sio.Unmarshal(r, &hits.Flags)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &hits.Params)
	if err != nil {
		return err
	}
	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}
	hits.Hits = make([]RawCalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		err = binary.Read(r, binary.BigEndian, &hit.CellID0)
		if err != nil {
			return err
		}
		if r.Version() == 8 || hits.Flags.Test(RChBitID1) {
			err = binary.Read(r, binary.BigEndian, &hit.CellID1)
			if err != nil {
				return err
			}
		}
		err = binary.Read(r, binary.BigEndian, &hit.Amplitude)
		if err != nil {
			return err
		}
		if hits.Flags.Test(RChBitTime) {
			err = binary.Read(r, binary.BigEndian, &hit.TimeStamp)
			if err != nil {
				return err
			}
		}
		if !hits.Flags.Test(RChBitNoPtr) {
			err = r.Tag(hit)
			if err != nil {
				return err
			}
		}
	}
	return err
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
	var err error
	err = sio.Unmarshal(r, &hits.Flags)
	if err != nil {
		return err
	}
	err = sio.Unmarshal(r, &hits.Params)
	if err != nil {
		return err
	}
	var n int32
	err = sio.Unmarshal(r, &n)
	if err != nil {
		return err
	}

	hits.Hits = make([]CalorimeterHit, int(n))
	for i := range hits.Hits {
		hit := &hits.Hits[i]
		err = sio.Unmarshal(r, &hit.CellID0)
		if err != nil {
			log.Panic(err)
			return err
		}
		err = sio.Unmarshal(r, &hit.CellID1)
		if err != nil {
			log.Panic(err)
			return err
		}
		err = sio.Unmarshal(r, &hit.Energy)
		if err != nil {
			log.Panic(err)
			return err
		}
		if r.Version() > 1009 && hits.Flags.Test(RChBitEnergyError) {
			err = sio.Unmarshal(r, &hit.EnergyErr)
			if err != nil {
				log.Panic(err)
				return err
			}
		}
		if r.Version() > 1002 && hits.Flags.Test(RChBitTime) {
			sio.Unmarshal(r, &hit.Time)
			if err != nil {
				return err
			}
		}
		if hits.Flags.Test(RChBitBarrel) {
			err = sio.Unmarshal(r, &hit.Pos)
			if err != nil {
				log.Panic(err)
				return err
			}
		}
		if r.Version() > 1002 {
			err = sio.Unmarshal(r, &hit.Type)
			if err != nil {
				log.Panic(err)
				return err
			}

			err = r.Pointer(&hit.Raw)
			if err != nil {
				return err
			}
		}
		if r.Version() > 1002 {
			// the logic of the pointer bit has been inverted in v1.3
			if hits.Flags.Test(RChBitNoPtr) {
				err = r.Tag(hit)
				if err != nil {
					return err
				}
			}
		} else {
			if !hits.Flags.Test(RChBitNoPtr) {
				err = r.Tag(hit)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
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
	var err error
	err = sio.Unmarshal(r, &obj.Flag)
	if err != nil {
		return err
	}

	err = sio.Unmarshal(r, &obj.Params)
	if err != nil {
		return err
	}

	var (
		ni32  int32
		nf32  int32
		nf64  int32
		nobjs int32
	)

	if obj.Flag.Test(GOBitFixed) {
		err = sio.Unmarshal(r, &ni32)
		if err != nil {
			return err
		}
		err = sio.Unmarshal(r, &nf32)
		if err != nil {
			return err
		}
		err = sio.Unmarshal(r, &nf64)
		if err != nil {
			return err
		}
	}
	err = sio.Unmarshal(r, &nobjs)
	if err != nil {
		return err
	}
	obj.Data = make([]GenericObjectData, int(nobjs))
	for iobj := range obj.Data {
		data := &obj.Data[iobj]
		if !obj.Flag.Test(GOBitFixed) {

			err = sio.Unmarshal(r, &ni32)
			if err != nil {
				return err
			}
			err = sio.Unmarshal(r, &nf32)
			if err != nil {
				return err
			}
			err = sio.Unmarshal(r, &nf64)
			if err != nil {
				return err
			}
		}
		data.I32s = make([]int32, int(ni32))
		for i := range data.I32s {
			err = sio.Unmarshal(r, &data.I32s[i])
			if err != nil {
				return err
			}
		}
		data.F32s = make([]float32, int(nf32))
		for i := range data.F32s {
			err = sio.Unmarshal(r, &data.F32s[i])
			if err != nil {
				return err
			}
		}
		data.F64s = make([]float64, int(nf64))
		for i := range data.F64s {
			err = sio.Unmarshal(r, &data.F64s[i])
			if err != nil {
				return err
			}
		}

		err = r.Tag(data)
		if err != nil {
			return err
		}
	}

	return err
}

var _ sio.Codec = (*Index)(nil)
var _ sio.Codec = (*McParticle)(nil)
var _ sio.Codec = (*McParticles)(nil)
var _ sio.Codec = (*GenericObject)(nil)
var _ sio.Codec = (*SimTrackerHits)(nil)
var _ sio.Codec = (*SimCalorimeterHits)(nil)
var _ sio.Codec = (*RawCalorimeterHits)(nil)
var _ sio.Codec = (*CalorimeterHits)(nil)
