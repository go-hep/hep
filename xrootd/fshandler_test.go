// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd_test // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"reflect"
	"sync"
	"testing"

	"go-hep.org/x/hep/xrootd"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
)

func getTCPAddr() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}

func createServer(errorHandler func(err error)) (srv *xrootd.Server, addr, baseDir string, err error) {
	baseDir, err = ioutil.TempDir("", "xrd-srv-")
	if err != nil {
		return nil, "", "", fmt.Errorf("xrd-srv: could not create test dir: %w", err)
	}

	addr, err = getTCPAddr()
	if err != nil {
		return nil, "", "", fmt.Errorf("xrd-srv: could not get free port to listen: %w", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", "", fmt.Errorf("xrd-srv: could not listen on %q: %w", addr, err)
	}

	srv = xrootd.NewServer(xrootd.NewFSHandler(baseDir), func(err error) {
		errorHandler(fmt.Errorf("xrd-srv: an error occured: %w", err))
	})

	go func() {
		if err = srv.Serve(listener); err != nil && err != xrootd.ErrServerClosed {
			errorHandler(fmt.Errorf("xrd-srv: could not serve: %w", err))
		}
	}()

	return srv, addr, baseDir, nil
}

func createClient(addr string) (*xrootd.Client, error) {
	return xrootd.NewClient(context.Background(), addr, "gopher")
}

func TestHandler_Dirlist(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	file := path.Join(baseDir, "file1.txt")
	err = ioutil.WriteFile(file, nil, 0777)
	if err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Fatalf("could not stat test file: %v", err)
	}

	cli, err := createClient(addr)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer cli.Close()

	got, err := cli.FS().Dirlist(context.Background(), "/")
	if err != nil {
		t.Fatalf("could not call Dirlist: %v", err)
	}

	want := []xrdfs.EntryStat{
		{
			EntryName:   "file1.txt",
			HasStatInfo: true,
			Flags:       xrdfs.StatIsWritable | xrdfs.StatIsReadable,
			Mtime:       fileInfo.ModTime().Unix(),
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("wrong Dirlist response:\ngot = %v\nwant = %v", got, want)
	}
}
func TestHandler_Dirlist_WhenPathIsInvalid(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	cli, err := createClient(addr)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer cli.Close()

	_, err = cli.FS().Dirlist(context.Background(), "/path/not/exist")
	serverError, ok := err.(xrdproto.ServerError)
	if !ok {
		t.Fatalf("could not cast err to ServerError: %v", err)
	}
	if serverError.Code != xrdproto.IOError {
		t.Fatalf("wrong error code:\ngot = %v\nwant = %v", serverError.Code, xrdproto.IOError)
	}
}

func TestHandler_Dirlist_With1000Requests(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	file := path.Join(baseDir, "file1.txt")
	err = ioutil.WriteFile(file, nil, 0777)
	if err != nil {
		t.Fatalf("could not create test file: %v", err)
	}

	cli, err := createClient(addr)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer cli.Close()

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			_, err := cli.FS().Dirlist(context.Background(), "/")
			if err != nil {
				t.Errorf("could not call Dirlist: %v", err)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkHandler_Dirlist(b *testing.B) {
	srv, addr, baseDir, err := createServer(func(err error) {
		b.Error(err)
	})
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	file := path.Join(baseDir, "file1.txt")
	err = ioutil.WriteFile(file, nil, 0777)
	if err != nil {
		b.Fatalf("could not create test file: %v", err)
	}

	cli, err := createClient(addr)
	if err != nil {
		b.Fatalf("could not create client: %v", err)
	}

	defer cli.Close()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err = cli.FS().Dirlist(context.Background(), "/")
		if err != nil {
			b.Fatalf("could not call Dirlist: %v", err)
		}
	}
}

func TestHandler_Open(t *testing.T) {
	for _, tc := range []struct {
		testName      string
		file          string
		mode          xrdfs.OpenMode
		options       xrdfs.OpenOptions
		createFile    bool
		errCode       xrdproto.ServerErrorCode
		checkStatInfo bool
	}{
		{
			testName:   "Readonly | created file",
			options:    xrdfs.OpenOptionsOpenRead,
			createFile: true,
			file:       "test1.txt",
		},
		{
			testName:   "Read & write | created file",
			options:    xrdfs.OpenOptionsOpenUpdate,
			createFile: true,
			file:       "test1.txt",
		},
		{
			testName:   "Append | created file",
			options:    xrdfs.OpenOptionsOpenAppend,
			createFile: true,
			file:       "test1.txt",
		},
		{
			testName: "Read & Write | new file",
			options:  xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsNew,
			file:     "test1.txt",
		},
		{
			testName:   "Read & Write | create existing file",
			options:    xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsNew,
			createFile: true,
			errCode:    xrdproto.IOError,
			file:       "test1.txt",
		},
		{
			testName:   "Read & Write | recreate file",
			options:    xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsDelete,
			createFile: true,
			file:       "test1.txt",
		},
		{
			testName: "Read & Write | new file in new directory without OpenOptionsMkPath",
			options:  xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsNew,
			file:     path.Join("testdir", "test1.txt"),
			errCode:  xrdproto.IOError,
		},
		{
			testName: "Read & Write | new file in new directory with OpenOptionsMkPath",
			options:  xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsNew | xrdfs.OpenOptionsMkPath,
			file:     path.Join("testdir", "test1.txt"),
			mode:     xrdfs.OpenModeOwnerRead | xrdfs.OpenModeOwnerWrite | xrdfs.OpenModeOwnerExecute,
		},
		{
			testName:      "Read & Write | created file | with stat info",
			options:       xrdfs.OpenOptionsOpenUpdate | xrdfs.OpenOptionsReturnStatus,
			file:          "test1.txt",
			createFile:    true,
			checkStatInfo: true,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			if tc.createFile {
				err = ioutil.WriteFile(path.Join(baseDir, tc.file), nil, 0777)
				if err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			got, err := cli.FS().Open(context.Background(), tc.file, tc.mode, tc.options)
			if err != nil {
				if serverError, ok := err.(xrdproto.ServerError); ok {
					if serverError.Code != tc.errCode {
						t.Fatalf("wrong error code:\ngot = %v\nwant = %v\nerror message = %q", serverError.Code, tc.errCode, serverError.Message)
					}
					return
				}
				t.Fatalf("could not call Open: %v", err)
			}
			if err == nil && tc.errCode != 0 {
				t.Fatalf("unexpected successfull call\nwant error code = %v", tc.errCode)
			}

			if tc.checkStatInfo {
				st, err := os.Stat(path.Join(baseDir, tc.file))
				if err != nil {
					t.Fatalf("could not read stat info: %v", err)
				}
				want := xrdfs.EntryStatFrom(st)
				want.EntryName = ""
				if !reflect.DeepEqual(*got.Info(), want) {
					t.Fatalf("wrong stat:\ngot = %v\nwant = %v", *got.Info(), want)
				}
			}

			err = got.Close(context.Background())
			if err != nil {
				t.Fatalf("could not call Close: %v", err)
			}
		})
	}
}

func TestHandler_Read(t *testing.T) {
	bigData := make([]byte, 10*1024)
	_, err := rand.Read(bigData)
	if err != nil {
		t.Fatalf("could not prepare test data: %v", err)
	}

	for _, tc := range []struct {
		testName string
		data     []byte
		want     []byte
		offset   int64
		length   int
	}{
		{
			testName: "Without offset",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			length:   6,
			want:     []byte{1, 2, 3, 4, 5, 6},
		},
		{
			testName: "With offset",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			length:   6,
			offset:   1,
			want:     []byte{2, 3, 4, 5, 6, 7},
		},
		{
			testName: "With offset with EOF",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			length:   20,
			offset:   1,
			want:     []byte{2, 3, 4, 5, 6, 7, 8},
		},
		{
			testName: "With offset larger than file size",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			length:   20,
			offset:   40,
			want:     []byte{},
		},
		{
			testName: "With big length",
			data:     bigData,
			length:   len(bigData),
			offset:   40,
			want:     bigData[40:],
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			file := path.Join(baseDir, "file1.txt")

			err = ioutil.WriteFile(file, tc.data, 0777)
			if err != nil {
				t.Fatalf("could not create test file: %v", err)
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			gotFile, err := cli.FS().Open(context.Background(), "file1.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
			if err != nil {
				t.Fatalf("could not call Open: %v", err)
			}
			defer gotFile.Close(context.Background())

			got := make([]byte, tc.length)
			_, err = gotFile.ReadAt(got, tc.offset)

			if err != nil {
				t.Fatalf("could not call ReadAt: %v", err)
			}

			if !reflect.DeepEqual(got[:len(tc.want)], tc.want) {
				t.Fatalf("wrong data:\ngot = %v\nwant = %v", got[:len(tc.want)], tc.want)
			}
		})
	}
}

func TestHandler_Write(t *testing.T) {
	bigData := make([]byte, 10*1024)
	_, err := rand.Read(bigData)
	if err != nil {
		t.Fatalf("could not prepare test data: %v", err)
	}

	for _, tc := range []struct {
		testName    string
		initialData []byte
		data        []byte
		want        []byte
		n           int
		offset      int64
	}{
		{
			testName: "Without offset",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			want:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			n:        8,
		},
		{
			testName:    "Without offset, with partial rewrite",
			initialData: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data:        []byte{9, 8, 7, 6, 0},
			want:        []byte{9, 8, 7, 6, 0, 6, 7, 8},
			n:           5,
		},
		{
			testName:    "With offset, with partial rewrite",
			initialData: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			data:        []byte{1, 2, 3, 4, 5, 6, 7, 8},
			offset:      1,
			want:        []byte{1, 1, 2, 3, 4, 5, 6, 7, 8},
			n:           8,
		},
		{
			testName: "With offset larger than file size",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			offset:   2,
			want:     []byte{0, 0, 1, 2, 3, 4, 5, 6, 7, 8},
			n:        8,
		},
		{
			testName: "With big length",
			data:     bigData,
			offset:   0,
			want:     bigData,
			n:        len(bigData),
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			file := path.Join(baseDir, "file1.txt")

			err = ioutil.WriteFile(file, tc.initialData, 0777)
			if err != nil {
				t.Fatalf("could not create test file: %v", err)
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			gotFile, err := cli.FS().Open(context.Background(), "file1.txt", xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsOpenUpdate)
			if err != nil {
				t.Fatalf("could not call Open: %v", err)
			}
			defer gotFile.Close(context.Background())

			n, err := gotFile.WriteAt(tc.data, tc.offset)
			if err != nil {
				t.Fatalf("could not call WriteAt: %v", err)
			}

			if n != tc.n {
				t.Fatalf("wrong length:\ngot = %v\nwant = %v", n, tc.n)
			}

			if err := gotFile.Sync(context.Background()); err != nil {
				t.Fatalf("could not call Sync: %v", err)
			}

			got, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatalf("could not read written data: %v", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("wrong data:\ngot = %v\nwant = %v", got, tc.want)
			}
		})
	}
}

func TestHandler_Stat(t *testing.T) {
	for _, tc := range []struct {
		testName string
		isFile   bool
		withPath bool
	}{
		{
			testName: "Stat of file by file handle",
			isFile:   true,
			withPath: false,
		},
		{
			testName: "Stat of file by path",
			isFile:   true,
			withPath: true,
		},
		{
			testName: "Stat of directory by path",
			isFile:   false,
			withPath: true,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			var entry string
			if tc.isFile {
				entry = "file1.txt"
				err := ioutil.WriteFile(path.Join(baseDir, entry), []byte{1, 2, 3, 4, 5}, 0777)
				if err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
			} else {
				entry = "dir1"
				err := os.MkdirAll(path.Join(baseDir, entry), 0777)
				if err != nil {
					t.Fatalf("could not create test directory: %v", err)
				}
			}

			stat, err := os.Stat(path.Join(baseDir, entry))
			if err != nil {
				t.Fatalf("could not read stat info: %v", err)
			}

			want := xrdfs.EntryStatFrom(stat)
			want.EntryName = ""

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			var got xrdfs.EntryStat
			if tc.withPath {
				got, err = cli.FS().Stat(context.Background(), entry)
				if err != nil {
					t.Fatalf("could not call Stat: %v", err)
				}
			} else {
				file, err := cli.FS().Open(context.Background(), entry, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
				if err != nil {
					t.Fatalf("could not call Open: %v", err)
				}
				got, err = file.Stat(context.Background())
				if err != nil {
					t.Fatalf("could not call Stat: %v", err)
				}
				if err := file.Close(context.Background()); err != nil {
					t.Fatalf("could not call Close: %v", err)
				}
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("wrong stat:\ngot = %v\nwant = %v", got, want)
			}
		})
	}
}

func TestHandler_Truncate(t *testing.T) {
	for _, tc := range []struct {
		testName string
		data     []byte
		want     []byte
		size     int64
		withPath bool
	}{
		{
			testName: "Truncate file by file handle",
			data:     []byte{1, 2, 3, 4, 5, 6},
			want:     []byte{1, 2, 3},
			size:     3,
			withPath: false,
		},
		{
			testName: "Truncate file by path",
			data:     []byte{1, 2, 3, 4, 5, 6},
			want:     []byte{1, 2, 3},
			size:     3,
			withPath: true,
		},
		{
			testName: "Extend file by file handle",
			data:     []byte{1, 2, 3, 4, 5, 6},
			want:     []byte{1, 2, 3, 4, 5, 6, 0, 0},
			size:     8,
			withPath: false,
		},
		{
			testName: "Extend file by path",
			data:     []byte{1, 2, 3, 4, 5, 6},
			want:     []byte{1, 2, 3, 4, 5, 6, 0, 0},
			size:     8,
			withPath: true,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			entry := "file1.txt"
			err = ioutil.WriteFile(path.Join(baseDir, entry), tc.data, 0777)
			if err != nil {
				t.Fatalf("could not create test file: %v", err)
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			if tc.withPath {
				err = cli.FS().Truncate(context.Background(), entry, tc.size)
				if err != nil {
					t.Fatalf("could not call Truncate: %v", err)
				}
			} else {
				file, err := cli.FS().Open(context.Background(), entry, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsOpenUpdate)
				if err != nil {
					t.Fatalf("could not call Open: %v", err)
				}
				if err = file.Truncate(context.Background(), tc.size); err != nil {
					t.Fatalf("could not call Truncate: %v", err)
				}
				if err := file.Sync(context.Background()); err != nil {
					t.Fatalf("could not call Sync: %v", err)
				}
				if err := file.Close(context.Background()); err != nil {
					t.Fatalf("could not call Close: %v", err)
				}
			}

			s, err := os.Stat(path.Join(baseDir, entry))
			if err != nil {
				t.Fatalf("could not read stat info: %v", err)
			}

			if !reflect.DeepEqual(s.Size(), tc.size) {
				t.Fatalf("wrong size:\ngot = %v\nwant = %v", s.Size(), tc.size)
			}
		})
	}
}

func TestHandler_Rename(t *testing.T) {
	for _, tc := range []struct {
		testName     string
		oldPathExist bool
		newPathExist bool
	}{
		{
			testName:     "Old path exists, new path doesn't exist",
			oldPathExist: true,
			newPathExist: false,
		},
		{
			testName:     "Old path exists, new path exists",
			oldPathExist: true,
			newPathExist: true,
		},
		{
			testName:     "Old path doesn't exist, new path doesn't exist",
			oldPathExist: true,
			newPathExist: false,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			oldName := "old.txt"
			newName := "new.txt"

			if tc.oldPathExist {
				if err := ioutil.WriteFile(path.Join(baseDir, oldName), nil, 0777); err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
			}
			if tc.newPathExist {
				if err := ioutil.WriteFile(path.Join(baseDir, newName), nil, 0777); err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			if err := cli.FS().Rename(context.Background(), oldName, newName); err != nil {
				t.Fatalf("could not call Rename: %v", err)
			}

			if _, err := os.Stat(path.Join(baseDir, oldName)); !os.IsNotExist(err) {
				t.Fatalf("old file was not removed after rename")
			}
			if _, err := os.Stat(path.Join(baseDir, newName)); os.IsNotExist(err) {
				t.Fatalf("new file was not created after rename")
			}
		})
	}
}

func TestHandler_Mkdir(t *testing.T) {
	for _, tc := range []struct {
		testName   string
		path       string
		createFile bool
		mkdirAll   bool
		errCode    xrdproto.ServerErrorCode
	}{
		{
			testName:   "existing dir",
			createFile: true,
			path:       "testdir",
			errCode:    xrdproto.IOError,
		},
		{
			testName:   "new dir",
			createFile: false,
			path:       "testdir",
		},
		{
			testName:   "nested dir",
			createFile: false,
			path:       "nested/testdir",
			errCode:    xrdproto.IOError,
		},
		{
			testName:   "nested dir and MkdirAll",
			createFile: false,
			path:       "nested/testdir",
			mkdirAll:   true,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			if tc.createFile {
				err = os.MkdirAll(path.Join(baseDir, tc.path), os.FileMode(0777))
				if err != nil {
					t.Fatalf("could not create test dir: %v", err)
				}
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			if tc.mkdirAll {
				err = cli.FS().MkdirAll(context.Background(), tc.path, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite|xrdfs.OpenModeOwnerExecute)
			} else {
				err = cli.FS().Mkdir(context.Background(), tc.path, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite|xrdfs.OpenModeOwnerExecute)
			}
			if err != nil {
				if serverError, ok := err.(xrdproto.ServerError); ok {
					if serverError.Code != tc.errCode {
						t.Fatalf("wrong error code:\ngot = %v\nwant = %v\nerror message = %q", serverError.Code, tc.errCode, serverError.Message)
					}
					return
				}
				t.Fatalf("could not call Mkdir: %v", err)
			}
			if err == nil && tc.errCode != 0 {
				t.Fatalf("unexpected successfull call\nwant error code = %v", tc.errCode)
			}
		})
	}
}

func TestHandler_Remove(t *testing.T) {
	for _, tc := range []struct {
		testName   string
		path       string
		createFile bool
		errCode    xrdproto.ServerErrorCode
	}{
		{
			testName:   "existing file",
			createFile: true,
			path:       "testfile",
		},
		{
			testName:   "non-existing file",
			createFile: false,
			path:       "testfile",
			errCode:    xrdproto.IOError,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			if tc.createFile {
				f, err := os.Create(path.Join(baseDir, tc.path))
				if err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("could not close test file: %v", err)
				}
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			err = cli.FS().RemoveFile(context.Background(), tc.path)
			if err != nil {
				if serverError, ok := err.(xrdproto.ServerError); ok {
					if serverError.Code != tc.errCode {
						t.Fatalf("wrong error code:\ngot = %v\nwant = %v\nerror message = %q", serverError.Code, tc.errCode, serverError.Message)
					}
					return
				}
				t.Fatalf("could not call RemoveFile: %v", err)
			}
			if err == nil && tc.errCode != 0 {
				t.Fatalf("unexpected successfull call\nwant error code = %v", tc.errCode)
			}
		})
	}
}

func TestHandler_RemoveDir(t *testing.T) {
	for _, tc := range []struct {
		testName   string
		path       string
		createFile bool
		createDir  bool
		errCode    xrdproto.ServerErrorCode
	}{
		{
			testName:   "empty existing dir",
			createFile: false,
			createDir:  true,
			path:       "testdir",
		},
		{
			testName:   "non-existing dir",
			createFile: false,
			createDir:  false,
			path:       "testdir",
			errCode:    xrdproto.IOError,
		},
		{
			testName:   "non-empty existing dir",
			createFile: true,
			createDir:  true,
			path:       "testdir",
			errCode:    xrdproto.IOError,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer func() {
				_ = srv.Shutdown(context.Background())
			}()

			dirPath := path.Join(baseDir, tc.path)
			if tc.createDir {
				err := os.MkdirAll(dirPath, os.FileMode(0777))
				if err != nil {
					t.Fatalf("could not create test dir: %v", err)
				}
			}

			if tc.createFile {
				f, err := os.Create(path.Join(dirPath, "file.txt"))
				if err != nil {
					t.Fatalf("could not create test file: %v", err)
				}
				err = f.Close()
				if err != nil {
					t.Fatalf("could not close test file: %v", err)
				}
			}

			cli, err := createClient(addr)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			defer cli.Close()

			err = cli.FS().RemoveDir(context.Background(), tc.path)
			if err != nil {
				if serverError, ok := err.(xrdproto.ServerError); ok {
					if serverError.Code != tc.errCode {
						t.Fatalf("wrong error code:\ngot = %v\nwant = %v\nerror message = %q", serverError.Code, tc.errCode, serverError.Message)
					}
					return
				}
				t.Fatalf("could not call RemoveDir: %v", err)
			}
			if err == nil && tc.errCode != 0 {
				t.Fatalf("unexpected successfull call\nwant error code = %v", tc.errCode)
			}
		})
	}
}

func TestHandler_Ping(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	cli, err := createClient(addr)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer cli.Close()

	_, err = cli.Send(context.Background(), nil, &ping.Request{})
	if err != nil {
		t.Fatalf("could not call Ping: %v", err)
	}
}
