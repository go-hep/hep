// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestRWHist(t *testing.T) {

	rootls := "rootls"
	if runtime.GOOS == "windows" {
		rootls = "rootls.exe"
	}

	rootls, err := exec.LookPath(rootls)
	withROOTCxx := err == nil

	dir, err := ioutil.TempDir("", "rootio-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	type histoer interface {
		Object
		ROOTMarshaler
		ROOTUnmarshaler
	}
	for i, tc := range []struct {
		name string
		skip bool
		want histoer
	}{
		{
			name: "TH1F",
			want: &H1F{
				rvers: 2,
				th1: th1{
					rvers:     7,
					tnamed:    tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000008}, name: "h1f", title: "my-title"},
					attline:   attline{rvers: 2, color: 602, style: 1, width: 1},
					attfill:   attfill{rvers: 2, color: 0, style: 1001},
					attmarker: attmarker{rvers: 2, color: 1, style: 1, width: 1},
					ncells:    102,
					xaxis: taxis{
						rvers:  10,
						tnamed: tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "xaxis", title: ""},
						attaxis: attaxis{
							rvers: 4,
							ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.005, lsize: 0.035, ticks: 0.03, toffset: 1, tsize: 0.035, tcolor: 1, tfont: 42,
						},
						nbins: 100, xmin: 0, xmax: 100,
						xbins: ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					yaxis: taxis{
						rvers:  10,
						tnamed: tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "yaxis", title: ""},
						attaxis: attaxis{
							rvers: 4,
							ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.005, lsize: 0.035, ticks: 0.03, toffset: 1, tsize: 0.035, tcolor: 1, tfont: 42,
						},
						nbins: 1, xmin: 0, xmax: 1,
						xbins: ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					zaxis: taxis{
						rvers:  10,
						tnamed: tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "zaxis", title: ""},
						attaxis: attaxis{
							rvers: 4,
							ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.005, lsize: 0.035, ticks: 0.03, toffset: 1, tsize: 0.035, tcolor: 1, tfont: 42,
						},
						nbins: 1, xmin: 0, xmax: 1,
						xbins: ArrayD{Data: nil},
						first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
						labels:  nil,
						modlabs: nil,
					},
					boffset: 0, bwidth: 1000,
					entries: 10,
					tsumw:   10, tsumw2: 16, tsumwx: 278, tsumwx2: 9286,
					max: -1111, min: -1111,
					norm:    0,
					contour: ArrayD{Data: nil},
					sumw2: ArrayD{
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
					opt: "",
					funcs: tlist{
						rvers: 5,
						obj:   tobject{id: 0x0, bits: 0x3000000},
						name:  "", objs: []Object{},
					},
					buffer: nil,
					erropt: 0,
				},
				arr: ArrayF{
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
				rvers: 3,
				th2: th2{
					rvers: 4,
					th1: th1{
						rvers: 7,
						tnamed: tnamed{
							rvers: 1,
							obj:   tobject{id: 0x0, bits: 0x3000008},
							name:  "h2f",
							title: "my title",
						},
						attline:   attline{rvers: 2, color: 602, style: 1, width: 1},
						attfill:   attfill{rvers: 2, color: 0, style: 1001},
						attmarker: attmarker{rvers: 2, color: 1, style: 1, width: 1},
						ncells:    144,
						xaxis: taxis{
							rvers: 10,
							tnamed: tnamed{
								rvers: 1,
								obj:   tobject{id: 0x0, bits: 0x3000000},
								name:  "xaxis",
								title: "",
							},
							attaxis: attaxis{
								rvers: 4,
								ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.004999999888241291, lsize: 0.03500000014901161,
								ticks: 0.029999999329447746, toffset: 1, tsize: 0.03500000014901161, tcolor: 1, tfont: 42,
							},
							nbins:   10,
							xmin:    0,
							xmax:    10,
							xbins:   ArrayD{},
							first:   0,
							last:    0,
							bits2:   0x0,
							time:    false,
							tfmt:    "",
							labels:  nil,
							modlabs: nil,
						},
						yaxis: taxis{
							rvers: 10,
							tnamed: tnamed{
								rvers: 1,
								obj:   tobject{id: 0x0, bits: 0x3000000},
								name:  "yaxis",
								title: "",
							},
							attaxis: attaxis{
								rvers: 4,
								ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.004999999888241291, lsize: 0.03500000014901161,
								ticks: 0.029999999329447746, toffset: 1, tsize: 0.03500000014901161, tcolor: 1, tfont: 42,
							},
							nbins:   10,
							xmin:    0,
							xmax:    10,
							xbins:   ArrayD{},
							first:   0,
							last:    0,
							bits2:   0x0,
							time:    false,
							tfmt:    "",
							labels:  nil,
							modlabs: nil,
						},
						zaxis: taxis{
							rvers: 10,
							tnamed: tnamed{
								rvers: 1,
								obj:   tobject{id: 0x0, bits: 0x3000000},
								name:  "zaxis",
								title: "",
							},
							attaxis: attaxis{
								rvers: 4,
								ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.004999999888241291, lsize: 0.03500000014901161,
								ticks: 0.029999999329447746, toffset: 1, tsize: 0.03500000014901161, tcolor: 1, tfont: 42,
							},
							nbins:   1,
							xmin:    0,
							xmax:    1,
							xbins:   ArrayD{},
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
						contour: ArrayD{},
						sumw2: ArrayD{
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
						opt: "",
						funcs: tlist{
							rvers: 5,
							obj:   tobject{id: 0x0, bits: 0x3000000},
							name:  "",
							objs:  []Object{},
						},
						buffer: nil,
						erropt: 0,
					},
					scale:   1,
					tsumwy:  21,
					tsumwy2: 55,
					tsumwxy: 55,
				},
				arr: ArrayF{
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
	} {
		fname := filepath.Join(dir, fmt.Sprintf("histos-%d.root", i))
		t.Run(tc.name, func(t *testing.T) {
			const kname = "my-key"

			w, err := Create(fname)
			if err != nil {
				t.Fatal(err)
			}

			err = w.Put(kname, tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(w.Keys()), 1; got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			err = w.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			r, err := Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			si := r.StreamerInfos()
			if len(si) == 0 {
				t.Fatalf("empty list of streamers")
			}

			if got, want := len(r.Keys()), 1; got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			rgot, err := r.Get(kname)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := rgot.(histoer), tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("error reading back objstring.\ngot = %#v\nwant = %#v", got, want)
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			if !withROOTCxx {
				t.Logf("skip test with ROOT/C++")
				return
			}

			cmd := exec.Command(rootls, "-l", fname)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("ROOT/C++ could not open file %q", fname)
			}
		})
	}
}
