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

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
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
		return xerrors.Errorf("could not create output ROOT file %q: %w", oname, err)
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
		return xerrors.Errorf("could not close output ROOT file %q: %w", oname, err)
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
		return xerrors.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	for _, k := range f.Keys() {
		if !re.MatchString(k.Name()) {
			continue
		}

		v, err := k.Object()
		if err != nil {
			return xerrors.Errorf("could not load object %q from file %q: %w", k.Name(), fname, err)
		}

		err = copyObj(o, k.Name(), v)
		if err != nil {
			return xerrors.Errorf("could not copy object %q to output file: %w", k.Name(), err)
		}
	}

	return nil
}

func copyObj(o *riofs.File, k string, obj root.Object) error {
	var err error
	switch obj := obj.(type) {
	case rtree.Tree:
		err = copyTree(o, k, obj)
	default:
		err = o.Put(k, obj)
	}

	if err != nil {
		return xerrors.Errorf("could not save object %q to output file: %w", k, err)
	}

	return nil
}

func copyTree(o *riofs.File, name string, tree rtree.Tree) error {
	dst, err := rtree.NewWriter(o, name, rtree.WriteVarsFromTree(tree))
	if err != nil {
		return xerrors.Errorf("could not create output copy tree: %w", err)
	}
	_, err = rtree.Copy(dst, tree)
	if err != nil {
		return xerrors.Errorf("could not copy tree %q: %w", name, err)
	}

	err = dst.Close()
	if err != nil {
		return xerrors.Errorf("could not close copy tree %q: %w", name, err)
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
		return "", "", xerrors.Errorf("root-cp: too many ':' in %q", cmd)
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
