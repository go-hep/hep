package ascii

import (
	"fmt"
	"io"
	"sync"

	"github.com/go-hep/hepmc"
)

const (
	genevent_start      = "HepMC::IO_GenEvent-START_EVENT_LISTING"
	ascii_start         = "HepMC::IO_Ascii-START_EVENT_LISTING"
	extendedascii_start = "HepMC::IO_ExtendedAscii-START_EVENT_LISTING"

	genevent_end      = "HepMC::IO_GenEvent-END_EVENT_LISTING"
	ascii_end         = "HepMC::IO_Ascii-END_EVENT_LISTING"
	extendedascii_end = "HepMC::IO_ExtendedAscii-END_EVENT_LISTING"

	pdt_start               = "HepMC::IO_Ascii-START_PARTICLE_DATA"
	extendedascii_pdt_start = "HepMC::IO_ExtendedAscii-START_PARTICLE_DATA"
	pdt_end                 = "HepMC::IO_Ascii-END_PARTICLE_DATA"
	extendedascii_pdt_end   = "HepMC::IO_ExtendedAscii-END_PARTICLE_DATA"
)

type Encoder struct {
	w    io.Writer
	once sync.Once
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(evt *hepmc.Event) error {
	var err error

	enc.once.Do(func() {
		_, err = fmt.Fprintf(
			enc.w,
			"\nHepMC::Version %s\n",
			hepmc.VersionName(),
		)
	})

	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(enc.w, "%s\n", genevent_start)
	if err != nil {
		return err
	}

	sig_bc := 0
	if evt.SignalVertex != nil {
		sig_bc = evt.SignalVertex.Barcode
	}
	// output the event data including the number of primary vertices
	// and the total number of vertices
	_, err = fmt.Fprintf(
		enc.w,
		"E %d %d %e %e %e %d %d %d",
		evt.EventNumber,
		evt.Mpi,
		evt.Scale,
		evt.AlphaQCD,
		evt.AlphaQED,
		evt.SignalProcessId,
		sig_bc,
		len(evt.RandomStates),
	)
	if err != nil {
		return err
	}
	for _, rndm := range evt.RandomStates {
		_, err = fmt.Fprintf(enc.w, " %e", rndm)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(enc.w, " %d", len(evt.Weights))
	if err != nil {
		return err
	}
	for _, weight := range evt.Weights {
		_, err = fmt.Fprintf(enc.w, " %e", weight)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(enc.w, "\n")
	if err != nil {
		return err
	}

	err = enc.encode_heavy_ion(evt.HeavyIon)
	if err != nil {
		return err
	}

	err = enc.encode_pdf_info(evt.PdfInfo)
	if err != nil {
		return err
	}

	// output all of the vertices
	for i, _ := range evt.Vertices {
		vtx := &evt.Vertices[i]
		err = enc.encode_vertex(vtx)
		if err != nil {
			return err
		}
	}
	return err
}

func (enc *Encoder) encode_vertex(vtx *hepmc.Vertex) error {
	var err error
	orphans := 0
	for _, p := range vtx.ParticlesIn {
		if p.ProdVertex == nil {
			orphans += 1
		}
	}

	_, err = fmt.Fprintf(
		enc.w,
		"V %d %d %e %e %e %e %d %d %d",
		vtx.Barcode,
		vtx.Id,
		vtx.Position.X(), vtx.Position.Y(), vtx.Position.Z(), vtx.Position.T(),
		orphans,
		len(vtx.ParticlesOut),
		len(vtx.Weights),
	)
	if err != nil {
		return err
	}
	for _, w := range vtx.Weights {
		_, err = fmt.Fprintf(enc.w, " %e", w)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(enc.w, "\n")
	if err != nil {
		return err
	}

	for _, p := range vtx.ParticlesIn {
		if p.ProdVertex == nil {
			err = enc.encode_particle(p)
			if err != nil {
				return err
			}
		}
	}
	for _, p := range vtx.ParticlesOut {
		err = enc.encode_particle(p)
		if err != nil {
			return err
		}
	}
	return err
}

func (enc *Encoder) encode_particle(p *hepmc.Particle) error {
	var err error

	end_bc := 0
	if p.EndVertex != nil {
		end_bc = p.EndVertex.Barcode
	}

	_, err = fmt.Fprintf(
		enc.w,
		"P %d %d %e %e %e %e %d %e %e %d",
		p.Barcode,
		p.PdgId,
		p.Momentum.Px(), p.Momentum.Py(), p.Momentum.Pz(), p.Momentum.E(),
		p.Status,
		p.Polarization.Theta,
		p.Polarization.Phi,
		end_bc,
	)
	if err != nil {
		return err
	}
	err = enc.encode_flow(&p.Flow)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(enc.w, "\n")
	return err
}

func (enc *Encoder) encode_flow(flow *hepmc.Flow) error {
	var err error
	_, err = fmt.Fprintf(enc.w, " %d", len(flow.Icode))
	if err != nil {
		return err
	}
	for k, v := range flow.Icode {
		_, err = fmt.Fprintf(enc.w, " %d %d", k, v)
		if err != nil {
			return err
		}
	}
	return err
}

func (enc *Encoder) encode_heavy_ion(hi *hepmc.HeavyIon) error {
	var err error
	if hi == nil {
		_, err = fmt.Fprintf(
			enc.w,
			"H %d %d %d %d %d %d %d %d %d %e %e %e %e\n",
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0., 0., 0., 0.,
		)
		return err
	}
	_, err = fmt.Fprintf(
		enc.w,
		"H %d %d %d %d %d %d %d %d %d %e %e %e %e\n",
		hi.Ncoll_hard,
		hi.Npart_proj,
		hi.Npart_targ,
		hi.Ncoll,
		hi.N_Nwounded_collisions,
		hi.Nwounded_N_collisions,
		hi.Nwounded_Nwounded_collisions,
		hi.Spectator_neutrons,
		hi.Spectator_protons,
		hi.Impact_parameter,
		hi.Event_plane_angle,
		hi.Eccentricity,
		hi.Sigma_inel_NN,
	)
	return err
}

func (enc *Encoder) encode_pdf_info(pdf *hepmc.PdfInfo) error {
	var err error
	if pdf == nil {
		_, err = fmt.Fprintf(
			enc.w,
			"F %d %d %e %e %e %e %e\n",
			0, 0, 0., 0., 0., 0., 0.,
		)
		return err
	}
	_, err = fmt.Fprintf(
		enc.w,
		"F %d %d %e %e %e %e %e\n",
		pdf.Id1, pdf.Id2,
		pdf.X1, pdf.X2,
		pdf.ScalePDF,
		pdf.Pdf1,
		pdf.Pdf2,
	)
	return err
}

// EOF
