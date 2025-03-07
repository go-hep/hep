// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio_test

import (
	"compress/flate"
	"reflect"
	"testing"

	"go-hep.org/x/hep/lcio"
)

func TestRWCluster(t *testing.T) {
	for _, test := range []struct {
		fname   string
		complvl int
	}{
		{"testdata/cluster.slcio", flate.NoCompression},
		{"testdata/cluster-compressed.slcio", flate.BestCompression},
	} {
		testRWCluster(t, test.complvl, test.fname)
	}
}

func testRWCluster(t *testing.T, compLevel int, fname string) {
	w, err := lcio.Create(fname)
	if err != nil {
		t.Error(err)
		return
	}
	defer w.Close()
	w.SetCompressionLevel(compLevel)

	const (
		nevents = 10
		N       = 100
		CLUS    = "Clusters"
		CLUSERR = "ClustersWithEnergyError"
	)

	for i := range nevents {
		evt := lcio.Event{
			RunNumber:   4711,
			EventNumber: int32(i),
		}

		var (
			clus    = lcio.ClusterContainer{Flags: lcio.BitsRChLong}
			clusErr = lcio.ClusterContainer{Flags: lcio.BitsRChLong | lcio.BitsRChEnergyError}
		)
		for j := range N {
			clu := lcio.Cluster{
				Energy: float32(i*j) + 117,
				Pos:    [3]float32{float32(i), float32(j), float32(i * j)},
			}
			clus.Clusters = append(clus.Clusters, clu)

			clu.EnergyErr = float32(i*j) * 0.117
			clusErr.Clusters = append(clusErr.Clusters, clu)
		}

		evt.Add(CLUS, &clus)
		evt.Add(CLUSERR, &clusErr)

		err = w.WriteEvent(&evt)
		if err != nil {
			t.Errorf("%s: error writing event %d: %v", fname, i, err)
			return
		}
	}

	err = w.Close()
	if err != nil {
		t.Errorf("%s: error closing file: %v", fname, err)
		return
	}

	r, err := lcio.Open(fname)
	if err != nil {
		t.Errorf("%s: error opening file: %v", fname, err)
	}
	defer r.Close()

	for i := range nevents {
		if !r.Next() {
			t.Errorf("%s: error reading event %d", fname, i)
			return
		}
		evt := r.Event()
		if got, want := evt.RunNumber, int32(4711); got != want {
			t.Errorf("%s: run-number error. got=%d. want=%d", fname, got, want)
			return
		}

		if got, want := evt.EventNumber, int32(i); got != want {
			t.Errorf("%s: run-number error. got=%d. want=%d", fname, got, want)
			return
		}

		if !evt.Has(CLUS) {
			t.Errorf("%s: no %s collection", fname, CLUS)
			return
		}

		if !evt.Has(CLUSERR) {
			t.Errorf("%s: no %s collection", fname, CLUSERR)
			return
		}

		clus := evt.Get(CLUS).(*lcio.ClusterContainer)
		clusErr := evt.Get(CLUSERR).(*lcio.ClusterContainer)

		for j := range N {
			got := clus.Clusters[j]
			want := lcio.Cluster{
				Clusters:   []*lcio.Cluster{},
				Hits:       []*lcio.CalorimeterHit{},
				PIDs:       []lcio.ParticleID{},
				Shape:      []float32{},
				SubDetEnes: []float32{},
				Weights:    []float32{},
				Energy:     float32(i*j) + 117,
				Pos:        [3]float32{float32(i), float32(j), float32(i * j)},
			}

			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: event %d clu[%d]\ngot= %#v\nwant=%#v\n",
					fname, i, j,
					got, want,
				)
				return
			}

			want.EnergyErr = float32(i*j) * 0.117
			got = clusErr.Clusters[j]
			if !reflect.DeepEqual(got, want) {
				t.Errorf("%s: event %d clu[%d]\ngot= %#v\nwant=%#v\n",
					fname, i, j,
					got, want,
				)
				return
			}
		}
	}

	err = r.Close()
	if err != nil {
		t.Errorf("%s: error closing file: %v", fname, err)
		return
	}

}
