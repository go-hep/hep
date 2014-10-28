// +build ignore

package slha

// Modsel holds program-independent model switches.
// e.g. which model of supersymetry breaking to use.
type Modsel struct {
	Model      int     // choice of SUSY breaking model
	Content    int     // choice of particle content
	GridPoints int     // number of points for a logarithmically spaced grid in Q
	QMax       float64 // largest Q scale
	PdgID      int     // PDG code for a particle
}

// SMInputs holds the measured values of SM parameters, used as boundary conditions in the spectrun calculation.
// These are also required for subsequent calculations to be consistent with the spectrim calculation.
type SMInputs struct {
	Data map[int]float64
}

// MinPar holds the input parameters for minimal/default models.
type MinPar struct {
	Data map[int]float64
}

// ExtPar holds the optional input parameters for non-minimal/non-universal models.
type ExtPar struct {
	Data map[int]float64
}

// Mass holds the mass spectrum parameters.
type Mass struct {
	Data map[int]float64
}

// NMix holds the Neutralino mixing matrix
type NMix struct{}

// UMix holds the Chargino U mixing matrix
type UMix struct{}

// VMix holds the Chargino V mixing matrix
type VMix struct{}

// StopMix holds the Stop mixing matrix
type StopMix struct{}

// SbotMix holds the Sbottom mixing matrix
type SbotMix struct{}

// StauMix holds the Stau mixing matrix
type StauMix struct{}

// Alpha holds the Higgs mixing angle alpha
type Alpha struct{}

// HMix holds Higgs parameters at scale Q.
type HMix struct{}

// Gauge holds gauge couplings at scale Q.
type Gauge struct{}

// MSoft holds soft SUSY breaking mass parameters at scale Q.
type MSoft struct{}

// AU holds trilinear couplings at scale Q.
type AU struct{}

// AD holds trilinear couplings at scale Q.
type AD struct{}

// AE holds trilinear couplings at scale Q.
type AE struct{}

// YU holds Yukawa couplings at scale Q.
type YU struct{}

// YD holds Yukawa couplings at scale Q.
type YD struct{}

// YE holds Yukawa couplings at scale Q.
type YE struct{}

// SpInfo holds informations from the spectrum calculator
type SpInfo struct{}
