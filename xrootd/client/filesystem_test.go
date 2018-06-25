// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

func tempdir(client *Client, dir, prefix string) (name string, err error) {
	name, err = ioutil.TempDir("", prefix)
	if err != nil {
		return "", err
	}
	os.RemoveAll(name)

	// Cross-platform way of obtaining the directory name.
	name = filepath.ToSlash(name)
	name = path.Base(name)

	name = path.Join(dir, name)

	fs := client.FS()
	err = fs.MkdirAll(context.Background(), name, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite|xrdfs.OpenModeOwnerExecute)
	if err != nil {
		return "", errors.Errorf("could not create tempdir: %v", err)
	}
	return name, nil
}

func testFileSystem_Dirlist(t *testing.T, addr string) {
	var want = []xrdfs.EntryStat{
		{
			EntryName:   "file1.txt",
			HasStatInfo: true,
			ID:          60129606914,
			EntrySize:   0,
			Mtime:       1529946481,
			Flags:       xrdfs.StatIsReadable | xrdfs.StatIsWritable,
		},
	}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	got, err := fs.Dirlist(context.Background(), "/tmp/dir1")
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FileSystem.Dirlist()\ngot = %v\nwant = %v", got, want)
	}
}

func TestFileSystem_Dirlist(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Dirlist(t, addr)
		})
	}
}

func testFileSystem_Open(t *testing.T, addr string, options xrdfs.OpenOptions, wantFileHandle xrdfs.FileHandle, wantFileCompression *xrdfs.FileCompression, wantFileInfo *xrdfs.EntryStat) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	gotFile, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOtherRead, options)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer gotFile.Close(context.Background())

	if !reflect.DeepEqual(gotFile.Handle(), wantFileHandle) {
		t.Errorf("Filesystem.Open()\ngotFile.Handle() = %v\nwantFileHandle = %v", gotFile.Handle(), wantFileHandle)
	}

	if !reflect.DeepEqual(gotFile.Compression(), wantFileCompression) {
		// TODO: Remove this workaround when fix for https://github.com/xrootd/xrootd/issues/721 will be released.
		skippedDefaultCompressionValue := reflect.DeepEqual(wantFileCompression, &xrdfs.FileCompression{}) && gotFile.Compression() == nil
		if !skippedDefaultCompressionValue {
			t.Errorf("Filesystem.Open()\ngotFile.Compression() = %v\nwantFileCompression = %v", gotFile.Compression(), wantFileCompression)
		}
	}

	if !reflect.DeepEqual(gotFile.Info(), wantFileInfo) {
		t.Errorf("Filesystem.Open()\ngotFile.Info() = %v\nwantFileInfo = %v", gotFile.Info(), wantFileInfo)
	}
}

func TestFileSystem_Open(t *testing.T) {
	emptyCompression := xrdfs.FileCompression{}
	entryStat := &xrdfs.EntryStat{
		HasStatInfo: true,
		ID:          60129606914,
		Mtime:       1529946481,
		Flags:       xrdfs.StatIsWritable | xrdfs.StatIsReadable,
	}

	testCases := []struct {
		name        string
		options     xrdfs.OpenOptions
		handle      xrdfs.FileHandle
		compression *xrdfs.FileCompression
		info        *xrdfs.EntryStat
	}{
		{"WithoutCompressionAndStat", xrdfs.OpenOptionsOpenRead, xrdfs.FileHandle{0, 0, 0, 0}, nil, nil},
		{"WithCompression", xrdfs.OpenOptionsOpenRead | xrdfs.OpenOptionsCompress, xrdfs.FileHandle{0, 0, 0, 0}, &emptyCompression, nil},
		{"WithStat", xrdfs.OpenOptionsOpenRead | xrdfs.OpenOptionsReturnStatus, xrdfs.FileHandle{0, 0, 0, 0}, &emptyCompression, entryStat},
	}

	for _, addr := range testClientAddrs {
		for _, tc := range testCases {
			t.Run(addr+"/"+tc.name, func(t *testing.T) {
				t.Parallel()

				testFileSystem_Open(t, addr, tc.options, tc.handle, tc.compression, tc.info)
			})
		}
	}
}

func testFileSystem_RemoveFile(t *testing.T, addr string) {
	fileName := "rm_test.txt"

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-rm")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(context.Background(), dir)
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsDelete)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	file.Close(context.Background())

	err = fs.RemoveFile(context.Background(), filePath)
	if err != nil {
		t.Fatalf("invalid rm call: %v", err)
	}

	got, err := fs.Dirlist(context.Background(), dir)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}

	found := false
	for _, entry := range got {
		if entry.Name() == fileName {
			found = true
		}
	}

	if found {
		t.Errorf("file '%s' is still present after fs.RemoveFile()", filePath)
	}
}

