// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio_test

import (
	"compress/flate"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func TestCreateRunHeader(t *testing.T) {
	const fname = "testdata/run-header.slcio"
	w, err := lcio.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	w.SetCompressionLevel(flate.NoCompression)

	rhdr := lcio.RunHeader{
		RunNbr:       42,
		Descr:        "a simple run header",
		Detector:     "my detector",
		SubDetectors: []string{"det-1", "det-2"},
		Params: lcio.Params{
			Floats: map[string][]float32{
				"floats-1": {1, 2, 3},
				"floats-2": {4, 5, 6},
			},
		},
	}

	err = w.WriteRunHeader(&rhdr)
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	chk, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile(strings.Replace(fname, ".slcio", "_golden.slcio", -1))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Fatalf("%s: differ with golden", fname)
	}

	os.Remove(fname)
}

func TestCreateCompressedRunHeader(t *testing.T) {
	const fname = "testdata/run-header-compressed.slcio"
	w, err := lcio.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	w.SetCompressionLevel(flate.BestCompression)

	rhdr := lcio.RunHeader{
		RunNbr:       42,
		Descr:        "a simple run header",
		Detector:     "my detector",
		SubDetectors: []string{"det-1", "det-2"},
		Params: lcio.Params{
			Floats: map[string][]float32{
				"floats-1": {1, 2, 3},
				"floats-2": {4, 5, 6},
			},
		},
	}

	err = w.WriteRunHeader(&rhdr)
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	chk, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile(strings.Replace(fname, ".slcio", "_golden.slcio", -1))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ref, chk) {
		t.Fatalf("%s: differ with golden", fname)
	}

	os.Remove(fname)
}

func TestCreate(t *testing.T) {
	const fname = "testdata/test.slcio"
	w, err := lcio.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	rhdr := lcio.RunHeader{
		RunNbr:       42,
		Descr:        "a simple run header",
		Detector:     "my detector",
		SubDetectors: []string{"det-1", "det-2"},
		Params: lcio.Params{
			Floats: map[string][]float32{
				"floats-1": {1, 2, 3},
				"floats-2": {4, 5, 6},
			},
		},
	}

	err = w.WriteRunHeader(&rhdr)
	if err != nil {
		t.Fatal(err)
	}

	evt := lcio.Event{
		RunNumber:   rhdr.RunNbr,
		Detector:    rhdr.Detector,
		EventNumber: 52,
		TimeStamp:   1234567890,
		Params: lcio.Params{
			Floats: map[string][]float32{
				"_weight": {42},
			},
			Strings: map[string][]string{
				"Descr": {"a description"},
			},
		},
	}

	if evt.Has("not-there") {
		t.Errorf("got an unexpected collection")
	}

	mcparts := lcio.McParticles{
		Flags: 0x1234ffff,
	}
	for i := 0; i < 3; i++ {
		i32 := int32(i+1) * 10
		f32 := float32(i+1) * 10
		f64 := float64(i+1) * 10
		mc := lcio.McParticle{
			PDG:       i32,
			Mass:      f64,
			Charge:    f32,
			P:         [3]float64{f64, f64, f64},
			PEndPoint: [3]float64{f64, f64, f64},
			GenStatus: int32(i + 1),
			SimStatus: 1 << 31,
			ColorFlow: [2]int32{i32, i32},
			Spin:      [3]float32{f32, f32, f32},
		}
		mcparts.Particles = append(mcparts.Particles, mc)
	}
	mcparts.Particles[1].Parents = []*lcio.McParticle{&mcparts.Particles[0], &mcparts.Particles[2]}
	mcparts.Particles[0].Children = []*lcio.McParticle{&mcparts.Particles[1]}
	mcparts.Particles[2].Children = []*lcio.McParticle{&mcparts.Particles[1]}

	simhits := lcio.SimCalorimeterHits{
		Flags: lcio.BitsChLong | lcio.BitsChID1 | lcio.BitsChStep | lcio.BitsChPDG,
		Hits: []lcio.SimCalorimeterHit{
			{
				CellID0: 1024, CellID1: 256, Energy: 42.666, Pos: [3]float32{1, 2, 3},
				Contributions: []lcio.Contrib{
					{Energy: 10, Mc: &mcparts.Particles[0]},
					{Energy: 11, Mc: &mcparts.Particles[1]},
				},
			},
			{
				CellID0: 1025, CellID1: 256, Energy: 42.666, Pos: [3]float32{1, 2, 3},
				Contributions: []lcio.Contrib{
					{Energy: 10, Mc: &mcparts.Particles[0]},
					{Energy: 11, Mc: &mcparts.Particles[1]},
					{Energy: 12, Mc: &mcparts.Particles[2]},
				},
			},
		},
	}

	calhits := lcio.CalorimeterHits{
		Flags: lcio.BitsRChLong | lcio.BitsRChID1 | lcio.BitsRChTime | lcio.BitsRChNoPtr | lcio.BitsRChEnergyError,
		Params: lcio.Params{
			Floats:  map[string][]float32{"f32": {1, 2, 3}},
			Ints:    map[string][]int32{"i32": {1, 2, 3}},
			Strings: map[string][]string{"str": {"1", "2", "3"}},
		},
		Hits: []lcio.CalorimeterHit{
			{
				CellID0:   1024,
				CellID1:   2048,
				Energy:    1000,
				EnergyErr: 0.1,
				Time:      1234,
				Pos:       [3]float32{11, 22, 33},
				Type:      42,
			},
		},
	}

	evt.Add("McParticles", &mcparts)
	evt.Add("SimCaloHits", &simhits)
	evt.Add("CaloHits", &calhits)

	if key := "McParticles"; !evt.Has(key) {
		t.Errorf("expected to have key %q", key)
	}

	err = w.WriteEvent(&evt)
	if err != nil {
		t.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	r, err := lcio.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	r.Next()
	if r.Err() != nil {
		t.Fatal(r.Err())
	}

	if got, want := r.RunHeader(), rhdr; !reflect.DeepEqual(got, want) {
		t.Fatalf("run-headers differ.\ngot= %#v\nwant=%#v\n", got, want)
	}

	if got, want := r.Event(), evt; !reflect.DeepEqual(got, want) {
		t.Fatalf("evts differ.\ngot:\n%v\nwant:\n%v\n", &got, &want)
	}

	os.Remove(fname)
}
