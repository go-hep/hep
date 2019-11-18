// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

const (
	defaultCellEncoding = "byte0:8,byte1:8,byte2:8,byte3:8,byte4:8,byte5:8,byte6:8,byte7:8"
	cellIDEncoding      = "CellIDEncoding"
)

// CellIDDecoder decodes cell IDs from a Hit cell-ID
type CellIDDecoder struct {
	bits  *bitField64
	codec string
}

func NewCellIDDecoder(codec string) *CellIDDecoder {
	return &CellIDDecoder{
		codec: codec,
		bits:  newBitField64(codec),
	}
}

func NewCellIDDecoderFrom(params Params) *CellIDDecoder {
	codec, ok := params.Strings[cellIDEncoding]
	if !ok {
		return nil
	}
	return NewCellIDDecoder(codec[0])
}

func (dec *CellIDDecoder) Get(hit Hit, name string) int64 {
	i := dec.bits.index[name]
	dec.bits.value = dec.value(hit)
	return dec.bits.fields[i].value()
}

func (dec *CellIDDecoder) Value(hit Hit) int64 {
	dec.bits.value = dec.value(hit)
	return dec.bits.value
}
func (dec *CellIDDecoder) ValueString(hit Hit) string {
	dec.bits.value = dec.value(hit)
	return dec.bits.valueString()
}

func (dec *CellIDDecoder) value(hit Hit) int64 {
	return int64(hit.GetCellID0())&0xffffffff | int64(hit.GetCellID1())<<32
}

type bitField64 struct {
	fields []bitFieldValue
	value  int64
	index  map[string]int
	joined int64
}

func newBitField64(codec string) *bitField64 {
	var bf bitField64
	toks := strings.Split(codec, ",")
	cur := 0
	for _, tok := range toks {
		subfields := strings.Split(tok, ":")
		var (
			field = bitFieldValue{bits: &bf.value}
			err   error
		)
		switch len(subfields) {
		case 2:
			field.name = subfields[0]
			field.width, err = strconv.Atoi(subfields[1])
			if err != nil {
				panic(err)
			}
			field.offset = cur
			cur += iabs(field.width)
		case 3:
			field.name = subfields[0]
			field.offset, err = strconv.Atoi(subfields[1])
			if err != nil {
				panic(err)
			}
			field.width, err = strconv.Atoi(subfields[2])
			if err != nil {
				panic(err)
			}
			cur = field.offset + iabs(field.width)
		default:
			panic(xerrors.Errorf("lcio: invalid number of subfields: %q", tok))
		}
		field.signed = field.width < 0
		field.width = iabs(field.width)
		field.mask = ((0x0001 << uint64(field.width)) - 1) << uint64(field.offset)
		bf.fields = append(bf.fields, field)
	}

	bf.index = make(map[string]int, len(bf.fields))
	for i, v := range bf.fields {
		bf.index[v.name] = i
	}
	return &bf
}

func (bf *bitField64) Description() string {
	o := new(bytes.Buffer)
	for i, v := range bf.fields {
		format := "%s:%d:%d"
		if i != 0 {
			format = "," + format
		}
		width := v.width
		if v.signed {
			width = -v.width
		}
		fmt.Fprintf(o, format, v.name, v.offset, width)
	}
	return string(o.Bytes())
}

func (bf *bitField64) valueString() string {
	o := new(bytes.Buffer)
	for i, v := range bf.fields {
		format := "%s:%d"
		if i != 0 {
			format = "," + format
		}
		fmt.Fprintf(o, format, v.name, v.value())
	}
	return string(o.Bytes())
}

type bitFieldValue struct {
	bits   *int64
	mask   int64
	name   string
	offset int
	width  int
	min    int
	max    int
	signed bool
}

func (bfv *bitFieldValue) value() int64 {
	bits := *bfv.bits
	val := bits & bfv.mask >> uint64(bfv.offset)

	if bfv.signed {
		if val&(1<<uint64(bfv.width-1)) != 0 { // negative value
			val -= 1 << uint64(bfv.width)
		}
	}
	return val
}

func iabs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
