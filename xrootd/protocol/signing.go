// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol // import "go-hep.org/x/hep/xrootd/protocol"

import (
	"go-hep.org/x/hep/xrootd/protocol/dirlist"
	"go-hep.org/x/hep/xrootd/protocol/protocol"
)

// SignRequirements implements a way to check if request should be signed
// according to XRootD protocol specification v. 3.1.0, p.75-76.
type SignRequirements struct {
	requirements map[uint16]protocol.RequestLevel
}

// Needed return whether the request should be signed.
// "Modifies" indicates that request modifies data or metadata
// and is used to handle the "signLikely" level which specifies that
// request should be signed only if it modifies data.
// For the list of actual examples see XRootD protocol specification v. 3.1.0, p.76.
func (sr *SignRequirements) Needed(requestID uint16, modifies bool) bool {
	v, exist := sr.requirements[requestID]
	if !exist || v == protocol.SignNone {
		return false
	}
	if v == protocol.SignLikely && !modifies {
		return false
	}
	return true
}

// DefaultSignRequirements creates a default SignRequirements with "None" security level.
func DefaultSignRequirements() SignRequirements {
	return NewSignRequirements(protocol.NoneLevel, nil)
}

// NewSignRequirements creates a SignRequirements according to provided security level and security overrides.
func NewSignRequirements(level protocol.SecurityLevel, overrides []protocol.SecurityOverride) SignRequirements {
	var sr = SignRequirements{make(map[uint16]protocol.RequestLevel)}

	if level >= protocol.Compatible {
		// TODO: set requirements
	}
	if level >= protocol.Standard {
		// TODO: set requirements
	}
	if level >= protocol.Intense {
		// TODO: set requirements
	}
	if level >= protocol.Pedantic {
		// TODO: set requirements
		sr.requirements[dirlist.RequestID] = protocol.SignNeeded
	}

	for _, override := range overrides {
		// TODO: use auth.RequestID instead of 3000.
		requestID := 3000 + uint16(override.RequestIndex)
		sr.requirements[requestID] = override.RequestLevel
	}

	return sr
}
