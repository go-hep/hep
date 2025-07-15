// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rbytes

import (
	"encoding/binary"
	"math"

	"go-hep.org/x/hep/groot/rvers"
)

func (w *WBuffer) WriteArrayU16(sli []uint16) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 2)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 2
		cur = end
		binary.BigEndian.PutUint16(w.w.p[beg:end], v)

	}
	w.w.c += 2 * len(sli)
}

func (w *WBuffer) WriteU16(v uint16) {
	if w.err != nil {
		return
	}
	w.w.grow(2)
	beg := w.w.c
	end := w.w.c + 2
	binary.BigEndian.PutUint16(w.w.p[beg:end], v)

	w.w.c += 2
}

func (w *WBuffer) WriteStdVectorU16(sli []uint16) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<uint16>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 2)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 2
		cur = end
		binary.BigEndian.PutUint16(w.w.p[beg:end], v)

	}
	w.w.c += 2 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayU32(sli []uint32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], v)

	}
	w.w.c += 4 * len(sli)
}

func (w *WBuffer) WriteU32(v uint32) {
	if w.err != nil {
		return
	}
	w.w.grow(4)
	beg := w.w.c
	end := w.w.c + 4
	binary.BigEndian.PutUint32(w.w.p[beg:end], v)

	w.w.c += 4
}

func (w *WBuffer) WriteStdVectorU32(sli []uint32) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<uint32>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], v)

	}
	w.w.c += 4 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayU64(sli []uint64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], v)

	}
	w.w.c += 8 * len(sli)
}

func (w *WBuffer) WriteU64(v uint64) {
	if w.err != nil {
		return
	}
	w.w.grow(8)
	beg := w.w.c
	end := w.w.c + 8
	binary.BigEndian.PutUint64(w.w.p[beg:end], v)

	w.w.c += 8
}

func (w *WBuffer) WriteStdVectorU64(sli []uint64) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<uint64>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], v)

	}
	w.w.c += 8 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayI16(sli []int16) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 2)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 2
		cur = end
		binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	}
	w.w.c += 2 * len(sli)
}

func (w *WBuffer) WriteI16(v int16) {
	if w.err != nil {
		return
	}
	w.w.grow(2)
	beg := w.w.c
	end := w.w.c + 2
	binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	w.w.c += 2
}

func (w *WBuffer) WriteStdVectorI16(sli []int16) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<int16>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 2)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 2
		cur = end
		binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	}
	w.w.c += 2 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayI32(sli []int32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], uint32(v))
	}
	w.w.c += 4 * len(sli)
}

func (w *WBuffer) WriteI32(v int32) {
	if w.err != nil {
		return
	}
	w.w.grow(4)
	beg := w.w.c
	end := w.w.c + 4
	binary.BigEndian.PutUint32(w.w.p[beg:end], uint32(v))
	w.w.c += 4
}

func (w *WBuffer) WriteStdVectorI32(sli []int32) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<int32>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], uint32(v))
	}
	w.w.c += 4 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayI64(sli []int64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], uint64(v))
	}
	w.w.c += 8 * len(sli)
}

func (w *WBuffer) WriteI64(v int64) {
	if w.err != nil {
		return
	}
	w.w.grow(8)
	beg := w.w.c
	end := w.w.c + 8
	binary.BigEndian.PutUint64(w.w.p[beg:end], uint64(v))
	w.w.c += 8
}

func (w *WBuffer) WriteStdVectorI64(sli []int64) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<int64>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], uint64(v))
	}
	w.w.c += 8 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayF32(sli []float32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], math.Float32bits(v))
	}
	w.w.c += 4 * len(sli)
}

func (w *WBuffer) WriteF32(v float32) {
	if w.err != nil {
		return
	}
	w.w.grow(4)
	beg := w.w.c
	end := w.w.c + 4
	binary.BigEndian.PutUint32(w.w.p[beg:end], math.Float32bits(v))
	w.w.c += 4
}

func (w *WBuffer) WriteStdVectorF32(sli []float32) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<float32>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 4)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 4
		cur = end
		binary.BigEndian.PutUint32(w.w.p[beg:end], math.Float32bits(v))
	}
	w.w.c += 4 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}

func (w *WBuffer) WriteArrayF64(sli []float64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], math.Float64bits(v))
	}
	w.w.c += 8 * len(sli)
}

func (w *WBuffer) WriteF64(v float64) {
	if w.err != nil {
		return
	}
	w.w.grow(8)
	beg := w.w.c
	end := w.w.c + 8
	binary.BigEndian.PutUint64(w.w.p[beg:end], math.Float64bits(v))
	w.w.c += 8
}

func (w *WBuffer) WriteStdVectorF64(sli []float64) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<float64>", rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(sli)))
	w.w.grow(len(sli) * 8)

	cur := w.w.c
	for _, v := range sli {
		beg := cur
		end := cur + 8
		cur = end
		binary.BigEndian.PutUint64(w.w.p[beg:end], math.Float64bits(v))
	}
	w.w.c += 8 * len(sli)

	if w.err != nil {
		return
	}
	_, w.err = w.SetHeader(hdr)
}
