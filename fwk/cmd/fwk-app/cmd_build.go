// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"

	"github.com/gonuts/commander"
	"go-hep.org/x/hep/fwk/utils/builder"
)

func fwk_make_cmd_build() *commander.Command {
	cmd := &commander.Command{
		Run:       fwk_run_cmd_build,
		UsageLine: "build [options] <config.go> [<config2.go> [...]]",
		Short:     "build a fwk job",
		Long: `
build builds a fwk-based job and produces a binary.

ex:
 $ fwk-app build config.go
 $ fwk-app build config1.go config2.go
 $ fwk-app build ./some-dir
 $ fwk-app build -o=my-binary config.go
`,
		Flag: *flag.NewFlagSet("fwk-app-build", flag.ExitOnError),
	}
	cmd.Flag.String("o", "", "name of the resulting binary (default=name of parent directory)")
	return cmd
}

func fwk_run_cmd_build(cmd *commander.Command, args []string) error {
	var err error
	n := "fwk-app-" + cmd.Name()

	fnames := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "" {
			continue
		}
		if arg[0] == '-' {
			continue
		}
		fnames = append(fnames, arg)
	}

	if len(fnames) <= 0 {
		return fmt.Errorf("%s: you need to give a list of files or a directory", n)
	}

	bldr, err := builder.NewBuilder(fnames...)
	if err != nil {
		return err
	}

	if o := cmd.Lookup("o").(string); o != "" {
		bldr.Name = o
	}

	err = bldr.Build()
	if err != nil {
		return err
	}

	return err
}
