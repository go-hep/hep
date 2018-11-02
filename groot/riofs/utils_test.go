// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"testing"
	"time"
)

func TestDatime2time(t *testing.T) {
	for _, tc := range []struct {
		datime uint32
		want   time.Time
	}{
		{
			datime: 1576331001,
			want:   time.Date(2018, time.July, 26, 14, 27, 57, 0, time.UTC),
		},
		{
			datime: 1576331001,
			want:   time.Date(2018, time.July, 26, 14, 27, 57, 0, time.UTC),
		},
		{
			datime: 347738243,
			want:   time.Date(2000, time.February, 29, 1, 2, 3, 0, time.UTC),
		},
		{
			datime: 4325376,
			want:   time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	} {
		t.Run(tc.want.String(), func(t *testing.T) {
			got := datime2time(tc.datime)
			if !got.Equal(tc.want) {
				t.Fatalf("got=%v. want=%v", got, tc.want)
			}

			datime := time2datime(tc.want)
			if datime != tc.datime {
				t.Fatalf("got=%v. want=%v", datime, tc.datime)
			}
		})
	}
}

func TestTime2Datime(t *testing.T) {
	for _, tc := range []struct {
		time    time.Time
		want    uint32
		panicks bool
	}{
		{
			time: time.Date(2018, time.July, 26, 14, 27, 57, 0, time.UTC),
			want: 1576331001,
		},
		{
			time: time.Date(2000, time.February, 29, 1, 2, 3, 0, time.UTC),
			want: 347738243,
		},
		{
			time: time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: 4325376,
		},
		{
			time:    time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
			panicks: true,
		},
	} {
		t.Run(tc.time.String(), func(t *testing.T) {
			if tc.panicks {
				defer func() {
					if recover() == nil {
						t.Fatalf("should have panicked.")
					}
				}()
			}
			{
				got := time2datime(tc.time)
				if got != tc.want {
					t.Fatalf("got=%v. want=%v", got, tc.want)
				}
			}

			{
				got := datime2time(tc.want)
				if !got.Equal(tc.time) {
					t.Fatalf("got=%v. want=%v", got, tc.time)
				}
			}
		})
	}
}

func TestDecodeNameCycle(t *testing.T) {
	for _, tc := range []struct {
		nc     string
		name   string
		cycle  int16
		panics bool
	}{
		{
			nc:    "name",
			name:  "name",
			cycle: 9999,
		},
		{
			nc:    "name;0",
			name:  "name",
			cycle: 0,
		},
		{
			nc:    "name;1",
			name:  "name",
			cycle: 1,
		},
		{
			nc:    "name;42",
			name:  "name",
			cycle: 42,
		},
		{
			nc:    "name;42.0",
			name:  "name",
			cycle: 9999,
		},
		{
			nc:    "name;e",
			name:  "name",
			cycle: 9999,
		},
		{
			nc:     "nam;e;1",
			name:   "name",
			cycle:  0,
			panics: true,
		},
	} {
		t.Run(tc.nc, func(t *testing.T) {
			if tc.panics {
				defer func() {
					e := recover()
					if e == nil {
						t.Fatalf("should have panicked.")
					}
				}()
			}
			n, c := decodeNameCycle(tc.nc)
			if n != tc.name {
				t.Fatalf("got=%q. want=%q", n, tc.name)
			}
			if c != tc.cycle {
				t.Fatalf("got=%d. want=%d", c, tc.cycle)
			}
		})
	}
}
