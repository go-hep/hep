// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run ./gen-data.go -f ../../testdata/streamers.root
//go:generate root-gen-type -p main -t Event,P3 -o ./testdata/streamers.txt ../../testdata/streamers.root

// Command root-gen-type generates a Go type from the StreamerInfo contained
// in a ROOT file.
package main // import "go-hep.org/x/hep/groot/cmd/root-gen-type"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func main() {
	log.SetPrefix("root-gen-type: ")
	log.SetFlags(0)

	var (
		typeNames = flag.String("t", ".*", "comma-separated list of (regexp) type names")
		pkgPath   = flag.String("p", "", "package import path")
		output    = flag.String("o", "", "output file name")
		verbose   = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-gen-type [options] input.root

ex:
 $> root-gen-type -p mypkg -t MyType -o streamers_gen.go ./input.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *typeNames == "" {
		flag.Usage()
		os.Exit(2)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")

	var (
		err error
		out io.WriteCloser
	)

	switch *output {
	case "":
		out = os.Stdout
	default:
		out, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	err = generate(out, *pkgPath, types, flag.Arg(0), *verbose)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func generate(w io.Writer, pkg string, types []string, fname string, verbose bool) error {
	f, err := groot.Open(fname)
	if err != nil {
		return err
	}

	g, err := rdict.NewGenGoType(pkg, f, verbose)
	if err != nil {
		return err
	}

	filters := make([]*regexp.Regexp, len(types))
	for i, t := range types {
		filters[i] = regexp.MustCompile(t)
	}

	accept := func(name string) string {
		for _, filter := range filters {
			if filter.MatchString(name) {
				return name
			}
		}
		return ""
	}
	for _, si := range f.StreamerInfos() {
		if t := accept(si.Name()); t != "" {
			err := g.Generate(t)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	buf, err := g.Format()
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
