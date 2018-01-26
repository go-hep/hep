// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hifact

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestBinData(t *testing.T) {
	for _, fname := range []string{
		"testdata/1bin_example1.json",
		"testdata/2bin_example1.json",
	} {
		t.Run(fname, func(t *testing.T) {
			raw, err := ioutil.ReadFile(fname)
			if err != nil {
				t.Fatal(err)
			}
			var data Data
			err = json.Unmarshal(raw, &data)
			if err != nil {
				t.Fatal(err)
			}
			log.Printf("%#v", data)
		})
	}
}

func TestHistoSys(t *testing.T) {
	src := Data{
		Binning: []float64{2, -0.5, 1.5},
		BinData: BinData{
			Data:       []float64{120.0, 180.0},
			Bkg:        []float64{100.0, 150.0},
			BkgSysUp:   []float64{102, 190},
			BkgSysDown: []float64{98, 100},
			Sig:        []float64{30.0, 95.0},
		},
	}
	samples := []Sample{
		{
			Data: src.BinData.Sig,
			Models: []ModelDescr{
				{
					Name: "mu",
					Type: NormFactorModel,
				},
			},
		},
		{
			Data: src.BinData.Bkg,
			Models: []ModelDescr{
				{
					Name: "bkg_norm",
					Type: HistoSysModel,
					Data: map[string][]float64{
						"lo_hist": src.BinData.BkgSysDown,
						"hi_hist": src.BinData.BkgSysUp,
					},
				},
			},
		},
	}
	pdf := New(samples)
	log.Printf("%#v", pdf)
	pars := []float64{0, 0}
	bounds := [][]float64{{0, 10}, {-5, 5}}
	log.Print(pars, bounds)
}
