// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package login contains the structures describing request and response for login request.
// Login request should be issued prior to most of the other
// requests (see http://xrootd.org/doc/dev45/XRdv310.pdf, p.10).
// As part of the response, SecurityInformation may be provided,
// indicating that an auth request is required. SecurityInformation
// defines the available authentication protocols together with some additional parameters.
// See XRootD protocol specification, page 127 for further information
// about the format of the SecurityInformation.
package login // import "go-hep.org/x/hep/xrootd/protocol/login"

import (
	"os"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3007

// ResponseLength is the length of the Response assuming that SecurityInformation is empty.
const ResponseLength = 16

// Response is a response for the login request, which contains the session id and the security information.
type Response struct {
	SessionID           [16]byte
	SecurityInformation []byte
}

func (resp *Response) MarshalXrd() ([]byte, error) {
	var enc xrdenc.Encoder
	enc.WriteBytes(resp.SessionID[:])
	enc.WriteBytes(resp.SecurityInformation)
	return enc.Bytes(), nil
}

func (resp *Response) UnmarshalXrd(data []byte) error {
	dec := xrdenc.NewDecoder(data)
	dec.ReadBytes(resp.SessionID[:])
	resp.SecurityInformation = append(resp.SecurityInformation, dec.Bytes()...)
	return nil
}

// Request holds the login request parameters.
type Request struct {
	Pid          int32   // Pid is the process number associated with this connection.
	Username     [8]byte // Username is the unauthenticated name of the user to be associated with the connection.
	_            byte    // Reserved is an area reserved for future use.
	Ability      byte    // Ability are the client's extended capabilities. See xrootd protocol specification, p. 56.
	Capabilities byte    // Capabilities are the Client capabilities. It is 4 for v3.1.0 client without async support.
	Role         byte    // Role is the role being assumed for this login: administrator or regular user.
	Token        []byte  // Token is the token supplied by the previous redirection response, plus optional elements.
}

// Capabilities for v3.1.0 client without async support.
const clientCapabilities byte = 4

// NewRequest forms a Request according to provided parameters.
func NewRequest(username, token string) *Request {
	var usernameBytes [8]byte
	copy(usernameBytes[:], username)

	return &Request{
		Pid:          int32(os.Getpid()),
		Username:     usernameBytes,
		Capabilities: clientCapabilities,
		Token:        []byte(token),
	}
}

// ReqID implements protocol.Request.ReqID
func (req *Request) ReqID() uint16 { return RequestID }

// MarshalXrd implements xrootd/protocol.Marshaler
func (o *Request) MarshalXrd() ([]byte, error) {
	var enc xrdenc.Encoder
	enc.WriteI32(o.Pid)
	enc.WriteBytes(o.Username[:])
	enc.WriteReserved(1)
	enc.WriteU8(o.Ability)
	enc.WriteU8(o.Capabilities)
	enc.WriteU8(o.Role)
	enc.WriteLen(len(o.Token))
	enc.WriteBytes(o.Token)
	return enc.Bytes(), nil
}

// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *Request) UnmarshalXrd(data []byte) error {
	dec := xrdenc.NewDecoder(data)
	o.Pid = dec.ReadI32()
	dec.ReadBytes(o.Username[:])
	dec.Skip(1)
	o.Ability = dec.ReadU8()
	o.Capabilities = dec.ReadU8()
	o.Role = dec.ReadU8()
	o.Token = make([]byte, dec.ReadLen())
	dec.ReadBytes(o.Token)
	return nil
}
