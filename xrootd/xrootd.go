// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrootd implements the XRootD protocol from
//  http://xrootd.org
//
// Package xrootd provides a Client and a Server.
//
// The NewClient function connects to a server:
//
//	ctx := context.Background()
//
//	client, err := xrootd.NewClient(ctx, addr, username)
//	if err != nil {
//		// handle error
//	}
//
//	// ...
//
//	if err := client.Close(); err != nil {
//		// handle error
//	}
//
// The NewServer function creates a server:
//
//  srv := xrootd.NewServer(xrootd.Default(), nil)
//  err := srv.Serve(listener)
package xrootd // import "go-hep.org/x/hep/xrootd"
