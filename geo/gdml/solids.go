// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdml

import (
	"encoding/xml"
)

type Box struct {
	XMLName xml.Name `xml:"box"`
	Name    string   `xml:"name,attr"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
	Z       string   `xml:"z,attr"`
	LUnit   string   `xml:"lunit,attr"`
}

type Cone struct {
	XMLName  xml.Name `xml:"cone"`
	Name     string   `xml:"name,attr"`
	RMin1    string   `xml:"rmin1,attr"`    // inner radius at base of cone
	RMax1    string   `xml:"rmax1,attr"`    // outer radius at base of cone
	RMin2    string   `xml:"rmin2,attr"`    // inner radius at top of cone
	RMax2    string   `xml:"rmax2,attr"`    // outer radius at top of cone
	Z        string   `xml:"z,attr"`        // height of cone segment
	StartPhi string   `xml:"startphi,attr"` // start angle of the segment
	DPhi     string   `xml:"deltaphi,attr"` // angle of the segment
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
}

type Ellipsoid struct {
	XMLName xml.Name `xml:"ellipsoid"`
	Name    string   `xml:"name,attr"`
	Ax      string   `xml:"ax,attr"` // x semi axis
	By      string   `xml:"by,attr"` // y semi axis
	Cz      string   `xml:"cz,attr"` // z semi axis
	ZCut1   string   `xml:"zcut1,attr"`
	ZCut2   string   `xml:"zcut2,attr"`
	LUnit   string   `xml:"lunit,attr"`
}

type EllipticalTube struct {
	XMLName xml.Name `xml:"eltube"`
	Name    string   `xml:"name,attr"`
	Dx      string   `xml:"dx,attr"` // x semi axis
	Dy      string   `xml:"dy,attr"` // y semi axis
	Dz      string   `xml:"dz,attr"` // z semi axis
	LUnit   string   `xml:"lunit,attr"`
}

type EllipticalCone struct {
	XMLName xml.Name `xml:"elcone"`
	Name    string   `xml:"name,attr"`
	Dx      string   `xml:"dx,attr"` // x semi axis
	Dy      string   `xml:"dy,attr"` // y semi axis
	ZMax    string   `xml:"zmax,attr"`
	ZCut    string   `xml:"zcut,attr"`
	LUnit   string   `xml:"lunit,attr"`
}

type Orb struct {
	XMLName xml.Name `xml:"orb"`
	Name    string   `xml:"name,attr"`
	R       string   `xml:"r,attr"` // radius
	LUnit   string   `xml:"lunit,attr"`
}

type Paraboloid struct {
	XMLName xml.Name `xml:"paraboloid"`
	Name    string   `xml:"name,attr"`
	Rlo     string   `xml:"rlo,attr"` // radius at -z
	Rhi     string   `xml:"rhi,attr"` // radius at +z
	Dz      string   `xml:"dz,attr"`  // z length
	AUnit   string   `xml:"aunit,attr"`
	LUnit   string   `xml:"lunit,attr"`
}

type Parallelepiped struct {
	XMLName xml.Name `xml:"para"`
	Name    string   `xml:"name,attr"`
	X       string   `xml:"x,attr"`     // length of x
	Y       string   `xml:"y,attr"`     // length of y
	Z       string   `xml:"z,attr"`     // length of z
	Alpha   string   `xml:"alpha,attr"` // angle between x and z planes
	Theta   string   `xml:"theta,attr"` // polar angle of the line joining the centres of the faces at -z and +z
	Phi     string   `xml:"phi,attr"`   // azimuthal angle of the line joining the centres of the faces at -x and +z
	AUnit   string   `xml:"aunit,attr"`
	LUnit   string   `xml:"lunit,attr"`
}

type PolyCone struct {
	XMLName  xml.Name `xml:"polycone"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	StartPhi string   `xml:"startphi,attr"` // start angle of the segment - default: 0.0
	DPhi     string   `xml:"deltaphi,attr"` // angle of the segment
	ZPlanes  []ZPlane `xml:"zplane"`
}

type ZPlane struct {
	XMLName xml.Name `xml:"zplane"`
	RMin    string   `xml:"rmin,attr"` // inner radius of cone at this point
	RMax    string   `xml:"rmax,attr"` // outer radius of cone at this point
	Z       string   `xml:"z,attr"`    // z coordinate of the plane
}

type GenericPolyCone struct {
	XMLName  xml.Name  `xml:"genericPolycone"`
	Name     string    `xml:"name,attr"`
	AUnit    string    `xml:"aunit,attr"`
	LUnit    string    `xml:"lunit,attr"`
	StartPhi string    `xml:"startphi,attr"` // start angle of the segment - default: 0.0
	DPhi     string    `xml:"deltaphi,attr"` // angle of the segment
	RZPoints []RZPoint `xml:"rzpoint"`
}

