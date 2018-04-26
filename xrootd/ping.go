// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/ping"
)

// Ping determines if the server is still alive
func (client *Client) Ping(ctx context.Context) error {
	_, err := client.call(ctx, ping.RequestID, ping.NewRequest())
	return err
}
