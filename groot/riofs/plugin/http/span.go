// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"slices"
	"sort"
)

type span struct {
	off int64
	len int64
}

func split(sp span, sps []span) []span {
	if len(sps) == 0 {
		return []span{sp}
	}

	var o []span

	for _, v := range sps {
		b1 := v.off
		e1 := v.off + v.len
		b2 := sp.off
		e2 := sp.off + sp.len
		switch {
		case e1 <= b2:
			//  [ s1=v ]
			//         [  s2=sp ]
			continue
		case e2 <= b1:
			//         [ s1 ]
			// [ s2=sp ]
			o = append(o, sp)
			sp.len = 0
		case b2 < b1 && e1 <= e2:
			//   [ s1=v ]
			//  [  s2=sp ]
			len := b1 - b2
			o = append(o, span{
				off: b2,
				len: len,
			})
			sp.off = e1
			sp.len -= v.len + len
		case b2 < b1 && b1 < e2 && e2 < e1:
			//   [ s1=v   ]
			//  [  s2=sp ]
			len := b1 - b2
			o = append(o, span{
				off: b2,
				len: len,
			})
			sp.len = 0
		case b1 <= b2 && e2 <= e1:
			//  [  s1=v   ]
			//    [s2=sp]
			sp.len = 0
		case b1 <= b2 && e1 < e2:
			//  [  s1=v  ]
			//    [s2=sp  ]
			sp.off = e1
			sp.len = e2 - e1
		}
		if sp.len == 0 {
			break
		}
	}
	if sp.len != 0 {
		o = append(o, sp)
	}

	return o
}

type spans []span

func (p *spans) consolidate() {
	for i := len(*p) - 1; i >= 1; i-- {
		ii := &(*p)[i]
		jj := &(*p)[i-1]
		jend := jj.off + jj.len
		iend := ii.off + ii.len
		if jend < ii.off {
			continue
		}
		if iend >= jend {
			jj.len += iend - jend
		}
		p.remove(i)
	}
}

func (p *spans) remove(i int) {
	list := *p
	*p = slices.Delete(list, i, i+1)
}

func (p *spans) add(sp span) {
	*p = append(*p, sp)
	sort.Slice(*p, func(i, j int) bool {
		pi := (*p)[i]
		pj := (*p)[j]
		if pi.off < pj.off {
			return true
		}
		if pi.off == pj.off {
			return pi.off+pi.len < pj.off+pj.len
		}
		return false
	})
	p.consolidate()
}