type RZPoint struct {
	XMLName xml.Name `xml:"rzpoint"`
	R       string   `xml:"r,attr"` // r coordinate of this point
	Z       string   `xml:"z,attr"` // z coordinate of this point
}

type PolyHedra struct {
	XMLName  xml.Name `xml:"polyhedra"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	StartPhi string   `xml:"startphi,attr"` // start angle of the segment - default: 0.0
	DPhi     string   `xml:"deltaphi,attr"` // angle of the segment
	NumSides string   `xml:"numsides,attr"` // number of sides
	ZPlanes  []ZPlane `xml:"zplane"`
}

type GenericPolyHedra struct {
	XMLName  xml.Name  `xml:"genericPolyhedra"`
	Name     string    `xml:"name,attr"`
	AUnit    string    `xml:"aunit,attr"`
	LUnit    string    `xml:"lunit,attr"`
	StartPhi string    `xml:"startphi,attr"` // start angle of the segment - default: 0.0
	DPhi     string    `xml:"deltaphi,attr"` // angle of the segment
	NumSides string    `xml:"numsides,attr"` // number of sides
	RZPoints []RZPoint `xml:"rzpoint"`
}

type Sphere struct {
	XMLName    xml.Name `xml:"sphere"`
	Name       string   `xml:"name,attr"`
	AUnit      string   `xml:"aunit,attr"`
	LUnit      string   `xml:"lunit,attr"`
	RMin       string   `xml:"rmin,attr"`       // inner radius
	RMax       string   `xml:"rmax,attr"`       // outer radius
	StartPhi   string   `xml:"startphi,attr"`   // start angle of the segment - default: 0.0
	DPhi       string   `xml:"deltaphi,attr"`   // angle of the segment
	StartTheta string   `xml:"starttheta,attr"` // start angle of the segment - default: 0.0
	DTheta     string   `xml:"deltatheta,attr"` // angle of the segment
}

type Torus struct {
	XMLName  xml.Name `xml:"torus"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	RMin     string   `xml:"rmin,attr"`     // inner radius of segment
	RMax     string   `xml:"rmax,attr"`     // outer radius of segment
	Rtor     string   `xml:"rtor,attr"`     // swept radius of torus
	StartPhi string   `xml:"startphi,attr"` // start angle of the segment - default: 0.0
	DPhi     string   `xml:"deltaphi,attr"` // angle of the segment
}

type Trapezoid struct {
	XMLName xml.Name `xml:"trd"`
	Name    string   `xml:"name,attr"`
	LUnit   string   `xml:"lunit,attr"`
	X1      string   `xml:"x1,attr"` // x length at -z
	X2      string   `xml:"x2,attr"` // x length at +z
	Y1      string   `xml:"y1,attr"` // y length at -z
	Y2      string   `xml:"y2,attr"` // y length at +z
	Z       string   `xml:"z,attr"`  // z length
}

type GeneralTrapezoid struct {
	XMLName xml.Name `xml:"trap"`
	Name    string   `xml:"name,attr"`
	AUnit   string   `xml:"aunit,attr"`
	LUnit   string   `xml:"lunit,attr"`
	Z       string   `xml:"z,attr"`      // length along z axis
	Theta   string   `xml:"theta,attr"`  // polar angle to faces joining at -/+ z
	Phi     string   `xml:"phi,attr"`    // azimuthal angle of line joining centre of -z to centre of +z face
	Y1      string   `xml:"y1,attr"`     // length along y at the face -z
	X1      string   `xml:"x1,attr"`     // length along x at side y = -y1 of the face at -z
	X2      string   `xml:"x2,attr"`     // length along x at side y = +y1 of the face at -z
	Alpha1  string   `xml:"alpha1,attr"` // angle with respect to the y axis from the centre of side at y = -y1 to centre of y = +y1 of the face at -z
	Y2      string   `xml:"y2,attr"`     // length along y at the face +z
	X3      string   `xml:"x3,attr"`     // length along x at side y = -y1 of the face at +z
	X4      string   `xml:"x4,attr"`     // length along x at side y = +y1 of the face at +z
	Alpha2  string   `xml:"alpha2,attr"` // angle with respect to the y axis from the centre of side at y = -y2 to centre of y = +y2 of the face at +z
}

