package slha

import (
	"fmt"
	"io"
	"math"
	"strings"
)

func blockFormat(name string, nargs int) string {
	switch name {

	case "MASS":
		//return "(1x,I9,3x,1P,E16.8,0P,3x,’#’,1x,A)"
		return " %9d   %16.8E   # %s\n"

	case "NMIX", "UMIX", "VMIX", "STOPMIX", "SBOTMIX", "STAUMIX",
		"AU", "AD", "AE",
		"YU", "YD", "YE":
		//return "(1x,I2,1x,I2,3x,1P,E16.8,0P,3x,’#’,1x,A)"
		return " %2d %2d   %16.8E   # %s\n"

	case "ALPHA":
		// return "(9x,1P,E16.8,0P,3x,’#’,1x,A)"
		return "         %16.8E   # %s\n"

	case "HMIX",
		"GAUGE",
		"MSOFT":
		//return "(1x,I5,3x,1P,E16.8,0P,3x,’#’,1x,A)"
		return " %5d   %16.8E   # %s\n"

	case "SMINPUTS":
		//return "(1x,I5,3x,A)"
		return " %5d   %16.8E   # %s\n"

	case "MINPAR":
		return " %5d   %16.8E   # %s\n"

	case "SPINFO", "DCINFO":
		return " %5d   %-8s    # %s\n"

	case "MODSEL":
		return " %5d %5d  # %s\n"

	}

	return strings.Repeat(" %v", nargs-1) + " # %s\n"
}

const (
	particleHeader = "#         PDG            Width\n"
	decayHeader    = "#          BR         NDA      ID1       ID2\n"
	decayLineFront = "   %16.8E   %2d  "
	decayLineID    = " %9d"
	decayLineBack  = "   # %s\n"
)

// Encoder writes SLHA objects to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new Encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the SLHA encoding of data to the stream, followed by a newline character.
func (enc *Encoder) Encode(data *SLHA) error {
	var err error
	for i := range data.Blocks {
		blk := &data.Blocks[i]
		str := []string{"BLOCK", blk.Name}
		if !math.IsNaN(blk.Q) {
			str = append(str, fmt.Sprintf("Q=%16.8E", blk.Q))
		}
		if blk.Comment != "" {
			str = append(str, " # "+blk.Comment)
		}
		_, err = fmt.Fprintf(enc.w, "%s\n", strings.Join(str, " "))
		if err != nil {
			return err
		}

		for _, item := range blk.Data {
			v := item.Value
			idx := item.Index.Index() //
			args := make([]interface{}, 0, len(idx)+2)
			if blk.Name != "ALPHA" {
				for _, v := range idx {
					args = append(args, v)
				}
			}
			args = append(args, v.Interface(), v.Comment())
			format := blockFormat(blk.Name, len(args))
			_, err = fmt.Fprintf(enc.w, format, args...)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(enc.w, "#\n")
		if err != nil {
			return err
		}
	}

	for i := range data.Particles {
		part := &data.Particles[i]

		_, err = fmt.Fprintf(enc.w, "%sDECAY %9d   %16.8E   # %s\n", particleHeader, part.PdgID, part.Width, part.Comment)
		if err != nil {
			return err
		}
		if len(part.Decays) <= 0 {
			_, err = fmt.Fprintf(enc.w, "#\n")
			if err != nil {
				return err
			}
			continue
		}

		_, err = fmt.Fprintf(enc.w, decayHeader)
		if err != nil {
			return err
		}

		for j := range part.Decays {
			decay := &part.Decays[j]
			_, err = fmt.Fprintf(enc.w, decayLineFront, decay.Br, len(decay.IDs))
			if err != nil {
				return err
			}
			for _, id := range decay.IDs {
				_, err = fmt.Fprintf(enc.w, decayLineID, id)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(enc.w, decayLineBack, decay.Comment)
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintf(enc.w, "#\n")
		if err != nil {
			return err
		}

	}
	return err
}
