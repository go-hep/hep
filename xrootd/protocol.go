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
func (client *Client) Protocol(ctx context.Context) (response *protocol.Response, securityInfo *protocol.SecurityInfo, err error) {
	serverResponse, err := client.call(ctx, protocol.RequestID, protocol.NewRequest(client.protocolVersion))
	if err != nil {
		return
	}

	response = &protocol.Response{}
	securityInfo = &protocol.SecurityInfo{}

	err = encoder.Unmarshal(serverResponse, response)
	if err != nil {
		return
	}

	if len(serverResponse) > 8 {
		err = encoder.Unmarshal(serverResponse, securityInfo)
	}
	return
}
