// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"codeberg.org/sbinet/npyio/npy"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func TestConvert(t *testing.T) {
	for _, tc := range []struct {
		name string
		want any
	}{
		// 4 scalars
		{
			name: "bool_4x1",
			want: [4]bool{true, true, false, true},
		},
		{
			name: "uint8_4x1",
			want: [4]uint8{0, 1, 2, 3},
		},
		{
			name: "uint16_4x1",
			want: [4]uint16{0, 1, 2, 3},
		},
		{
			name: "uint32_4x1",
			want: [4]uint32{0, 1, 2, 3},
		},
		{
			name: "uint64_4x1",
			want: [4]uint64{0, 1, 2, 3},
		},
		{
			name: "int8_4x1",
			want: [4]int8{0, 1, 2, 3},
		},
		{
			name: "int16_4x1",
			want: [4]int16{0, 1, 2, 3},
		},
		{
			name: "int32_4x1",
			want: [4]int32{0, 1, 2, 3},
		},
		{
			name: "int64_4x1",
			want: [4]int64{0, 1, 2, 3},
		},
		{
			name: "float32_4x1",
			want: [4]float32{0, 1, 2, 3},
		},
		{
			name: "float64_4x1",
			want: [4]float64{0, 1, 2, 3},
		},
		{
			name: "nans_4x1",
			want: [4]float64{math.Inf(-1), math.Inf(+1), math.NaN(), 0},
		},
		// 4 1d-arrays
		{
			name: "bool_4x2",
			want: [4][2]bool{{true, false}, {true, false}, {false, true}, {true, false}},
		},
		{
			name: "uint8_4x2",
			want: [4][2]uint8{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint16_4x2",
			want: [4][2]uint16{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint32_4x2",
			want: [4][2]uint32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint64_4x2",
			want: [4][2]uint64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int8_4x2",
			want: [4][2]int8{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int16_4x2",
			want: [4][2]int16{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int32_4x2",
			want: [4][2]int32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int64_4x2",
			want: [4][2]int64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "float32_4x2",
			want: [4][2]float32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "float64_4x2",
			want: [4][2]float64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		// 4 2d-arrays
		{
			name: "bool_4x3x2",
			want: [4][3][2]bool{
				{{true, false}, {true, false}, {true, false}},
				{{false, true}, {false, true}, {false, true}},
				{{true, false}, {true, false}, {true, false}},
				{{false, true}, {false, true}, {false, true}},
			},
		},
		{
			name: "uint8_4x3x2",
			want: [4][3][2]uint8{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "uint16_4x3x2",
			want: [4][3][2]uint16{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "uint32_4x3x2",
			want: [4][3][2]uint32{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "uint64_4x3x2",
			want: [4][3][2]uint64{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "int8_4x3x2",
			want: [4][3][2]int8{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "int16_4x3x2",
			want: [4][3][2]int16{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "int32_4x3x2",
			want: [4][3][2]int32{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "int64_4x3x2",
			want: [4][3][2]int64{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "float32_4x3x2",
			want: [4][3][2]float32{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name: "float64_4x3x2",
			want: [4][3][2]float64{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		// 3d-array
		{
			name: "float64_4x3x2x1",
			want: [4][3][2][1]float64{
				{{{10}, {11}}, {{12}, {13}}, {{14}, {15}}},
				{{{16}, {17}}, {{18}, {19}}, {{20}, {21}}},
				{{{22}, {23}}, {{24}, {25}}, {{26}, {27}}},
				{{{28}, {29}}, {{30}, {31}}, {{32}, {33}}},
			},
		},
		// 4d-array
		{
			name: "float64_4x3x2x1x2",
			want: [4][3][2][1][2]float64{
				{{{{10, 1}}, {{11, 2}}}, {{{12, 3}}, {{13, 4}}}, {{{14, 5}}, {{15, 6}}}},
				{{{{16, 1}}, {{17, 2}}}, {{{18, 3}}, {{19, 4}}}, {{{20, 5}}, {{21, 6}}}},
				{{{{22, 1}}, {{23, 2}}}, {{{24, 3}}, {{25, 4}}}, {{{26, 5}}, {{27, 6}}}},
				{{{{28, 1}}, {{29, 2}}}, {{{30, 3}}, {{31, 4}}}, {{{32, 5}}, {{33, 6}}}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp, err := os.MkdirTemp("", "npy2root-")
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer os.RemoveAll(tmp)

			fname := filepath.Join(tmp, "data.npy")
			src, err := os.Create(fname)
			if err != nil {
				t.Fatalf("could not create NumPy data file: %+v", err)
			}
			defer src.Close()

			err = npy.Write(src, tc.want)
			if err != nil {
				t.Fatalf("could not save NumPy data file: %+v", err)
			}

			err = src.Close()
			if err != nil {
				t.Fatalf("could not close NumPy data file: %+v", err)
			}

			oname := filepath.Join(tmp, "out.root")
			err = process(oname, "tree", fname)
			if err != nil {
				t.Fatalf("could not create ROOT data file: %+v", err)
			}

			f, err := groot.Open(oname)
			if err != nil {
				t.Fatalf("could not open ROOT file: %+v", err)
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatalf("could not get ROOT tree: %+v", err)
			}

			tree := obj.(rtree.Tree)
			rvars := rtree.NewReadVars(tree)
			r, err := rtree.NewReader(tree, rvars)
			if err != nil {
				t.Fatalf("could not create tree reader: %+v", err)
			}
			defer r.Close()

			want := reflect.ValueOf(tc.want)
			n := 0
			err = r.Read(func(ctx rtree.RCtx) error {
				i := int(ctx.Entry)
				want := want.Index(i).Interface()
				got := reflect.ValueOf(rvars[0].Value).Elem().Interface()
				ok := false
				switch want := want.(type) {
				case float64:
					got := got.(float64)
					switch {
					case math.IsNaN(want):
						ok = math.IsNaN(got)
					case math.IsInf(want, +1):
						ok = math.IsInf(got, +1)
					case math.IsInf(want, -1):
						ok = math.IsInf(got, -1)
					default:
						ok = got == want
					}
				default:
					ok = reflect.DeepEqual(got, want)
				}

				if !ok {
					return fmt.Errorf("invalid value for entry %d:\ngot= %v\nwant=%v", ctx.Entry, got, want)
				}
				n++
				return nil
			})
			if err != nil {
				t.Fatalf("could not read tree: %+v", err)
			}

			if got, want := n, want.Len(); got != want {
				t.Fatalf("invalid number of events: got=%d, want=%d", got, want)
			}
		})
	}
}
