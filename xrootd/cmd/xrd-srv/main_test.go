// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/xrootd/cmd/xrd-srv"

import (
	"context"
	"io/ioutil"
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
	defer srv.Shutdown(context.Background())
	defer os.RemoveAll(baseDir)

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
	defer srv.Shutdown(context.Background())
	defer os.RemoveAll(baseDir)

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
	if serverError.Code != xrdproto.IOErrorCode {
		t.Fatalf("wrong error code:\ngot = %v\nwant = %v", serverError.Code, xrdproto.IOErrorCode)
	}
}

func TestHandler_Dirlist_With1000Requests(t *testing.T) {
	srv, addr, baseDir, err := createServer(func(err error) {
		t.Error(err)
	})
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Shutdown(context.Background())
	defer os.RemoveAll(baseDir)

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
	defer srv.Shutdown(context.Background())
	defer os.RemoveAll(baseDir)

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
