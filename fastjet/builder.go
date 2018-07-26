// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

// Builder builds jets out of 4-vectors
type Builder interface {
	// InclusiveJets returns all jets (in the sense of
	// the inclusive algorithm) with pt >= ptmin
	InclusiveJets(ptmin float64) ([]Jet, error)

	// ExclusiveJets

	// Constituents retrieves the constituents of a jet
	Constituents(jet *Jet) ([]Jet, error)
}
