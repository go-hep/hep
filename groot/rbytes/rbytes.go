// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rbytes contains the definitions of types useful for
// serializing and deserializing ROOT data buffers.
//
// rbytes also defines the interfaces to interact with ROOT's metadata classes
// such as StreamerInfo and StreamerElements.
package rbytes // import "go-hep.org/x/hep/groot/rbytes"

import (
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
)

// RVersioner is the interface implemented by an object that
// can tell the ROOT system what is its current version.
type RVersioner interface {
	RVersion() int16
}

// StreamerInfo describes a ROOT Streamer.
type StreamerInfo interface {
	root.Named

	CheckSum() int
	ClassVersion() int
	Elements() []StreamerElement

	// BuildStreamers builds the r/w streamers.
	BuildStreamers() error
}

// StreamerElement describes a ROOT StreamerElement
type StreamerElement interface {
	root.Named

	ArrayDim() int
	ArrayLen() int
	Type() rmeta.Enum
	Offset() uintptr
	Size() uintptr
	TypeName() string
	XMin() float64
	XMax() float64
	Factor() float64
}

// StreamerInfoContext defines the protocol to retrieve a ROOT StreamerInfo
// metadata type by name.
//
// Implementations should make sure the protocol is goroutine safe.
type StreamerInfoContext interface {
	// StreamerInfo returns the named StreamerInfo.
	// If version is negative, the latest version should be returned.
	StreamerInfo(name string, version int) (StreamerInfo, error)
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
type Unmarshaler interface {
	UnmarshalROOT(r *RBuffer) error
}

// Marshaler is the interface implemented by an object that can
// marshal itself into a ROOT buffer
type Marshaler interface {
	MarshalROOT(w *WBuffer) (int, error)
}

// WStreamer is the interface implemented by types that can stream themselves
// to a ROOT buffer.
type WStreamer interface {
	WStreamROOT(*WBuffer) error
}

// RStreamer is the interface implemented by types that can stream themselves
// from a ROOT buffer.
type RStreamer interface {
	RStreamROOT(*RBuffer) error
}

// Streamer is the interface implemented by types that can stream themselves
// to and from a ROOT buffer.
type Streamer interface {
	WStreamer
	RStreamer
}

const (
	BypassStreamer                  uint32 = 1 << 12
	CannotHandleMemberWiseStreaming uint32 = 1 << 17
)
