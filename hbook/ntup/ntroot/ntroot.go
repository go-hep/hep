// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntroot provides convenience functions to access ROOT trees as n-tuple
// data.
//
// Example:
//
//	nt, err := ntroot.Open("testdata/simple.root", "mytree")
//	if err != nil {
//	    log.Fatalf("%+v", err)
//	}
//	defer nt.DB().Close()
package ntroot // import "go-hep.org/x/hep/hbook/ntup/ntroot"

import (
	"fmt"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rsql/rsqldrv"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook/ntup"
)

// Open opens the named ROOT file in read-only mode and returns an n-tuple
// connected to the named tree.
func Open(name, tree string) (*ntup.Ntuple, error) {
	f, err := groot.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not open ROOT file: %w", err)
	}
	defer f.Close()

	obj, err := riofs.Dir(f).Get(tree)
	if err != nil {
		return nil, fmt.Errorf("could not find ROOT tree %q: %w", tree, err)
	}
	if _, ok := obj.(rtree.Tree); !ok {
		return nil, fmt.Errorf("ROOT object %q is not a tree", tree)
	}

	db, err := rsqldrv.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not open ROOT db: %w", err)
	}

	nt, err := ntup.Open(db, tree)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("could not open n-tuple %q: %w", tree, err)
	}
	return nt, nil
}
