package rootio

import (
	"bytes"
	"reflect"
)

type objarray []Object

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (arr *objarray) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
	return err
}

func init() {
	f := func() reflect.Value {
		o := make(objarray, 0)
		return reflect.ValueOf(o)
	}
	Factory.db["TObjArray"] = f
	Factory.db["*rootio.objarray"] = f
}

// check interfaces
//var _ Object = (*objarray)(nil)
var _ ROOTUnmarshaler = (*objarray)(nil)

// EOF
