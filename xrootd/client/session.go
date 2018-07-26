// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
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
//
// If the request that supports sending data over a separate socket is issued,
// the session tries to obtain a sub-session to the same server using a `bind` request.
// If the connection is successful, the request is sent specifying that socket for the data exchange.
// Otherwise, a default socket connected to the server is used.
type session struct {
	ctx              context.Context
	cancel           context.CancelFunc
	conn             net.Conn
	mux              *mux.Mux
	protocolVersion  int32
	signRequirements signing.Requirements
	seqID            int64
	mu               sync.RWMutex
	requests         map[xrdproto.StreamID]pendingRequest

	subCreateMu sync.Mutex   // subCreateMu is used to serialize the creation of sub-sessions.
	subsMu      sync.RWMutex // subsMu is used to serialize the access to the subs map.
	subs        map[xrdproto.PathID]*session

	maxSubs   int
	freeSubs  chan xrdproto.PathID
	isSub     bool // indicates whether this session is a sub-session.
	client    *Client
	sessionID string
	addr      string
	loginID   [16]byte
	pathID    xrdproto.PathID
}

// pendingRequest is a request that has been sent to the remote server.
type pendingRequest struct {
	// Header is the header part of the request.
	// It may contain all of the request content if there is no data that is
	// intended to be sent over a separate socket.
	Header []byte

	// Data is the data part of the request that is intended to be sent over a separate socket.
	Data []byte

	// PathID is the identifier of the socket which should be used to read or write a data.
	PathID xrdproto.PathID
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
		ctx:       ctx,
		cancel:    cancel,
		conn:      conn,
		mux:       mux.New(),
		subs:      make(map[xrdproto.PathID]*session),
		freeSubs:  make(chan xrdproto.PathID),
		requests:  make(map[xrdproto.StreamID]pendingRequest),
		client:    client,
		sessionID: addr,
		addr:      addr,
		maxSubs:   8, // TODO: The value of 8 is just a guess. Change it?
	}

	go sess.consume()

	if err := sess.handshake(ctx); err != nil {
		sess.Close()
		return nil, err
	}

	securityInfo, err := sess.Login(ctx, username, token)
	if err != nil {
		sess.Close()
		return nil, err
	}

	sess.loginID = securityInfo.SessionID

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

	var errs []error
	for _, child := range sess.subs {
		err := child.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if !sess.isSub {
		sess.mux.Close()
	}

	// TODO: should we remove session here somehow?
	err := sess.conn.Close()
	if err != nil {
		errs = append(errs, err)
	}
	if errs != nil {
		return errors.Errorf("xrootd: errors occured during closing of the session: %v", errs)
	}
	return nil
}

// handleReadError handles an error encountered while reading and parsing a response.
// If the current session is equal to the initial, the error is considered critical and handleReadError panics.
// Otherwise, the current session is closed and all requests are redirected to the initial session.
// See http://xrootd.org/doc/dev45/XRdv310.pdf, p. 11 for details.
func (sess *session) handleReadError(err error) {
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
}

// handleWaitResponse handles a "kXR_wait" response by re-issuing the request with streamID
// after the number of seconds encoded in data.
// See http://xrootd.org/doc/dev45/XRdv310.pdf, p. 35 for the specification of the response.
func (sess *session) handleWaitResponse(streamID xrdproto.StreamID, data []byte) error {
	var resp xrdproto.WaitResponse
	rBuffer := xrdenc.NewRBuffer(data)
	if err := resp.UnmarshalXrd(rBuffer); err != nil {
		return err
	}

	sess.mu.RLock()
	req, ok := sess.requests[streamID]
	sess.mu.RUnlock()
	if !ok {
		return errors.Errorf("xrootd: could not find a request with stream id equal to %v", streamID)
	}

	go func(req pendingRequest) {
		time.Sleep(resp.Duration)
		if err := sess.writeRequest(req); err != nil {
			resp := mux.ServerResponse{Err: errors.WithMessage(err, "xrootd: could not send data to the server")}
			err := sess.mux.SendData(streamID, resp)
			// TODO: should we log error somehow? We have nowhere to send it.
			_ = err
			sess.cleanupRequest(streamID)
		}
	}(req)

	return nil
}

