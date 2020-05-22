// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rsqldrv registers a database/sql/driver.Driver implementation for ROOT files.
package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/xwb1989/sqlparser"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

const driverName = "root"

func init() {
	sql.Register(driverName, &rootDriver{})
}

// Open is a ROOT/SQL-driver helper function for sql.Open.
//
// It opens a database connection to the ROOT/SQL driver.
func Open(name string) (*sql.DB, error) {
	return sql.Open(driverName, name)
}

// Create is a ROOT/SQL-driver helper function for sql.Open.
//
// It creates a new ROOT file, connected via the ROOT/SQL driver.
func Create(name string) (*sql.DB, error) {
	panic("not implemented") // FIXME(sbinet)
}

// rootDriver implements the interface required by database/sql/driver.
type rootDriver struct {
	mu   sync.Mutex
	dbs  map[string]*driverConn
	owns map[string]bool // whether the driver owns the ROOT files (and needs to close it)
}

func (drv *rootDriver) open(fname string) (driver.Conn, error) {
	drv.mu.Lock()
	defer drv.mu.Unlock()
	if drv.dbs == nil {
		drv.dbs = make(map[string]*driverConn)
	}
	if drv.owns == nil {
		drv.owns = make(map[string]bool)
	}

	conn := drv.dbs[fname]
	if conn == nil {
		f, err := riofs.Open(fname)
		if err != nil {
			return nil, fmt.Errorf("rsqldriver: could not open file: %w", err)
		}

		conn = &driverConn{
			f: f,
			// cfg:  c,
			drv:  drv,
			stop: make(map[*driverStmt]struct{}),
			refs: 0,
		}

		drv.dbs[fname] = conn
		drv.owns[fname] = true
	}
	conn.refs++

	return conn, nil
}

func (drv *rootDriver) connect(f *riofs.File) driver.Conn {
	drv.mu.Lock()
	defer drv.mu.Unlock()
	if drv.dbs == nil {
		drv.dbs = make(map[string]*driverConn)
	}
	if drv.owns == nil {
		drv.owns = make(map[string]bool)
	}

	conn := drv.dbs[f.Name()]
	if conn == nil {
		conn = &driverConn{
			f: f,
			//cfg:  c,
			drv:  drv,
			stop: make(map[*driverStmt]struct{}),
			refs: 0,
		}
		drv.dbs[f.Name()] = conn
		drv.owns[f.Name()] = false
	}
	conn.refs++

	return conn
}

// Open returns a new connection to the database.
// The name is a string in a driver-specific format.
//
// Open may return a cached connection (one previously
// closed), but doing so is unnecessary; the sql package
// maintains a pool of idle connections for efficient re-use.
//
// The returned connection is only used by one goroutine at a
// time.
func (drv *rootDriver) Open(name string) (driver.Conn, error) {
	return drv.open(name)
}

type driverConn struct {
	f    *riofs.File
	drv  *rootDriver
	stop map[*driverStmt]struct{}
	refs int

	tx driver.Tx
}

// Prepare returns a prepared statement, bound to this connection.
func (conn *driverConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	s := &driverStmt{conn: conn, stmt: stmt}
	conn.stop[s] = struct{}{}
	return s, nil
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
func (conn *driverConn) Close() error {
	conn.drv.mu.Lock()
	defer conn.drv.mu.Unlock()

	if conn.refs > 1 {
		conn.refs--
		return nil
	}

	for s := range conn.stop {
		err := s.Close()
		if err != nil {
			return fmt.Errorf("rsqldrv: could not close statement %v: %w", s, err)
		}
	}

	var err error
	if conn.drv.owns[conn.f.Name()] {
		err = conn.f.Close()
		if err != nil {
			return err
		}
	}

	if conn.refs == 1 {
		delete(conn.drv.dbs, conn.f.Name())
	}
	conn.refs = 0

	return err
}

// Begin starts and returns a new transaction.
func (conn *driverConn) Begin() (driver.Tx, error) {
	panic("conn-begin: not implemented")
}

func (conn *driverConn) Commit() error {
	panic("conn-commit: not implemented")
}

func (conn *driverConn) Rollback() error {
	panic("conn-rollback: not implemented")
}

func (conn *driverConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	return conn.exec(ctx, stmt, args)
}

func (conn *driverConn) exec(ctx context.Context, stmt sqlparser.Statement, args []driver.NamedValue) (driver.Result, error) {
	panic("not implemented")
}

func (conn *driverConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}
	return conn.query(ctx, stmt, args)
}

