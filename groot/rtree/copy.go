// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
)

// Copy copies from src to dst until either the reader is depleted or
// an error occurs. It returns the number of bytes copied and the first error
// encountered while copying, if any.
func Copy(dst Writer, src *Reader) (int64, error) {
	// FIXME(sbinet): optimize for the case(s) where:
	//  - all branches are being copied
	//  - compression is the same
	// => it might not be needed to uncompress-recompress the baskets.

	var (
		tot int64
		err error
	)

	wvars := dst.(*wtree).wvars
	rvars := make([]ReadVar, len(wvars))
	for i, wvar := range wvars {
		rvars[i] = ReadVar{
			Name:  wvar.Name,
			Value: wvar.Value,
		}
	}

	orig := src.rvars
	defer func() {
		src.rvars = orig
		src.dirty = true
	}()
	src.rvars = rvars
	src.dirty = true

	err = src.Read(func(ctx RCtx) error {
		written, err := dst.Write()
		if err != nil {
			return fmt.Errorf("rtree: could not write entry %d to tree: %w", ctx.Entry, err)
		}
		tot += int64(written)
		return nil
	})
	if err != nil {
		return tot, fmt.Errorf("rtree: could not read through tree: %w", err)
	}

	return tot, err
}
