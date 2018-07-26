// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"

	"go-hep.org/x/hep/sio"
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
	fname := flag.String("fname", "runhdr.sio", "file to create")
	compr := flag.Bool("compr", false, "enable records compression")
	flag.Parse()

	f, err := sio.Create(*fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("*** error closing file [%s]: %v\n", *fname, err)
			os.Exit(1)
		}
	}()

	var runhdr RunHeader
	rec := f.Record("RioRunHeader")
	if rec == nil {
		fmt.Printf("*** error: could not fetch record [RioRunHeader]\n")
		os.Exit(1)
	}

	rec.SetCompress(*compr)
	err = rec.Connect("RunHeader", &runhdr)
	if err != nil {
		fmt.Printf("error connecting [RunHeader]: %v", err)
		os.Exit(1)
	}

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
		err = f.WriteRecord(rec)
		if err != nil {
			fmt.Printf("error writing record: %v (irec=%d)", err, irec)
			os.Exit(1)
		}

		err = f.Sync()
		if err != nil {
			fmt.Printf("error flushing record: %v (irec=%d)", err, irec)
			os.Exit(1)
		}
	}
}
