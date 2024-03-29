// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rootcnv provides tools to convert between HepMC2 data types
// and ROOT data structures.
package rootcnv // import "go-hep.org/x/hep/hepmc/rootcnv"

import (
	"sort"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/sliceop"
)

type event struct {
	SignalProcessID  int32   `groot:"Event_processID"` // id of the signal process
	Event_number     int32   `groot:"Event_nbr"`       // event number
	Event_mpi        int32   `groot:"Event_mpi"`       // number of multi particle interactions
	Event_scale      float64 `groot:"Event_scale"`     // energy scale,
	Event_alphaQCD   float64 `groot:"Event_alphaQCD"`  // QCD coupling, see hep-ph/0109068
	Event_alphaQED   float64 `groot:"Event_alphaQED"`  // QED coupling, see hep-ph/0109068
	Event_barcodeSPV int32   `groot:"Event_barcodeSPV"`
	Event_barcodeBP1 int32   `groot:"Event_barcodeBP1"`
	Event_barcodeBP2 int32   `groot:"Event_barcodeBP2"`
	Event_nvtx       int32   `groot:"Event_nvtx"`
	Event_npart      int32   `groot:"Event_npart"`
	Event_inbcs      []int32 `groot:"Event_inbcs"`  // Event barcodes of (p-in) for each vertex
	Event_outbcs     []int32 `groot:"Event_outbcs"` // Event barcodes of (p-out) for each vertex

	WeightsSlice    []float64 `groot:"Weights_slice"`
	WeightsMapKeys  []string  `groot:"Weights_keys"`
	WeightsMapNames []int32   `groot:"Weights_names"`
	RandomStates    []int64   `groot:"Random_states"`

	XsectValue float64 `groot:"Xsection_value"`
	XsectError float64 `groot:"Xsection_error"`

	HI_ncollHard         int32   `groot:"HI_ncoll_hard"`
	HI_npartProj         int32   `groot:"HI_npart_proj"`
	HI_npartTarg         int32   `groot:"HI_npart_targ"`
	HI_ncoll             int32   `groot:"HI_ncoll"`
	HI_nnwColl           int32   `groot:"HI_nnw_coll"`
	HI_nwNColl           int32   `groot:"HI_nwn_coll"`
	HI_nwNwColl          int32   `groot:"HI_nwnw_coll"`
	HI_spectatorNeutrons int32   `groot:"HI_spect_neutrons"`
	HI_spectatorProtons  int32   `groot:"HI_spect_protons"`
	HI_impactParameter   float32 `groot:"HI_impact_param"`
	HI_eventPlaneAngle   float32 `groot:"HI_evt_plane_angle"`
	HI_eccentricity      float32 `groot:"HI_eccentricity"`
	HI_sigmaInelNN       float32 `groot:"HI_sigma_inel_nn"`

	PDF_Parton1 int32   `groot:"PDF_parton1"`
	PDF_Parton2 int32   `groot:"PDF_parton2"`
	PDF_X1      float64 `groot:"PDF_x1"`
	PDF_X2      float64 `groot:"PDF_x2"`
	PDF_Q2      float64 `groot:"PDF_Q2"`
	PDF_X1f     float64 `groot:"PDF_x1f"`
	PDF_X2f     float64 `groot:"PDF_x2f"`
	PDF_ID1     int32   `groot:"PDF_id1"`
	PDF_ID2     int32   `groot:"PDF_id2"`

	MomentumUnit int8 `groot:"Momentum_unit"`
	LengthUnit   int8 `groot:"Length_unit"`

	Vertex_x    []float64   `groot:"Vertex_x"`
	Vertex_y    []float64   `groot:"Vertex_y"`
	Vertex_z    []float64   `groot:"Vertex_z"`
	Vertex_t    []float64   `groot:"Vertex_t"`
	Vertex_id   []int32     `groot:"Vertex_id"`
	Vertex_wsli [][]float64 `groot:"Vertex_wsli"`
	Vertex_wkey [][]string  `groot:"Vertex_wkey"`
	Vertex_wval [][]int32   `groot:"Vertex_wval"`
	Vertex_bc   []int32     `groot:"Vertex_bc"`
	Vertex_nin  []int32     `groot:"Vertex_nin"`
	Vertex_nout []int32     `groot:"Vertex_nout"`

	Particle_bc   []int32   `groot:"Particle_bc"`
	Particle_pid  []int64   `groot:"Particle_pid"`
	Particle_px   []float64 `groot:"Particle_px"`
	Particle_py   []float64 `groot:"Particle_py"`
	Particle_pz   []float64 `groot:"Particle_pz"`
	Particle_ene  []float64 `groot:"Particle_ene"`
	Particle_mass []float64 `groot:"Particle_mass"`
	//	Particle_nflow  []int32   `groot:"Particle_nflow"`
	Particle_flow   [][]int32 `groot:"Particle_flow"`
	Particle_theta  []float64 `groot:"Particle_theta"`
	Particle_phi    []float64 `groot:"Particle_phi"`
	Particle_status []int32   `groot:"Particle_status"`
	Particle_pvtx   []int32   `groot:"Particle_pvtx"`
	Particle_evtx   []int32   `groot:"Particle_evtx"`

	barcodes []int // work buffer
}

