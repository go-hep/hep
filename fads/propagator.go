// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
)

type Propagator struct {
	fwk.TaskBase

	radius  float64
	radius2 float64
	halflen float64
	bz      float64

	input  string
	output string

	hadrons string
	eles    string
	muons   string
}

func (tsk *Propagator) Configure(ctx fwk.Context) error {
	var err error

	tsk.radius2 = tsk.radius * tsk.radius
	if tsk.radius < 1.0e-2 {
		return fmt.Errorf("%s: too small radius value (%v)", tsk.Name(), tsk.radius)
	}

	if tsk.halflen < 1.0e-2 {
		return fmt.Errorf("%s: too small 1/2-length value (%v)", tsk.Name(), tsk.halflen)
	}

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.hadrons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.eles, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.muons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *Propagator) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Propagator) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Propagator) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	msg.Debugf(">>> candidates: %v\n", len(input))

	output := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	hadrons := make([]Candidate, 0)
	eles := make([]Candidate, 0)
	muons := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.hadrons, hadrons)
		if err != nil {
			return
		}
		err = store.Put(tsk.eles, eles)
		if err != nil {
			return
		}
		err = store.Put(tsk.muons, muons)
		if err != nil {
			return
		}
	}()

	const (
		cLight  = 2.99792458e8
		cLight2 = cLight * cLight
	)

	for i := range input {
		cand := &input[i]
		x := cand.Pos.X() * 1e-3
		y := cand.Pos.Y() * 1e-3
		z := cand.Pos.Z() * 1e-3

		// is particle inside cylinder ?
		if math.Hypot(x, y) > tsk.radius || math.Abs(z) > tsk.halflen {
			continue
		}

		px := cand.Mom.Px()
		py := cand.Mom.Py()
		pt2 := px*px + py*py
		if pt2 < 1e-9 {
			continue
		}

		q := float64(cand.Charge())
		if math.Abs(q) < 1e-9 || math.Abs(tsk.bz) < 1e-9 {
			// solve pt2*t^2 + 2(px.x + py.y)*t + (radius2 - x.x - y.y) = 0
			v := px*y - py*x
			discr2 := pt2*tsk.radius2 - v*v
			if discr2 < 0 {
				// no solution
				continue
			}

			v = px*x + py*y
			discr := math.Sqrt(discr2)
			t1 := (-v + discr) / pt2
			t2 := (-v - discr) / pt2
			t := 0.0
			switch t1 < 0 {
			case true:
				t = t2
			case false:
				t = t1
			}

			pz := cand.Mom.Pz()
			e := cand.Mom.E()

			zt := z + pz*t
			if math.Abs(zt) > tsk.halflen {
				invpz := 1.0 / pz
				t3 := (+tsk.halflen - z) * invpz
				t4 := (-tsk.halflen - z) * invpz
				if t3 < 0 {
					t = t4
				} else {
					t = t3
				}
			}

			xt := x + px*t
			yt := y + py*t
			zt = z + pz*t

			mother := cand
			c := cand.Clone()
			c.Pos = fmom.NewPxPyPzE(xt*1e3, yt*1e3, zt*1e3, cand.Pos.T()+t*e*1e3)
			c.Add(mother)

			output = append(output, *c)
			if math.Abs(q) > 1e-9 {
				pid := c.Pid
				if pid < 0 {
					pid = 0
				}
				switch pid {
				case 11:
					eles = append(eles, *c)
				case 13:
					muons = append(muons, *c)
				default:
					hadrons = append(hadrons, *c)
				}
			}
		} else {
			// 1.  initial transverse momentum p_{T0} : Part->pt
			//     initial transverse momentum direction \phi_0 = -atan(p_X0/p_Y0)
			//     relativistic gamma : gamma = E/mc² ; gammam = gamma \times m
			//     giration frequency \omega = q/(gamma m) fBz
			//     helix radius r = p_T0 / (omega gamma m)

			e := cand.Mom.E()
			pt := cand.Mom.Pt()
			pz := cand.Mom.Pz()

			gammam := e * 1e9 / cLight2  // in eV/c^2
			omega := q * tsk.bz / gammam // omega is in [89875518 / s]
			r := pt / (q * tsk.bz) * 1e9 / cLight

			phi0 := math.Atan2(py, px) // [rad] in [-pi,pi)

			// 2. helix axis coordinates
			xc := x + r*math.Sin(phi0)
			yc := y - r*math.Cos(phi0)
			rc := math.Hypot(xc, yc)
			phic := math.Atan2(yc, xc)
			phi := phic
			if xc < 0 {
				phi += math.Pi
			}

			// 3. time evaluation t = TMath::Min(t_r, t_z)
			//    t_r : time to exit from the sides
			//    t_z : time to exit from the front or the back
			tr := 0.0 // in [ns]
			signpz := -1.0
			if pz > 0 {
				signpz = +1.0
			}

			tz := 0.0
			if pz == 0 {
				tz = 1e99
			} else {
				tz = gammam / (pz * 1e9 / cLight) * (-z + tsk.halflen*signpz)
			}

			absr := math.Abs(r)
			t := 0.0
			if rc+absr < tsk.radius {
				// helix does not cross the cylinder sides
				t = tz
			} else {
				asinrho := math.Asin((tsk.radius*tsk.radius - rc*rc - r*r) * 0.5 / (absr * rc))
				delta := phi0 - phi
				if delta < -math.Pi {
					delta += 2.0 * math.Pi
				}
				if delta > math.Pi {
					delta -= 2.0 * math.Pi
				}

				t1 := (delta + asinrho) / omega
				t2 := (delta + math.Pi - asinrho) / omega
				t3 := (delta + math.Pi + asinrho) / omega
				t4 := (delta - asinrho) / omega
				t5 := (delta - math.Pi - asinrho) / omega
				t6 := (delta - math.Pi + asinrho) / omega

				if t1 < 0 {
					t1 = 1.0e99
				}
				if t2 < 0 {
					t2 = 1.0e99
				}
				if t3 < 0 {
					t3 = 1.0e99
				}
				if t4 < 0 {
					t4 = 1.0e99
				}
				if t5 < 0 {
					t5 = 1.0e99
				}
				if t6 < 0 {
					t6 = 1.0e99
				}

				tra := math.Min(t1, math.Min(t2, t3))
				trb := math.Min(t4, math.Min(t5, t6))
				tr = math.Min(tra, trb)
				t = math.Min(tr, tz)
			}

			// 4. position in terms of x(t), y(t), z(t)
			xt := xc + r*math.Sin(omega*t-phi0)
			yt := yc + r*math.Cos(omega*t-phi0)
			zt := z + pz*1.0e9/cLight/gammam*t
			rt := math.Hypot(xt, yt)

			if rt > 0.0 {
				mother := cand
				c := cand.Clone()
				c.Pos = fmom.NewPxPyPzE(xt*1e3, yt*1e3, zt*1e3, cand.Pos.T()+t*cLight*1e3)
				c.Add(mother)

				output = append(output, *c)
				if math.Abs(q) > 1e-9 {
					pid := c.Pid
					if pid < 0 {
						pid = 0
					}
					switch pid {
					case 11:
						eles = append(eles, *c)
					case 13:
						muons = append(muons, *c)
					default:
						hadrons = append(hadrons, *c)
					}
				}

			}

		}
	}

	sort.Sort(ByPt(output))
	msg.Debugf(">>> output:     %v\n", len(output))
	return err
}

