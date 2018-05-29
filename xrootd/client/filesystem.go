// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/dirlist"
)

// Filesystem contains filesystem-related methods of the XRootD protocol.
type Filesystem struct {
	c *Client
}

// Dirlist returns the contents of a directory together with the stat information.
func (fs *Filesystem) Dirlist(ctx context.Context, path string) ([]protocol.EntryStat, error) {
	serverResponse, err := fs.c.call(ctx, dirlist.RequestID, dirlist.NewRequest(path))
	if err != nil {
		return nil, err
	}

	var result = &dirlist.Response{}
	err = protocol.Unmarshal(serverResponse, result)
	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, nil
	}

	result.Data = bytes.TrimRight(result.Data, "\x00")
	linesBytes := bytes.Split(result.Data, []byte{'\n'})
	lines := make([]string, len(linesBytes))

	for i, v := range linesBytes {
		lines[i] = string(v)
	}

	if !bytes.HasPrefix(result.Data, []byte(".\n0 0 0 0\n")) {
		// That means that the server doesn't support returning stat information.
		fileStats := make([]protocol.EntryStat, len(lines))
		for i, v := range lines {
			fileStats[i] = protocol.EntryStat{Name: v}
		}
		return fileStats, nil
	}

	if len(lines)%2 != 0 {
		return nil, errors.Errorf("xrootd: wrong response size for the dirlist request: want even number of lines, got %d", len(lines))
	}

	lines = lines[2:]
	fileStats := make([]protocol.EntryStat, len(lines)/2)

	for i := 0; i < len(lines); i += 2 {
		fileStats[i/2], err = protocol.NewEntryStat(lines[i], lines[i+1])
		if err != nil {
			return nil, err
		}
	}

	return fileStats, nil
}