func (conn *driverConn) query(ctx context.Context, stmt sqlparser.Statement, args []driver.NamedValue) (driver.Rows, error) {
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		rows, err := newDriverRows(ctx, conn, stmt, args)
		return rows, err
	}
	panic("not implemented")
}

type driverResult struct {
	id   int64 // last inserted ID
	rows int64 // rows affected
}

func (res *driverResult) LastInsertId() (int64, error) { return res.id, nil } // -golint
func (res *driverResult) RowsAffected() (int64, error) { return res.rows, nil }

// driverRows is an iterator over an executed query's results.
type driverRows struct {
	conn  *driverConn
	args  []driver.NamedValue
	cols  []string
	types []colDescr    // types of the columns
	deps  []string      // names of the columns to be read
	vars  []interface{} // values of the columns that were read

	cursor *rtree.TreeScanner
	eval   expression
	filter expression
}

type colDescr struct {
	Name     string
	Len      int64 // -1 if no length.
	Nullable bool
	Type     reflect.Type
}

func newDriverRows(ctx context.Context, conn *driverConn, stmt *sqlparser.Select, args []driver.NamedValue) (*driverRows, error) {
	var (
		name = ""
		f    = conn.f
	)

	switch len(stmt.From) {
	case 1:
		switch from := stmt.From[0].(type) {
		case *sqlparser.AliasedTableExpr:
			switch expr := from.Expr.(type) {
			case sqlparser.TableName:
				name = expr.Name.CompliantName()
			default:
				panic(fmt.Errorf("unknown FROM expression type: %#v", expr))
			}

		default:
			panic(fmt.Errorf("unknown table expression: %#v", from))
		}

	default:
		return nil, fmt.Errorf("rsqldrv: invalid number of tables (got=%d, want=1)", len(stmt.From))
	}

	obj, err := riofs.Dir(f).Get(name)
	if err != nil {
		return nil, err
	}

	tree, ok := obj.(rtree.Tree)
	if !ok {
		return nil, fmt.Errorf("rsqldrv: object %q is not a Tree", name)
	}

	rows := &driverRows{conn: conn, args: args}

	rows.cols, err = rows.extractColsFromSelect(tree, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("could not extract columns: %w", err)
	}
	rows.types = make([]colDescr, len(rows.cols))
	for i, name := range rows.cols {
		if name == "" {
			rows.types[i].Type = reflect.TypeOf(new(interface{})).Elem()
			continue
		}
		rows.types[i].Name = name
		branch := tree.Branch(name)
		if branch == nil {
			rows.types[i].Type = reflect.TypeOf(new(interface{})).Elem()
			continue
		}
		rows.types[i] = colDescrFromLeaf(branch.Leaves()[0]) // FIXME(sbinet): multi-leaves' branches
	}

	vars, err := rows.extractDepsFromSelect(tree, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("could not extract read-vars: %w", err)
	}
	rows.vars = varsFrom(vars)
	for _, v := range vars {
		rows.deps = append(rows.deps, v.Name)
	}

	rows.cursor, err = rtree.NewTreeScannerVars(tree, vars...)
	if err != nil {
		return nil, err
	}

	switch expr := stmt.SelectExprs[0].(type) { // FIXME(sbinet): handle multiple select-expressions
	case *sqlparser.AliasedExpr:
		rows.eval, err = newExprFrom(expr.Expr, args)
		if err != nil {
			return nil, fmt.Errorf("could not generate row expression: %w", err)
		}
	case *sqlparser.StarExpr:
		tuple := make(sqlparser.ValTuple, len(rows.cols))
		for i, name := range rows.cols {
			tuple[i] = &sqlparser.ColName{Name: sqlparser.NewColIdent(name)}
		}
		rows.eval, err = newExprFrom(tuple, args)
		if err != nil {
			return nil, fmt.Errorf("could not generate row expression from 'select *': %w", err)
		}
	}

	if stmt.Where != nil {
		switch stmt.Where.Type {
		case sqlparser.WhereStr:
			rows.filter, err = newExprFrom(stmt.Where.Expr, args)
			if err != nil {
				return nil, err
			}
		default:
			panic(fmt.Errorf("unknown 'where' type: %q", stmt.Where.Type))
		}
	}

	return rows, nil
}

