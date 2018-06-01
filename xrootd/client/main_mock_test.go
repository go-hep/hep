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
	"go-hep.org/x/hep/xrootd/protocol"
)

var testClientAddrs []string

func testClientWithMockServer(serverFunc func(cancel func(), conn net.Conn), clientFunc func(cancel func(), client *Client)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, conn := net.Pipe()
	defer server.Close()
	defer conn.Close()

	client := &Client{cancel: cancel, conn: conn, mux: mux.New(), signRequirements: protocol.DefaultSignRequirements()}
	defer client.Close()

	go serverFunc(func() { client.Close() }, server)
	go client.consume(ctx)

	clientFunc(cancel, client)
}

func readRequest(conn net.Conn) ([]byte, error) {
	// 16 is for the request options and 4 is for the data length
	const requestSize = protocol.RequestHeaderLength + 16 + 4
	var request = make([]byte, requestSize)
	if _, err := io.ReadFull(conn, request); err != nil {
		return nil, err
	}

	dataLength := binary.BigEndian.Uint32(request[protocol.RequestHeaderLength+16:])
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
		pData, err := protocol.Marshal(p)
		if err != nil {
			return nil, err
		}
		data = append(data, pData...)
	}
	return data, nil
}

// TODO: move unmarshalRequest outside of main_mock_test.go and use it for server implementation.
func unmarshalRequest(data []byte, request interface{}) (protocol.RequestHeader, error) {
	var header protocol.RequestHeader
	if err := protocol.Unmarshal(data[:protocol.RequestHeaderLength], &header); err != nil {
		return protocol.RequestHeader{}, err
	}
	if err := protocol.Unmarshal(data[protocol.RequestHeaderLength:], request); err != nil {
		return protocol.RequestHeader{}, err
	}

	return header, nil
}
