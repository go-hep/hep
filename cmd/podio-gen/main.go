// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command podio-gen generates a complete EDM from a PODIO YAML file definition.
//
// Usage: podio-gen [OPTIONS] edm.yaml
//
// Example:
//
//   $> podio-gen -p myedm -o out.go -r 'edm4hep::->edm_,ExNamespace::->exns_' edm.yaml
//
// Options:
//   -o string
//     	path to the output file containing the generated code (default "out.go")
//   -p string
//     	package name for the PODIO generated types (default "podio")
//   -r string
//     	comma-separated list of rewrite rules (e.g., 'edm4hep::->edm_')
package main // import "go-hep.org/x/hep/cmd/podio-gen"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetPrefix("podio: ")
	log.SetFlags(0)

	var (
		pkg     = flag.String("p", "podio", "package name for the PODIO generated types")
		oname   = flag.String("o", "out.go", "path to the output file containing the generated code")
		rewrite = flag.String("r", "", "comma-separated list of rewrite rules (e.g., 'edm4hep::->edm_')")
	)

	flag.Usage = func() {
		fmt.Printf(`podio-gen generates a complete EDM from a PODIO YAML file definition.

Usage: podio-gen [OPTIONS] edm.yaml

Example:

  $> podio-gen -p myedm -o out.go -r 'edm4hep::->edm_,ExNamespace::->exns_' edm.yaml

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	switch flag.NArg() {
	case 0:
		flag.Usage()
		log.Fatalf("missing input PODIO file")
	case 1:
		// ok
	default:
		flag.Usage()
		log.Fatalf("too many input PODIO files")
	}

	out, err := os.Create(*oname)
	if err != nil {
		log.Fatalf("could not create output file %q: %+v", *oname, err)
	}
	defer out.Close()

	fname := flag.Arg(0)

	err = process(out, *pkg, fname, *rewrite)
	if err != nil {
		log.Fatalf("could not process %q: %+v", fname, err)
	}

	err = out.Close()
	if err != nil {
		log.Fatalf("could not save output file %q: %+v", *oname, err)
	}
}

func process(w io.Writer, pkg, fname, rules string) error {
	g, err := newGenerator(w, pkg, fname, rules)
	if err != nil {
		return fmt.Errorf("could not create PODIO generator: %w", err)
	}

	err = g.generate()
	if err != nil {
		return fmt.Errorf("could not generate PODIO code for %q: %w", fname, err)
	}

	return nil
}