func (evt *event) write(h *hepmc.Event) error {
	h.SignalProcessID = int(evt.SignalProcessID)
	h.EventNumber = int(evt.Event_number)
	h.Mpi = int(evt.Event_mpi)
	h.Scale = evt.Event_scale
	h.AlphaQCD = evt.Event_alphaQCD
	h.AlphaQED = evt.Event_alphaQED
	h.Vertices = make(map[int]*hepmc.Vertex, evt.Event_nvtx)
	h.Particles = make(map[int]*hepmc.Particle, evt.Event_npart)

	h.Weights.Slice = sliceop.Resize(h.Weights.Slice, len(evt.WeightsSlice))
	copy(h.Weights.Slice, evt.WeightsSlice)
	h.Weights.Map = make(map[string]int, len(evt.WeightsMapKeys))
	for i, k := range evt.WeightsMapKeys {
		v := evt.WeightsMapNames[i]
		h.Weights.Map[k] = int(v)
	}
	h.RandomStates = sliceop.Resize(h.RandomStates, len(evt.RandomStates))
	copy(h.RandomStates, evt.RandomStates)

	switch {
	case !evt.hasXSect():
		h.CrossSection = nil
	default:
		h.CrossSection = &hepmc.CrossSection{
			Value: evt.XsectValue,
			Error: evt.XsectError,
		}
	}

	switch {
	case !evt.hasHI():
		h.HeavyIon = nil
	default:
		h.HeavyIon = &hepmc.HeavyIon{
			NCollHard:         int(evt.HI_ncollHard),
			NPartProj:         int(evt.HI_npartProj),
			NPartTarg:         int(evt.HI_npartTarg),
			NColl:             int(evt.HI_ncoll),
			NNwColl:           int(evt.HI_nnwColl),
			NwNColl:           int(evt.HI_nwNColl),
			NwNwColl:          int(evt.HI_nwNwColl),
			SpectatorNeutrons: int(evt.HI_spectatorNeutrons),
			SpectatorProtons:  int(evt.HI_spectatorProtons),
			ImpactParameter:   evt.HI_impactParameter,
			EventPlaneAngle:   evt.HI_eventPlaneAngle,
			Eccentricity:      evt.HI_eccentricity,
			SigmaInelNN:       evt.HI_sigmaInelNN,
		}
	}

	if h.PdfInfo == nil {
		h.PdfInfo = new(hepmc.PdfInfo)
	}
	h.PdfInfo.ID1 = int(evt.PDF_Parton1)
	h.PdfInfo.ID2 = int(evt.PDF_Parton2)
	h.PdfInfo.X1 = evt.PDF_X1
	h.PdfInfo.X2 = evt.PDF_X2
	h.PdfInfo.ScalePDF = evt.PDF_Q2
	h.PdfInfo.Pdf1 = evt.PDF_X1f
	h.PdfInfo.Pdf2 = evt.PDF_X2f
	h.PdfInfo.LHAPdf1 = int(evt.PDF_ID1)
	h.PdfInfo.LHAPdf2 = int(evt.PDF_ID2)

	h.MomentumUnit = hepmc.MomentumUnit(evt.MomentumUnit)
	h.LengthUnit = hepmc.LengthUnit(evt.LengthUnit)

	// flow := evt.Particle_flow
	for i := range evt.Particle_bc {
		flow := evt.Particle_flow[i]
		p := &hepmc.Particle{
			Momentum: fmom.NewPxPyPzE(
				evt.Particle_px[i],
				evt.Particle_py[i],
				evt.Particle_pz[i],
				evt.Particle_ene[i],
			),
			PdgID:  evt.Particle_pid[i],
			Status: int(evt.Particle_status[i]),
			Flow: hepmc.Flow{
				Icode: make(map[int]int, len(flow)),
			},
			Polarization: hepmc.Polarization{
				Theta: evt.Particle_theta[i],
				Phi:   evt.Particle_phi[i],
			},
			Barcode:       int(evt.Particle_bc[i]),
			GeneratedMass: evt.Particle_mass[i],
		}
		p.Flow.Particle = p
		for ii := 0; ii < len(flow); ii += 2 {
			p.Flow.Icode[int(flow[ii])] = int(flow[ii+1])
		}
		h.Particles[p.Barcode] = p
	}

	var (
		pin  = evt.Event_inbcs
		pout = evt.Event_outbcs
	)
	for i := range evt.Vertex_bc {
		wsli := evt.Vertex_wsli[i]
		wkey := evt.Vertex_wkey[i]
		wval := evt.Vertex_wval[i]
		vtx := &hepmc.Vertex{
			Position: fmom.NewPxPyPzE(
				evt.Vertex_x[i],
				evt.Vertex_y[i],
				evt.Vertex_z[i],
				evt.Vertex_t[i],
			),
			ParticlesIn:  make([]*hepmc.Particle, evt.Vertex_nin[i]),
			ParticlesOut: make([]*hepmc.Particle, evt.Vertex_nout[i]),
			ID:           int(evt.Vertex_id[i]),
			Weights: hepmc.Weights{
				Slice: make([]float64, len(wsli)),
				Map:   make(map[string]int, len(wkey)),
			},
			Event:   h,
			Barcode: int(evt.Vertex_bc[i]),
		}
		for i := range vtx.ParticlesIn {
			p := h.Particles[int(pin[i])]
			p.EndVertex = vtx
			vtx.ParticlesIn[i] = p
		}
		pin = pin[len(vtx.ParticlesIn):]

		for i := range vtx.ParticlesOut {
			p := h.Particles[int(pout[i])]
			p.ProdVertex = vtx
			vtx.ParticlesOut[i] = p
		}
		pout = pout[len(vtx.ParticlesOut):]

		copy(vtx.Weights.Slice, wsli)
		for i, k := range wkey {
			v := wval[i]
			vtx.Weights.Map[k] = int(v)
		}

		h.Vertices[vtx.Barcode] = vtx
	}

	switch evt.Event_barcodeSPV {
	case 0:
		h.SignalVertex = nil
	default:
		h.SignalVertex = h.Vertices[int(evt.Event_barcodeSPV)]
	}

	switch evt.Event_barcodeBP1 {
	case 0:
		h.Beams[0] = nil
	default:
		h.Beams[0] = h.Particles[int(evt.Event_barcodeBP1)]
	}

	switch evt.Event_barcodeBP2 {
	case 0:
		h.Beams[1] = nil
	default:
		h.Beams[1] = h.Particles[int(evt.Event_barcodeBP2)]
	}

	return nil
}

