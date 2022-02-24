// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httpio provides types for basic I/O interfaces over HTTP.
package httpio // import "go-hep.org/x/hep/groot/internal/httpio"

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	defaultClient *http.Client

	errAcceptRange = errors.New("httpio: accept-range not supported")
)

type config struct {
	ctx  context.Context
	cli  *http.Client
	auth struct {
		usr string
		pwd string
	}
}

func newConfig() *config {
	return &config{
		ctx: context.Background(),
		cli: defaultClient,
	}
}

type Option func(*config) error

// WithClient sets up Reader to use a user-provided HTTP client.
//
// By default, Reader uses an httpio-local default client.
func WithClient(cli *http.Client) Option {
	return func(c *config) error {
		c.cli = cli
		return nil
	}
}

// WithBasicAuth sets up a basic authentification scheme.
func WithBasicAuth(usr, pwd string) Option {
	return func(c *config) error {
		c.auth.usr = usr
		c.auth.pwd = pwd
		return nil
	}
}

// WithContext configures Reader to use a user-provided context.
//
// By default, Reader uses context.Background.
func WithContext(ctx context.Context) Option {
	return func(c *config) error {
		c.ctx = ctx
		return nil
	}
}

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	defaultClient = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}
