package hepmc

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

type hepmc_ftype int

const (
	_ hepmc_ftype = iota
	hepmc_genevent
	hepmc_ascii
	hepmc_extendedascii
	hepmc_ascii_pdt
	hepmc_extendedascii_pdt
)

// EOF
