// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

// Encoder encodes data into the underlying io.Writer
type Encoder interface {
	Encode(data interface{}) error
}

// Save saves a job's configuration description using the Encoder enc.
func Save(stmts []Stmt, enc Encoder) error {
	var err error
	err = enc.Encode(stmts)
	return err
}

// Decoder decodes data from the unerlying io.Reader
type Decoder interface {
	Decode(ptr interface{}) error
}

// Load loads a job's configuration description using the Decoder dec.
func Load(dec Decoder) ([]Stmt, error) {
	var err error
	stmts := make([]Stmt, 0)
	err = dec.Decode(&stmts)
	if err != nil {
		return nil, err
	}

	return stmts, err
}
