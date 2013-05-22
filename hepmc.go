package hepmc

import (
	"errors"
	"fmt"
	"io"
)

var errNilVtx = errors.New("hepmc: nil Vertex")
var errNilParticle = errors.New("hepmc: nil Particle")

// Event represents a record for MC generators (for use at any stage of generation)
//
// This type is intended as both a "container class" ( to store a MC
//  event for interface between MC generators and detector simulation )
//  and also as a "work in progress class" ( that could be used inside
//  a generator and modified as the event is built ).
type Event struct {
	SignalProcessId int     // id of the signal process
	EventNumber     int     // event number
	Mpi             int     // number of multi particle interactions
	Scale           float64 // energy scale,
	AlphaQCD        float64 // QCD coupling, see hep-ph/0109068
	AlphaQED        float64 // QED coupling, see hep-ph/0109068

	SignalVertex *Vertex      // signal vertex
	Beams        [2]*Particle // incoming beams
	Weights      Weights      // weights for this event. first weight is used by default for hit and miss
	RandomStates []int64      // container of random number generator states

	Vertices  map[int]*Vertex
	Particles map[int]*Particle

	CrossSection *CrossSection
	HeavyIon     *HeavyIon
	PdfInfo      *PdfInfo
	MomentumUnit MomentumUnit
	LengthUnit   LengthUnit
}

// AddVertex adds a vertex to this event
func (evt *Event) AddVertex(vtx *Vertex) error {
	if vtx == nil {
		return errNilVtx
	}
	if vtx.Event != nil && vtx.Event != evt {
		//TODO: warn and remove from previous event
	}
	return vtx.set_parent_event(evt)
}

func (evt *Event) Print(w io.Writer) error {
	var err error
	_, err = fmt.Fprintf(
		w,
		"________________________________________________________________________________\n")
	if err != nil {
		return err
	}

	sig_vtx := 0
	if evt.SignalVertex != nil {
		sig_vtx = evt.SignalVertex.Barcode
	}
	_, err = fmt.Fprintf(
		w,
		"GenEvent: #%04d ID=%5d SignalProcessGenVertex Barcode: %d\n",
		evt.EventNumber,
		evt.SignalProcessId,
		sig_vtx,
	)
	if err != nil {
		return err
	}
	return err
}

// Particle represents a generator particle within an event coming in/out of a vertex
//
// Particle is the basic building block of the event record
type Particle struct {
	Momentum      FourVector   // momentum vector
	PdgId         int          // id according to PDG convention
	Status        int          // status code as defined for HEPEVT
	Flow          Flow         // flow of this particle
	Polarization  Polarization // polarization of this particle
	ProdVertex    *Vertex      // pointer to production vertex (nil if vacuum or beam)
	EndVertex     *Vertex      // pointer to decay vertex (nil if not-decayed)
	Barcode       int          // unique identifier in the event
	GeneratedMass float64      // mass of this particle when it was generated
}

// Vertex represents a generator vertex within an event
// A vertex is indirectly (via particle "edges") linked to other
//   vertices ("nodes") to form a composite "graph"
type Vertex struct {
	Position     FourVector  // 4-vector of vertex [mm]
	ParticlesIn  []*Particle // all incoming particles
	ParticlesOut []*Particle // all outgoing particles
	Id           int         // vertex id
	Weights      Weights     // weights for this vertex
	Event        *Event      // pointer to event owning this vertex
	Barcode      int         // unique identifier in the event
}

func (vtx *Vertex) set_parent_event(evt *Event) error {
	var err error
	orig_evt := vtx.Event
	vtx.Event = evt
	if orig_evt == evt {
		return err
	}
	if evt != nil {
		evt.Vertices[vtx.Barcode] = vtx
	}
	if orig_evt != nil {
		delete(orig_evt.Vertices, vtx.Barcode)
	}
	// we also need to loop over all the particles which are owned by
	// this vertex and remove their barcodes from the old event.
	for _, p := range vtx.ParticlesIn {
		if p.ProdVertex == nil {
			if evt != nil {
				evt.Particles[p.Barcode] = p
			}
			if orig_evt != nil {
				delete(orig_evt.Particles, p.Barcode)
			}
		}
	}

	for _, p := range vtx.ParticlesOut {
		if evt != nil {
			evt.Particles[p.Barcode] = p
		}
		if orig_evt != nil {
			delete(orig_evt.Particles, p.Barcode)
		}
	}
	return err
}

// AddParticleIn adds a particle to the list of in-coming particles to this vertex
func (vtx *Vertex) AddParticleIn(p *Particle) error {
	var err error
	if p == nil {
		return errNilParticle
	}
	// if p had a decay vertex, remove it from that vertex's list
	if p.EndVertex != nil {
		err = p.EndVertex.remove_particle_in(p)
		if err != nil {
			return err
		}
	}
	// make sure we don't add it twice...
	err = vtx.remove_particle_in(p)
	if err != nil {
		return err
	}
	p.EndVertex = vtx
	vtx.ParticlesIn = append(vtx.ParticlesIn, p)
	return err
}

// AddParticleOut adds a particle to the list of out-going particles to this vertex
func (vtx *Vertex) AddParticleOut(p *Particle) error {
	var err error
	if p == nil {
		return errNilParticle
	}
	// if p had a production vertex, remove it from that vertex's list
	if p.ProdVertex != nil {
		err = p.ProdVertex.remove_particle_out(p)
		if err != nil {
			return err
		}
	}
	// make sure we don't add it twice...
	err = vtx.remove_particle_out(p)
	if err != nil {
		return err
	}
	p.ProdVertex = vtx
	vtx.ParticlesOut = append(vtx.ParticlesOut, p)
	return err
}

