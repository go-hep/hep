// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/bind"
)

// Bind binds the client's socket to a pre-existing session ID.
func (client *Client) Bind(ctx context.Context, sessionID [16]byte) (byte, error) {
	resp, err := client.call(ctx, bind.RequestID, bind.NewRequest(sessionID))
	if err != nil {
		return 0, err
	}

	var result bind.Response
	err = encoder.Unmarshal(resp, &result)
	if err != nil {
		return 0, err
	}

	return result.PathID, nil
}
