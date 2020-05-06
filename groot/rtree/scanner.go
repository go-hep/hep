// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/root"
)

type baseScanner struct {
	tree  Tree
	i     int64 // number of entries iterated over
	n     int64 // number of entries to iterate over
	cur   int64 // current entry index
	chain bool  // whether we are scanning through a chain
	off   int64 // entry offset. 0 for TTree.
	tot   int64 // tot-entries
	err   error // last error

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
	if s.chain {
		ch := s.tree.(*tchain)
		if i >= ch.off+ch.tree.Entries() || i < ch.off {
			itree := s.findTree(i)
			if itree == -1 {
				s.err = fmt.Errorf("rtree: could not find Tree containing entry %d", i)
				return s.err
			}
			s.loadTree(itree)
		}
	}
	s.i = i
	s.cur = i - 1
	return s.err
}

// findTree finds the tree number in the chain that contains entry i
func (s *baseScanner) findTree(i int64) int {
	ch := s.tree.(*tchain)
	for j := range ch.trees {
		if i <= ch.tots[j] {
			return j
		}
	}
	return -1
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, false if there is no next result row.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (s *baseScanner) Next() bool {
	if s.closed {
		return false
	}
	next := s.i < s.n
	s.cur++
	s.i++

	if s.chain {
		if s.cur >= s.tot {
			ch := s.tree.(*tchain)
			s.loadTree(ch.cur + 1)
		}
	}

	return next
}

func (s *baseScanner) loadTree(i int) {
	ch := s.tree.(*tchain)
	ch.loadTree(i)
	s.off = ch.off
	s.tot = ch.tot
	if ch.tree == nil {
		// tchain exhausted.
		return
	}
	// reconnect branches
	for i, v := range s.ibr {
		name := v.br.Name()
		br := ch.Branch(name)
		br.setAddress(v.ptr)
		s.ibr[i].br = br
		s.mbr[i] = br
		if v.lcnt >= 0 {
			leaf := br.Leaves()[0]
			lcnt := leaf.LeafCount()
			lbr := ch.Branch(lcnt.Name())
			s.cbr[v.lcnt] = lbr
		}
	}
}

func (s *baseScanner) icur() int64 {
	return s.cur - s.off
}

// scanField associates a Branch with a struct's field index
type scanField struct {
	br   Branch
	i    int         // field index
	ptr  interface{} // field address
	lcnt int         // index of dependant leaf-count (if any)
	dup  bool        // whether the field is already read via leaf-count
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
	cbrset := make(map[string]bool)
	lset := make(map[Leaf]struct{})
	clset := make(map[Leaf]struct{})

	rt := reflect.TypeOf(ptr).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("rtree: NewTreeScanner expects a pointer to a struct (got: %T)", ptr)
	}
	rv := reflect.New(rt).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Name != strings.Title(f.Name) {
			return nil, fmt.Errorf("rtree: field[%d] %q from %T is not exported", i, f.Name, rv.Interface())
		}
		name := f.Tag.Get("groot")
		if name == "" {
			name = f.Name
		}
		if i := strings.Index(name, "["); i > 0 {
			name = name[:i]
		}
		br := t.Branch(name)
		if br == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no branch named %q", t.Name(), name)
		}
		leaf := br.Leaves()[0]
		lset[leaf] = struct{}{}
		lidx := -1
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Leaf(lcnt.Name())
			if lbr == nil {
				return nil, fmt.Errorf("rtree: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			lidx = len(cbr)
			bbr := lbr.Branch()
			if !cbrset[bbr.Name()] {
				cbr = append(cbr, bbr)
				cbrset[bbr.Name()] = true
				clset[lcnt] = struct{}{}
			}
		}
		fptr := rv.Field(i).Addr().Interface()
		err := br.setAddress(fptr)
		if err != nil {
			return nil, err
		}
		mbr = append(mbr, br)
		ibr = append(ibr, scanField{br: br, i: i, ptr: fptr, lcnt: lidx})
	}

	// setup addresses for leaf-count not explicitly requested by user
	for leaf := range clset {
		_, ok := lset[leaf]
		if ok {
			continue
		}
		err := leaf.setAddress(nil)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not set leaf-count address for %q: %w", leaf.Name(), err)
		}
	}

	// remove branches already loaded via leaf-count
	for i, ib := range ibr {
		if _, dup := cbrset[ib.br.Name()]; dup {
			ibr[i].dup = true
		}
	}

	base := baseScanner{
		tree: t,
		i:    0,
		n:    t.Entries(),
		cur:  -1,
		tot:  t.Entries(),
		err:  nil,
		ibr:  ibr,
		mbr:  mbr,
		cbr:  cbr,
	}
	if ch, ok := t.(*tchain); ok {
		base.chain = ok
		base.off = ch.off
		base.tot = ch.tot
	}
	return &TreeScanner{
		scan: base,
		typ:  rt,
		ptr:  rv.Addr(),
	}, nil
}

