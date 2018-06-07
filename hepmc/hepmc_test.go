// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepmc_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"testing"

	"go-hep.org/x/hep/hepmc"
)

func TestEventRW(t *testing.T) {

	var mu sync.Mutex
	for _, table := range []struct {
		fname    string
		outfname string
		nevts    int
	}{
		{"testdata/small.hepmc", "out.small.hepmc", 1},
		{"testdata/test.hepmc", "out.hepmc", 6},
		{"out.hepmc", "rb.out.hepmc", 6},
	} {
		t.Run(table.fname, func(t *testing.T) {
			mu.Lock()
			testEventRW(t, table.fname, table.outfname, table.nevts)
			mu.Unlock()
		})
	}

	// clean-up
	for _, fname := range []string{"out.small.hepmc", "out.hepmc", "rb.out.hepmc"} {
		err := os.Remove(fname)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testEventRW(t *testing.T, fname, outfname string, nevts int) {
	f, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := hepmc.NewDecoder(f)
	if dec == nil {
		t.Fatalf("hepmc.decoder: nil decoder")
	}

	const NEVTS = 10
	evts := make([]*hepmc.Event, 0, NEVTS)
	for ievt := 0; ievt < NEVTS; ievt++ {
		var evt hepmc.Event
		err = dec.Decode(&evt)
		if err != nil {
			if err == io.EOF && ievt == nevts {
				break
			}
			t.Fatalf("file: %s. ievt=%d err=%v\n", fname, ievt, err)
		}
		evts = append(evts, &evt)
		defer hepmc.Delete(&evt)
	}

	o, err := os.Create(outfname)
	if err != nil {
		t.Fatal(err)
	}
	defer o.Close()

	enc := hepmc.NewEncoder(o)
	if enc == nil {
		t.Fatal(fmt.Errorf("hepmc.encoder: nil encoder"))
	}
	for _, evt := range evts {
		err = enc.Encode(evt)
		if err != nil {
			t.Fatalf("file: %s. err=%v\n", fname, err)
		}
	}
	err = enc.Close()
	if err != nil {
		t.Fatalf("file: %s. err=%v\n", fname, err)
	}

	err = o.Close()
	if err != nil {
		t.Fatalf("error closing output file %q: %v", outfname, err)
	}

	// test output files
	cmd := exec.Command("diff", "-urN", fname, outfname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("file: %s. err=%v\n", fname, err)
	}
}

type reader struct {
	header []byte
	footer []byte
	event  []byte
	data   []byte
	pos    int
	nevts  int
}

func newReader() *reader {
	r := &reader{
		header: []byte(`
HepMC::Version 2.06.09
HepMC::IO_GenEvent-START_EVENT_LISTING
`),
		data: []byte(`E 1 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 20 -3 4 0 0 0 0
U GEV MM
F 0 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0 0
V -1 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 1 1 0
P 1 2212 0.0000000000000000e+00 0.0000000000000000e+00 7.0000000000000000e+03 7.0000000000000000e+03 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -1 0
P 3 1 7.5000000000000000e-01 -1.5690000000000000e+00 3.2191000000000003e+01 3.2238000000000000e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -3 0
V -2 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 1 1 0
P 2 2212 0.0000000000000000e+00 0.0000000000000000e+00 -7.0000000000000000e+03 7.0000000000000000e+03 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -2 0
P 4 -2 -3.0470000000000002e+00 -1.9000000000000000e+01 -5.4628999999999998e+01 5.7920000000000002e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -3 0
V -3 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0 2 0
P 5 22 -3.8130000000000002e+00 1.1300000000000000e-01 -1.8330000000000000e+00 4.2329999999999997e+00 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
P 6 -24 1.5169999999999999e+00 -2.0680000000000000e+01 -2.0605000000000000e+01 8.5924999999999997e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -4 0
V -4 0 1.2000000000000000e-01 -2.9999999999999999e-01 5.0000000000000003e-02 4.0000000000000001e-03 0 2 0
P 7 1 -2.4449999999999998e+00 2.8815999999999999e+01 6.0819999999999999e+00 2.9552000000000000e+01 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
P 8 -2 3.9620000000000002e+00 -4.9497999999999998e+01 -2.6687000000000001e+01 5.6372999999999998e+01 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
`),
		footer: []byte(`HepMC::IO_GenEvent-END_EVENT_LISTING
`),
	}

	return r
}

func (r *reader) Read(data []byte) (int, error) {
	return r.read(data)
}

func printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}

func (r *reader) read(dst []byte) (int, error) {
	n := 0
	for n < len(dst) {
		// printf("::read - n=%d dst=%d...\n", n, len(dst))
		if len(r.header) > 0 {
			// printf(":: header...\n")
			sz := copy(dst, r.header)
			n += sz
			r.header = nil
		}
		end := len(dst[n:])
		if end > len(r.data[r.pos:]) {
			end = len(r.data[r.pos:]) + r.pos
		}
		// printf(":: copy(dst[%d:], data[%d:%d])...\n", n, r.pos, end)
		sz := copy(dst[n:], r.data[r.pos:end])
		// printf(":: copy(dst[%d:], data[%d:%d]) -> %d\n", n, r.pos, end, sz)
		n += sz
		if r.pos+sz >= len(r.data) {
			r.pos = 0
		} else {
			r.pos += sz
		}
	}
	return n, nil
}

func TestRead(t *testing.T) {
	r := newReader()
	dec := hepmc.NewDecoder(r)
	if dec == nil {
		t.Fatal(fmt.Errorf("hepmc.decoder: nil decoder"))
	}

	const NEVTS = 10
	const fname = "small.hepmc"
	evts := make([]*hepmc.Event, 0, NEVTS)
	for ievt := 0; ievt < NEVTS; ievt++ {
		var evt hepmc.Event
		err := dec.Decode(&evt)
		if err != nil {
			if err == io.EOF && ievt == NEVTS {
				break
			}
			t.Fatalf("file: %s. ievt=%d err=%v\n", fname, ievt, err)
		}
		evts = append(evts, &evt)
		defer hepmc.Delete(&evt)
	}
}

func BenchmarkDecode(b *testing.B) {
	r := newReader()
	dec := hepmc.NewDecoder(r)
	if dec == nil {
		b.Fatalf("hepmc.decoder: nil decoder")
	}

	const fname = "small.hepmc"

	{
		var evt hepmc.Event
		err := dec.Decode(&evt)
		if err != nil {
			b.Fatalf("error: %v\n", err)
		}
		hepmc.Delete(&evt)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var evt hepmc.Event
		err := dec.Decode(&evt)
		if err != nil {
			if err == io.EOF {
				break
			}
			b.Fatalf("file: %s. ievt=%d err=%v\n", fname, i, err)
		}
		//defer hepmc.Delete(&evt)
	}

}
