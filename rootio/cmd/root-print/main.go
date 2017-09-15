// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-print prints ROOT files contents to PDF, PNG, ... files.
//
// Examples:
//
//  $> root-print -f pdf ./testdata/histos.root
//  $> root-print -f pdf ./testdata/histos.root:h1
//  $> root-print -f pdf ./testdata/histos.root:h.*
//  $> root-print -f pdf -o output ./testdata/histos.root:h1
//
//  $> root-print -h
//  Usage: root-print [options] file.root [file.root [...]]
//
//  options:
//    -f string
//      	output format for plots (pdf, png, svg, ...) (default "pdf")
//    -o string
//      	output directory for plots
//    -v	enable verbose mode
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/rootio"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	odirFlag    = flag.String("o", "", "output directory for plots")
	fmtFlag     = flag.String("f", "pdf", "output format for plots (pdf, png, svg, ...)")
	verboseFlag = flag.Bool("v", false, "enable verbose mode")

	colors = plotutil.SoftColors
)

func main() {
	log.SetPrefix("root-print: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-print [options] file.root [file.root [...]]
ex:
 $> root-print -f pdf ./testdata/histos.root
 $> root-print -f pdf ./testdata/histos.root:h1
 $> root-print -f pdf ./testdata/histos.root:h.*
 $> root-print -f pdf -o output ./testdata/histos.root:h1

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		log.Fatalf("need at least 1 input ROOT file")
	}

	_ = os.MkdirAll(*odirFlag, 0755)

	for _, fname := range flag.Args() {
		err := process(fname)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func process(name string) error {
	var (
		fname = ""
		hname = ""
	)

	toks := strings.Split(name, ":")
	switch len(toks) {
	case 2:
		fname = toks[0]
		hname = toks[1]
	case 1:
		fname = toks[0]
		hname = ".*"
	case 0:
		fname = name
		hname = ".*"
	default:
		return fmt.Errorf(
			"invalid input file format. got %q. want: \"file.root:histo\"",
			name,
		)
	}

	f, err := rootio.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	re, err := regexp.CompilePOSIX(hname)
	if err != nil {
		return err
	}

	var objs []rootio.Object
	for _, k := range f.Keys() {
		if !re.MatchString(k.Name()) {
			continue
		}
		o, err := k.Object()
		if err != nil {
			return err
		}
		if !filter(o) {
			continue
		}
		objs = append(objs, o)
	}

	for _, obj := range objs {
		err := printObject(f, obj)
		if err != nil {
			return err
		}
	}
	return err
}

func printObject(f *rootio.File, obj rootio.Object) error {
	p, err := hplot.New()
	if err != nil {
		return err
	}

	name := obj.(rootio.Named).Name()
	title := obj.(rootio.Named).Title()
	if title == "" {
		title = name
	}
	p.Title.Text = title

	oname := filepath.Join(*odirFlag, name+"."+*fmtFlag)
	if *verboseFlag {
		log.Printf("printing %q to %s...", name, oname)
	}

	switch o := obj.(type) {
	case *rootio.H1D, *rootio.H1F, *rootio.H1I:
		h, err := rootcnv.H1D(o.(yodacnv.Marshaler))
		if err != nil {
			return err
		}
		hh, err := hplot.NewH1D(h)
		if err != nil {
			return err
		}
		hh.Color = colors[2]
		hh.LineStyle.Color = colors[2]
		hh.LineStyle.Width = vg.Points(1.5)
		hh.Infos.Style = hplot.HInfoSummary

		p.Add(hh)

	case *rootio.H2D, *rootio.H2F, *rootio.H2I:
		h, err := rootcnv.H2D(o.(yodacnv.Marshaler))
		if err != nil {
			return err
		}
		p.Add(hplot.NewH2D(h, nil))

	case rootio.GraphErrors:
		h, err := rootcnv.S2D(o)
		if err != nil {
			return err
		}
		if name := h.Name(); name != "" {
			p.Title.Text = name
		}
		g := hplot.NewS2D(h, hplot.WithXErrBars, hplot.WithYErrBars)
		g.Color = colors[0]
		p.Add(g)

	case rootio.Graph:
		h, err := rootcnv.S2D(o)
		if err != nil {
			return err
		}
		if name := h.Name(); name != "" {
			p.Title.Text = name
		}
		g := hplot.NewS2D(h)
		g.Color = colors[0]
		p.Add(g)

	default:
		return fmt.Errorf("unknown type %T for %q", o, name)
	}

	p.Add(hplot.NewGrid())

	ext := strings.ToLower(filepath.Ext(oname))
	if len(ext) > 0 {
		ext = ext[1:]
	}

	err = p.Save(20*vg.Centimeter, -1, oname)
	if err != nil {
		return err
	}

	return err
}

func filter(obj rootio.Object) bool {
	switch obj.(type) {
	case *rootio.H1D, *rootio.H1F, *rootio.H1I:
		return true

	case *rootio.H2D, *rootio.H2F, *rootio.H2I:
		return true

	case rootio.Graph, rootio.GraphErrors:
		return true
	}
	return false
}
