// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes // import "go-hep.org/x/hep/groot/rbytes"

import (
	"go-hep.org/x/hep/groot/root"
)

func ResizeBool(sli []bool, n int) []bool {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]bool, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeI8(sli []int8, n int) []int8 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]int8, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeI16(sli []int16, n int) []int16 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]int16, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeI32(sli []int32, n int) []int32 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]int32, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeI64(sli []int64, n int) []int64 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]int64, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeU8(sli []uint8, n int) []uint8 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]uint8, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeU16(sli []uint16, n int) []uint16 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]uint16, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeU32(sli []uint32, n int) []uint32 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]uint32, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeU64(sli []uint64, n int) []uint64 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]uint64, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeF32(sli []float32, n int) []float32 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]float32, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeF64(sli []float64, n int) []float64 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]float64, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeD32(sli []root.Double32, n int) []root.Double32 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]root.Double32, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeF16(sli []root.Float16, n int) []root.Float16 {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]root.Float16, n-m)...)
	}
	sli = sli[:n]
	return sli
}

func ResizeStr(sli []string, n int) []string {
	if m := cap(sli); m < n {
		sli = sli[:m]
		sli = append(sli, make([]string, n-m)...)
	}
	sli = sli[:n]
	return sli
}
