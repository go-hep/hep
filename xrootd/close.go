// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/close"
)

// Close a previously opened file by handle
func (client *Client) Close(ctx context.Context, fileHandle [4]byte, fileSize int64) error {
	_, err := client.call(ctx, close.RequestID, close.NewRequest(fileHandle, fileSize))
	return err
}
