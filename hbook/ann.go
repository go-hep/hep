// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Annotation is a bag of attributes that are attached to a histogram.
type Annotation map[string]interface{}

func (ann Annotation) clone() Annotation {
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(&ann)
	if err != nil {
		panic(err)
	}
	out := make(Annotation, len(ann))
	err = gob.NewDecoder(buf).Decode(&out)
	if err != nil {
		panic(err)
	}
	return out
}

// MarshalYODA implements the YODAMarshaler interface.
func (ann Annotation) MarshalYODA() ([]byte, error) {
	return ann.marshalYODAv2()
}

// marshalYODAv1 implements the YODAMarshaler interface.
func (ann Annotation) marshalYODAv1() ([]byte, error) {
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

// marshalYODAv2 implements the YODAMarshaler interface.
func (ann Annotation) marshalYODAv2() ([]byte, error) {
	return yaml.Marshal(ann)
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (ann *Annotation) UnmarshalYODA(data []byte) error {
	return ann.unmarshalYODAv2(data)
}

// unmarshalYODAv1 unmarshal YODA v1.
func (ann *Annotation) unmarshalYODAv1(data []byte) error {
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

// unmarshalYODAv2 unmarshal YODA v2.
func (ann *Annotation) unmarshalYODAv2(data []byte) error {
	return yaml.Unmarshal(data, ann)
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
