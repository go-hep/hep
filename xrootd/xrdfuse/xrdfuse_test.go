// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package xrdfuse // import "go-hep.org/x/hep/xrootd/xrdfuse"

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"go-hep.org/x/hep/xrootd"
)

var testClientAddrs []string

func mount(t *testing.T, addr string) (mountPoint string, server *fuse.Server, err error) {
	tmp, err := ioutil.TempDir("", "xrdfuse-")
	if err != nil {
		return "", nil, err
	}

	c, err := xrootd.NewClient(context.Background(), addr, "gopher")
	if err != nil {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remove %q: %v", tmp, err)
		}
		return "", nil, err
	}

	fs := NewFS(c, "/tmp", func(e error) {
		t.Errorf("got error: %v", e)
	})

	nfs := pathfs.NewPathNodeFs(fs, nil)
	server, _, err = nodefs.MountRoot(tmp, nfs.Root(), &nodefs.Options{
		Debug: true,
	})

	return tmp, server, err
}

func testFS_Mkdir(t *testing.T, addr string) {
	mnt, server, err := mount(t, addr)
	if err != nil {
		t.Fatalf("could not mount: %v", err)
	}
	defer func() {
		err := os.RemoveAll(mnt)
		if err != nil {
			t.Logf("could not remove %q: %v", mnt, err)
		}
	}()
	defer server.Unmount()
	go server.Serve()

	tmp, err := ioutil.TempDir(mnt, "xrdfuse-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remove %q: %v", tmp, err)
		}
	}()
}

func TestFS_Mkdir(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFS_Mkdir(t, addr)
		})
	}
}

func testFS_OpenDir(t *testing.T, addr string) {
	mnt, server, err := mount(t, addr)
	if err != nil {
		t.Fatalf("could not mount: %v", err)
	}
	defer func() {
		err := os.RemoveAll(mnt)
		if err != nil {
			t.Logf("could not remove %q: %v", mnt, err)
		}
	}()
	defer server.Unmount()
	go server.Serve()

	tmp, err := ioutil.TempDir(mnt, "xrdfuse-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remove %q: %v", tmp, err)
		}
	}()

	f, err := os.Open(mnt)
	if err != nil {
		t.Fatalf("could not open dir: %v", err)
	}
	defer f.Close()

	tmpName := path.Base(tmp)

	dirs, err := f.Readdirnames(0)
	if err != nil {
		t.Fatalf("could not readdir: %v", err)
	}
	for _, d := range dirs {
		if d == tmpName {
			return
		}
	}
	t.Fatalf("could not find child with name %q", tmpName)
}

func TestFS_OpenDir(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFS_OpenDir(t, addr)
		})
	}
}

func testFS_Rename(t *testing.T, addr string) {
	mnt, server, err := mount(t, addr)
	if err != nil {
		t.Fatalf("could not mount: %v", err)
	}
	defer func() {
		err := os.RemoveAll(mnt)
		if err != nil {
			t.Logf("could not remove %q: %v", mnt, err)
		}
	}()
	defer server.Unmount()
	go server.Serve()

	tmp, err := ioutil.TempDir(mnt, "xrdfuse-")
	if err != nil {
		t.Fatalf("could not create dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remove %q: %v", tmp, err)
		}
	}()

	tmpName := path.Base(tmp)
	newTmpName := tmpName + "-renamed"
	newTmp := path.Join(mnt, newTmpName)

	err = os.Rename(tmp, newTmp)
	if err != nil {
		t.Fatalf("could not rename dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(newTmp)
		if err != nil {
			t.Logf("could not remove %q: %v", newTmp, err)
		}
	}()

	f, err := os.Open(mnt)
	if err != nil {
		t.Fatalf("could not open dir: %v", err)
	}
	defer f.Close()

	dirs, err := f.Readdirnames(0)
	if err != nil {
		t.Fatalf("could not readdir: %v", err)
	}

	ok := false
	for _, d := range dirs {
		if d == tmpName {
			t.Errorf("dir %q was not renamed", tmpName)
		}
		if d == newTmpName {
			ok = true
			break
		}
	}
	if !ok {
		t.Fatalf("could not find dir with name %q", newTmpName)
	}
}

func TestFS_Rename(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFS_Rename(t, addr)
		})
	}
}

func testFS_Mknod(t *testing.T, addr string) {
	mnt, server, err := mount(t, addr)
	if err != nil {
		t.Fatalf("could not mount: %v", err)
	}
	defer func() {
		err := os.RemoveAll(mnt)
		if err != nil {
			t.Logf("could not remove %q: %v", mnt, err)
		}
	}()
	defer server.Unmount()
	go server.Serve()

	tmp, err := ioutil.TempFile(mnt, "xrdfuse-")
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	defer os.Remove(tmp.Name())

	err = tmp.Close()
	if err != nil {
		t.Fatalf("could not close %q: %v", tmp.Name(), err)
	}

	err = os.Remove(tmp.Name())
	if err != nil {
		t.Fatalf("could not remove %q: %v", tmp.Name(), err)
	}
}

func TestFS_Mknod(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFS_Mknod(t, addr)
		})
	}
}

func testFS_Chmod(t *testing.T, addr string) {
	mnt, server, err := mount(t, addr)
	if err != nil {
		t.Fatalf("could not mount: %v", err)
	}
	defer func() {
		err := os.RemoveAll(mnt)
		if err != nil {
			t.Logf("could not remove %q: %v", mnt, err)
		}
	}()
	defer server.Unmount()
	go server.Serve()

	tmp, err := ioutil.TempFile(mnt, "xrdfuse-")
	if err != nil {
		t.Fatalf("could not create file: %v", err)
	}
	defer os.Remove(tmp.Name())

	err = tmp.Close()
	if err != nil {
		t.Fatalf("could not close %q: %v", tmp.Name(), err)
	}

	want := os.FileMode(0222)
	err = os.Chmod(tmp.Name(), want)
	if err != nil {
		t.Fatalf("could not chmod %q: %v", tmp.Name(), err)
	}

	stat, err := os.Stat(tmp.Name())
	if err != nil {
		t.Fatalf("could not stat %q: %v", tmp.Name(), err)
	}

	if stat.Mode() != want {
		t.Fatalf("mode doesn't match:\ngot = %v\nwant = %v\n", stat.Mode(), want)
	}
}

func TestFS_Chmod(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFS_Chmod(t, addr)
		})
	}
}
