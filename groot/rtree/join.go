// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"strings"

	"go-hep.org/x/hep/groot/root"
)

type join struct {
	name     string
	title    string
	trees    []Tree
	branches []Branch
	leaves   []Leaf
	bmap     map[string]Branch
	lmap     map[string]Leaf
}

// Join returns a new Tree that represents the logical join of the input trees.
// The returned tree will contain all the columns of all the input trees.
// Join errors out if the input slice of trees is empty.
// Join errors out if the input trees do not have the same amount of entries.
// Join errors out if two trees have each a column with the same name.
func Join(trees ...Tree) (Tree, error) {
	if len(trees) == 0 {
		return nil, fmt.Errorf("rtree: no trees to join")
	}

	nevts := trees[0].Entries()
	for _, t := range trees[1:] {
		if t.Entries() != nevts {
			return nil, fmt.Errorf(
				"rtree: invalid number of entries in tree %s (got=%d, want=%d)",
				t.Name(), t.Entries(), nevts,
			)
		}
	}

	var (
		bset     = make([]map[string]struct{}, len(trees))
		branches []Branch
		leaves   []Leaf
		names    = make([]string, len(trees))
		titles   = make([]string, len(trees))
	)
	for i, t := range trees {
		names[i] = t.Name()
		titles[i] = t.Title()
		bset[i] = make(map[string]struct{}, len(t.Branches()))
		for _, b := range t.Branches() {
			bset[i][b.Name()] = struct{}{}
		}
		branches = append(branches, t.Branches()...)
		leaves = append(leaves, t.Leaves()...)
	}

	for i, ti := range trees {
		bsi := bset[i]
		for j, tj := range trees[i+1:] {
			bsj := bset[j+i+1]
			for ki := range bsi {
				if _, dup := bsj[ki]; dup {
					return nil, fmt.Errorf(
						"rtree: trees %s and %s both have a branch named %s",
						ti.Name(), tj.Name(), ki,
					)
				}
			}
		}
	}

	tree := &join{
		name:     "join_" + strings.Join(names, "_"),
		title:    strings.Join(titles, ", "),
		trees:    trees,
		branches: branches,
		leaves:   leaves,
		bmap:     make(map[string]Branch, len(branches)),
		lmap:     make(map[string]Leaf, len(leaves)),
	}

	for _, b := range tree.branches {
		tree.bmap[b.Name()] = b
	}
	for _, l := range tree.leaves {
		tree.lmap[l.Name()] = l
	}

	return tree, nil
}

// Class returns the ROOT class of the argument.
func (*join) Class() string {
	return "TJoin"
}

// Name returns the name of the ROOT objet in the argument.
func (t *join) Name() string {
	return t.name
}

// Title returns the title of the ROOT object in the argument.
func (t *join) Title() string {
	return t.title
}

// Entries returns the total number of entries.
func (t *join) Entries() int64 {
	return t.trees[0].Entries()
}

// Branches returns the list of branches.
func (t *join) Branches() []Branch {
	return t.branches
}

// Branch returns the branch whose name is the argument.
func (t *join) Branch(name string) Branch {
	return t.bmap[name]
}

// Leaves returns direct pointers to individual branch leaves.
func (t *join) Leaves() []Leaf {
	return t.leaves
}

// Leaf returns the leaf whose name is the argument.
func (t *join) Leaf(name string) Leaf {
	return t.lmap[name]
}

var (
	_ root.Object = (*chain)(nil)
	_ root.Named  = (*chain)(nil)
	_ Tree        = (*chain)(nil)
)
var (
	_ Tree = (*join)(nil)
)
