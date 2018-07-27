// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/xrootd/cmd/xrd-srv"

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"reflect"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/server"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
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

func createServer(errorHandler func(err error)) (srv *server.Server, addr, baseDir string, err error) {
	baseDir, err = ioutil.TempDir("", "xrd-srv-")
	if err != nil {
		return nil, "", "", errors.Errorf("xrd-srv: could not create test dir: %v", err)
	}

	addr, err = getTCPAddr()
	if err != nil {
		return nil, "", "", errors.Errorf("xrd-srv: could not get free port to listen: %v", err)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, "", "", errors.Errorf("xrd-srv: could not listen on %q: %v", addr, err)
	}

	srv = server.New(newHandler(baseDir), func(err error) {
		errorHandler(errors.Wrap(err, "xrd-srv: an error occured"))
	})

	go func() {
		if err = srv.Serve(listener); err != nil && err != server.ErrServerClosed {
			errorHandler(errors.Wrap(err, "xrd-srv: could not serve"))
		}
	}()

	return srv, addr, baseDir, nil
}

func createClient(addr string) (*client.Client, error) {
	return client.NewClient(context.Background(), addr, "gopher")
}

func TestHandler_Dirlist(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(baseDir)
	defer srv.Shutdown(context.Background())

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
	defer srv.Shutdown(context.Background())

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
	defer srv.Shutdown(context.Background())

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
				t.Fatalf("could not call Dirlist: %v", err)
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
	defer srv.Shutdown(context.Background())

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
		testName   string
		file       string
		mode       xrdfs.OpenMode
		options    xrdfs.OpenOptions
		createFile bool
		errCode    xrdproto.ServerErrorCode
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
	} {
		t.Run(tc.testName, func(t *testing.T) {
			srv, addr, baseDir, err := createServer(func(err error) {
				t.Error(err)
			})
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(baseDir)
			defer srv.Shutdown(context.Background())

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
			defer srv.Shutdown(context.Background())

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
