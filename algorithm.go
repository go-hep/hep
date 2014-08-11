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
