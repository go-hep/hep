package main

const g_task_template = `package {{.Package}}

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type {{.Name}} struct {
    fwk.TaskBase
}

func (tsk *{{.Name}}) Configure(ctx fwk.Context) error {
    var err error

	// err = tsk.DeclInPort(tsk.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

    return err
}

func (tsk *{{.Name}}) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *{{.Name}}) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *{{.Name}}) Process(ctx fwk.Context) error {
	var err error

	return err
}

func new{{.Name}}(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &{{.Name}}{
		TaskBase: fwk.NewTask(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
	}

	// err = tsk.DeclProp("Input", &tsk.input)
	// if err != nil {
	// 	return nil, err
	// }

	// err = tsk.DeclProp("Output", &tsk.output)
	// if err != nil {
	//	return nil, err
	// }

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf({{.Name}}{}), new{{.Name}})
}
`
