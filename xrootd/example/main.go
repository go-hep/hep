// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/xrootd/example"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go-hep.org/x/hep/xrootd"
	"go-hep.org/x/hep/xrootd/requests/open"
)

var addr = flag.String("addr", "0.0.0.0:9001", "address of xrootd server")

func main() {
	flag.Parse()

	client, err := xrootd.NewClient(context.Background(), *addr)
	checkError(err)

	response, securityInfo, err := client.Protocol(context.Background())
	checkError(err)
	log.Printf("Protocol binary version is %d. Security level is %d.", response.BinaryProtocolVersion, securityInfo.SecurityLevel)

	loginResult, err := client.Login(context.Background(), "gopher")
	checkError(err)
	log.Printf("Logged in! Security information length is %d. Value is \"%s\"\n", len(loginResult.SecurityInformation), loginResult.SecurityInformation)
	log.Printf("Session id is %x\n", loginResult.SessionID)

	err = client.Ping(context.Background())
	checkError(err)
	log.Print("Pong!")

	dirs, err := client.Dirlist(context.Background(), "/tmp/")
	checkError(err)
	log.Printf("dirlist /tmp: %s", dirs)

	fileHandle, err := client.Open(context.Background(), "/tmp/test", open.ModeOwnerWrite, open.OptionsOpenAppend|open.OptionsOpenUpdate)
	checkError(err)
	log.Printf("Open /tmp/test... File handle: %x", fileHandle)

	err = client.Write(context.Background(), fileHandle, 0, 0, []byte("Works! Hello from ematirov!"))
	checkError(err)
	log.Print("Wrote!")

	err = client.Sync(context.Background(), fileHandle)
	checkError(err)

	data, err := client.Read(context.Background(), fileHandle, 0, 27)
	checkError(err)
	log.Printf("Read /tmp/test... Content: %s", data)

	err = client.Close(context.Background(), fileHandle, 0)
	checkError(err)

	stat, err := client.Stat(context.Background(), "/tmp/test")
	checkError(err)
	log.Printf("Stat /tmp/test... Id: %d, Size: %d, Flags: %d, Modification Time: %s", stat.ID, stat.Size, stat.Flags, time.Unix(stat.ModificationTime, 0))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
