// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command root-gen-rfunc generates a rfunc.Formula based on a function
// signature or an already existing function.
package main // import "go-hep.org/x/hep/groot/cmd/root-gen-rfunc"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	log.SetPrefix("gen-rfunc: ")
	log.SetFlags(0)

	var (
		pkg = flag.String("p", "", "import path of the package holding the function definition (if any)")
		fct = flag.String("f", "", "name of the function definition or signature of the function")
		n   = flag.String("n", "", "name of the output function")
		o   = flag.String("o", "", "path to output file (if any)")
	)

	flag.Parse()
	err := generate(*o, *pkg, *fct, *n, flag.Usage)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func generate(oname, pkg, fct, name string, usage func()) error {
	switch {
	case pkg == "" && fct == "":
		usage()
		return fmt.Errorf("missing package import path and/or function name")
	case pkg == "" && !strings.Contains(fct, "func"):
		usage()
		return fmt.Errorf("missing function signature")
	}

	var out io.Writer
	switch oname {
	case "":
		out = os.Stdout
	default:
		f, err := os.Create(oname)
		if err != nil {
			return fmt.Errorf("could not create output file: %w", err)
		}
		defer f.Close()
		out = f
	}

	err := genroot.GenRFunc(out, genroot.RFunc{
		Path: pkg,
		Name: name,
		Def:  fct,
	})
	if err != nil {
		return fmt.Errorf("could not generate rfunc formula: %w", err)
	}

	return nil
}