func TestFileSystem_RemoveFile(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_RemoveFile(t, addr)
		})
	}
}

func testFileSystem_Truncate(t *testing.T, addr string) {
	fileName := "test_truncate_fs.txt"
	write := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	want := write[:4]

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	dir, err := tempdir(client, "/tmp/", "xrd-test-truncate")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(context.Background(), dir)
	filePath := path.Join(dir, fileName)

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	_, err = file.WriteAt(write, 0)
	if err != nil {
		t.Fatalf("invalid write call: %v", err)
	}

	err = file.Sync(context.Background())
	if err != nil {
		t.Fatalf("invalid sync call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}

	err = fs.Truncate(context.Background(), filePath, int64(len(want)))
	if err != nil {
		t.Fatalf("invalid truncate call: %v", err)
	}

	file, err = fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	got := make([]uint8, len(want)+10)
	n, err := file.ReadAt(got, 0)
	if err != nil {
		t.Fatalf("invalid read call: %v", err)
	}

	if n != len(want) {
		t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
	}

	if !reflect.DeepEqual(got[:n], want) {
		t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
	}

}

func TestFileSystem_Truncate(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Truncate(t, addr)
		})
	}
}

func testFileSystem_Stat(t *testing.T, addr string) {
	want := xrdfs.EntryStat{
		HasStatInfo: true,
		ID:          60129606914,
		EntrySize:   0,
		Mtime:       1529946481,
		Flags:       xrdfs.StatIsWritable | xrdfs.StatIsReadable,
	}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	got, err := fs.Stat(context.Background(), "/tmp/dir1/file1.txt")
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filesystem.Open()\ngot = %v\nwant = %v", got, want)
	}
}

func TestFileSystem_Stat(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Stat(t, addr)
		})
	}
}

func testFileSystem_VirtualStat(t *testing.T, addr string) {
	want := xrdfs.VirtualFSStat{
		NumberRW:      1,
		FreeRW:        365,
		UtilizationRW: 23,
	}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	got, err := fs.VirtualStat(context.Background(), "/tmp/dir1/file1.txt")
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if got.NumberRW != want.NumberRW {
		t.Errorf("wrong NumberRW:\ngot = %v\nwant = %v", got.NumberRW, want.NumberRW)
	}

	if got.FreeRW <= 0 || got.FreeRW > 500 {
		t.Errorf("wrong FreeRW:\ngot = %v\nwant to be between 0 and 500", got.FreeRW)
	}

	if got.UtilizationRW <= 0 || got.UtilizationRW > 100 {
		t.Errorf("wrong UtilizationRW:\ngot = %v\nwant to be between 0 and 100", got.UtilizationRW)
	}

	if got.NumberStaging != want.NumberStaging {
		t.Errorf("wrong NumberStaging:\ngot = %v\nwant = %v", got.NumberStaging, want.NumberStaging)
	}
	if got.FreeStaging != want.FreeStaging {
		t.Errorf("wrong FreeStaging:\ngot = %v\nwant = %v", got.FreeStaging, want.FreeStaging)
	}
	if got.UtilizationStaging != want.UtilizationStaging {
		t.Errorf("wrong UtilizationStaging:\ngot = %v\nwant = %v", got.UtilizationStaging, want.UtilizationStaging)
	}
}

func TestFileSystem_VirtualStat(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_VirtualStat(t, addr)
		})
	}
}

func testFileSystem_RemoveDir(t *testing.T, addr string) {
	dirName := "test_remove_dir"

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	parent, err := tempdir(client, "/tmp/", "xrd-test-removedir")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveDir(context.Background(), parent)
	dir := path.Join(parent, dirName)

	err = fs.Mkdir(context.Background(), dir, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite)
	if err != nil {
		t.Fatalf("invalid mkdir call: %v", err)
	}

	dirs, err := fs.Dirlist(context.Background(), parent)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}

	found := false
	for _, d := range dirs {
		if d.EntryName == dirName {
			found = true
		}
	}

	if !found {
		t.Fatalf("dir '%s' has not been created", dir)
	}

	err = fs.RemoveDir(context.Background(), dir)
	if err != nil {
		t.Fatalf("invalid rmdir call: %v", err)
	}

	dirs, err = fs.Dirlist(context.Background(), parent)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}
	for _, d := range dirs {
		if d.EntryName == dirName {
			t.Fatalf("dir '%s' has not been deleted", dir)
		}
	}

}

func TestFileSystem_RemoveDir(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_RemoveDir(t, addr)
		})
	}
}

