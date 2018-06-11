// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdproto // import "go-hep.org/x/hep/xrootd/xrdproto"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"go-hep.org/x/hep/xrootd/xrdproto/chmod"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/mkdir"
	"go-hep.org/x/hep/xrootd/xrdproto/mv"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/rmdir"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/statx"
	"go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/verifyw"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

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

// MarshalXrd implements xrdproto.Marshaler
func (o SecurityOverride) MarshalXrd(enc *xrdenc.WBuffer) error {
	enc.WriteU8(o.RequestIndex)
	enc.WriteU8(byte(o.RequestLevel))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
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
// For the list of actual examples see XRootD protocol specification v. 3.1.0, p.76.
func (sr *SignRequirements) Needed(request Request) bool {
	v, exist := sr.requirements[request.ReqID()]
	if !exist || v == SignNone {
		return false
	}
	if v == SignLikely && !request.ShouldSign() {
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
		sr.requirements[chmod.RequestID] = SignNeeded
		sr.requirements[mv.RequestID] = SignNeeded
		sr.requirements[open.RequestID] = SignLikely
		sr.requirements[rm.RequestID] = SignNeeded
		sr.requirements[rmdir.RequestID] = SignNeeded
		sr.requirements[truncate.RequestID] = SignNeeded
	}
	if level >= Standard {
		// TODO: set requirements
		sr.requirements[mkdir.RequestID] = SignNeeded
		sr.requirements[open.RequestID] = SignNeeded
	}
	if level >= Intense {
		// TODO: set requirements
		sr.requirements[xrdclose.RequestID] = SignNeeded
		sr.requirements[verifyw.RequestID] = SignNeeded
		sr.requirements[write.RequestID] = SignNeeded
	}
	if level >= Pedantic {
		// TODO: set requirements
		sr.requirements[dirlist.RequestID] = SignNeeded
		sr.requirements[read.RequestID] = SignNeeded
		sr.requirements[stat.RequestID] = SignNeeded
		sr.requirements[statx.RequestID] = SignNeeded
		sr.requirements[sync.RequestID] = SignNeeded
	}

	for _, override := range overrides {
		requestID := auth.RequestID + uint16(override.RequestIndex)
		sr.requirements[requestID] = override.RequestLevel
	}

	return sr
}
