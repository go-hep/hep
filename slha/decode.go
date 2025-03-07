// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slha

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	reBlock = regexp.MustCompile(`BLOCK\s+(\w+)(\s+Q\s*=\s*.+)?`)
	reDecay = regexp.MustCompile(`DECAY\s+(-?\d+)\s+([\d\.E+-]+|NAN).*`)
)

// Decode reads SLHA informations from r and returns them as a *slha.SLHA.
func Decode(r io.Reader) (*SLHA, error) {
	var err error
	type stateType int
	const (
		stBlock stateType = 1
		stDecay stateType = 2
	)
	var state stateType
	var blk *Block
	var part *Particle
	var data SLHA
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		bline := scan.Bytes()
		if len(bline) <= 0 {
			continue
		}
		if bline[0] == '#' {
			// TODO(sbinet) store block/entry comments
			continue
		}
		bup := bytes.ToUpper(bline)
		switch bup[0] {
		case 'B':
			idx := bytes.Index(bup, []byte("#"))
			comment := ""
			if idx > 0 {
				comment = strings.TrimSpace(string(bline[idx+1:]))
				bup = bytes.TrimSpace(bup[:idx])
			}
			// fmt.Printf("Block> %q (comment=%q)\n", string(bup), comment)
			state = stBlock
			all := reBlock.FindStringSubmatch(string(bup))
			if all == nil {
				return nil, fmt.Errorf("slha.decode: invalid block: %q", string(bline))
			}
			groups := all[1:]
			// for i, v := range groups {
			// 	fmt.Printf("  %d/%d: %q\n", i+1, len(groups), v)
			// }
			i := len(data.Blocks)
			data.Blocks = append(data.Blocks, Block{
				Name:    groups[0],
				Comment: comment,
				Q:       math.NaN(),
				Data:    make(DataArray, 0),
			})
			blk = &data.Blocks[i]
			if len(groups) > 1 && groups[1] != "" {
				qstr := groups[1]
				idx := strings.Index(qstr, "=")
				if idx > 0 {
					qstr = strings.TrimSpace(qstr[idx+1:])
				}
				blk.Q, err = strconv.ParseFloat(qstr, 64)
				if err != nil {
					return nil, err
				}
			}
			// fmt.Printf("Block> %v\n", blk)

		case 'D':
			state = stDecay
			idx := bytes.Index(bup, []byte("#"))
			comment := ""
			if idx > 0 {
				comment = strings.TrimSpace(string(bline[idx+1:]))
				bup = bytes.TrimSpace(bup[:idx])
			}
			all := reDecay.FindStringSubmatch(string(bup))
			if all == nil {
				return nil, fmt.Errorf("slha.decode: invalid decay: %q", string(bline))
			}
			groups := all[1:]
			pdgid, err := strconv.Atoi(groups[0])
			if err != nil {
				return nil, err
			}
			width := math.NaN()
			if len(groups) > 1 && groups[1] != "" {
				width, err = strconv.ParseFloat(groups[1], 64)
				if err != nil {
					return nil, err
				}
			}
			i := len(data.Particles)
			data.Particles = append(data.Particles, Particle{
				PdgID:   pdgid,
				Width:   width,
				Mass:    math.NaN(),
				Comment: comment,
				Decays:  make(Decays, 0, 2),
			})
			part = &data.Particles[i]

		case '\t', ' ':
			// data line
			switch state {
			case stBlock:
				err = addBlockEntry(bline, blk)
				if err != nil {
					return nil, err
				}
			case stDecay:
				err = addDecayEntry(bline, part)
				if err != nil {
					return nil, err
				}
			}
		default:

			fmt.Fprintf(os.Stderr, "**WARN** ignoring unknown section [%s]\n", string(bup))
		}
	}
	err = scan.Err()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		err = nil
	}

	// try to populate particles' masses from the MASS block
	if blk := data.Blocks.Get("MASS"); blk != nil {
		for i := range data.Particles {
			part := &data.Particles[i]
			pdgid := part.PdgID
			val, err := blk.Get(pdgid)
			if err == nil {
				part.Mass = val.Float()
			}
		}
	}
	return &data, err
}

