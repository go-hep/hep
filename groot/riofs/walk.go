// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"errors"
	"fmt"
	stdpath "path"
	"strings"

	"go-hep.org/x/hep/groot/root"
)

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = errors.New("riofs: skip this directory") //lint:ignore ST1012 EOF-like sentry

// Walk walks the ROOT file tree rooted at dir, calling walkFn for each ROOT object
// or Directory in the ROOT file tree, including dir.
//
// If an object exists with multiple cycle values, only the latest one is considered.
func Walk(dir Directory, walkFn WalkFunc) error {
	// prepare a "stable" top directory.
	// depending on whether the dir is rooted in a file that was created
	// with an absolute path, the call to Name() may return a path like:
	//   ./data/file.root
	// the first call to walkFn will be given "./data/file.root" as a 'path'
	// argument.
	// but the subsequent calls (walking through directories' hierarchy) will
	// be given "data/file.root/dir11", instead of the probably expected
	// "./data/file.root/dir11".
	//
	// side-step this by providing directly the "stable" top directory in a
	// more regularized form.
	top := stdpath.Join(dir.(root.Named).Name(), ".")
	err := walk(top, dir.(root.Object), walkFn)
	if err == SkipDir {
		return nil
	}
	return err
}

// walk recursively descends path, calling walkFn.
func walk(path string, obj root.Object, walkFn WalkFunc) error {
	dir, ok := obj.(Directory)
	if !ok {
		return walkFn(path, obj, nil)
	}

	err := walkFn(path, obj, nil)
	if err != nil {
		return err
	}

	keys := dir.Keys()
	set := make(map[string]int, len(keys))
	for _, key := range keys {
		if cycle, dup := set[key.Name()]; dup && key.Cycle() < cycle {
			continue
		}
		set[key.Name()] = key.Cycle()
	}

	for _, key := range keys {
		if key.Cycle() != set[key.Name()] {
			continue
		}
		dirname := stdpath.Join(path, key.Name())
		obj, err := dir.Get(key.Name())
		switch err {
		case nil:
			err = walk(dirname, obj, walkFn)
			if err != nil && err != SkipDir {
				return err
			}
		default:
			err := walkFn(dirname, obj, err)
			if err != nil && err != SkipDir {
				return err
			}
		}

	}

	return nil
}

// WalkFunc is the type of the function called for each object or directory
// visited by Walk. The path argument contains the argument to Walk as a
// prefix; that is, if Walk is called with "dir", which is a directory
// containing the file "a", the walk function will be called with argument
// "dir/a". The obj argument is the root.Object for the named path.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how
// to handle that error (and Walk will not descend into that directory). In the
// case of an error, the obj argument will be nil. If an error is returned,
// processing stops. The sole exception is when the function returns the special
// value SkipDir. If the function returns SkipDir when invoked on a directory,
// Walk skips the directory's contents entirely. If the function returns SkipDir
// when invoked on a non-directory root.Object, Walk skips the remaining keys in the
// containing directory.
type WalkFunc func(path string, obj root.Object, err error) error

// recDir handles nested paths.
type recDir struct {
	dir Directory
}

func (dir *recDir) Get(namecycle string) (root.Object, error) { return dir.get(namecycle) }
func (dir *recDir) Put(name string, v root.Object) error      { return dir.put(name, v) }
func (dir *recDir) Keys() []Key                               { return dir.dir.Keys() }
func (dir *recDir) Mkdir(name string) (Directory, error)      { return dir.mkdir(name) }
func (dir *recDir) Parent() Directory                         { return dir.dir.Parent() }

func (dir *recDir) get(namecycle string) (root.Object, error) {
	switch namecycle {
	case "", "/":
		return dir.dir.(root.Object), nil
	}
	name, cycle := decodeNameCycle(namecycle)
	name = strings.TrimPrefix(name, "/")
	path := strings.Split(name, "/")
	return dir.walk(dir.dir, path, cycle)
}

func (dir *recDir) put(name string, v root.Object) error {
	pdir, n := stdpath.Split(name)
	pdir = strings.TrimRight(pdir, "/")
	switch pdir {
	case "":
		return dir.dir.Put(name, v)
	default:
		p, err := dir.mkdir(pdir)
		if err != nil {
			return fmt.Errorf("riofs: could not create parent directory %q for %q: %w", pdir, name, err)
		}
		return p.Put(n, v)
	}
}

func (dir *recDir) mkdir(path string) (Directory, error) {
	if path == "" || path == "/" {
		return nil, fmt.Errorf("riofs: invalid path %q to Mkdir", path)
	}

	if o, err := dir.get(path); err == nil {
		d, ok := o.(Directory)
		if ok {
			return d, nil
		}
		return nil, keyTypeError{key: path, class: d.(root.Object).Class()}
	}

	ps := strings.Split(path, "/")
	if len(ps) == 1 {
		return dir.dir.Mkdir(path)
	}
	for i := range ps {
		p := strings.Join(ps[:i+1], "/")
		_, err := dir.get(p)
		if err == nil {
			continue
		}
		switch {
		case errors.As(err, &noKeyError{}):
			pname, name := stdpath.Split(p)
			pname = strings.TrimRight(pname, "/")
			d, err := dir.get(pname)
			if err != nil {
				return nil, err
			}
			pdir := d.(Directory)
			_, err = pdir.Mkdir(name)
			if err != nil {
				return nil, err
			}
			_, err = pdir.Get(name)
			if err != nil {
				return nil, err
			}
			continue

		default:
			return nil, fmt.Errorf("riofs: unknown error accessing %q: %w", p, err)
		}
	}
	o, err := dir.get(path)
	if err != nil {
		return nil, err
	}
	d, ok := o.(Directory)
	if !ok {
		return nil, fmt.Errorf("riofs: could not create directory %q", path)
	}
	return d, nil
}

func (rd *recDir) walk(dir Directory, path []string, cycle int16) (root.Object, error) {
	if len(path) == 1 {
		name := fmt.Sprintf("%s;%d", path[0], cycle)
		return dir.Get(name)
	}

	o, err := dir.Get(path[0])
	if err != nil {
		return nil, err
	}
	sub, ok := o.(Directory)
	if ok {
		return rd.walk(sub, path[1:], cycle)
	}
	return nil, fmt.Errorf("riofs: not a directory %q", strings.Join([]string{dir.(root.Named).Name(), path[0]}, "/"))
}

// Dir wraps the given directory to handle fully specified directory names:
//
//	rdir := Dir(dir)
//	obj, err := rdir.Get("some/dir/object/name;1")
func Dir(dir Directory) Directory {
	return &recDir{dir}
}

func fileOf(d Directory) *File {
	const max = 1<<31 - 1
	for i := 0; i < max; i++ {
		p := d.Parent()
		if p == nil {
			switch d := d.(type) {
			case *File:
				return d
			case *recDir:
				return fileOf(d.dir)
			default:
				panic(fmt.Errorf("riofs: unknown Directory type %T", d))
			}
		}
		d = p
	}
	panic("impossible")
}

// Get retrieves the named key from the provided directory.
func Get[T any](dir Directory, key string) (T, error) {
	obj, err := Dir(dir).Get(key)
	if err != nil {
		var v T
		return v, err
	}

	v, ok := obj.(T)
	if !ok {
		return v, fmt.Errorf("riofs: could not convert %q (%T) to %T", key, obj, *new(T))
	}

	return v, nil
}

var (
	_ Directory = (*recDir)(nil)
)
