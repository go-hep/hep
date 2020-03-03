// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepmc

import (
	"fmt"
)

// MomentumUnit describes the units of momentum quantities (MeV or GeV)
type MomentumUnit int

// LengthUnit describes the units of length quantities (mm or cm)
type LengthUnit int

const (
	// MEV is a Momentum in MeV (default)
	MEV MomentumUnit = iota
	// GEV is a Momentum in GeV
	GEV
)

const (
	// MM is a Length in mm (default)
	MM LengthUnit = iota
	// CM is a Length in cm
	CM
)

func (mu MomentumUnit) String() string {
	switch mu {
	case MEV:
		return "MEV"
	case GEV:
		return "GEV"
	}
	panic(fmt.Errorf("hepmc.units: invalid MomentumUnit value (%d)", int(mu)))
}

// MomentumUnitFromString creates a MomentumUnit value from its string representation
func MomentumUnitFromString(s string) (MomentumUnit, error) {
	switch s {
	case "MEV":
		return MEV, nil
	case "GEV":
		return GEV, nil
	}
	return -1, fmt.Errorf("hepmc.units: invalid MomentumUnit string-value (%s)", s)
}

func (lu LengthUnit) String() string {
	switch lu {
	case MM:
		return "MM"
	case CM:
		return "CM"
	}
	panic(fmt.Errorf("hepmc.units: invalid LengthUnit value (%d)", int(lu)))
}

// LengthUnitFromString creates a LengthUnit value from its string representation
func LengthUnitFromString(s string) (LengthUnit, error) {
	switch s {
	case "MM":
		return MM, nil
	case "CM":
		return CM, nil
	}
	return -1, fmt.Errorf("hepmc.units: invalid LengthUnit string-value (%s)", s)
}
