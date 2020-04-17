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
	var err error
	err = gen(os.Stdout, g_task_template, c)
	return err
}

func gen_svc(c Component) error {
	var err error
	err = gen(os.Stdout, g_svc_template, c)
	return err
}

func gen(w io.Writer, text string, data interface{}) error {
	t := template.New("fwk")
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