func TestFileSystem_RemoveAll(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			dirName := "test_remove_all"

			client, err := NewClient(context.Background(), addr, "gopher")
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer client.Close()
			fs := client.FS()

			parent, err := tempdir(client, "/tmp/", "xrd-test-remove-all")
			if err != nil {
				t.Fatal(err)
			}
			defer fs.RemoveAll(context.Background(), parent)
			dir := path.Join(parent, dirName)

			err = fs.Mkdir(context.Background(), dir, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite)
			if err != nil {
				t.Fatalf("invalid mkdir call: %v", err)
			}

			dirs, err := fs.Dirlist(context.Background(), parent)
			if err != nil {
				t.Fatalf("invalid dirlist call: %v", err)
			}

			found := false
			for _, d := range dirs {
				if d.EntryName == dirName {
					found = true
				}
			}

			if !found {
				t.Fatalf("dir '%s' has not been created", dir)
			}

			err = fs.RemoveAll(context.Background(), parent)
			if err != nil {
				t.Fatalf("invalid rmdir call: %v", err)
			}

			dirs, err = fs.Dirlist(context.Background(), "/tmp")
			if err != nil {
				t.Fatalf("invalid dirlist call: %v", err)
			}
			for _, d := range dirs {
				if d.EntryName == path.Base(parent) {
					t.Fatalf("dir '%s' has not been deleted", dir)
				}
			}
		})
	}
}

func testFileSystem_Rename(t *testing.T, addr string) {
	oldName := "test_rename_before"
	newName := "test_rename_after"

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	parent, err := tempdir(client, "/tmp/", "xrd-test-rename")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(context.Background(), parent)
	oldpath := path.Join(parent, oldName)
	newpath := path.Join(parent, newName)

	defer fs.RemoveDir(context.Background(), oldpath)
	defer fs.RemoveDir(context.Background(), newpath)

	err = fs.Mkdir(context.Background(), oldpath, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite)
	if err != nil {
		t.Fatalf("invalid mkdir call: %v", err)
	}

	dirs, err := fs.Dirlist(context.Background(), parent)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}

	found := false
	for _, d := range dirs {
		if d.EntryName == oldName {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("dir %q has not been created", oldpath)
	}

	err = fs.Rename(context.Background(), oldpath, newpath)
	if err != nil {
		t.Fatalf("invalid rmdir call: %v", err)
	}

	dirs, err = fs.Dirlist(context.Background(), parent)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}

	found = false
	for _, d := range dirs {
		if d.EntryName == oldName {
			t.Fatalf("dir %q has not been renamed", oldpath)
		}
		if d.EntryName == newName {
			found = true
		}
	}

	if !found {
		t.Fatalf("dir %q has not been renamed to %q", oldpath, newpath)
	}
}

func TestFileSystem_Rename(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Rename(t, addr)
		})
	}
}

func testFileSystem_Chmod(t *testing.T, addr string) {
	name := "test_chmod"
	oldPerm := xrdfs.OpenModeOwnerWrite | xrdfs.OpenModeOwnerRead
	newPerm := xrdfs.OpenModeOwnerRead

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()
	fs := client.FS()

	parent, err := tempdir(client, "/tmp/", "xrd-test-chmod")
	if err != nil {
		t.Fatal(err)
	}
	defer fs.RemoveAll(context.Background(), parent)
	file := path.Join(parent, name)

	f, err := fs.Open(context.Background(), file, oldPerm, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	err = f.Close(context.Background())
	if err != nil {
		t.Fatalf("could not close file: %v", err)
	}
	defer fs.RemoveFile(context.Background(), file)

	s, err := fs.Stat(context.Background(), file)
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if s.Flags&xrdfs.StatIsWritable == 0 {
		t.Fatalf("invalid mode: file should be writable")
	}

	err = fs.Chmod(context.Background(), file, newPerm)

	s, err = fs.Stat(context.Background(), file)
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if s.Flags&xrdfs.StatIsWritable != 0 {
		t.Fatalf("invalid mode: file should not be writable")
	}

	err = fs.Chmod(context.Background(), file, oldPerm)

	s, err = fs.Stat(context.Background(), file)
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if s.Flags&xrdfs.StatIsWritable == 0 {
		t.Fatalf("invalid mode: file should be writable")
	}
}

func TestFileSystem_Chmod(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Chmod(t, addr)
		})
	}
}

func testFileSystem_Statx(t *testing.T, addr string) {
	want := []xrdfs.StatFlags{xrdfs.StatIsFile, xrdfs.StatIsDir}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	got, err := fs.Statx(context.Background(), []string{"/tmp/dir1/file1.txt", "/tmp/dir1"})
	if err != nil {
		t.Fatalf("invalid statx call: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filesystem.Statx()\ngot = %v\nwant = %v", got, want)
	}
}

func TestFileSystem_Statx(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			t.Parallel()

			testFileSystem_Statx(t, addr)
		})
	}
}
