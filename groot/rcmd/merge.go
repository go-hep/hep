// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"fmt"
	"log"
	stdpath "path"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

// Merge merges all input fnames ROOT files into the output oname one.
func Merge(oname string, fnames []string, verbose bool) error {
	o, err := groot.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	cmd := mergeCmd{verbose: verbose}
	tsks, err := cmd.mergeTasksFrom(o, fnames[0])
	if err != nil {
		return fmt.Errorf("could not create merge tasks: %w", err)
	}

	for _, fname := range fnames[1:] {
		err := cmd.process(tsks, fname)
		if err != nil {
			return fmt.Errorf("could not process ROOT file %q: %w", fname, err)
		}
	}

	for i := range tsks {
		tsk := &tsks[i]
		err := tsk.close(o)
		if err != nil {
			return fmt.Errorf("could not close task %d (%s): %w", i, tsk.path(), err)
		}
	}

	err = o.Close()
	if err != nil {
		return fmt.Errorf("could not close output ROOT file %q: %w", oname, err)
	}

	return nil
}

type mergeCmd struct {
	verbose bool
}

func (mergeCmd) acceptObj(obj root.Object) bool {
	switch obj.(type) {
	case rtree.Tree:
		// need to specially handle rtree.Tree.
		// rtree.Tree does not implement root.Merger: only rtree.Writer does.
		return true
	case rhist.H1, rhist.H2:
		return true
	case root.Merger:
		return true
	default:
		return false
	}
}

func (cmd mergeCmd) process(tsks []task, fname string) error {
	if cmd.verbose {
		log.Printf("merging [%s]...", fname)
	}

	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	for i := range tsks {
		tsk := &tsks[i]
		err = tsk.merge(f)
		if err != nil {
			return fmt.Errorf("could not merge task %d (%s) for file %q: %w", i, tsk.path(), fname, err)
		}
	}

	return nil
}

type task struct {
	dir string
	key string
	obj root.Object

	verbose bool
}

func (cmd *mergeCmd) mergeTasksFrom(o *riofs.File, fname string) ([]task, error) {
	f, err := groot.Open(fname)
	if err != nil {
		return nil, fmt.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	var tsks []task
	err = riofs.Walk(f, func(path string, obj root.Object, err error) error {
		if err != nil {
			return err
		}
		name := path[len(f.Name()):]
		if name == "" {
			return nil
		}

		if _, ok := obj.(riofs.Directory); ok {
			_, err := riofs.Dir(o).Mkdir(name)
			if err != nil {
				return fmt.Errorf("could not create dir %q in output ROOT file: %w", name, err)
			}
			if cmd.verbose {
				log.Printf("selecting %q", name)
			}
			return nil
		}

		if !cmd.acceptObj(obj) {
			return nil
		}
		if cmd.verbose {
			log.Printf("selecting %q", name)
		}

		var (
			dirName = stdpath.Dir(name)
			objName = stdpath.Base(name)
			dir     = riofs.Directory(o)
		)

		if dirName != "/" && dirName != "" {
			obj, err := riofs.Dir(o).Get(dirName)
			if err != nil {
				return fmt.Errorf("could not get dir %q from output ROOT file: %w", dirName, err)
			}
			dir = obj.(riofs.Directory)
		}

		switch oo := obj.(type) {
		case rtree.Tree:
			w, err := rtree.NewWriter(dir, objName, rtree.WriteVarsFromTree(oo), rtree.WithTitle(oo.Title()))
			if err != nil {
				return fmt.Errorf("could not create output ROOT tree %q: %w", name, err)
			}

			r, err := rtree.NewReader(oo, nil)
			if err != nil {
				return fmt.Errorf(
					"could not create input ROOT tree reader %q: %w",
					name, err,
				)
			}
			defer r.Close()

			_, err = rtree.Copy(w, r)
			if err != nil {
				return fmt.Errorf("could not seed output ROOT tree %q: %w", name, err)
			}
			obj = w
		}

		tsks = append(tsks, task{
			dir:     dirName,
			key:     objName,
			obj:     obj,
			verbose: cmd.verbose,
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not inspect input ROOT file: %w", err)
	}

	if cmd.verbose {
		log.Printf("merging [%s]...", fname)
	}

	return tsks, nil
}

func (tsk *task) path() string {
	return stdpath.Join(tsk.dir, tsk.key)
}

func (tsk *task) merge(f *riofs.File) error {
	name := tsk.path()
	obj, err := riofs.Dir(f).Get(name)
	if err != nil {
		return fmt.Errorf("could not get %q: %w", name, err)
	}

	err = tsk.mergeObj(tsk.obj, obj)
	if err != nil {
		return fmt.Errorf("could not merge %q: %w", name, err)
	}

	return nil
}

func (tsk *task) close(f *riofs.File) error {
	var err error
	switch obj := tsk.obj.(type) {
	case rtree.Writer:
		err = obj.Close()
	default:
		err = riofs.Dir(f).Put(tsk.path(), tsk.obj)
	}

	if err != nil {
		return fmt.Errorf("could not save %q (%T) to output ROOT file: %w", tsk.path(), tsk.obj, err)
	}

	return nil
}

func (tsk *task) mergeObj(dst, src root.Object) error {
	var (
		rdst = dst.Class()
		rsrc = src.Class()
	)
	if rdst != rsrc {
		return fmt.Errorf("types differ: dst=%T, src=%T", dst, src)
	}

	switch dst := dst.(type) {
	case rhist.H2:
		return tsk.mergeH2(dst, src.(rhist.H2))
	case root.Merger:
		return dst.ROOTMerge(src)
	default:
		return fmt.Errorf("could not find suitable merge-API for (dst=%T, src=%T)", dst, src)
	}
}

func (tsk *task) mergeH2(dst, src rhist.H2) error {
	panic("not implemented")
}
