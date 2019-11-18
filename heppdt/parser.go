// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"golang.org/x/xerrors"
)

// parse fills a Table from the content of r
func parse(r io.Reader, table *Table) error {
	var err error
	s := bufio.NewScanner(r)
	lineno := 0
	for s.Scan() {
		lineno++
		err = s.Err()
		if err != nil {
			break
		}
		bline := s.Bytes()
		bline = bytes.Trim(bline, " \t\r\n")
		if len(bline) <= 0 {
			continue
		}
		if bline[0] == '#' {
			continue
		}
		if len(bline) >= 2 && string(bline[:2]) == "//" {
			continue
		}
		toks := bytes.Split(bline, []byte(" "))
		tokens := make([]string, 0, len(toks))
		for _, tok := range toks {
			//fmt.Printf("--> [%s] => ", string(tok))
			tok = bytes.Trim(tok, " \t\r\n")
			//fmt.Printf("[%s]\n", string(tok))
			if len(tok) > 0 {
				tokens = append(tokens, string(tok))
			}
		}
		if len(tokens) != 6 {
			stoks := ""
			for i, tok := range tokens {
				stoks += fmt.Sprintf("%q", tok)
				if i != len(tokens)-1 {
					stoks += ", "
				}
			}
			//fmt.Printf("** error: line %d (%d): %v\n", lineno, len(tokens), stoks)
			return xerrors.Errorf("heppdt: malformed line:%d: %v", lineno, string(bline))
		}
		var id int64
		id, err = strconv.ParseInt(tokens[0], 10, 64)
		if err != nil {
			return xerrors.Errorf("heppdt: line:%d: %w", lineno, err)
		}
		pid := PID(id)

		name := string(tokens[1])

		var icharge int64
		icharge, err = strconv.ParseInt(tokens[2], 10, 64)
		if err != nil {
			return xerrors.Errorf("heppdt: line:%d: %w", lineno, err)
		}
		var charge float64
		// allow for Q-balls
		if pid.IsQBall() {
			// 10x the charge
			charge = float64(icharge) * 0.01
		} else {
			// 3x the charge
			const onethird = 1. / 3.0
			charge = float64(icharge) * onethird
		}
		var mass float64
		mass, err = strconv.ParseFloat(tokens[3], 64)
		if err != nil {
			return xerrors.Errorf("heppdt: line:%d: %w", lineno, err)
		}

		var totwidth float64
		totwidth, err = strconv.ParseFloat(tokens[4], 64)
		if err != nil {
			return xerrors.Errorf("heppdt: line:%d: %w", lineno, err)
		}

		var lifetime float64
		lifetime, err = strconv.ParseFloat(tokens[5], 64)
		if err != nil {
			return xerrors.Errorf("heppdt: line:%d: %w", lineno, err)
		}

		res := Resonance{
			Mass: Measurement{
				Value: mass,
				Sigma: 0,
			},
		}

		// either width or lifetime is defined. not both.
		switch {
		case totwidth > 0.:
			res.Width = Measurement{
				Value: totwidth,
				Sigma: 0,
			}
		case totwidth == -1.:
			res.Width = Measurement{
				Value: -1.,
				Sigma: 0.,
			}
		case lifetime > 0.:
			res.Width = Measurement{
				Value: calcWidthFromLifetime(lifetime),
				Sigma: 0.,
			}
		default:
			res.Width = Measurement{}
		}

		part := Particle{
			ID:        pid,
			Name:      name,
			PDG:       int(id),
			Mass:      mass,
			Charge:    charge,
			Resonance: res,
		}
		table.pdt[pid] = &part
		table.pid[name] = pid
		//fmt.Printf(">>> %d: [%s]\n", lineno, tokens)
	}
	if err == io.EOF {
		err = nil
	}

	return err
}
