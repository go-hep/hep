// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// hplot is a simple gnuplot-like command to create plots
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var g_persist = flag.Bool("p", false, "lets plot windows survive after hplot exists")
var g_cmdlist = flag.String("e", "", "executes the requested commands before loading the next input file")

func main() {
	fmt.Printf(":: welcome to hplot\n")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "**error** you need to give an input data file\n")
		os.Exit(1)
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** opening file [%s]: %v\n", flag.Args()[0], err)
		os.Exit(1)
	}
	defer f.Close()

	p, err := hplot.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** creating plotter: %v\n", err)
		os.Exit(1)
	}

	values := make(plotter.Values, 0, 10)
	for err == nil {
		var v float64
		_, err = fmt.Fscanf(f, "%f\n", &v)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "** error ** %v\n", err)
			os.Exit(1)
		}

		values = append(values, v)
	}

	// create histogram from data
	h, err := hplot.NewH1FromValuer(values, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** creating histogram: %v\n", err)
		os.Exit(1)
	}

	p.Add(h)
	p.Add(hplot.NewGrid())

	// Save the plot to a PDF file.
	if err := p.Save(-1, 10*vg.Centimeter, "hist.pdf"); err != nil {
		fmt.Fprintf(os.Stderr, "**error** saving plot: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(":: bye.\n")
}
