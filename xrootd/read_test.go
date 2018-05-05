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
	"go-hep.org/x/hep/xrootd/requests/open"
)

func TestClient_Read(t *testing.T) {
	client, err := NewClient(context.Background(), *Addr)
	assert.NoError(t, err)

	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)

	handle, err := client.Open(context.Background(), "/tmp/testFiles/read", open.ModeOwnerWrite, open.OptionsMkPath|open.OptionsDelete)
	assert.NoError(t, err)
	assert.NotNil(t, handle)

	message := []byte("Hello")
	err = client.Write(context.Background(), handle, 0, 0, message)
	assert.NoError(t, err)

	err = client.Sync(context.Background(), handle)
	assert.NoError(t, err)

	err = client.Close(context.Background(), handle, int64(len(message)))
	assert.NoError(t, err)

	handle, err = client.Open(context.Background(), "/tmp/testFiles/read", open.ModeOwnerRead, open.OptionsNone)
	assert.NoError(t, err)
	assert.NotNil(t, handle)

	data, err := client.Read(context.Background(), handle, 0, 5)
	assert.NoError(t, err)
	assert.Equal(t, message, data)

	err = client.Close(context.Background(), handle, int64(len(message)))
	assert.NoError(t, err)
}

func ExampleClient_Read() {
	client, _ := NewClient(context.Background(), *Addr)
	client.Login(context.Background(), "gopher")
	handle, _ := client.Open(context.Background(), "/tmp/testFiles/read_example", open.ModeOwnerWrite, open.OptionsMkPath|open.OptionsDelete)
	message := []byte("Hello")
	client.Write(context.Background(), handle, 0, 0, message)
	client.Sync(context.Background(), handle)

	data, _ := client.Read(context.Background(), handle, 0, int32(len(message)))
	fmt.Printf("%s", data)

	client.Close(context.Background(), handle, 0)
	// Output: Hello
}
