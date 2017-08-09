// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.8

package server

import "net/url"

// PathEscape escapes the string so it can be safely placed
// inside a URL path segment.
func urlPathEscape(s string) string {
	return url.PathEscape(s)
}

// PathUnescape does the inverse transformation of PathEscape, converting
// %AB into the byte 0xAB. It returns an error if any % is not followed by
// two hexadecimal digits.
//
// PathUnescape is identical to QueryUnescape except that it does not unescape '+' to ' ' (space).
func urlPathUnescape(s string) (string, error) {
	return url.PathUnescape(s)
}
