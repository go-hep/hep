hplt
====

`hplt` is a work-in-progress package providing helper functions to create histograms from n-tuples, using the `sql.DB` as a workhorse.

## Installation

```sh
$> go get github.com/go-hep/hbook/hplt
```

## Documentation

Documentation is available from [godoc](https://godoc.org):
[go-hep/hbook/hplt](https://godoc.org/github.com/go-hep/hbook/hplt)

## Example

### Creating a hbook.H1D from a database

```go
// filling a hbook.H1D with value of 'x'
h, err := hplt.Plot1D(db, "select x from ntuple", nil)

// filling a hbook.H1D with value of 'x' when 'y>10'
h, err := hplt.Plot1D(db, "select x from ntuple where y>10", nil)

// filling an already existing hbook.H1D
h := hbook.NewH1D(100, -10, 10)
h, err := hplt.Plot1D(db, "select x from ntuple", h)

// filling a hbook.H1D with a complex query
h := hbook.NewH1D(100, -10, 10)
err := hplt.Plot(db, "select x from ntuple", func(x float64) error {
	h.Fill(math.Sqrt(x))
	return nil
})

// filling a hbook.H1D with an even more complex query
h := hbook.NewH1D(100, -10, 10)
err := hplt.Plot(db, "select (x,y) from ntuple", func (x, y float64) error {
	h.Fill(math.Sqrt(x*x+y*y))
	return nil
})
```
