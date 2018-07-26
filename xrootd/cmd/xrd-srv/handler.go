// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/xrootd/cmd/xrd-srv"

import (
	"fmt"
	"io/ioutil"
	"path"

	"go-hep.org/x/hep/xrootd/server"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
)

// handler implements server.Handler API by making request to the backing filesystem at basePath.
type handler struct {
	server.Handler
	basePath string
}

func newHandler(basePath string) server.Handler {
	return &handler{Handler: server.Default(), basePath: basePath}
}

// Dirlist implements server.Handler.Dirlist.
func (h *handler) Dirlist(sessionID [16]byte, request *dirlist.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	files, err := ioutil.ReadDir(path.Join(h.basePath, request.Path))
	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	resp := &dirlist.Response{
		WithStatInfo: request.Options&dirlist.WithStatInfo != 0,
		Entries:      make([]xrdfs.EntryStat, 0, len(files)),
	}

	for _, file := range files {
		entry := xrdfs.EntryStatFrom(file)
		entry.HasStatInfo = resp.WithStatInfo
		resp.Entries = append(resp.Entries, entry)
	}

	return resp, xrdproto.Ok
}
