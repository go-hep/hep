// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"bytes"
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/dirlist"
)

// Dirlist returns contents of a directory
func (client *Client) Dirlist(ctx context.Context, path string) ([]string, error) {
	serverResponse, err := client.call(ctx, dirlist.RequestID, dirlist.NewRequest(path))
	if err != nil {
		return nil, err
	}

	var result = &dirlist.Response{}
	err = encoder.Unmarshal(serverResponse, result)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return []string{}, nil
	}

	strings := bytes.Split(result.Data, []byte{'\n'})

	resultStrings := make([]string, len(strings))

	for i := 0; i < len(strings); i++ {
		strings[i] = bytes.Trim(strings[i], "\x00")
		resultStrings[i] = string(strings[i])
	}

	return resultStrings, nil
}
