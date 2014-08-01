package main

import (
	"io"
	"os"
	"text/template"
)

func gen_task(c Component) error {
	var err error
	err = gen(os.Stderr, g_task_template, c)
	return err
}

func gen_svc(c Component) error {
	var err error
	err = gen(os.Stderr, g_svc_template, c)
	return err
}

func gen(w io.Writer, text string, data interface{}) error {
	t := template.New("fwk")
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
