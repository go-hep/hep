// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

// Scanner scans, selects and iterates over Tree entries.
type Scanner struct {
	tree Tree
	i    int64 // number of entries iterated over
	n    int64 // number of entries to iterate over
	cur  int64 // current entry index
	err  error // last error

	mbr []Branch    // activated branches
	ibr []scanField // indices of activated branches

	typ reflect.Type // type bound to this scanner (a struct, a map, a slice of types)

	closed bool
}

// scanField associates a Branch with a struct's field index
type scanField struct {
	br Branch
	i  int // field index
}

// NewScanner creates a new Scanner connecting the pointer to some
// user provided type to the given Tree.
func NewScanner(t Tree, ptr interface{}) (*Scanner, error) {
	mbr := make([]Branch, 0, len(t.Branches()))
	ibr := make([]scanField, 0, cap(mbr))
	rt := reflect.TypeOf(ptr).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, errorf("rootio: NewScanner expects a pointer to a struct (got: %T)", ptr)
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		name := f.Tag.Get("rootio")
		if name == "" {
			name = f.Name
		}
		br := t.Branch(name)
		if br == nil {
			return nil, errorf("rootio: Tree %q has no branch named %q", t.Name(), name)
		}
		mbr = append(mbr, br)
		ibr = append(ibr, scanField{br: br, i: i})
	}
	return &Scanner{
		tree: t,
		i:    0,
		n:    t.Entries(),
		cur:  -1,
		err:  nil,
		ibr:  ibr,
		mbr:  mbr,
		typ:  reflect.TypeOf(ptr).Elem(),
	}, nil
}

// Close closes the Scanner, preventing further iteration.
// Close is idempotent and does not affect the result of Err.
func (s *Scanner) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	s.tree = nil
	s.mbr = nil
	s.ibr = nil
	return nil
}

// Err returns the error, if any, that was encountered during iteration.
func (s *Scanner) Err() error {
	return s.err
}

// Entry returns the entry number of the last read row.
func (s *Scanner) Entry() int64 {
	return s.cur
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (s *Scanner) Next() bool {
	if s.closed {
		return false
	}
	next := s.i < s.n
	s.cur++
	s.i++
	if !next {
		s.err = s.Close()
	}
	return next
}

// Scan copies data loaded from the underlying Tree into the values pointed at by args.
func (s *Scanner) Scan(args ...interface{}) (err error) {
	defer func(err error) {
		if err != nil && s.err == nil {
			s.err = err
		}
	}(err)

	switch len(args) {
	case 0:
		return fmt.Errorf("rootio: Scanner.Scan needs at least one argument")

	case 1:
		// maybe special case: map? struct?
		rt := reflect.TypeOf(args[0]).Elem()
		switch rt.Kind() {
		case reflect.Map:
			err = s.scanMap(*args[0].(*map[string]interface{}))
			return err
		case reflect.Struct:
			err = s.scanStruct(args[0])
			return err
		}
	}

	err = s.scan(args...)
	return err
}

func (s *Scanner) scanMap(data map[string]interface{}) error {
	panic("not implemented")
}

func (s *Scanner) scan(args ...interface{}) error {
	panic("not implemented")
}

func (s *Scanner) scanStruct(data interface{}) error {
	var err error

	rt := reflect.TypeOf(data).Elem()
	rv := reflect.ValueOf(data).Elem()
	if rt != s.typ {
		return errorf("rootio: Scanner.Scan: types do not match (got: %T, want: %T)", rv.Interface(), reflect.New(s.typ).Elem().Interface())
	}
	for _, br := range s.ibr {
		fv := rv.Field(br.i)
		err = br.br.loadEntry(s.cur)
		if err != nil {
			// FIXME(sbinet): properly decorate error
			return err
		}
		err = br.br.scan(fv.Addr().Interface())
		if err != nil {
			return err
		}
	}
	return err
}
