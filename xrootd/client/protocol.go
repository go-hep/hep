// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
)

// Protocol obtains the protocol version number, type of the server and security information, such as:
// the security version, the security options, the security level, and the list of alterations
// needed to the specified predefined security level.
func (client *Client) Protocol(ctx context.Context) (protocol.Response, error) {
	var resp protocol.Response
	err := client.Send(ctx, &resp, protocol.NewRequest(client.protocolVersion, true))
	return resp, err
}
