// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"fmt"
	"strings"

	"go-hep.org/x/hep/sio"
)

// VertexContainer is a collection of vertices
type VertexContainer struct {
	Flags  Flags
	Params Params
	Vtxs   []Vertex
}

type Vertex struct {
	Primary int32      // primary vertex of the event
	AlgType int32      // algorithm type
	Chi2    float32    // Chi^2 of vertex
	Prob    float32    // probability of the fit
	Pos     [3]float32 // position of the vertex (Px,Py,Pz)
	Cov     [6]float32 // covariance matrix
	Params  []float32
	RecPart *RecParticle // reconstructed particle associated to the vertex
}

func (vtx *Vertex) AlgName() string {
	// FIXME(sbinet)
	return fmt.Sprintf("Unknown (id=%d)", vtx.AlgType)
}

func (vtxs VertexContainer) String() string {
	o := new(strings.Builder)
	fmt.Fprintf(o, "%[1]s print out of Vertex collection %[1]s\n\n", strings.Repeat("-", 15))
	fmt.Fprintf(o, "%v", vtxs.Params)

	fmt.Fprintf(o, "\n")

	const (
		head = " [   id   ] |pri|     alg. type     |    chi2   |    prob.  |       position ( x, y, z)       | [par] |  [idRecP]  \n"
		tail = "------------|---|-------------------|-----------|-----------|---------------------------------|-------|------------\n"
	)
	o.WriteString(head)
	o.WriteString(tail)

	for _, vtx := range vtxs.Vtxs {
		fmt.Fprintf(o, " [%08d] | %d | %-17s | %+.2e | %+.2e | %+.2e, %+.2e, %+.2e | [%03d] | [%06d]\n",
			0, // id
			vtx.Primary, vtx.AlgName(), vtx.Chi2, vtx.Prob,
			vtx.Pos[0], vtx.Pos[1], vtx.Pos[2],
			len(vtx.Params),
			0, // vtx.RecPart
		)
	}

	return o.String()
}

func (*VertexContainer) VersionSio() uint32 {
	return Version
}

func (vtxs *VertexContainer) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&vtxs.Flags)
	enc.Encode(&vtxs.Params)
	enc.Encode(int32(len(vtxs.Vtxs)))
	for i := range vtxs.Vtxs {
		vtx := &vtxs.Vtxs[i]
		enc.Encode(&vtx.Primary)
		enc.Encode(&vtx.AlgType)
		enc.Encode(&vtx.Chi2)
		enc.Encode(&vtx.Prob)
		enc.Encode(&vtx.Pos)
		enc.Encode(&vtx.Cov)
		enc.Encode(&vtx.Params)
		enc.Pointer(&vtx.RecPart)
		enc.Tag(vtx)
	}
	return enc.Err()
}

func (vtxs *VertexContainer) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&vtxs.Flags)
	dec.Decode(&vtxs.Params)
	var n int32
	dec.Decode(&n)
	vtxs.Vtxs = make([]Vertex, int(n))
	for i := range vtxs.Vtxs {
		vtx := &vtxs.Vtxs[i]
		dec.Decode(&vtx.Primary)
		dec.Decode(&vtx.AlgType)
		dec.Decode(&vtx.Chi2)
		dec.Decode(&vtx.Prob)
		dec.Decode(&vtx.Pos)
		dec.Decode(&vtx.Cov)
		dec.Decode(&vtx.Params)
		dec.Pointer(&vtx.RecPart)
		dec.Tag(vtx)
	}
	return dec.Err()
}

var (
	_ sio.Versioner = (*VertexContainer)(nil)
	_ sio.Codec     = (*VertexContainer)(nil)
)
