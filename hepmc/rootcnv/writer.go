// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootcnv

import (
	"fmt"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hepmc"
)

// FlatTreeWriter writes HepMC events as a flat ROOT TTree.
type FlatTreeWriter struct {
	w rtree.Writer

	evt   event
	wvars []rtree.WriteVar
}

// NewFlatTreeWriter creates a new named tree under the dir directory.
func NewFlatTreeWriter(dir riofs.Directory, name string, opts ...rtree.WriteOption) (*FlatTreeWriter, error) {
	var w FlatTreeWriter
	w.wvars = rtree.WriteVarsFromStruct(&w.evt)
	tree, err := rtree.NewWriter(dir, name, w.wvars, opts...)
	if err != nil {
		return nil, fmt.Errorf("hepmc: could not create flat-tree writer %q: %w", name, err)
	}

	w.w = tree

	return &w, nil
}

func (w *FlatTreeWriter) Close() error {
	w.wvars = nil
	return w.w.Close()
}

func (w *FlatTreeWriter) Write(evt hepmc.Event) error {
	w.evt.reset()

	err := w.evt.read(&evt)
	if err != nil {
		return fmt.Errorf("hepmc: could not encode event to ROOT: %w", err)
	}

	_, err = w.w.Write()
	if err != nil {
		return fmt.Errorf("hepmc: could not write event to ROOT: %w", err)
	}

	return nil
}

var (
	_ hepmc.Writer = (*FlatTreeWriter)(nil)
)
