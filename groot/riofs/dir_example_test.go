// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs_test

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func Example_mkdir() {
	const fname = "../testdata/subdirs.root"
	defer os.Remove(fname)

	{
		w, err := groot.Create(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer w.Close()

		dir1, err := w.Mkdir("dir1")
		if err != nil {
			log.Fatal(err)
		}

		dir11, err := dir1.Mkdir("dir11")
		if err != nil {
			log.Fatal(err)
		}

		err = dir11.Put("obj1", rbase.NewObjString("data-obj1"))
		if err != nil {
			log.Fatal(err)
		}

		dir2, err := w.Mkdir("dir2")
		if err != nil {
			log.Fatal(err)
		}

		err = dir2.Put("obj2", rbase.NewObjString("data-obj2"))
		if err != nil {
			log.Fatal(err)
		}

		err = w.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	riofs.Walk(r, func(path string, obj root.Object, err error) error {
		fmt.Printf(">> %v\n", path)
		return err
	})

	// Output:
	// >> ../testdata/subdirs.root
	// >> ../testdata/subdirs.root/dir1
	// >> ../testdata/subdirs.root/dir1/dir11
	// >> ../testdata/subdirs.root/dir1/dir11/obj1
	// >> ../testdata/subdirs.root/dir2
	// >> ../testdata/subdirs.root/dir2/obj2
}
