// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes // import "go-hep.org/x/hep/groot/rbytes"

import "testing"

func TestResizeBool(t *testing.T) {
	{
		sli := ResizeBool(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeBool(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeBool(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeBool(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeU8(t *testing.T) {
	{
		sli := ResizeU8(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeU8(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU8(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU8(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeU16(t *testing.T) {
	{
		sli := ResizeU16(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeU16(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU16(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU16(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeU32(t *testing.T) {
	{
		sli := ResizeU32(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeU32(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU32(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU32(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeU64(t *testing.T) {
	{
		sli := ResizeU64(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeU64(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU64(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeU64(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeI8(t *testing.T) {
	{
		sli := ResizeI8(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeI8(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI8(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI8(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeI16(t *testing.T) {
	{
		sli := ResizeI16(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeI16(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI16(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI16(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 16; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeI32(t *testing.T) {
	{
		sli := ResizeI32(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeI32(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI32(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI32(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeI64(t *testing.T) {
	{
		sli := ResizeI64(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeI64(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI64(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeI64(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeF32(t *testing.T) {
	{
		sli := ResizeF32(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeF32(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF32(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF32(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeF64(t *testing.T) {
	{
		sli := ResizeF64(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeF64(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF64(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF64(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeF16(t *testing.T) {
	{
		sli := ResizeF16(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeF16(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF16(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeF16(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 12; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeD32(t *testing.T) {
	{
		sli := ResizeD32(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeD32(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeD32(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeD32(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}

func TestResizeStr(t *testing.T) {
	{
		sli := ResizeStr(nil, 0)
		if sli != nil {
			t.Fatalf("expected a nil slice")
		}
	}

	sli := ResizeStr(nil, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeStr(sli, 10)
	if got, want := len(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	sli = ResizeStr(sli, 5)
	if got, want := len(sli), 5; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
	if got, want := cap(sli), 10; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}
}
