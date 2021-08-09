// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs_test

import (
	"fmt"
	"log"
	"os"
	stdpath "path"

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

	err = riofs.Walk(r, func(path string, obj root.Object, err error) error {
		fmt.Printf(">> %v\n", path)
		return err
	})
	if err != nil {
		log.Fatalf("could not walk ROOT file: %+v", err)
	}

	// Output:
	// >> ../testdata/subdirs.root
	// >> ../testdata/subdirs.root/dir1
	// >> ../testdata/subdirs.root/dir1/dir11
	// >> ../testdata/subdirs.root/dir1/dir11/obj1
	// >> ../testdata/subdirs.root/dir2
	// >> ../testdata/subdirs.root/dir2/obj2
}

func Example_recursivePut() {
	dir, err := os.MkdirTemp("", "groot-riofs-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := stdpath.Join(dir, "dirs.root")
	f, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rd := riofs.Dir(f)

	// create obj-1, put it under dir1/dir11, create all intermediate directories.
	err = rd.Put("dir1/dir11/obj-1", rbase.NewObjString("obj-1"))
	if err != nil {
		log.Fatal(err)
	}

	// create obj-2
	err = rd.Put("/dir2/dir22/dir222/obj-2", rbase.NewObjString("obj-2"))
	if err != nil {
		log.Fatal(err)
	}

	// update obj-1
	err = rd.Put("dir1/dir11/obj-1", rbase.NewObjString("obj-1-1"))
	if err != nil {
		log.Fatal(err)
	}

	o, err := rd.Get("dir1/dir11/obj-1;1")
	if err != nil {
		log.Fatal(err)
	}
	if got, want := o.(*rbase.ObjString).String(), "obj-1"; got != want {
		log.Fatalf("invalid obj-1;1 value. got=%q, want=%q", got, want)
	}

	o, err = rd.Get("dir1/dir11/obj-1;2")
	if err != nil {
		log.Fatal(err)
	}
	if got, want := o.(*rbase.ObjString).String(), "obj-1-1"; got != want {
		log.Fatalf("invalid obj-1;1 value. got=%q, want=%q", got, want)
	}

	o, err = rd.Get("dir1/dir11/obj-1")
	if err != nil {
		log.Fatal(err)
	}
	if got, want := o.(*rbase.ObjString).String(), "obj-1-1"; got != want {
		log.Fatalf("invalid obj-1 value. got=%q, want=%q", got, want)
	}

	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		name := path[len(fname):]
		if name == "" {
			return err
		}
		switch o := obj.(type) {
		case *rbase.ObjString:
			fmt.Printf(">> %v -- value=%q\n", name, o)
		default:
			fmt.Printf(">> %v\n", name)
		}
		return err
	})
	if err != nil {
		log.Fatalf("could not walk ROOT file: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close ROOT file: %v", err)
	}

	// Output:
	// >> /dir1
	// >> /dir1/dir11
	// >> /dir1/dir11/obj-1 -- value="obj-1-1"
	// >> /dir2
	// >> /dir2/dir22
	// >> /dir2/dir22/dir222
	// >> /dir2/dir22/dir222/obj-2 -- value="obj-2"
}

func Example_recursiveMkdir() {
	dir, err := os.MkdirTemp("", "groot-riofs-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fname := stdpath.Join(dir, "dirs.root")
	f, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rd := riofs.Dir(f)

	for _, path := range []string{
		"dir1/dir11/dir111",
		"/dir2/dir22/dir000",
		"dir2/dir22/dir222",
	} {
		_, err = rd.Mkdir(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rd.Put("/dir1/dir11/obj-1", rbase.NewObjString("obj-1"))
	if err != nil {
		log.Fatal(err)
	}

	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		name := path[len(fname):]
		if name == "" {
			return err
		}
		switch o := obj.(type) {
		case *rbase.ObjString:
			fmt.Printf(">> %v -- value=%q\n", name, o)
		default:
			fmt.Printf(">> %v\n", name)
		}
		return err
	})
	if err != nil {
		log.Fatalf("could not walk ROOT file: %+v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close ROOT file: %v", err)
	}

	// Output:
	// >> /dir1
	// >> /dir1/dir11
	// >> /dir1/dir11/dir111
	// >> /dir1/dir11/obj-1 -- value="obj-1"
	// >> /dir2
	// >> /dir2/dir22
	// >> /dir2/dir22/dir000
	// >> /dir2/dir22/dir222
}
