// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	xrdclient "go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

func openFile(path string) (Reader, error) {
	switch {
	case strings.HasPrefix(path, "http://"), strings.HasPrefix(path, "https://"):
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		f, err := ioutil.TempFile("", "rootio-remote-")
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			f.Close()
			return nil, err
		}
		_, err = f.Seek(0, 0)
		if err != nil {
			f.Close()
			return nil, err
		}
		return &tmpFile{f}, nil

	case strings.HasPrefix(path, "file://"):
		return os.Open(path)

	case strings.HasPrefix(path, "xroot://"), strings.HasPrefix(path, "root://"):
		f, err := xrdOpen(path)
		return f, err

	default:
		return os.Open(path)
	}
}

// tmpFile wraps a regular os.File to automatically remove it when closed.
type tmpFile struct {
	*os.File
}

func (f *tmpFile) Close() error {
	os.Remove(f.File.Name())
	return f.File.Close()
}

func xrdOpen(name string) (*xrdFile, error) {
	urn, err := url.Parse(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	host := urn.Hostname()
	port := urn.Port()

	path := urn.Path
	if strings.HasPrefix(path, "//") {
		path = path[1:]
	}

	user := ""
	if urn.User != nil {
		user = urn.User.Username()
	}

	addr := host
	if port != "" {
		addr += ":" + port
	}

	xrd, err := xrdclient.NewClient(context.Background(), addr, user)
	if err != nil {
		return nil, errors.Errorf("rootio: could not connect to xrootd server %q: %v", host, err)
	}

	fs := xrd.FS()
	f, err := fs.Open(context.Background(), path, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		xrd.Close()
		return nil, errors.Errorf("rootio: could not open %q: %v", name, err)
	}

	return &xrdFile{cli: xrd, fs: fs, f: f, name: name}, nil
}

type xrdFile struct {
	cli *xrdclient.Client
	fs  xrdfs.FileSystem
	f   xrdfs.File

	name string
	pos  int64
}

// Name returns the name of the file.
func (f *xrdFile) Name() string {
	return f.name
}

func (f *xrdFile) Close() error {
	err1 := f.f.Close(context.Background())
	err2 := f.cli.Close()
	if err1 != nil {
		return errors.WithStack(err1)
	}
	if err2 != nil {
		return errors.WithStack(err2)
	}
	return nil
}

func (f *xrdFile) Read(data []byte) (int, error) {
	n, err := f.f.ReadAt(data, f.pos)
	f.pos += int64(n)
	return n, err
}

func (f *xrdFile) ReadAt(data []byte, offset int64) (int, error) {
	return f.f.ReadAt(data, offset)
}

func (f *xrdFile) Write(data []byte) (int, error) {
	n, err := f.f.WriteAt(data, f.pos)
	f.pos += int64(n)
	return n, err
}

func (f *xrdFile) WriteAt(data []byte, offset int64) (int, error) {
	return f.f.WriteAt(data, offset)
}

func (f *xrdFile) Seek(offset int64, whence int) (int64, error) {
	var err error
	switch whence {
	case io.SeekStart:
		f.pos = offset
	case io.SeekEnd:
		st, err := f.Stat()
		if err != nil {
			return 0, errors.Errorf("rootio: could not xrootd-stat %q: %v", f.Name(), err)
		}
		f.pos = st.Size() - offset
	case io.SeekCurrent:
		f.pos += offset
	}
	return f.pos, err
}

func (f *xrdFile) Stat() (os.FileInfo, error) {
	v, err := f.f.Stat(context.Background())
	return v, err
}

var (
	_ Reader = (*tmpFile)(nil)
	_ Writer = (*tmpFile)(nil)

	_ Reader = (*xrdFile)(nil)
	_ Writer = (*xrdFile)(nil)
)
