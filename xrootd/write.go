// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/requests/write"
)

// Write writes the data to an open file
func (client *Client) Write(ctx context.Context, fileHandle [4]byte, offset int64, pathID byte, data []byte) error {
	_, err := client.call(ctx, write.RequestID, write.NewRequest(fileHandle, offset, pathID, data))
	return err
}
