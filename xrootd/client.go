// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/streammanager"
)

var logger = log.New(os.Stderr, "xrootd: ", log.LstdFlags)

// A Client to xrootd server
type Client struct {
	connection      *net.TCPConn
	sm              *streammanager.StreamManager
	protocolVersion int32
}

type serverError struct {
	Code    int32
	Message string
}

func (err serverError) Error() string {
	return fmt.Sprintf("Server error %d: %s", err.Code, err.Message)
}

const responseHeaderSize = 2 + 2 + 4

type responseHeader struct {
	StreamID   streammanager.StreamID
	Status     uint16
	DataLength int32
}

// New creates a client to xrootd server at address
func New(ctx context.Context, address string) (*Client, error) {
	conn, err := createTCPConnection(address)
	if err != nil {
		return nil, err
	}

	client := &Client{conn, streammanager.New(), 0}

	go client.consume()

	err = client.handshake(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createTCPConnection(address string) (connection *net.TCPConn, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return
	}

	connection, err = net.DialTCP("tcp", nil, tcpAddr)
	return
}

func (client *Client) consume() {
	for {
		var header = &responseHeader{}

		var headerBytes = make([]byte, responseHeaderSize)
		if _, err := io.ReadFull(client.connection, headerBytes); err != nil {
			logger.Panic(err)
		}

		if err := encoder.Unmarshal(headerBytes, header); err != nil {
			logger.Panic(err)
		}

		data := make([]byte, header.DataLength)
		if _, err := io.ReadFull(client.connection, data); err != nil {
			logger.Panic(err)
		}

		response := &streammanager.ServerResponse{data, nil}
		if header.Status != 0 {
			response.Error = extractError(header, data)
		}

		if err := client.sm.SendData(header.StreamID, response); err != nil {
			logger.Panic(err)
		}

		if header.Status != 4000 { // oksofar
			client.sm.Unclaim(header.StreamID)
		}
	}
}

func extractError(header *responseHeader, data []byte) error {
	if header.Status == 4003 {
		code := int32(binary.BigEndian.Uint32(data[0:4]))
		message := string(data[4 : len(data)-1]) // Skip \0 character at the end

		return serverError{code, message}
	}
	return nil
}

func (client *Client) callWithBytesAndResponseChannel(ctx context.Context, responseChannel streammanager.DataReceiveChannel, requestData []byte) (responseBytes []byte, err error) {
	if _, err = client.connection.Write(requestData); err != nil {
		return nil, err
	}

	more := true
	var serverResponse *streammanager.ServerResponse
	for more {
		select {
		case serverResponse, more = <-responseChannel:
			if serverResponse != nil {
				responseBytes = append(responseBytes, serverResponse.Data...)
				err = serverResponse.Error
				if err != nil {
					return
				}
			}
		case <-ctx.Done():
			err = ctx.Err()
		}
	}

	return
}

func (client *Client) call(ctx context.Context, requestID uint16, request interface{}) (responseBytes []byte, err error) {
	streamID, responseChannel, err := client.sm.Claim()
	if err != nil {
		return nil, err
	}

	requestData, err := encoder.MarshalRequest(requestID, streamID, request)
	if err != nil {
		return nil, err
	}

	return client.callWithBytesAndResponseChannel(ctx, responseChannel, requestData)
}
