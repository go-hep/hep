// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/mkdir"
	"go-hep.org/x/hep/xrootd/xrdproto/mv"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/rmdir"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	xrdsync "go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// ErrServerClosed is returned by the Server's Serve method after a call to Shutdown.
var ErrServerClosed = errors.New("xrootd: server closed")

// ErrorHandler is the function which handles occurred error (e.g. logs it).
type ErrorHandler func(error)

// Server implements the XRootD server following protocol from http://xrootd.org.
// The Server uses a Handler to handle incoming requests.
// To listen for incoming connections, Serve method must be called.
// It is possible to configure to listen on several ports simultaneously
// by calling Serve with different net.Listeners.
type Server struct {
	handler      Handler
	errorHandler ErrorHandler

	mu        sync.Mutex
	listeners []net.Listener

	closedMu sync.RWMutex
	closed   bool

	connMu     sync.Mutex
	activeConn map[net.Conn]struct{}
}

// NewServer creates a XRootD server which uses specified handler to handle requests
// and errorHandler to handle errors. If the errorHandler is nil,
// then a default error handler is used that does nothing.
func NewServer(handler Handler, errorHandler ErrorHandler) *Server {
	if errorHandler == nil {
		errorHandler = func(error) {}
	}
	return &Server{
		handler:      handler,
		errorHandler: errorHandler,
		activeConn:   make(map[net.Conn]struct{}),
	}
}

// Shutdown stops Server and closes all listeners and active connections.
// Shutdown returns the first non nil error while closing listeners and connections.
func (s *Server) Shutdown(ctx context.Context) error {
	var err error

	s.closedMu.Lock()
	s.closed = true
	s.closedMu.Unlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.listeners {
		if cerr := s.listeners[i].Close(); cerr != nil && err == nil {
			err = cerr
		}
	}

	// TODO: wait for active requests to be processed as long as ctx is not done.
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for conn := range s.activeConn {
		if cerr := conn.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each. The service goroutines read requests and
// then call s.handler to handle them.
func (s *Server) Serve(l net.Listener) error {
	s.mu.Lock()
	s.listeners = append(s.listeners, l)
	s.mu.Unlock()
	for {
		conn, err := l.Accept()
		if err != nil {
			s.closedMu.RLock()
			defer s.closedMu.RUnlock()
			if s.closed {
				return ErrServerClosed
			}
			return err
		}

		s.connMu.Lock()
		s.activeConn[conn] = struct{}{}
		s.connMu.Unlock()

		go s.handleConnection(conn)
	}
}

// handleConnection handles the client connection.
// handleConnection reads the handshake and checks it correctness.
// In case of success, main loop is started that reads requests and
// handles them. Otherwise, connection is aborted.
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		s.connMu.Lock()
		delete(s.activeConn, conn)
		s.connMu.Unlock()
	}()

	var sessionID [16]byte
	if _, err := rand.Read(sessionID[:]); err != nil {
		s.errorHandler(fmt.Errorf("could not read session ID: %w", err))
	}
	defer func() {
		if err := s.handler.CloseSession(sessionID); err != nil {
			s.errorHandler(fmt.Errorf("could not close session ID %q: %w", sessionID, err))
		}
	}()

	if err := s.handleHandshake(conn); err != nil {
		s.errorHandler(fmt.Errorf("could not handle handshake: %w", err))
		// Abort the connection if the handshake was malformed.
		return
	}

	for {
		// We are using conn for read access only in that place
		// and only once at time for each conn, so no additional
		// serialization is needed.
		reqData, err := xrdproto.ReadRequest(conn)
		if err == io.EOF || err == io.ErrClosedPipe {
			// Client closed the connection.
			return
		}
		if err != nil {
			s.closedMu.RLock()
			defer s.closedMu.RUnlock()
			// TODO: wait for active requests to be processed while closing.
			if !s.closed {
				s.errorHandler(fmt.Errorf("could not close connection: %w", err))
			}
			// Abort the connection if an error occurred during
			// the reading phase because we can't recover from it.
			return
		}

		// Performing a request may take some time so we are running it
		// in the separate goroutine. We follow the XRootD protocol and
		// write results back with StreamID provided in the request,
		// so Client will match the responses to the corresponding request calls.
		go func(req []byte) {
			var (
				reqHeader xrdproto.RequestHeader
				resp      xrdproto.Marshaler
				status    xrdproto.ResponseStatus
			)

			rBuffer := xrdenc.NewRBuffer(req)
			if err := reqHeader.UnmarshalXrd(rBuffer); err != nil {
				resp, status = newUnmarshalingErrorResponse(err)
			} else {
				resp, status = s.handleRequest(sessionID, reqHeader.RequestID, rBuffer)
			}

			if err := xrdproto.WriteResponse(conn, reqHeader.StreamID, status, resp); err != nil {
				s.closedMu.RLock()
				defer s.closedMu.RUnlock()
				// TODO: wait for active requests to be processed while closing.
				if !s.closed {
					s.errorHandler(fmt.Errorf("could not close connection: %w", err))
				}
				// Abort the connection if an error occurred during
				// the writing phase because we can't recover from it.
				return
			}
		}(reqData)
	}
}

