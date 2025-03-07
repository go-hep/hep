// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
)

// Dump dumps the content of the fname ROOT file to the provided io.Writer.
// If deep is true, Dump will recursively inspect directories and trees.
// Dump only display the content of ROOT objects satisfying the provided filter function.
//
// If filter is nil, Dump will consider all ROOT objects.
func Dump(w io.Writer, fname string, deep bool, filter func(name string) bool) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open file with read-access: %w", err)
	}
	defer f.Close()

	if filter == nil {
		filter = func(string) bool { return true }
	}

	cmd := dumpCmd{
		w:     w,
		deep:  deep,
		match: filter,
	}
	return cmd.dumpDir(f)
}

type dumpCmd struct {
	w     io.Writer
	deep  bool
	match func(name string) bool
}

func (cmd *dumpCmd) dumpDir(dir riofs.Directory) error {
	for i, key := range dir.Keys() {
		fmt.Fprintf(cmd.w, "key[%03d]: %s;%d %q (%s)", i, key.Name(), key.Cycle(), key.Title(), key.ClassName())
		if !(cmd.deep && cmd.match(key.Name())) {
			fmt.Fprint(cmd.w, "\n")
			continue
		}
		obj, err := key.Object()
		if err != nil {
			return fmt.Errorf("could not decode object %q from dir %q: %w", key.Name(), dir.(root.Named).Name(), err)
		}
		err = cmd.dumpObj(obj)
		if errors.Is(err, errIgnoreKey) {
			continue
		}
		if err != nil {
			return fmt.Errorf("error dumping key %q: %w", key.Name(), err)
		}
	}
	return nil
}

var errIgnoreKey = fmt.Errorf("rcmd: ignore key")

func (cmd *dumpCmd) dumpObj(obj root.Object) error {
	var err error
	switch obj := obj.(type) {
	case rtree.Tree:
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpTree(obj)
	case riofs.Directory:
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpDir(obj)
	case rhist.H2:
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpH2(obj)
	case rhist.H1: // keep after rhist.H2
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpH1(obj)
	case rhist.Graph:
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpGraph(obj)
	case rhist.MultiGraph:
		for _, g := range obj.Graphs() {
			fmt.Fprintf(cmd.w, "\n")
			err = cmd.dumpGraph(g)
			if err != nil {
				return err
			}
		}
	case root.List:
		fmt.Fprintf(cmd.w, "\n")
		err = cmd.dumpList(obj)
	case *rdict.Object:
		fmt.Fprintf(cmd.w, " => %v\n", obj)
	case fmt.Stringer:
		fmt.Fprintf(cmd.w, " => %q\n", obj.String())
	case root.ObjectFinder:
		keys := obj.Keys()
		if len(keys) == 0 {
			return err
		}
		fmt.Fprintf(cmd.w, "\n")
		for _, name := range keys {
			sub, err := obj.Get(name)
			if err != nil {
				return err
			}
			err = cmd.dumpObj(sub)
			if err != nil {
				return err
			}
		}
	default:
		fmt.Fprintf(cmd.w, " => ignoring key of type %T\n", obj)
		return errIgnoreKey
	}
	return err
}

func (cmd *dumpCmd) dumpList(lst root.List) error {
	for i := range lst.Len() {
		fmt.Fprintf(cmd.w, "lst[%s][%d]: ", lst.Name(), i)
		err := cmd.dumpObj(lst.At(i))
		if err != nil && !errors.Is(err, errIgnoreKey) {
			return fmt.Errorf("could not dump list: %w", err)
		}
	}
	return nil
}

func (cmd *dumpCmd) dumpTree(t rtree.Tree) error {
	vars := rtree.NewReadVars(t)
	r, err := rtree.NewReader(t, vars)
	if err != nil {
		return fmt.Errorf("could not create reader: %w", err)
	}
	defer r.Close()

	names := make([][]byte, len(vars))
	for i, v := range vars {
		name := v.Name
		if v.Leaf != "" && v.Leaf != v.Name {
			name = v.Name + "." + v.Leaf
		}
		names[i] = []byte(name)
	}

	// FIXME(sbinet): don't use a "global" buffer for when rtree.Reader reads multiple
	// events in parallel.
	buf := make([]byte, 0, 8*1024)
	hdr := make([]byte, 0, 6)
	err = r.Read(func(rctx rtree.RCtx) error {
		hdr = hdr[:0]
		hdr = append(hdr, '[')
		switch {
		case rctx.Entry < 10:
			hdr = append(hdr, '0', '0')
		case rctx.Entry < 100:
			hdr = append(hdr, '0')
		}
		hdr = strconv.AppendInt(hdr, rctx.Entry, 10)
		hdr = append(hdr, ']', '[')
		for i, v := range vars {
			buf = buf[:0]
			buf = append(buf, hdr...)
			buf = append(buf, names[i]...)
			buf = append(buf, ']', ':', ' ')
			rv := reflect.Indirect(reflect.ValueOf(v.Value))
			// All of this is a convoluted (but efficient) way to do:
			//  fmt.Fprintf(cmd.w, "[%03d][%s]: %v\n", rctx.Entry, name, rv.Interface())
			buf = append(buf, fmt.Sprintf("%v\n", rv.Interface())...)
			_, err = cmd.w.Write(buf)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("rcmd: could not read through tree: %w", err)
	}
	return nil
}

func (cmd *dumpCmd) dumpH1(h1 rhist.H1) error {
	h := rootcnv.H1D(h1)
	return yodacnv.Write(cmd.w, h)
}

func (cmd *dumpCmd) dumpH2(h2 rhist.H2) error {
	h := rootcnv.H2D(h2)
	return yodacnv.Write(cmd.w, h)
}

func (cmd *dumpCmd) dumpGraph(gr rhist.Graph) error {
	g := rootcnv.S2D(gr)
	return yodacnv.Write(cmd.w, g)
}
