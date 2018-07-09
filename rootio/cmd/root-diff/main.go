// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-diff compares the content of 2 ROOT files, including the content of
// their Trees (for all entries), if any.
//
// Example:
//
//  $> root-diff ./ref.root ./chk.root
//  $> root-diff -k=key1,tree,my-tree ./ref.root ./chk.root
//
//  $> root-diff -h
//  Usage: root-diff [options] a.root b.root
//
//  ex:
//   $> root-diff ./testdata/small-flat-tree.root ./testdata/small-flat-tree.root
//
//  options:
//    -k string
//      	comma-separated list of keys to inspect and compare (default=all common keys)
//
package main // import "go-hep.org/x/hep/rootio/cmd/root-diff"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"go-hep.org/x/hep/rootio"
)

func main() {
	keysFlag := flag.String("k", "", "comma-separated list of keys to inspect and compare (default=all common keys)")

	log.SetPrefix("root-diff: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: root-diff [options] a.root b.root

ex:
 $> root-diff ./testdata/small-flat-tree.root ./testdata/small-flat-tree.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		log.Fatalf("need 2 input ROOT files to compare")
	}

	fref, err := rootio.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer fref.Close()

	fchk, err := rootio.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer fchk.Close()

	keys, err := calcKeys(*keysFlag, fchk, fref)
	if err != nil {
		log.Fatal(err)
	}

	err = diffFiles(keys, fref, fchk)
	if err != nil {
		log.Fatal(err)
	}
}

func calcKeys(kstr string, fchk, fref *rootio.File) ([]string, error) {
	var (
		err   error
		ukeys []string
	)

	if kstr != "" {
		toks := strings.Split(kstr, ",")
		for _, tok := range toks {
			tok = strings.TrimSpace(tok)
			if tok == "" {
				continue
			}
			ukeys = append(ukeys, tok)
		}

		if len(ukeys) == 0 {
			return nil, fmt.Errorf("empty key set")
		}
	} else {
		for _, k := range fchk.Keys() {
			ukeys = append(ukeys, k.Name())
		}
	}

	allgood := true
	var keys []string
	for _, k := range ukeys {
		_, err = fref.Get(k)
		if err != nil {
			allgood = false
			log.Printf("key %q is missing from ref-file=%q", k, fref.Name())
		}

		_, err = fchk.Get(k)
		if err != nil {
			allgood = false
			log.Printf("key %q is missing from chk-file=%q", k, fchk.Name())
		}

		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("empty key set")
	}

	if !allgood {
		return nil, fmt.Errorf("key set differ")
	}

	sort.Strings(keys)
	return keys, err
}

func diffFiles(keys []string, fref, fchk *rootio.File) error {
	for _, key := range keys {
		ref, err := fref.Get(key)
		if err != nil {
			return err
		}

		chk, err := fchk.Get(key)
		if err != nil {
			return err
		}

		err = diffObject(key, ref, chk)
		if err != nil {
			return err
		}
	}

	return nil
}

func diffObject(key string, ref, chk rootio.Object) error {
	refType := reflect.TypeOf(ref)
	chkType := reflect.TypeOf(chk)

	if !reflect.DeepEqual(refType, chkType) {
		return fmt.Errorf("%s: type of keys differ: ref=%v chk=%v", key, refType, chkType)
	}

	switch ref := ref.(type) {
	case rootio.Tree:
		return diffTree(key, ref, chk.(rootio.Tree))
	default:
		return fmt.Errorf("unhandled type %T (key=%v)", ref, key)

	}
}

func diffTree(key string, ref, chk rootio.Tree) error {
	if eref, echk := ref.Entries(), chk.Entries(); eref != echk {
		return fmt.Errorf("%s: number of entries differ: ref=%v chk=%v", key, eref, echk)
	}

	refVars, err := treeVars(ref)
	if err != nil {
		return err
	}

	chkVars, err := treeVars(chk)
	if err != nil {
		return err
	}

	quit := make(chan struct{})
	defer close(quit)

	refc := make(chan treeEntry)
	chkc := make(chan treeEntry)

	go treeDump(quit, refc, ref, refVars)
	go treeDump(quit, chkc, chk, chkVars)

	allgood := true
	n := chk.Entries()
	for i := int64(0); i < n; i++ {
		ref := <-refc
		chk := <-chkc
		if ref.err != nil {
			return fmt.Errorf("%s: error reading ref-tree: %v", key, ref.err)
		}
		if chk.err != nil {
			return fmt.Errorf("%s: error reading chk-tree: %v", key, chk.err)
		}
		if chk.n != ref.n {
			return fmt.Errorf("%s: tree out of sync (ref=%d, chk=%d)", key, ref.n, chk.n)
		}

		for ii := range refVars {
			ref := reflect.Indirect(reflect.ValueOf(refVars[ii].Value)).Interface()
			chk := reflect.Indirect(reflect.ValueOf(chkVars[ii].Value)).Interface()
			diff := cmp.Diff(ref, chk)
			if diff != "" {
				fmt.Printf("key[%s][%04d].%s -- (-ref +chk)\n%s", key, i, refVars[ii].Name, diff)
				allgood = false
				// return fmt.Errorf("%s: trees differ", key)
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

func treeVars(t rootio.Tree) ([]rootio.ScanVar, error) {
	var vars []rootio.ScanVar
	for _, b := range t.Branches() {
		for _, leaf := range b.Leaves() {
			if cls := leaf.Class(); cls == "TLeafElement" {
				return nil, fmt.Errorf("trees with TLeafElement(s) not handled (leaf=%q)", leaf.Name())
			}
			ptr := newValue(leaf)
			if leaf.LeafCount() != nil && false {
				continue
			}
			vars = append(vars, rootio.ScanVar{Name: b.Name(), Leaf: leaf.Name(), Value: ptr})
		}
	}

	return vars, nil
}

func newValue(leaf rootio.Leaf) interface{} {
	etype := leaf.Type()
	switch {
	case leaf.LeafCount() != nil:
		etype = reflect.SliceOf(etype)
	case leaf.Len() > 1 && leaf.Kind() != reflect.String:
		etype = reflect.ArrayOf(leaf.Len(), etype)
	}
	return reflect.New(etype).Interface()
}

type treeEntry struct {
	n   int64
	val interface{}
	err error
	ok  chan int
}

func treeDump(quit chan struct{}, out chan treeEntry, t rootio.Tree, vars []rootio.ScanVar) {
	sc, err := rootio.NewScannerVars(t, vars...)
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
