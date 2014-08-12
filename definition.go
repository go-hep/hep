package fastjet

import "fmt"

// RecombinationScheme defines the recombination choice for the 4-momenta of
// pseudo-jets during the clustering procedure
type RecombinationScheme int

const (
	EScheme   RecombinationScheme = iota // summing the 4-momenta
	PtScheme                             // pt-weighted recombination of y,phi
	Pt2Scheme                            // pt^2 weighted recombination of y,phi
	EtScheme
	Et2Scheme
	BIPtScheme
	BIPt2Scheme

	ExternalScheme RecombinationScheme = 99
)

func (s RecombinationScheme) String() string {
	switch s {
	case EScheme:
		return "E"
	case PtScheme:
		return "Pt"
	case Pt2Scheme:
		return "Pt2"
	case EtScheme:
		return "Et"
	case Et2Scheme:
		return "Et2"
	case BIPtScheme:
		return "BIPt"
	case BIPt2Scheme:
		return "BIPt2"

	case ExternalScheme:
		return "External"

	default:
		panic(fmt.Errorf("fastjet: invalid RecombinationScheme (%d)", int(s)))
	}

	panic("unreachable")
}

// Strategy defines the algorithmic strategy used while clustering.
type Strategy int

const (
	N2MinHeapTiledStrategy Strategy = -4
	N2TiledStrategy                 = -3
	N2PoorTiledStrategy             = -2
	N2PlainStrategy                 = -1
	N3DumbStrategy                  = 0
	BestStrategy                    = 1
	NlnNStrategy                    = 2
	NlnN3piStrategy                 = 3
	NlnN4piStrategy                 = 4
	NlnNCam4piStrategy              = 14
	NlnNCam2pi2RStrategy            = 13
	NlnNCamStrategy                 = 12

	PluginStrategy Strategy = 999
)

func (s Strategy) String() string {
	switch s {
	case N2MinHeapTiledStrategy:
		return "N2MinHeapTiled"
	case N2TiledStrategy:
		return "N2Tiled"
	case N2PoorTiledStrategy:
		return "N2PoorTiled"
	case N2PlainStrategy:
		return "N2Plain"
	case N3DumbStrategy:
		return "N3Dumb"
	case BestStrategy:
		return "Best"
	case NlnNStrategy:
		return "NlnN"
	case NlnN3piStrategy:
		return "NlnN3pi"
	case NlnN4piStrategy:
		return "NlnN4pi"
	case NlnNCam4piStrategy:
		return "NlnNCam4pi"
	case NlnNCam2pi2RStrategy:
		return "NlnNCam2pi2R"
	case NlnNCamStrategy:
		return "NlnNCam"

	case PluginStrategy:
		return "Plugin"

	default:
		panic(fmt.Errorf("fastjet: invalid Strategy (%d)", int(s)))
	}

	panic("unreachable")
}

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

	panic("unreachable")
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
