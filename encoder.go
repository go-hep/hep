package hepmc

import (
	"fmt"
	"io"
	"sort"
)

// Encoder encodes a hepmc Event into a stream.
type Encoder struct {
	w            io.Writer
	seen_evt_hdr bool
}

// NewEncoder returns a new hepmc Encoder that writes into the io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Close closes the encoder and adds a footer to the stream.
func (enc *Encoder) Close() error {
	var err error
	if enc.seen_evt_hdr {
		_, err = fmt.Fprintf(
			enc.w,
			"%s\n",
			endGenEvent,
		)
		if err != nil {
			return err
		}
	}
	return err
}

// Encode writes evt into the stream.
func (enc *Encoder) Encode(evt *Event) error {
	var err error

	if !enc.seen_evt_hdr {
		_, err = fmt.Fprintf(
			enc.w,
			"\nHepMC::Version %s\n",
			VersionName(),
		)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(enc.w, "%s\n", startGenEvent)
		if err != nil {
			return err
		}

		enc.seen_evt_hdr = true
	}

	sig_bc := 0
	if evt.SignalVertex != nil {
		sig_bc = evt.SignalVertex.Barcode
	}
	bp1 := 0
	if evt.Beams[0] != nil {
		bp1 = evt.Beams[0].Barcode
	}
	bp2 := 0
	if evt.Beams[1] != nil {
		bp2 = evt.Beams[1].Barcode
	}
	// output the event data including the number of primary vertices
	// and the total number of vertices
	_, err = fmt.Fprintf(
		enc.w,
		"E %d %d %1.16e %1.16e %1.16e %d %d %d %d %d %d",
		evt.EventNumber,
		evt.Mpi,
		evt.Scale,
		evt.AlphaQCD,
		evt.AlphaQED,
		evt.SignalProcessId,
		sig_bc,
		len(evt.Vertices),
		bp1,
		bp2,
		len(evt.RandomStates),
	)
	if err != nil {
		return err
	}
	for _, rndm := range evt.RandomStates {
		_, err = fmt.Fprintf(enc.w, " %1.16e", rndm)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(enc.w, " %d", len(evt.Weights.Slice))
	if err != nil {
		return err
	}
	// we need to iterate over the weights in the same order than their names
	// (we'll make sure of that in the 'N' line)
	for _, weight := range evt.Weights.Slice {
		_, err = fmt.Fprintf(enc.w, " %1.16e", weight)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(enc.w, "\n")
	if err != nil {
		return err
	}
	if len(evt.Weights.Slice) > 0 {
		nn := len(evt.Weights.Slice)
		names := make(map[int]string, nn)
		for k, v := range evt.Weights.Map {
			names[v] = k
		}
		_, err = fmt.Fprintf(enc.w, "N %d ", nn)
		if err != nil {
			return err
		}
		for iw := 0; iw < nn; iw++ {
			_, err = fmt.Fprintf(enc.w, "%q ", names[iw])
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(enc.w, "\n")
		if err != nil {
			return err
		}
	}

	// units
	_, err = fmt.Fprintf(
		enc.w,
		"U %s %s\n",
		evt.MomentumUnit,
		evt.LengthUnit,
	)
	if err != nil {
		return err
	}

	// cross-section
	if evt.CrossSection != nil {
		err = enc.encode_cross_section(evt.CrossSection)
		if err != nil {
			return err
		}
	}

	if evt.HeavyIon != nil {
		err = enc.encode_heavy_ion(evt.HeavyIon)
		if err != nil {
			return err
		}
	}

	err = enc.encode_pdf_info(evt.PdfInfo)
	if err != nil {
		return err
	}

	// output all of the vertices
	vertices := make([]*Vertex, 0, len(evt.Vertices))
	for _, vtx := range evt.Vertices {
		vertices = append(vertices, vtx)
	}
	sort.Sort(sort.Reverse(t_vertices(vertices)))
	for _, vtx := range vertices {
		err = enc.encode_vertex(vtx)
		if err != nil {
			return err
		}
	}
	return err
}

func (enc *Encoder) encode_vertex(vtx *Vertex) error {
	var err error
	orphans := 0
	for _, p := range vtx.ParticlesIn {
		if p.ProdVertex == nil {
			orphans++
		}
	}

	_, err = fmt.Fprintf(
		enc.w,
		"V %d %d %1.16e %1.16e %1.16e %1.16e %d %d %d",
		vtx.Barcode,
		vtx.Id,
		vtx.Position.X(), vtx.Position.Y(), vtx.Position.Z(), vtx.Position.T(),
		orphans,
		len(vtx.ParticlesOut),
		len(vtx.Weights.Slice),
	)
	if err != nil {
		return err
	}
	for _, w := range vtx.Weights.Slice {
		_, err = fmt.Fprintf(enc.w, " %1.16e", w)
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

func (enc *Encoder) encode_particle(p *Particle) error {
	var err error

	end_bc := 0
	if p.EndVertex != nil {
		end_bc = p.EndVertex.Barcode
	}

	_, err = fmt.Fprintf(
		enc.w,
		"P %d %d %1.16e %1.16e %1.16e %1.16e %1.16e %d %1.16e %1.16e %d",
		p.Barcode,
		p.PdgId,
		p.Momentum.Px(), p.Momentum.Py(), p.Momentum.Pz(), p.Momentum.E(),
		p.GeneratedMass,
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

func (enc *Encoder) encode_flow(flow *Flow) error {
	var err error
	_, err = fmt.Fprintf(enc.w, " %d", len(flow.Icode))
	if err != nil {
		return err
	}
	icodes := make([]int, 0, len(flow.Icode))
	for k := range flow.Icode {
		icodes = append(icodes, k)
	}
	sort.Ints(icodes)
	for _, k := range icodes {
		v := flow.Icode[k]
		_, err = fmt.Fprintf(enc.w, " %d %d", k, v)
		if err != nil {
			return err
		}
	}
	return err
}

func (enc *Encoder) encode_cross_section(x *CrossSection) error {
	var err error
	_, err = fmt.Fprintf(
		enc.w,
		"C %1.16e %1.16e\n",
		x.Value,
		x.Error,
	)
	return err
}

func (enc *Encoder) encode_heavy_ion(hi *HeavyIon) error {
	var err error
	_, err = fmt.Fprintf(
		enc.w,
		"H %d %d %d %d %d %d %d %d %d %1.16e %1.16e %1.16e %1.16e\n",
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

func (enc *Encoder) encode_pdf_info(pdf *PdfInfo) error {
	var err error
	if pdf == nil {
		_, err = fmt.Fprintf(
			enc.w,
			"F %d %d %1.16e %1.16e %1.16e %1.16e %1.16e %d %d\n",
			0, 0, 0., 0., 0., 0., 0., 0, 0,
		)
		return err
	}
	_, err = fmt.Fprintf(
		enc.w,
		"F %d %d %1.16e %1.16e %1.16e %1.16e %1.16e %d %d\n",
		pdf.ID1,
		pdf.ID2,
		pdf.X1,
		pdf.X2,
		pdf.ScalePDF,
		pdf.Pdf1,
		pdf.Pdf2,
		pdf.LHAPdf1,
		pdf.LHAPdf2,
	)
	return err
}

// EOF
