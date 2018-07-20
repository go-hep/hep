// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package server provides a high level API for implementing
// the XRootD server following protocol from http://xrootd.org.
//
// This package contains an implementation of the general requests such
// as handshake, protocol, and login inside the default Handler which
// can be obtained via Default.
package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
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

// New creates a XRootD server which uses specified handler to handle requests
// and errorHandler to handle errors. If the errorHandler is nil,
// then a default error handler is used that does nothing.
func New(handler Handler, errorHandler ErrorHandler) *Server {
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
		s.errorHandler(errors.WithStack(err))
	}
	defer s.handler.CloseSession(sessionID)

	if err := s.handleHandshake(conn); err != nil {
		s.errorHandler(errors.WithStack(err))
		// Abort the connection if the handshake was malformed.
		return
	}

	for {
		// We are using conn for read access only in that place
		// and only once at time for each conn, so no additional
		// serialization is needed.
		reqData, err := ReadRequest(conn)
		if err == io.EOF || err == io.ErrClosedPipe {
			// Client closed the connection.
			return
		}
		if err != nil {
			s.closedMu.RLock()
			defer s.closedMu.RUnlock()
			// TODO: wait for active requests to be processed while closing.
			if !s.closed {
				s.errorHandler(errors.WithStack(err))
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

			if err := WriteResponse(conn, reqHeader.StreamID, status, resp); err != nil {
				s.closedMu.RLock()
				defer s.closedMu.RUnlock()
				// TODO: wait for active requests to be processed while closing.
				if !s.closed {
					s.errorHandler(errors.WithStack(err))
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
		return errors.Errorf("xrootd: connection %v: wrong handshake\ngot = %v\nwant = %v", conn.RemoteAddr(), req, correctHandshake)
	}

	resp, status := s.handler.Handshake()
	return WriteResponse(conn, xrdproto.StreamID{0, 0}, status, resp)
}

func newUnmarshalingErrorResponse(err error) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	response := xrdproto.ServerError{
		Code:    xrdproto.InvalidRequestCode,
		Message: fmt.Sprintf("An error occurred while parsing the request: %v", err),
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
	default:
		response := xrdproto.ServerError{
			Code:    xrdproto.InvalidRequestCode,
			Message: fmt.Sprintf("Unknown request id: %d", requestID),
		}
		return response, xrdproto.Error
	}
}

// ReadRequest reads a XRootD request from r.
// ReadRequest returns entire payload of the request including header.
// ReadRequest requires serialization since multiple ReadFull calls are made.
func ReadRequest(r io.Reader) ([]byte, error) {
	// 16 is for the request options and 4 is for the data length
	const requestSize = xrdproto.RequestHeaderLength + 16 + 4
	request := make([]byte, requestSize)
	if _, err := io.ReadFull(r, request); err != nil {
		return nil, err
	}

	dataLength := binary.BigEndian.Uint32(request[xrdproto.RequestHeaderLength+16:])
	if dataLength == 0 {
		return request, nil
	}

	data := make([]byte, dataLength)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	return append(request, data...), nil
}

// WriteResponse writes a XRootD response resp to the w.
// The response is directed to the stream with id equal to the streamID.
// The status is sent as part of response header.
// WriteResponse writes all data to the w as single Write call, so no
// serialization is required.
func WriteResponse(w io.Writer, streamID xrdproto.StreamID, status xrdproto.ResponseStatus, resp xrdproto.Marshaler) error {
	var respWBuffer xrdenc.WBuffer
	if resp != nil {
		if err := resp.MarshalXrd(&respWBuffer); err != nil {
			return err
		}
	}

	header := xrdproto.ResponseHeader{
		StreamID:   streamID,
		Status:     status,
		DataLength: int32(len(respWBuffer.Bytes())),
	}

	var headerWBuffer xrdenc.WBuffer
	if err := header.MarshalXrd(&headerWBuffer); err != nil {
		return err
	}

	response := append(headerWBuffer.Bytes(), respWBuffer.Bytes()...)
	if _, err := w.Write(response); err != nil {
		return err
	}
	return nil
}