// ScanVar describes a variable to be read out of a tree.
//
// DEPRECATED: please use ReadVar instead.
type ScanVar = ReadVar

// NewScanVars returns the complete set of ReadVars to read all the data
// contained in the provided Tree.
//
// DEPRECATED: please use NewReadVars instead.
func NewScanVars(t Tree) []ScanVar { return NewReadVars(t) }

// NewTreeScannerVars creates a new Scanner from a list of branches.
// It will return an error if the provided type does not match the
// type stored in the corresponding branch.
func NewTreeScannerVars(t Tree, vars ...ReadVar) (*TreeScanner, error) {
	if len(vars) <= 0 {
		return nil, fmt.Errorf("rtree: NewTreeScannerVars expects at least one branch name")
	}

	mbr := make([]Branch, len(vars))
	ibr := make([]scanField, cap(mbr))
	cbr := make([]Branch, 0)
	cbrset := make(map[string]bool)
	lset := make(map[Leaf]struct{})
	clset := make(map[Leaf]struct{})

	for i, sv := range vars {
		br := t.Branch(sv.Name)
		if br == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no branch named %q", t.Name(), sv.Name)
		}
		mbr[i] = br
		ibr[i] = scanField{br: br, i: 0, lcnt: -1}
		leaf := br.Leaves()[0]
		if sv.Leaf != "" {
			leaf = br.Leaf(sv.Leaf)
		}
		if leaf == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no leaf named %q", t.Name(), sv.Leaf)
		}
		lset[leaf] = struct{}{}
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Leaf(lcnt.Name())
			if lbr == nil {
				return nil, fmt.Errorf("rtree: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			bbr := lbr.Branch()
			if !cbrset[bbr.Name()] {
				cbr = append(cbr, bbr)
				cbrset[bbr.Name()] = true
				clset[lcnt] = struct{}{}
			}
		}
		if sv.Value == nil {
			sv.Value = newValue(leaf)
		}
		arg := sv.Value
		if rv := reflect.ValueOf(arg); rv.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("rtree: ReadVar %d (name=%v) has non pointer Value", i, sv.Name)
		}
		err := br.setAddress(arg)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not set branch address for %q: %w", br.Name(), err)
		}
	}

	// setup addresses for leaf-count not explicitly requested by user
	for leaf := range clset {
		_, ok := lset[leaf]
		if ok {
			continue
		}
		err := leaf.setAddress(nil)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not set leaf-count address for %q: %w", leaf.Name(), err)
		}
	}

	// remove branches already loaded via leaf-count
	for i, ib := range ibr {
		if _, dup := cbrset[ib.br.Name()]; dup {
			ibr[i].dup = true
		}
	}

	base := baseScanner{
		tree: t,
		i:    0,
		n:    t.Entries(),
		cur:  -1,
		err:  nil,
		ibr:  ibr,
		mbr:  mbr,
		cbr:  cbr,
	}

	if ch, ok := t.(*tchain); ok {
		base.chain = ok
		base.off = ch.off
		base.tot = ch.tot
	}

	return &TreeScanner{
		scan: base,
		typ:  nil,
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
	return s.scan.Next()
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
		return fmt.Errorf("rtree: TreeScanner.Scan needs at least one argument")

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

	ientry := s.scan.icur()

	// load leaf count data
	for _, br := range s.scan.cbr {
		err = br.loadEntry(ientry)
		if err != nil {
			// FIXME(sbinet): properly decorate error
			return err
		}
	}

	for i, ptr := range args {
		br := s.scan.ibr[i]
		fv := reflect.ValueOf(ptr).Elem()
		err = br.br.loadEntry(ientry)
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

	ientry := s.scan.icur()

	// load leaf count data
	for _, br := range s.scan.cbr {
		s.scan.err = br.loadEntry(ientry)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
	}

	rt := reflect.TypeOf(data).Elem()
	rv := reflect.ValueOf(data).Elem()
	if rt != s.typ {
		return fmt.Errorf("rtree: Scanner.Scan: types do not match (got: %v, want: %v)", rt, s.typ)
	}
	for _, br := range s.scan.ibr {
		fv := rv.Field(br.i)
		err = br.br.loadEntry(ientry)
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
func NewScannerVars(t Tree, vars ...ReadVar) (*Scanner, error) {
	mbr := make([]Branch, len(vars))
	ibr := make([]scanField, cap(mbr))
	cbr := make([]Branch, 0)
	cbrset := make(map[string]bool)
	lset := make(map[Leaf]struct{})
	clset := make(map[Leaf]struct{})

	args := make([]interface{}, len(vars))
	for i, sv := range vars {
		br := t.Branch(sv.Name)
		if br == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no branch named %q", t.Name(), sv.Name)
		}
		mbr[i] = br
		ibr[i] = scanField{br: br, i: 0, lcnt: -1}

		leaf := br.Leaves()[0]
		if sv.Leaf != "" {
			leaf = br.Leaf(sv.Leaf)
		}
		if leaf == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no leaf named %q", t.Name(), sv.Leaf)
		}
		lset[leaf] = struct{}{}
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Leaf(lcnt.Name())
			if lbr == nil {
				return nil, fmt.Errorf("rtree: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			bbr := lbr.Branch()
			if !cbrset[bbr.Name()] {
				cbr = append(cbr, bbr)
				cbrset[bbr.Name()] = true
				clset[lcnt] = struct{}{}
			}
		}
		arg := sv.Value
		if arg == nil {
			return nil, fmt.Errorf("rtree: ReadVar %d (name=%v) has nil Value", i, sv.Name)
		}
		if rv := reflect.ValueOf(arg); rv.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("rtree: ReadVar %d (name=%v) has non pointer Value", i, sv.Name)
		}
		var err error
		switch br := br.(type) {
		case *tbranchElement:
			err = br.setAddress(arg)
		case *tbranch:
			err = leaf.setAddress(arg)
		default:
			panic(fmt.Errorf("rtree: unknown Branch type %T", br))
		}
		if err != nil {
			panic(err)
		}
		args[i] = arg
		ibr[i].ptr = arg
	}

	// setup addresses for leaf-count not explicitly requested by user
	for leaf := range clset {
		_, ok := lset[leaf]
		if ok {
			continue
		}
		err := leaf.setAddress(nil)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not set leaf-count address for %q: %w", leaf.Name(), err)
		}
	}

	// remove branches already loaded via leaf-count
	for i, ib := range ibr {
		if _, dup := cbrset[ib.br.Name()]; dup {
			ibr[i].dup = true
		}
	}

	base := baseScanner{
		tree: t,
		i:    0,
		n:    t.Entries(),
		cur:  -1,
		err:  nil,
		ibr:  ibr,
		mbr:  mbr,
		cbr:  cbr,
	}

	if ch, ok := t.(*tchain); ok {
		base.chain = ok
		base.off = ch.off
		base.tot = ch.tot
	}

	return &Scanner{
		scan: base,
		args: args,
	}, nil
}

