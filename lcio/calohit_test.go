// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio_test

import (
	"compress/flate"
	"reflect"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func TestRWCalo(t *testing.T) {
	for _, test := range []struct {
		fname   string
		complvl int
	}{
		{"testdata/calohit.slcio", flate.NoCompression},
		{"testdata/calohit-compressed.slcio", flate.BestCompression},
	} {
		testRWCalo(t, test.complvl, test.fname)
	}
}

func testRWCalo(t *testing.T, compLevel int, fname string) {
	w, err := lcio.Create(fname)
	if err != nil {
		t.Error(err)
		return
	}
	defer w.Close()
	w.SetCompressionLevel(compLevel)

	const (
		nevents    = 10
		nhits      = 100
		CALHITS    = "CalorimeterHits"
		CALHITSERR = "CalorimeterHitsWithEnergyError"
	)

	for i := 0; i < nevents; i++ {
		evt := lcio.Event{
			RunNumber:   4711,
			EventNumber: int32(i),
		}

		var (
			calhits    = lcio.CalorimeterHitContainer{Flags: lcio.BitsRChLong}
			calhitsErr = lcio.CalorimeterHitContainer{Flags: lcio.BitsRChLong | lcio.BitsRChEnergyError}
		)
		for j := 0; j < nhits; j++ {
			hit := lcio.CalorimeterHit{
				CellID0:   int32(i*100000 + j),
				Energy:    float32(i*j) + 117,
				EnergyErr: float32(i*j) * 0.117,
				Pos:       [3]float32{float32(i), float32(j), float32(i * j)},
			}
			calhits.Hits = append(calhits.Hits, hit)

			hit.EnergyErr = float32(i*j) * 0.117
			calhitsErr.Hits = append(calhitsErr.Hits, hit)
		}

		evt.Add(CALHITS, &calhits)
		evt.Add(CALHITSERR, &calhitsErr)

		err = w.WriteEvent(&evt)
		if err != nil {
			t.Errorf("%s: error writing event %d: %v", fname, i, err)
			return
		}
	}

	err = w.Close()
	if err != nil {
		t.Errorf("%s: error closing file: %v", fname, err)
		return
	}

	r, err := lcio.Open(fname)
	if err != nil {
		t.Errorf("%s: error opening file: %v", fname, err)
	}
	defer r.Close()

	for i := 0; i < nevents; i++ {
		if !r.Next() {
			t.Errorf("%s: error reading event %d", fname, i)
			return
		}
		evt := r.Event()
		if got, want := evt.RunNumber, int32(4711); got != want {
			t.Errorf("%s: run-number error. got=%d. want=%d", fname, got, want)
			return
		}

		if got, want := evt.EventNumber, int32(i); got != want {
			t.Errorf("%s: run-number error. got=%d. want=%d", fname, got, want)
			return
		}

		if !evt.Has(CALHITS) {
			t.Errorf("%s: no %s collection", fname, CALHITS)
			return
		}

		if !evt.Has(CALHITSERR) {
			t.Errorf("%s: no %s collection", fname, CALHITSERR)
			return
		}

		calhits := evt.Get(CALHITS).(*lcio.CalorimeterHitContainer)
		calhitsErr := evt.Get(CALHITSERR).(*lcio.CalorimeterHitContainer)

		for j := 0; j < nhits; j++ {
			got := calhits.Hits[j]
			want := lcio.CalorimeterHit{
				CellID0: int32(i*100000 + j),
				Energy:  float32(i*j) + 117,
				Pos:     [3]float32{float32(i), float32(j), float32(i * j)},
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: event %d hit[%d]\ngot= %#v\nwant=%#v\n",
					fname, i, j,
					got, want,
				)
				return
			}

			want.EnergyErr = float32(i*j) * 0.117
			got = calhitsErr.Hits[j]
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: event %d hit[%d]\ngot= %#v\nwant=%#v\n",
					fname, i, j,
					got, want,
				)
				return
			}
		}
	}

	err = r.Close()
	if err != nil {
		t.Errorf("%s: error closing file: %v", fname, err)
		return
	}
}
