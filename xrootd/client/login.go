// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/login"
)

// Login initializes a server connection using username
// and token which can be supplied by the previous redirection response.
func (client *Client) Login(ctx context.Context, username string, token string) (login.Response, error) {
	serverResponse, err := client.call(ctx, login.NewRequest(username, token))
	if err != nil {
		return login.Response{}, err
	}

	var response login.Response
	err = protocol.Unmarshal(serverResponse, &response)
	if err != nil {
		return login.Response{}, err
	}

	return response, nil
}
