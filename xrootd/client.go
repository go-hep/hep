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
	"go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/streammanager"
)

var logger = log.New(os.Stderr, "xrootd: ", log.LstdFlags)

// A Client to xrootd server which allows to send requests and receive responses.
// Concurrent requests are supported.
// Zero value is invalid, Client should be instantiated using NewClient.
type Client struct {
	conn            net.Conn
	smgr            *streammanager.StreamManager
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
	StreamID   protocol.StreamID
	Status     uint16
	DataLength int32
}

// NewClient creates a client to xrootd server at address.
// ctx defines a lifetime of the client, once it is done response handling is stopped.
func NewClient(ctx context.Context, address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{conn, streammanager.New(), 0}

	go client.consume(ctx)

	if err := client.handshake(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (client *Client) consume(ctx context.Context) {
	var header = &responseHeader{}
	var headerBytes = make([]byte, responseHeaderSize)

	for {
		if _, err := io.ReadFull(client.conn, headerBytes); err != nil {
			// TODO: handle EOF by redirection as specified at http://xrootd.org/doc/dev45/XRdv310.pdf, page 11
		}

		if err := encoder.Unmarshal(headerBytes, header); err != nil {
			// TODO: should redirect in case if is not possible to decode a header as well?
		}

		resp := &streammanager.ServerResponse{make([]byte, header.DataLength), nil}
		if _, err := io.ReadFull(client.conn, resp.Data); err != nil {
			resp.Error = err
		} else if header.Status != 0 {
			resp.Error = extractError(header, resp.Data)
		}

		if err := client.smgr.SendData(header.StreamID, resp); err != nil {
			// TODO: should we just ignore responses to unclaimed stream IDs?
		}

		if header.Status != protocol.OkSoFar {
			client.smgr.Unclaim(header.StreamID)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func extractError(header *responseHeader, data []byte) error {
	if header.Status == protocol.Error {
		code := int32(binary.BigEndian.Uint32(data[0:4]))
		message := string(data[4 : len(data)-1]) // Skip \0 character at the end

		return serverError{code, message}
	}
	return nil
}

func (client *Client) callWithBytesAndResponseChannel(ctx context.Context, responseChannel streammanager.DataReceiveChannel, requestData []byte) (data []byte, err error) {
	if _, err = client.conn.Write(requestData); err != nil {
		return nil, err
	}

	more := true
	var serverResponse *streammanager.ServerResponse
	for more {
		select {
		case serverResponse, more = <-responseChannel:
			if serverResponse != nil {
				data = append(data, serverResponse.Data...)
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
	streamID, responseChannel, err := client.smgr.Claim()
	if err != nil {
		return nil, err
	}

	requestData, err := encoder.MarshalRequest(requestID, streamID, request)
	if err != nil {
		return nil, err
	}

	return client.callWithBytesAndResponseChannel(ctx, responseChannel, requestData)
}
