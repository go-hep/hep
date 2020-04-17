// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

type TrackerPulseContainer struct {
	Flags  Flags
	Params Params
	Pulses []TrackerPulse
}

type TrackerPulse struct {
	CellID0 int32
	CellID1 int32
	Time    float32      // time of pulse
	Charge  float32      // charge of pulse
	Cov     [3]float32   // covariance matrix of charge (c) and time (t) measurements
	Quality int32        // quality flag word
	TPC     *TrackerData // TPC corrected data: spectrum used to create this pulse
}

func (hit *TrackerPulse) GetCellID0() int32 { return hit.CellID0 }
func (hit *TrackerPulse) GetCellID1() int32 { return hit.CellID1 }

func (ps *TrackerPulseContainer) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%[1]s print out of TrackerPulse collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "  flag:  0x%x\n", ps.Flags)
	fmt.Fprintf(o, "     LCIO::TRAWBIT_ID1 : %v\n", ps.Flags.Test(BitsTRawID1))
	fmt.Fprintf(o, "     LCIO::TRAWBIT_CM  : %v\n", ps.Flags.Test(BitsTRawCM))

	fmt.Fprintf(o, "%v\n", ps.Params)

	const (
		head = " [   id   ] | cellid0  | cellid1  |  time | charge | quality  | [corr.Data] |  cellid-fields: | cov(c,c), cov(t,c), cov(t,t) \n"
		tail = "------------|----------|----------|-------|--------|----------|-------------|-----------------|------------------------------\n"
	)
	fmt.Fprintf(o, head)
	fmt.Fprintf(o, tail)
	for i := range ps.Pulses {
		p := &ps.Pulses[i]
		fmt.Fprintf(o,
			"[%09d] | %08d | %08d |%+.4f| %+.4f| %8d | [%09d] | unknown/default | %+.2e, %+.2e, %.2e|\n",
			ID(p),
			p.CellID0, p.CellID1,
			p.Time, p.Charge,
			p.Quality,
			ID(p.TPC),
			p.Cov[0], p.Cov[1], p.Cov[2],
		)
	}
	return string(o.Bytes())
}

func (*TrackerPulseContainer) VersionSio() uint32 {
	return Version
}

func (ps *TrackerPulseContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&ps.Flags)
	enc.Encode(&ps.Params)
	enc.Encode(int32(len(ps.Pulses)))
	for i := range ps.Pulses {
		p := &ps.Pulses[i]
		enc.Encode(&p.CellID0)
		if ps.Flags.Test(BitsTRawID1) {
			enc.Encode(&p.CellID1)
		}
		enc.Encode(&p.Time)
		enc.Encode(&p.Charge)
		if ps.Flags.Test(BitsTRawCM) {
			enc.Encode(&p.Cov)
		}
		enc.Encode(&p.Quality)
		enc.Pointer(&p.TPC)
		enc.Tag(p)
	}
	return enc.Err()
}

func (ps *TrackerPulseContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&ps.Flags)
	dec.Decode(&ps.Params)
	var n int32
	dec.Decode(&n)
	ps.Pulses = make([]TrackerPulse, int(n))
	for i := range ps.Pulses {
		p := &ps.Pulses[i]
		dec.Decode(&p.CellID0)
		if ps.Flags.Test(BitsTRawID1) {
			dec.Decode(&p.CellID1)
		}
		dec.Decode(&p.Time)
		dec.Decode(&p.Charge)
		if r.VersionSio() > 1012 && ps.Flags.Test(BitsTRawCM) {
			dec.Decode(&p.Cov)
		}
		dec.Decode(&p.Quality)
		dec.Pointer(&p.TPC)
		dec.Tag(p)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*TrackerPulseContainer)(nil)
	_ sio.Codec     = (*TrackerPulseContainer)(nil)
	_ Hit           = (*TrackerPulse)(nil)
)
