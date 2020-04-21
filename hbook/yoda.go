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
func readYODAHeader(r *rbuffer, hdr string) (string, int, error) {
	pos := bytes.Index(r.Bytes(), []byte("\n"))
	if pos < 0 {
		return "", 0, fmt.Errorf("hbook: could not find %s line", hdr)
	}
	var (
		path = string(r.next(pos + 1))
		vers int
	)
	switch {
	case strings.HasPrefix(path, hdr+"_V2 "):
		hdr += "_V2"
		vers = 2
	case strings.HasPrefix(path, hdr+" "):
		vers = 1
	default:
		return "", 0, fmt.Errorf("hbook: could not find %s mark", hdr)
	}

	return path[len(hdr)+1 : len(path)-1], vers, nil
}
