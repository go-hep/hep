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

![hist-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h1d_plot_golden.png)

[embedmd]:# (h1d_example_test.go go /func ExampleH1D/ /\n}/)
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

![hist-yerrs-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h1d_yerrs_golden.png)

[embedmd]:# (h1d_example_test.go go /func ExampleH1D_withYErrBars/ /\n}/)
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

![hist-glyphs-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h1d_glyphs_golden.png)

[embedmd]:# (h1d_example_test.go go /func ExampleH1D_withYErrBarsAndData/ /\n}/)
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

### 1D histogram with y-error bars and error bands

![hist-yerrs-band-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h1d_yerrs_band_golden.png)

[embedmd]:# (h1d_example_test.go go /func ExampleH1D_withYErrBars_withBand/ /\n}/)
```go
func ExampleH1D_withYErrBars_withBand() {
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
		hplot.WithBand(true),
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
	if err := p.Save(6*vg.Inch, -1, "testdata/h1d_yerrs_band.png"); err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

### Tiles of 1D histograms

![tiled-plot](https://github.com/go-hep/hep/raw/main/hplot/testdata/tiled_plot_histogram_golden.png)

[embedmd]:# (tiledplot_example_test.go go /func ExampleTiledPlot/ /\n}/)
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
			p := tp.Plot(j, i)
			p.X.Min = -5
			p.X.Max = +5
			newHist(p)
			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", j, i)
		}
	}

	// remove plot at (1,0)
	tp.Plots[1] = nil

	err := tp.Save(15*vg.Centimeter, -1, "testdata/tiled_plot_histogram.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
```

