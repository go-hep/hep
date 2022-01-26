// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs_test

import (
	"fmt"
	"log"
	stdpath "path"
	"strings"

	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func ExampleWalk() {
	f, err := riofs.Open("../testdata/dirs-6.14.00.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Printf("visit all ROOT file tree:\n")
	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		fmt.Printf("%s (%s)\n", path, obj.Class())
		return nil
	})
	if err != nil {
		log.Fatalf("could not walk through file: %v", err)
	}

	fmt.Printf("visit only dir1:\n")
	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		if !(strings.HasPrefix(path, stdpath.Join(f.Name(), "dir1")) || path == f.Name()) {
			return riofs.SkipDir
		}
		fmt.Printf("%s (%s)\n", path, obj.Class())
		return nil
	})
	if err != nil {
		log.Fatalf("could not walk through file: %v", err)
	}

	// Output:
	// visit all ROOT file tree:
	// dirs-6.14.00.root (TFile)
	// dirs-6.14.00.root/dir1 (TDirectoryFile)
	// dirs-6.14.00.root/dir1/dir11 (TDirectoryFile)
	// dirs-6.14.00.root/dir1/dir11/h1 (TH1F)
	// dirs-6.14.00.root/dir2 (TDirectoryFile)
	// dirs-6.14.00.root/dir3 (TDirectoryFile)
	// visit only dir1:
	// dirs-6.14.00.root (TFile)
	// dirs-6.14.00.root/dir1 (TDirectoryFile)
	// dirs-6.14.00.root/dir1/dir11 (TDirectoryFile)
	// dirs-6.14.00.root/dir1/dir11/h1 (TH1F)
}

func ExampleGet() {
	f, err := riofs.Open("../testdata/dirs-6.14.00.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h1, err := riofs.Get[rhist.H1](f, "dir1/dir11/h1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("histo: %s (%s)", h1.Name(), h1.Class())

	// Output:
	// histo: h1 (TH1F)
}
