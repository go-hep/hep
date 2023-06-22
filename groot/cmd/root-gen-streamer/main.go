// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command root-gen-streamer generates a StreamerInfo for ROOT and user types.
package main // import "go-hep.org/x/hep/groot/cmd/root-gen-streamer"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go-hep.org/x/hep/groot/rdict"
)

var (
	typeNames = flag.String("t", "", "comma-separated list of type names")
	pkgPath   = flag.String("p", "", "package import path")
	output    = flag.String("o", "", "output file name")
	verbose   = flag.Bool("v", false, "enable verbose mode")
)

func main() {
	log.SetPrefix("root-gen-streamer: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-gen-streamer [options]

ex:
 $> root-gen-streamer -p image -t Point -o streamers_gen.go
 $> root-gen-streamer -p go-hep.org/x/hep/hbook -t Dist0D,Dist1D,Dist2D -o foo_streamer_gen.go

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

	types := strings.Split(*typeNames, ",")

	var (
		err error
		out io.WriteCloser
		buf = new(bytes.Buffer)
	)

	err = generate(buf, *pkgPath, types)
	if err != nil {
		log.Fatal(err)
	}

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

	_, err = io.Copy(out, buf)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func generate(w io.Writer, pkg string, types []string) error {
	g, err := rdict.NewGenStreamer(pkg, *verbose)
	if err != nil {
		return err
	}

	for _, t := range types {
		err := g.Generate(t)
		if err != nil {
			log.Fatal(err)
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
