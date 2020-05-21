// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
)

var fstest = map[string]*xrdfs.EntryStat{
	"/tmp/dir1/file1.txt": {
		HasStatInfo: true,
		ID:          139698106334466,
		EntrySize:   0,
		Mtime:       1530559859,
		Flags:       xrdfs.StatIsReadable,
	},
}

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
		return "", fmt.Errorf("could not create tempdir: %w", err)
	}
	return name, nil
}

func testFileSystem_Dirlist(t *testing.T, addr string) {
	var want = []xrdfs.EntryStat{
		*fstest["/tmp/dir1/file1.txt"],
	}
	want[0].EntryName = "file1.txt"

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
		t.Errorf("FileSystem.Open()\ngotFile.Handle() = %v\nwantFileHandle = %v", gotFile.Handle(), wantFileHandle)
	}

	if !reflect.DeepEqual(gotFile.Compression(), wantFileCompression) {
		// TODO: Remove this workaround when fix for https://github.com/xrootd/xrootd/issues/721 will be released.
		skippedDefaultCompressionValue := reflect.DeepEqual(wantFileCompression, &xrdfs.FileCompression{}) && gotFile.Compression() == nil
		if !skippedDefaultCompressionValue {
			t.Errorf("FileSystem.Open()\ngotFile.Compression() = %v\nwantFileCompression = %v", gotFile.Compression(), wantFileCompression)
		}
	}

	if !reflect.DeepEqual(gotFile.Info(), wantFileInfo) {
		t.Errorf("FileSystem.Open()\ngotFile.Info() = %v\nwantFileInfo = %v", gotFile.Info(), wantFileInfo)
	}
}

func TestFileSystem_Open(t *testing.T) {
	emptyCompression := xrdfs.FileCompression{}
	entryStat := fstest["/tmp/dir1/file1.txt"]

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
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
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
	defer func() {
		_ = fs.RemoveAll(context.Background(), dir)
	}()
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
	want := *fstest["/tmp/dir1/file1.txt"]

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
		t.Errorf("FileSystem.Stat()\ngot = %v\nwant = %v", got, want)
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
	defer func() {
		_ = fs.RemoveDir(context.Background(), parent)
	}()
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
			defer func() {
				_ = fs.RemoveAll(context.Background(), parent)
			}()
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
	oldpath := path.Join(parent, oldName)
	newpath := path.Join(parent, newName)

	defer func() {
		_ = fs.RemoveDir(context.Background(), newpath)
		_ = fs.RemoveDir(context.Background(), oldpath)
		_ = fs.RemoveAll(context.Background(), parent)
	}()

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
	defer func() {
		_ = fs.RemoveAll(context.Background(), parent)
	}()
	file := path.Join(parent, name)

	f, err := fs.Open(context.Background(), file, oldPerm, xrdfs.OpenOptionsNew)
	if err != nil {
		t.Fatalf("could not open file: %v", err)
	}
	err = f.Close(context.Background())
	if err != nil {
		t.Fatalf("could not close file: %v", err)
	}
	defer func() {
		_ = fs.RemoveFile(context.Background(), file)
	}()

	s, err := fs.Stat(context.Background(), file)
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if s.Flags&xrdfs.StatIsWritable == 0 {
		t.Fatalf("invalid mode: file should be writable")
	}

	err = fs.Chmod(context.Background(), file, newPerm)
	if err != nil {
		t.Fatalf("could not chmod %q: %v", file, err)
	}

	s, err = fs.Stat(context.Background(), file)
	if err != nil {
		t.Fatalf("invalid stat call: %v", err)
	}

	if s.Flags&xrdfs.StatIsWritable != 0 {
		t.Fatalf("invalid mode: file should not be writable")
	}

	err = fs.Chmod(context.Background(), file, oldPerm)
	if err != nil {
		t.Fatalf("could not chmod %q: %v", file, err)
	}

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
		t.Errorf("FileSystem.Statx()\ngot = %v\nwant = %v", got, want)
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

func ExampleClient_dirlist() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	entries, err := client.FS().Dirlist(ctx, "/tmp/dir1")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		fmt.Printf("Name: %s, size: %d\n", entry.Name(), entry.Size())
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Name: file1.txt, size: 0
}

func ExampleClient_open() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	file, err := client.FS().Open(ctx, "/tmp/test.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		log.Fatal(err)
	}

	if err := file.Close(ctx); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_removeFile() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().RemoveFile(ctx, "/tmp/test.txt"); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_truncate() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().Truncate(ctx, "/tmp/test.txt", 10); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_stat() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	info, err := client.FS().Stat(ctx, "/tmp/test.txt")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Name: %s, size: %d", info.Name(), info.Size())

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_virtualStat() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	info, err := client.FS().VirtualStat(ctx, "/tmp/")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("RW: %d%% is free", info.FreeRW)

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_mkdir() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().Mkdir(ctx, "/tmp/testdir", xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_mkdirAll() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().MkdirAll(ctx, "/tmp/testdir/subdir", xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_removeDir() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().RemoveDir(ctx, "/tmp/testdir"); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_removeAll() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().RemoveAll(ctx, "/tmp/testdir"); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_rename() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().Rename(ctx, "/tmp/old.txt", "/tmp/new.txt"); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}

func ExampleClient_chmod() {
	ctx := context.Background()
	const username = "gopher"
	client, err := NewClient(ctx, "ccxrootdgotest.in2p3.fr:9001", username)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.FS().Chmod(ctx, "/tmp/test.txt", xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite); err != nil {
		log.Fatal(err)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
