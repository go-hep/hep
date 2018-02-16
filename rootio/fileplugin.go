// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

var (
	_ Reader = (*tmpFile)(nil)
	_ Writer = (*tmpFile)(nil)
)
