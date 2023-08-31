// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"

	"github.com/gonuts/commander"
)

func handle_err(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	g_cmd *commander.Command
)

func init() {
	g_cmd = &commander.Command{
		UsageLine: "fwk-app <sub-command> [options] [args [...]]",
		Short:     "builds and runs fwk-based applications",
		Subcommands: []*commander.Command{
			fwk_make_cmd_run(),
			fwk_make_cmd_build(),
		},
		Flag: *flag.NewFlagSet("fwk-app", flag.ExitOnError),
	}
}

func main() {

	var err error

	err = g_cmd.Flag.Parse(os.Args[1:])
	handle_err(err)

	args := g_cmd.Flag.Args()
	err = g_cmd.Dispatch(context.Background(), args)
	handle_err(err)
}