func (vtx *Vertex) remove_particle_in(p *Particle) error {
	var err error
	nparts := len(vtx.ParticlesIn)
	switch nparts {
	case 0:
		//FIXME: logical error ?
		return err
	case 1:
		vtx.ParticlesIn = make([]*Particle, 0)
		return err
	}
	parts := make([]*Particle, 0, nparts-1)
	for _, pp := range vtx.ParticlesIn {
		if pp == p {
			continue
		}
		parts = append(parts, pp)
	}
	vtx.ParticlesIn = parts
	return err
}

func (vtx *Vertex) remove_particle_out(p *Particle) error {
	var err error
	nparts := len(vtx.ParticlesOut)
	switch nparts {
	case 0:
		//FIXME: logical error ?
		return err
	case 1:
		vtx.ParticlesOut = make([]*Particle, 0)
		return err
	}
	parts := make([]*Particle, 0, nparts-1)
	for _, pp := range vtx.ParticlesOut {
		if pp == p {
			continue
		}
		parts = append(parts, pp)
	}
	vtx.ParticlesOut = parts
	return err
}

type HeavyIon struct {
	Ncoll_hard                   int     // number of hard scatterings
	Npart_proj                   int     // number of projectile participants
	Npart_targ                   int     // number of target participants
	Ncoll                        int     // number of NN (nucleon-nucleon) collisions
	N_Nwounded_collisions        int     // Number of N-Nwounded collisions
	Nwounded_N_collisions        int     // Number of Nwounded-N collisons
	Nwounded_Nwounded_collisions int     // Number of Nwounded-Nwounded collisions
	Spectator_neutrons           int     // Number of spectators neutrons
	Spectator_protons            int     // Number of spectators protons
	Impact_parameter             float32 // Impact Parameter(fm) of collision
	Event_plane_angle            float32 // Azimuthal angle of event plane
	Eccentricity                 float32 // eccentricity of participating nucleons in the transverse plane (as in phobos nucl-ex/0510031)
	Sigma_inel_NN                float32 // nucleon-nucleon inelastic (including diffractive) cross-section
}

// CrossSection is used to store the generated cross section.
// This type is meant to be used to pass, on an event by event basis,
// the current best guess of the total cross section.
// It is expected that the final cross section will be stored elsewhere.
type CrossSection struct {
	Value float64 // value of the cross-section (in pb)
	Error float64 // error on the value of the cross-section (in pb)
	//IsSet bool
}

type PdfInfo struct {
	Id1      int     // flavour code of first parton
	Id2      int     // flavour code of second parton
	LHAPdf1  int     // LHA PDF id of first parton
	LHAPdf2  int     // LHA PDF id of second parton
	X1       float64 // fraction of beam momentum carried by first parton ("beam side")
	X2       float64 // fraction of beam momentum carried by second parton ("target side")
	ScalePDF float64 //  Q-scale used in evaluation of PDF's   (in GeV)
	Pdf1     float64 // PDF (id1, x1, Q)
	Pdf2     float64 // PDF (id2, x2, Q)
}

// Flow represents a particle's flow and keeps track of an arbitrary number of flow patterns within a graph (i.e. color flow, charge flow, lepton number flow,...)
//
// Flow patterns are coded with an integer, in the same manner as in Herwig.
// Note: 0 is NOT allowed as code index nor as flow code since it
//       is used to indicate null.
//
// This class can be used to keep track of flow patterns within
//  a graph. An example is color flow. If we have two quarks going through
//  an s-channel gluon to form two more quarks:
//
//  \q1       /q3   then we can keep track of the color flow with the
//   \_______/      HepMC::Flow class as follows:
//   /   g   \.
//  /q2       \q4
//
//  lets say the color flows from q2-->g-->q3  and q1-->g-->q4
//  the individual colors are unimportant, but the flow pattern is.
//  We can capture this flow by assigning the first pattern (q2-->g-->q3)
//  a unique (arbitrary) flow code 678 and the second pattern (q1-->g-->q4)
//  flow code 269  ( you can ask HepMC::Flow to choose
//  a unique code for you using Flow::set_unique_icode() ).
//  these codes with the particles as follows:
//    q2->flow().set_icode(1,678);
//    g->flow().set_icode(1,678);
//    q3->flow().set_icode(1,678);
//    q1->flow().set_icode(1,269);
//    g->flow().set_icode(2,269);
//    q4->flow().set_icode(1,269);
//  later on if we wish to know the color partner of q1 we can ask for a list
//  of all particles connected via this code to q1 which do have less than
//  2 color partners using:
//    vector<GenParticle*> result=q1->dangling_connected_partners(q1->icode(1),1,2);
//  this will return a list containing q1 and q4.
//    vector<GenParticle*> result=q1->connected_partners(q1->icode(1),1,2);
//  would return a list containing q1, g, and q4.
type Flow struct {
	Particle *Particle   // the particle this flow describes
	Icode    map[int]int // flow patterns as (code_index, icode)
}

type Polarization struct {
	Theta float64 // polar angle of polarization in radians [0, math.Pi)
	Phi   float64 // azimuthal angle of polarization in radians [0, 2*math.Pi)
}

type Weights struct {
	Slice []float64      // the slice of weight values
	Map   map[string]int // the map of name->index-in-the-slice
}

func (w Weights) At(n string) float64 {
	idx, ok := w.Map[n]
	if ok {
		return w.Slice[idx]
	}
	panic("hepmc.Weights.At: invalid name [" + n + "]")
}

func NewWeights() Weights {
	return Weights{
		Slice: make([]float64, 0, 1),
		Map:   make(map[string]int),
	}
}

// EOF
