// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package riofs contains the types and low-level functions to deal with opening
// and creating ROOT files, and decoding the internal structure of ROOT files.
//
// Users should prefer to use the groot package to open or create ROOT files instead of this one.
package riofs // import "go-hep.org/x/hep/groot/riofs"

import (
	"go-hep.org/x/hep/groot/root"
)

//go:generate go run ./gen-code.go
//go:generate go run ./gendata/gen-dirs.go -f ../testdata/dirs.root
//go:generate go run ./gendata/gen-evnt-tree.go -f ../testdata/small-evnt-tree-nosplit.root
//go:generate go run ./gendata/gen-evnt-tree.go -f ../testdata/small-evnt-tree-fullsplit.root -split=99
//go:generate go run ./gendata/gen-flat-tree.go -f ../testdata/leaves.root
//go:generate go run ./gendata/gen-map-tree.go -f ../testdata/std-map-split0.root -split=0
//go:generate go run ./gendata/gen-map-tree.go -f ../testdata/std-map-split1.root -split=1
//go:generate go run ./gendata/gen-map-tree.go -f ../testdata/std-map-split2.root -split=2
//go:generate go run ./gendata/gen-multi-leaves-tree.go -f ../testdata/padding.root
//go:generate go run ./gendata/gen-join-trees.go -d ../testdata
//go:generate go run ./gendata/gen-bitset-tree.go -f ../testdata/std-bitset.root
//go:generate go run ./gendata/gen-std-containers.go -f ../testdata/std-containers-split00.root -split=0
//go:generate go run ./gendata/gen-tlv.go -f ../testdata/tlv-split00.root -split=0
//go:generate go run ./gendata/gen-tlv.go -f ../testdata/tlv-split01.root -split=1
//go:generate go run ./gendata/gen-tlv.go -f ../testdata/tlv-split99.root -split=99
//go:generate go run ./gendata/gen-tntuple.go  -f ../testdata/tntuple.root
//go:generate go run ./gendata/gen-tntupled.go -f ../testdata/tntupled.root
//go:generate go run ./gendata/gen-tformula.go -f ../testdata/tformula.root
//go:generate go run ./gendata/gen-tconflvl.go -f ../testdata/tconfidence-level.root
//go:generate go run ./gendata/gen-tprofile.go -f ../testdata/tprofile.root
//go:generate go run ./gendata/gen-tgme.go -f ../testdata/tgme.root
//go:generate go run ./gendata/gen-tdatime.go -f ../testdata/tdatime.root
//go:generate go run ./gendata/gen-base.go -f ../testdata/tbase.root
//go:generate go run ./gendata/gen-tcanvas.go -f ../testdata/tcanvas.root

// Directory describes a ROOT directory structure in memory.
type Directory interface {
	// Get returns the object identified by namecycle
	//   namecycle has the format name;cycle
	//   name  = * is illegal, cycle = * is illegal
	//   cycle = "" or cycle = 9999 ==> apply to a memory object
	//
	//   examples:
	//     foo   : get object named foo in memory
	//             if object is not in memory, try with highest cycle from file
	//     foo;1 : get cycle 1 of foo on file
	Get(namecycle string) (root.Object, error)

	// Put puts the object v under the key with the given name.
	Put(name string, v root.Object) error

	// Keys returns the list of keys being held by this directory.
	Keys() []Key

	// Mkdir creates a new subdirectory
	Mkdir(name string) (Directory, error)

	// Parent returns the directory holding this directory.
	// Parent returns nil if this is the top-level directory.
	Parent() Directory
}

// SetFiler is a simple interface to establish File ownership.
type SetFiler interface {
	SetFile(f *File)
}
