package hepmc

import (
	"fmt"
)

type MomentumUnit int
type LengthUnit int

const (
	MEV MomentumUnit = iota // Momentum in MeV (default)
	GEV                     // Momentum in GeV
)

const (
	MM LengthUnit = iota // Length in mm (default)
	CM                   // Length in cm
)

func (mu MomentumUnit) String() string {
	switch mu {
	case MEV:
		return "MEV"
	case GEV:
		return "GEV"
	}
	err := fmt.Errorf("hepmc.units: invalid MomentumUnit value (%d)", int(mu))
	panic(err.Error())
}

func MomentumUnitFromString(s string) (MomentumUnit, error) {
	switch s {
	case "MEV":
		return MEV, nil
	case "GEV":
		return GEV, nil
	}
	err := fmt.Errorf("hepmc.units: invalid MomentumUnit string-value (%s)", s)
	return -1, err
}

func (lu LengthUnit) String() string {
	switch lu {
	case MM:
		return "MM"
	case CM:
		return "CM"
	}
	err := fmt.Errorf("hepmc.units: invalid LengthUnit value (%d)", int(lu))
	panic(err.Error())
}

func LengthUnitFromString(s string) (LengthUnit, error) {
	switch s {
	case "MM":
		return MM, nil
	case "CM":
		return CM, nil
	}
	err := fmt.Errorf("hepmc.units: invalid LengthUnit string-value (%s)", s)
	return -1, err
}

// EOF
