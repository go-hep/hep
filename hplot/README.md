hplot
====

[![GoDoc](https://godoc.org/go-hep.org/x/hep/hplot?status.svg)](https://godoc.org/go-hep.org/x/hep/hplot)

`hplot` is a WIP package relying on `gonum/plot` to plot histograms,
n-tuples and functions.

## Installation

```sh
$ go get go-hep.org/x/hep/hplot
```

## Documentation

Is available on ``godoc``:

https://godoc.org/go-hep.org/x/hep/hplot


## Examples

### 1D histogram

![hist-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/h1d_plot_golden.png)

[embedmd]:# (example_h1d_test.go go /func ExampleH1D/ /\n}/)
```go
func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_plot.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

### 1D histogram with y-error bars

![hist-yerrs-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/h1d_yerrs_golden.png)

[embedmd]:# (example_h1d_test.go go /func ExampleH1D_withYErrBars/ /\n}/)
```go
func ExampleH1D_withYErrBars() {
	const npoints = 100

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist,
		hplot.WithHInfo(hplot.HInfoSummary),
		hplot.WithYErrBars(true),
	)
	h.YErrs.LineStyle.Color = color.RGBA{R: 255, A: 255}
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_yerrs.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

### 1D histogram with y-error bars, no lines

![hist-glyphs-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/h1d_glyphs_golden.png)

[embedmd]:# (example_h1d_test.go go /func ExampleH1D_withYErrBarsAndData/ /\n}/)
```go
func ExampleH1D_withYErrBarsAndData() {
	const npoints = 100

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.Legend.Top = true
	p.Legend.Left = true

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist,
		hplot.WithHInfo(hplot.HInfoSummary),
		hplot.WithYErrBars(true),
		hplot.WithGlyphStyle(draw.GlyphStyle{
			Shape:  draw.CrossGlyph{},
			Color:  color.Black,
			Radius: vg.Points(2),
		}),
	)
	h.GlyphStyle.Shape = draw.CircleGlyph{}
	h.YErrs.LineStyle.Color = color.Black
	h.LineStyle.Width = 0 // disable histogram lines
	p.Add(h)
	p.Legend.Add("data", h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)
	p.Legend.Add("model", norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_glyphs.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

### Tiles of 1D histograms

![tiled-plot](https://github.com/go-hep/hep/raw/master/hplot/testdata/tiled_plot_histogram_golden.png)

[embedmd]:# (example_tiledplot_test.go go /func ExampleTiledPlot/ /\n}/)
```go
func ExampleTiledPlot() {
	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 3, Rows: 2})

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	newHist := func(p *hplot.Plot) {
		const npoints = 10000
		hist := hbook.NewH1D(20, -4, +4)
		for i := 0; i < npoints; i++ {
			v := dist.Rand()
			hist.Fill(v, 1)
		}

		h := hplot.NewH1D(hist)
		p.Add(h)
	}

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(i, j)
			p.X.Min = -5
			p.X.Max = +5
			newHist(p)
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", i, j)
		}
	}

	// remove plot at (0,1)
	tp.Plots[1] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
```

![tiled-plot-aligned](https://github.com/go-hep/hep/raw/master/hplot/testdata/tiled_plot_aligned_histogram_golden.png)

[embedmd]:# (example_tiledplot_test.go go /func ExampleTiledPlot_align/ /\n}/)
```go
func ExampleTiledPlot_align() {
	tp := hplot.NewTiledPlot(draw.Tiles{
		Cols: 3, Rows: 3,
		PadX: 20, PadY: 20,
	})
	tp.Align = true

	points := func(i, j int) []hbook.Point2D {
		n := i*tp.Tiles.Cols + j + 1
		i += 1
		j = int(math.Pow(10, float64(n)))

		var pts []hbook.Point2D
		for ii := 0; ii < 10; ii++ {
			pts = append(pts, hbook.Point2D{
				X: float64(i + ii),
				Y: float64(j + ii + 1),
			})
		}
		return pts

	}

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(i, j)
			p.X.Min = -5
			p.X.Max = +5
			s := hplot.NewS2D(hbook.NewS2D(points(i, j)...))
			s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
			s.GlyphStyle.Radius = vg.Points(4)
			p.Add(s)

			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", i, j)
		}
	}

	// remove plot at (1,1)
	tp.Plots[4] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_aligned_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
