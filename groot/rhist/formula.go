// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// Formula describes a ROOT TFormula.
type Formula struct {
	named rbase.Named

	clingParams  []float64 // parameter values
	allParamsSet bool      // flag to control if all parameters are set

	params      map[string]int32 // list of parameter names
	formula     string           // string representing the formula expression
	ndim        int32            // Dimension - needed for lambda expressions
	linearParts []root.Object    // vector of linear functions
	vectorized  bool             // whether we should use vectorized or regular variables
}

func newFormula() *Formula {
	return &Formula{
		named: *rbase.NewNamed("", ""),
	}
}

func (*Formula) RVersion() int16 {
	return rvers.Formula
}

func (*Formula) Class() string {
	return "TFormula"
}

// Name returns the name of the instance
func (f *Formula) Name() string {
	return f.named.Name()
}

// Title returns the title of the instance
func (f *Formula) Title() string {
	return f.named.Title()
}

func (f *Formula) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.Formula {
		panic(fmt.Errorf("rhist: invalid TFormula version=%d > %d", vers, rvers.Formula))
	}

	if vers < 12 || vers > 13 {
		// tested with v12 and v13
		panic(fmt.Errorf("rhist: too old TFormula version=%d < 12", vers))
	}

	for _, v := range []rbytes.Unmarshaler{
		&f.named,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	r.ReadStdVectorF64(&f.clingParams)
	f.allParamsSet = r.ReadBool()
	f.params = func() map[string]int32 {
		if r.Err() != nil {
			return nil
		}
		const typename = "map<TString,int,TFormulaParamOrder>"
		beg := r.Pos()
		vers, pos, bcnt := r.ReadVersion(typename)
		if vers != rvers.StreamerInfo {
			r.SetErr(fmt.Errorf("rbytes: invalid %s version: got=%d, want=%d",
				typename, vers, rvers.StreamerInfo,
			))
			return nil
		}
		n := int(r.ReadI32())
		o := make(map[string]int32, n)
		for i := 0; i < n; i++ {
			k := r.ReadString()
			v := r.ReadI32()
			o[k] = v
		}
		r.CheckByteCount(pos, bcnt, beg, typename)
		return o
	}()
	f.formula = r.ReadString()
	f.ndim = r.ReadI32()

	f.linearParts = func() []root.Object {
		if r.Err() != nil {
			return nil
		}
		const typename = "vector<TObject*>"
		beg := r.Pos()
		vers, pos, bcnt := r.ReadVersion(typename)
		if vers != rvers.StreamerInfo {
			r.SetErr(fmt.Errorf("rbytes: invalid %s version: got=%d, want=%d",
				typename, vers, rvers.StreamerInfo,
			))
			return nil
		}
		n := int(r.ReadI32())
		o := make([]root.Object, n)
		for i := range o {
			o[i] = r.ReadObjectAny()
		}
		r.CheckByteCount(pos, bcnt, beg, typename)
		return o
	}()
	f.vectorized = r.ReadBool()

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

func (f *Formula) String() string {
	return fmt.Sprintf("TFormula{%s}", f.formula)
}

func init() {
	{
		f := func() reflect.Value {
			o := newFormula()
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TFormula", f)
	}
}

var (
	_ root.Object        = (*Formula)(nil)
	_ root.Named         = (*Formula)(nil)
	_ rbytes.Unmarshaler = (*Formula)(nil)
)
