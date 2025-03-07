// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnpy

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rtree"
)

// NewColumns returns all the ReadVars of the provided Tree as
// a slice of Columns.
//
// ReadVars that can not be represented as NumPy arrays are silently discarded.
func NewColumns(tree rtree.Tree) []Column {
	var (
		rvars = rtree.NewReadVars(tree)
		cols  []Column
	)

	for _, rvar := range rvars {
		rv := reflect.ValueOf(rvar.Value).Elem()
		switch rv.Kind() {
		case reflect.Chan, reflect.Interface,
			reflect.Struct, reflect.Slice, reflect.Map,
			reflect.Ptr, reflect.UnsafePointer:
			continue
		}
		cols = append(cols, Column{
			tree: tree,
			rvar: rvar,
			etyp: reflect.TypeOf(rvar.Value).Elem(),
		})
	}

	return cols
}

// Column provides a NumPy representation of a Branch or Leaf.
type Column struct {
	tree rtree.Tree
	rvar rtree.ReadVar
	etyp reflect.Type
}

// NewColumn returns the column with the provided name and tree.
//
// NewColumn returns an error if no branch or leaf could be found.
// NewColumn returns an error if the branch or leaf is of an unsupported type.
func NewColumn(tree rtree.Tree, rvar rtree.ReadVar) (Column, error) {
	var (
		rvars = rtree.NewReadVars(tree)
		idx   = -1
		col   Column
	)

	for i := range rvars {
		if rvars[i].Name == rvar.Name && (rvars[i].Leaf == rvar.Leaf || rvar.Leaf == "") {
			idx = i
			break
		}
	}

	if idx < 0 {
		name := rvar.Name
		if rvar.Leaf != "" {
			name += "." + rvar.Leaf
		}
		return col, fmt.Errorf("rnpy: no rvar named %q", name)
	}
	rvar = rvars[idx]

	rv := reflect.ValueOf(rvar.Value).Elem()
	switch rv.Kind() {
	case reflect.Chan, reflect.Interface,
		reflect.Struct, reflect.Slice, reflect.Map,
		reflect.Ptr, reflect.UnsafePointer:
		return col, fmt.Errorf("rnpy: invalid branch or leaf type %T", rv.Interface())
	}

	col = Column{
		tree: tree,
		rvar: rvar,
		etyp: reflect.TypeOf(rvar.Value).Elem(),
	}
	return col, nil
}

// Name returns the branch name this Column is bound to.
func (col Column) Name() string {
	return col.rvar.Name
}

// Slice reads the whole data slice from the underlying ROOT Tree
// into memory.
func (col Column) Slice() (sli any, err error) {
	r, err := rtree.NewReader(col.tree, []rtree.ReadVar{col.rvar})
	if err != nil {
		return nil, fmt.Errorf(
			"rnpy: could not create ROOT reader for %q: %w",
			col.rvar.Name, err,
		)
	}
	defer r.Close()

	var (
		n     = col.tree.Entries()
		rtyp  = reflect.SliceOf(col.etyp)
		data  = reflect.ValueOf(col.rvar.Value).Elem()
		slice = reflect.MakeSlice(rtyp, int(n), int(n))
		i     int
	)

	err = r.Read(func(ctx rtree.RCtx) error {
		slice.Index(i).Set(data)
		i++
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf(
			"rnpy: could not read ROOT data for %q: %w",
			col.rvar.Name, err,
		)
	}

	return slice.Interface(), nil
}