```

### Subplots

![sub-plot](https://github.com/go-hep/hep/raw/master/hplot/testdata/sub_plot_golden.png)

https://godoc.org/go-hep.org/x/hep/hplot#example-package--Subplot

### Ratio-plots

![ratio-plot](https://github.com/go-hep/hep/raw/master/hplot/testdata/diff_plot_golden.png)

[embedmd]:# (example_ratioplot_test.go go /func ExampleRatioPlot/ /\n}/)
```go
func ExampleRatioPlot() {

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	hist1 := hbook.NewH1D(20, -4, +4)
	hist2 := hbook.NewH1D(20, -4, +4)

	for i := 0; i < npoints; i++ {
		v1 := dist.Rand() - 0.5
		v2 := dist.Rand() + 0.5
		hist1.Fill(v1, 1)
		hist2.Fill(v2, 1)
	}

	rp := hplot.NewRatioPlot()
	rp.Ratio = 0.3

	// Make a plot and set its title.
	rp.Top.Title.Text = "Histos"
	rp.Top.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h1 := hplot.NewH1D(hist1)
	h1.FillColor = color.NRGBA{R: 255, A: 100}
	rp.Top.Add(h1)

	h2 := hplot.NewH1D(hist2)
	h2.FillColor = color.NRGBA{B: 255, A: 100}
	rp.Top.Add(h2)

	rp.Top.Add(hplot.NewGrid())

	hist3 := hbook.NewH1D(20, -4, +4)
	for i := 0; i < hist3.Len(); i++ {
		v1 := hist1.Value(i)
		v2 := hist2.Value(i)
		x1, _ := hist1.XY(i)
		hist3.Fill(x1, v1-v2)
	}

	hdiff := hplot.NewH1D(hist3)

	rp.Bottom.X.Label.Text = "X"
	rp.Bottom.Y.Label.Text = "Delta-Y"
	rp.Bottom.Add(hdiff)
	rp.Bottom.Add(hplot.NewGrid())

	const (
		width  = 15 * vg.Centimeter
		height = width / math.Phi
	)

	err := hplot.Save(rp, width, height, "testdata/diff_plot.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}
```

### LaTeX-plots

[latex-plot (PDF)](https://github.com/go-hep/hep/raw/master/hplot/testdata/latex_plot_golden.pdf)

https://godoc.org/go-hep.org/x/hep/hplot#example-package--Latexplot

### 2D histogram

[embedmd]:# (example_h2d_test.go go /func ExampleH2D/ /\n}/)
```go
func ExampleH2D() {
	h := hbook.NewH2D(100, -10, 10, 100, -10, 10)

	const npoints = 10000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		h.Fill(v[0], v[1], 1)
	}

	p := hplot.New()
	p.Title.Text = "Hist-2D"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	p.Add(hplot.NewH2D(h, nil))
	p.Add(plotter.NewGrid())
	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/h2d_plot.png")
	if err != nil {
		log.Fatal(err)
	}
}
```
![h2d-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/h2d_plot_golden.png)

### Scatter2D

[embedmd]:# (example_s2d_test.go go /func ExampleS2D/ /\n}/)
```go
func ExampleS2D() {
	const npoints = 1000

	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	s2d := hbook.NewS2D()

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		s2d.Fill(hbook.Point2D{X: v[0], Y: v[1]})
	}

	p := hplot.New()
	p.Title.Text = "Scatter-2D"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s := hplot.NewS2D(s2d)
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(2)

	p.Add(s)

	err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/s2d.png")
	if err != nil {
		log.Fatal(err)
	}
}
```
![s2d-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/s2d_golden.png)
![s2d-errbars-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/s2d_errbars_golden.png)
![s2d-band-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/s2d_band_golden.png)

### Vertical lines

[embedmd]:# (example_line_test.go go /func ExampleVLine/ /\n}/)
```go
func ExampleVLine() {
	p := hplot.New()
	p.Title.Text = "vlines"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	var (
		left  = color.RGBA{B: 255, A: 255}
		right = color.RGBA{R: 255, A: 255}
	)

	p.Add(
		hplot.VLine(2.5, left, nil),
		hplot.VLine(5, nil, nil),
		hplot.VLine(7.5, nil, right),
	)

	err := p.Save(10*vg.Centimeter, -1, "testdata/vline.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
```
![vline-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/vline_golden.png)

### Horizontal lines

[embedmd]:# (example_line_test.go go /func ExampleHLine/ /\n}/)
```go
func ExampleHLine() {
	p := hplot.New()
	p.Title.Text = "hlines"
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	var (
		top    = color.RGBA{B: 255, A: 255}
		bottom = color.RGBA{R: 255, A: 255}
	)

	p.Add(
		hplot.HLine(2.5, nil, bottom),
		hplot.HLine(5, nil, nil),
		hplot.HLine(7.5, top, nil),
	)

	err := p.Save(10*vg.Centimeter, -1, "testdata/hline.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
```
![hline-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/hline_golden.png)

### Band between lines

[embedmd]:# (example_band_test.go go /func ExampleBand/ /\n}/)
```go
func ExampleBand() {
	const (
		npoints = 100
		xmax    = 10
	)

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	topData := make(plotter.XYs, npoints)
	botData := make(plotter.XYs, npoints)

	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		x := float64(i+1) / xmax

		v1 := dist.Rand()
		v2 := dist.Rand()

		topData[i].X = x
		topData[i].Y = 1/x + v1 + 10

		botData[i].X = x
		botData[i].Y = math.Log(x) + v2
	}

	top, err := hplot.NewLine(topData)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	top.LineStyle.Color = color.RGBA{R: 255, A: 255}

	bot, err := hplot.NewLine(botData)
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	bot.LineStyle.Color = color.RGBA{B: 255, A: 255}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 2})

	tp.Plots[0].Title.Text = "Band"
	tp.Plots[0].Add(
		top,
		bot,
		hplot.NewBand(color.Gray{200}, topData, botData),
	)

	tp.Plots[1].Title.Text = "Band"
	var (
		blue = color.RGBA{B: 255, A: 255}
		grey = color.Gray{200}
		band = hplot.NewBand(grey, topData, botData)
	)
	band.LineStyle = plotter.DefaultLineStyle
	band.LineStyle.Color = blue
	tp.Plots[1].Add(band)

	err = tp.Save(10*vg.Centimeter, -1, "testdata/band.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
```
![band-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/band_golden.png)

### Plot with borders

One can specify extra-space between the image borders (the physical file canvas) and the actual plot data.

![plot-border-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/h1d_borders_golden.png)

[embedmd]:# (example_h1d_test.go go /func ExampleH1D_withPlotBorders/ /\n}/)
```go
func ExampleH1D_withPlotBorders() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	hist := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		hist.Fill(v, 1)
	}

	// normalize histogram
	area := 0.0
	for _, bin := range hist.Binning.Bins {
		area += bin.SumW() * bin.XWidth()
	}
	hist.Scale(1 / area)

	// Make a plot and set its title.
	p := hplot.New()
	p.Title.Text = "Histogram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Create a histogram of our values drawn
	// from the standard normal.
	h := hplot.NewH1D(hist)
	h.Infos.Style = hplot.HInfoSummary
	p.Add(h)

	// The normal distribution function
	norm := hplot.NewFunction(dist.Prob)
	norm.Color = color.RGBA{R: 255, A: 255}
	norm.Width = vg.Points(2)
	p.Add(norm)

	// draw a grid
	p.Add(hplot.NewGrid())

	fig := hplot.Figure(p,
		hplot.WithDPI(96),
		hplot.WithBorder(hplot.Border{
			Right:  25,
			Left:   20,
			Top:    25,
			Bottom: 20,
		}),
	)

	// Save the plot to a PNG file.
	if err := hplot.Save(fig, 6*vg.Inch, -1, "testdata/h1d_borders.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

### Stack of 1D histograms

![hstack-example](https://github.com/go-hep/hep/raw/master/hplot/testdata/hstack_golden.png)

[embedmd]:# (example_hstack_test.go go /func ExampleHStack/ /\n}/)
```go
func ExampleHStack() {
	h1 := hbook.NewH1D(100, -10, 10)
	h2 := hbook.NewH1D(100, -10, 10)
	h3 := hbook.NewH1D(100, -10, 10)

	const seed = 1234
	fillH1(h1, 10000, -2, 1, seed)
	fillH1(h2, 10000, +3, 3, seed)
	fillH1(h3, 10000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}

	hh1 := hplot.NewH1D(h1)
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black

	hh2 := hplot.NewH1D(h2)
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3)
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black

	hs := []*hplot.H1D{hh1, hh2, hh3}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 1, Rows: 3})
	tp.Align = true

	{
		p := tp.Plots[0]
		p.Title.Text = "Histograms"
		p.Y.Label.Text = "Y"
		p.Add(hh1, hh2, hh3, hplot.NewGrid())
		p.Legend.Add("h1", hh1)
		p.Legend.Add("h2", hh2)
		p.Legend.Add("h3", hh3)
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "HStack - stack: OFF"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs)
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(2, 0)
		p.Title.Text = "Hstack - stack: ON"
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithLogY(false))
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(15*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}

}
```

