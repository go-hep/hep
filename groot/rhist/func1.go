// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// F1 is a ROOT 1-dim function.
type F1 struct {
	named   rbase.Named
	attline rbase.AttLine
	attfill rbase.AttFill
	attmark rbase.AttMarker

	xmin   float64 // Lower bounds for the range
	xmax   float64 // Upper bounds for the range
	npar   int32   // Number of parameters
	ndim   int32   // Function dimension
	npx    int32   // Number of points used for the graphical representation
	typ    int32
	npfits int32   // Number of points used in the fit
	ndf    int32   // Number of degrees of freedom in the fit
	chi2   float64 // Function fit chisquare
	fmin   float64 // Minimum value for plotting
	fmax   float64 // Maximum value for plotting

	parErrs []float64 // Array of errors of the fNpar parameters
	parMin  []float64 // Array of lower limits of the fNpar parameters
	parMax  []float64 // Array of upper limits of the fNpar parameters
	save    []float64 // Array of fNsave function values

	normalized   bool    // Normalization option (false by default)
	normIntegral float64 // Integral of the function before being normalized

	formula *Formula // Pointer to TFormula in case when user define formula

	params *F1Parameters // Pointer to Function parameters object (exists only for not-formula functions)
	compos F1Composition // saved pointer (unique_ptr is transient)

}

func newF1() *F1 {
	return &F1{
		named:   *rbase.NewNamed("", ""),
		attline: *rbase.NewAttLine(),
		attfill: *rbase.NewAttFill(),
		attmark: *rbase.NewAttMarker(),
	}
}

func (*F1) RVersion() int16 {
	return rvers.F1
}

func (*F1) Class() string {
	return "TF1"
}

// Name returns the name of the instance
func (f *F1) Name() string {
	return f.named.Name()
}

// Title returns the title of the instance
func (f *F1) Title() string {
	return f.named.Title()
}

// MarshalROOT implements rbytes.Marshaler
func (f *F1) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(f.RVersion())
	w.WriteObject(&f.named)
	w.WriteObject(&f.attline)
	w.WriteObject(&f.attfill)
	w.WriteObject(&f.attmark)

	w.WriteF64(f.xmin)
	w.WriteF64(f.xmax)
	w.WriteI32(f.npar)
	w.WriteI32(f.ndim)
	w.WriteI32(f.npx)
	w.WriteI32(f.typ)
	w.WriteI32(f.npfits)
	w.WriteI32(f.ndf)
	w.WriteF64(f.chi2)
	w.WriteF64(f.fmin)
	w.WriteF64(f.fmax)

	w.WriteStdVectorF64(f.parErrs)
	w.WriteStdVectorF64(f.parMin)
	w.WriteStdVectorF64(f.parMax)
	w.WriteStdVectorF64(f.save)

	w.WriteBool(f.normalized)
	w.WriteF64(f.normIntegral)

	w.WriteObjectAny(f.formula)
	w.WriteObjectAny(f.params)
	w.WriteObjectAny(f.compos)

	return w.SetByteCount(pos, f.Class())
}

func (f *F1) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.F1 {
		panic(fmt.Errorf("rhist: invalid TF1 version=%d > %d", vers, rvers.F1))
	}

	if vers < 10 {
		// tested with v10.
		panic(fmt.Errorf("rhist: invalid TF1 version=%d < 10", vers))
	}

	r.ReadObject(&f.named)
	r.ReadObject(&f.attline)
	r.ReadObject(&f.attfill)
	r.ReadObject(&f.attmark)

	f.xmin = r.ReadF64()
	f.xmax = r.ReadF64()
	f.npar = r.ReadI32()
	f.ndim = r.ReadI32()
	f.npx = r.ReadI32()
	f.typ = r.ReadI32()
	f.npfits = r.ReadI32()
	f.ndf = r.ReadI32()
	f.chi2 = r.ReadF64()
	f.fmin = r.ReadF64()
	f.fmax = r.ReadF64()

	r.ReadStdVectorF64(&f.parErrs)
	r.ReadStdVectorF64(&f.parMin)
	r.ReadStdVectorF64(&f.parMax)
	r.ReadStdVectorF64(&f.save)

	f.normalized = r.ReadBool()
	f.normIntegral = r.ReadF64()

	if obj := r.ReadObjectAny(); obj != nil {
		f.formula = obj.(*Formula)
	}

	if obj := r.ReadObjectAny(); obj != nil {
		f.params = obj.(*F1Parameters)
	}

	if obj := r.ReadObjectAny(); obj != nil {
		f.compos = obj.(F1Composition)
	}

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

