// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go18

package server

import "net/url"

// PathEscape escapes the string so it can be safely placed
// inside a URL path segment.
func urlPathEscape(s string) string {
	return url.PathEscape(s)
}
