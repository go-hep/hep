// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsrv

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"go-hep.org/x/hep/groot/riofs"
)

type DB struct {
	sync.RWMutex
	dir   string
	files map[string]*riofs.File // a map of URI -> ROOT file
}

func NewDB(dir string) *DB {
	os.MkdirAll(dir, 0755)
	return &DB{
		dir:   dir,
		files: make(map[string]*riofs.File),
	}
}

func (db *DB) Close() {
	db.Lock()
	defer db.Unlock()
	for _, f := range db.files {
		f.Close()
	}
	db.files = nil
	os.RemoveAll(db.dir)
}

func (db *DB) Files() []string {
	db.RLock()
	defer db.RUnlock()
	uris := make([]string, 0, len(db.files))
	for uri := range db.files {
		uris = append(uris, uri)
	}
	sort.Strings(uris)
	return uris
}

func (db *DB) Tx(uri string, fct func(f *riofs.File) error) error {
	db.RLock()
	defer db.RUnlock()
	f := db.files[uri]
	if f == nil {
		return fmt.Errorf("rsrv: no such file %q", uri)
	}
	return fct(db.files[uri])
}

func (db *DB) get(uri string) *riofs.File {
	db.RLock()
	defer db.RUnlock()
	return db.files[uri]
}

func (db *DB) set(uri string, f *riofs.File) {
	db.Lock()
	defer db.Unlock()
	if old, dup := db.files[uri]; dup {
		old.Close()
	}
	db.files[uri] = f
}

func (db *DB) del(uri string) {
	db.Lock()
	defer db.Unlock()

	f, ok := db.files[uri]
	if !ok {
		return
	}
	f.Close()
	delete(db.files, uri)
}
