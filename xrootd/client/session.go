// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/signing"
	"go-hep.org/x/hep/xrootd/xrdproto/sigver"
)

// session is a connection to the specific XRootD server
// which allows to send requests and receive responses.
// Concurrent requests are supported.
// Zero value is invalid, session should be instantiated using newSession.
//
// The session is used by the Client to send requests to the particular server
// specified by the name and port. If the current server cannot
// handle a request, it responds with the redirect to the new server.
// After that, Client obtains a session associated with that server and
// re-issues the request. Stream ID may be different during these 2 requests
// because it is used to identify requests among one particular server
// and is not shared between servers in any way.
type session struct {
	cancel           context.CancelFunc
	conn             net.Conn
	mux              *mux.Mux
	protocolVersion  int32
	signRequirements signing.Requirements
	seqID            int64
	mu               sync.RWMutex
	requests         map[xrdproto.StreamID][]byte

	client    *Client
	sessionID string
}

func newSession(ctx context.Context, address, username, token string, client *Client) (*session, error) {
	ctx, cancel := context.WithCancel(ctx)

	var d net.Dialer
	addr := parseAddr(address)
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		cancel()
		return nil, err
	}

	sess := &session{
		cancel:    cancel,
		conn:      conn,
		mux:       mux.New(),
		requests:  make(map[xrdproto.StreamID][]byte),
		client:    client,
		sessionID: addr,
	}

	go sess.consume(ctx)

	if err := sess.handshake(ctx); err != nil {
		sess.Close()
		return nil, err
	}

	securityInfo, err := sess.Login(ctx, username, token)
	if err != nil {
		sess.Close()
		return nil, err
	}

	if len(securityInfo.SecurityInformation) > 0 {
		err = sess.auth(ctx, securityInfo.SecurityInformation)
		if err != nil {
			sess.Close()
			return nil, err
		}
	}

	protocolInfo, err := sess.Protocol(ctx)
	if err != nil {
		sess.Close()
		return nil, err
	}

	sess.signRequirements = signing.New(protocolInfo.SecurityLevel, protocolInfo.SecurityOverrides)

	return sess, nil
}

// Close closes the connection. Any blocked operation will be unblocked and return error.
func (sess *session) Close() error {
	sess.cancel()

	sess.mux.Close()
	// TODO: should we remove session here somehow?
	return sess.conn.Close()
}

func (sess *session) consume(ctx context.Context) {
	var header xrdproto.ResponseHeader
	var headerBytes = make([]byte, xrdproto.ResponseHeaderLength)

	for {
		select {
		case <-ctx.Done():
			// TODO: Should wait for active requests to be completed?
			return
		default:
			if _, err := io.ReadFull(sess.conn, headerBytes); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				if sess.sessionID == sess.client.initialSessionID {
					// TODO: what should we do in case initial session is aborted?
					// Should we try to reconnect to the server and re-issue all requests?
					panic(err)
				}
				sess.mu.RLock()
				resp := mux.ServerResponse{Redirection: &mux.Redirection{Addr: sess.client.initialSessionID}}
				for streamID := range sess.requests {
					err := sess.mux.SendData(streamID, resp)
					// TODO: should we log error somehow? We have nowhere to send it.
					_ = err
				}
				sess.mu.RUnlock()
				sess.Close()
				return
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
			if _, err := io.ReadFull(sess.conn, resp.Data); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				resp.Err = err
			}

			switch header.Status {
			case xrdproto.Error:
				resp.Err = header.Error(resp.Data)
			case xrdproto.Wait:
				if len(resp.Data) < 4 {
					resp.Err = errors.Errorf("xrootd: error decoding wait duration, want 4 bytes, got: %v", resp.Data)
				}
				duration := time.Duration(binary.BigEndian.Uint32(resp.Data)) * time.Second
				sess.mu.RLock()
				req := sess.requests[header.StreamID]
				sess.mu.RUnlock()
				go func(req []byte) {
					time.Sleep(duration)
					if _, err := sess.conn.Write(req); err != nil {
						resp := mux.ServerResponse{Err: errors.WithMessage(err, "xrootd: could not send data to the server")}
						err := sess.mux.SendData(header.StreamID, resp)
						// TODO: should we log error somehow? We have nowhere to send it.
						_ = err
						sess.cleanupRequest(header.StreamID)
					}
				}(req)
				continue
			case xrdproto.Redirect:
				redirection, err := mux.ParseRedirection(resp.Data)
				if err != nil {
					resp.Err = err
				} else {
					resp.Redirection = redirection
				}
			}

			if err := sess.mux.SendData(header.StreamID, resp); err != nil {
				if ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					continue
				}
				panic(err)
				// TODO: should we just ignore responses to unclaimed stream IDs?
			}

			if header.Status != xrdproto.OkSoFar {
				sess.cleanupRequest(header.StreamID)
			}
		}
	}
}

func (sess *session) cleanupRequest(streamID xrdproto.StreamID) {
	sess.mux.Unclaim(streamID)
	sess.mu.Lock()
	delete(sess.requests, streamID)
	sess.mu.Unlock()
}

func (sess *session) send(ctx context.Context, streamID xrdproto.StreamID, responseChannel mux.DataRecvChan, request []byte) ([]byte, *mux.Redirection, error) {
	sess.mu.Lock()
	sess.requests[streamID] = request
	sess.mu.Unlock()
	if _, err := sess.conn.Write(request); err != nil {
		return nil, nil, err
	}

	var data []byte

	for {
		select {
		case resp, more := <-responseChannel:
			if !more {
				return data, nil, nil
			}

			if resp.Err != nil {
				return nil, resp.Redirection, resp.Err
			}

			if resp.Redirection != nil {
				return nil, resp.Redirection, nil
			}

			data = append(data, resp.Data...)
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return nil, nil, err
			}
		}
	}
	panic("unreachable")
}

// Send sends the request to the server and stores the response inside the resp.
func (sess *session) Send(ctx context.Context, resp xrdproto.Response, req xrdproto.Request) (*mux.Redirection, error) {
	streamID, responseChannel, err := sess.mux.Claim()
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

	if sess.signRequirements.Needed(req) {
		data, err = sess.sign(streamID, req.ReqID(), data)
		if err != nil {
			return nil, err
		}
	}

	data, redirection, err := sess.send(ctx, streamID, responseChannel, data)
	if err != nil || redirection != nil || resp == nil {
		return redirection, err
	}

	return nil, resp.UnmarshalXrd(xrdenc.NewRBuffer(data))
}

func (sess *session) sign(streamID xrdproto.StreamID, requestID uint16, data []byte) ([]byte, error) {
	seqID := atomic.AddInt64(&sess.seqID, 1)
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
