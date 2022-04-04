// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// WriteVar describes a variable to be written out to a tree.
type WriteVar struct {
	Name  string      // name of the variable
	Value interface{} // pointer to the value to write
	Count string      // name of the branch holding the count-leaf value for slices
}

// WriteVarsFromStruct creates a slice of WriteVars from the ptr value.
// WriteVarsFromStruct panics if ptr is not a pointer to a struct value.
// WriteVarsFromStruct ignores fields that are not exported.
func WriteVarsFromStruct(ptr interface{}, opts ...WriteOption) []WriteVar {
	cfg := wopt{
		splitlvl: defaultSplitLevel,
	}
	for _, opt := range opts {
		_ = opt(&cfg)
	}

	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic(fmt.Errorf("rtree: expect a pointer value, got %T", ptr))
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("rtree: expect a pointer to struct value, got %T", ptr))
	}
	var (
		reDims = regexp.MustCompile(`\w*?\[(\w*)\]+?`)
	)

	split := func(s string) (string, []string) {
		n := s
		if i := strings.Index(s, "["); i > 0 {
			n = s[:i]
		}

		out := reDims.FindAllStringSubmatch(s, -1)
		if len(out) == 0 {
			return n, nil
		}

		dims := make([]string, len(out))
		for i := range out {
			dims[i] = out[i][1]
		}
		return n, dims
	}

	if cfg.splitlvl == 0 {
		name := cfg.title
		if name == "" {
			panic(fmt.Errorf("rtree: expect a title for split-less struct"))
		}
		return []WriteVar{
			{Name: name, Value: ptr},
		}
	}

	var (
		rt    = rv.Type()
		wvars = make([]WriteVar, 0, rt.NumField())
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			ft = rt.Field(i)
			fv = rv.Field(i)
		)
		if ft.Name != toTitle(ft.Name) {
			// not exported. ignore.
			continue
		}
		wvar := WriteVar{
			Name:  ft.Tag.Get("groot"),
			Value: fv.Addr().Interface(),
		}
		if wvar.Name == "" {
			wvar.Name = ft.Name
		}

		if strings.Contains(wvar.Name, "[") {
			switch ft.Type.Kind() {
			case reflect.Slice:
				sli, dims := split(wvar.Name)
				if len(dims) > 1 {
					panic(fmt.Errorf("rtree: invalid number of slice-dimensions for field %q: %q", ft.Name, wvar.Name))
				}
				wvar.Name = sli
				wvar.Count = dims[0]

			case reflect.Array:
				arr, dims := split(wvar.Name)
				if len(dims) > 3 {
					panic(fmt.Errorf("rtree: invalid number of array-dimension for field %q: %q", ft.Name, wvar.Name))
				}
				wvar.Name = arr
			default:
				panic(fmt.Errorf("rtree: invalid field type for %q, or invalid struct-tag %q: %T", ft.Name, wvar.Name, fv.Interface()))
			}
		}
		switch ft.Type.Kind() {
		case reflect.Int, reflect.Uint, reflect.UnsafePointer, reflect.Uintptr, reflect.Chan, reflect.Interface:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T", ft.Name, fv.Interface()))
		case reflect.Map:
			panic(fmt.Errorf("rtree: invalid field type for %q: %T (not yet supported)", ft.Name, fv.Interface()))
		}

		wvars = append(wvars, wvar)
	}

	return wvars
}

// WriteVarsFromTree creates a slice of WriteVars from the tree value.
func WriteVarsFromTree(t Tree) []WriteVar {
	rvars := NewReadVars(t)
	wvars := make([]WriteVar, len(rvars))
	for i, rvar := range rvars {
		wvars[i] = WriteVar{
			Name:  rvar.Name,
			Value: reflect.New(reflect.TypeOf(rvar.Value).Elem()).Interface(),
			Count: rvar.count,
		}
	}
	return wvars
}
