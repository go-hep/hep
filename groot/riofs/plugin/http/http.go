// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package http is a plugin for riofs.Open to support opening ROOT files over http(s).
package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"go-hep.org/x/hep/groot/riofs"
)

func init() {
	riofs.Register("http", openFile)
	riofs.Register("https", openFile)
}

func openFile(path string) (riofs.Reader, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	f, err := ioutil.TempFile("", "riofs-remote-")
	if err != nil {
		return nil, err
	}
	_, err = io.CopyBuffer(f, resp.Body, make([]byte, 16*1024*1024))
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

}

// tmpFile wraps a regular os.File to automatically remove it when closed.
type tmpFile struct {
	*os.File
}

func (f *tmpFile) Close() error {
	err1 := f.File.Close()
	err2 := os.Remove(f.File.Name())
	if err1 != nil {
		return err1
	}
	return err2
}

var (
	_ riofs.Reader = (*tmpFile)(nil)
	_ riofs.Writer = (*tmpFile)(nil)
)
