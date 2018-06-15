// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package client implements the XRootD client following protocol from http://xrootd.org.

The NewClient function connects to a server:

	ctx := context.Background()

	client, err := NewClient(ctx, addr, username)
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
	"sync/atomic"

	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/krb5"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/unix"
	"go-hep.org/x/hep/xrootd/xrdproto/sigver"
)

// A Client to xrootd server which allows to send requests and receive responses.
// Concurrent requests are supported.
// Zero value is invalid, Client should be instantiated using NewClient.
type Client struct {
	cancel           context.CancelFunc
	conn             net.Conn
	mux              *mux.Mux
	protocolVersion  int32
	signRequirements xrdproto.SignRequirements
	seqID            int64
	auths            map[string]auth.Auther
}

// Option configures an XRootD client.
type Option func(*Client) error

// WithAuth adds an authentication mechanism to the XRootD client.
// If an authentication mechanism was already registered for that provider,
// it will be silently replaced.
func WithAuth(a auth.Auther) Option {
	return func(client *Client) error {
		return client.addAuth(a)
	}
}

func (client *Client) addAuth(auth auth.Auther) error {
	client.auths[auth.Provider()] = auth
	return nil
}

func (client *Client) initSecurityProviders() {
	providers := []auth.Auther{krb5.Default, unix.Default}
	for _, provider := range providers {
		if provider == nil {
			continue
		}
		client.auths[provider.Provider()] = provider
	}
}

// NewClient creates a new xrootd client that connects to the given address using username.
// Options opts configure the client and are applied in the order they were specified.
// When the context expires, a response handling is stopped, however, it is
// necessary to call Cancel to correctly free resources.
func NewClient(ctx context.Context, address string, username string, opts ...Option) (*Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		cancel()
		return nil, err
	}

	client := &Client{cancel: cancel, conn: conn, mux: mux.New(), auths: make(map[string]auth.Auther)}
	client.initSecurityProviders()

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(client); err != nil {
			client.Close()
			return nil, err
		}
	}

	go client.consume(ctx)

	if err := client.handshake(ctx); err != nil {
		client.Close()
		return nil, err
	}

	securityInfo, err := client.Login(ctx, username, "")
	if err != nil {
		client.Close()
		return nil, err
	}

	if len(securityInfo.SecurityInformation) > 0 {
		err = client.auth(ctx, securityInfo.SecurityInformation)
		if err != nil {
			client.Close()
			return nil, err
		}
	}

	protocolInfo, err := client.Protocol(ctx)
	if err != nil {
		client.Close()
		return nil, err
	}

	client.signRequirements = xrdproto.NewSignRequirements(protocolInfo.SecurityLevel, protocolInfo.SecurityOverrides)

	return client, nil
}

// Close closes the connection. Any blocked operation will be unblocked and return error.
func (client *Client) Close() error {
	client.cancel()

	client.mux.Close()
	return client.conn.Close()
}

func (client *Client) consume(ctx context.Context) {
	var header xrdproto.ResponseHeader
	var headerBytes = make([]byte, xrdproto.ResponseHeaderLength)

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

			if err := xrdproto.Unmarshal(headerBytes, &header); err != nil {
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
			} else if header.Status != xrdproto.Ok {
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

			if header.Status != xrdproto.OkSoFar {
				client.mux.Unclaim(header.StreamID)
			}
		}
	}
}

// Send sends the request to the server and stores the response inside the resp.
func (client *Client) Send(ctx context.Context, resp xrdproto.Response, req xrdproto.Request) error {
	data, err := client.call(ctx, req)
	if err != nil {
		return err
	}

	return resp.UnmarshalXrd(xrdenc.NewRBuffer(data))
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

func (client *Client) call(ctx context.Context, req xrdproto.Request) ([]byte, error) {
	streamID, responseChannel, err := client.mux.Claim()
	if err != nil {
		return nil, err
	}

	var wBuffer xrdenc.WBuffer
	header := xrdproto.RequestHeader{streamID, req.ReqID()}
	if err = header.MarshalXrd(&wBuffer); err != nil {
		return nil, err
	}
	if err = req.MarshalXrd(&wBuffer); err != nil {
		return nil, err
	}
	data := wBuffer.Bytes()

	if client.signRequirements.Needed(req) {
		data, err = client.sign(streamID, req.ReqID(), data)
		if err != nil {
			return nil, err
		}
	}

	return client.send(ctx, responseChannel, data)
}

func (client *Client) sign(streamID xrdproto.StreamID, requestID uint16, data []byte) ([]byte, error) {
	seqID := atomic.AddInt64(&client.seqID, 1)
	signRequest := sigver.NewRequest(requestID, seqID, data)
	header := xrdproto.RequestHeader{streamID, signRequest.ReqID()}

	var wBuffer xrdenc.WBuffer
	if err := header.MarshalXrd(&wBuffer); err != nil {
		return nil, err
	}
	if err := signRequest.MarshalXrd(&wBuffer); err != nil {
		return nil, err
	}
	wBuffer.WriteBytes(data)

	return wBuffer.Bytes(), nil
}
