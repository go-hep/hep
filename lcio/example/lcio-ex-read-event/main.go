// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lcio-ex-read-event is the hep/x/lcio example equivalent to:
//  https://github.com/iLCSoft/LCIO/blob/master/examples/cpp/rootDict/readEventTree.C
//
// example:
//
//  $> lcio-ex-read-event ./DST01-06_ppr004_bbcsdu.slcio
//  lcio-ex-read-event: read 50 events from file "./DST01-06_ppr004_bbcsdu.slcio"
//  $> open out.png
//
package main

import (
	"flag"
	"image/color"
	"io"
	"log"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/lcio"
	"gonum.org/v1/plot/vg"
)

func main() {
	log.SetPrefix("lcio-ex-read-event: ")
	log.SetFlags(0)

	var (
		fname  = ""
		h      = hbook.NewH1D(100, 0., 100.)
		nevts  = 0
		mcname = flag.String("mc", "MCParticlesSkimmed", "name of the MCParticle collection to read")
	)

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := lcio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for f.Next() {
		evt := f.Event()
		mcs := evt.Get(*mcname).(*lcio.McParticleContainer)
		for _, mc := range mcs.Particles {
			h.Fill(mc.Energy(), 1)
		}
		nevts++
	}

	err = f.Err()
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("read %d events from file %q", nevts, fname)

	p := hplot.New()
	p.Title.Text = "LCIO -- McParticles"
	p.X.Label.Text = "E (GeV)"

	hh := hplot.NewH1D(h)
	hh.Color = color.RGBA{R: 255, A: 255}
	hh.Infos.Style = hplot.HInfoSummary

	p.Add(hh)
	p.Add(hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, -1, "out.png")
	if err != nil {
		log.Fatal(err)
	}
}