func (f *F1) String() string {
	switch {
	case f.formula != nil:
		return fmt.Sprintf("TF1{Formula: %v}", f.formula)
	case f.params != nil:
		return fmt.Sprintf("TF1{Params: %v}", f.params)
	default:
		return "TF1{...}"
	}
}

type F1Parameters struct {
	params []float64
	names  []string
}

func (*F1Parameters) RVersion() int16 {
	return rvers.F1Parameters
}

func (*F1Parameters) Class() string {
	return "TF1Parameters"
}

func (f *F1Parameters) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.F1Parameters {
		panic(fmt.Errorf("rhist: invalid TF1Parameters version=%d > %d", vers, rvers.F1Parameters))
	}

	if vers < 1 {
		// tested with v1.
		panic(fmt.Errorf("rhist: invalid TF1Parameters version=%d < 1", vers))
	}

	r.ReadStdVectorF64(&f.params)
	r.ReadStdVectorStrs(&f.names)

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

func (f *F1Parameters) String() string {
	return fmt.Sprintf(
		"TF1Parameters{Values: %v, Names: %v}",
		f.params,
		f.names,
	)
}

type f1Composition struct {
	base rbase.Object
}

func (*f1Composition) RVersion() int16 {
	return rvers.F1AbsComposition
}

func (*f1Composition) Class() string {
	return "TF1AbsComposition"
}

func (*f1Composition) isF1Composition() {}

func (f *f1Composition) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.F1AbsComposition {
		panic(fmt.Errorf("rhist: invalid TF1AbsComposition version=%d > %d", vers, rvers.F1AbsComposition))
	}

	if vers < 1 {
		// tested with v1.
		panic(fmt.Errorf("rhist: invalid TF1AbsComposition version=%d < 1", vers))
	}

	r.ReadObject(&f.base)

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

// F1Convolution is a ROOT composition function describing a
// convolution between two TF1 ROOT functions.
type F1Convolution struct {
	base f1Composition

	func1 F1 // First function to be convolved
	func2 F1 // Second function to be convolved

	params1  []float64
	params2  []float64
	parNames []string // Parameters' names

	xmin     float64 // Minimal bound of the range of the convolution
	xmax     float64 // Maximal bound of the range of the convolution
	nParams1 int32
	nParams2 int32
	cstIndex int32 // Index of the constant parameter f the first function
	nPoints  int32 // Number of point for FFT array
	flagFFT  bool  // Choose FFT or numerical convolution
}

func (*F1Convolution) RVersion() int16 {
	return rvers.F1Convolution
}

func (*F1Convolution) Class() string {
	return "TF1Convolution"
}

func (*F1Convolution) isF1Composition() {}

func (f *F1Convolution) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.F1Convolution {
		panic(fmt.Errorf("rhist: invalid TF1Convolution version=%d > %d", vers, rvers.F1Convolution))
	}

	if vers < 1 {
		// tested with v1.
		panic(fmt.Errorf("rhist: invalid TF1Convolution version=%d < 1", vers))
	}

	r.ReadObject(&f.base)

	for i, v := range []*F1{&f.func1, &f.func2} {
		obj := r.ReadObjectAny()
		if obj == nil {
			r.SetErr(fmt.Errorf("rhist: could not read fFunction%d TF1 of TF1Convolution", i+1))
			return r.Err()
		}
		*v = *(obj.(*F1))
	}

	r.ReadStdVectorF64(&f.params1)
	r.ReadStdVectorF64(&f.params2)
	r.ReadStdVectorStrs(&f.parNames)

	f.xmin = r.ReadF64()
	f.xmax = r.ReadF64()
	f.nParams1 = r.ReadI32()
	f.nParams2 = r.ReadI32()
	f.cstIndex = r.ReadI32()
	f.nPoints = r.ReadI32()
	f.flagFFT = r.ReadBool()

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

func (f *F1Convolution) String() string {
	return fmt.Sprintf(
		"TF1Convolution{Func1: %v, Func2: %v}",
		f.func1.String(),
		f.func2.String(),
	)
}

