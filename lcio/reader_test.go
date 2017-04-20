// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio_test

import (
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func TestOpen(t *testing.T) {
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

	for _, fname := range []string{
		"testdata/run-header_golden.slcio",
		"testdata/run-header-compressed_golden.slcio",
	} {
		r, err := lcio.Open(fname)
		if err != nil {
			t.Fatalf("%s: error opening file: %v", fname, err)
		}
		defer r.Close()

		r.Next()
		if err := r.Err(); err != nil && err != io.EOF {
			t.Fatalf("%s: %v", fname, err)
		}

		if got, want := r.RunHeader(), rhdr; !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: run-headers differ.\ngot= %#v\nwant=%#v\n", fname, got, want)
		}

		err = r.Close()
		if err != nil {
			t.Fatalf("%s: error closing file: %v", fname, err)
		}
	}
}