func varsFrom(vars []rtree.ReadVar) []interface{} {
	vs := make([]interface{}, len(vars))
	for i, v := range vars {
		vs[i] = v.Value
	}
	return vs
}

// extractDepsFromSelect analyses the query and extracts the branches that need to be read
// for the query to be properly executed.
func (rows *driverRows) extractDepsFromSelect(tree rtree.Tree, stmt *sqlparser.Select, args []driver.NamedValue) ([]rtree.ReadVar, error) {
	var (
		vars []rtree.ReadVar

		set  = make(map[string]struct{})
		cols []string
	)

	markBranch := func(name string) {
		if name != "" {
			if _, dup := set[name]; !dup {
				set[name] = struct{}{}
				cols = append(cols, name)
			}
		}
	}

	collectCols := func(node sqlparser.SQLNode) (bool, error) {
		switch node := node.(type) {
		case *sqlparser.StarExpr:
			other := node.TableName.Name.CompliantName()
			switch other {
			case "", tree.Name():
				for _, b := range tree.Branches() {
					markBranch(b.Name())
				}
			default:
				panic(fmt.Errorf("rsqldrv: star-expression with other table name not supported"))
			}
			return false, nil

		case sqlparser.ColIdent:
			name := node.CompliantName()
			markBranch(name)
			return false, nil

		default:
			return true, nil
		}
	}

	nodes := make([]sqlparser.SQLNode, len(stmt.SelectExprs))
	for i, expr := range stmt.SelectExprs {
		nodes[i] = expr
	}

	if stmt.Where != nil {
		nodes = append(nodes, stmt.Where.Expr)
	}

	err := sqlparser.Walk(collectCols, nodes...)
	if err != nil {
		return nil, err
	}

	for _, name := range cols {
		branch := tree.Branch(name)
		if branch == nil {
			return nil, fmt.Errorf("rsqldrv: could not find branch/leaf %q in tree %q", name, tree.Name())
		}
		leaf := branch.Leaves()[0] // FIXME(sbinet): handle sub-leaves
		etyp := leaf.Type()
		switch etyp.Kind() {
		case reflect.Int8:
			if leaf.IsUnsigned() {
				etyp = reflect.TypeOf(uint8(0))
			}
		case reflect.Int16:
			if leaf.IsUnsigned() {
				etyp = reflect.TypeOf(uint16(0))
			}
		case reflect.Int32:
			if leaf.IsUnsigned() {
				etyp = reflect.TypeOf(uint32(0))
			}
		case reflect.Int64:
			if leaf.IsUnsigned() {
				etyp = reflect.TypeOf(uint64(0))
			}
		}
		switch {
		case leaf.LeafCount() != nil:
			etyp = reflect.SliceOf(etyp)
		case leaf.Len() > 1 && leaf.Kind() != reflect.String:
			etyp = reflect.ArrayOf(leaf.Len(), etyp)
		}
		vars = append(vars, rtree.ReadVar{
			Name:  branch.Name(),
			Leaf:  leaf.Name(),
			Value: reflect.New(etyp).Interface(),
		})
	}

	return vars, nil
}