func (s *Server) handleHandshake(conn net.Conn) error {
	data := make([]byte, handshake.RequestLength)
	if _, err := io.ReadFull(conn, data); err != nil {
		return err
	}

	var req handshake.Request
	rBuffer := xrdenc.NewRBuffer(data)
	err := req.UnmarshalXrd(rBuffer)
	if err != nil {
		return err
	}

	correctHandshake := handshake.NewRequest()
	if !reflect.DeepEqual(req, correctHandshake) {
		return fmt.Errorf("xrootd: connection %v: wrong handshake\ngot = %v\nwant = %v", conn.RemoteAddr(), req, correctHandshake)
	}

	resp, status := s.handler.Handshake()
	return xrdproto.WriteResponse(conn, xrdproto.StreamID{0, 0}, status, resp)
}

func newUnmarshalingErrorResponse(err error) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	response := xrdproto.ServerError{
		Code:    xrdproto.InvalidRequest,
		Message: fmt.Errorf("An error occurred while parsing the request: %w", err).Error(),
	}
	return response, xrdproto.Error
}

func (s *Server) handleRequest(sessionID [16]byte, requestID uint16, rBuffer *xrdenc.RBuffer) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	switch requestID {
	case login.RequestID:
		var request login.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Login(sessionID, &request)
	case protocol.RequestID:
		var request protocol.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Protocol(sessionID, &request)
	case dirlist.RequestID:
		var request dirlist.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Dirlist(sessionID, &request)
	case open.RequestID:
		var request open.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Open(sessionID, &request)
	case xrdclose.RequestID:
		var request xrdclose.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Close(sessionID, &request)
	case read.RequestID:
		var request read.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Read(sessionID, &request)
	case write.RequestID:
		var request write.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Write(sessionID, &request)
	case stat.RequestID:
		var request stat.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Stat(sessionID, &request)
	case xrdsync.RequestID:
		var request xrdsync.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Sync(sessionID, &request)
	case truncate.RequestID:
		var request truncate.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Truncate(sessionID, &request)
	case mv.RequestID:
		var request mv.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Rename(sessionID, &request)
	case mkdir.RequestID:
		var request mkdir.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Mkdir(sessionID, &request)
	case ping.RequestID:
		var request ping.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Ping(sessionID, &request)
	case rm.RequestID:
		var request rm.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.Remove(sessionID, &request)
	case rmdir.RequestID:
		var request rmdir.Request
		err := request.UnmarshalXrd(rBuffer)
		if err != nil {
			return newUnmarshalingErrorResponse(err)
		}
		return s.handler.RemoveDir(sessionID, &request)
	default:
		response := xrdproto.ServerError{
			Code:    xrdproto.InvalidRequest,
			Message: fmt.Sprintf("Unknown request id: %d", requestID),
		}
		return response, xrdproto.Error
	}
}
