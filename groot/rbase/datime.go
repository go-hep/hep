// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"fmt"
	"reflect"
	"time"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// Datime is a ROOT date + time.
// Note that ROOT's TDatime is relative to 1995.
type Datime time.Time

func (*Datime) Class() string {
	return "TDatime"
}

func (*Datime) RVersion() int16 {
	return rvers.Datime
}

// MarshalROOT implements rbytes.Marshaler
func (dt *Datime) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	// TDatime does not write a version header.
	w.WriteU32(time2datime(time.Time(*dt)))

	return 4, w.Err()
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (dt *Datime) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	// TDatime does not write a version header.
	*dt = Datime(datime2time(r.ReadU32()))

	return r.Err()
}

func (dt Datime) String() string {
	return dt.Time().String()
}

func (dt Datime) Time() time.Time { return time.Time(dt) }

func init() {
	f := func() reflect.Value {
		var o Datime
		return reflect.ValueOf(&o)
	}
	rtypes.Factory.Add("TDatime", f)
}

// time2datime converts a time.Time to a uint32 representing a ROOT's TDatime.
func time2datime(t time.Time) uint32 {
	var (
		year  = uint32(t.Year())
		month = uint32(t.Month())
		day   = uint32(t.Day())
		hour  = uint32(t.Hour())
		min   = uint32(t.Minute())
		sec   = uint32(t.Second())
	)

	if year < 1995 {
		panic(fmt.Errorf("rbase: TDatime year must be >= 1995"))
	}

	return (year-1995)<<26 | month<<22 | day<<17 | hour<<12 | min<<6 | sec
}

// datime2time converts a uint32 holding a ROOT's TDatime into a time.Time
func datime2time(d uint32) time.Time {

	// ROOT's TDatime begins in January 1995...
	var year uint32 = (d >> 26) + 1995
	var month uint32 = (d << 6) >> 28
	var day uint32 = (d << 10) >> 27
	var hour uint32 = (d << 15) >> 27
	var min uint32 = (d << 20) >> 26
	var sec uint32 = (d << 26) >> 26
	nsec := 0
	return time.Date(int(year), time.Month(month), int(day),
		int(hour), int(min), int(sec), nsec, time.UTC)
}
