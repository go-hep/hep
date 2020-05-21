// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepmc_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/hepmc"
)

func TestEventRW(t *testing.T) {
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
			testEventRW(t, table.fname, table.outfname, table.nevts)
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
		defer func() {
			err = hepmc.Delete(&evt)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
		}()
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
	data   []byte
	pos    int
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
	for ievt := 0; ievt < NEVTS; ievt++ {
		var evt hepmc.Event
		err := dec.Decode(&evt)
		if err != nil {
			if err == io.EOF && ievt == NEVTS {
				break
			}
			t.Fatalf("file: %s. ievt=%d err=%v\n", fname, ievt, err)
		}
		err = hepmc.Delete(&evt)
		if err != nil {
			t.Fatalf("error: %+v", err)
		}
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
		_ = hepmc.Delete(&evt)
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
		_ = hepmc.Delete(&evt)
	}

}

// This example will place the following event into HepMC "by hand"
//
//     name status pdg_id  parent Px       Py    Pz       Energy      Mass
//  1  !p+!    3   2212    0,0    0.000    0.000 7000.000 7000.000    0.938
//  2  !p+!    3   2212    0,0    0.000    0.000-7000.000 7000.000    0.938
//=========================================================================
//  3  !d!     3      1    1,1    0.750   -1.569   32.191   32.238    0.000
//  4  !u~!    3     -2    2,2   -3.047  -19.000  -54.629   57.920    0.000
//  5  !W-!    3    -24    1,2    1.517   -20.68  -20.605   85.925   80.799
//  6  !gamma! 1     22    1,2   -3.813    0.113   -1.833    4.233    0.000
//  7  !d!     1      1    5,5   -2.445   28.816    6.082   29.552    0.010
//  8  !u~!    1     -2    5,5    3.962  -49.498  -26.687   56.373    0.006
//
// now we build the graph, which will look like
//  #                       p7                         #
//  # p1                   /                           #
//  #   \v1__p3      p5---v4                           #
//  #         \_v3_/       \                           #
//  #         /    \        p8                         #
//  #    v2__p4     \                                  #
//  #   /            p6                                #
//  # p2                                               #
//  #                                                  #
func ExampleEvent_buildFromScratch() {
	var err error

	// first create the event container, with signal process 20, event number 1
	evt := hepmc.Event{
		SignalProcessID: 20,
		EventNumber:     1,
		Particles:       make(map[int]*hepmc.Particle),
		Vertices:        make(map[int]*hepmc.Vertex),
	}
	defer func() {
		err = evt.Delete()
		if err != nil {
			log.Fatalf("could not clean-up event: %+v", err)
		}
	}()

	// define the units
	evt.MomentumUnit = hepmc.GEV
	evt.LengthUnit = hepmc.MM

	// create vertex 1 and 2, together with their in-particles
	v1 := &hepmc.Vertex{}
	err = evt.AddVertex(v1)
	if err != nil {
		log.Fatal(err)
	}

	err = v1.AddParticleIn(&hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(0, 0, 7000, 7000),
		PdgID:    2212,
		Status:   3,
	})
	if err != nil {
		log.Fatal(err)
	}

	v2 := &hepmc.Vertex{}
	err = evt.AddVertex(v2)
	if err != nil {
		log.Fatal(err)
	}

	err = v2.AddParticleIn(&hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(0, 0, -7000, 7000),
		PdgID:    2212,
		Status:   3,
		//Barcode:  2,
	})
	if err != nil {
		log.Fatal(err)
	}

	// create the outgoing particles of v1 and v2
	p3 := &hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(.750, -1.569, 32.191, 32.238),
		PdgID:    1,
		Status:   3,
		// Barcode: 3,
	}
	err = v1.AddParticleOut(p3)
	if err != nil {
		log.Fatal(err)
	}

	p4 := &hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(-3.047, -19., -54.629, 57.920),
		PdgID:    -2,
		Status:   3,
		// Barcode: 4,
	}
	err = v2.AddParticleOut(p4)
	if err != nil {
		log.Fatal(err)
	}

	// create v3
	v3 := &hepmc.Vertex{}
	err = evt.AddVertex(v3)
	if err != nil {
		log.Fatal(err)
	}

	err = v3.AddParticleIn(p3)
	if err != nil {
		log.Fatal(err)
	}

	err = v3.AddParticleIn(p4)
	if err != nil {
		log.Fatal(err)
	}

	err = v3.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(-3.813, 0.113, -1.833, 4.233),
		PdgID:    22,
		Status:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	p5 := &hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(1.517, -20.68, -20.605, 85.925),
		PdgID:    -24,
		Status:   3,
	}
	err = v3.AddParticleOut(p5)
	if err != nil {
		log.Fatal(err)
	}

	// create v4
	v4 := &hepmc.Vertex{
		Position: fmom.NewPxPyPzE(0.12, -0.3, 0.05, 0.004),
	}
	err = evt.AddVertex(v4)
	if err != nil {
		log.Fatal(err)
	}

	err = v4.AddParticleIn(p5)
	if err != nil {
		log.Fatal(err)
	}

	err = v4.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(-2.445, 28.816, 6.082, 29.552),
		PdgID:    1,
		Status:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = v4.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.NewPxPyPzE(3.962, -49.498, -26.687, 56.373),
		PdgID:    -2,
		Status:   1,
	})
	if err != nil {
		log.Fatal(err)
	}

	evt.SignalVertex = v3

	err = evt.Print(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	out := new(bytes.Buffer)
	enc := hepmc.NewEncoder(out)
	err = enc.Encode(&evt)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out.Bytes())

	// Output:
	// ________________________________________________________________________________
	// GenEvent: #0001 ID=   20 SignalProcessGenVertex Barcode: -3
	//  Momentum units:     GEV     Position units:      MM
	//  Entries this event: 4 vertices, 8 particles.
	//  Beam Particles are not defined.
	//  RndmState(0)=
	//  Wgts(0)=
	//  EventScale 0.00000 [energy] 	 alphaQCD=0.00000000	 alphaQED=0.00000000
	//                                     GenParticle Legend
	//         Barcode   PDG ID      ( Px,       Py,       Pz,     E ) Stat  DecayVtx
	// ________________________________________________________________________________
	// GenVertex:       -1 ID:    0 (X,cT):0
	//  I: 1         1     2212 +0.00e+00,+0.00e+00,+7.00e+03,+7.00e+03   3        -1
	//  O: 1         3        1 +7.50e-01,-1.57e+00,+3.22e+01,+3.22e+01   3        -3
	// GenVertex:       -2 ID:    0 (X,cT):0
	//  I: 1         2     2212 +0.00e+00,+0.00e+00,-7.00e+03,+7.00e+03   3        -2
	//  O: 1         4       -2 -3.05e+00,-1.90e+01,-5.46e+01,+5.79e+01   3        -3
	// GenVertex:       -3 ID:    0 (X,cT):0
	//  I: 2         3        1 +7.50e-01,-1.57e+00,+3.22e+01,+3.22e+01   3        -3
	//               4       -2 -3.05e+00,-1.90e+01,-5.46e+01,+5.79e+01   3        -3
	//  O: 2         5       22 -3.81e+00,+1.13e-01,-1.83e+00,+4.23e+00   1
	//               6      -24 +1.52e+00,-2.07e+01,-2.06e+01,+8.59e+01   3        -4
	// Vertex:       -4 ID:    0 (X,cT)=+1.20e-01,-3.00e-01,+5.00e-02,+4.00e-03
	//  I: 1         6      -24 +1.52e+00,-2.07e+01,-2.06e+01,+8.59e+01   3        -4
	//  O: 2         7        1 -2.44e+00,+2.88e+01,+6.08e+00,+2.96e+01   1
	//               8       -2 +3.96e+00,-4.95e+01,-2.67e+01,+5.64e+01   1
	// ________________________________________________________________________________
	//
	// HepMC::Version 2.06.09
	// HepMC::IO_GenEvent-START_EVENT_LISTING
	// E 1 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 20 -3 4 0 0 0 0
	// U GEV MM
	// F 0 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0 0
	// V -1 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 1 1 0
	// P 1 2212 0.0000000000000000e+00 0.0000000000000000e+00 7.0000000000000000e+03 7.0000000000000000e+03 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -1 0
	// P 3 1 7.5000000000000000e-01 -1.5690000000000000e+00 3.2191000000000003e+01 3.2238000000000000e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -3 0
	// V -2 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 1 1 0
	// P 2 2212 0.0000000000000000e+00 0.0000000000000000e+00 -7.0000000000000000e+03 7.0000000000000000e+03 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -2 0
	// P 4 -2 -3.0470000000000002e+00 -1.9000000000000000e+01 -5.4628999999999998e+01 5.7920000000000002e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -3 0
	// V -3 0 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0.0000000000000000e+00 0 2 0
	// P 5 22 -3.8130000000000002e+00 1.1300000000000000e-01 -1.8330000000000000e+00 4.2329999999999997e+00 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
	// P 6 -24 1.5169999999999999e+00 -2.0680000000000000e+01 -2.0605000000000000e+01 8.5924999999999997e+01 0.0000000000000000e+00 3 0.0000000000000000e+00 0.0000000000000000e+00 -4 0
	// V -4 0 1.2000000000000000e-01 -2.9999999999999999e-01 5.0000000000000003e-02 4.0000000000000001e-03 0 2 0
	// P 7 1 -2.4449999999999998e+00 2.8815999999999999e+01 6.0819999999999999e+00 2.9552000000000000e+01 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
	// P 8 -2 3.9620000000000002e+00 -4.9497999999999998e+01 -2.6687000000000001e+01 5.6372999999999998e+01 0.0000000000000000e+00 1 0.0000000000000000e+00 0.0000000000000000e+00 0 0
	//
}
