// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook // import "go-hep.org/x/hep/hbook"

import (
	"bytes"
	"fmt"
	"strings"
)

// readYODAHeader parses the input buffer and extracts the YODA header line
// from that buffer.
// readYODAHeader returns the associated YODA path and an error if any.
func readYODAHeader(r *bytes.Buffer, hdr string) (string, error) {
	pos := bytes.Index(r.Bytes(), []byte("\n"))
	if pos < 0 {
		return "", fmt.Errorf("hbook: could not find %s line", hdr)
	}
	path := string(r.Next(pos + 1))
	if !strings.HasPrefix(path, hdr+" ") {
		return "", fmt.Errorf("hbook: could not find %s mark", hdr)
	}

	return path[len(hdr)+1 : len(path)-1], nil
}
