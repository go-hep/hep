// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/xwb1989/sqlparser"
)

type expression interface {
	sql() sqlparser.Expr

	eval(ectx *execCtx, vctx map[interface{}]interface{}) (v interface{}, err error)
	isStatic() bool
}

type execCtx struct {
	db    *driverConn
	args  []interface{}
	cache map[interface{}]interface{}
	mu    sync.RWMutex
}

func newExecCtx(db *driverConn, args []driver.NamedValue) *execCtx {
	ectx := execCtx{db: db}
	return &ectx
}

type binExpr struct {
	expr sqlparser.Expr

	op operator
	l  expression
	r  expression
}

func newBinExpr(expr sqlparser.Expr, op operator, x, y expression) (v expression, err error) {
	be := &binExpr{
		expr: expr,
		op:   op,
		l:    x,
		r:    y,
	}

	return be, nil
}

func (expr *binExpr) sql() sqlparser.Expr { return expr.expr }
func (expr *binExpr) eval(ectx *execCtx, vctx map[interface{}]interface{}) (r interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case error:
				r, err = nil, x
			default:
				r, err = nil, fmt.Errorf("%v", x)
			}
		}
	}()

	switch expr.op {
	case opAndAnd:
		l, err := expr.l.eval(ectx, vctx)
		if err != nil {
			return nil, err
		}

		switch l := l.(type) {
		case bool:
			if !l {
				return false, nil
			}

			r, err := expr.r.eval(ectx, vctx)
			if err != nil {
				return nil, err
			}

			switch r := r.(type) {
			case bool:
				return r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}

		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opOrOr:
		l, err := expr.l.eval(ectx, vctx)
		if err != nil {
			return nil, err
		}

		switch l := l.(type) {
		case bool:
			if l {
				return true, nil
			}

			r, err := expr.r.eval(ectx, vctx)
			if err != nil {
				return nil, err
			}

			switch r := r.(type) {
			case bool:
				return r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}

		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opLT:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l < r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opLE:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l <= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opGT:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l > r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %#v: %#v", expr.l, l))
		}

	case opGE:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l >= r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opNotEq:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case bool:
			switch r := r.(type) {
			case bool:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l != r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opEq:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case bool:
			switch r := r.(type) {
			case bool:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l == r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opAdd:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return idealUint(uint64(l) + uint64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return idealInt(int64(l) + int64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return idealFloat(float64(l) + float64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case string:
			switch r := r.(type) {
			case string:
				return l + r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opSub:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return idealUint(uint64(l) - uint64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return idealInt(int64(l) - int64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return idealFloat(float64(l) - float64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l - r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opMul:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return idealUint(uint64(l) * uint64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return idealInt(int64(l) * int64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return idealFloat(float64(l) * float64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l * r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	case opDiv:
		l, r := expr.load(ectx, vctx)
		switch l := l.(type) {
		case idealUint:
			switch r := r.(type) {
			case idealUint:
				return idealUint(uint64(l) / uint64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealInt:
			switch r := r.(type) {
			case idealInt:
				return idealInt(int64(l) / int64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case idealFloat:
			switch r := r.(type) {
			case idealFloat:
				return idealFloat(float64(l) / float64(r)), nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v %T", expr.r, r, r))
			}
		case int8:
			switch r := r.(type) {
			case int8:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int16:
			switch r := r.(type) {
			case int16:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int32:
			switch r := r.(type) {
			case int32:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case int64:
			switch r := r.(type) {
			case int64:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint8:
			switch r := r.(type) {
			case uint8:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint16:
			switch r := r.(type) {
			case uint16:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint32:
			switch r := r.(type) {
			case uint32:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case uint64:
			switch r := r.(type) {
			case uint64:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float32:
			switch r := r.(type) {
			case float32:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		case float64:
			switch r := r.(type) {
			case float64:
				return l / r, nil
			default:
				panic(fmt.Errorf("sqldrv: invalid right-operand value/type in %v: %#v", expr.r, r))
			}
		default:
			panic(fmt.Errorf("sqldrv: invalid left-operand value/type in %v: %#v", expr.l, l))
		}

	}

	panic("impossible")
}

func (expr *binExpr) isStatic() bool {
	return expr.l.isStatic() && expr.r.isStatic()
}

func (expr *binExpr) load(ectx *execCtx, vctx map[interface{}]interface{}) (l, r interface{}) {
	l, r = eval2(expr.l, expr.r, ectx, vctx)
	return coerce(l, r)
}

func eval(v expression, ectx *execCtx, vctx map[interface{}]interface{}) (y interface{}) {
	y, err := v.eval(ectx, vctx)
	if err != nil {
		panic(err) // panic ok here
	}
	return
}

func eval2(a, b expression, ectx *execCtx, vctx map[interface{}]interface{}) (x, y interface{}) {
	return eval(a, ectx, vctx), eval(b, ectx, vctx)
}

type operator byte

const (
	opInvalid operator = iota
	opAnd
	opOr
	opXor
	opAdd
	opSub
	opMul
	opDiv
	opAndAnd
	opOrOr
	opLT
	opGT
	opLE
	opGE
	opEq
	opNotEq
	opRShift
	opLShift
)

func operatorFrom(opstr string) operator {
	switch opstr {
	case sqlparser.LessThanStr:
		return opLT
	case sqlparser.LessEqualStr:
		return opLE
	case sqlparser.GreaterThanStr:
		return opGT
	case sqlparser.GreaterEqualStr:
		return opGE
	case sqlparser.NotEqualStr:
		return opNotEq
	case sqlparser.EqualStr:
		return opEq
	case sqlparser.BitAndStr:
		return opAnd
	case sqlparser.BitXorStr:
		return opXor
	case sqlparser.BitOrStr:
		return opOr

	case sqlparser.PlusStr:
		return opAdd
	case sqlparser.MinusStr:
		return opSub
	case sqlparser.MultStr:
		return opMul
	case sqlparser.DivStr:
		return opDiv

	case sqlparser.ShiftLeftStr:
		return opLShift
	case sqlparser.ShiftRightStr:
		return opRShift
	}
	return opInvalid
}

func (op operator) String() string {
	switch op {
	case opAnd:
		return "&"
	case opOr:
		return "|"
	case opXor:
		return "XOr" // FIXME
	case opAdd:
		return "+"
	case opSub:
		return "-"
	case opMul:
		return "*"
	case opDiv:
		return "/"
	case opAndAnd:
		return "&&"
	case opOrOr:
		return "||"
	case opLT:
		return "<"
	case opGT:
		return ">"
	case opLE:
		return "<="
	case opGE:
		return ">="
	case opEq:
		return "=="
	case opNotEq:
		return "!="
	case opRShift:
		return ">>"
	case opLShift:
		return "<<"
	}
	return fmt.Sprintf("%d", byte(op))
}

type identExpr struct {
	expr sqlparser.Expr
	name string
}

func (expr *identExpr) sql() sqlparser.Expr { return expr.expr }
func (expr *identExpr) isStatic() bool      { return false }

func (expr *identExpr) eval(ectx *execCtx, vctx map[interface{}]interface{}) (r interface{}, err error) {
	r, ok := vctx[expr.name]
	if !ok {
		err = fmt.Errorf("unknown field %q", expr.name)
	}
	return r, err
}

type valueExpr struct {
	expr sqlparser.Expr
	v    interface{}
}

func newValueExpr(expr *sqlparser.SQLVal, args []driver.NamedValue) (expression, error) {
	s := string(expr.Val)
	switch expr.Type {
	//	case sqlparser.HexVal: // FIXME(sbinet): difference with HexNum?
	//		v, err := strconv.ParseInt(s, 16, 64)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return &valueExpr{expr: expr, v: idealInt(v)}, nil

	case sqlparser.HexNum:
		v, err := strconv.ParseInt(s[len("0x"):], 16, 64)
		if err != nil {
			return nil, err
		}
		return &valueExpr{expr: expr, v: idealInt(v)}, nil

	case sqlparser.IntVal:
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return &valueExpr{expr: expr, v: idealInt(v)}, nil

	case sqlparser.FloatVal:
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return &valueExpr{expr: expr, v: idealFloat(v)}, nil

	case sqlparser.StrVal:
		return &valueExpr{expr: expr, v: s}, nil

	case sqlparser.ValArg:
		if !strings.HasPrefix(s, ":v") {
			return nil, fmt.Errorf("rsqldrv: invalid ValArg name %q", s)
		}
		i, err := strconv.ParseInt(s[len(":v"):], 10, 64)
		if err != nil {
			return nil, err
		}
		i-- // :v1 --> index-0
		return &valueExpr{
			expr: expr,
			v:    idealValArgFrom(args[i].Value), // FIXME(sbinet): unwrap driver.Value?
		}, nil

	default:
		panic(fmt.Errorf("invalid SQLVal type %#v (%T)", expr, expr))
	}
}

func idealValArgFrom(v interface{}) interface{} {
	switch rv := reflect.ValueOf(v); rv.Kind() {
	case reflect.Int,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return idealInt(rv.Int())

	case reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return idealUint(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return idealFloat(rv.Float())
	case reflect.String:
		return rv.String()
	}
	panic(fmt.Errorf("rsqldrv: invalid ValArg type %#v", v))
}

func (expr *valueExpr) sql() sqlparser.Expr { return expr.expr }
func (expr *valueExpr) isStatic() bool      { return true }
func (expr *valueExpr) eval(ectx *execCtx, vctx map[interface{}]interface{}) (interface{}, error) {
	return expr.v, nil
}

type tupleExpr struct {
	expr  sqlparser.Expr
	exprs []expression
}

func (expr *tupleExpr) sql() sqlparser.Expr { return expr.expr }
func (expr *tupleExpr) isStatic() bool {
	for _, e := range expr.exprs {
		if !e.isStatic() {
			return false
		}
	}
	return true
}
func (expr *tupleExpr) eval(ectx *execCtx, vctx map[interface{}]interface{}) (interface{}, error) {
	var (
		o   = make([]interface{}, len(expr.exprs))
		err error
	)
	for i, e := range expr.exprs {
		o[i], err = e.eval(ectx, vctx)
		if err != nil {
			return nil, err
		}
	}
	return o, nil
}

var (
	_ expression = (*binExpr)(nil)
	_ expression = (*identExpr)(nil)
	_ expression = (*valueExpr)(nil)
	_ expression = (*tupleExpr)(nil)
)