type HyperbolicTube struct {
	XMLName   xml.Name `xml:"hype"`
	Name      string   `xml:"name,attr"`
	LUnit     string   `xml:"lunit,attr"`
	RMin      string   `xml:"rmin,attr"`  // inside radius of tube
	RMax      string   `xml:"rmax,attr"`  // outside radius of tube
	InStereo  string   `xml:"inst,attr"`  // inner stereo
	OutStereo string   `xml:"outst,attr"` // outer stereo
	Z         string   `xml:"z,attr"`     // z length
}

type CutTube struct {
	XMLName  xml.Name `xml:"cutTube"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	Z        string   `xml:"z,attr"`        // z length
	RMin     string   `xml:"rmin,attr"`     // inside radius, default:0.0
	RMax     string   `xml:"rmax,attr"`     // outside radius
	StartPhi string   `xml:"startphi,attr"` // starting phi angle of segment, default:0.0
	DPhi     string   `xml:"deltaphi,attr"` // delta phi of angle
	LowX     string   `xml:"lowX,attr"`     // normal at lower z plane
	LowY     string   `xml:"lowY,attr"`     // normal at lower z plane
	LowZ     string   `xml:"lowZ,attr"`     // normal at lower z plane
	HighX    string   `xml:"highX,attr"`    // normal at upper z plane
	HighY    string   `xml:"highY,attr"`    // normal at upper z plane
	HighZ    string   `xml:"highZ,attr"`    // normal at upper z plane
}

type Tube struct {
	XMLName  xml.Name `xml:"tube"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	Z        string   `xml:"z,attr"`        // z length
	RMin     string   `xml:"rmin,attr"`     // inside radius, default:0.0
	RMax     string   `xml:"rmax,attr"`     // outside radius
	StartPhi string   `xml:"startphi,attr"` // starting phi angle of segment, default:0.0
	DPhi     string   `xml:"deltaphi,attr"` // delta phi of angle
}

type TwistedBox struct {
	XMLName  xml.Name `xml:"twistedbox"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	PhiTwist string   `xml:"PhiTwist,attr"`
	X        string   `xml:"x,attr"`
	Y        string   `xml:"y,attr"`
	Z        string   `xml:"z,attr"`
}

type TwistedTrapezoid struct {
	XMLName  xml.Name `xml:"twistedtrd"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	PhiTwist string   `xml:"PhiTwist,attr"`
	X1       string   `xml:"x1,attr"` // x length at -z
	X2       string   `xml:"x2,attr"` // x length at +z
	Y1       string   `xml:"y1,attr"` // y length at -z
	Y2       string   `xml:"y2,attr"` // y length at +z
	Z        string   `xml:"z,attr"`  // z length
}

type TwistedGeneralTrapezoid struct {
	XMLName  xml.Name `xml:"twistedtrap"`
	Name     string   `xml:"name,attr"`
	AUnit    string   `xml:"aunit,attr"`
	LUnit    string   `xml:"lunit,attr"`
	PhiTwist string   `xml:"PhiTwist,attr"`
	Z        string   `xml:"z,attr"`     // length along z axis
	Theta    string   `xml:"theta,attr"` // polar angle to faces joining at -/+ z
	Phi      string   `xml:"phi,attr"`   // azimuthal angle of line joining centre of -z to centre of +z face
	Y1       string   `xml:"y1,attr"`    // length along y at the face -z
	X1       string   `xml:"x1,attr"`    // length along x at side y = -y1 of the face at -z
	X2       string   `xml:"x2,attr"`    // length along x at side y = +y1 of the face at -z
	Y2       string   `xml:"y2,attr"`    // length along y at the face +z
	X3       string   `xml:"x3,attr"`    // length along x at side y = -y1 of the face at +z
	X4       string   `xml:"x4,attr"`    // length along x at side y = +y1 of the face at +z
	Alpha    string   `xml:"Alph,attr"`  // angle with respect to the y axis from the centre of side
}

type TwistedTube struct {
	XMLName      xml.Name `xml:"twistedtube"`
	Name         string   `xml:"name,attr"`
	AUnit        string   `xml:"aunit,attr"`
	LUnit        string   `xml:"lunit,attr"`
	EndInnerRad  string   `xml:"endinnerrad,attr"`  // inside radius at end of segment
	EndOuterRad  string   `xml:"endouterrad,attr"`  // outside radius at end of segment
	ZLen         string   `xml:"zlen,attr"`         // z length of tube segment
	TwistedAngle string   `xml:"twistedangle,attr"` // twist angle
	Phi          string   `xml:"phi,attr"`          // phi angle of segment
	MidInnerRad  string   `xml:"midinnerrad,attr"`  // inner radius at z=0
	MidOuterRad  string   `xml:"midouterrad,attr"`  // outer radius at z=0
	NSeg         string   `xml:"nseg,attr"`         // number of segments in TotPhi
	TotPhi       string   `xml:"totphi,attr"`       // total angle of all segments
}

