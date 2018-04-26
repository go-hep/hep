// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/open"
)

// Open returns file handle for a file.
func (client *Client) Open(ctx context.Context, path string, mode open.Mode, options open.Options) ([4]byte, error) {
	serverResponse, err := client.call(ctx, open.RequestID, open.NewRequest(path, mode, options))
	if err != nil {
		return [4]byte{}, err
	}

	var result = &open.Response{}
	err = encoder.Unmarshal(serverResponse, result)
	if err != nil {
		return [4]byte{}, err
	}

	return result.FileHandle, nil
}
