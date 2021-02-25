// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"compress/zlib"
	"flag"
	"log"
	"os"

	"go-hep.org/x/hep/rio"
)

type RunHeader struct {
	RunNbr   int32
	Detector string
	Descr    string
	SubDets  []string
	Ints     []int64
	Floats   []float64
}

func main() {
	fname := flag.String("fname", "runhdr.rio", "file to create")
	compr := flag.Bool("compr", false, "enable records compression")
	flag.Parse()

	f, err := os.Create(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	w, err := rio.NewWriter(f)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	if *compr {
		err = w.SetCompressor(rio.CompressZlib, zlib.DefaultCompression)
		if err != nil {
			log.Fatalf("error setting compressor: %v\n", err)
		}
	}

	var runhdr RunHeader
	rec := w.Record("RioRunHeader")
	if rec == nil {
		log.Fatal("could not fetch record [RioRunHeader]\n")
	}

	err = rec.Connect("RunHeader", &runhdr)
	if err != nil {
		log.Fatalf("error connecting [RunHeader]: %v", err)
	}
	blk := rec.Block("RunHeader")

	for irec := 0; irec < 10; irec++ {
		runhdr = RunHeader{
			RunNbr:   int32(irec),
			Detector: "MyDetector",
			Descr:    "dummy run number",
			SubDets:  []string{"subdet 0", "subdet 1"},
			Floats: []float64{
				float64(irec) + 100,
				float64(irec) + 200,
				float64(irec) + 300,
			},
			Ints: []int64{
				int64(irec) + 100,
				int64(irec) + 200,
				int64(irec) + 300,
			},
		}
		err = blk.Write(&runhdr)
		if err != nil {
			log.Fatalf("error writing block: %v (irec=%d)", err, irec)
		}
		err = rec.Write()
		if err != nil {
			log.Fatalf("error writing record: %v (irec=%d)", err, irec)
		}
	}
}

// EOF
