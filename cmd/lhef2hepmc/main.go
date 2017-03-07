// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lhef2hepmc converts a LHEF input file into a HepMC file.
//
// Example:
//
//  $> lhef2hepmc -i in.lhef -o out.hepmc
//  $> lhef2hepmc < in.lhef > out.hepmc
package main // import "go-hep.org/x/hep/cmd/lhef2hepmc"

import (
	"flag"
	"io"
	"log"
	"math"
	"os"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/lhef"
)

var (
	ifname = flag.String("i", "", "path to LHEF input file (default: STDIN)")
	ofname = flag.String("o", "", "path to HEPMC output file (default: STDOUT)")

	// in case IDWTUP == +/-4, one has to keep track of the accumulated
	// weights and event numbers to evaluate the cross section on-the-fly.
	// The last evaluation is the one used.
	// Better to be sure that crossSection() is never used to fill the
	// histograms, but only in the finalization stage, by reweighting the
	// histograms with crossSection()/sumOfWeights()
	sumw  = 0.0
	sumw2 = 0.0
	nevt  = 0
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("lhef2hepmc: ")
	log.SetOutput(os.Stderr)

	flag.Parse()

	var r io.Reader
	if *ifname == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(*ifname)
		if err != nil {
			log.Fatal(err)
		}
		r = f
		defer f.Close()
	}

	var w io.Writer
	if *ofname == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(*ofname)
		if err != nil {
			log.Fatal(err)
		}
		w = f
		defer f.Close()
	}

	dec, err := lhef.NewDecoder(r)
	if err != nil {
		log.Fatalf("error creating LHEF decoder: %v", err)
	}

	enc := hepmc.NewEncoder(w)
	if enc == nil {
		log.Fatalf("error creating HepMC encoder: %v", err)
	}
	defer enc.Close()

	for ievt := 0; ; ievt++ {
		lhevt, err := dec.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error at event #%d: %v", ievt, err)
		}

		evt := hepmc.Event{
			EventNumber: ievt + 1,
			Particles:   make(map[int]*hepmc.Particle),
			Vertices:    make(map[int]*hepmc.Vertex),
			Weights:     hepmc.NewWeights(),
		}

		// define the units
		evt.MomentumUnit = hepmc.GEV
		evt.LengthUnit = hepmc.MM

		weight := lhevt.XWGTUP
		evt.Weights.Add("0", weight)

		xsecval := -1.0
		xsecerr := -1.0
		switch math.Abs(float64(dec.Run.IDWTUP)) {
		case 3:
			xsecval = dec.Run.XSECUP[0]
			xsecerr = dec.Run.XSECUP[1]

		case 4:
			sumw += weight
			sumw2 += weight * weight
			nevt++
			xsecval = sumw / float64(nevt)
			xsecerr2 := (sumw2/float64(nevt) - xsecval*xsecval) / float64(nevt)

			if xsecerr2 < 0 {
				log.Printf("WARNING: xsecerr^2 < 0. forcing to zero. (%f)\n", xsecerr2)
				xsecerr2 = 0.
			}
			xsecerr = math.Sqrt(xsecerr2)

		default:
			log.Fatalf("IDWTUP=%v value not handled yet", dec.Run.IDWTUP)
		}

		evt.CrossSection = &hepmc.CrossSection{Value: xsecval, Error: xsecerr}
		vtx := hepmc.Vertex{
			Event:   &evt,
			Barcode: -1,
		}
		p1 := hepmc.Particle{
			Momentum: fmom.PxPyPzE{
				0, 0,
				dec.Run.EBMUP[0],
				dec.Run.EBMUP[0],
			},
			PdgID:   dec.Run.IDBMUP[0],
			Status:  4,
			Barcode: 1,
		}
		p2 := hepmc.Particle{
			Momentum: fmom.PxPyPzE{
				0, 0,
				dec.Run.EBMUP[1],
				dec.Run.EBMUP[1],
			},
			PdgID:   dec.Run.IDBMUP[1],
			Status:  4,
			Barcode: 2,
		}
		err = vtx.AddParticleIn(&p1)
		if err != nil {
			log.Fatalf("error at event #%d: %v", ievt, err)
		}
		err = vtx.AddParticleIn(&p2)
		if err != nil {
			log.Fatalf("error at event #%d: %v\n", ievt, err)
		}
		evt.Beams[0] = &p1
		evt.Beams[1] = &p2

		nmax := 2
		imax := int(lhevt.NUP)
		for i := 0; i < imax; i++ {
			if lhevt.ISTUP[i] != 1 {
				continue
			}
			nmax += 1
			vtx.AddParticleOut(&hepmc.Particle{
				Momentum: fmom.PxPyPzE{
					lhevt.PUP[i][0],
					lhevt.PUP[i][1],
					lhevt.PUP[i][2],
					lhevt.PUP[i][3],
				},
				GeneratedMass: lhevt.PUP[i][4],
				PdgID:         lhevt.IDUP[i],
				Status:        1,
				Barcode:       3 + i,
			})
		}
		err = evt.AddVertex(&vtx)
		if err != nil {
			log.Fatalf("error at event #%d: %v\n", ievt, err)
		}

		nparts := len(evt.Particles)
		if nmax != nparts {
			log.Printf("error at event #%d: LHEF/HEPMC inconsistency (LHEF particles: %d, HEPMC particles: %d)\n", ievt, nmax, nparts)
			for _, p := range evt.Particles {
				log.Printf("part: %v\n", p)
			}
			log.Printf("p1: %v\n", p1)
			log.Printf("p2: %v\n", p2)
			os.Exit(1)
		}
		if len(evt.Vertices) != 1 {
			log.Fatalf("error at event #%d: inconsistent number of vertices in HEPMC (got %d, expected 1)\n", ievt, len(evt.Vertices))
		}

		err = enc.Encode(&evt)
		if err != nil {
			log.Fatalf("error at event #%d: %v\n", ievt, err)
		}
	}
}
