package fastjet

import (
	"fmt"
)

// JetDefinition contains a full specification of how to carry out jet clustering.
type JetDefinition struct {
	alg        JetAlgorithm
	r          float64
	extra      float64
	strategy   Strategy
	recombiner Recombiner
	plugin     Plugin
}

func NewJetDefinition(alg JetAlgorithm, r float64, scheme RecombinationScheme, strategy Strategy) JetDefinition {
	return JetDefinition{
		alg:        alg,
		r:          r,
		recombiner: NewRecombiner(scheme),
		strategy:   strategy,
	}
}

func (def JetDefinition) Description() string {
	switch def.alg {
	case PluginAlgorithm:
		return def.plugin.Description()

	case KtAlgorithm:
		return fmt.Sprintf("Longitudinally invariant kt algorithm with R = %v and %s",
			def.R(), def.Recombiner().Description(),
		)

	case CambridgeAlgorithm:
		return fmt.Sprintf("Longitudinally invariant Cambridge/Aachen algorithm with R = %v and %s",
			def.R(), def.Recombiner().Description(),
		)

	case AntiKtAlgorithm:
		return fmt.Sprintf("Longitudinally invariant anti-kt algorithm with R = %v and %s",
			def.R(), def.Recombiner().Description(),
		)

	case GenKtAlgorithm:
		return fmt.Sprintf("Longitudinally invariant generalised kt algorithm with R = %v, p = %v and %s",
			def.R(), def.ExtraParam(), def.Recombiner().Description(),
		)

	case CambridgeForPassiveAlgorithm:
		return fmt.Sprintf("Longitudinally invariant Cambridge/Aache algorithm with R = %v, kt<%v as ghosts",
			def.R(), def.ExtraParam(),
		)

	case EeKtAlgorithm:
		return fmt.Sprintf("e+e- kt (Durham) algorithm with %s", def.Recombiner().Description())

	case EeGenKtAlgorithm:
		return fmt.Sprintf("e+e- generalised kt algorithm with R = %v, p = %s and %s",
			def.R(), def.ExtraParam(), def.Recombiner().Description(),
		)

	case UndefinedJetAlgorithm:
		return "uninitialised JetDefinition"

	default:
		panic(fmt.Errorf("fastjet.Description: invalid jet algorithm (%d)", int(def.alg)))
	}
}

func (def JetDefinition) R() float64 {
	return def.r
}

func (def JetDefinition) ExtraParam() float64 {
	return def.extra
}

func (def JetDefinition) Strategy() Strategy {
	return def.strategy
}

func (def JetDefinition) Recombiner() Recombiner {
	return def.recombiner
}

func (def JetDefinition) RecombinationScheme() RecombinationScheme {
	return def.recombiner.Scheme()
}

func (def JetDefinition) Algorithm() JetAlgorithm {
	return def.alg
}

func (def JetDefinition) Plugin() Plugin {
	return def.plugin
}

// to impl:
//  - ClusterSequence
//  - PseudoJet
//  - Selector + JetMedianBkgEstimator (only if compute-rho)