// NewScanner creates a new Scanner bound to a (pointer to a) struct value.
// Scanner will read the branches' data during Scan() and load them into the fields of the struct value.
func NewScanner(t Tree, ptr interface{}) (*Scanner, error) {
	mbr := make([]Branch, 0, len(t.Branches()))
	ibr := make([]scanField, 0, cap(mbr))
	cbr := make([]Branch, 0)
	cbrset := make(map[string]bool)
	lset := make(map[Leaf]struct{})
	clset := make(map[Leaf]struct{})

	args := make([]interface{}, 0, cap(mbr))
	rt := reflect.TypeOf(ptr).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("rtree: NewScanner expects a pointer to a struct (got: %T)", ptr)
	}
	rv := reflect.ValueOf(ptr).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Name != strings.Title(f.Name) {
			return nil, fmt.Errorf("rtree: field[%d] %q from %T is not exported", i, f.Name, rv.Interface())
		}
		name := f.Tag.Get("groot")
		if name == "" {
			name = f.Name
		}
		if i := strings.Index(name, "["); i > 0 {
			name = name[:i]
		}
		br := t.Branch(name)
		if br == nil {
			return nil, fmt.Errorf("rtree: Tree %q has no branch named %q", t.Name(), name)
		}
		leaf := br.Leaves()[0]
		lset[leaf] = struct{}{}
		if lcnt := leaf.LeafCount(); lcnt != nil {
			lbr := t.Leaf(lcnt.Name())
			if lbr == nil {
				return nil, fmt.Errorf("rtree: Tree %q has no (count) branch named %q", t.Name(), lcnt.Name())
			}
			bbr := lbr.Branch()
			if !cbrset[bbr.Name()] {
				cbr = append(cbr, bbr)
				cbrset[bbr.Name()] = true
				clset[lcnt] = struct{}{}
			}
		}
		fptr := rv.Field(i).Addr().Interface()
		err := br.setAddress(fptr)
		if err != nil {
			panic(err)
		}
		args = append(args, fptr)
		mbr = append(mbr, br)
		ibr = append(ibr, scanField{br: br, i: i, ptr: fptr, lcnt: -1})
	}

	// setup addresses for leaf-count not explicitly requested by user
	for leaf := range clset {
		_, ok := lset[leaf]
		if ok {
			continue
		}
		err := leaf.setAddress(nil)
		if err != nil {
			return nil, fmt.Errorf("rtree: could not set leaf-count address for %q: %w", leaf.Name(), err)
		}
	}

	// remove branches already loaded via leaf-count
	for i, ib := range ibr {
		if _, dup := cbrset[ib.br.Name()]; dup {
			ibr[i].dup = true
		}
	}

	base := baseScanner{
		tree: t,
		i:    0,
		n:    t.Entries(),
		cur:  -1,
		err:  nil,
		ibr:  ibr,
		mbr:  mbr,
		cbr:  cbr,
	}

	if ch, ok := t.(*tchain); ok {
		base.chain = ok
		base.off = ch.off
		base.tot = ch.tot
	}

	return &Scanner{
		scan: base,
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

	ientry := s.scan.icur()

	// load leaf count data
	for _, br := range s.scan.cbr {
		s.scan.err = br.loadEntry(ientry)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
	}

	for _, br := range s.scan.ibr {
		if br.dup {
			continue
		}
		s.scan.err = br.br.loadEntry(ientry)
		if s.scan.err != nil {
			// FIXME(sbinet): properly decorate error
			return s.scan.err
		}
	}
	return s.scan.err
}

