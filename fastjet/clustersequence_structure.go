// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

// ClusterSequenceStructure is a ClusterSequence that implements
// the JetStructure interface.
type ClusterSequenceStructure struct {
	cs *ClusterSequence
}

func (css ClusterSequenceStructure) Constituents(jet *Jet) ([]Jet, error) {
	return css.cs.Constituents(jet)
}
