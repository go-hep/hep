// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

// Encoder encodes data into the underlying io.Writer
type Encoder interface {
	Encode(data any) error
}

// Save saves a job's configuration description using the Encoder enc.
func Save(stmts []Stmt, enc Encoder) error {
	return enc.Encode(stmts)
}

// Decoder decodes data from the unerlying io.Reader
type Decoder interface {
	Decode(ptr any) error
}

// Load loads a job's configuration description using the Decoder dec.
func Load(dec Decoder) ([]Stmt, error) {
	stmts := make([]Stmt, 0)
	err := dec.Decode(&stmts)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}
