// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootio provides a pure-go read-access to ROOT files.
// rootio might, with time, provide write-access too.
//
// A typical usage is as follow:
//
//   f, err := rootio.Open("ntup.root")
//   obj, err := f.Get("tree")
//   tree := obj.(*rootio.Tree)
//   fmt.Printf("entries= %v\n", t.Entries())
package rootio // import "go-hep.org/x/hep/rootio"
