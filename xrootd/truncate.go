// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/truncate"
)

// Truncate a file to a particular size.
func (client *Client) Truncate(ctx context.Context, path string, size int64) error {
	_, err := client.call(ctx, truncate.RequestID, truncate.NewRequestWithPath(path, size))
	return err
}