func (sess *session) consume() {
	var header xrdproto.ResponseHeader
	var headerBytes = make([]byte, xrdproto.ResponseHeaderLength)
	var resp mux.ServerResponse

	for {
		select {
		case <-sess.ctx.Done():
			// TODO: Should wait for active requests to be completed?
			return
		default:
			var err error
			resp.Data, err = xrdproto.ReadResponseWithReuse(sess.conn, headerBytes, &header)
			if err != nil {
				if sess.ctx.Err() != nil {
					// something happened to the context.
					// ignore this error.
					return
				}
				sess.handleReadError(err)
			}
			resp.Err = nil
			resp.Redirection = nil

			switch header.Status {
			case xrdproto.Error:
				resp.Err = header.Error(resp.Data)
			case xrdproto.Wait:
				resp.Err = sess.handleWaitResponse(header.StreamID, resp.Data)
				if resp.Err == nil {
					continue
				}
			case xrdproto.Redirect:
				resp.Redirection, resp.Err = mux.ParseRedirection(resp.Data)
			}

			if err := sess.mux.SendData(header.StreamID, resp); err != nil {
				if sess.ctx.Err() != nil {
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

func (sess *session) writeRequest(request pendingRequest) error {
	if request.PathID == 0 {
		request.Header = append(request.Header, request.Data...)
	}

	if _, err := sess.conn.Write(request.Header); err != nil {
		return err
	}

	if request.PathID != 0 && len(request.Data) > 0 {
		sess.subsMu.RLock()
		conn, ok := sess.subs[request.PathID]
		sess.subsMu.RUnlock()
		if !ok {
			return errors.Errorf("xrootd: connection with wrong pathID = %v was requested", request.PathID)
		}
		if _, err := conn.conn.Write(request.Data); err != nil {
			return err
		}
	}
	return nil
}

func (sess *session) send(ctx context.Context, streamID xrdproto.StreamID, responseChannel mux.DataRecvChan, header, body []byte, pathID xrdproto.PathID) ([]byte, *mux.Redirection, error) {
	if pathID == 0 {
		header = append(header, body...)
	}
	request := pendingRequest{Header: header, Data: body, PathID: pathID}
	sess.mu.Lock()
	sess.requests[streamID] = request
	sess.mu.Unlock()

	if err := sess.writeRequest(request); err != nil {
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

	var pathID xrdproto.PathID = 0
	var pathData []byte
	if dr, ok := req.(xrdproto.DataRequest); ok {
		var err error
		pathID, err = sess.claimPathID(ctx)
		if err != nil {
			// Should we log error somehow?
			// Fallback to sending the data over a single connection.
			pathID = 0
		}
		defer sess.unclaimPathID(pathID)
		dr.SetPathID(pathID)
		pathData = dr.PathData()
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

	data, redirection, err := sess.send(ctx, streamID, responseChannel, data, pathData, pathID)
	if err != nil || redirection != nil || resp == nil {
		return redirection, err
	}

	return nil, resp.UnmarshalXrd(xrdenc.NewRBuffer(data))
}

func (sess *session) claimPathID(ctx context.Context) (xrdproto.PathID, error) {
	select {
	case child := <-sess.freeSubs:
		return child, nil
	default:
		sess.subCreateMu.Lock()
		defer sess.subCreateMu.Unlock()

		sess.subsMu.RLock()
		if len(sess.subs) >= sess.maxSubs {
			sess.subsMu.RUnlock()
			return 0, errors.Errorf("xrootd: could not claimPathID: all of %d connections are taken", sess.maxSubs)
		}
		sess.subsMu.RUnlock()

		ds, err := newSubSession(ctx, sess)
		if err != nil {
			return 0, err
		}
		sess.subsMu.Lock()
		sess.subs[ds.pathID] = ds
		sess.subsMu.Unlock()

		return ds.pathID, nil
	}
}

func (sess *session) unclaimPathID(pathID xrdproto.PathID) {
	if pathID == 0 {
		return
	}
	go func() {
		select {
		case <-sess.ctx.Done():
			return
		case sess.freeSubs <- pathID:
		}
	}()
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

func newSubSession(ctx context.Context, parent *session) (*session, error) {
	ctx, cancel := context.WithCancel(ctx)

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", parent.addr)
	if err != nil {
		cancel()
		return nil, err
	}

	sess := &session{
		ctx:       ctx,
		cancel:    cancel,
		conn:      conn,
		mux:       parent.mux,
		subs:      make(map[xrdproto.PathID]*session),
		requests:  make(map[xrdproto.StreamID]pendingRequest),
		client:    parent.client,
		sessionID: parent.addr,
		addr:      parent.addr,
		isSub:     true,
	}

	go sess.consume()

	if err := sess.handshake(ctx); err != nil {
		sess.Close()
		return nil, err
	}

	pathID, err := sess.bind(ctx, parent.loginID)
	if err != nil {
		sess.Close()
		return nil, err
	}

	sess.pathID = pathID
	return sess, nil
}
