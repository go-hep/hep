// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

type TrackerDataContainer struct {
	Flags  Flags // bits 0-15 are user/detector specific
	Params Params
	Data   []TrackerData
}

type TrackerData struct {
	CellID0 int32
	CellID1 int32
	Time    float32
	Charges []float32
}

func (data *TrackerData) GetCellID0() int32 { return data.CellID0 }
func (data *TrackerData) GetCellID1() int32 { return data.CellID1 }

func (tds *TrackerDataContainer) String() string {
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
			"[%09d] | %08d | %08d | %+.2e| unknown/default charge-ADC: %v\n",
			ID(data),
			data.CellID0, data.CellID1,
			data.Time,
			data.Charges,
		)
	}
	return string(o.Bytes())
}

func (*TrackerDataContainer) VersionSio() uint32 {
	return Version
}

func (tds *TrackerDataContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&tds.Flags)
	enc.Encode(&tds.Params)
	enc.Encode(int32(len(tds.Data)))
	for i := range tds.Data {
		data := &tds.Data[i]
		enc.Encode(&data.CellID0)
		if tds.Flags.Test(BitsTRawID1) {
			enc.Encode(&data.CellID1)
		}
		enc.Encode(&data.Time)
		enc.Encode(&data.Charges)
		enc.Tag(data)
	}
	return enc.Err()
}

func (tds *TrackerDataContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&tds.Flags)
	dec.Decode(&tds.Params)
	var n int32
	dec.Decode(&n)
	tds.Data = make([]TrackerData, int(n))
	for i := range tds.Data {
		data := &tds.Data[i]
		dec.Decode(&data.CellID0)
		if tds.Flags.Test(BitsTRawID1) {
			dec.Decode(&data.CellID1)
		}
		dec.Decode(&data.Time)
		dec.Decode(&data.Charges)
		dec.Tag(data)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackerDataContainer)(nil)
	_ sio.Codec     = (*TrackerDataContainer)(nil)
	_ Hit           = (*TrackerData)(nil)
)
