// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook // import "go-hep.org/x/hep/hbook"

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

//go:generate brio-gen -p go-hep.org/x/hep/hbook -t dist0D,dist1D,dist2D -o dist_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t Range,binning1D,binningP1D,Bin1D,BinP1D,binning2D,Bin2D -o binning_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t Point2D -o points_brio.go
//go:generate brio-gen -p go-hep.org/x/hep/hbook -t H1D,H2D,P1D,S2D -o hbook_brio.go

// Bin models 1D, 2D, ... bins.
type Bin interface {
	Rank() int           // Number of dimensions of the bin
	Entries() int64      // Number of entries in the bin
	EffEntries() float64 // Effective number of entries in the bin
	SumW() float64       // sum of weights
	SumW2() float64      // sum of squared weights
}

// Range is a 1-dim interval [Min, Max].
type Range struct {
	Min float64
	Max float64
}

// Width returns the size of the range.
func (r Range) Width() float64 {
	return math.Abs(r.Max - r.Min)
}

// Annotation is a bag of attributes that are attached to a histogram.
type Annotation map[string]interface{}

// Histogram is an n-dim histogram (with weighted entries)
type Histogram interface {
	// Annotation returns the annotations attached to the
	// histogram. (e.g. name, title, ...)
	Annotation() Annotation

	// Name returns the name of this histogram
	Name() string

	// Rank returns the number of dimensions of this histogram.
	Rank() int

	// Entries returns the number of entries of this histogram.
	Entries() int64
}

// MarshalYODA implements the YODAMarshaler interface.
func (ann Annotation) MarshalYODA() ([]byte, error) {
	keys := make([]string, 0, len(ann))
	for k := range ann {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := new(bytes.Buffer)
	for _, k := range keys {
		fmt.Fprintf(buf, "%s=%v\n", k, ann[k])
	}
	return buf.Bytes(), nil
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (ann *Annotation) UnmarshalYODA(data []byte) error {
	var err error
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		txt := s.Text()
		i := strings.Index(txt, "=")
		k := txt[:i]
		v := txt[i+1:]
		(*ann)[k] = v
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (ann *Annotation) MarshalBinary() ([]byte, error) {
	var v map[string]interface{} = *ann
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(v)
	return buf.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (ann *Annotation) UnmarshalBinary(data []byte) error {
	var v = make(map[string]interface{})
	buf := bytes.NewReader(data)
	err := gob.NewDecoder(buf).Decode(&v)
	if err != nil {
		return err
	}
	*ann = Annotation(v)
	return nil
}
