// Copyright 2017 The go-hep Authors. All rights reserved.
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
package main // import "go-hep.org/x/hep/groot/cmd/root-print"

import (
	"flag"
	"fmt"
	"log"
	"os"
	stdpath "path"
	"path/filepath"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	colors = plotutil.SoftColors
)

func main() {
	log.SetPrefix("root-print: ")
	log.SetFlags(0)

	var (
		odirFlag    = flag.String("o", "", "output directory for plots")
		fmtFlag     = flag.String("f", "pdf", "output format for plots (pdf, png, svg, ...)")
		verboseFlag = flag.Bool("v", false, "enable verbose mode")
	)

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

	err := rootprint(*odirFlag, flag.Args(), *fmtFlag, *verboseFlag)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func rootprint(odir string, fnames []string, otype string, verbose bool) error {
	err := os.MkdirAll(odir, 0755)
	if err != nil {
		return fmt.Errorf("could not create output directory %q: %w", odir, err)
	}

	for _, fname := range fnames {
		err := process(odir, fname, otype, verbose)
		if err != nil {
			return fmt.Errorf("could not process %q: %w", fname, err)
		}
	}

	return nil
}

func process(odir, name, otyp string, verbose bool) error {
	fname, hname, err := splitArg(name)
	if err != nil {
		return fmt.Errorf(
			"invalid input file format. got %q. want: \"file.root:histo\"",
			name,
		)
	}

	f, err := groot.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	re, err := regexp.CompilePOSIX(hname)
	if err != nil {
		return err
	}

	var objs []root.Object
	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		if err != nil {
			return err
		}
		name := path[len(f.Name()):]
		if name == "" {
			return nil
		}
		if !re.MatchString(name) {
			return nil
		}

		if !filter(obj) {
			return nil
		}

		objs = append(objs, obj)
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not inspect input ROOT file: %w", err)
	}

	for _, obj := range objs {
		err := printObject(odir, otyp, obj, verbose)
		if err != nil {
			return err
		}
	}
	return err
}

func printObject(odir, otyp string, obj root.Object, verbose bool) error {
	p := hplot.New()
	name := obj.(root.Named).Name()
	title := obj.(root.Named).Title()
	if title == "" {
		title = name
	}
	p.Title.Text = title

	oname := stdpath.Join(odir, name+"."+otyp)
	if verbose {
		log.Printf("printing %q to %s...", name, oname)
	}

	switch o := obj.(type) {
	case rhist.H2:
		h, err := rootcnv.H2D(o)
		if err != nil {
			return fmt.Errorf("could not convert %q to hbook.H2D: %w", name, err)
		}
		p.Add(hplot.NewH2D(h, nil))

	case rhist.H1:
		h, err := rootcnv.H1D(o)
		if err != nil {
			return fmt.Errorf("could not convert %q to hbook.H1D: %w", name, err)
		}
		hh := hplot.NewH1D(h)
		hh.Color = colors[2]
		hh.LineStyle.Color = colors[2]
		hh.LineStyle.Width = vg.Points(1.5)
		hh.Infos.Style = hplot.HInfoSummary

		p.Add(hh)

	case rhist.GraphErrors:
		h, err := rootcnv.S2D(o)
		if err != nil {
			return fmt.Errorf("could not convert %q to hbook.S2D: %w", name, err)
		}
		if name := h.Name(); name != "" {
			p.Title.Text = name
		}
		g := hplot.NewS2D(h, hplot.WithXErrBars(true), hplot.WithYErrBars(true))
		g.Color = colors[0]
		p.Add(g)

	case rhist.Graph:
		h, err := rootcnv.S2D(o)
		if err != nil {
			return fmt.Errorf("could not convert %q to hbook.S2D: %w", name, err)
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

	ext := strings.ToLower(stdpath.Ext(oname))
	if len(ext) > 0 {
		ext = ext[1:]
	}

	err := p.Save(20*vg.Centimeter, -1, oname)
	if err != nil {
		return fmt.Errorf("could not print %q to %q: %w", name, oname, err)
	}

	return nil
}

func filter(obj root.Object) bool {
	switch obj.(type) {
	case rhist.H1:
		return true

	case rhist.H2:
		return true

	case rhist.Graph, rhist.GraphErrors:
		return true
	}
	return false
}

func splitArg(cmd string) (fname, sel string, err error) {
	fname = cmd
	prefix := ""
	for _, p := range []string{"https://", "http://", "root://", "file://"} {
		if strings.HasPrefix(cmd, p) {
			prefix = p
			break
		}
	}
	fname = fname[len(prefix):]

	vol := filepath.VolumeName(fname)
	if vol != fname {
		fname = fname[len(vol):]
	}

	if strings.Count(fname, ":") > 1 {
		return "", "", fmt.Errorf("root-cp: too many ':' in %q", cmd)
	}

	i := strings.LastIndex(fname, ":")
	switch {
	case i > 0:
		sel = fname[i+1:]
		fname = fname[:i]
	default:
		sel = ".*"
	}
	if sel == "" {
		sel = ".*"
	}
	fname = prefix + vol + fname
	switch {
	case strings.HasPrefix(sel, "/"):
	case strings.HasPrefix(sel, "^/"):
	case strings.HasPrefix(sel, "^"):
		sel = "^/" + sel[1:]
	default:
		sel = "/" + sel
	}
	return fname, sel, err
}
