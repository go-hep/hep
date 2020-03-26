// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
)

// Copy copies from src to dst until either the source tree is depleted or
// an error occurs. It returns the number of bytes copied and the first error
// encountered while copying, if any.
func Copy(dst Writer, src Tree) (int64, error) {
	return CopyN(dst, src, src.Entries())
}

// Copy copies n events (or until an error) from src to dst until either the source tree is depleted or
// an error occurs. It returns the number of bytes copied and the first error
// encountered while copying, if any.
func CopyN(dst Writer, src Tree, n int64) (int64, error) {
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
			Name: wvar.Name,
			//			Leaf:  wvar.Name, // ??
			Value: wvar.Value,
		}
	}

	scan, err := NewScannerVars(src, rvars...)
	if err != nil {
		return 0, fmt.Errorf("rtree: could not create scanner: %w", err)
	}

	for scan.Next() && scan.Entry() < n {
		err = scan.Scan()
		if err != nil {
			return tot, fmt.Errorf("rtree: could not read entry %d from tree: %w", scan.Entry(), err)
		}

		written, err := dst.Write()
		if err != nil {
			return tot, fmt.Errorf("rtree: could not write entry %d to tree: %w", scan.Entry(), err)
		}
		tot += int64(written)
	}

	err = scan.Err()
	if err != nil {
		return tot, fmt.Errorf("rtree: could not scan through tree: %w", err)
	}

	return tot, err
}