func (rows *driverRows) extractColsFromSelect(tree rtree.Tree, stmt *sqlparser.Select, args []driver.NamedValue) ([]string, error) {
	var cols []string

	collect := func(node sqlparser.SQLNode) (bool, error) {
		switch node := node.(type) {
		case *sqlparser.ColName:
			return true, nil
		case sqlparser.ColIdent:
			cols = append(cols, node.CompliantName())
			return false, nil
		case *sqlparser.ParenExpr:
			return true, nil
		case sqlparser.ValTuple:
			return true, nil
		case sqlparser.Exprs:
			return true, nil
		case *sqlparser.BinaryExpr:
			// not a simple select query.
			// add a dummy column name and stop recursion
			cols = append(cols, "")
			return false, nil
		case *sqlparser.UnaryExpr:
			return true, nil
		case *sqlparser.SQLVal:
			// not a simple select query.
			// add a dummy column name and stop recursion
			cols = append(cols, "")
			return false, nil
		}
		return false, nil
	}

	switch expr := stmt.SelectExprs[0].(type) { // FIXME(sbinet): handle multiple select-expressions
	case *sqlparser.AliasedExpr:
		err := sqlparser.Walk(collect, expr.Expr)
		return cols, err

	case *sqlparser.StarExpr:
		branches := make([]string, len(tree.Branches()))
		for i, b := range tree.Branches() {
			branches[i] = b.Name()
		}
		return branches, nil

	default:
		panic(fmt.Errorf("rsqldrv: invalid select-expr type %#v", expr))
	}
}

// Columns returns the names of the columns. The number of columns of the
// result is inferred from the length of the slice.  If a particular column
// name isn't known, an empty string should be returned for that entry.
func (r *driverRows) Columns() []string {
	cols := make([]string, len(r.cols))
	copy(cols, r.cols)
	return cols
}

// ColumnTypeScanType returns the value type that can be used to scan types into.
//
// See database/sql/driver.RowsColumnTypeScanType.
func (r *driverRows) ColumnTypeScanType(i int) reflect.Type {
	return r.types[i].Type
}

// ColumnTypeLength returns the column type length for variable length column types such
// as text and binary field types. If the type length is unbounded the value will
// be math.MaxInt64 (any database limits will still apply).
// If the column type is not variable length, such as an int, or if not supported
// by the driver ok is false.
func (r *driverRows) ColumnTypeLength(i int) (length int64, ok bool) {
	col := r.types[i]
	switch col.Len {
	case -1:
		return 0, false
	}
	return col.Len, true
}

// ColumnTypeNullable reports whether the column may be null.
func (r *driverRows) ColumnTypeNullable(i int) (nullable, ok bool) {
	return r.types[i].Nullable, true
}

// Close closes the rows iterator.
func (r *driverRows) Close() error {
	return r.cursor.Close()
}

// Next is called to populate the next row of data into
// the provided slice. The provided slice will be the same
// size as the Columns() are wide.
//
// Next should return io.EOF when there are no more rows.
//
// The dest should not be written to outside of Next. Care
// should be taken when closing Rows not to modify
// a buffer held in dest.
func (r *driverRows) Next(dest []driver.Value) error {
	if !r.cursor.Next() {
		return io.EOF
	}
	err := r.cursor.Scan(r.vars...)
	if err != nil {
		return err
	}

	ectx := newExecCtx(r.conn, r.args)
	vctx := make(map[interface{}]interface{})
	for i, v := range r.vars {
		vctx[r.deps[i]] = reflect.Indirect(reflect.ValueOf(v)).Interface()
	}

	switch r.filter {
	case nil:
		// no filter
	default:
		ok, err := r.filter.eval(ectx, vctx)
		if err != nil {
			//log.Printf("filter.eval: ok=%#v err=%v", ok, err)
			return err
		}
		if !ok.(bool) {
			return r.Next(dest)
		}
	}

	vs, err := r.eval.eval(ectx, vctx)
	// log.Printf("row.eval: v=%#v, err=%v n=%d", vs, err, len(dest))
	if err != nil {
		return fmt.Errorf("could not evaluate row values: %w", err)
	}

	switch vs := vs.(type) {
	case []interface{}:
		for i, v := range vs {
			switch v := v.(type) {
			case string:
				dest[i] = []byte(v)
			default:
				dest[i] = v
			}
		}
	case string:
		dest[0] = []byte(vs)
	default:
		dest[0] = vs
	}

	return nil
}

