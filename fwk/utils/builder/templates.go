// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io"
	"strings"
	"text/template"
)

const (
	tmpl = `// template for a fwk-app main

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	_ "go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/job"
)

var (
	g_lvl      = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	g_evtmax   = flag.Int("evtmax", -1, "number of events to process")
	g_nprocs   = flag.Int("nprocs", 0, "number of concurrent events to process")
	g_cpu_prof = flag.Bool("cpu-prof", false, "enable CPU profiling")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, {{.Usage | gen_usage}}, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	fmt.Printf("::: {{.Name}}...\n")
	if *g_cpu_prof {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app := job.New(job.P{
		"EvtMax":   int64(*g_evtmax),
		"NProcs":   *g_nprocs,
		"MsgLevel": job.MsgLevel(*g_lvl),
	})

    {{with .SetupFuncs}}{{. | gen_setups}}{{end}}

	app.Run()
	fmt.Printf("::: {{.Name}}... [done]\n")
}
`
)

func render(w io.Writer, text string, data interface{}) error {
	t := template.New("fwk-main")
	t.Funcs(template.FuncMap{
		"trim":       strings.TrimSpace,
		"gen_setups": gen_setups,
		"gen_usage":  gen_usage,
	})
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

func gen_setups(setups []string) string {
	str := make([]string, 0, len(setups))
	for _, setup := range setups {
		str = append(str,
			"\t"+setup+"(app)",
		)
	}
	return strings.Join(str, "\n")
}

func gen_usage(usage string) string {
	return "`" + usage + "`"
}
