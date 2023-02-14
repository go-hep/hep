// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/mkdir"
	"go-hep.org/x/hep/xrootd/xrdproto/mv"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/rmdir"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	xrdsync "go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// fshandler implements server.Handler API by making request to the backing filesystem at basePath.
type fshandler struct {
	Handler
	basePath string

	// map + RWMutex works a bit faster and with significant lower memory usage under Linux
	// than sync.Map for given scenarios (write to map once per session and a lot of reads per session).
	mu       sync.RWMutex
	sessions map[[16]byte]*srvSession
}

type srvSession struct {
	mu      sync.Mutex
	handles map[xrdfs.FileHandle]*os.File
}

// NewFSHandler creates a Handler that passes requests to the backing filesystem at basePath.
func NewFSHandler(basePath string) Handler {
	return &fshandler{
		Handler:  Default(),
		basePath: basePath,
		sessions: make(map[[16]byte]*srvSession),
	}
}

// Dirlist implements server.Handler.Dirlist.
func (h *fshandler) Dirlist(sessionID [16]byte, request *dirlist.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	files, err := os.ReadDir(path.Join(h.basePath, request.Path))
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
		info, err := file.Info()
		if err != nil {
			return xrdproto.ServerError{
				Code:    xrdproto.IOError,
				Message: fmt.Sprintf("An IO error occurred: %+v", err),
			}, xrdproto.Error
		}
		entry := xrdfs.EntryStatFrom(info)
		entry.HasStatInfo = resp.WithStatInfo
		resp.Entries = append(resp.Entries, entry)
	}

	return resp, xrdproto.Ok
}

// Open implements server.Handler.Open.
func (h *fshandler) Open(sessionID [16]byte, request *open.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
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
			sess = &srvSession{handles: make(map[xrdfs.FileHandle]*os.File)}
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
		rand.Read(handle[:])
		if _, dup := sess.handles[handle]; !dup {
			resp := open.Response{FileHandle: handle}
			if request.Options&xrdfs.OpenOptionsReturnStatus != 0 {
				st, err := file.Stat()
				if err != nil {
					return xrdproto.ServerError{
						Code:    xrdproto.IOError,
						Message: fmt.Sprintf("An IO error occurred: %v", err),
					}, xrdproto.Error
				}
				es := xrdfs.EntryStatFrom(st)
				resp.Stat = &es
				if request.Options&xrdfs.OpenOptionsCompress == 0 {
					resp.Compression = &xrdfs.FileCompression{}
				}
			}
			// TODO: return compression info if requested.
			sess.handles[handle] = file

			return resp, xrdproto.Ok
		}
	}

	return xrdproto.ServerError{
		Code:    xrdproto.InvalidRequest,
		Message: "handle limit exceeded",
	}, xrdproto.Error
}

// Close implements server.Handler.Close.
func (h *fshandler) Close(sessionID [16]byte, request *xrdclose.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
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
func (h *fshandler) Read(sessionID [16]byte, request *read.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	file := h.getFile(sessionID, request.Handle)
	if file == nil {
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}

	buf := make([]byte, request.Length)
	n, err := file.ReadAt(buf, request.Offset)
	if err != nil && err != io.EOF {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return read.Response{Data: buf[:n]}, xrdproto.Ok
}

// Write implements server.Handler.Write.
func (h *fshandler) Write(sessionID [16]byte, request *write.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	file := h.getFile(sessionID, request.Handle)
	if file == nil {
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}

	_, err := file.WriteAt(request.Data, request.Offset)
	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return nil, xrdproto.Ok
}

func (h *fshandler) getFile(sessionID [16]byte, handle xrdfs.FileHandle) *os.File {
	h.mu.RLock()
	sess, ok := h.sessions[sessionID]
	h.mu.RUnlock()
	if !ok {
		return nil
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	file, ok := sess.handles[handle]
	if !ok {
		return nil
	}
	return file
}

// Stat implements server.Handler.Stat.
func (h *fshandler) Stat(sessionID [16]byte, request *stat.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	if request.Options&stat.OptionsVFS != 0 {
		// TODO: handle virtual stat info.
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: "Stat request with OptionsVFS is not implemented",
		}, xrdproto.Error
	}

	var fi os.FileInfo
	var err error
	if len(request.Path) == 0 {
		file := h.getFile(sessionID, request.FileHandle)
		if file == nil {
			return xrdproto.ServerError{
				Code:    xrdproto.InvalidRequest,
				Message: fmt.Sprintf("Invalid file handle: %v", request.FileHandle),
			}, xrdproto.Error
		}
		fi, err = file.Stat()
	} else {
		fi, err = os.Stat(path.Join(h.basePath, request.Path))
	}

	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return stat.DefaultResponse{EntryStat: xrdfs.EntryStatFrom(fi)}, xrdproto.Ok
}

// Truncate implements server.Handler.Truncate.
func (h *fshandler) Truncate(sessionID [16]byte, request *truncate.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	var err error
	if len(request.Path) == 0 {
		file := h.getFile(sessionID, request.Handle)
		if file == nil {
			return xrdproto.ServerError{
				Code:    xrdproto.InvalidRequest,
				Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
			}, xrdproto.Error
		}
		err = file.Truncate(request.Size)
	} else {
		err = os.Truncate(path.Join(h.basePath, request.Path), request.Size)
	}

	if err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return nil, xrdproto.Ok
}

// Sync implements server.Handler.Sync.
func (h *fshandler) Sync(sessionID [16]byte, request *xrdsync.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	file := h.getFile(sessionID, request.Handle)
	if file == nil {
		return xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Invalid file handle: %v", request.Handle),
		}, xrdproto.Error
	}

	if err := file.Sync(); err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return nil, xrdproto.Ok
}

// Rename implements server.Handler.Rename.
func (h *fshandler) Rename(sessionID [16]byte, request *mv.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	if err := os.Rename(path.Join(h.basePath, request.OldPath), path.Join(h.basePath, request.NewPath)); err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}

	return nil, xrdproto.Ok
}

// Mkdir implements server.Handler.Mkdir.
func (h *fshandler) Mkdir(sessionID [16]byte, request *mkdir.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	mkdirFunc := os.Mkdir
	if request.Options&mkdir.OptionsMakePath != 0 {
		mkdirFunc = os.MkdirAll
	}

	if err := mkdirFunc(path.Join(h.basePath, request.Path), os.FileMode(request.Mode)); err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}
	return nil, xrdproto.Ok
}

// Remove implements server.Handler.Remove.
func (h *fshandler) Remove(sessionID [16]byte, request *rm.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	if err := os.Remove(path.Join(h.basePath, request.Path)); err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}
	return nil, xrdproto.Ok
}

// RemoveDir implements server.Handler.RemoveDir.
func (h *fshandler) RemoveDir(sessionID [16]byte, request *rmdir.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	if err := os.Remove(path.Join(h.basePath, request.Path)); err != nil {
		return xrdproto.ServerError{
			Code:    xrdproto.IOError,
			Message: fmt.Sprintf("An IO error occurred: %v", err),
		}, xrdproto.Error
	}
	return nil, xrdproto.Ok
}

// CloseSession implements server.Handler.CloseSession.
func (h *fshandler) CloseSession(sessionID [16]byte) error {
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
