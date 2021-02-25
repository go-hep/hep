// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.12
// +build go1.12

package hep

import (
	"reflect"
	"runtime/debug"
	"testing"
)

func TestVersion(t *testing.T) {
	type vsum struct {
		version, sum string
	}

	for _, tc := range []struct {
		b    *debug.BuildInfo
		want vsum
	}{
		{
			b:    nil,
			want: vsum{"", ""},
		},
		{
			b:    &debug.BuildInfo{},
			want: vsum{"", ""},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
				},
			},
			want: vsum{"", ""},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
					&debug.Module{
						Path:    "go-hep.org/x/hep",
						Version: "v0.18.0",
						Sum:     "12345",
					},
				},
			},
			want: vsum{"v0.18.0", "12345"},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
					&debug.Module{
						Path:    "go-hep.org/x/hep",
						Version: "v0.18.0",
						Sum:     "12345",
						Replace: &debug.Module{
							Path:    "go-hep.org/x/hep-fixup",
							Version: "v0.18.0-fixup",
							Sum:     "11111",
						},
					},
				},
			},
			want: vsum{"go-hep.org/x/hep-fixup v0.18.0-fixup", "11111"},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
					&debug.Module{
						Path:    "go-hep.org/x/hep",
						Version: "v0.18.0",
						Sum:     "12345",
						Replace: &debug.Module{
							Version: "v0.18.0-fixup",
							Sum:     "11111",
						},
					},
				},
			},
			want: vsum{"v0.18.0-fixup", "11111"},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
					&debug.Module{
						Path:    "go-hep.org/x/hep",
						Version: "v0.18.0",
						Sum:     "12345",
						Replace: &debug.Module{
							Path: "go-hep.org/x/hep-fixup",
							Sum:  "11111",
						},
					},
				},
			},
			want: vsum{"go-hep.org/x/hep-fixup", "11111"},
		},
		{
			b: &debug.BuildInfo{
				Deps: []*debug.Module{
					&debug.Module{
						Path:    "gonum.org/v1/gonum",
						Version: "v0.1.0",
						Sum:     "12345XYZ",
					},
					&debug.Module{
						Path:    "go-hep.org/x/hep",
						Version: "v0.18.0",
						Sum:     "12345",
						Replace: &debug.Module{
							Sum: "11111",
						},
					},
				},
			},
			want: vsum{"v0.18.0*", ""},
		},
	} {
		t.Run("", func(t *testing.T) {
			version, sum := versionOf(tc.b)
			got := vsum{
				version: version,
				sum:     sum,
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error:\ngot= %#v\nwant=%#v\n", got, tc.want)
			}
		})
	}

	vers, sum := Version()
	if got, want := (vsum{vers, sum}), (vsum{}); !reflect.DeepEqual(got, want) {
		// this might fail when GO111MODULE=on and module support percolates everywhere.
		t.Fatalf("error:\ngot= %#v\nwant=%#v\n", got, want)
	}
}
