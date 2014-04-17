package heppdt

import (
	"bytes"
	"io"
)

var defaultTable Table

type Table struct {
	name string
	pdt  map[PID]*Particle
	pid  map[string]PID
}

// New returns a new particle data table, initialized from r
func New(r io.Reader, n string) (Table, error) {
	t := Table{
		name: n,
		pdt:  make(map[PID]*Particle),
		pid:  make(map[string]PID),
	}
	err := parse(r, &t)
	return t, err
}

// Name returns the name of this particle data table
func (t *Table) Name() string {
	return t.name
}

// Len returns the size of the particle data table
func (t *Table) Len() int {
	return len(t.pdt)
}

// Pdt returns the particle data table
func (t *Table) Pdt() map[PID]*Particle {
	return t.pdt
}

// ParticleByID returns the particle information via particle ID
func (t *Table) ParticleByID(pid PID) *Particle {
	p, ok := t.pdt[pid]
	if !ok {
		return nil
	}
	return p
}

// ParticleByName returns the particle information via particle name
func (t *Table) ParticleByName(n string) *Particle {
	pid, ok := t.pid[n]
	if !ok {
		return nil
	}

	p, ok := t.pdt[pid]
	if !ok {
		return nil
	}
	return p
}

func init() {
	var err error
	defaultTable, err = New(bytes.NewBufferString(tabledata), "particle.tbl")
	if err != nil {
		panic(err)
	}
}

// EOF
