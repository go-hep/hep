#!/usr/bin/env bash

# Copyright 2018 The go-hep Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

echo ">>> using Go: `which go`"
go version

echo ">>> go env:"
go env

echo ">>> build from: `pwd`"

go get -d -t -v ./...
go install -v $TAGS ./...
