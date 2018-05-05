// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"

	"go-hep.org/x/hep/xrootd/encoder"
	"go-hep.org/x/hep/xrootd/requests/protocol"
)

// Protocol obtains the protocol version number, type of server and security information
func (client *Client) Protocol(ctx context.Context) (result *protocol.Response, securityInfo *protocol.SecurityInfo, err error) {
	resp, err := client.call(ctx, protocol.RequestID, protocol.NewRequest(client.protocolVersion))
	if err != nil {
		return
	}

	result = &protocol.Response{}
	securityInfo = &protocol.SecurityInfo{}

	if err = encoder.Unmarshal(resp, result); err != nil {
		return
	}

	if len(resp) > 8 {
		err = encoder.Unmarshal(resp, securityInfo)
	}
	return
}
