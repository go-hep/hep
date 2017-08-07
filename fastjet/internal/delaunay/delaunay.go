// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math/rand"
)

type Delaunay struct {
}

func HierarchicalDelaunay() *Delaunay {
	panic(fmt.Errorf("delaunay: HierarchicalDelaunay not implemented"))
}

func WalkDelaunay(points []*Point, r *rand.Rand) *Delaunay {
	panic(fmt.Errorf("delaunay: WalkDelaunay not implemented"))
}

func (d *Delaunay) Triangles() []*Triangle {
	panic(fmt.Errorf("delaunay: Triangles not implemented"))
}

func (d *Delaunay) Insert(p *Point) (updatedNearestNeighbor []*Point) {
	panic(fmt.Errorf("delaunay: Insert not implemented"))
}

func (d *Delaunay) Remove(p *Point) (updatedNearestNeighbor []*Point) {
	panic(fmt.Errorf("delaunay: Remove not implemented"))
}
