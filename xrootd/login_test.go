// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build xrootd_test_with_server

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	client, err := NewClient(context.Background(), *Addr)
	assert.NoError(t, err)

	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)
}

func ExampleClient_Login() {
	client, _ := NewClient(context.Background(), *Addr)
	loginResult, _ := client.Login(context.Background(), "gopher")
	fmt.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)
	// Output: Logged in! Security information length is 0. Value is ""
}
