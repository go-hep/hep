// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/login"
)

// Login initializes a server connection using username
func (client *Client) Login(ctx context.Context, username string) (*login.Response, error) {
	serverResponse, err := client.call(ctx, login.RequestID, login.NewRequest(username))
	if err != nil {
		return nil, err
	}

	var response = &login.Response{}
	err = encoder.Unmarshal(serverResponse, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
