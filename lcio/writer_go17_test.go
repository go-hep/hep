// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package lcio_test

import (
	"compress/flate"
	"encoding/hex"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func TestCreateCompressedRunHeader(t *testing.T) {
	const fname = "testdata/run-header-compressed.slcio"
	w, err := lcio.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	w.SetCompressionLevel(flate.BestCompression)

	rhdr := lcio.RunHeader{
		RunNumber:    42,
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
		t.Errorf("%s: --- ref ---\n%s\n", fname, hex.Dump(ref))
		t.Errorf("%s: --- chk ---\n%s\n", fname, hex.Dump(chk))
		t.Fatalf("%s: differ with golden", fname)
	}

	os.Remove(fname)
}
