package main

const g_svc_template = `package {{.Package}}

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type {{.Name}} struct {
    fwk.SvcBase
}

func (svc *{{.Name}}) Configure(ctx fwk.Context) error {
    var err error

	// err = svc.DeclInPort(svc.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = svc.DeclOutPort(svc.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

    return err
}

func (svc *{{.Name}}) StartSvc(ctx fwk.Context) error {
	var err error

	return err
}

func (svc *{{.Name}}) StopSvc(ctx fwk.Context) error {
	var err error

	return err
}

func new{{.Name}}(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	svc := &{{.Name}}{
		SvcBase: fwk.NewSvc(typ, name, mgr),
		// input:    "Input",
		// output:   "Output",
	}

	// err = svc.DeclProp("Input", &svc.input)
	// if err != nil {
	// 	return nil, err
	// }

	// err = svc.DeclProp("Output", &svc.output)
	// if err != nil {
	//	return nil, err
	// }

	return svc, err
}

func init() {
	fwk.Register(reflect.TypeOf({{.Name}}{}), new{{.Name}})
}
`
