// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd

import (
	"fmt"
	"io"
	"log"
	"os"
	stdpath "path"
	"reflect"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

// Diff compares the values of the list of keys between the two provided ROOT files.
// Diff writes the differing data (if any) to w.
//
// if w is nil, os.Stdout is used.
// if the slice of keys is nil, all keys are considered.
func Diff(w io.Writer, ref, chk *riofs.File, keys []string) error {
	cmd, err := newDiffCmd(w, ref, chk, keys)
	if err != nil {
		return fmt.Errorf("could not compute keys to compare: %w", err)
	}

	return cmd.diffFiles()
}

type diffCmd struct {
	w    io.Writer
	fref *riofs.File
	fchk *riofs.File
	keys []string
}

func newDiffCmd(w io.Writer, fref, fchk *riofs.File, keys []string) (*diffCmd, error) {
	var (
		err   error
		ukeys []string
		cmd   = &diffCmd{fref: fref, fchk: fchk, w: w}
	)

	if w == nil {
		cmd.w = os.Stdout
	}

	if len(keys) != 0 {
		for _, k := range keys {
			k = strings.TrimSpace(k)
			if k == "" {
				continue
			}
			ukeys = append(ukeys, k)
		}

		if len(ukeys) == 0 {
			return nil, fmt.Errorf("empty key set")
		}
	} else {
		for _, k := range cmd.fref.Keys() {
			ukeys = append(ukeys, k.Name())
		}
	}

	allgood := true
	for _, k := range ukeys {
		_, err = cmd.fref.Get(k)
		if err != nil {
			allgood = false
			fmt.Fprintf(cmd.w, "key[%s] -- missing from ref-file\n", k)
			log.Printf("key %q is missing from ref-file=%q", k, cmd.fref.Name())
		}

		_, err = cmd.fchk.Get(k)
		if err != nil {
			allgood = false
			fmt.Fprintf(cmd.w, "key[%s] -- missing from chk-file\n", k)
			log.Printf("key %q is missing from chk-file=%q", k, cmd.fchk.Name())
		}

		cmd.keys = append(cmd.keys, k)
	}

	if len(cmd.keys) == 0 {
		return nil, fmt.Errorf("empty key set")
	}

	if !allgood {
		return nil, fmt.Errorf("key set differ")
	}

	sort.Strings(cmd.keys)
	return cmd, nil
}

func (cmd *diffCmd) diffFiles() error {
	for _, key := range cmd.keys {
		ref, err := cmd.fref.Get(key)
		if err != nil {
			return err
		}

		chk, err := cmd.fchk.Get(key)
		if err != nil {
			return err
		}

		err = cmd.diffObject(key, ref, chk)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *diffCmd) diffObject(key string, ref, chk root.Object) error {
	refType := reflect.TypeOf(ref)
	chkType := reflect.TypeOf(chk)

	if !reflect.DeepEqual(refType, chkType) {
		return fmt.Errorf("%s: type of keys differ: ref=%v chk=%v", key, refType, chkType)
	}

	switch ref := ref.(type) {
	case rtree.Tree:
		return cmd.diffTree(key, ref, chk.(rtree.Tree))
	case riofs.Directory:
		return cmd.diffDir(key, ref, chk.(riofs.Directory))

	case root.Object:
		ok := reflect.DeepEqual(ref, chk)
		if !ok {
			fmt.Fprintf(cmd.w, "key[%s] (%T) -- (-ref +chk)\n-%v\n+%v\n", key, ref, ref, chk)
			return fmt.Errorf("%s: keys differ", key)
		}
		return nil
	default:
		return fmt.Errorf("unhandled type %T (key=%v)", ref, key)
	}
}

func (cmd *diffCmd) diffDir(key string, ref, chk riofs.Directory) error {
	kref := ref.Keys()
	kchk := chk.Keys()
	if len(kref) != len(kchk) {
		return fmt.Errorf("%s: number of keys in directory differ: ref=%d, chk=%d", key, len(kref), len(kchk))
	}

	krefset := make(map[string]struct{})
	kchkset := make(map[string]struct{})
	for _, k := range kref {
		krefset[k.Name()] = struct{}{}
	}
	for _, k := range kchk {
		kchkset[k.Name()] = struct{}{}
	}
	refnames := make([]string, 0, len(krefset))
	for k := range krefset {
		refnames = append(refnames, k)
	}
	chknames := make([]string, 0, len(kchkset))
	for k := range kchkset {
		chknames = append(chknames, k)
	}
	sort.Strings(refnames)
	sort.Strings(chknames)
	if len(krefset) != len(kchkset) {
		return fmt.Errorf("%s: keys in directory differ: ref=%s, chk=%s", key, refnames, chknames)
	}

	for _, k := range refnames {
		oref, err := ref.Get(k)
		if err != nil {
			return fmt.Errorf("%s: could not retrieve %s from ref-directory", key, k)
		}
		ochk, err := chk.Get(k)
		if err != nil {
			return fmt.Errorf("%s: could not retrieve %s from chk-directory", key, k)
		}

		err = cmd.diffObject(stdpath.Join(key, k), oref, ochk)
		if err != nil {
			return fmt.Errorf("%s: values for %s in directory differ: %w", key, k, err)
		}
	}

	return nil
}

func (cmd *diffCmd) diffTree(key string, ref, chk rtree.Tree) error {
	if eref, echk := ref.Entries(), chk.Entries(); eref != echk {
		return fmt.Errorf("%s: number of entries differ: ref=%v chk=%v", key, eref, echk)
	}

	refVars := rtree.NewScanVars(ref)
	chkVars := rtree.NewScanVars(chk)

	quit := make(chan struct{})
	defer close(quit)

	refc := make(chan treeEntry)
	chkc := make(chan treeEntry)

	go cmd.treeDump(quit, refc, ref, refVars)
	go cmd.treeDump(quit, chkc, chk, chkVars)

	allgood := true
	n := chk.Entries()
	for i := int64(0); i < n; i++ {
		ref := <-refc
		chk := <-chkc
		if ref.err != nil {
			return fmt.Errorf("%s: error reading ref-tree: %w", key, ref.err)
		}
		if chk.err != nil {
			return fmt.Errorf("%s: error reading chk-tree: %w", key, chk.err)
		}
		if chk.n != ref.n {
			return fmt.Errorf("%s: tree out of sync (ref=%d, chk=%d)", key, ref.n, chk.n)
		}

		for ii := range refVars {
			var (
				ref  = reflect.Indirect(reflect.ValueOf(refVars[ii].Value)).Interface()
				chk  = reflect.Indirect(reflect.ValueOf(chkVars[ii].Value)).Interface()
				diff = cmp.Diff(ref, chk)
			)
			if diff != "" {
				fmt.Fprintf(cmd.w, "key[%s][%04d].%s -- (-ref +chk)\n%s", key, i, refVars[ii].Name, diff)
				allgood = false
			}
		}
		ref.ok <- 1
		chk.ok <- 1
	}

	if !allgood {
		return fmt.Errorf("%s: trees differ", key)
	}

	return nil
}

type treeEntry struct {
	n   int64
	err error
	ok  chan int
}

func (cmd *diffCmd) treeDump(quit chan struct{}, out chan treeEntry, t rtree.Tree, vars []rtree.ScanVar) {
	sc, err := rtree.NewScannerVars(t, vars...)
	if err != nil {
		out <- treeEntry{err: err}
		return
	}
	defer sc.Close()

	defer close(out)

	next := make(chan int)
	for sc.Next() {
		err = sc.Scan()
		select {
		case <-quit:
			return
		case out <- treeEntry{err: err, n: sc.Entry(), ok: next}:
			<-next
			continue
		}
	}
}
