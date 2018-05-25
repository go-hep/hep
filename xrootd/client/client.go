// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package client implements the XRootD client following protocol from http://xrootd.org.

The NewClient function connects to a server:

	ctx := context.Background()

	client, err := NewClient(ctx, *Addr)
	if err != nil {
		// handle error
	}

	// ...

	if err := client.Close(); err != nil {
		// handle error
	}
*/
package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"io"
	"net"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/protocol"
)

// A Client to xrootd server which allows to send requests and receive responses.
// Concurrent requests are supported.
// Zero value is invalid, Client should be instantiated using NewClient.
type Client struct {
	cancel          context.CancelFunc
	conn            net.Conn
	mux             *mux.Mux
	protocolVersion int32
}

// NewClient creates a new xrootd client that connects to the given address.
// When the context expires, a response handling is stopped, however, it is
// necessary to call Cancel to correctly free resources.
func NewClient(ctx context.Context, address string) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		cancel()
		return nil, err
	}

	client := &Client{cancel, conn, mux.New(), 0}

	go client.consume(ctx)

	if err := client.handshake(ctx); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// Close closes the connection. Any blocked operation will be unblocked and return error.
func (client *Client) Close() error {
	client.cancel()

	client.mux.Close()
	return client.conn.Close()
}

func (client *Client) consume(ctx context.Context) {
	var header protocol.ResponseHeader
	var headerBytes = make([]byte, protocol.ResponseHeaderLength)

	for {
		select {
		case <-ctx.Done():
			// TODO: Should wait for active requests to be completed?
			return
		default:
			if _, err := io.ReadFull(client.conn, headerBytes); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				panic(err)
				// TODO: handle EOF by redirection as specified at http://xrootd.org/doc/dev45/XRdv310.pdf, page 11
			}

			if err := encoder.Unmarshal(headerBytes, &header); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				panic(err)
				// TODO: should redirect in case if is not possible to decode a header as well?
			}

			resp := mux.ServerResponse{Data: make([]byte, header.DataLength)}
			if _, err := io.ReadFull(client.conn, resp.Data); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				resp.Err = err
			} else if header.Status != protocol.Ok {
				resp.Err = header.Error(resp.Data)
			}

			if err := client.mux.SendData(header.StreamID, resp); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				panic(err)
				// TODO: should we just ignore responses to unclaimed stream IDs?
			}

			if header.Status != protocol.OkSoFar {
				client.mux.Unclaim(header.StreamID)
			}
		}
	}
}

func (client *Client) send(ctx context.Context, responseChannel mux.DataRecvChan, request []byte) ([]byte, error) {
	if _, err := client.conn.Write(request); err != nil {
		return nil, err
	}

	var data []byte

	for {
		select {
		case resp, more := <-responseChannel:
			if !more {
				return data, nil
			}

			if resp.Err != nil {
				return nil, resp.Err
			}

			data = append(data, resp.Data...)
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
	}
	panic("unreachable")
}

func (client *Client) call(ctx context.Context, requestID uint16, request interface{}) ([]byte, error) {
	streamID, responseChannel, err := client.mux.Claim()
	if err != nil {
		return nil, err
	}

	data, err := encoder.MarshalRequest(requestID, streamID, request)
	if err != nil {
		return nil, err
	}

	return client.send(ctx, responseChannel, data)
}
