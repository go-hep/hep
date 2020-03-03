// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// decodeNameCycle decodes a namecycle "aap;2" into name "aap" and cycle "2"
func decodeNameCycle(namecycle string) (string, int16) {
	var name string
	var cycle int16

	toks := strings.Split(namecycle, ";")
	switch len(toks) {
	case 1:
		name = toks[0]
		cycle = 9999
	case 2:
		name = toks[0]
		i, err := strconv.Atoi(toks[1])
		if err != nil {
			// not a number
			cycle = 9999
		} else {
			cycle = int16(i)
		}
	default:
		panic(fmt.Errorf("invalid namecycle format [%v]", namecycle))
	}

	return name, cycle
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
		panic("riofs: TDatime year must be >= 1995")
	}

	return (year-1995)<<26 | month<<22 | day<<17 | hour<<12 | min<<6 | sec
}
