package rootio

import (
	"bytes"
	"reflect"
)

type attfill struct {
	color int16
	style int16
}

func (a *attfill) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	dec := NewDecoder(data)

	start := dec.Pos()
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	}

	err = dec.readBin(&a.color)
	if err != nil {
		return err
	}

	err = dec.readBin(&a.style)
	if err != nil {
		return err
	}

	err = dec.checkByteCount(pos, bcnt, start, "TAttFill")
	return err
}

//
func init() {
	f := func() reflect.Value {
		o := &attfill{}
		return reflect.ValueOf(o)
	}
	Factory.db["TAttFill"] = f
	Factory.db["*rootio.attfill"] = f
}

// ifaces
var _ ROOTUnmarshaler = (*attfill)(nil)
