// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/stat"
)

// Stat obtains status information for a path
func (client *Client) Stat(ctx context.Context, path string) (*stat.Response, error) {
	serverResponse, err := client.call(ctx, stat.RequestID, stat.NewRequest(path))
	if err != nil {
		return nil, err
	}

	return stat.ParseReponsee(serverResponse)
}
