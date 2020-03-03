// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"bytes"
	"context"
	"fmt"

	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/host"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/krb5"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/unix"
)

// defaultProviders is the list of authentification providers a xrootd client will use by default.
var defaultProviders = []auth.Auther{
	krb5.Default,
	unix.Default,
	host.Default,
}

func (sess *cliSession) auth(ctx context.Context, securityInformation []byte) error {
	securityInformation = bytes.TrimLeft(securityInformation, "&")
	providerInfos := bytes.Split(securityInformation, []byte{'&'})

	var errs []error
	for _, providerInfo := range providerInfos {
		providerInfo = bytes.TrimLeft(providerInfo, "P=")[:]
		paramsData := bytes.Split(providerInfo, []byte{','})
		params := make([]string, len(paramsData))
		for i := range paramsData {
			params[i] = string(paramsData[i])
		}
		provider := params[0]
		params = params[1:]

		auther, ok := sess.client.auths[provider]
		if !ok {
			errs = append(errs, fmt.Errorf("xrootd: could not authorize using %s: provider was not found", provider))
			continue
		}
		r, err := auther.Request(params)
		if err != nil {
			errs = append(errs, fmt.Errorf("xrootd: could not authorize using %s: %w", provider, err))
			continue
		}
		_, err = sess.Send(ctx, nil, r)
		// TODO: should we react somehow to redirection?
		if err != nil {
			errs = append(errs, fmt.Errorf("xrootd: could not authorize using %s: %w", provider, err))
			continue
		}
		return nil
	}

	return fmt.Errorf("xrootd: could not authorize:\n%v", errs)
}
