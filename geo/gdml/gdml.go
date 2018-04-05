// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gdml parses and interprets GDML files.
// Geometry Description Markup Language (GDML) files are specialized XML-based
// language files designed to describe the geometries of detectors associated
// with physics measurements.
//
// See:
//  http://gdml.web.cern.ch/GDML/doc/GDMLmanual.pdf
//
package gdml

import "encoding/xml"

// Constant describes a named constant in a GDML file.
type Constant struct {
	Name  string  `xml:"name,attr"`
	Value float64 `xml:"value,attr"`
}

// Quantity is a constant with a unit.
type Quantity struct {
	Name  string  `xml:"name,attr"`
	Type  string  `xml:"type,attr"`
	Value float64 `xml:"value,attr"`
	Unit  string  `xml:"unit,attr"`
}

// Variable is a named value in a GDML file.
// Once defined, a variable can be used anywhere inside the file.
// The value of a variable is evaluated each time it is used.
type Variable struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Position struct {
	Name string `xml:"name,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
	Z    string `xml:"z,attr"`
	Unit string `xml:"unit,attr"`
}

type Rotation struct {
	Name string `xml:"name,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
	Z    string `xml:"z,attr"`
	Unit string `xml:"unit,attr"`
}

type Scale struct {
	Name string `xml:"name,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
	Z    string `xml:"z,attr"`
}

type Matrix struct {
	Name   string `xml:"name,attr"`
	Cols   int    `xml:"coldim,attr"`
	Values string `xml:"values,attr"`
}

type Isotope struct {
	Name string `xml:"name,attr"`
	Z    int    `xml:"Z,attr"`
	N    int    `xml:"N,attr"`
	Atom Atom   `xml:"atom"`
}

type Element struct {
	Name      string     `xml:"name,attr"`
	Z         int        `xml:"Z,attr"`
	Formula   string     `xml:"formula,attr"`
	Atom      Atom       `xml:"atom"`
	Fractions []Fraction `xml:"fraction"`
}

type Atom struct {
	Type  string  `xml:"type,attr"`
	Value float64 `xml:"value,attr"`
}

type Fraction struct {
	Ref string  `xml:"ref,attr"`
	N   float64 `xml:"n,attr"`
}

type Material struct {
	Name       string      `xml:"name,attr"`
	Z          float64     `xml:"Z,attr"`
	Formula    string      `xml:"formula,attr"`
	Density    Density     `xml:"D"`
	Atom       Atom        `xml:"atom"`
	Composites []Composite `xml:"composite"`
	Fractions  []Fraction  `xml:"fraction"`
}

type Density struct {
	Value float64 `xml:"value,attr"`
}

type Composite struct {
	N   int    `xml:"n,attr"`
	Ref string `xml:"ref,attr"`
}

type Schema struct {
	XMLName xml.Name `xml:"data"`
}
