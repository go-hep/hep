// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"encoding/json"
	"io"
	"os"
)

// NewJSONEncoder returns a new encoder that writes to w
func NewJSONEncoder(w io.Writer) *json.Encoder {
	if w == nil {
		w = os.Stdout
	}
	return json.NewEncoder(w)
}

// NewJSONDecoder returns a new decoder that reads from r.
func NewJSONDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}
