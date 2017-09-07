// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"go-hep.org/x/hep/fwk/utils/builder"
)

func fwk_make_cmd_run() *commander.Command {
	cmd := &commander.Command{
		Run:       fwk_run_cmd_run,
		UsageLine: "run [options] <config.go> [<config2.go> [...]]",
		Short:     "run a fwk job",
		Long: `
run runs a fwk-based job.

ex:
 $ fwk-app run config.go
 $ fwk-app run config1.go config2.go
 $ fwk-app run ./some-dir
 $ fwk-app run -l=INFO -nprocs=4 -evtmax=-1 config.go
`,
		Flag: *flag.NewFlagSet("fwk-app-run", flag.ExitOnError),
	}

	cmd.Flag.String("o", "", "name of the resulting binary (default=name of parent directory)")
	cmd.Flag.Bool("k", false, "whether to keep the resulting binary after a successful run")

	// flags passed to sub-process
	cmd.Flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	cmd.Flag.Int("evtmax", -1, "number of events to process")
	cmd.Flag.Int("nprocs", 0, "number of concurrent events to process")
	cmd.Flag.Bool("cpu-prof", false, "enable CPU profiling")
	return cmd
}

func fwk_run_cmd_run(cmd *commander.Command, args []string) error {
	var err error

	n := "fwk-app-" + cmd.Name()

	subargs := make([]string, 0, len(args))
	for _, nn := range []string{"l", "evtmax", "nprocs", "cpu-prof"} {
		val := cmd.Flag.Lookup(nn)
		if val == nil {
			continue
		}
		subargs = append(
			subargs,
			fmt.Sprintf("-"+nn+"=%v", val.Value.Get()),
		)
	}

	fnames := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "" || arg == "--" {
			continue
		}

		if arg[0] == '-' {
			subargs = append(subargs, arg)
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

	if o := cmd.Flag.Lookup("o").Value.Get().(string); o != "" {
		bldr.Name = o
	}

	err = bldr.Build()
	if err != nil {
		return err
	}

	bin := bldr.Name
	if !filepath.IsAbs(bin) {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		bin = filepath.Join(pwd, bin)
	}

	if !cmd.Flag.Lookup("k").Value.Get().(bool) {
		defer os.Remove(bin)
	}

	sub := exec.Command(bin, subargs...)
	sub.Stdout = os.Stdout
	sub.Stderr = os.Stderr
	sub.Stdin = os.Stdin

	err = sub.Run()

	return err
}