type driverStmt struct {
	conn *driverConn
	stmt sqlparser.Statement
}

func (stmt *driverStmt) Close() error {
	panic("not implemented")
}

func (stmt *driverStmt) NumInput() int {
	panic("not implemented")
}

func (stmt *driverStmt) Exec(args []driver.Value) (driver.Result, error) {
	panic("not implemented")
}

func (stmt *driverStmt) Query(args []driver.Value) (driver.Rows, error) {
	panic("not implemented")
}

func newExprFrom(expr sqlparser.Expr, args []driver.NamedValue) (expression, error) {
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		op := operatorFrom(expr.Operator)
		if op == opInvalid {
			return nil, fmt.Errorf("rsqldrv: invalid comparison operator %q", expr.Operator)
		}

		l, err := newExprFrom(expr.Left, args)
		if err != nil {
			return nil, err
		}
		r, err := newExprFrom(expr.Right, args)
		if err != nil {
			return nil, err
		}
		return newBinExpr(expr, op, l, r)

	case *sqlparser.ParenExpr:
		return newExprFrom(expr.Expr, args)

	case *sqlparser.AndExpr:
		l, err := newExprFrom(expr.Left, args)
		if err != nil {
			return nil, err
		}
		r, err := newExprFrom(expr.Right, args)
		if err != nil {
			return nil, err
		}
		return newBinExpr(expr, opAndAnd, l, r)

	case *sqlparser.OrExpr:
		l, err := newExprFrom(expr.Left, args)
		if err != nil {
			return nil, err
		}
		r, err := newExprFrom(expr.Right, args)
		if err != nil {
			return nil, err
		}
		return newBinExpr(expr, opOrOr, l, r)

	case *sqlparser.ColName:
		return &identExpr{
			expr: expr,
			name: expr.Name.CompliantName(),
		}, nil

	case *sqlparser.SQLVal:
		return newValueExpr(expr, args)

	case sqlparser.BoolVal:
		return &valueExpr{expr: expr, v: bool(expr)}, nil

	case *sqlparser.BinaryExpr:
		l, err := newExprFrom(expr.Left, args)
		if err != nil {
			return nil, err
		}
		r, err := newExprFrom(expr.Right, args)
		if err != nil {
			return nil, err
		}
		op := operatorFrom(expr.Operator)
		if op == opInvalid {
			return nil, fmt.Errorf("rsqldrv: invalid binary-expression operator %q", expr.Operator)
		}
		return newBinExpr(expr, op, l, r)

	case sqlparser.ValTuple:
		vs := make([]expression, len(expr))
		for i, e := range expr {
			v, err := newExprFrom(e, args)
			if err != nil {
				return nil, err
			}
			vs[i] = v
		}
		return &tupleExpr{expr: expr, exprs: vs}, nil
	}
	return nil, fmt.Errorf("rsqldrv: invalid filter expression %#v %T", expr, expr)
}

var (
	_ driver.Driver         = (*rootDriver)(nil)
	_ driver.Conn           = (*driverConn)(nil)
	_ driver.ExecerContext  = (*driverConn)(nil)
	_ driver.QueryerContext = (*driverConn)(nil)
	_ driver.Tx             = (*driverConn)(nil)

	_ driver.Result = (*driverResult)(nil)
	_ driver.Rows   = (*driverRows)(nil)
)

var (
	_ driver.RowsColumnTypeLength   = (*driverRows)(nil)
	_ driver.RowsColumnTypeNullable = (*driverRows)(nil)
	_ driver.RowsColumnTypeScanType = (*driverRows)(nil)
)