![tiled-plot-aligned](https://github.com/go-hep/hep/raw/main/hplot/testdata/tiled_plot_aligned_histogram_golden.png)

[embedmd]:# (tiledplot_example_test.go go /func ExampleTiledPlot_align/ /\n}/)
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
			p := tp.Plot(j, i)
			p.X.Min = -5
			p.X.Max = +5
			s := hplot.NewS2D(hbook.NewS2D(points(i, j)...))
			s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
			s.GlyphStyle.Radius = vg.Points(4)
			p.Add(s)

			p.Title.Text = fmt.Sprintf("hist - (%02d, %02d)", j, i)
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

![sub-plot](https://github.com/go-hep/hep/raw/main/hplot/testdata/sub_plot_golden.png)

https://godoc.org/go-hep.org/x/hep/hplot#example-package--Subplot

### Ratio-plots

![ratio-plot](https://github.com/go-hep/hep/raw/main/hplot/testdata/diff_plot_golden.png)

[embedmd]:# (ratioplot_example_test.go go /func ExampleRatioPlot/ /\n}/)
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

[latex-plot (PDF)](https://github.com/go-hep/hep/raw/main/hplot/testdata/latex_plot_golden.pdf)

https://godoc.org/go-hep.org/x/hep/hplot#example-package--Latexplot

### 2D histogram

[embedmd]:# (h2d_example_test.go go /func ExampleH2D/ /\n}/)
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
![h2d-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h2d_plot_golden.png)

### Scatter2D

[embedmd]:# (s2d_example_test.go go /func ExampleS2D/ /\n}/)
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
![s2d-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/s2d_golden.png)
![s2d-errbars-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/s2d_errbars_golden.png)
![s2d-band-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/s2d_band_golden.png)
![s2d-steps-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/s2d_steps_golden.png)
![s2d-steps-band-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/s2d_steps_band_golden.png)

### Vertical lines

[embedmd]:# (line_example_test.go go /func ExampleVLine/ /\n}/)
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
![vline-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/vline_golden.png)

### Horizontal lines

[embedmd]:# (line_example_test.go go /func ExampleHLine/ /\n}/)
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
![hline-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/hline_golden.png)

### Band between lines

[embedmd]:# (band_example_test.go go /func ExampleBand/ /\n}/)
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
![band-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/band_golden.png)

### Plot with borders

One can specify extra-space between the image borders (the physical file canvas) and the actual plot data.

![plot-border-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/h1d_borders_golden.png)

[embedmd]:# (h1d_example_test.go go /func ExampleH1D_withPlotBorders/ /\n}/)
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

![hstack-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/hstack_golden.png)

[embedmd]:# (hstack_example_test.go go /func ExampleHStack/ /\n}/)
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
		p := tp.Plot(0, 1)
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
		p := tp.Plot(0, 2)
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

### Stack of 1D histograms with a band

![hstack-band-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/hstack_band_golden.png)

[embedmd]:# (hstack_example_test.go go /func ExampleHStack_withBand/ /\n}/)
```go
func ExampleHStack_withBand() {
	h1 := hbook.NewH1D(50, -8, 12)
	h2 := hbook.NewH1D(50, -8, 12)
	h3 := hbook.NewH1D(50, -8, 12)

	const seed = 1234
	fillH1(h1, 2000, -2, 1, seed)
	fillH1(h2, 2000, +3, 3, seed)
	fillH1(h3, 2000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}

	hh1 := hplot.NewH1D(h1, hplot.WithBand(true))
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black
	hh1.Band.FillColor = color.NRGBA{G: 210, A: 200}

	hh2 := hplot.NewH1D(h2, hplot.WithBand(false))
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3, hplot.WithBand(true))
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black
	hh3.Band.FillColor = color.NRGBA{R: 220, A: 200}

	hs := []*hplot.H1D{hh1, hh2, hh3}

	hh4 := hplot.NewH1D(h1)
	hh4.FillColor = colors[0]
	hh4.LineStyle.Color = color.Black

	hh5 := hplot.NewH1D(h2)
	hh5.FillColor = colors[1]
	hh5.LineStyle.Width = 0

	hh6 := hplot.NewH1D(h3)
	hh6.FillColor = colors[2]
	hh6.LineStyle.Color = color.Black

	hsHistoNoBand := []*hplot.H1D{hh4, hh5, hh6}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 2})
	tp.Align = true

	{
		p := tp.Plot(0, 0)
		p.Title.Text = "Histos With or Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true))
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "Histos Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true))
		hstack.Stack = hplot.HStackOff
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(0, 1)
		p.Title.Text = "Histos With or Without Band, Stack: ON"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true))
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 1)
		p.Title.Text = "Histos Without Band, Stack: ON"
		p.Title.Padding = 10
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true))
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(25*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack_band.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
```

### Stack of 1D histograms with a band, with a log-y scale

![hstack-logy-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/hstack_logy_golden.png)

[embedmd]:# (hstack_example_test.go go /func ExampleHStack_withLogY/ /\n}/)
```go
func ExampleHStack_withLogY() {
	h1 := hbook.NewH1D(50, -8, 12)
	h2 := hbook.NewH1D(50, -8, 12)
	h3 := hbook.NewH1D(50, -8, 12)

	const seed = 1234
	fillH1(h1, 2000, -2, 1, seed)
	fillH1(h2, 2000, +3, 3, seed)
	fillH1(h3, 2000, +4, 1, seed)

	colors := []color.Color{
		color.NRGBA{122, 195, 106, 150},
		color.NRGBA{90, 155, 212, 150},
		color.NRGBA{250, 167, 91, 150},
	}
	logy := hplot.WithLogY(true)

	hh1 := hplot.NewH1D(h1, hplot.WithBand(true), logy)
	hh1.FillColor = colors[0]
	hh1.LineStyle.Color = color.Black
	hh1.Band.FillColor = color.NRGBA{G: 210, A: 200}

	hh2 := hplot.NewH1D(h2, hplot.WithBand(false), logy)
	hh2.FillColor = colors[1]
	hh2.LineStyle.Width = 0

	hh3 := hplot.NewH1D(h3, hplot.WithBand(true), logy)
	hh3.FillColor = colors[2]
	hh3.LineStyle.Color = color.Black
	hh3.Band.FillColor = color.NRGBA{R: 220, A: 200}

	hs := []*hplot.H1D{hh1, hh2, hh3}

	hh4 := hplot.NewH1D(h1, logy)
	hh4.FillColor = colors[0]
	hh4.LineStyle.Color = color.Black

	hh5 := hplot.NewH1D(h2, logy)
	hh5.FillColor = colors[1]
	hh5.LineStyle.Width = 0

	hh6 := hplot.NewH1D(h3, logy)
	hh6.FillColor = colors[2]
	hh6.LineStyle.Color = color.Black

	hsHistoNoBand := []*hplot.H1D{hh4, hh5, hh6}

	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 2})
	tp.Align = true

	{
		p := tp.Plot(0, 0)
		p.Title.Text = "Histos With or Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true), logy)
		hstack.Stack = hplot.HStackOff
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 0)
		p.Title.Text = "Histos Without Band, Stack: OFF"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true), logy)
		hstack.Stack = hplot.HStackOff
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(0, 1)
		p.Title.Text = "Histos With or Without Band, Stack: ON"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hs, hplot.WithBand(true), logy)
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	{
		p := tp.Plot(1, 1)
		p.Title.Text = "Histos Without Band, Stack: ON"
		p.Title.Padding = 10
		p.Y.Scale = plot.LogScale{}
		p.Y.Tick.Marker = plot.LogTicks{}
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		hstack := hplot.NewHStack(hsHistoNoBand, hplot.WithBand(true), logy)
		hstack.Band.FillColor = color.NRGBA{R: 100, G: 100, B: 100, A: 200}
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true
	}

	err := tp.Save(25*vg.Centimeter, 15*vg.Centimeter, "testdata/hstack_logy.png")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
}
```

## Labels

![label-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/label_plot_golden.png)

[embedmd]:# (label_example_test.go go /func ExampleLabel/ /\n}/)
```go
func ExampleLabel() {

	// Creating a new plot
	p := hplot.New()
	p.Title.Text = "Plot labels"
	p.X.Min = -10
	p.X.Max = +10
	p.Y.Min = -10
	p.Y.Max = +10

	// Default labels
	l1 := hplot.NewLabel(-8, 5, "(-8,5)\nDefault label")
	p.Add(l1)

	// Label with normalized coordinates.
	l3 := hplot.NewLabel(
		0.5, 0.5,
		"(0.5,0.5)\nLabel with relative coords",
		hplot.WithLabelNormalized(true),
	)
	p.Add(l3)

	// Label with normalized coordinates and auto-adjustement.
	l4 := hplot.NewLabel(
		0.95, 0.95,
		"(0.95,0.95)\nLabel at the canvas edge, with AutoAdjust",
		hplot.WithLabelNormalized(true),
		hplot.WithLabelAutoAdjust(true),
	)
	p.Add(l4)

	// Label with a customed TextStyle
	usrFont := font.Font{
		Typeface: "Liberation",
		Variant:  "Mono",
		Weight:   xfnt.WeightBold,
		Style:    xfnt.StyleNormal,
		Size:     12,
	}
	sty := text.Style{
		Color: plotutil.Color(2),
		Font:  usrFont,
	}
	l5 := hplot.NewLabel(
		0.0, 0.1,
		"(0.0,0.1)\nLabel with a user-defined font",
		hplot.WithLabelTextStyle(sty),
		hplot.WithLabelNormalized(true),
	)
	p.Add(l5)

	p.Add(plotter.NewGlyphBoxes())
	p.Add(hplot.NewGrid())

	// Save the plot to a PNG file.
	err := p.Save(15*vg.Centimeter, -1, "testdata/label_plot.png")
	if err != nil {
		log.Fatalf("error saving plot: %v\n", err)
	}
}
```

## Time series

![timeseries-example](https://github.com/go-hep/hep/raw/main/hplot/testdata/timeseries_monthly_golden.png)

[embedmd]:# (ticks_example_test.go go /func ExampleTicks_monthly/ /\n}/)
```go
func ExampleTicks_monthly() {
	cnv := epok.UTCUnixTimeConverter{}

	p := hplot.New()
	p.Title.Text = "Time series (monthly)"
	p.Y.Label.Text = "Goroutines"

	p.Y.Min = 0
	p.Y.Max = 4
	p.X.AutoRescale = true
	p.X.Tick.Marker = epok.Ticks{
		Ruler: epok.Rules{
			Major: epok.Rule{
				Freq:  epok.Monthly,
				Range: epok.RangeFrom(1, 13, 2),
			},
		},
		Format:    "2006\nJan-02\n15:04:05",
		Converter: cnv,
	}

	xysFrom := func(vs ...float64) plotter.XYs {
		o := make(plotter.XYs, len(vs))
		for i := range o {
			o[i].X = vs[i]
			o[i].Y = float64(i + 1)
		}
		return o
	}
	data := xysFrom(
		cnv.FromTime(parse("2010-01-02 01:02:03")),
		cnv.FromTime(parse("2010-02-01 01:02:03")),
		cnv.FromTime(parse("2010-02-04 11:22:33")),
		cnv.FromTime(parse("2010-03-04 01:02:03")),
		cnv.FromTime(parse("2010-04-05 01:02:03")),
		cnv.FromTime(parse("2010-04-05 01:02:03")),
		cnv.FromTime(parse("2010-05-01 00:02:03")),
		cnv.FromTime(parse("2010-05-04 04:04:04")),
		cnv.FromTime(parse("2010-05-08 11:12:13")),
		cnv.FromTime(parse("2010-06-15 01:02:03")),
		cnv.FromTime(parse("2010-07-04 04:04:43")),
		cnv.FromTime(parse("2010-07-14 14:17:09")),
		cnv.FromTime(parse("2010-08-04 21:22:23")),
		cnv.FromTime(parse("2010-08-15 11:12:13")),
		cnv.FromTime(parse("2010-09-01 21:52:53")),
		cnv.FromTime(parse("2010-10-25 01:19:23")),
		cnv.FromTime(parse("2010-11-30 11:32:53")),
		cnv.FromTime(parse("2010-12-24 23:59:59")),
		cnv.FromTime(parse("2010-12-31 23:59:59")),
		cnv.FromTime(parse("2011-01-12 01:02:03")),
	)

	line, pnts, err := hplot.NewLinePoints(data)
	if err != nil {
		log.Fatalf("could not create plotter: %+v", err)
	}

	line.Color = color.RGBA{B: 255, A: 255}
	pnts.Shape = draw.CircleGlyph{}
	pnts.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, pnts, hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, "testdata/timeseries_monthly.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}
```
