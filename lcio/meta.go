// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"go-hep.org/x/hep/sio"
)

const (
	MajorVersion uint32 = 2
	MinorVersion uint32 = 8
	Version      uint32 = (MajorVersion << 16) + MinorVersion
)

var Records = struct {
	Index        string
	RandomAccess string
	RunHeader    string
	EventHeader  string
	Event        string
}{
	Index:        "LCIOIndex",
	RandomAccess: "LCIORandomAccess",
	RunHeader:    "LCRunHeader",
	EventHeader:  "LCEventHeader",
	Event:        "LCEvent",
}

var Blocks = struct {
	Index        string
	RandomAccess string
	RunHeader    string
	EventHeader  string
}{
	Index:        "LCIOndex",
	RandomAccess: "LCIORandomAccess",
	RunHeader:    "RunHeader",
	EventHeader:  "EventHeader",
}

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

var _ sio.Codec = (*Index)(nil)
