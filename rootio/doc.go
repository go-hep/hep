// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootio provides a pure-go read-access to ROOT files.
// rootio might, with time, provide write-access too.
//
// A typical usage is as follows:
//
//   f, err := rootio.Open("ntup.root")
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer f.Close()
//
//   obj, err := f.Get("tree")
//   if err != nil {
//       log.Fatal(err)
//   }
//   tree := obj.(rootio.Tree)
//   fmt.Printf("entries= %v\n", tree.Entries())
//
// More complete examples on how to iterate over the content of a Tree can
// be found in the examples attached to rootio.TreeScanner and rootio.Scanner:
// https://godoc.org/go-hep.org/x/hep/rootio#pkg-examples
//
// Another possibility is to look at:
// https://godoc.org/go-hep.org/x/hep/rootio/cmd/root-ls,
// a command that inspects the content of ROOT files.
package rootio // import "go-hep.org/x/hep/rootio"