func newValue(leaf Leaf) interface{} {
	etype := leaf.Type()
	unsigned := leaf.IsUnsigned()

	switch etype.Kind() {
	case reflect.Interface, reflect.Map, reflect.Chan:
		panic(fmt.Errorf("rtree: type %T not supported", reflect.New(etype).Elem().Interface()))
	case reflect.Int8:
		if unsigned {
			etype = reflect.TypeOf(uint8(0))
		}
	case reflect.Int16:
		if unsigned {
			etype = reflect.TypeOf(uint16(0))
		}
	case reflect.Int32:
		if unsigned {
			etype = reflect.TypeOf(uint32(0))
		}
	case reflect.Int64:
		if unsigned {
			etype = reflect.TypeOf(uint64(0))
		}
	case reflect.Float32:
		if _, ok := leaf.(*LeafF16); ok {
			etype = reflect.TypeOf(root.Float16(0))
		}
	case reflect.Float64:
		if _, ok := leaf.(*LeafD32); ok {
			etype = reflect.TypeOf(root.Double32(0))
		}
	}

	switch {
	case leaf.LeafCount() != nil:
		etype = reflect.SliceOf(etype)
	case leaf.Len() > 1:
		switch leaf.Kind() {
		case reflect.String:
			switch dims := leaf.ArrayDim(); dims {
			case 0, 1:
				// interpret as a single string.
			default:
				// FIXME(sbinet): properly handle [N]string (but ROOT doesn't support that.)
				// see: https://root-forum.cern.ch/t/char-t-in-a-branch/5591/2
				// etype = reflect.ArrayOf(leaf.Len(), etype)
				panic(fmt.Errorf("groot/rtree: invalid number of dimensions (%d)", dims))
			}
		default:
			var shape []int
			switch leaf.(type) {
			case *LeafF16, *LeafD32:
				// workaround for https://sft.its.cern.ch/jira/browse/ROOT-10149
				shape = []int{leaf.Len()}
			default:
				shape = leafDims(leaf.Title())
			}
			for i := range shape {
				etype = reflect.ArrayOf(shape[len(shape)-1-i], etype)
			}

		}
	}
	return reflect.New(etype).Interface()
}
