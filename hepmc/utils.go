package hepmc

import (
	"strconv"
	"strings"
)

type tokens struct {
	toks []string
	pos  int
}

func newtokens(toks []string) tokens {
	return tokens{
		toks: toks,
		pos:  0,
	}
}

func (t *tokens) next() string {
	str := t.toks[t.pos]
	t.pos++
	return str
}

func (t *tokens) at(i int) string {
	return t.toks[i]
}

func (t *tokens) float64() (float64, error) {
	s := t.next()
	return strconv.ParseFloat(s, 64)
}

func (t *tokens) float32() (float32, error) {
	s := t.next()
	v, err := strconv.ParseFloat(s, 64)
	return float32(v), err
}

func (t *tokens) int() (int, error) {
	s := t.next()
	return strconv.Atoi(s)
}

func (t *tokens) int64() (int64, error) {
	s := t.next()
	return strconv.ParseInt(s, 10, 0)
}

func (t *tokens) String() string {
	return strings.Join(t.toks, " ")
}