func init() {
	fwk.Register(reflect.TypeOf(Propagator{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &Propagator{
				TaskBase: fwk.NewTask(typ, name, mgr),
			}
			tsk.radius = 1.0
			err = tsk.DeclProp("Radius", &tsk.radius)
			if err != nil {
				return nil, err
			}
			tsk.radius2 = tsk.radius * tsk.radius

			tsk.halflen = 3.0
			err = tsk.DeclProp("HalfLength", &tsk.halflen)
			if err != nil {
				return nil, err
			}

			tsk.bz = 0.0
			err = tsk.DeclProp("Bz", &tsk.bz)
			if err != nil {
				return nil, err
			}

			tsk.input = "/fads/StableParticles"
			err = tsk.DeclProp("Input", &tsk.input)
			if err != nil {
				return nil, err
			}

			tsk.output = "StableParticles"
			err = tsk.DeclProp("Output", &tsk.output)
			if err != nil {
				return nil, err
			}

			tsk.hadrons = "ChargedHadrons"
			err = tsk.DeclProp("ChargedHadrons", &tsk.hadrons)
			if err != nil {
				return nil, err
			}

			tsk.eles = "Electrons"
			err = tsk.DeclProp("Electrons", &tsk.eles)
			if err != nil {
				return nil, err
			}

			tsk.muons = "Muons"
			err = tsk.DeclProp("Muons", &tsk.muons)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
