// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type baseScanner struct {
	tree Tree
	i    int64 // number of entries iterated over
	n    int64 // number of entries to iterate over
	cur  int64 // current entry index
	err  error // last error

	mbr []Branch    // activated branches
	cbr []Branch    // branches activated because holding slice index
	ibr []scanField // indices of activated branches

	closed bool
}

// Close closes the Scanner, preventing further iteration.
// Close is idempotent and does not affect the result of Err.
func (s *baseScanner) Close() error {
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
func (s *baseScanner) Err() error {
	return s.err
}

// Entry returns the entry number of the last read row.
func (s *baseScanner) Entry() int64 {
	return s.cur
}

// SeekEntry points the scanner to the i-th entry, ready to call Next.
func (s *baseScanner) SeekEntry(i int64) error {
	if s.err != nil {
		return s.err
	}
	s.i = i
	s.cur = i - 1
	return s.err
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (s *baseScanner) Next() bool {
	if s.closed {
		return false
	}
	next := s.i < s.n
	switch t := s.tree.(type) {
	case tchain:
		s.cur++
		s.i++
		if s.cur == t.Trees[t.Icur].Entries() {
			t.Icur++
			t.Curtree = t.Trees[t.Icur]
			s.tree = t.Curtree
			s.cur = 0
		}
	case Tree:
		s.cur++
		s.i++
	}
	return next
}

// scanField associates a Branch with a struct's field index
type scanField struct {
	br Branch
	i  int // field index
}

// TreeScanner scans, selects and iterates over Tree entries.
type TreeScanner struct {
	scan baseScanner

	typ reflect.Type  // type bound to this scanner (a struct, a map, a slice of types)
	ptr reflect.Value // pointer to value bound to this scanner
}

// NewTreeScanner creates a new Scanner connecting the pointer to some
// user provided type to the given Tree.
func NewTreeScanner(t Tree, ptr interface{}) (*TreeScanner, error) {
	mbr := make([]Branch, 0, len(t.Branches()))
	ibr := make([]scanField, 0, cap(mbr))
	cbr := make([]Branch, 0)
	rt := reflect.TypeOf(ptr).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, errorf("rootio: NewTreeScanner expects a pointer to a struct (got: %T)", ptr)
	}
	rv := reflect.New(rt).Elem()
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
		leaf := br.Leaves()[0]
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Branch(lcnt.Name())
			if lbr == nil {
				return nil, errorf("rootio: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			cbr = append(cbr, lbr)
		}
		fptr := rv.Field(i).Addr().Interface()
		err := br.setAddress(fptr)
		if err != nil {
			return nil, err
		}
		mbr = append(mbr, br)
		ibr = append(ibr, scanField{br: br, i: i})
	}
	return &TreeScanner{
		scan: baseScanner{
			tree: t,
			i:    0,
			n:    t.Entries(),
			cur:  -1,
			err:  nil,
			ibr:  ibr,
			mbr:  mbr,
			cbr:  cbr,
		},
		typ: rt,
		ptr: rv.Addr(),
	}, nil
}

// ScanVar describes a variable to be read out of a tree during a scan.
type ScanVar struct {
	Name  string      // name of the branch to read
	Leaf  string      // name of the leaf to read
	Value interface{} // pointer to the value to fill
}

// NewTreeScannerVars creates a new Scanner from a list of branches.
// It will return an error if the provided type does not match the
// type stored in the corresponding branch.
func NewTreeScannerVars(t Tree, vars ...ScanVar) (*TreeScanner, error) {
	if len(vars) <= 0 {
		return nil, errorf("rootio: NewTreeScannerVars expects at least one branch name")
	}

	mbr := make([]Branch, len(vars))
	ibr := make([]scanField, cap(mbr))
	cbr := make([]Branch, 0)
	for i, sv := range vars {
		br := t.Branch(sv.Name)
		if br == nil {
			return nil, errorf("rootio: Tree %q has no branch named %q", t.Name(), sv.Name)
		}
		mbr[i] = br
		ibr[i] = scanField{br: br, i: 0}
		leaf := br.Leaves()[0]
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Branch(lcnt.Name())
			if lbr == nil {
				return nil, errorf("rootio: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			cbr = append(cbr, lbr)
		}
	}
	return &TreeScanner{
		scan: baseScanner{
			tree: t,
			i:    0,
			n:    t.Entries(),
			cur:  -1,
			err:  nil,
			ibr:  ibr,
			mbr:  mbr,
			cbr:  cbr,
		},
		typ: nil,
	}, nil
}

// Close closes the TreeScanner, preventing further iteration.
// Close is idempotent and does not affect the result of Err.
func (s *TreeScanner) Close() error {
	return s.scan.Close()
}

// Err returns the error, if any, that was encountered during iteration.
func (s *TreeScanner) Err() error {
	return s.scan.Err()
}

// Entry returns the entry number of the last read row.
func (s *TreeScanner) Entry() int64 {
	return s.scan.Entry()
}

// SeekEntry points the scanner to the i-th entry, ready to call Next.
func (s *TreeScanner) SeekEntry(i int64) error {
	return s.scan.SeekEntry(i)
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (s *TreeScanner) Next() bool {
	err := s.scan.Next()
	/*if s.scan.i == s.scan.n {

	}*/
	//fmt.Printf("\n%v, %v ", s.ptr, s.typ)
	return err
}

// Scan copies data loaded from the underlying Tree into the values pointed at by args.
func (s *TreeScanner) Scan(args ...interface{}) (err error) {
	defer func(err error) {
		if err != nil && s.scan.err == nil {
			s.scan.err = err
		}
	}(err)

	switch len(args) {
	case 0:
		return fmt.Errorf("rootio: TreeScanner.Scan needs at least one argument")

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

	err = s.scanArgs(args...)
	return err
}

func (s *TreeScanner) scanMap(data map[string]interface{}) error {
	panic("not implemented")
}

func (s *TreeScanner) scanArgs(args ...interface{}) error {
	var err error
	// load leaf count data
	for _, br := range s.scan.cbr {
		err = br.loadEntry(s.scan.cur)
		if err != nil {
			// FIXME(sbinet): properly decorate error
			return err
		}
	}

	for i, ptr := range args {
		fv := reflect.ValueOf(ptr).Elem()
		br := s.scan.ibr[i]
		err = br.br.loadEntry(s.scan.cur)
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

func (s *TreeScanner) scanStruct(data interface{}) error {
	var err error
	// load leaf count data
	for _, br := range s.scan.cbr {
		s.scan.err = br.loadEntry(s.scan.cur)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
	}

	rt := reflect.TypeOf(data).Elem()
	rv := reflect.ValueOf(data).Elem()
	if rt != s.typ {
		return errorf("rootio: Scanner.Scan: types do not match (got: %T, want: %T)", rv.Interface(), reflect.New(s.typ).Elem().Interface())
	}

	for _, br := range s.scan.ibr {
		fv := rv.Field(br.i)
		err = br.br.loadEntry(s.scan.cur)
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

// Scanner scans, selects and iterates over Tree entries.
// Scanner is bound to values the user provides, Scanner will
// then read data into these values during the tree scan.
type Scanner struct {
	scan baseScanner
	args []interface{} // a slice of pointers to read data into
}

// NewScannerVars creates a new Scanner from a list of pairs (branch-name, target-address).
// Scanner will read the branches' data during Scan() and load them into these target-addresses.
func NewScannerVars(t Tree, vars ...ScanVar) (*Scanner, error) {
	if len(vars) <= 0 {
		return nil, errorf("rootio: NewScannerVars expects at least one branch name")
	}

	mbr := make([]Branch, len(vars))
	ibr := make([]scanField, cap(mbr))
	cbr := make([]Branch, 0)
	args := make([]interface{}, len(vars))
	for i, sv := range vars {
		br := t.Branch(sv.Name)
		if br == nil {
			return nil, errorf("rootio: Tree %q has no branch named %q", t.Name(), sv.Name)
		}
		mbr[i] = br
		ibr[i] = scanField{br: br, i: 0}

		leaf := br.Leaves()[0]
		if sv.Leaf != "" {
			leaf = br.Leaf(sv.Leaf)
		}
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Branch(lcnt.Name())
			if lbr == nil {
				return nil, errorf("rootio: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			cbr = append(cbr, lbr)
		}
		arg := sv.Value
		if arg == nil {
			return nil, errorf("rootio: ScanVar %d (name=%v) has nil Value", i, sv.Name)
		}
		if rv := reflect.ValueOf(arg); rv.Kind() != reflect.Ptr {
			return nil, errorf("rootio: ScanVar %d (name=%v) has non pointer Value", i, sv.Name)
		}
		err := br.setAddress(arg)
		if err != nil {
			panic(err)
		}
		args[i] = arg
	}

	return &Scanner{
		scan: baseScanner{
			tree: t,
			i:    0,
			n:    t.Entries(),
			cur:  -1,
			err:  nil,
			ibr:  ibr,
			mbr:  mbr,
			cbr:  cbr,
		},
		args: args,
	}, nil
}

// NewScanner creates a new Scanner bound to a (pointer to a) struct value.
// Scanner will read the branches' data during Scan() and load them into the fields of the struct value.
func NewScanner(t Tree, ptr interface{}) (*Scanner, error) {
	mbr := make([]Branch, 0, len(t.Branches()))
	ibr := make([]scanField, 0, cap(mbr))
	cbr := make([]Branch, 0)
	args := make([]interface{}, 0, cap(mbr))
	rt := reflect.TypeOf(ptr).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, errorf("rootio: NewScanner expects a pointer to a struct (got: %T)", ptr)
	}
	rv := reflect.ValueOf(ptr).Elem()
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
		leaf := br.Leaves()[0]
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Branch(lcnt.Name())
			if lbr == nil {
				return nil, errorf("rootio: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			cbr = append(cbr, lbr)
		}
		fptr := rv.Field(i).Addr().Interface()
		mbr = append(mbr, br)
		ibr = append(ibr, scanField{br: br, i: i})
		err := br.setAddress(fptr)
		if err != nil {
			panic(err)
		}
		args = append(args, fptr)
	}

	return &Scanner{
		scan: baseScanner{
			tree: t,
			i:    0,
			n:    t.Entries(),
			cur:  -1,
			err:  nil,
			ibr:  ibr,
			mbr:  mbr,
			cbr:  cbr,
		},
		args: args,
	}, nil
}

// Close closes the Scanner, preventing further iteration.
// Close is idempotent and does not affect the result of Err.
func (s *Scanner) Close() error {
	return s.scan.Close()
}

// Err returns the error, if any, that was encountered during iteration.
func (s *Scanner) Err() error {
	return s.scan.Err()
}

// Entry returns the entry number of the last read row.
func (s *Scanner) Entry() int64 {
	return s.scan.Entry()
}

// SeekEntry points the scanner to the i-th entry, ready to call Next.
func (s *Scanner) SeekEntry(i int64) error {
	return s.scan.SeekEntry(i)
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (s *Scanner) Next() bool {
	return s.scan.Next()
}

// Scan copies data loaded from the underlying Tree into the values the Scanner is bound to.
// The values bound to the Scanner are valid until the next call to Scan.
func (s *Scanner) Scan() error {
	if s.scan.err != nil {
		return s.scan.err
	}

	// load leaf count data
	for _, br := range s.scan.cbr {
		s.scan.err = br.loadEntry(s.scan.cur)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
	}

	for i, ptr := range s.args {
		br := s.scan.ibr[i]
		s.scan.err = br.br.loadEntry(s.scan.cur)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
		s.scan.err = br.br.scan(ptr)
		if s.scan.err != nil {
			return s.scan.err
		}
	}
	return s.scan.err
}