func addBlockEntry(line []byte, blk *Block) error {
	var err error
	var val Value
	hidx := bytes.Index(line, []byte("#"))
	if hidx > 0 {
		val.c = strings.TrimSpace(string(line[hidx+1:]))
		line = line[:hidx]
	}
	line = bytes.TrimSpace(line)
	tokens := make([][]byte, 0, 3)
	for _, tok := range bytes.Split(line, []byte(" ")) {
		if len(tok) <= 0 || bytes.Equal(tok, []byte("")) {
			continue
		}
		tokens = append(tokens, tok)
	}

	// switch blk.Name {
	// case "DCINFO":
	// 	tokens
	// }

	ntokens := len(tokens) - 1
	index := make([]int, ntokens)
	for i := range index {
		tok := string(tokens[i])
		index[i], err = strconv.Atoi(tok)
		if err != nil {
			return fmt.Errorf("slha.decode: invalid index %q. err=%v", tok, err)
		}
	}

	sval := string(tokens[len(index)])
	switch blk.Name {
	case "MODSEL":
		v, err := strconv.Atoi(sval)
		if err != nil {
			return err
		}
		val.v = reflect.ValueOf(v)

	case "SPINFO", "DCINFO":
		val.v = reflect.ValueOf(sval)

	default:
		v, err := anyvalue(sval)
		if err != nil {
			return err
		}
		val.v = reflect.ValueOf(v)
	}

	// fmt.Printf("--- %q (comment=%q) len=%d indices=%v val=%#v\n", string(line), val.c, len(tokens), index, val.Interface())

	blk.Data = append(blk.Data, DataItem{
		Index: NewIndex(index...),
		Value: val,
	})
	return err
}

func anyvalue(str string) (any, error) {
	var err error
	var vv any

	vfloat, err := strconv.ParseFloat(str, 64)
	if err == nil {
		vv = reflect.ValueOf(vfloat).Interface()
		return vv, err
	}

	if strings.Count(str, "D") == 1 {
		vfloat, err = strconv.ParseFloat(strings.Replace(str, "D", "E", 1), 64)
		if err == nil {
			vv = reflect.ValueOf(vfloat).Interface()
			return vv, err
		}
	}

	if strings.Count(str, "d") == 1 {
		vfloat, err = strconv.ParseFloat(strings.Replace(str, "d", "E", 1), 64)
		if err == nil {
			vv = reflect.ValueOf(vfloat).Interface()
			return vv, err
		}
	}

	vint, err := strconv.Atoi(str)
	if err == nil {
		vv = reflect.ValueOf(int64(vint)).Interface()
		return vv, err
	}

	vv = str
	return vv, err
}

func addDecayEntry(line []byte, part *Particle) error {
	var err error
	comment := ""
	hidx := bytes.Index(line, []byte("#"))
	if hidx > 0 {
		comment = strings.TrimSpace(string(line[hidx+1:]))
		line = line[:hidx]
	}
	line = bytes.TrimSpace(line)
	tokens := make([][]byte, 0, 3)
	for _, tok := range bytes.Split(line, []byte(" ")) {
		if len(tok) <= 0 || bytes.Equal(tok, []byte("")) {
			continue
		}
		tokens = append(tokens, tok)
	}
	br, err := strconv.ParseFloat(string(tokens[0]), 64)
	if err != nil {
		return fmt.Errorf("slha.decode: invalid decay line %q. err=%v", string(line), err)
	}
	nda, err := strconv.Atoi(string(tokens[1]))
	if err != nil {
		return fmt.Errorf("slha.decode: invalid decay line %q. err=%v", string(line), err)
	}
	ids := make([]int, nda)
	toks := tokens[2:]
	for i := range ids {
		ids[i], err = strconv.Atoi(string(toks[i]))
		if err != nil {
			return fmt.Errorf("slha.decode: invalid decay line %q. err=%v", string(line), err)
		}
	}
	part.Decays = append(part.Decays, Decay{
		Br:      br,
		IDs:     ids,
		Comment: comment,
	})
	// i := len(part.Decays) - 1
	// fmt.Printf("--- %q (comment=%q) len=%d decay=%#v\n", string(line), comment, len(tokens), part.Decays[i])
	return err
}
