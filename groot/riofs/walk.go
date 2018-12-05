// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	stdpath "path"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/root"
)

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = errors.New("riofs: skip this directory")

// Walk walks the ROOT file tree rooted at dir, calling walkFn for each ROOT object
// or Directory in the ROOT file tree, including dir.
func Walk(dir Directory, walkFn WalkFunc) error {
	err := walk(dir.(root.Named).Name(), dir.(root.Object), walkFn)
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

	keys := dir.Keys()
	err := walkFn(path, obj, nil)
	if err != nil {
		return err
	}

	for _, key := range keys {
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
func (dir *recDir) Put(name string, v root.Object) error      { return dir.dir.Put(name, v) }
func (dir *recDir) Keys() []Key                               { return dir.dir.Keys() }
func (dir *recDir) Mkdir(name string) (Directory, error)      { return dir.dir.Mkdir(name) }

func (dir *recDir) get(namecycle string) (root.Object, error) {
	switch namecycle {
	case "", "/":
		return dir.dir.(root.Object), nil
	}
	name, cycle := decodeNameCycle(namecycle)
	if strings.HasPrefix(name, "/") {
		name = name[1:]
	}
	path := strings.Split(name, "/")
	return dir.walk(dir.dir, path, cycle)
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
	return nil, errors.Errorf("riofs: not a directory %q", strings.Join([]string{dir.(root.Named).Name(), path[0]}, "/"))
}

// Dir wraps the given directory to handle fully specified directory names:
//  rdir := Dir(dir)
//  obj, err := rdir.Get("some/dir/object/name;1")
func Dir(dir Directory) Directory {
	return &recDir{dir}
}

var (
	_ Directory = (*recDir)(nil)
)
