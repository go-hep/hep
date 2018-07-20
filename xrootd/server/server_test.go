// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"context"
	"io"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
)

type pipeListener struct {
	conns  chan net.Conn
	close  chan struct{}
	closed bool
}

func (pl *pipeListener) Close() error {
	if pl.closed {
		return nil
	}
	pl.closed = true
	pl.close <- struct{}{}
	return nil
}

func (pl *pipeListener) Addr() net.Addr {
	panic("implement me")
}

func (pl *pipeListener) Accept() (net.Conn, error) {
	select {
	case conn := <-pl.conns:
		return conn, nil
	case <-pl.close:
		return nil, closedError{}
	}
}

type closedError struct {
}

func (closedError) Error() string {
	return "xrootd: pipe listener closed"
}

func readResponse(r io.Reader) ([]byte, error) {
	const responseSize = xrdproto.ResponseHeaderLength
	var responseData = make([]byte, responseSize)
	if _, err := io.ReadFull(r, responseData); err != nil {
		return nil, err
	}

	rBuffer := xrdenc.NewRBuffer(responseData)
	var responseHdr xrdproto.ResponseHeader

	if err := responseHdr.UnmarshalXrd(rBuffer); err != nil {
		return nil, err
	}

	if responseHdr.DataLength == 0 {
		return responseData, nil
	}

	var data = make([]byte, responseHdr.DataLength)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	return append(responseData, data...), nil
}

