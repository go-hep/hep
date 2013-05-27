package lhef

// XSecInfo contains information given in the xsecinfo tag.
type XSecInfo struct {
	Neve       int64   // the number of events.
	TotXSec    float64 // the total cross section in pb.
	MaxWeight  float64 // the maximum weight.
	MeanWeight float64 // the average weight.
	NegWeights bool    // does the file contain negative weights ?
	VarWeights bool    // does the file contain varying weights ?
}

// Cut represents a cut used by the Matrix Element generator.
type Cut struct {
	Type string  // the variable in which to cut.
	NP1  string  // symbolic name for p1.
	NP2  string  // symbolic name for p2.
	P1   []int64 // the first types particle types for which this cut applies.
	P2   []int64 // the second types particle types for which this cut applies.
	Min  float64 // the minimum value of the variable
	Max  float64 // the maximum value of the variable
}

// ProcInfo represents the information in a procinfo tag.
type ProcInfo struct {
	Iproc       int32  // the id number for the process.
	Loops       int32  // the number of loops.
	QcdOrder    int32  // the number of QCD vertices.
	EwOrder     int32  // the number of electro-weak vertices.
	Fscheme     string // the factorization scheme used.
	Rscheme     string // the renormalization scheme used.
	Scheme      string // the NLO scheme used.
	Description string // Description of the process.
}

// MergeInfo represents the information in a mergeinfo tag.
type MergeInfo struct {
	Iproc        int32   // the id number for the process.
	Scheme       string  // the scheme used to reweight events.
	MergingScale float64 // the merging scale used if different from the cut definitions.
	MaxMult      bool    // is this event reweighted as if it was the maximum multiplicity.
}

// Weight represents the information in a weight tag.
type Weight struct {
	Name    string    // the identifier for this set of weights.
	Born    float64   // the relative size of the born cross section of this event.
	Sudakov float64   // the relative size of the sudakov applied to this event.
	Weights []float64 // the weights of this event.
}

// Clus represents a clustering of two particle entries into one as
// defined in a clustering tag.
type Clus struct {
	P1     int32   // the first particle entry that has been clustered.
	P2     int32   // the second particle entry that has been clustered.
	P0     int32   // the particle entry corresponding to the clustered particles.
	Scale  float64 // the scale in GeV associated with the clustering.
	Alphas float64 // the alpha_s used in the corresponding vertex, if this was used in the cross section.
}

// PDFInfo represents the information in a pdfinfo tag.
type PDFInfo struct {
	P1    int64   // type of the incoming particle 1.
	P2    int64   // type of the incoming particle 2.
	X1    float64 // x-value used for the incoming particle 1.
	X2    float64 // x-value used for the incoming particle 2.
	XF1   float64 // value of the PDF for the incoming particle 1.
	XF2   float64 // value of the PDF for the incoming particle 2.
	Scale float64 // scale used in the PDFs
}

// HEPRUP is a simple container corresponding to the Les Houches accord common block (User Process Run common block.)
// http://arxiv.org/abs/hep-ph/0109068 has more details.
// The members are named in the same way as in the common block.
// However, FORTRAN arrays are represented by slices, except for the arrays of
// length 2 which are represented as arrays (of size 2.)
type HEPRUP struct {
	IDBMUP     [2]int64            // PDG id's of beam particles.
	EBMUP      [2]float64          // Energy of beam particles (in GeV.)
	PDFGUP     [2]int32            // Author group for the PDF used for the beams according to the PDFLib specifications.
	PDFSUP     [2]int32            // Id number of the PDF used for the beams according to the PDFLib specifications.
	IDWTUP     int32               // Master switch indicating how the ME generator envisages the events weights should be interpreted according to the Les Houches accord.
	NPRUP      int32               // number of different subprocesses in this file.
	XSECUP     []float64           // cross-sections for the different subprocesses in pb.
	XERRUP     []float64           // statistical error in the cross sections for the different subprocesses in pb.
	XMAXUP     []float64           // maximum event weights (in HEPEUP.XWGTUP) for different subprocesses.
	LPRUP      []int32             // subprocess code for the different subprocesses.
	XSecInfo   XSecInfo            // contents of the xsecinfo tag
	Cuts       []Cut               // contents of the cuts tag.
	PTypes     map[string][]int64  // a map of codes for different particle types.
	ProcInfo   map[int64]ProcInfo  // contents of the procinfo tags
	MergeInfo  map[int64]MergeInfo // contents of the mergeinfo tags
	GenName    string              // name of the generator which produced the file.
	GenVersion string              // version of the generator which produced the file.
}

// EventGroup represents a set of events which are to be considered together.
type EventGroup struct {
	Events   []HEPEUP // the list of events to be considered together
	Nreal    int32    // number of real event in this event group.
	Ncounter int32    // number of counter events in this event group.
}

// HEPEUP is a simple container corresponding to the Les Houches accord common block (User Process Event common block.)
// http://arxiv.org/abs/hep-ph/0109068 has more details.
// The members are named in the same way as in the common block.
// However, FORTRAN arrays are represented by slices, except for the arrays of length 2 which are represented as arrays of size 2.
type HEPEUP struct {
	NUP        int32        // number of particle entries in the current event.
	IDPRUP     int32        // subprocess code for this event (as given in LPRUP)
	XWGTUP     float64      // weight for this event.
	XPDWUP     [2]float64   // PDF weights for the 2 incomong partons. Note that this variable is not present in the current LesHouches accord (http://arxiv.org/abs/hep-ph/0109068) but will hopefully be present in a future accord.
	SCALUP     float64      // scale in GeV used in the calculation of the PDFs in this event
	AQEDUP     float64      // value of the QED coupling used in this event.
	AQCDUP     float64      // value of the QCD coupling used in this event.
	IDUP       []int64      // PDG id's for the particle entries in this event.
	ISTUP      []int32      // status codes for the particle entries in this event.
	MOTHUP     [][2]int32   // indices for the first and last mother for the particle entries in this event.
	ICOLUP     [][2]int32   // colour-line indices (first(second) is (anti)colour) for the particle entries in this event.
	PUP        [][5]float64 // lab frame momentum (Px, Py, Pz, E and M in GeV) for the particle entries in this event.
	VTIMUP     []float64    // invariant lifetime (c*tau, distance from production to decay in mm) for the particle entries in this event.
	SPINUP     []float64    // spin info for the particle entries in this event given as the cosine of the angle between the spin vector of a particle and the 3-momentum of the decaying particle, specified in the lab frame.
	Weights    []Weight     // weights associated with this event.
	Clustering []Clus       // contents of the clustering tag.
	PdfInfo    PDFInfo      // contents of the pdfinfo tag.
	SubEvents  EventGroup   // events included in the group if this is not a single event.
}

// EOF
