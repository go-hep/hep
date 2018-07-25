// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-hepmc-dump is a simple command to dump in an almost human-friendly format
// the content of a hepmc file.
// ex:
//  $ go-hepmc-dump foo.hepmc | head -n20
// ________________________________________________________________________________
// GenEvent: #0000 ID=  111 SignalProcessGenVertex Barcode: 0
//  Momentum units:      GEV     Position units:       MM
//  Cross Section: 2.666668e+12 +/- 2.666668e+12
//  Entries this event: 129 vertices, 241 particles.
//  Beam Particle barcodes: 1 2
//  RndmState(0)=
//  Wgts(1)=(0,1.000000)
//  EventScale 1.438780e+00 [energy] 	 alphaQCD=4.519907e-01	 alphaQED=7.472465e-03
//                                     GenParticle Legend
//         Barcode   PDG ID      ( Px,       Py,       Pz,     E ) Stat  DecayVtx
// ________________________________________________________________________________
// GenVertex:       -1 ID:    0 (X,cT):0
//  I: 1         7        21 -5.87e-14, 9.60e-15, 3.41e+01, 3.41e+01 42        -1
//  O: 1         3        21  0.00e+00, 0.00e+00, 3.41e+01, 3.41e+01 21        -3
// GenVertex:       -2 ID:    0 (X,cT):0
//  I: 1         8        21  3.20e-14, 0.00e+00,-7.23e+01, 7.23e+01 41        -2
//  O: 2         4        21  0.00e+00, 0.00e+00,-1.19e+00, 1.19e+00 21        -3
//              11        21  9.42e-01,-1.56e-01,-7.11e+01, 7.11e+01 43       -12
//
package main

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/hepmc"
)

func main() {
	log.SetPrefix("hepmc-dump: ")
	log.SetFlags(0)

	var (
		err error
		r   io.ReadCloser
		w   io.Writer = os.Stdout
	)

	switch len(os.Args) {
	case 1:
		r = os.Stdin
	case 2:
		r, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
	default:
	}

	err = dump(w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func dump(w io.Writer, r io.Reader) error {
	var err error
	dec := hepmc.NewDecoder(r)
	for {
		var evt hepmc.Event
		err = dec.Decode(&evt)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.Errorf("error decoding event: %v", err)
		}
		err = evt.Print(w)
		if err != nil {
			return errors.Errorf("error printing event: %v", err)
		}

		err = hepmc.Delete(&evt)
		if err != nil {
			return errors.Errorf("error deleting event: %v", err)
		}
	}
	panic("unreachable")
}
