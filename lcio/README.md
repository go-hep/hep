# lcio

[![GoDoc](https://godoc.org/go-hep.org/x/hep/lcio?status.svg)](https://godoc.org/go-hep.org/x/hep/lcio)

`lcio` is a pure `Go` implementation of [LCIO](https://github.com/iLCSoft/LCIO).

## Installation

```sh
$ go get go-hep.org/x/hep/lcio
```

## Documentation

The documentation is browsable at godoc.org:

- https://godoc.org/go-hep.org/x/hep/lcio

## Example

### Reading a LCIO event file

[embedmd]:# (reader_test.go go /func ExampleReader/ /\n}/)
```go
func ExampleReader() {
	r, err := lcio.Open("testdata/event_golden.slcio")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for r.Next() {
		evt := r.Event()
		fmt.Printf("event number = %d (weight=%+e)\n", evt.EventNumber, evt.Weight())
		fmt.Printf("run   number = %d\n", evt.RunNumber)
		fmt.Printf("detector     = %q\n", evt.Detector)
		fmt.Printf("collections  = %v\n", evt.Names())
		calohits := evt.Get("CaloHits").(*lcio.CalorimeterHitContainer)
		fmt.Printf("calohits: %d\n", len(calohits.Hits))
		for i, hit := range calohits.Hits {
			fmt.Printf(" calohit[%d]: cell-id0=%d cell-id1=%d ene=%+e ene-err=%+e\n",
				i, hit.CellID0, hit.CellID1, hit.Energy, hit.EnergyErr,
			)
		}
	}

	err = r.Err()
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// event number = 52 (weight=+4.200000e+01)
	// run   number = 42
	// detector     = "my detector"
	// collections  = [McParticles SimCaloHits CaloHits]
	// calohits: 1
	//  calohit[0]: cell-id0=1024 cell-id1=2048 ene=+1.000000e+03 ene-err=+1.000000e-01
}
```

### Writing a LCIO event file

[embedmd]:# (writer_test.go go /func ExampleWriter/ /\n}/)
```go
func ExampleWriter() {
	w, err := lcio.Create("out.slcio")
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	run := lcio.RunHeader{
		RunNumber:    42,
		Descr:        "a simple run header",
		Detector:     "my detector",
		SubDetectors: []string{"det-1", "det-2"},
		Params: lcio.Params{
			Floats: map[string][]float32{
				"floats-1": {1, 2, 3},
				"floats-2": {4, 5, 6},
			},
		},
	}

	err = w.WriteRunHeader(&run)
	if err != nil {
		log.Fatal(err)
	}

	const NEVENTS = 1
	for ievt := 0; ievt < NEVENTS; ievt++ {
		evt := lcio.Event{
			RunNumber:   run.RunNumber,
			Detector:    run.Detector,
			EventNumber: 52 + int32(ievt),
			TimeStamp:   1234567890 + int64(ievt),
			Params: lcio.Params{
				Floats: map[string][]float32{
					"_weight": {42},
				},
				Strings: map[string][]string{
					"Descr": {"a description"},
				},
			},
		}

		calhits := lcio.CalorimeterHitContainer{
			Flags: lcio.BitsRChLong | lcio.BitsRChID1 | lcio.BitsRChTime | lcio.BitsRChNoPtr | lcio.BitsRChEnergyError,
			Params: lcio.Params{
				Floats:  map[string][]float32{"f32": {1, 2, 3}},
				Ints:    map[string][]int32{"i32": {1, 2, 3}},
				Strings: map[string][]string{"str": {"1", "2", "3"}},
			},
			Hits: []lcio.CalorimeterHit{
				{
					CellID0:   1024,
					CellID1:   2048,
					Energy:    1000,
					EnergyErr: 0.1,
					Time:      1234,
					Pos:       [3]float32{11, 22, 33},
					Type:      42,
				},
			},
		}

		evt.Add("CaloHits", &calhits)

		fmt.Printf("evt has key %q: %v\n", "CaloHits", evt.Has("CaloHits"))

		err = w.WriteEvent(&evt)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// evt has key "CaloHits": true
}
```

### Reading and plotting McParticles' energy

```sh
$> lcio-ex-read-event ./DST01-06_ppr004_bbcsdu.slcio
lcio-ex-read-event: read 50 events from file "./DST01-06_ppr004_bbcsdu.slcio"

$> open out.png
```

![hist-example](https://github.com/go-hep/hep/raw/master/lcio/example/lcio-ex-read-event/out.png)

[embedmd]:# (example/lcio-ex-read-event/main.go go /func main/ /\n}/)
```go
func main() {
	log.SetPrefix("lcio-ex-read-event: ")
	log.SetFlags(0)

	var (
		fname  = ""
		h      = hbook.NewH1D(100, 0., 100.)
		nevts  = 0
		mcname = flag.String("mc", "MCParticlesSkimmed", "name of the MCParticle collection to read")
	)

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := lcio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for f.Next() {
		evt := f.Event()
		mcs := evt.Get(*mcname).(*lcio.McParticleContainer)
		for _, mc := range mcs.Particles {
			h.Fill(mc.Energy(), 1)
		}
		nevts++
	}

	err = f.Err()
	if err == io.EOF {
		err = nil
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("read %d events from file %q", nevts, fname)

	p := hplot.New()
	p.Title.Text = "LCIO -- McParticles"
	p.X.Label.Text = "E (GeV)"

	hh := hplot.NewH1D(h)
	hh.Color = color.RGBA{R: 255, A: 255}
	hh.Infos.Style = hplot.HInfoSummary

	p.Add(hh)
	p.Add(hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, -1, "out.png")
	if err != nil {
		log.Fatal(err)
	}
}
```

