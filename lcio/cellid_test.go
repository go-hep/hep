// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"testing"
)

func TestCellIDDecoder(t *testing.T) {
	{
		dec := NewCellIDDecoder("layer:7,system:-3,barrel:3,theta:32:11,phi:11")
		hit := CalorimeterHit{
			CellID0: 0x00ffffaa,
			CellID1: 0x00aaffff,
		}

		for _, test := range []struct {
			name string
			want int64
		}{
			{"layer", 42},
			{"system", -1},
			{"barrel", 7},
			{"theta", 2047},
			{"phi", 1375},
		} {
			got := dec.Get(&hit, test.name)
			if got != test.want {
				t.Errorf("%s=%d. want=%d", test.name, got, test.want)
			}
		}

		if str, want := dec.ValueString(&hit), "layer:42,system:-1,barrel:7,theta:2047,phi:1375"; str != want {
			t.Errorf("value-string error.\n got=%q\nwant=%q", str, want)
		}
	}
	{
		codec := "M:3,S-1:3,I:9,J:9,K-1:6"
		dec := NewCellIDDecoder(codec)
		hit := CalorimeterHit{
			CellID0: 541065232,
		}
		if val, want := dec.Value(&hit), int64(hit.CellID0); val != want {
			t.Errorf("value=%d. want=%d", val, want)
		}
		if str, want := dec.ValueString(&hit), "M:0,S-1:2,I:0,J:128,K-1:32"; str != want {
			t.Errorf("value-string error.\n got=%q\nwant=%q", str, want)
		}

	}
}

func TestBitField(t *testing.T) {
	const codec = "layer:7,system:-3,barrel:3,theta:32:11,phi:11"
	bf := newBitField64(codec)
	if descr, want := bf.Description(), "layer:0:7,system:7:-3,barrel:10:3,theta:32:11,phi:43:11"; descr != want {
		t.Errorf("description error.\n got=%q\nwant=%q", descr, want)
	}
	if str, want := bf.valueString(), "layer:0,system:0,barrel:0,theta:0,phi:0"; str != want {
		t.Errorf("value-string error.\n got=%q\nwant=%q", str, want)
	}
	bf.value = 0x00aaffff00ffffaa
	if str, want := bf.valueString(), "layer:42,system:-1,barrel:7,theta:2047,phi:1375"; str != want {
		t.Errorf("value-string error.\n got=%q\nwant=%q", str, want)
	}
}
