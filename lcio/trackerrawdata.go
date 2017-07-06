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

type TrackerRawDataContainer struct {
	Flags  Flags // bits 0-15 are user/detector specific
	Params Params
	Data   []TrackerRawData
}

type TrackerRawData struct {
	CellID0 int32
	CellID1 int32
	Time    int32
	ADCs    []uint16
}

func (data *TrackerRawData) GetCellID0() int32 { return data.CellID0 }
func (data *TrackerRawData) GetCellID1() int32 { return data.CellID1 }

func (tds *TrackerRawDataContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of TrackerData collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n", tds.Flags)
	fmt.Fprintf(o, "     LCIO::TRAWBIT_ID1 : %v\n", tds.Flags.Test(BitsTRawID1))

	fmt.Fprintf(o, "%v\n", tds.Params)

	const (
		head = " [   id   ] |  cellid0 |  cellid1 |   time   | cellid-fields  \n"
		tail = "------------|----------|----------|----------|----------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range tds.Data {
		data := &tds.Data[i]
		fmt.Fprintf(o,
			"[%09d] | %08d | %08d | %+8d | unknown/default ADCs: %v\n",
			ID(data),
			data.CellID0, data.CellID1,
			data.Time,
			data.ADCs,
		)
	}
	return string(o.Bytes())
}

func (*TrackerRawDataContainer) VersionSio() uint32 {
	return Version
}

func (trs *TrackerRawDataContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&trs.Flags)
	enc.Encode(&trs.Params)
	enc.Encode(int32(len(trs.Data)))
	for i := range trs.Data {
		data := &trs.Data[i]
		enc.Encode(&data.CellID0)
		if trs.Flags.Test(BitsTRawID1) {
			enc.Encode(&data.CellID1)
		}
		enc.Encode(&data.Time)
		nADCs := int32(len(data.ADCs))
		enc.Encode(&nADCs)
		for _, value := range data.ADCs {
			enc.Encode(&value)
		}
		if nADCs%2 == 1 {
			pad := uint16(0)
			enc.Encode(&pad)
		}
		enc.Tag(data)
	}
	return enc.Err()
}

func (trs *TrackerRawDataContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&trs.Flags)
	dec.Decode(&trs.Params)
	var n int32
	dec.Decode(&n)
	trs.Data = make([]TrackerRawData, int(n))
	for i := range trs.Data {
		data := &trs.Data[i]
		dec.Decode(&data.CellID0)
		if trs.Flags.Test(BitsTRawID1) {
			dec.Decode(&data.CellID1)
		}
		dec.Decode(&data.Time)
		var nADCs int32
		dec.Decode(&nADCs)
		data.ADCs = make([]uint16, nADCs)
		for j := range data.ADCs {
			dec.Decode(&data.ADCs[j])
		}
		if nADCs%2 == 1 {
			var pad uint16
			dec.Decode(&pad)
		}
		dec.Tag(data)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackerRawDataContainer)(nil)
	_ sio.Codec     = (*TrackerRawDataContainer)(nil)
)
