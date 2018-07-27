// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/xrootd/cmd/xrd-srv"

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"sync"

	"go-hep.org/x/hep/xrootd/server"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// handler implements server.Handler API by making request to the backing filesystem at basePath.
type handler struct {
	server.Handler
	basePath string

	// map + RWMutex works a bit faster and with significant lower memory usage under Linux
	// than sync.Map for given scenarios (write to map once per session and a lot of reads per session).
	mu       sync.RWMutex
	sessions map[[16]byte]*session
}

type session struct {
	mu      sync.Mutex
	handles map[xrdfs.FileHandle]*os.File
}

func newHandler(basePath string) server.Handler {
	return &handler{
		Handler:  server.Default(),
		basePath: basePath,
		sessions: make(map[[16]byte]*session),
	}
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

// Open implements server.Handler.Open.
func (h *handler) Open(sessionID [16]byte, request *open.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	var flag int
	if request.Options&xrdfs.OpenOptionsOpenRead != 0 {
		flag |= os.O_RDONLY
	}
	if request.Options&xrdfs.OpenOptionsOpenUpdate != 0 {
		flag |= os.O_RDWR
	}
	if request.Options&xrdfs.OpenOptionsOpenAppend != 0 {
		flag |= os.O_APPEND
	}
	if request.Options&xrdfs.OpenOptionsNew != 0 || request.Options&xrdfs.OpenOptionsDelete != 0 {
		flag |= os.O_CREATE
		if request.Options&xrdfs.OpenOptionsDelete == 0 {
			flag |= os.O_EXCL
		} else {
			flag |= os.O_TRUNC
		}
	}

	filePath := path.Join(h.basePath, request.Path)
	if request.Options&xrdfs.OpenOptionsMkPath != 0 {
		if err := os.MkdirAll(path.Dir(filePath), os.FileMode(request.Mode)); err != nil {
			return xrdproto.ServerError{
				Code:    xrdproto.IOError,
				Message: fmt.Sprintf("An IO error occurred: %v", err),
			}, xrdproto.Error
		}
	}

	file, err := os.OpenFile(filePath, flag, os.FileMode(request.Mode))
	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	h.mu.RLock()
	sess, ok := h.sessions[sessionID]
	h.mu.RUnlock()
	if !ok {
		h.mu.Lock()
		// Check that there was no change in state during h.mu.RUnlock and h.mu.Lock.
		sess, ok = h.sessions[sessionID]
		if !ok {
			sess = &session{handles: make(map[xrdfs.FileHandle]*os.File)}
			h.sessions[sessionID] = sess
		}
		h.mu.Unlock()
	}

	sess.mu.Lock()
	defer sess.mu.Unlock()
	var handle xrdfs.FileHandle

	// TODO: make handle obtain more deterministic.
	// Right now, we hope that even if 1000000000 of 256*256*256*256 handles are obtained by single user,
	// we have appr. 0.7 probability to find a free handle by the random guess.
	// Then, probability that no free handle is found by 100 tries is something near pow(0.3,100) = 1e-53.
	for i := 0; i < 100; i++ {
		// TODO: use crypto/rand under Windows (4 times faster than math/rand) if handle generation is on the hot path.
		rand.Read(handle[:])
		if _, dup := sess.handles[handle]; !dup {
			sess.handles[handle] = file
			// TODO: return stat info if requested.
			return open.Response{FileHandle: handle}, xrdproto.Ok
		}
	}

	return xrdproto.ServerError{
		Code:    xrdproto.InvalidRequest,
		Message: "handle limit exceeded",
	}, xrdproto.Error
}

// Close implements server.Handler.Close.
func (h *handler) Close(sessionID [16]byte, request *xrdclose.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	h.mu.RLock()
	sess, ok := h.sessions[sessionID]
	h.mu.RUnlock()
	if !ok {
		// This situation can appear if user tries to close without opening any file at all.
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	file, ok := sess.handles[request.Handle]
	if !ok {
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}
	delete(sess.handles, request.Handle)
	err := file.Close()
	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}
	return nil, xrdproto.Ok
}

// Read implements server.Handler.Read.
func (h *handler) Read(sessionID [16]byte, request *read.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	h.mu.RLock()
	sess, ok := h.sessions[sessionID]
	h.mu.RUnlock()
	if !ok {
		// This situation can appear if user tries to read without opening any file at all.
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	file, ok := sess.handles[request.Handle]
	if !ok {
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}

	buf := make([]byte, request.Length)
	_, err := file.ReadAt(buf, request.Offset)
	if err != nil && err != io.EOF {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return read.Response{Data: buf}, xrdproto.Ok
}

// CloseSession implements server.Handler.CloseSession.
func (h *handler) CloseSession(sessionID [16]byte) error {
	h.mu.Lock()
	sess, ok := h.sessions[sessionID]
	if !ok {
		// That means that no files were opened in that session and we have nothing to clear.
		h.mu.Unlock()
		return nil
	}
	delete(h.sessions, sessionID)
	h.mu.Unlock()
	sess.mu.Lock()
	defer sess.mu.Unlock()

	var err error
	for _, f := range sess.handles {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}
