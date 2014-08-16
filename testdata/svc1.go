package testdata

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type svc1 struct {
	fwk.SvcBase
}

func (svc *svc1) Configure(ctx fwk.Context) error {
	var err error

	return err
}

func (svc *svc1) StartSvc(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("-- start svc --\n")
	return err
}

func (svc *svc1) StopSvc(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("-- stop svc --\n")
	return err
}

func newsvc1(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	svc := &svc1{
		SvcBase: fwk.NewSvc(typ, name, mgr),
	}

	return svc, err
}

func init() {
	fwk.Register(reflect.TypeOf(svc1{}), newsvc1)
}
