// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lhef

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Decoder represents an LHEF parser reading a particular input stream.
//
// A Decoder is initialized with an input io.Reader from which to read a version 1.0
// Les Houches Accord event file.
type Decoder struct {
	r       io.Reader
	dec     *xml.Decoder
	evt     xml.StartElement // the current xml.Token holding a HEPEUP
	Version int              // LHEF file version
	Run     HEPRUP           // User process run common block
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	var err error
	dec := xml.NewDecoder(r)
	d := &Decoder{
		r:       r,
		dec:     dec,
		Version: -42,
		Run:     HEPRUP{},
	}
	// err := dec.dec.Decode(&dec.lhevt)
	// if err != nil {
	// 	return nil
	// }
	// fmt.Printf(">>> version=%v\n", dec.lhevt.Version)

	tok, err := dec.Token()
	if err != nil || tok == nil {
		return nil, err
	}

	// make sure we are reading a LHEF file
	start, ok := tok.(xml.StartElement)
	if !ok || start.Name.Local != "LesHouchesEvents" {
		return nil, fmt.Errorf("lhef.Decoder: missing LesHouchesEvent start-tag")
	}

	//fmt.Printf(">>> %v\n", start)
	version := start.Attr[0].Value
	//fmt.Printf("    version=[%s]\n", version)
	switch version {
	default:
		d.Version = -42
	case "1.0":
		d.Version = 1
	case "2.0":
		d.Version = 2
	}

	var (
		init   xml.StartElement
		header xml.StartElement
	)

Loop:
	for {
		tok, err = dec.Token()
		if err != nil || tok == nil {
			return nil, err
		}
		switch tt := tok.(type) {
		case xml.Comment:
			// skip
		case xml.CharData:
			// skip
		case xml.StartElement:
			switch tt.Name.Local {
			case "init":
				init = tt
				break Loop
			case "header":
				header = tt //FIXME
				panic(fmt.Errorf("header not implemented: %v", header))
			}
		}
	}
	if init.Name.Local != "init" {
		return nil, fmt.Errorf("lhef.Decoder: missing init start-tag")
	}

	// extract compulsory initialization information
	tok, err = dec.Token()
	if err != nil {
		return nil, err
	}
	data, ok := tok.(xml.CharData)
	if !ok || len(data) == 0 {
		return nil, fmt.Errorf("lhef.Decoder: missing init payload")
	}
	buf := bytes.NewBuffer(data)
	_, err = fmt.Fscanf(
		buf,
		"\n%d %d %f %f %d %d %d %d %d %d\n",
		&d.Run.IDBMUP[0],
		&d.Run.IDBMUP[1],
		&d.Run.EBMUP[0],
		&d.Run.EBMUP[1],
		&d.Run.PDFGUP[0],
		&d.Run.PDFGUP[1],
		&d.Run.PDFSUP[0],
		&d.Run.PDFSUP[1],
		&d.Run.IDWTUP,
		&d.Run.NPRUP,
	)
	if err != nil {
		return nil, err
	}

	d.Run.XSECUP = make([]float64, int(d.Run.NPRUP))
	d.Run.XERRUP = make([]float64, int(d.Run.NPRUP))
	d.Run.XMAXUP = make([]float64, int(d.Run.NPRUP))
	d.Run.LPRUP = make([]int32, int(d.Run.NPRUP))

	for i := 0; i < int(d.Run.NPRUP); i++ {
		_, err = fmt.Fscanf(
			buf,
			"%f %f %f %d\n",
			&d.Run.XSECUP[i],
			&d.Run.XERRUP[i],
			&d.Run.XMAXUP[i],
			&d.Run.LPRUP[i],
		)
		if err != nil {
			return nil, errors.Errorf("lhef: failed to decode NPRUP %d: %v", i, err)
		}
	}

	if d.Version >= 2 {
		// do version-2 specific stuff
	}

	tok, err = dec.Token()
	if err != nil {
		return nil, err
	}
	if end, ok := tok.(xml.EndElement); !ok || end.Name.Local != "init" {
		return nil, fmt.Errorf("lhef.Decoder: missing init end-tag")
	}

	return d, nil
}

// advance to next event
func (d *Decoder) next() error {
LoopEvt:
	for {
		tok, err := d.dec.Token()
		if err != nil {
			return err
		}
		switch tt := tok.(type) {
		case xml.Comment:
			// skip
		case xml.CharData:
			// skip
		case xml.StartElement:
			switch tt.Name.Local {
			case "event":
				d.evt = tt
				break LoopEvt
			}
		}
	}

	return nil
}

// Read an event from the file
func (d *Decoder) Decode() (*HEPEUP, error) {

	// check whether the initialization was successful
	if d.Run.NPRUP < 0 {
		return nil, fmt.Errorf("lhef.Decode: initialization failed (no particle entries)")
	}

	err := d.next()
	if err != nil {
		return nil, err
	}

	// extract payload data
	tok, err := d.dec.Token()
	if err != nil {
		return nil, err
	}
	data, ok := tok.(xml.CharData)
	if !ok {
		return nil, fmt.Errorf("lhef.Decode: invalid token (%T: %v)", tok, tok)
	}
	if len(data) <= 0 {
		return nil, fmt.Errorf("lhef.Decode: empty event")
	}
	buf := bytes.NewBuffer(data)

	evt := &HEPEUP{
		NUP: 0,
	}

	_, err = fmt.Fscanf(
		buf,
		"\n%d %d %f %f %f %f\n",
		&evt.NUP,
		&evt.IDPRUP,
		&evt.XWGTUP,
		&evt.SCALUP,
		&evt.AQEDUP,
		&evt.AQCDUP,
	)
	if err != nil {
		fmt.Printf("--> 2 err: %v\n", err)
		fmt.Printf("  data:    %v\n", string(data))
		fmt.Printf("  token:   (%T: %v)\n", tok, tok)
		return nil, err
	}

	n := int(evt.NUP)
	evt.IDUP = make([]int64, n)
	evt.ISTUP = make([]int32, n)
	evt.MOTHUP = make([][2]int32, n)
	evt.ICOLUP = make([][2]int32, n)
	evt.PUP = make([][5]float64, n)
	evt.VTIMUP = make([]float64, n)
	evt.SPINUP = make([]float64, n)
	for i := 0; i < n; i++ {
		_, err = fmt.Fscanf(
			buf,
			"%d %d %d %d %d %d %f %f %f %f %f %f %f\n",
			&evt.IDUP[i],
			&evt.ISTUP[i],
			&evt.MOTHUP[i][0],
			&evt.MOTHUP[i][1],
			&evt.ICOLUP[i][0],
			&evt.ICOLUP[i][1],
			&evt.PUP[i][0], &evt.PUP[i][1], &evt.PUP[i][2],
			&evt.PUP[i][3], &evt.PUP[i][4],
			&evt.VTIMUP[i],
			&evt.SPINUP[i],
		)
		if err != nil {
			fmt.Printf("--> 3-%d err: %v\n", i, err)
			return nil, err
		}
	}

	// read any additional comments...
	_ /*evtComments*/ = buf.Bytes()

	// do

	// put "cursor" to next event...
	tok, err = d.dec.Token()
	if err != nil {
		return nil, err
	}
	if end, ok := tok.(xml.EndElement); !ok || end.Name.Local != "event" {
		return nil, fmt.Errorf("lhef.Decoder: missing event end-tag")
	}

	return evt, nil
}
