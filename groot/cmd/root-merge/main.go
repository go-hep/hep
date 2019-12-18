// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-merge merges ROOT files' content into a merged ROOT file.
//
// Usage: root-merge [options] file1.root [file2.root [file3.root [...]]]
//
// ex:
//  $> root-merge -o out.root ./testdata/chain.flat.1.root ./testdata/chain.flat.2.root
//
// options:
//   -o string
//     	path to merged output ROOT file (default "out.root")
//   -v	enable verbose mode
//
package main // import "go-hep.org/x/hep/groot/cmd/root-merge"

import (
	"flag"
	"fmt"
	"log"
	"os"
	stdpath "path"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
)

func main() {
	log.SetPrefix("root-merge: ")
	log.SetFlags(0)

	var (
		oname   = flag.String("o", "out.root", "path to merged output ROOT file")
		verbose = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-merge [options] file1.root [file2.root [file3.root [...]]]

ex:
 $> root-merge -o out.root ./testdata/chain.flat.1.root ./testdata/chain.flat.2.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		log.Fatalf("missing input files")
	}

	fnames := flag.Args()

	err := rootmerge(*oname, fnames, *verbose)
	if err != nil {
		log.Fatalf("could not merge ROOT files: %+v", err)
	}
}

func rootmerge(oname string, fnames []string, verbose bool) error {
	o, err := groot.Create(oname)
	if err != nil {
		return xerrors.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	tsks, err := tasksFrom(o, fnames[0], verbose)
	if err != nil {
		return xerrors.Errorf("could not create merge tasks: %w", err)
	}

	for _, fname := range fnames[1:] {
		err := process(tsks, fname, verbose)
		if err != nil {
			return xerrors.Errorf("could not process ROOT file %q: %w", fname, err)
		}
	}

	for i := range tsks {
		tsk := &tsks[i]
		err := tsk.close(o)
		if err != nil {
			return xerrors.Errorf("could not close task %d (%s): %w", i, tsk.path(), err)
		}
	}

	err = o.Close()
	if err != nil {
		return xerrors.Errorf("could not close output ROOT file %q: %w", oname, err)
	}

	return nil
}

func process(tsks []task, fname string, verbose bool) error {
	if verbose {
		log.Printf("merging [%s]...", fname)
	}

	f, err := groot.Open(fname)
	if err != nil {
		return xerrors.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	for i := range tsks {
		tsk := &tsks[i]
		err = tsk.merge(f)
		if err != nil {
			return xerrors.Errorf("could not merge task %d (%s) for file %q: %w", i, tsk.path(), err)
		}
	}

	return nil
}

type task struct {
	dir string
	key string
	obj root.Object

	verbose bool
}

func tasksFrom(o *riofs.File, fname string, verbose bool) ([]task, error) {
	f, err := groot.Open(fname)
	if err != nil {
		return nil, xerrors.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	var tsks []task
	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		if err != nil {
			return err
		}
		name := path[len(f.Name()):]
		if name == "" {
			return nil
		}

		if _, ok := obj.(riofs.Directory); ok {
			_, err := riofs.Dir(o).Mkdir(name)
			if err != nil {
				return xerrors.Errorf("could not create dir %q in output ROOT file: %w", name, err)
			}
			if verbose {
				log.Printf("selecting %q", name)
			}
			return nil
		}

		if !acceptObj(obj) {
			return nil
		}
		if verbose {
			log.Printf("selecting %q", name)
		}

		var (
			dirName = stdpath.Dir(name)
			objName = stdpath.Base(name)
			dir     = riofs.Directory(o)
		)

		if dirName != "/" && dirName != "" {
			obj, err := riofs.Dir(o).Get(dirName)
			if err != nil {
				return xerrors.Errorf("could not get dir %q from output ROOT file: %w", dirName, err)
			}
			dir = obj.(riofs.Directory)
		}

		switch oo := obj.(type) {
		case rtree.Tree:
			w, err := rtree.NewWriter(dir, objName, rtree.WriteVarsFromTree(oo), rtree.WithTitle(oo.Title()))
			if err != nil {
				return xerrors.Errorf("could not create output ROOT tree %q: %w", name, err)
			}
			_, err = rtree.Copy(w, oo)
			if err != nil {
				return xerrors.Errorf("could not seed output ROOT tree %q: %w", name, err)
			}
			obj = w
		}

		tsks = append(tsks, task{
			dir:     dirName,
			key:     objName,
			obj:     obj,
			verbose: verbose,
		})
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf("could not inspect input ROOT file: %w", err)
	}

	if verbose {
		log.Printf("merging [%s]...", fname)
	}

	return tsks, nil
}

func (tsk *task) path() string {
	return stdpath.Join(tsk.dir, tsk.key)
}

func (tsk *task) merge(f *riofs.File) error {
	name := tsk.path()
	obj, err := riofs.Dir(f).Get(name)
	if err != nil {
		return xerrors.Errorf("could not get %q: %w", name, err)
	}

	err = mergeObj(tsk.obj, obj)
	if err != nil {
		return xerrors.Errorf("could not merge %q: %w", name, err)
	}

	return nil
}

func (tsk *task) close(f *riofs.File) error {
	var err error
	switch obj := tsk.obj.(type) {
	case rtree.Writer:
		err = obj.Close()
	default:
		err = riofs.Dir(f).Put(tsk.path(), tsk.obj)
	}

	if err != nil {
		return xerrors.Errorf("could not save %q (%T) to output ROOT file: %w", tsk.path(), tsk.obj, err)
	}

	return nil
}

func mergeObj(dst, src root.Object) error {
	var (
		rdst = dst.Class()
		rsrc = src.Class()
	)
	if rdst != rsrc {
		return xerrors.Errorf("types differ: dst=%T, src=%T", dst, src)
	}

	switch dst := dst.(type) {
	case rhist.H2:
		return mergeH2(dst, src.(rhist.H2))
	case rhist.H1:
		return mergeH1(dst, src.(rhist.H1))
	case root.Merger:
		return dst.ROOTMerge(src)
	default:
		return xerrors.Errorf("could not find suitable merge-API for (dst=%T, src=%T)", dst, src)
	}
}

func mergeH1(dst, src rhist.H1) error {
	panic("not implemented")
}

func mergeH2(dst, src rhist.H2) error {
	panic("not implemented")
}

func acceptObj(obj root.Object) bool {
	switch obj.(type) {
	case rtree.Tree:
		// need to specially handle rtree.Tree.
		// rtree.Tree does not implement root.Merger: only rtree.Writer does.
		return true
	case rhist.H1, rhist.H2:
		return true
	case root.Merger:
		return true
	default:
		return false
	}
}
