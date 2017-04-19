// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"go-hep.org/x/hep/sio"
)

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
	fmt.Fprintf(o, " date       %v\n", time.Unix(0, evt.TimeStamp).UTC().Format("02.01.2006 15:04:05.999999999"))
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

var _ sio.Codec = (*FloatVec)(nil)
var _ sio.Codec = (*IntVec)(nil)
var _ sio.Codec = (*StrVec)(nil)
var _ sio.Codec = (*GenericObject)(nil)
