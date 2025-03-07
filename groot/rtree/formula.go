// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"

	"go-hep.org/x/hep/groot/rtree/rfunc"
)

func newFormula(r *Reader, f rfunc.Formula) (rfunc.Formula, error) {
	names := f.RVars()
	rvs, missing := formulaAutoLoad(r, names)
	if len(rvs) != len(names) {
		return nil, fmt.Errorf("rtree: could not find all needed ReadVars (missing: %v)", missing)
	}
	args := make([]any, len(rvs))
	for i, rv := range rvs {
		args[i] = rv.Value
	}
	err := f.Bind(args)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not bind formula to rvars: %w", err)
	}

	return f, nil
}

func formulaAutoLoad(r *Reader, idents []string) ([]*ReadVar, []string) {
	var (
		loaded  = make(map[string]*ReadVar, len(r.rvars))
		needed  = make([]*ReadVar, 0, len(idents))
		rvars   = NewReadVars(r.tree)
		all     = make(map[string]*ReadVar, len(rvars))
		missing []string
	)

	for i := range r.rvars {
		rvar := &r.rvars[i]
		loaded[rvar.Name] = rvar
		all[rvar.Name] = rvar
	}
	for i := range rvars {
		rvar := &rvars[i]
		if _, ok := all[rvar.Name]; ok {
			continue
		}
		all[rvar.Name] = rvar
	}
	for _, name := range idents {
		rvar, ok := all[name]
		if !ok {
			missing = append(missing, name)
			continue
		}
		if _, ok := loaded[name]; !ok {
			r.rvars = append(r.rvars, *rvar)
			rvar = &r.rvars[len(r.rvars)-1]
			loaded[name] = rvar
		}
		needed = append(needed, rvar)
	}

	return needed, missing
}
