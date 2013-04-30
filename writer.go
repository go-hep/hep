package lhef

import (
	"fmt"
	"io"
)

type Encoder struct {
	w   io.Writer
	Run HEPRUP // User process run common block
}

func NewEncoder(w io.Writer, run HEPRUP) (*Encoder, error) {
	var err error
	enc := &Encoder{
		w:   w,
		Run: run,
	}

	version := float64(1)
	if run.XSecInfo.Neve > 0 {
		version = 2.0
	}
	_, err = fmt.Fprintf(
		w,
		`<LesHouchesEvents version="%0.1f">\n`,
		version,
	)

	if err != nil {
		return nil, err
	}

	return enc, err
}

// EOF
