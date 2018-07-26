// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"testing"
)

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
