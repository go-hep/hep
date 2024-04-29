// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"runtime"
	"testing"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/vgop"
	"go-hep.org/x/hep/internal/diff"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

func TestJSON(t *testing.T) {
	sr := font.Font{Typeface: "Liberation", Variant: "Serif"}
	tr := font.From(sr, 12)
	ft13 := hplot.DefaultStyle.Fonts.Cache.Lookup(tr, 13)
	ft20 := hplot.DefaultStyle.Fonts.Cache.Lookup(tr, 20)

	c := vgop.NewJSON(vgop.WithSize(10, 20))

	c.Push()
	c.SetLineWidth(2)
	c.SetLineDash([]vg.Length{1, 2}, 4)
	c.SetColor(color.RGBA{R: 255, A: 255})
	c.SetColor(color.Gray{Y: 100})
	c.Rotate(math.Pi / 2)
	c.Translate(vg.Point{X: 10, Y: 20})
	c.Scale(15, 25)
	c.Pop()
	p0 := vg.Path(nil)
	p1 := vg.Path([]vg.PathComp{
		{Type: vg.MoveComp, Pos: vg.Point{X: 1, Y: 2}},
		{Type: vg.LineComp, Pos: vg.Point{X: 2, Y: 3}},
		{Type: vg.ArcComp, Pos: vg.Point{X: 3, Y: 4}, Radius: 5, Start: 6, Angle: 7},
		{Type: vg.CurveComp, Pos: vg.Point{X: 4, Y: 5}, Control: []vg.Point{{X: 6, Y: 7}}},
		{Type: vg.CurveComp, Pos: vg.Point{X: 5, Y: 6}, Control: []vg.Point{{X: 7, Y: 8}, {X: 9, Y: 10}}},
		{Type: vg.CloseComp},
	})
	c.Stroke(p0)
	c.Stroke(p1)
	c.Fill(p0)
	c.Fill(p1)

	c.FillString(ft13, vg.Point{X: 10, Y: 20}, "hello\nworld")
	c.FillString(ft20, vg.Point{X: 20, Y: 30}, "BYE.")

	img := image.NewRGBA(image.Rect(0, 0, 20, 30))
	draw.Draw(img, img.Rect, image.NewUniform(color.RGBA{0x66, 0x66, 0x66, 0xff}), image.Point{}, draw.Src)

	c.DrawImage(vg.Rectangle{Min: vg.Point{X: 1, Y: 2}, Max: vg.Point{X: 3, Y: 4}}, img)

	f, err := os.Create("testdata/simple.json")
	if err != nil {
		t.Fatalf("could not create output JSON file: %+v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")

	err = enc.Encode(c)
	if err != nil {
		t.Fatalf("could not encode canvas: %+v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("could not close output JSON file: %+v", err)
	}

	err = diff.Files("testdata/simple.json", "testdata/simple_golden.json")
	if err != nil {
		t.Fatalf("JSON files differ:\n%s", err)
	}

	c = vgop.NewJSON()
	got, err := os.ReadFile("testdata/simple.json")
	if err != nil {
		t.Fatalf("could not read-back JSON file: %+v", err)
	}
	dec := json.NewDecoder(bytes.NewReader(got))
	err = dec.Decode(c)
	if err != nil {
		t.Fatalf("could not decode JSON canvas: %+v", err)
	}

	bak := new(bytes.Buffer)
	enc = json.NewEncoder(bak)
	enc.SetIndent("", "  ")

	err = enc.Encode(c)
	if err != nil {
		t.Fatalf("could not re-encode JSON canvas: %+v", err)
	}

	if got, want := bak.String(), string(got); got != want {
		o := diff.Format(got, want)
		t.Fatalf("JSON roundtrip failed:\n%s", o)
	}

	defer os.Remove("testdata/simple.json")
}

func TestSaveJSON(t *testing.T) {
	p := hplot.New()
	p.Title.Text = "Title"
	p.X.Min = -1
	p.X.Max = +1
	p.X.Label.Text = "X"
	p.Y.Min = -10
	p.Y.Max = +10
	p.Y.Label.Text = "Y"

	err := hplot.Save(p, 10*vg.Centimeter, 20*vg.Centimeter, "testdata/plot.json")
	if err != nil {
		t.Fatalf("could not save plot to JSON: %+v", err)
	}

	err = diff.Files("testdata/plot.json", "testdata/plot_golden.json")
	if err != nil {
		fatalf := t.Fatalf
		if runtime.GOOS == "darwin" {
			// ignore errors for darwin and mac-silicon
			fatalf = t.Logf
		}
		fatalf("JSON files differ:\n%s", err)
	}

	defer os.Remove("testdata/plot.json")
}
