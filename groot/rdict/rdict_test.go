// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"math"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rmeta"
)

func TestElementGetRange(t *testing.T) {
	for _, tc := range []struct {
		name               string
		title              string
		rtype              rmeta.Enum
		xmin, xmax, factor float64
	}{
		{
			name:  "empty",
			title: "",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-d32",
			title: "var/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-f64",
			title: "var/D",
			rtype: rmeta.Float64,
		},
		{
			name:  "normal-f64-ndims",
			title: "var[10][20][30]/D",
			rtype: rmeta.Float64,
		},
		{
			name:  "normal-1d",
			title: "var[3]/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-2d",
			title: "var[3][4]/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-3d",
			title: "var[3][4][5]/d",
			rtype: rmeta.Double32,
		},
		{
			name:   "range",
			title:  "[ 0 , 100 ]",
			rtype:  rmeta.Double32,
			xmin:   0,
			xmax:   100,
			factor: float64(0xffffffff) / 100,
		},
		{
			name:   "range-ndim",
			title:  "var[3]/d[ 0 , 100 ]",
			rtype:  rmeta.Double32,
			xmin:   0,
			xmax:   100,
			factor: float64(0xffffffff) / 100,
		},
		{
			name:   "range-nbits",
			title:  "[ 10 , 100, 30 ]",
			rtype:  rmeta.Double32,
			xmin:   10,
			xmax:   100,
			factor: float64(1<<30) / 90,
		},
		{
			name:   "range-nbits-1d",
			title:  "var[3]/d[ 10 , 100, 30 ]",
			rtype:  rmeta.Double32,
			xmin:   10,
			xmax:   100,
			factor: float64(1<<30) / 90,
		},
		{
			name:   "range-pi",
			title:  "[ -pi , pi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi,
			xmax:   +math.Pi,
			factor: float64(0xffffffff) / (2 * math.Pi),
		},
		{
			name:   "range-pi/2",
			title:  "[ -pi/2 , 2pi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi / 2,
			xmax:   2 * math.Pi,
			factor: float64(0xffffffff) / (2*math.Pi + math.Pi/2),
		},
		{
			name:   "range-twopi/4",
			title:  "[ -pi/4 , twopi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi / 4,
			xmax:   2 * math.Pi,
			factor: float64(0xffffffff) / (2*math.Pi + math.Pi/4),
		},
		{
			name:   "range-2pi",
			title:  "[ -2*pi , 2*pi ]",
			rtype:  rmeta.Double32,
			xmin:   -2 * math.Pi,
			xmax:   +2 * math.Pi,
			factor: float64(0xffffffff) / (4 * math.Pi),
		},
		{
			name:  "float32-15bits",
			title: "[ 0 , 0 , 15 ]",
			rtype: rmeta.Double32,
		},
		{
			name:  "float32-14bits",
			title: "[ 0 , 0 , 14 ]",
			rtype: rmeta.Double32,
			xmin:  float64(14) + 0.1,
		},
		{
			name:  "float32-3bits",
			title: "[ 10 , 10 , 3 ]",
			rtype: rmeta.Double32,
			xmin:  float64(3) + 0.1,
			xmax:  10,
		},
		{
			name:  "float32-2bits",
			title: "[ 0 , 0 , 2 ]",
			rtype: rmeta.Double32,
			xmin:  float64(2) + 0.1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			elmt := Element{
				Name: *rbase.NewNamed(tc.name, tc.title),
				Type: tc.rtype,
			}
			elmt.parse()
			if got, want := elmt.XMin, tc.xmin; got != want {
				t.Fatalf("invalid xmin: got=%v, want=%v", got, want)
			}
			if got, want := elmt.XMax, tc.xmax; got != want {
				t.Fatalf("invalid xmax: got=%v, want=%v", got, want)
			}
			if got, want := elmt.Factor, tc.factor; got != want {
				t.Fatalf("invalid factor: got=%v, want=%v", got, want)
			}
		})
	}
}
