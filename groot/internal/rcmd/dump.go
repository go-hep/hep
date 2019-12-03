// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"fmt"
	"io"
	"reflect"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"golang.org/x/xerrors"
)

// Dump dumps the content of the fname ROOT file to the provided io.Writer.
// If deep is true, Dump will recursively inspect directories and trees.
// Dump only display the content of ROOT objects satisfying the provided filter function.
//
// If filter is nil, Dump will consider all ROOT objects.
func Dump(w io.Writer, fname string, deep bool, filter func(name string) bool) error {
	f, err := groot.Open(fname)
	if err != nil {
		return xerrors.Errorf("could not open file with read-access: %w", err)
	}
	defer f.Close()

	if filter == nil {
		filter = func(string) bool { return true }
	}
	return dumpDir(w, f, deep, filter)
}

func dumpDir(w io.Writer, dir riofs.Directory, deep bool, match func(name string) bool) error {
	for i, key := range dir.Keys() {
		fmt.Fprintf(w, "key[%03d]: %s;%d %q (%s)", i, key.Name(), key.Cycle(), key.Title(), key.ClassName())
		if !(deep && match(key.Name())) {
			fmt.Fprint(w, "\n")
			continue
		}
		obj, err := key.Object()
		if err != nil {
			return xerrors.Errorf("could not decode object %q from dir %q: %w", key.Name(), dir.(root.Named).Name(), err)
		}
		err = dumpObj(w, obj, deep, match)
		if xerrors.Is(err, errIgnoreKey) {
			continue
		}
		if err != nil {
			return xerrors.Errorf("error dumping key %q: %w", key.Name(), err)
		}
	}
	return nil
}

var errIgnoreKey = xerrors.Errorf("rcmd: ignore key")

func dumpObj(w io.Writer, obj root.Object, deep bool, match func(name string) bool) error {
	var err error
	switch obj := obj.(type) {
	case rtree.Tree:
		fmt.Fprintf(w, "\n")
		err = dumpTree(w, obj)
	case riofs.Directory:
		fmt.Fprintf(w, "\n")
		err = dumpDir(w, obj, deep, match)
	case rhist.H2:
		fmt.Fprintf(w, "\n")
		err = dumpH2(w, obj)
	case rhist.H1: // keep after rhist.H2
		fmt.Fprintf(w, "\n")
		err = dumpH1(w, obj)
	case rhist.Graph:
		fmt.Fprintf(w, "\n")
		err = dumpGraph(w, obj)
	case root.List:
		fmt.Fprintf(w, "\n")
		err = dumpList(w, obj, deep, match)
	case *rdict.Object:
		fmt.Fprintf(w, " => %v\n", obj)
	case fmt.Stringer:
		fmt.Fprintf(w, " => %q\n", obj.String())
	default:
		fmt.Fprintf(w, " => ignoring key of type %T\n", obj)
		return errIgnoreKey
	}
	return err
}

func dumpList(w io.Writer, lst root.List, deep bool, match func(name string) bool) error {
	for i := 0; i < lst.Len(); i++ {
		fmt.Fprintf(w, "lst[%s][%d]: ", lst.Name(), i)
		err := dumpObj(w, lst.At(i), deep, match)
		if err != nil && !xerrors.Is(err, errIgnoreKey) {
			return xerrors.Errorf("could not dump list: %w", err)
		}
	}
	return nil
}

func dumpTree(w io.Writer, t rtree.Tree) error {

	vars := rtree.NewScanVars(t)
	sc, err := rtree.NewScannerVars(t, vars...)
	if err != nil {
		return xerrors.Errorf("could not create scanner-vars: %w", err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan()
		if err != nil {
			return xerrors.Errorf("error scanning entry %d: %w", sc.Entry(), err)
		}
		for _, v := range vars {
			rv := reflect.Indirect(reflect.ValueOf(v.Value))
			fmt.Fprintf(w, "[%03d][%s]: %v\n", sc.Entry(), v.Name, rv.Interface())
		}
	}
	return nil
}

func dumpH1(w io.Writer, h1 rhist.H1) error {
	h, err := rootcnv.H1D(h1)
	if err != nil {
		return xerrors.Errorf("could not convert TH1x to hbook: %w", err)
	}
	return yodacnv.Write(w, h)
}

func dumpH2(w io.Writer, h2 rhist.H2) error {
	h, err := rootcnv.H2D(h2)
	if err != nil {
		return xerrors.Errorf("could not convert TH2x to hbook: %w", err)
	}
	return yodacnv.Write(w, h)
}

func dumpGraph(w io.Writer, gr rhist.Graph) error {
	g, err := rootcnv.S2D(gr)
	if err != nil {
		return xerrors.Errorf("could not convert TGraph to hbook: %w", err)
	}
	return yodacnv.Write(w, g)
}
