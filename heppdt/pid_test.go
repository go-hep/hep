package heppdt_test

import (
	"testing"

	"go-hep.org/x/hep/heppdt"
)

func TestPID(t *testing.T) {
	/*
		ids := []int{ 5, 25, 15, 213, -3214, 10213, 9050225, -200543, 129050225,
			2000025, 3101, 3301, -2212, 1000020040, -1000060120, 555,
			5000040, 5100005, 24, 5100024, 5100025, 9221132,
			4111370, -4120240, 4110050, 10013730,
			1000993, 1000612, 1000622, 1000632, 1006213, 1000652,
			1009113, 1009213, 1009323,
			1093114, 1009333, 1006313, 1092214, 1006223,
		}
	*/

	for _, table := range []struct {
		id                            int
		nx, nr, n1, nq1, nq2, nq3, nj int
		extra                         int
		jspin, lspin, spin            int
		fid                           int
		charge                        float64
		valid                         bool

		ishadron, ismeson, isrhadron, isbaryon bool
		islepton, isdiquark, isnucleus         bool

		issusy, isdyon, isqlb, isunknown bool

		hasup, hasdown, hasstrange, hascharm, hasbottom, hastop bool
	}{
		{
			id: 5,
			nx: 0, nr: 0, n1: 0, nq1: 0, nq2: 0, nq3: 0, nj: 5,
			extra: 0,
			jspin: 2, lspin: 0, spin: 0,
			fid:    0,
			charge: -1. / 3.,
			valid:  true,
			isqlb:  true,
		},
		{
			id: 25,
			nx: 0, nr: 0, n1: 0, nq1: 0, nq2: 0, nq3: 2, nj: 5,
			extra: 0,
			jspin: 0, lspin: 0, spin: 0,
			fid:    0,
			charge: 0,
			valid:  true,
			isqlb:  true,
		},
		{
			id: 15,
			nx: 0, nr: 0, n1: 0, nq1: 0, nq2: 0, nq3: 1, nj: 5,
			extra: 0,
			jspin: 2, lspin: 0, spin: 0,
			fid:      0,
			charge:   -1,
			valid:    true,
			isqlb:    true,
			islepton: true,
		},
		// {
		// 	id: 213,
		// 	nx: 0, nr:0, n1:0, nq1:0, nq2:2, nq3:1, nj:3,
		// 	extra: 0,
		// 	jspin: 3, lspin:0, spin:1,
		// 	fid: 0,
		// 	charge: 1,
		// 	valid: true,
		// 	ishadron: true,
		// 	ismeson: true,
		// 	hasup: true,
		// 	hasdown: true,
		// },
	} {
		id := table.id
		pid := heppdt.PID(id)
		nx := pid.Digit(heppdt.N)
		if nx != table.nx {
			t.Fatalf("pid=%d. expected nx=%d. got=%d", id, table.nx, nx)
		}

		nr := pid.Digit(heppdt.Nr)
		if nr != table.nr {
			t.Fatalf("pid=%d. expected nr=%d. got=%d", id, table.nr, nr)
		}

		nq1 := pid.Digit(heppdt.Nq1)
		if nq1 != table.nq1 {
			t.Fatalf("pid=%d. expected nq1=%d. got=%d", id, table.nq1, nq1)
		}

		nq2 := pid.Digit(heppdt.Nq2)
		if nq2 != table.nq2 {
			t.Fatalf("pid=%d. expected nq2=%d. got=%d", id, table.nq2, nq2)
		}

		nq3 := pid.Digit(heppdt.Nq3)
		if nq3 != table.nq3 {
			t.Fatalf("pid=%d. expected nq3=%d. got=%d", id, table.nq3, nq3)
		}

		extra := pid.ExtraBits()
		if extra != table.extra {
			t.Fatalf("pid=%d. expected extra=%d. got=%d", id, table.extra, extra)
		}

		jspin := pid.JSpin()
		if jspin != table.jspin {
			t.Fatalf("pid=%d. expected jspin=%d. got=%d", id, table.jspin, jspin)
		}

		lspin := pid.LSpin()
		if lspin != table.lspin {
			t.Fatalf("pid=%d. expected lspin=%d. got=%d", id, table.lspin, lspin)
		}

		spin := pid.SSpin()
		if spin != table.spin {
			t.Fatalf("pid=%d. expected sspin=%d. got=%d", id, table.spin, spin)
		}

		charge := pid.Charge()
		if charge != table.charge {
			t.Fatalf("pid=%d. expected charge=%v. got=%v", id, table.charge, charge)
		}

		if pid.IsValid() != table.valid {
			t.Fatalf("expected pid=%d valid=%v. got=%v", int(pid), table.valid, pid.IsValid())
		}

		if table.ishadron != pid.IsHadron() {
			t.Fatalf("expected pid=%d ishadron=%v. got=%v", int(pid), table.ishadron, pid.IsHadron())
		}

		if table.ismeson != pid.IsMeson() {
			t.Fatalf("expected pid=%d ismeson=%v. got=%v", int(pid), table.ismeson, pid.IsMeson())
		}

		if table.isrhadron != pid.IsRhadron() {
			t.Fatalf("expected pid=%d isrhadron=%v. got=%v", int(pid), table.isrhadron, pid.IsRhadron())
		}

		if table.isbaryon != pid.IsBaryon() {
			t.Fatalf("expected pid=%d isbaryon=%v. got=%v", int(pid), table.isbaryon, pid.IsBaryon())
		}

		if table.islepton != pid.IsLepton() {
			t.Fatalf("expected pid=%d islepton=%v. got=%v", int(pid), table.islepton, pid.IsLepton())
		}

		if table.isdiquark != pid.IsDiQuark() {
			t.Fatalf("expected pid=%d isdiquark=%v. got=%v", int(pid), table.isdiquark, pid.IsDiQuark())
		}

		if table.isnucleus != pid.IsNucleus() {
			t.Fatalf("expected pid=%d isnucleus=%v. got=%v", int(pid), table.isnucleus, pid.IsNucleus())
		}

		if table.issusy != pid.IsSUSY() {
			t.Fatalf("expected pid=%d issusy=%v. got=%v", int(pid), table.issusy, pid.IsSUSY())
		}

		if table.isdyon != pid.IsDyon() {
			t.Fatalf("expected pid=%d isdyon=%v. got=%v", int(pid), table.isdyon, pid.IsDyon())
		}

		if table.isqlb != (pid.AbsPID() < 100) {
			t.Fatalf("expected pid=%d is-quark/lepton/boson=%v. got=%v", int(pid), table.isqlb, pid.AbsPID() < 100)
		}

		if table.isunknown != (!pid.IsSUSY() && !pid.IsDyon() && !pid.IsHadron() &&
			!pid.IsLepton() && !pid.IsDiQuark() && !pid.IsNucleus() &&
			!(pid.AbsPID() < 100)) {
			t.Fatalf("expected pid=%d isunknown=%v. got=%v", int(pid), table.isunknown, !table.isunknown)
		}

		if table.hasup != pid.HasUp() {
			t.Fatalf("expected pid=%d hasup=%v. got=%v", int(pid), table.hasup, pid.HasUp())
		}

		if table.hasdown != pid.HasDown() {
			t.Fatalf("expected pid=%d hasdown=%v. got=%v", int(pid), table.hasdown, pid.HasDown())
		}

		if table.hasstrange != pid.HasStrange() {
			t.Fatalf("expected pid=%d hasstrange=%v. got=%v", int(pid), table.hasstrange, pid.HasStrange())
		}

		if table.hascharm != pid.HasCharm() {
			t.Fatalf("expected pid=%d hascharm=%v. got=%v", int(pid), table.hascharm, pid.HasCharm())
		}

		if table.hasbottom != pid.HasBottom() {
			t.Fatalf("expected pid=%d hasbottom=%v. got=%v", int(pid), table.hasbottom, pid.HasBottom())
		}

		if table.hastop != pid.HasTop() {
			t.Fatalf("expected pid=%d hastop=%v. got=%v", int(pid), table.hastop, pid.HasTop())
		}
	}
}
