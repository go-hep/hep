// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

// ListOption controls how List behaves.
type ListOption func(*lsCmd)

type lsCmd struct {
	w io.Writer

	streamers bool
	trees     bool
}

// ListStreamers enables the display of streamer informations
// contained in the provided ROOT file.
func ListStreamers(v bool) ListOption {
	return func(cmd *lsCmd) {
		cmd.streamers = v
	}
}

// ListTrees enables the detailed display of trees contained in
// the provided ROOT file.
func ListTrees(v bool) ListOption {
	return func(cmd *lsCmd) {
		cmd.trees = v
	}
}

// List displays the summary content of the named ROOT file into the
// provided io Writer.
//
// List's behaviour can be customized with a set of optional ListOptions.
func List(w io.Writer, fname string, opts ...ListOption) error {
	cmd := lsCmd{
		w:         w,
		streamers: false,
		trees:     false,
	}

	for _, opt := range opts {
		opt(&cmd)
	}

	return cmd.ls(fname)
}

func (ls lsCmd) ls(fname string) error {
	fmt.Fprintf(ls.w, "=== [%s] ===\n", fname)
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(ls.w, "version: %v\n", f.Version())
	if ls.streamers {
		fmt.Fprintf(ls.w, "streamer-infos:\n")
		sinfos := f.StreamerInfos()
		for _, v := range sinfos {
			name := v.Name()
			fmt.Fprintf(ls.w, " StreamerInfo for %q version=%d title=%q\n", name, v.ClassVersion(), v.Title())
			w := tabwriter.NewWriter(ls.w, 8, 4, 1, ' ', 0)
			for _, elm := range v.Elements() {
				fmt.Fprintf(w, "  %s\t%s\toffset=%3d\ttype=%3d\tsize=%3d\t %s\n", elm.TypeName(), elm.Name(), elm.Offset(), elm.Type(), elm.Size(), elm.Title())
			}
			w.Flush()
		}
		fmt.Fprintf(ls.w, "---\n")
	}

	w := tabwriter.NewWriter(ls.w, 8, 4, 1, ' ', 0)
	for _, k := range f.Keys() {
		ls.walk(w, k)
	}
	w.Flush()

	return nil
}

func (ls lsCmd) walk(w io.Writer, k riofs.Key) {
	if ls.trees && isTreelike(k.ClassName()) {
		obj := k.Value()
		tree, ok := obj.(rtree.Tree)
		if ok {
			w := newWindent(2, w)
			fmt.Fprintf(w, "%s\t%s\t%s\t(entries=%d)\n", k.ClassName(), k.Name(), k.Title(), tree.Entries())
			displayBranches(w, tree, 2)
			w.Flush()
			return
		}
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t(cycle=%d)\n", k.ClassName(), k.Name(), k.Title(), k.Cycle())
	if isDirlike(k.ClassName()) {
		obj := k.Value()
		if dir, ok := obj.(riofs.Directory); ok {
			w := newWindent(2, w)
			for _, k := range dir.Keys() {
				ls.walk(w, k)
			}
			w.Flush()
		}
	}
}

func isDirlike(class string) bool {
	switch class {
	case "TDirectory", "TDirectoryFile":
		return true
	}
	return false
}

func isTreelike(class string) bool {
	switch class {
	case "TTree", "TTreeSQL", "TChain", "TNtuple", "TNtupleD":
		return true
	}
	return false
}

type windent struct {
	hdr []byte
	w   io.Writer
}

func newWindent(n int, w io.Writer) *windent {
	return &windent{
		hdr: bytes.Repeat([]byte(" "), n),
		w:   w,
	}
}

func (w *windent) Write(data []byte) (int, error) {
	return w.w.Write(append(w.hdr, data...))
}

func (w *windent) Flush() error {
	ww, ok := w.w.(interface{ Flush() error })
	if !ok {
		return nil
	}
	return ww.Flush()
}

type brancher interface {
	Branches() []rtree.Branch
}

func displayBranches(w io.Writer, bres brancher, indent int) {
	branches := bres.Branches()
	if len(branches) <= 0 {
		return
	}
	ww := newWindent(indent, w)
	for _, b := range branches {
		var (
			name  = clip(b.Name(), 60)
			title = clip(b.Title(), 50)
			class = clip(b.Class(), 20)
		)
		fmt.Fprintf(ww, "%s\t%q\t%v\n", name, title, class)
		displayBranches(ww, b, 2)
	}
	ww.Flush()
}

func clip(s string, n int) string {
	if len(s) > n {
		s = s[:n-5] + "[...]"
	}
	return s
}
