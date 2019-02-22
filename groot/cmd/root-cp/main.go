// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-cp selects and copies keys from a ROOT file to another ROOT file.
//
// Usage: root-cp [options] file1.root[:REGEXP] [file2.root[:REGEXP] [...]] out.root
//
// ex:
//
//  $> root-cp f.root out.root
//  $> root-cp f1.root f2.root f3.root out.root
//  $> root-cp f1.root:hist.* f2.root:h2 out.root
//
package main // import "go-hep.org/x/hep/groot/cmd/root-cp"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
)

func main() {
	log.SetPrefix("root-cp: ")
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-cp [options] file1.root[:REGEXP] [file2.root[:REGEXP] [...]] out.root

ex:
 $> root-cp f.root out.root
 $> root-cp f1.root f2.root f3.root out.root
 $> root-cp f1.root:hist.* f2.root:h2 out.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 2 {
		log.Printf("error: you need to give input and output ROOT files\n\n")
		flag.Usage()
		os.Exit(1)
	}

	dst := flag.Arg(flag.NArg() - 1)
	srcs := flag.Args()[:flag.NArg()-1]

	err := rootcp(dst, srcs)
	if err != nil {
		log.Fatal(err)
	}
}

func rootcp(oname string, fnames []string) error {
	o, err := groot.Create(oname)
	if err != nil {
		return errors.Errorf("could not create output ROOT file %q: %v", oname, err)
	}
	defer o.Close()

	for _, arg := range fnames {
		err := process(o, arg)
		if err != nil {
			return err
		}
	}

	err = o.Close()
	if err != nil {
		return errors.Errorf("could not close output ROOT file %q: %v", oname, err)
	}
	return nil
}

func process(o *riofs.File, arg string) error {
	log.Printf("copying %q...", arg)

	fname, sel, err := splitArg(arg)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(sel)

	f, err := groot.Open(fname)
	if err != nil {
		return errors.Errorf("could not open input ROOT file %q: %v", fname, err)
	}
	defer f.Close()

	for _, k := range f.Keys() {
		if !re.MatchString(k.Name()) {
			continue
		}

		v, err := k.Object()
		if err != nil {
			return errors.Errorf("could not load object %q from file %q: %v", k.Name(), fname, err)
		}

		err = o.Put(k.Name(), v)
		if err != nil {
			return errors.Errorf("could not save object %q to output file: %v", k.Name(), err)
		}
	}

	return nil
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
		return "", "", errors.Errorf("root-cp: too many ':' in %q", cmd)
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
	return fname, sel, err
}
