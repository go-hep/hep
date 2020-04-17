// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lhef

import (
	"fmt"
	"io"
	"sync"
)

// Encoder encodes a LHEF event to the underlying writer, following the
// Les Houches Event File format.
type Encoder struct {
	w      io.Writer
	once   sync.Once
	Run    HEPRUP // User process run common block
	Header []byte // header block data

}

// NewEncoder creates a new Encoder connected to the given writer.
func NewEncoder(w io.Writer) (*Encoder, error) {
	enc := &Encoder{
		w: w,
	}
	return enc, nil
}

func (e *Encoder) init() error {

	var err error
	run := &e.Run

	version := 1.0
	if run.XSecInfo.Neve > 0 {
		version = 2.0
		panic("not implemented")
	}
	_, err = fmt.Fprintf(
		e.w,
		"<LesHouchesEvents version=\"%0.1f\">\n",
		version,
	)

	if err != nil {
		return err
	}

	if len(e.Header) > 0 {
		hdr := string(e.Header)
		if hdr[len(hdr)-1] == '\n' {
			hdr = hdr[:len(hdr)-1]
		}
		_, err = fmt.Fprintf(
			e.w,
			"<header>\n%v\n</header>\n",
			hdr,
		)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(
		e.w,
		"<init>\n %7d %7d %13.6E %13.6E %5d %5d %5d %5d %5d %5d\n",
		run.IDBMUP[0], run.IDBMUP[1],
		run.EBMUP[0], run.EBMUP[1],
		run.PDFGUP[0], run.PDFGUP[1],
		run.PDFSUP[0], run.PDFSUP[1],
		run.IDWTUP, run.NPRUP,
	)
	if err != nil {
		return err
	}

	for i := 0; i < int(run.NPRUP); i++ {
		_, err = fmt.Fprintf(
			e.w,
			" %13.6E %13.6E %13.6E %5d\n",
			run.XSECUP[i],
			run.XERRUP[i],
			run.XMAXUP[i],
			run.LPRUP[i],
		)
		if err != nil {
			return err
		}
	}

	if run.XSecInfo.Neve <= 0 {
		_, err = fmt.Fprintf(
			e.w,
			"#%s\n</init>\n",
			"",
		)
		return err
	}

	return err
}

func (e *Encoder) Encode(evt *HEPEUP) error {
	var err error
	e.once.Do(func() { err = e.init() })
	if err != nil {
		return err
	}

	if len(evt.SubEvents.Events) > 0 {
		_, err = fmt.Fprintf(e.w, "<eventgroup")
		if err != nil {
			return err
		}
		if evt.SubEvents.Nreal > 0 {
			_, err = fmt.Fprintf(e.w, " nreal=\"%d\"", evt.SubEvents.Nreal)
			if err != nil {
				return err
			}
		}
		if evt.SubEvents.Ncounter > 0 {
			_, err = fmt.Fprintf(e.w, " ncounter=\"%d\"", evt.SubEvents.Ncounter)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(e.w, "</eventgroup>\n")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(
		e.w,
		"<event>\n %5d %5d %13.6E %13.6E %13.6E %13.6E\n",
		evt.NUP,
		evt.IDPRUP,
		evt.XWGTUP,
		evt.SCALUP,
		evt.AQEDUP,
		evt.AQCDUP,
	)
	if err != nil {
		return err
	}

	for i := 0; i < int(evt.NUP); i++ {
		_, err = fmt.Fprintf(
			e.w,
			" %7d %4d %4d %4d %4d %4d %17.10E %17.10E %17.10E %17.10E %17.10E %1.f. %1.f.\n",
			evt.IDUP[i],
			evt.ISTUP[i],
			evt.MOTHUP[i][0], evt.MOTHUP[i][1],
			evt.ICOLUP[i][0], evt.ICOLUP[i][1],
			evt.PUP[i][0], evt.PUP[i][1], evt.PUP[i][2], evt.PUP[i][3], evt.PUP[i][4],
			evt.VTIMUP[i],
			evt.SPINUP[i],
		)
		if err != nil {
			return err
		}
	}

	if e.Run.XSecInfo.Neve > 0 {
		panic("not implemented")
	}

	_, err = fmt.Fprintf(
		e.w,
		"#%s\n</event>\n",
		"",
	)

	return err
}

func (e *Encoder) Close() error {

	_, err := fmt.Fprintf(
		e.w,
		"</LesHouchesEvents>\n",
	)
	if err != nil {
		return err
	}

	if w, ok := e.w.(io.WriteCloser); ok {
		err = w.Close()
	}
	return err
}
