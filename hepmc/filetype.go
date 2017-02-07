package hepmc

const (
	startGenEvent      = "HepMC::IO_GenEvent-START_EVENT_LISTING"
	startASCII         = "HepMC::IO_Ascii-START_EVENT_LISTING"
	startExtendedASCII = "HepMC::IO_ExtendedAscii-START_EVENT_LISTING"

	endGenEvent      = "HepMC::IO_GenEvent-END_EVENT_LISTING"
	endASCII         = "HepMC::IO_Ascii-END_EVENT_LISTING"
	endExtendedASCII = "HepMC::IO_ExtendedAscii-END_EVENT_LISTING"

	startPdt              = "HepMC::IO_Ascii-START_PARTICLE_DATA"
	startExtendedASCIIPdt = "HepMC::IO_ExtendedAscii-START_PARTICLE_DATA"
	endPdt                = "HepMC::IO_Ascii-END_PARTICLE_DATA"
	endExtendedASCIIPdt   = "HepMC::IO_ExtendedAscii-END_PARTICLE_DATA"
)

type hepmcFileType int

const (
	_ hepmcFileType = iota
	hepmcGenEvent
	hepmcASCII
	hepmcExtendedASCII
	hepmcASCIIPdt
	hepmcExtendedASCIIPdt
)

// EOF
