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

type FloatVec struct {
	Flags    Flags
	Params   Params
	Elements [][]float32
}

func (vec FloatVec) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of LCFloatVec collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", vec.Flags, vec.Params)
	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] | val0, val1, ...\n"
		tail = "------------|----------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, slice := range vec.Elements {
		fmt.Fprintf(o, " [%08d] |",
			0, //id
		)
		for i, v := range slice {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			if i+1%10 == 0 {
				fmt.Fprintf(o, "\n     ")
			}
			fmt.Fprintf(o, "%8f", v)
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
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

func (vec IntVec) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of LCIntVec collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", vec.Flags, vec.Params)
	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] | val0, val1, ...\n"
		tail = "------------|----------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, slice := range vec.Elements {
		fmt.Fprintf(o, " [%08d] |",
			0, //id
		)
		for i, v := range slice {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			if i+1%10 == 0 {
				fmt.Fprintf(o, "\n     ")
			}
			fmt.Fprintf(o, "%8d", v)
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
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

func (vec StrVec) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of LCStrVec collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n%v", vec.Flags, vec.Params)
	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] | val0, val1, ...\n"
		tail = "------------|----------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for _, slice := range vec.Elements {
		fmt.Fprintf(o, " [%08d] |",
			0, //id
		)
		for i, v := range slice {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			if i+1%10 == 0 {
				fmt.Fprintf(o, "\n     ")
			}
			fmt.Fprintf(o, "%s", v)
		}
		fmt.Fprintf(o, "\n")
	}
	fmt.Fprintf(o, tail)
	return string(o.Bytes())
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

type GenericObjectData struct {
	I32s []int32
	F32s []float32
	F64s []float64
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
			obj.Flag.Test(BitsGOFixed),
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

func (*GenericObject) VersionSio() uint32 {
	return Version
}

func (obj *GenericObject) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&obj.Flag)
	enc.Encode(&obj.Params)

	if obj.Flag.Test(BitsGOFixed) {
		var (
			ni32 int32
			nf32 int32
			nf64 int32
		)

		if len(obj.Data) > 0 {
			data := obj.Data[0]
			ni32 = int32(len(data.I32s))
			nf32 = int32(len(data.F32s))
			nf64 = int32(len(data.F64s))
		}
		enc.Encode(&ni32)
		enc.Encode(&nf32)
		enc.Encode(&nf64)
	}
	enc.Encode(int32(len(obj.Data)))
	for iobj := range obj.Data {
		data := &obj.Data[iobj]
		if !obj.Flag.Test(BitsGOFixed) {
			enc.Encode(int32(len(data.I32s)))
			enc.Encode(int32(len(data.F32s)))
			enc.Encode(int32(len(data.F64s)))
		}
		for i := range data.I32s {
			enc.Encode(&data.I32s[i])
		}
		for i := range data.F32s {
			enc.Encode(&data.F32s[i])
		}
		for i := range data.F64s {
			enc.Encode(&data.F64s[i])
		}
		enc.Tag(data)
	}

	return enc.Err()
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

	if obj.Flag.Test(BitsGOFixed) {
		dec.Decode(&ni32)
		dec.Decode(&nf32)
		dec.Decode(&nf64)
	}
	dec.Decode(&nobjs)
	obj.Data = make([]GenericObjectData, int(nobjs))
	for iobj := range obj.Data {
		data := &obj.Data[iobj]
		if !obj.Flag.Test(BitsGOFixed) {
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

var (
	_ sio.Versioner = (*FloatVec)(nil)
	_ sio.Codec     = (*FloatVec)(nil)
	_ sio.Versioner = (*IntVec)(nil)
	_ sio.Codec     = (*IntVec)(nil)
	_ sio.Versioner = (*StrVec)(nil)
	_ sio.Codec     = (*StrVec)(nil)
	_ sio.Versioner = (*GenericObject)(nil)
	_ sio.Codec     = (*GenericObject)(nil)
)
