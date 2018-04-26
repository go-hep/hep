// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/sync"
)

// Sync commits all pending writes to an open file
func (client *Client) Sync(ctx context.Context, fileHandle [4]byte) error {
	_, err := client.call(ctx, sync.RequestID, sync.NewRequest(fileHandle))
	return err
}
