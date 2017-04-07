// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
		"USQMIX", "DSQMIX", "SELMIX", "SNUMIX", // FIXME(sbinet) right format ?
		"AU", "AD", "AE",
		"YU", "YD", "YE",
		"TU", "TD", "TE",
		"TUIN", "TDIN", "TEIN",
		"MSQ2IN", "MSU2IN", "MSD2IN", "MSL2IN", "MSE2IN",
		"MSQ2", "MSU2", "MSD2", "MSL2", "MSE2",
		"VCKM", "IMVCKM",
		"UPMNS", "IMUPMNS":
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

	case "SMINPUTS", "VCKMIN", "UPMNSIN":
		//return "(1x,I5,3x,A)"
		return " %5d   %16.8E   # %s\n"

	case "MINPAR", "EXTPAR", "QEXTPAR":
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

// Encode writes the SLHA informations to w.
func Encode(w io.Writer, data *SLHA) error {
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
		_, err = fmt.Fprintf(w, "%s\n", strings.Join(str, " "))
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
			_, err = fmt.Fprintf(w, format, args...)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(w, "#\n")
		if err != nil {
			return err
		}
	}

	for i := range data.Particles {
		part := &data.Particles[i]

		_, err = fmt.Fprintf(w, "%sDECAY %9d   %16.8E   # %s\n", particleHeader, part.PdgID, part.Width, part.Comment)
		if err != nil {
			return err
		}
		if len(part.Decays) <= 0 {
			_, err = fmt.Fprintf(w, "#\n")
			if err != nil {
				return err
			}
			continue
		}

		_, err = fmt.Fprintf(w, decayHeader)
		if err != nil {
			return err
		}

		for j := range part.Decays {
			decay := &part.Decays[j]
			_, err = fmt.Fprintf(w, decayLineFront, decay.Br, len(decay.IDs))
			if err != nil {
				return err
			}
			for _, id := range decay.IDs {
				_, err = fmt.Fprintf(w, decayLineID, id)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintf(w, decayLineBack, decay.Comment)
			if err != nil {
				return err
			}
		}

		_, err = fmt.Fprintf(w, "#\n")
		if err != nil {
			return err
		}

	}
	return err
}
