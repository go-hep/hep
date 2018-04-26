// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/read"
)

// Read data from an open file
func (client *Client) Read(ctx context.Context, fileHandle [4]byte, offset int64, length int32) ([]byte, error) {
	serverResponse, err := client.call(ctx, read.RequestID, read.NewRequest(fileHandle, offset, length))
	if err != nil {
		return nil, err
	}

	var result = &read.Response{}
	if err = encoder.Unmarshal(serverResponse, result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
