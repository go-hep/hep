// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"text/template"
)

func gen_task(c Component) error {
	return gen(os.Stdout, g_task_template, c)
}

func gen_svc(c Component) error {
	return gen(os.Stdout, g_svc_template, c)
}

func gen(w io.Writer, text string, data interface{}) error {
	t := template.Must(template.New("fwk").Parse(text))
	return t.Execute(w, data)
}
