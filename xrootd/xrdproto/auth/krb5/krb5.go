// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package krb5 contains the implementation of krb5 (Kerberos) security provider.
package krb5

import (
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"gopkg.in/jcmturner/gokrb5.v5/client"
	"gopkg.in/jcmturner/gokrb5.v5/config"
	"gopkg.in/jcmturner/gokrb5.v5/credentials"
	"gopkg.in/jcmturner/gokrb5.v5/crypto"
	"gopkg.in/jcmturner/gokrb5.v5/messages"
	"gopkg.in/jcmturner/gokrb5.v5/types"
)

// Default is a Kerberos 5 client configured from cached credentials.
// If the credentials could not be correctly configured, Default will be nil.
var Default auth.Auther

func init() {
	v, err := WithCredCache()
	if err == nil {
		Default = v
	}
}

// Auth implements krb5 (Kerberos) security provider.
type Auth struct {
	client *client.Client
}

// WithPassword creates a new Auth configured from the provided user, realm and password.
func WithPassword(user, realm, password string) (*Auth, error) {
	krb := client.NewClientWithPassword(user, realm, password)

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not load kerberos-5 configuration")
	}

	krb.WithConfig(cfg)

	err = krb.Login()
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not login")
	}

	return &Auth{client: &krb}, nil
}

// WithCredCache creates a new Auth configured from cached credentials.
func WithCredCache() (*Auth, error) {
	cred, err := credentials.LoadCCache(cachePath())
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not load kerberos-5 cached credentials")
	}

	krb, err := client.NewClientFromCCache(cred)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not create kerberos-5 client from cached credentials")
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not load kerberos-5 configuration")
	}

	krb.WithConfig(cfg)

	return &Auth{client: &krb}, nil
}

// WithClient creates a new Auth using the provided krb5 client.
func WithClient(client *client.Client) *Auth {
	return &Auth{client: client}

}

// Provider implements auth.Auther
func (*Auth) Provider() string {
	return "krb5"
}

// Type indicates that krb5 (Kerberos) authentication protocol is used.
var Type = [4]byte{'k', 'r', 'b', '5'}

// Request implements auth.Auther
func (a *Auth) Request(params []string) (*auth.Request, error) {
	if len(params) == 0 {
		return nil, errors.New("auth/krb5: want at least 1 parameter, got 0")
	}
	serviceName := string(params[0])
	if strings.Contains(serviceName, "@") {
		// Service name from the XRootD server may be in the following format: "xrootd/server.example.com@example.com"
		// While gokrb5 expects server name in that format: "xrootd/server.example.com".
		// The "@example.com" part (realm) will be guessed from the instance name "server.example.com".
		index := strings.Index(serviceName, "@")
		serviceName = serviceName[:index]
	}
	tkt, key, err := a.client.GetServiceTicket(serviceName)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not retrieve kerberos service ticket")
	}
	authenticator, err := types.NewAuthenticator(a.client.Credentials.Realm, a.client.Credentials.CName)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not create kerberos authenticator")
	}
	etype, err := crypto.GetEtype(key.KeyType)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not retrieve crypto key type")
	}
	err = authenticator.GenerateSeqNumberAndSubKey(key.KeyType, etype.GetKeyByteSize())
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not generate sequence number or sub key")
	}
	APReq, err := messages.NewAPReq(tkt, key, authenticator)
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not generate AP request")
	}
	request, err := APReq.Marshal()
	if err != nil {
		return nil, errors.WithMessage(err, "auth/krb5: could not marshal AP request")
	}

	return &auth.Request{Type: Type, Credentials: "krb5\000" + string(request)}, nil
}

var (
	_ auth.Auther = (*Auth)(nil)
)
