// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepevt

import (
	"fmt"
	"io"

	"golang.org/x/xerrors"
)

// Encoder encodes ASCII files in the HEPEVT format.
type Encoder struct {
	w io.Writer
}

// NewEncoder create a new Encoder, writing to the provided io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encodes a full HEPEVT event to the underlying writer.
func (enc *Encoder) Encode(evt *Event) error {
	_, err := fmt.Fprintf(enc.w, "%d %d\n", evt.Nevhep, evt.Nhep)
	if err != nil {
		return xerrors.Errorf("could not encode event header line: %w", err)
	}

	for i := 0; i < evt.Nhep; i++ {
		_, err = fmt.Fprintf(
			enc.w,
			"%d %d %d %d %d %d %E %E %E %E %E %E %E %E %E\n",
			evt.Isthep[i],
			evt.Idhep[i],
			// convert 0-based indices to 1-based ones
			evt.Jmohep[i][0]+1, evt.Jmohep[i][1]+1,
			evt.Jdahep[i][0]+1, evt.Jdahep[i][1]+1,
			// <---
			evt.Phep[i][0], evt.Phep[i][1], evt.Phep[i][2], evt.Phep[i][3],
			evt.Phep[i][4],
			evt.Vhep[i][0], evt.Vhep[i][1], evt.Vhep[i][2], evt.Vhep[i][3],
		)
		if err != nil {
			return xerrors.Errorf("could not encode event particle line[%d]: %w", i, err)
		}
	}
	return nil
}

// Decoder decodes ASCII files in the HEPEVT format.
type Decoder struct {
	r io.Reader
}

// NewDecoder creates a new Decoder, reading from the provided io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode decodes a full HEPEVT event from the underlying reader.
func (dec *Decoder) Decode(evt *Event) error {
	_, err := fmt.Fscanf(dec.r, "%d %d\n", &evt.Nevhep, &evt.Nhep)
	if err != nil {
		return xerrors.Errorf("could not decode event header line: %w", err)
	}

	// resize
	if len(evt.Isthep) > evt.Nhep {
		evt.Isthep = evt.Isthep[:evt.Nhep]
		evt.Idhep = evt.Idhep[:evt.Nhep]
		evt.Jmohep = evt.Jmohep[:evt.Nhep]
		evt.Jdahep = evt.Jdahep[:evt.Nhep]
		evt.Phep = evt.Phep[:evt.Nhep]
		evt.Vhep = evt.Vhep[:evt.Nhep]
	} else {
		sz := evt.Nhep - len(evt.Isthep)
		evt.Isthep = append(evt.Isthep, make([]int, sz)...)
		evt.Idhep = append(evt.Idhep, make([]int, sz)...)
		evt.Jmohep = append(evt.Jmohep, make([][2]int, sz)...)
		evt.Jdahep = append(evt.Jdahep, make([][2]int, sz)...)
		evt.Phep = append(evt.Phep, make([][5]float64, sz)...)
		evt.Vhep = append(evt.Vhep, make([][4]float64, sz)...)
	}

	for i := 0; i < evt.Nhep; i++ {
		_, err = fmt.Fscanf(
			dec.r,
			"%d %d %d %d %d %d %E %E %E %E %E %E %E %E %E\n",
			&evt.Isthep[i],
			&evt.Idhep[i],
			&evt.Jmohep[i][0], &evt.Jmohep[i][1],
			&evt.Jdahep[i][0], &evt.Jdahep[i][1],
			&evt.Phep[i][0], &evt.Phep[i][1], &evt.Phep[i][2], &evt.Phep[i][3],
			&evt.Phep[i][4],
			&evt.Vhep[i][0], &evt.Vhep[i][1], &evt.Vhep[i][2], &evt.Vhep[i][3],
		)
		if err != nil {
			return xerrors.Errorf("could not decode event line[%d]: %w", i, err)
		}
		// convert 0-based indices to 1-based ones
		evt.Jmohep[i][0] -= 1
		evt.Jmohep[i][1] -= 1
		evt.Jdahep[i][0] -= 1
		evt.Jdahep[i][1] -= 1
	}
	return nil
}
