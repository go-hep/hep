// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

// JetAlgorithm defines the algorithm used for clustering jets
type JetAlgorithm int

const (
	UndefinedJetAlgorithm JetAlgorithm = iota
	KtAlgorithm
	CambridgeAlgorithm
	AntiKtAlgorithm
	GenKtAlgorithm
	CambridgeForPassiveAlgorithm
	GenKtForPassiveAlgorithm
	EeKtAlgorithm
	EeGenKtAlgorithm
	PluginAlgorithm

	AachenAlgorithm          = CambridgeAlgorithm
	CambridgeAachenAlgorithm = CambridgeAlgorithm
)