type Extruded struct {
	XMLName  xml.Name   `xml:"xtru"`
	Name     string     `xml:"name,attr"`
	LUnit    string     `xml:"lunit,attr"`
	Vertices []Vertex2D `xml:"twoDimVertex"` // vertices of unbound blueprint polygon
	Sections []Section  `xml:"section"`      // z sections
}

type Vertex2D struct {
	XMLName xml.Name `xml:"twoDimVertex"`
	X       string   `xml:"x,attr"`
	Y       string   `xml:"y,attr"`
}

type Section struct {
	XMLName xml.Name `xml:"section"`
	ZOrder  string   `xml:"zOrder,attr"`        // index of the section
	ZPos    string   `xml:"zPosition,attr"`     // distance from the plane z=0
	XOff    string   `xml:"xOffset,attr"`       // x offset from centre point of original plane
	YOff    string   `xml:"yOffset,attr"`       // y offset from centre point of original plane
	Fact    string   `xml:"scalingFactor,attr"` // proportion to original blueprint
}

type ArbitraryTrapezoid struct {
	XMLName xml.Name `xml:"arb8"`
	Name    string   `xml:"name,attr"`
	LUnit   string   `xml:"lunit,attr"`
	V1x     string   `xml:"v1x,attr"` // vertex 1 x position
	V1y     string   `xml:"v1y,attr"` // vertex 1 y position
	V2x     string   `xml:"v2x,attr"` // vertex 2 x position
	V2y     string   `xml:"v2y,attr"` // vertex 2 y position
	V3x     string   `xml:"v3x,attr"` // vertex 3 x position
	V3y     string   `xml:"v3y,attr"` // vertex 3 y position
	V4x     string   `xml:"v4x,attr"` // vertex 4 x position
	V4y     string   `xml:"v4y,attr"` // vertex 4 y position
	V5x     string   `xml:"v5x,attr"` // vertex 5 x position
	V5y     string   `xml:"v5y,attr"` // vertex 5 y position
	V6x     string   `xml:"v6x,attr"` // vertex 6 x position
	V6y     string   `xml:"v6y,attr"` // vertex 6 y position
	V7x     string   `xml:"v7x,attr"` // vertex 7 x position
	V7y     string   `xml:"v7y,attr"` // vertex 7 y position
	V8x     string   `xml:"v8x,attr"` // vertex 8 x position
	V8y     string   `xml:"v8y,attr"` // vertex 8 y position
	Dz      string   `xml:"dz,attr"`  // half z length
}

type Tesselated struct {
	XMLName xml.Name       `xml:"tesselated"`
	Name    string         `xml:"name,attr"`
	Tris    []Triangular   `xml:"triangular"`
	Quads   []Quadrangular `xml:"quadrangular"`
}

type Triangular struct {
	XMLName xml.Name `xml:"triangular"`
	Vtx1    string   `xml:"vertex1,attr"`
	Vtx2    string   `xml:"vertex2,attr"`
	Vtx3    string   `xml:"vertex3,attr"`
	Type    string   `xml:"type,attr"` // vertex type: ABSOLUTE (default) or RELATIVE
}

type Quadrangular struct {
	XMLName xml.Name `xml:"quadrangular"`
	Vtx1    string   `xml:"vertex1,attr"`
	Vtx2    string   `xml:"vertex2,attr"`
	Vtx3    string   `xml:"vertex3,attr"`
	Vtx4    string   `xml:"vertex4,attr"`
	Type    string   `xml:"type,attr"` // vertex type: ABSOLUTE (default) or RELATIVE
}

type TetraHedron struct {
	XMLName xml.Name `xml:"tet"`
	Name    string   `xml:"name,attr"`
	Vtx1    string   `xml:"vertex1,attr"`
	Vtx2    string   `xml:"vertex2,attr"`
	Vtx3    string   `xml:"vertex3,attr"`
	Vtx4    string   `xml:"vertex4,attr"`
}

type ScaledSolid struct {
	XMLName xml.Name `xml:"scaledSolid"`
	Name    string   `xml:"name,attr"`
	Ref     SolidRef `xml:"solidref"`
	Scale   Scale    `xml:"scale"`
}

type SolidRef struct {
	XMLName xml.Name `xml:"solidref"`
	Ref     string   `xml:"ref,attr"`
}
