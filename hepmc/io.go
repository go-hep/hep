// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepmc

import (
	"errors"
	"io"
)

// Reader is the interface that wraps the Read method.
//
// Read reads an HepMC event from the underlying storage into
// the provided event.
// Read returns io.EOF when no more events are available.
type Reader interface {
	Read(evt *Event) error
}

// Writer is the interface that wraps the Write method.
//
// Write writes the provided HepMC event to the underlying storage.
type Writer interface {
	Write(evt Event) error
}

// Copy copies events from r to w until either EOF is reached or an
// error occurs. It returns the number of events copied and the first error
// encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == io.EOF.
func Copy(w Writer, r Reader) (n int64, err error) {
	var evt Event

loop:
	for {
		err = r.Read(&evt)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			break loop
		}

		err = w.Write(evt)
		if err != nil {
			break loop
		}

		n++
	}

	return n, err
}

// ASCIIReader reads ASCII HepMC-v2 data.
type ASCIIReader struct {
	dec *Decoder
}

// NewASCIIReader creates a HepMC reader that reads data
// from r, in the HepMC-v2 ASCII format.
func NewASCIIReader(r io.Reader) *ASCIIReader {
	return &ASCIIReader{dec: NewDecoder(r)}
}

func (r *ASCIIReader) Read(evt *Event) error {
	return r.dec.Decode(evt)
}

// ASCIIWriter writes ASCII HepMC-v2 data.
type ASCIIWriter struct {
	enc *Encoder
}

// NewASCIIWriter creates a HepMC writer that writes data
// to w, in the HepMC-v2 ASCII format.
func NewASCIIWriter(w io.Writer) *ASCIIWriter {
	return &ASCIIWriter{enc: NewEncoder(w)}
}

func (w *ASCIIWriter) Write(evt Event) error {
	return w.enc.Encode(&evt)
}

func (w *ASCIIWriter) Close() error {
	return w.enc.Close()
}

var (
	_ Reader = (*ASCIIReader)(nil)

	_ Writer    = (*ASCIIWriter)(nil)
	_ io.Closer = (*ASCIIWriter)(nil)
)