func (evt *event) read(h *hepmc.Event) error {
	evt.SignalProcessID = int32(h.SignalProcessID)
	evt.Event_number = int32(h.EventNumber)
	evt.Event_mpi = int32(h.Mpi)
	evt.Event_scale = h.Scale
	evt.Event_alphaQCD = h.AlphaQCD
	evt.Event_alphaQED = h.AlphaQED
	switch {
	case h.SignalVertex != nil:
		evt.Event_barcodeSPV = int32(h.SignalVertex.Barcode)
	default:
		evt.Event_barcodeSPV = 0
	}
	evt.Event_nvtx = int32(len(h.Vertices))
	evt.Event_npart = int32(len(h.Particles))
	switch {
	case h.Beams[0] != nil:
		evt.Event_barcodeBP1 = int32(h.Beams[0].Barcode)
	default:
		evt.Event_barcodeBP1 = 0
	}
	switch {
	case h.Beams[1] != nil:
		evt.Event_barcodeBP2 = int32(h.Beams[1].Barcode)
	default:
		evt.Event_barcodeBP2 = 0
	}

	evt.WeightsSlice = sliceop.Resize(evt.WeightsSlice, len(h.Weights.Slice))
	copy(evt.WeightsSlice, h.Weights.Slice)
	evt.WeightsMapKeys = sliceop.Resize(evt.WeightsMapKeys, len(h.Weights.Map))[:0]
	evt.WeightsMapNames = sliceop.Resize(evt.WeightsMapNames, len(h.Weights.Map))[:0]
	for k, v := range h.Weights.Map {
		evt.WeightsMapKeys = append(evt.WeightsMapKeys, k)
		evt.WeightsMapNames = append(evt.WeightsMapNames, int32(v))
	}
	evt.RandomStates = sliceop.Resize(evt.RandomStates, len(h.RandomStates))
	copy(evt.RandomStates, h.RandomStates)

	switch xsect := h.CrossSection; xsect {
	case nil:
		evt.XsectValue = 0
		evt.XsectError = 0
	default:
		evt.XsectValue = h.CrossSection.Value
		evt.XsectError = h.CrossSection.Error
	}

	switch hi := h.HeavyIon; hi {
	case nil:
		evt.HI_ncollHard = 0
		evt.HI_npartProj = 0
		evt.HI_npartTarg = 0
		evt.HI_ncoll = 0
		evt.HI_nnwColl = 0
		evt.HI_nwNColl = 0
		evt.HI_nwNwColl = 0
		evt.HI_spectatorNeutrons = 0
		evt.HI_spectatorProtons = 0
		evt.HI_impactParameter = 0
		evt.HI_eventPlaneAngle = 0
		evt.HI_eccentricity = 0
		evt.HI_sigmaInelNN = 0
	default:
		evt.HI_ncollHard = int32(hi.NCollHard)
		evt.HI_npartProj = int32(hi.NPartProj)
		evt.HI_npartTarg = int32(hi.NPartTarg)
		evt.HI_ncoll = int32(hi.NColl)
		evt.HI_nnwColl = int32(hi.NNwColl)
		evt.HI_nwNColl = int32(hi.NwNColl)
		evt.HI_nwNwColl = int32(hi.NwNwColl)
		evt.HI_spectatorNeutrons = int32(hi.SpectatorNeutrons)
		evt.HI_spectatorProtons = int32(hi.SpectatorProtons)
		evt.HI_impactParameter = hi.ImpactParameter
		evt.HI_eventPlaneAngle = hi.EventPlaneAngle
		evt.HI_eccentricity = hi.Eccentricity
		evt.HI_sigmaInelNN = hi.SigmaInelNN
	}

	evt.PDF_Parton1 = int32(h.PdfInfo.ID1)
	evt.PDF_Parton2 = int32(h.PdfInfo.ID2)
	evt.PDF_X1 = h.PdfInfo.X1
	evt.PDF_X2 = h.PdfInfo.X2
	evt.PDF_Q2 = h.PdfInfo.ScalePDF
	evt.PDF_X1f = h.PdfInfo.Pdf1
	evt.PDF_X2f = h.PdfInfo.Pdf2
	evt.PDF_ID1 = int32(h.PdfInfo.LHAPdf1)
	evt.PDF_ID2 = int32(h.PdfInfo.LHAPdf2)

	evt.MomentumUnit = int8(h.MomentumUnit)
	evt.LengthUnit = int8(h.LengthUnit)

	evt.barcodes = sliceop.Resize(evt.barcodes, len(h.Vertices))[:0]
	for bc := range h.Vertices {
		evt.barcodes = append(evt.barcodes, bc)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(evt.barcodes)))

	n := len(h.Vertices)
	evt.Vertex_x = sliceop.Resize(evt.Vertex_x, n)[:0]
	evt.Vertex_y = sliceop.Resize(evt.Vertex_y, n)[:0]
	evt.Vertex_z = sliceop.Resize(evt.Vertex_z, n)[:0]
	evt.Vertex_t = sliceop.Resize(evt.Vertex_t, n)[:0]
	evt.Vertex_id = sliceop.Resize(evt.Vertex_id, n)[:0]
	evt.Vertex_wsli = sliceop.Resize(evt.Vertex_wsli, n)[:0]
	evt.Vertex_wkey = sliceop.Resize(evt.Vertex_wkey, n)[:0]
	evt.Vertex_wval = sliceop.Resize(evt.Vertex_wval, n)[:0]
	evt.Vertex_bc = sliceop.Resize(evt.Vertex_bc, n)[:0]
	evt.Vertex_nin = sliceop.Resize(evt.Vertex_nin, n)[:0]
	evt.Vertex_nout = sliceop.Resize(evt.Vertex_nout, n)[:0]

	for _, bc := range evt.barcodes {
		vtx := h.Vertices[bc]
		evt.Vertex_x = append(evt.Vertex_x, vtx.Position.X())
		evt.Vertex_y = append(evt.Vertex_y, vtx.Position.Y())
		evt.Vertex_z = append(evt.Vertex_z, vtx.Position.Z())
		evt.Vertex_t = append(evt.Vertex_t, vtx.Position.T())
		evt.Vertex_id = append(evt.Vertex_id, int32(vtx.ID))
		evt.Vertex_bc = append(evt.Vertex_bc, int32(vtx.Barcode))
		evt.Vertex_nin = append(evt.Vertex_nin, int32(len(vtx.ParticlesIn)))
		evt.Vertex_nout = append(evt.Vertex_nout, int32(len(vtx.ParticlesOut)))
		for _, p := range vtx.ParticlesIn {
			evt.Event_inbcs = append(evt.Event_inbcs, int32(p.Barcode))
		}
		for _, p := range vtx.ParticlesOut {
			evt.Event_outbcs = append(evt.Event_outbcs, int32(p.Barcode))
		}

		wsli := make([]float64, len(vtx.Weights.Slice))
		copy(wsli, vtx.Weights.Slice)
		wkey := make([]string, 0, len(vtx.Weights.Map))
		wval := make([]int32, 0, len(vtx.Weights.Map))
		for k, v := range vtx.Weights.Map {
			wkey = append(wkey, k)
			wval = append(wval, int32(v))
		}
		evt.Vertex_wsli = append(evt.Vertex_wsli, wsli)
		evt.Vertex_wkey = append(evt.Vertex_wkey, wkey)
		evt.Vertex_wval = append(evt.Vertex_wval, wval)
	}

	evt.barcodes = sliceop.Resize(evt.barcodes, len(h.Particles))[:0]
	for bc := range h.Particles {
		evt.barcodes = append(evt.barcodes, bc)
	}
	sort.Ints(evt.barcodes)

	n = len(h.Particles)
	evt.Particle_bc = sliceop.Resize(evt.Particle_bc, n)[:0]
	evt.Particle_pid = sliceop.Resize(evt.Particle_pid, n)[:0]
	evt.Particle_px = sliceop.Resize(evt.Particle_px, n)[:0]
	evt.Particle_py = sliceop.Resize(evt.Particle_py, n)[:0]
	evt.Particle_pz = sliceop.Resize(evt.Particle_pz, n)[:0]
	evt.Particle_ene = sliceop.Resize(evt.Particle_ene, n)[:0]
	evt.Particle_mass = sliceop.Resize(evt.Particle_mass, n)[:0]
	//	evt.Particle_nflow = sliceop.Resize(evt.Particle_nflow, n)[:0]
	evt.Particle_flow = sliceop.Resize(evt.Particle_flow, n)[:0]
	evt.Particle_theta = sliceop.Resize(evt.Particle_theta, n)[:0]
	evt.Particle_phi = sliceop.Resize(evt.Particle_phi, n)[:0]
	evt.Particle_status = sliceop.Resize(evt.Particle_status, n)[:0]
	evt.Particle_pvtx = sliceop.Resize(evt.Particle_pvtx, n)[:0]
	evt.Particle_evtx = sliceop.Resize(evt.Particle_evtx, n)[:0]

	for _, bc := range evt.barcodes {
		p := h.Particles[bc]
		switch vtx := p.ProdVertex; vtx {
		case nil:
			evt.Particle_pvtx = append(evt.Particle_pvtx, 0)
		default:
			evt.Particle_pvtx = append(evt.Particle_pvtx, int32(vtx.Barcode))
		}

		evt.Particle_bc = append(evt.Particle_bc, int32(p.Barcode))
		evt.Particle_pid = append(evt.Particle_pid, p.PdgID)
		evt.Particle_px = append(evt.Particle_px, p.Momentum.Px())
		evt.Particle_py = append(evt.Particle_py, p.Momentum.Py())
		evt.Particle_pz = append(evt.Particle_pz, p.Momentum.Pz())
		evt.Particle_ene = append(evt.Particle_ene, p.Momentum.E())
		evt.Particle_mass = append(evt.Particle_mass, p.GeneratedMass)
		//		evt.Particle_nflow = append(evt.Particle_nflow, int32(len(p.Flow.Icode)))
		flow := make([]int32, 0, 2*len(p.Flow.Icode))
		for k, v := range p.Flow.Icode {
			flow = append(flow, int32(k), int32(v))
		}
		evt.Particle_flow = append(evt.Particle_flow, flow)
		evt.Particle_theta = append(evt.Particle_theta, p.Polarization.Theta)
		evt.Particle_phi = append(evt.Particle_phi, p.Polarization.Phi)
		evt.Particle_status = append(evt.Particle_status, int32(p.Status))
		switch vtx := p.EndVertex; vtx {
		case nil:
			evt.Particle_evtx = append(evt.Particle_evtx, 0)
		default:
			evt.Particle_evtx = append(evt.Particle_evtx, int32(vtx.Barcode))
		}
	}

	return nil
}

