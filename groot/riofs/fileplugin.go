// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/go-mmap/mmap"
)

var drivers = struct {
	sync.RWMutex
	db map[string]func(path string) (Reader, error)
}{
	db: make(map[string]func(path string) (Reader, error)),
}

// Register registers a plugin to open ROOT files.
// Register panics if it is called twice with the same name of if the plugin
// function is nil.
func Register(name string, f func(path string) (Reader, error)) {
	drivers.Lock()
	defer drivers.Unlock()
	if f == nil {
		panic("riofs: plugin function is nil")
	}
	if _, dup := drivers.db[name]; dup {
		panic(fmt.Errorf("riofs: Register called twice for plugin %q", name))
	}
	drivers.db[name] = f
}

// Drivers returns a sorted list of the names of the registered plugins
// to open ROOT files.
func Drivers() []string {
	drivers.RLock()
	defer drivers.RUnlock()
	names := make([]string, 0, len(drivers.db))
	for name := range drivers.db {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func openFile(path string) (Reader, error) {
	drivers.RLock()
	defer drivers.RUnlock()

	if f, err := openLocalFile(path); err == nil {
		return f, nil
	}

	scheme := "file"
	if u, err := url.Parse(path); err == nil {
		scheme = u.Scheme
	}
	if open, ok := drivers.db[scheme]; ok {
		return open(path)
	}

	return nil, fmt.Errorf("riofs: no ROOT plugin to open [%s] (scheme=%s)", path, scheme)
}

func openLocalFile(path string) (Reader, error) {
	path = strings.TrimPrefix(path, "file://")
	return os.Open(path)
}

func mmapLocalFile(path string) (Reader, error) {
	path = strings.TrimPrefix(path, "file+mmap://")
	return mmap.Open(path)
}

func init() {
	Register("file", openLocalFile)
	Register("file+mmap", mmapLocalFile)
}
