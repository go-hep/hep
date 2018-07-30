// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

// StreamControl provides concurrency-safe control to input and output streamers.
type StreamControl struct {
	Ports []Port        // list of ports streamers will read-from or write-to
	Ctx   chan Context  // contexts to read-from or write-to
	Err   chan error    // errors encountered during reading-from or writing-to
	Quit  chan struct{} // closed to signify in/out-streamers should stop reading-from/writing-to
}

// InputStreamer reads data from the underlying io.Reader
// and puts it into fwk's Context
type InputStreamer interface {

	// Connect connects the InputStreamer to the underlying io.Reader,
	// and configure it to only read-in the data specified in ports.
	Connect(ports []Port) error

	// Read reads the data from the underlying io.Reader
	// and puts it in the store associated with the fwk.Context ctx
	Read(ctx Context) error

	// Disconnect disconnects the InputStreamer from the underlying io.Reader,
	// possibly computing some statistics data.
	// It does not (and can not) close the underlying io.Reader.
	Disconnect() error
}

// OutputStreamer gets data from the Context
// and writes it to the underlying io.Writer
type OutputStreamer interface {

	// Connect connects the OutputStreamer to the underlying io.Writer,
	// and configure it to only write-out the data specified in ports.
	Connect(ports []Port) error

	// Write gets the data from the store associated with the fwk.Context ctx
	// and writes it to the underlying io.Writer
	Write(ctx Context) error

	// Disconnect disconnects the OutputStreamer from the underlying io.Writer,
	// possibly computing some statistics data.
	// It does not (and can not) close the underlying io.Writer.
	Disconnect() error
}
