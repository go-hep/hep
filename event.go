package fads

import (
	"time"
)

type Event struct {
	Header RunHeader
}

type RunHeader struct {
	RunNbr int64     // run number
	EvtNbr int64     // event number
	Time   time.Time // time stamp

	Descr   string   // description of the simulation conditions (e.g. physics channels)
	Det     string   // detector name
	Subdets []string // active subdetectors
}

type ParticleID struct {
	Type   int       // type of this PID (user defined)
	PDG    int       // Particle Data Group id
	Prob   float64   // likelihood of this hypothesis (user defined)
	Author int       // author of this PID
	Params []float64 // parameters associated with this hypothesis
}

// EOF
