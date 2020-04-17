// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"fmt"
	"log"
	stdpath "path"
	"regexp"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

// Copy copies the content of the ROOT files fnames into the output
// ROOT file named oname.
func Copy(oname string, fnames []string) error {
	o, err := groot.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	var cmd copyCmd
	for _, arg := range fnames {
		err := cmd.process(o, arg)
		if err != nil {
			return err
		}
	}

	err = o.Close()
	if err != nil {
		return fmt.Errorf("could not close output ROOT file %q: %w", oname, err)
	}
	return nil
}

type copyCmd struct{}

func (cmd copyCmd) process(o *riofs.File, arg string) error {
	log.Printf("copying %q...", arg)

	fname, sel, err := splitArg(arg)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(sel)

	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		if err != nil {
			return err
		}
		name := path[len(f.Name()):]
		if !re.MatchString(name) {
			return nil
		}

		var (
			dst riofs.Directory
			dir = stdpath.Dir(name)
		)

		odst, err := riofs.Dir(o).Get(dir)
		if err != nil {
			v, err := riofs.Dir(o).Mkdir(dir)
			if err != nil {
				return fmt.Errorf("could not create directory %q: %w", dir, err)
			}
			odst = v.(root.Object)
		}
		dst = odst.(riofs.Directory)

		return cmd.copyObj(dst, stdpath.Base(name), obj)
	})
	if err != nil {
		return fmt.Errorf("could not copy input ROOT file: %w", err)
	}
	return nil
}

func (cmd copyCmd) copyObj(odir riofs.Directory, k string, obj root.Object) error {
	var err error
	switch obj := obj.(type) {
	case rtree.Tree:
		err = cmd.copyTree(odir, k, obj)
	case riofs.Directory:
		_, err = odir.Mkdir(k)
	default:
		err = odir.Put(k, obj)
	}

	if err != nil {
		return fmt.Errorf("could not save object %q to output file: %w", k, err)
	}

	return nil
}

func (cmd copyCmd) copyTree(dir riofs.Directory, name string, tree rtree.Tree) error {
	dst, err := rtree.NewWriter(dir, name, rtree.WriteVarsFromTree(tree))
	if err != nil {
		return fmt.Errorf("could not create output copy tree: %w", err)
	}
	_, err = rtree.Copy(dst, tree)
	if err != nil {
		return fmt.Errorf("could not copy tree %q: %w", name, err)
	}

	err = dst.Close()
	if err != nil {
		return fmt.Errorf("could not close copy tree %q: %w", name, err)
	}

	return nil
}
