// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lhef_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"go-hep.org/x/hep/lhef"
)

const r_debug = false
const ifname = "testdata/ttbar.lhe"

func TestLhefReading(t *testing.T) {
	f, err := os.Open(ifname)
	if err != nil {
		t.Error(err)
	}

	dec, err := lhef.NewDecoder(f)
	if err != nil {
		t.Error(err)
	}

	for i := 0; ; i++ {
		if r_debug {
			fmt.Printf("===[%d]===\n", i)
		}
		evt, err := dec.Decode()
		if err == io.EOF {
			if r_debug {
				fmt.Printf("** EOF **\n")
			}
			break
		}
		if err != nil {
			t.Error(err)
		}
		if r_debug {
			fmt.Printf("evt: %v\n", *evt)
		}
	}
}