func (evt *event) reset() {
	evt.Event_inbcs = evt.Event_inbcs[:0]
	evt.Event_outbcs = evt.Event_outbcs[:0]
	evt.WeightsSlice = evt.WeightsSlice[:0]
	evt.WeightsMapKeys = evt.WeightsMapKeys[:0]
	evt.WeightsMapNames = evt.WeightsMapNames[:0]
	evt.RandomStates = evt.RandomStates[:0]

	evt.Vertex_x = evt.Vertex_x[:0]
	evt.Vertex_y = evt.Vertex_y[:0]
	evt.Vertex_z = evt.Vertex_z[:0]
	evt.Vertex_t = evt.Vertex_t[:0]
	evt.Vertex_wsli = evt.Vertex_wsli[:0]
	evt.Vertex_wkey = evt.Vertex_wkey[:0]
	evt.Vertex_wval = evt.Vertex_wval[:0]
	evt.Vertex_id = evt.Vertex_id[:0]
	evt.Vertex_bc = evt.Vertex_bc[:0]
	evt.Vertex_nin = evt.Vertex_nin[:0]
	evt.Vertex_nout = evt.Vertex_nout[:0]

	evt.Particle_bc = evt.Particle_bc[:0]
	evt.Particle_pid = evt.Particle_pid[:0]
	evt.Particle_px = evt.Particle_px[:0]
	evt.Particle_py = evt.Particle_py[:0]
	evt.Particle_pz = evt.Particle_pz[:0]
	evt.Particle_ene = evt.Particle_ene[:0]
	evt.Particle_mass = evt.Particle_mass[:0]
	//	evt.Particle_nflow = evt.Particle_nflow[:0]
	evt.Particle_flow = evt.Particle_flow[:0]
	evt.Particle_theta = evt.Particle_theta[:0]
	evt.Particle_phi = evt.Particle_phi[:0]
	evt.Particle_status = evt.Particle_status[:0]
	evt.Particle_pvtx = evt.Particle_pvtx[:0]
	evt.Particle_evtx = evt.Particle_evtx[:0]
}

func (evt *event) hasXSect() bool {
	return !(evt.XsectError == 0 && evt.XsectValue == 0)
}

func (evt *event) hasHI() bool {
	return !(evt.HI_ncollHard == 0 && evt.HI_npartProj == 0 && evt.HI_npartTarg == 0 && evt.HI_ncoll == 0 &&
		evt.HI_nnwColl == 0 && evt.HI_nwNColl == 0 && evt.HI_nwNwColl == 0 && evt.HI_spectatorNeutrons == 0 &&
		evt.HI_spectatorProtons == 0 && evt.HI_impactParameter == 0 && evt.HI_eventPlaneAngle == 0 &&
		evt.HI_eccentricity == 0 && evt.HI_sigmaInelNN == 0)
}
