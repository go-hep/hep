// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio_test

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func ExampleReader() {
	r, err := lcio.Open("testdata/event_golden.slcio")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for r.Next() {
		evt := r.Event()
		fmt.Printf("event number = %d (weight=%+e)\n", evt.EventNumber, evt.Weight())
		fmt.Printf("run   number = %d\n", evt.RunNumber)
		fmt.Printf("detector     = %q\n", evt.Detector)
		fmt.Printf("collections  = %v\n", evt.Names())
		calohits := evt.Get("CaloHits").(*lcio.CalorimeterHitContainer)
		fmt.Printf("calohits: %d\n", len(calohits.Hits))
		for i, hit := range calohits.Hits {
			fmt.Printf(" calohit[%d]: cell-id0=%d cell-id1=%d ene=%+e ene-err=%+e\n",
				i, hit.CellID0, hit.CellID1, hit.Energy, hit.EnergyErr,
			)
		}
	}

	err = r.Err()
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// event number = 52 (weight=+4.200000e+01)
	// run   number = 42
	// detector     = "my detector"
	// collections  = [McParticles SimCaloHits CaloHits]
	// calohits: 1
	//  calohit[0]: cell-id0=1024 cell-id1=2048 ene=+1.000000e+03 ene-err=+1.000000e-01
}

func TestOpen(t *testing.T) {
	ref := lcio.RunHeader{
		RunNumber:    42,
		Descr:        "a simple run header",
		Detector:     "my detector",
		SubDetectors: []string{"det-1", "det-2"},
		Params: lcio.Params{
			Floats: map[string][]float32{
				"floats-1": {1, 2, 3},
				"floats-2": {4, 5, 6},
			},
			Ints: map[string][]int32{
				"ints-1": {1, 2, 3},
				"ints-2": {4, 5, 6},
			},
			Strings: map[string][]string{
				"strs-1": {"1", "2", "3"},
			},
		},
	}

	for _, fname := range []string{
		"testdata/run-header_golden.slcio",
		"testdata/run-header-compressed_golden.slcio",
	} {
		t.Run(fname, func(t *testing.T) {
			r, err := lcio.Open(fname)
			if err != nil {
				t.Fatalf("%s: error opening file: %v", fname, err)
			}
			defer r.Close()

			r.Next()
			if err := r.Err(); err != nil && err != io.EOF {
				t.Fatalf("%s: %v", fname, err)
			}

			rhdr := r.RunHeader()
			if got, want := rhdr, ref; !reflect.DeepEqual(got, want) {
				t.Fatalf("%s: run-headers differ.\ngot= %#v\nwant=%#v\n", fname, got, want)
			}

			if got, want := rhdr.String(), ref.String(); got != want {
				t.Fatalf("%s: run-headers differ.\ngot= %q\nwant=%q\n", fname, got, want)
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("%s: error closing file: %v", fname, err)
			}
		})
	}
}
