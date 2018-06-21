// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdproto/login"
)

// Login initializes a server connection using username
// and token which can be supplied by the previous redirection response.
func (sess *session) Login(ctx context.Context, username, token string) (login.Response, error) {
	var resp login.Response
	_, err := sess.Send(ctx, &resp, login.NewRequest(username, token))
	// TODO: should we react somehow to redirection?
	return resp, err
}
