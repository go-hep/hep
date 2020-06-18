// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/auth"
)

// A Client to xrootd server which allows to send requests and receive responses.
// Concurrent requests are supported.
// Zero value is invalid, Client should be instantiated using NewClient.
type Client struct {
	cancel   context.CancelFunc
	auths    map[string]auth.Auther
	username string
	// initialSessionID is the sessionID of the server which is used as default
	// for all requests that don't specify sessionID explicitly.
	// Any failed request with another sessionID should be redirected to the initialSessionID.
	// See http://xrootd.org/doc/dev45/XRdv310.pdf, page 11 for details.
	initialSessionID string
	mu               sync.RWMutex
	sessions         map[string]*cliSession

	maxRedirections int
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
	for _, provider := range defaultProviders {
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

	client := &Client{
		cancel:          cancel,
		auths:           make(map[string]auth.Auther),
		username:        username,
		sessions:        make(map[string]*cliSession),
		maxRedirections: 10,
	}

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

	_, err := client.getSession(ctx, address, "")
	if err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

// Close closes the connection. Any blocked operation will be unblocked and return error.
func (client *Client) Close() error {
	if client == nil {
		return os.ErrInvalid
	}

	defer client.cancel()
	client.mu.Lock()
	defer client.mu.Unlock()
	var errs []error
	for _, session := range client.sessions {
		err := session.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		return fmt.Errorf("xrootd: could not close client: %v", errs)
	}
	return nil
}

// Send sends the request to the server and stores the response inside the resp.
// If the resp is nil, then no response is stored.
// Send returns a session id which identifies the server that provided response.
func (client *Client) Send(ctx context.Context, resp xrdproto.Response, req xrdproto.Request) (string, error) {
	return client.sendSession(ctx, client.initialSessionID, resp, req)
}

func (client *Client) sendSession(ctx context.Context, sessionID string, resp xrdproto.Response, req xrdproto.Request) (string, error) {
	client.mu.RLock()
	session, ok := client.sessions[sessionID]
	client.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("xrootd: session with id = %q was not found", sessionID)
	}

	redirection, err := session.Send(ctx, resp, req)
	if err != nil {
		return sessionID, err
	}

	for cnt := client.maxRedirections; redirection != nil && cnt > 0; cnt-- {
		sessionID = redirection.Addr
		session, err = client.getSession(ctx, sessionID, redirection.Token)
		if err != nil {
			return sessionID, err
		}
		if fp, ok := req.(xrdproto.FilepathRequest); ok {
			fp.SetOpaque(redirection.Opaque)
		}
		// TODO: we should check if the request contains file handle and re-issue open request in that case.
		redirection, err = session.Send(ctx, resp, req)
		if err != nil {
			return sessionID, err
		}
	}

	if redirection != nil {
		err = fmt.Errorf("xrootd: received %d redirections in a row, aborting request", client.maxRedirections)
	}

	return sessionID, err
}

func (client *Client) getSession(ctx context.Context, address, token string) (*cliSession, error) {
	client.mu.RLock()
	v, ok := client.sessions[address]
	client.mu.RUnlock()
	if ok {
		return v, nil
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	session, err := newSession(ctx, address, client.username, token, client)
	if err != nil {
		return nil, err
	}
	client.sessions[address] = session

	if len(client.initialSessionID) == 0 {
		client.initialSessionID = address
	}
	// TODO: check if initial sessionID should be changed.
	// See http://xrootd.org/doc/dev45/XRdv310.pdf, p. 11 for details.

	return session, nil
}
