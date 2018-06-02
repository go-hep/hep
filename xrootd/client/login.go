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
func (client *Client) Login(ctx context.Context, username, token string) (login.Response, error) {
	var resp login.Response
	err := client.Send(ctx, &resp, login.NewRequest(username, token))
	return resp, err
}
