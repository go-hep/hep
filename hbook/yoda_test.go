// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"fmt"
	"testing"
)

func TestReadYODAHeader(t *testing.T) {
	const mark = "BEGIN YODA_HISTO1D"
	for _, tc := range []struct {
		str  string
		want string
		vers int
		err  error
	}{
		{
			str:  "BEGIN YODA_HISTO1D /name\n",
			want: "/name",
			vers: 1,
		},
		{
			str:  "BEGIN YODA_HISTO1D /name with whitespace\n",
			want: "/name with whitespace",
			vers: 1,
		},
		{
			str:  "BEGIN YODA_HISTO1D_V2 /name\n",
			want: "/name",
			vers: 2,
		},
		{
			str:  "BEGIN YODA_HISTO1D_V2 /name with whitespace\n",
			want: "/name with whitespace",
			vers: 2,
		},
		{
			str:  "BEGIN YODA /name",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s line", mark),
		},
		{
			str:  "BEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  "\nBEGIN YODA /name",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  "\nBEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  " BEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			name, vers, err := readYODAHeader(newRBuffer([]byte(tc.str)), mark)
			if err == nil && tc.err != nil {
				t.Fatalf("got err=nil, want=%v", tc.err.Error())
			}
			if err != nil && tc.err == nil {
				t.Fatalf("got=%v, want=nil", err.Error())
			}
			if err != nil && tc.err != nil {
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("got error=%v, want=%v", got, want)
				}
			}
			if got, want := name, tc.want; got != want {
				t.Fatalf("invalid name: got: %q, want: %q", got, want)
			}
			if got, want := vers, tc.vers; got != want {
				t.Fatalf("invalid version: got: %d, want: %d", got, want)
			}
		})
	}
}
