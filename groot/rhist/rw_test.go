// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestWRBuffer(t *testing.T) {
	loadFrom := func(fname, key string) rtests.ROOTer {
		f, err := riofs.Open(fname)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		obj, err := riofs.Dir(f).Get(key)
		if err != nil {
			t.Fatal(err)
		}
		return obj.(rtests.ROOTer)
	}

	for _, tc := range []struct {
		name string
		want rtests.ROOTer
	}{
		{
			name: "TH1F",
			want: &H1F{
				th1: th1{
					Named:     *rbase.NewNamed("h1f", "my-title"),
					attline:   rbase.AttLine{Color: 602, Style: 1, Width: 1},
					attfill:   rbase.AttFill{Color: 0, Style: 1001},
					attmarker: rbase.AttMarker{Color: 1, Style: 1, Width: 1},
					ncells:    102,
					xaxis: taxis{
						Named: *rbase.NewNamed("xaxis", ""),
						attaxis: rbase.AttAxis{
							Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.005, LabelSize: 0.035, Ticks: 0.03, TitleOffset: 1, TitleSize: 0.035, TitleColor: 1, TitleFont: 42,
						},
						nbins: 100, xmin: 0, xmax: 100,
						xbins: rcont.ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					yaxis: taxis{
						Named: *rbase.NewNamed("yaxis", ""),
						attaxis: rbase.AttAxis{
							Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.005, LabelSize: 0.035, Ticks: 0.03, TitleOffset: 1, TitleSize: 0.035, TitleColor: 1, TitleFont: 42,
						},
						nbins: 1, xmin: 0, xmax: 1,
						xbins: rcont.ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					zaxis: taxis{
						Named: *rbase.NewNamed("zaxis", ""),
						attaxis: rbase.AttAxis{
							Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.005, LabelSize: 0.035, Ticks: 0.03, TitleOffset: 1, TitleSize: 0.035, TitleColor: 1, TitleFont: 42,
						},
						nbins: 1, xmin: 0, xmax: 1,
						xbins: rcont.ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					boffset: 0, bwidth: 1000,
					entries: 10,
					tsumw:   10, tsumw2: 16, tsumwx: 278, tsumwx2: 9286,
					max: -1111, min: -1111,
					norm:    0,
					contour: rcont.ArrayD{Data: nil},
					sumw2: rcont.ArrayD{
						Data: []float64{
							1,
							0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
							9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0,
							1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
							0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
							0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
							1,
						},
					},
					opt:    "",
					funcs:  *rcont.NewList("", []root.Object{}),
					buffer: nil,
					erropt: 0,
				},
				arr: rcont.ArrayF{
					Data: []float32{
						1,
						0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						1,
					},
				},
			},
		},
		{
			name: "TH2F",
			want: &H2F{
				th2: th2{
					th1: th1{
						Named:     *rbase.NewNamed("h2f", "my title"),
						attline:   rbase.AttLine{Color: 602, Style: 1, Width: 1},
						attfill:   rbase.AttFill{Color: 0, Style: 1001},
						attmarker: rbase.AttMarker{Color: 1, Style: 1, Width: 1},
						ncells:    144,
						xaxis: taxis{
							Named: *rbase.NewNamed("xaxis", ""),
							attaxis: rbase.AttAxis{
								Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.004999999888241291, LabelSize: 0.03500000014901161,
								Ticks: 0.029999999329447746, TitleOffset: 1, TitleSize: 0.03500000014901161, TitleColor: 1, TitleFont: 42,
							},
							nbins:   10,
							xmin:    0,
							xmax:    10,
							xbins:   rcont.ArrayD{},
							first:   0,
							last:    0,
							bits2:   0x0,
							time:    false,
							tfmt:    "",
							labels:  nil,
							modlabs: nil,
						},
						yaxis: taxis{
							Named: *rbase.NewNamed("yaxis", ""),
							attaxis: rbase.AttAxis{
								Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.004999999888241291, LabelSize: 0.03500000014901161,
								Ticks: 0.029999999329447746, TitleOffset: 1, TitleSize: 0.03500000014901161, TitleColor: 1, TitleFont: 42,
							},
							nbins:   10,
							xmin:    0,
							xmax:    10,
							xbins:   rcont.ArrayD{},
							first:   0,
							last:    0,
							bits2:   0x0,
							time:    false,
							tfmt:    "",
							labels:  nil,
							modlabs: nil,
						},
						zaxis: taxis{
							Named: *rbase.NewNamed("zaxis", ""),
							attaxis: rbase.AttAxis{
								Ndivs: 510, AxisColor: 1, LabelColor: 1, LabelFont: 42, LabelOffset: 0.004999999888241291, LabelSize: 0.03500000014901161,
								Ticks: 0.029999999329447746, TitleOffset: 1, TitleSize: 0.03500000014901161, TitleColor: 1, TitleFont: 42,
							},
							nbins:   1,
							xmin:    0,
							xmax:    1,
							xbins:   rcont.ArrayD{},
							first:   0,
							last:    0,
							bits2:   0x0,
							time:    false,
							tfmt:    "",
							labels:  nil,
							modlabs: nil,
						},
						boffset: 0,
						bwidth:  1000,
						entries: 13,
						tsumw:   9,
						tsumw2:  29,
						tsumwx:  21,
						tsumwx2: 55,
						max:     -1111,
						min:     -1111,
						norm:    0,
						contour: rcont.ArrayD{},
						sumw2: rcont.ArrayD{
							Data: []float64{
								1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2,
								0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 25, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
								0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1,
							},
						},
						opt:    "",
						funcs:  *rcont.NewList("", []root.Object{}),
						buffer: nil,
						erropt: 0,
					},
					scale:   1,
					tsumwy:  21,
					tsumwy2: 55,
					tsumwxy: 55,
				},
				arr: rcont.ArrayF{
					Data: []float32{
						1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0,
						0, 1, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
						0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0,
						0, 0, 0, 0, 0, 1,
					},
				},
			},
		},
		{
			name: "TConfidenceLevel",
			want: &ConfidenceLevel{
				base:   *rbase.NewObject(),
				fNNMC:  3,
				fDtot:  5,
				fStot:  2,
				fBtot:  3,
				fTSD:   10,
				fNMC:   11,
				fMCL3S: 12,
				fMCL5S: 13,
				fTSB:   []float64{11, 12, 13},
				fTSS:   []float64{21, 22, 23},
				fLRS:   []float64{31, 32, 33},
				fLRB:   []float64{41, 42, 43},
				fISB:   []int32{1, 2, 3},
				fISS:   []int32{4, 5, 6},
			},
		},
		{
			name: "TLimit",
			want: &Limit{},
		},
		{
			name: "TLimitDataSource",
			want: &LimitDataSource{
				base:     *rbase.NewObject(),
				sig:      newObjArray("sig", "sig1"),
				bkg:      newObjArray("bkg", "bkg1"),
				data:     newObjArray("data", "data1"),
				sigErr:   newObjArray("0", "1", "2"),
				bkgErr:   newObjArray("3", "4", "5"),
				ids:      newObjArray("6", "7", "8"),
				dummyTA:  newObjArray("00", "11", "22"),
				dummyIDs: newObjArray("11", "22", "33"),
			},
		},
		{
			name: "TEfficiency",
			want: loadFrom("../testdata/tconfidence-level.root", "eff"),
		},
		{
			name: "TProfile",
			want: loadFrom("../testdata/tprofile.root", "p1d"),
		},
		{
			name: "TProfile2D",
			want: loadFrom("../testdata/tprofile.root", "p2d"),
		},
		{
			name: "TGraphMultiErrors",
			want: loadFrom("../testdata/tgme.root", "gme"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			class := tc.want.Class()
			obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
			{
				rbuf.SetErr(io.EOF)
				err = obj.UnmarshalROOT(rbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
				rbuf.SetErr(nil)
			}
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}

func TestReadF1(t *testing.T) {
	f, err := riofs.Open("../testdata/tformula.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, key := range []string{
		"func1", "func2", "func3", "func4",
		"fconv",
		"fnorm",
	} {
		t.Run(key, func(t *testing.T) {
			obj, err := f.Get(key)
			if err != nil {
				t.Fatalf("could not read object %q: %+v", key, err)
			}
			switch v := obj.(type) {
			case *F1:
				if got, want := v.Name(), key; got != want {
					t.Fatalf("invalid name: got=%q, want=%q", got, want)
				}
				if got, want := v.Class(), "TF1"; got != want {
					t.Fatalf("invalid class: got=%q, want=%q", got, want)
				}
				if got, want := v.chi2, 0.2; got != want {
					t.Fatalf("invalid chi2: got=%v, want=%v", got, want)
				}
			case F1Composition:
				// ok.
			default:
				t.Fatalf("invalid object type for %q", key)
			}
		})
	}
}

func newObjArray(vs ...string) rcont.ObjArray {
	elems := make([]root.Object, len(vs))
	for i, v := range vs {
		elems[i] = rbase.NewObjString(v)
	}
	o := rcont.NewObjArray()
	o.SetElems(elems)
	return *o
}