// F1NormSum is a ROOT composition function describing the linear
// combination of two TF1 ROOT functions.
type F1NormSum struct {
	base f1Composition

	nFuncs uint32  // Number of functions to add
	scale  float64 // Fixed Scale parameter to normalize function (e.g. bin width)
	xmin   float64 // Minimal bound of range of NormSum
	xmax   float64 // Maximal bound of range of NormSum

	funcs      []*F1     // Vector of size fNOfFunctions containing TF1 functions
	coeffs     []float64 // Vector of size fNOfFunctions containing coefficients in front of each function
	cstIndices []int32   // Vector with size of fNOfFunctions containing the index of the constant parameter/ function (the removed ones)
	parNames   []string  // Parameter names
}

func (*F1NormSum) RVersion() int16 {
	return rvers.F1NormSum
}

func (*F1NormSum) Class() string {
	return "TF1NormSum"
}

func (*F1NormSum) isF1Composition() {}

func (f *F1NormSum) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(f.Class())
	if vers > rvers.F1NormSum {
		panic(fmt.Errorf("rhist: invalid TF1NormSum version=%d > %d", vers, rvers.F1NormSum))
	}

	if vers < 1 {
		// tested with v1.
		panic(fmt.Errorf("rhist: invalid TF1NormSum version=%d < 1", vers))
	}

	r.ReadObject(&f.base)

	f.nFuncs = r.ReadU32()
	f.scale = r.ReadF64()
	f.xmin = r.ReadF64()
	f.xmax = r.ReadF64()

	readF1s(r, &f.funcs)
	r.ReadStdVectorF64(&f.coeffs)
	r.ReadStdVectorI32(&f.cstIndices)
	r.ReadStdVectorStrs(&f.parNames)

	r.CheckByteCount(pos, bcnt, beg, f.Class())
	return r.Err()
}

func (f *F1NormSum) String() string {
	o := new(strings.Builder)
	o.WriteString("TF1Convolution{Funcs: []{")
	for i, fct := range f.funcs {
		if i > 0 {
			o.WriteString(", ")
		}
		o.WriteString(fct.String())
	}
	o.WriteString("}, Coeffs: ")
	fmt.Fprintf(o, "%v", f.coeffs)
	o.WriteString("}")
	return o.String()
}

func readF1s(r *rbytes.RBuffer, sli *[]*F1) {
	if r.Err() != nil {
		return
	}
	const typename = "vector<TF1*>"
	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(typename)
	if vers != rvers.StreamerInfo {
		r.SetErr(fmt.Errorf(
			"rbytes: invalid %s version: got=%d, want=%d",
			typename, vers, rvers.StreamerInfo,
		))
		return
	}
	n := int(r.ReadI32())
	{
		if m := cap(*sli); m < n {
			*sli = (*sli)[:m]
			*sli = append(*sli, make([]*F1, n-m)...)
		}
		*sli = (*sli)[:n]

	}
	for i := range *sli {
		obj := r.ReadObjectAny()
		if obj == nil {
			(*sli)[i] = nil
			continue
		}
		(*sli)[i] = obj.(*F1)
	}
	r.CheckByteCount(pos, bcnt, beg, typename)
}

func init() {
	{
		f := func() reflect.Value {
			o := newF1()
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TF1", f)
	}
	{
		f := func() reflect.Value {
			o := new(F1Parameters)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TF1Parameters", f)
	}
	{
		f := func() reflect.Value {
			o := new(f1Composition)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TF1AbsComposition", f)
	}
	{
		f := func() reflect.Value {
			o := new(F1Convolution)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TF1Convolution", f)
	}
	{
		f := func() reflect.Value {
			o := new(F1NormSum)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TF1NormSum", f)
	}
}

var (
	_ root.Object        = (*F1)(nil)
	_ root.Named         = (*F1)(nil)
	_ rbytes.Marshaler   = (*F1)(nil)
	_ rbytes.Unmarshaler = (*F1)(nil)

	_ root.Object        = (*F1Parameters)(nil)
	_ rbytes.Unmarshaler = (*F1Parameters)(nil)

	_ root.Object        = (*f1Composition)(nil)
	_ rbytes.Unmarshaler = (*f1Composition)(nil)
	_ F1Composition      = (*f1Composition)(nil)

	_ root.Object        = (*F1Convolution)(nil)
	_ rbytes.Unmarshaler = (*F1Convolution)(nil)
	_ F1Composition      = (*F1Convolution)(nil)

	_ root.Object        = (*F1NormSum)(nil)
	_ rbytes.Unmarshaler = (*F1NormSum)(nil)
	_ F1Composition      = (*F1NormSum)(nil)
)
