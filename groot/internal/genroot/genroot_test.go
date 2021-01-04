// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genroot

import "testing"

func TestExtractYear(t *testing.T) {
	for _, tc := range []struct {
		fname string
		year  int
	}{
		{
			fname: "./genroot.go",
			year:  2018,
		},
		{
			fname: "./genroot_test.go",
			year:  2020,
		},
		{
			fname: "../../cmd/root-gen-rfunc/testdata/func1_golden.txt",
			year:  2021,
		},
	} {
		t.Run(tc.fname, func(t *testing.T) {
			year := ExtractYear(tc.fname)
			if year != tc.year {
				t.Fatalf("invalid year: got=%d, want=%d", year, tc.year)
			}
		})
	}
}
