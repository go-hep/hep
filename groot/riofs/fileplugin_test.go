// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestLocalFile(t *testing.T) {
	local, err := filepath.Abs("../testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	f, err := openFile("file://" + local)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
}

func TestRegister(t *testing.T) {
	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Fatalf("expected a panic")
			}
		}()
		Register("file1", nil)
	}()

	func() {
		defer func() {
			e := recover()
			if e == nil {
				t.Fatalf("expected a panic")
			}
		}()
		Register("test-register", openLocalFile)
		Register("test-register", openLocalFile)
	}()
}

func TestDrivers(t *testing.T) {
	list := Drivers()
	const name = "test-drivers"
	defer func() {
		drivers.Lock()
		defer drivers.Unlock()
		delete(drivers.db, name)
	}()

	Register(name, openLocalFile)
	list = append(list, name)
	sort.Strings(list)

	if got, want := Drivers(), list; !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v, want=%v", got, want)
	}
}
