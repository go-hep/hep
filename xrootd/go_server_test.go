// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build xrootd_test_with_go_server
// +build xrootd_test_with_go_server

package client // import "go-hep.org/x/hep/xrootd/client"

func init() {
	testClientAddrs = append(testClientAddrs, "0.0.0.0:9001")
}
