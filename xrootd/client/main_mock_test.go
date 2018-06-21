// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/signing"
)

var testClientAddrs []string

func testClientWithMockServer(serverFunc func(cancel func(), conn net.Conn), clientFunc func(cancel func(), client *Client)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, conn := net.Pipe()
	defer server.Close()
	defer conn.Close()

	client := &Client{cancel: cancel, sessions: make(map[string]*session)}
	session := &session{cancel: cancel, conn: conn, mux: mux.New(), requests: make(map[xrdproto.StreamID][]byte), client: client, signRequirements: signing.Default()}
	client.initialSessionID = "test.org:1234"
	client.sessions[client.initialSessionID] = session
	defer client.Close()

	go serverFunc(func() { client.Close() }, server)
	go session.consume(ctx)

	clientFunc(cancel, client)
}

func readRequest(conn net.Conn) ([]byte, error) {
	// 16 is for the request options and 4 is for the data length
	const requestSize = xrdproto.RequestHeaderLength + 16 + 4
	var request = make([]byte, requestSize)
	if _, err := io.ReadFull(conn, request); err != nil {
		return nil, err
	}

	dataLength := binary.BigEndian.Uint32(request[xrdproto.RequestHeaderLength+16:])
	if dataLength == 0 {
		return request, nil
	}

	var data = make([]byte, dataLength)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return append(request, data...), nil
}

func writeResponse(conn net.Conn, data []byte) error {
	n, err := conn.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return errors.Errorf("could not write all %d bytes: wrote %d", len(data), n)
	}
	return nil
}

// TODO: move marshalResponse outside of main_mock_test.go and use it for server implementation.
func marshalResponse(responseParts ...interface{}) ([]byte, error) {
	var data []byte
	for _, p := range responseParts {
		pData, err := xrdproto.Marshal(p)
		if err != nil {
			return nil, err
		}
		data = append(data, pData...)
	}
	return data, nil
}

// TODO: move unmarshalRequest outside of main_mock_test.go and use it for server implementation.
func unmarshalRequest(data []byte, request interface{}) (xrdproto.RequestHeader, error) {
	var header xrdproto.RequestHeader
	if err := xrdproto.Unmarshal(data[:xrdproto.RequestHeaderLength], &header); err != nil {
		return xrdproto.RequestHeader{}, err
	}
	if err := xrdproto.Unmarshal(data[xrdproto.RequestHeaderLength:], request); err != nil {
		return xrdproto.RequestHeader{}, err
	}

	return header, nil
}
