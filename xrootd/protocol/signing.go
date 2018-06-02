// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol // import "go-hep.org/x/hep/xrootd/protocol"

import "go-hep.org/x/hep/xrootd/internal/xrdenc"

// RequestLevel is the security requirement that the associated request is to have.
type RequestLevel byte

const (
	SignNone   RequestLevel = 0 // SignNone indicates that the request need not to be signed.
	SignLikely RequestLevel = 1 // SignLikely indicates that the request must be signed if it modifies data.
	SignNeeded RequestLevel = 2 // SignNeeded indicates that the request mush be signed.
)

// SecurityLevel is the predefined security level that specifies which requests should be signed.
// See specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, p. 75.
type SecurityLevel byte

const (
	// NoneLevel indicates that no request needs to be signed.
	NoneLevel SecurityLevel = 0
	// Compatible indicates that only potentially destructive requests need to be signed.
	Compatible SecurityLevel = 1
	// Standard indicates that potentially destructive requests
	// as well as certain non-destructive requests need to be signed.
	Standard SecurityLevel = 2
	// Intense indicates that request that may reveal metadata or modify data need to be signed.
	Intense SecurityLevel = 3
	// Pedantic indicates that all requests need to be signed.
	Pedantic SecurityLevel = 4
)

// SecurityOverrideLength is the length of SecurityOverride in bytes.
const SecurityOverrideLength = 2

// SecurityOverride is an alteration needed to the specified predefined security level.
// It consists of the request index and the security requirement the associated request should have.
// Request index is calculated as:
//     (request code) - (request code of Auth request)
// according to xrootd protocol specification.
type SecurityOverride struct {
	RequestIndex byte
	RequestLevel RequestLevel
}

// MarshalXrd implements xrootd/protocol.Marshaler
func (o SecurityOverride) MarshalXrd(enc *xrdenc.WBuffer) error {
	enc.WriteU8(o.RequestIndex)
	enc.WriteU8(byte(o.RequestLevel))
	return nil
}

// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *SecurityOverride) UnmarshalXrd(dec *xrdenc.RBuffer) error {
	o.RequestIndex = dec.ReadU8()
	o.RequestLevel = RequestLevel(dec.ReadU8())
	return nil
}

// SignRequirements implements a way to check if request should be signed
// according to XRootD protocol specification v. 3.1.0, p.75-76.
type SignRequirements struct {
	requirements map[uint16]RequestLevel
}

// Needed return whether the request should be signed.
// "Modifies" indicates that request modifies data or metadata
// and is used to handle the "signLikely" level which specifies that
// request should be signed only if it modifies data.
// For the list of actual examples see XRootD protocol specification v. 3.1.0, p.76.
func (sr *SignRequirements) Needed(requestID uint16, modifies bool) bool {
	v, exist := sr.requirements[requestID]
	if !exist || v == SignNone {
		return false
	}
	if v == SignLikely && !modifies {
		return false
	}
	return true
}

// DefaultSignRequirements creates a default SignRequirements with "None" security level.
func DefaultSignRequirements() SignRequirements {
	return NewSignRequirements(NoneLevel, nil)
}

// NewSignRequirements creates a SignRequirements according to provided security level and security overrides.
func NewSignRequirements(level SecurityLevel, overrides []SecurityOverride) SignRequirements {
	var sr = SignRequirements{make(map[uint16]RequestLevel)}

	if level >= Compatible {
		// TODO: set requirements
	}
	if level >= Standard {
		// TODO: set requirements
	}
	if level >= Intense {
		// TODO: set requirements
	}
	if level >= Pedantic {
		// TODO: set requirements
	}

	for _, override := range overrides {
		// TODO: use auth.RequestID instead of 3000.
		requestID := 3000 + uint16(override.RequestIndex)
		sr.requirements[requestID] = override.RequestLevel
	}

	return sr
}