func TestServe_Handshake(t *testing.T) {
	connsCh := make(chan net.Conn, 1)
	p1, p2 := net.Pipe()
	defer p1.Close()
	defer p2.Close()

	connsCh <- p1
	listener := &pipeListener{conns: connsCh, close: make(chan struct{})}
	defer listener.Close()

	server := New(Default(), func(err error) {
		t.Error(err)
	})
	defer server.Shutdown(context.Background())

	go func() {
		req := handshake.NewRequest()
		var wBuffer xrdenc.WBuffer
		req.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())
		resp, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		var (
			respHeader    xrdproto.ResponseHeader
			handshakeResp handshake.Response
			rBuffer       = xrdenc.NewRBuffer(resp)
		)

		if err := respHeader.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal header: %v", err)
		}
		if err := handshakeResp.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal: %v", err)
		}

		wantHeader := xrdproto.ResponseHeader{StreamID: xrdproto.StreamID{0}, Status: xrdproto.Ok, DataLength: 8}
		if !reflect.DeepEqual(wantHeader, respHeader) {
			t.Errorf("wrong response header:\ngot = %v\nwant = %v", respHeader, wantHeader)
		}

		want := handshake.Response{ProtocolVersion: 0x310, ServerType: xrdproto.DataServer}
		if !reflect.DeepEqual(want, handshakeResp) {
			t.Errorf("wrong handshake response:\ngot = %v\nwant = %v", handshakeResp, want)
		}

		server.Shutdown(context.Background())
	}()

	if err := server.Serve(listener); err != nil && err != ErrServerClosed {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServe_Login(t *testing.T) {
	connsCh := make(chan net.Conn, 1)
	p1, p2 := net.Pipe()
	defer p1.Close()
	defer p2.Close()

	connsCh <- p1
	listener := &pipeListener{conns: connsCh, close: make(chan struct{})}
	defer listener.Close()

	server := New(Default(), func(err error) {
		t.Error(err)
	})
	defer server.Shutdown(context.Background())

	go func() {
		handshakeReq := handshake.NewRequest()
		var wBuffer xrdenc.WBuffer
		handshakeReq.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())
		_, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		req := login.NewRequest("gopher", "")
		streamID := [2]byte{0, 1}
		reqHeader := xrdproto.RequestHeader{RequestID: login.RequestID, StreamID: streamID}
		wBuffer = xrdenc.WBuffer{}
		reqHeader.MarshalXrd(&wBuffer)
		req.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())
		resp, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		var (
			respHeader xrdproto.ResponseHeader
			loginResp  login.Response
			rBuffer    = xrdenc.NewRBuffer(resp)
		)

		if err := respHeader.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal header: %v", err)
		}
		if err := loginResp.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal: %v", err)
		}

		wantHeader := xrdproto.ResponseHeader{StreamID: streamID, Status: xrdproto.Ok, DataLength: 16}
		if !reflect.DeepEqual(wantHeader, respHeader) {
			t.Errorf("wrong response header:\ngot = %v\nwant = %v", respHeader, wantHeader)
		}

		// TODO: validate loginResp.

		server.Shutdown(context.Background())
	}()

	if err := server.Serve(listener); err != nil && err != ErrServerClosed {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServe_Protocol(t *testing.T) {
	connsCh := make(chan net.Conn, 1)
	p1, p2 := net.Pipe()
	defer p1.Close()
	defer p2.Close()

	connsCh <- p1
	listener := &pipeListener{conns: connsCh, close: make(chan struct{})}
	defer listener.Close()

	server := New(Default(), func(err error) {
		t.Error(err)
	})
	defer server.Shutdown(context.Background())

	go func() {
		handshakeReq := handshake.NewRequest()
		var wBuffer xrdenc.WBuffer
		handshakeReq.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())
		_, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		req := protocol.NewRequest(0x310, false)
		streamID := [2]byte{0, 2}
		reqHeader := xrdproto.RequestHeader{RequestID: protocol.RequestID, StreamID: streamID}
		wBuffer = xrdenc.WBuffer{}
		reqHeader.MarshalXrd(&wBuffer)
		req.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())

		resp, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		var (
			respHeader   xrdproto.ResponseHeader
			protocolResp protocol.Response
			rBuffer      = xrdenc.NewRBuffer(resp)
		)

		if err := respHeader.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal header: %v", err)
		}
		if err := protocolResp.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal: %v", err)
		}

		wantHeader := xrdproto.ResponseHeader{StreamID: streamID, Status: xrdproto.Ok, DataLength: 8}
		if !reflect.DeepEqual(wantHeader, respHeader) {
			t.Errorf("wrong response header:\ngot = %v\nwant = %v", respHeader, wantHeader)
		}

		want := protocol.Response{BinaryProtocolVersion: 0x310, Flags: protocol.IsServer}
		if !reflect.DeepEqual(want, protocolResp) {
			t.Errorf("wrong response:\ngot = %v\nwant = %v", protocolResp, want)
		}

		server.Shutdown(context.Background())
	}()

	if err := server.Serve(listener); err != nil && err != ErrServerClosed {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServe_Dirlist(t *testing.T) {
	connsCh := make(chan net.Conn, 1)
	p1, p2 := net.Pipe()
	defer p1.Close()
	defer p2.Close()

	connsCh <- p1
	listener := &pipeListener{conns: connsCh, close: make(chan struct{})}
	defer listener.Close()

	server := New(Default(), func(err error) {
		t.Error(err)
	})
	defer server.Shutdown(context.Background())

	go func() {
		handshakeReq := handshake.NewRequest()
		var wBuffer xrdenc.WBuffer
		handshakeReq.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())
		_, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		req := dirlist.NewRequest("/tmp")
		streamID := [2]byte{0, 2}
		reqHeader := xrdproto.RequestHeader{RequestID: dirlist.RequestID, StreamID: streamID}
		wBuffer = xrdenc.WBuffer{}
		reqHeader.MarshalXrd(&wBuffer)
		req.MarshalXrd(&wBuffer)
		p2.Write(wBuffer.Bytes())

		resp, err := readResponse(p2)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		}

		var (
			respHeader xrdproto.ResponseHeader
			errorResp  xrdproto.ServerError
			rBuffer    = xrdenc.NewRBuffer(resp)
		)

		if err := respHeader.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal header: %v", err)
		}
		if err := errorResp.UnmarshalXrd(rBuffer); err != nil {
			t.Errorf("could not unmarshal: %v", err)
		}

		wantHeader := xrdproto.ResponseHeader{StreamID: streamID, Status: xrdproto.Error, DataLength: 39}
		if !reflect.DeepEqual(wantHeader, respHeader) {
			t.Errorf("wrong response header:\ngot = %v\nwant = %v", respHeader, wantHeader)
		}

		want := xrdproto.ServerError{Code: xrdproto.InvalidRequestCode, Message: "Dirlist request is not implemented"}
		if !reflect.DeepEqual(want, errorResp) {
			t.Errorf("wrong response:\ngot = %v\nwant = %v", errorResp, want)
		}

		server.Shutdown(context.Background())
	}()

	if err := server.Serve(listener); err != nil && err != ErrServerClosed {
		t.Fatalf("unexpected error: %v", err)
	}
}
