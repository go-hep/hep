// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"bytes"
	"context"
	"os/user"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/xrdproto/auth"
)

func (client *Client) auth(ctx context.Context, securityInformation []byte) error {
	securityInformation = bytes.TrimLeft(securityInformation, "&")
	securityProviders := bytes.Split(securityInformation, []byte{'&'})

	var errs []error
	for _, securityProvider := range securityProviders {
		securityProvider = bytes.TrimLeft(securityProvider, "P=")
		if bytes.Equal(securityProvider, auth.UnixType[:]) {
			u, err := user.Current()
			if err != nil {
				errs = append(errs, errors.WithMessage(err, "xrootd: could not authorize using unix"))
				continue
			}
			g, err := lookupGroupID(u)
			if err != nil {
				errs = append(errs, errors.WithMessage(err, "xrootd: could not authorize using unix"))
				continue
			}

			_, err = client.call(ctx, auth.NewUnixRequest(u.Username, g))
			if err != nil {
				errs = append(errs, errors.WithMessage(err, "xrootd: could not authorize using unix"))
				continue
			} else {
				return nil
			}
		}
	}

	return errors.Errorf("xrootd: could not authorize:\n%v", errs)
}
