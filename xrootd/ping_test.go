// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkPing(t *testing.T, client *Client, done chan<- bool) {
	err := client.Ping(context.Background())
	assert.NoError(t, err)
	done <- true
}

func TestClient_Ping_100(t *testing.T) {
	client, err := New(context.Background(), *Addr)
	assert.NoError(t, err)
	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(t, err)

	count := 100
	done := make(chan bool, count)

	for i := 0; i < count; i++ {
		go checkPing(t, client, done)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}

func BenchmarkHundredPings(b *testing.B) {
	client, err := New(context.Background(), *Addr)
	assert.NoError(b, err)
	_, err = client.Login(context.Background(), "gopher")
	assert.NoError(b, err)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := client.Ping(context.Background()); err != nil {
			b.Error(err)
		}
	}
}

func ExampleClient_Ping() {
	client, _ := New(context.Background(), *Addr)

	client.Login(context.Background(), "gopher")
	client.Ping(context.Background())

	fmt.Print("Pong!")
	// Output: Pong!
}
