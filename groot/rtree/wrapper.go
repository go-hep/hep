// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

// wobject wrapps a type created from a Streamer and implements the
// following interfaces:
//  - root.Object
//  - rbytes.Marshaler
//  - rbytes.Unmarshaler
type wobject struct {
	v         interface{}
	class     func(recv interface{}) string
	unmarshal func(recv interface{}, r *rbytes.RBuffer) error
	marshal   func(recv interface{}, w *rbytes.WBuffer) (int, error)
}

func (obj *wobject) Class() string {
	return obj.class(obj.v)
}

func (obj *wobject) UnmarshalROOT(r *rbytes.RBuffer) error {
	return obj.unmarshal(obj.v, r)
}

func (obj *wobject) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	return obj.marshal(obj.v, w)
}

var (
	_ root.Object        = (*wobject)(nil)
	_ rbytes.Marshaler   = (*wobject)(nil)
	_ rbytes.Unmarshaler = (*wobject)(nil)
)

var builtins = map[string]reflect.Type{
	"TObject":        reflect.TypeOf((*rbase.Object)(nil)).Elem(),
	"TString":        reflect.TypeOf(""),
	"TNamed":         reflect.TypeOf((*rbase.Named)(nil)).Elem(),
	"TList":          reflect.TypeOf((*rcont.List)(nil)).Elem(),
	"TObjArray":      reflect.TypeOf((*rcont.ObjArray)(nil)).Elem(),
	"TObjString":     reflect.TypeOf((*rbase.ObjString)(nil)).Elem(),
	"TTree":          reflect.TypeOf((*ttree)(nil)).Elem(),
	"TBranch":        reflect.TypeOf((*tbranch)(nil)).Elem(),
	"TBranchElement": reflect.TypeOf((*tbranchElement)(nil)).Elem(),
}
