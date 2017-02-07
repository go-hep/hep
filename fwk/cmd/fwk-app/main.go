package main

import (
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
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
	err = g_cmd.Dispatch(args)
	handle_err(err)
}
