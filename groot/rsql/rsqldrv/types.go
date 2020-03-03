// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"fmt"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/rtree"
)

type (
	idealFloat float64
	idealInt   int64
	idealUint  uint64
)

func coerce(a, b interface{}) (x, y interface{}) {
	if reflect.TypeOf(a) == reflect.TypeOf(b) {
		return a, b
	}

	switch a.(type) {
	case idealFloat, idealInt, idealUint:
		switch b.(type) {
		case idealFloat, idealInt, idealUint:
			x, y = coerce1(a, b), b
			if reflect.TypeOf(x) == reflect.TypeOf(y) {
				return
			}

			return a, coerce1(b, a)
		default:
			return coerce1(a, b), b
		}
	default:
		switch b.(type) {
		case idealFloat, idealInt, idealUint:
			return a, coerce1(b, a)
		default:
			return a, b
		}
	}
}

func coerce1(inVal, otherVal interface{}) (coercedInVal interface{}) {
	coercedInVal = inVal
	if otherVal == nil {
		return
	}

	switch x := inVal.(type) {
	case nil:
		return
	case idealFloat:
		switch otherVal.(type) {
		case idealFloat:
			return idealFloat(float64(x))
		//case idealInt:
		//case idealRune:
		//case idealUint:
		//case bool:
		case float32:
			return float32(float64(x))
		case float64:
			return float64(x)
			//case int8:
			//case int16:
			//case int32:
			//case int64:
			//case string:
			//case uint8:
			//case uint16:
			//case uint32:
			//case uint64:
		}
	case idealInt:
		switch otherVal.(type) {
		case idealFloat:
			return idealFloat(int64(x))
		case idealInt:
			return idealInt(int64(x))
		//case idealRune:
		case idealUint:
			if x >= 0 {
				return idealUint(int64(x))
			}
		//case bool:
		case float32:
			return float32(int64(x))
		case float64:
			return float64(int64(x))
		case int8:
			if x >= math.MinInt8 && x <= math.MaxInt8 {
				return int8(int64(x))
			}
		case int16:
			if x >= math.MinInt16 && x <= math.MaxInt16 {
				return int16(int64(x))
			}
		case int32:
			if x >= math.MinInt32 && x <= math.MaxInt32 {
				return int32(int64(x))
			}
		case int64:
			return int64(x)
		//case string:
		case uint8:
			if x >= 0 && x <= math.MaxUint8 {
				return uint8(int64(x))
			}
		case uint16:
			if x >= 0 && x <= math.MaxUint16 {
				return uint16(int64(x))
			}
		case uint32:
			if x >= 0 && x <= math.MaxUint32 {
				return uint32(int64(x))
			}
		case uint64:
			if x >= 0 {
				return uint64(int64(x))
			}
		}
	case idealUint:
		switch otherVal.(type) {
		case idealFloat:
			return idealFloat(uint64(x))
		case idealInt:
			if x <= math.MaxInt64 {
				return idealInt(int64(x))
			}
		//case idealRune:
		case idealUint:
			return idealUint(uint64(x))
		//case bool:
		case float32:
			return float32(uint64(x))
		case float64:
			return float64(uint64(x))
		case int8:
			if x <= math.MaxInt8 {
				return int8(int64(x))
			}
		case int16:
			if x <= math.MaxInt16 {
				return int16(int64(x))
			}
		case int32:
			if x <= math.MaxInt32 {
				return int32(int64(x))
			}
		case int64:
			if x <= math.MaxInt64 {
				return int64(x)
			}
		//case string:
		case uint8:
			if x >= 0 && x <= math.MaxUint8 {
				return uint8(int64(x))
			}
		case uint16:
			if x >= 0 && x <= math.MaxUint16 {
				return uint16(int64(x))
			}
		case uint32:
			if x >= 0 && x <= math.MaxUint32 {
				return uint32(int64(x))
			}
		case uint64:
			return uint64(x)
		}
	}
	return
}

func colDescrFromLeaf(leaf rtree.Leaf) colDescr {
	name := leaf.Name()
	etyp := leaf.Type()
	kind := leaf.Kind()
	hasCount := leaf.LeafCount() != nil
	unsigned := leaf.IsUnsigned()

	size := 1
	if !hasCount {
		size = leaf.Len()
	}

	return colDescrFrom(name, etyp, kind, hasCount, size, unsigned)
}

func colDescrFrom(name string, etyp reflect.Type, kind reflect.Kind, hasCount bool, size int, unsigned bool) colDescr {
	col := colDescr{
		Name: name,
		Len:  -1,
	}

	switch {
	case hasCount:
		// slice
		col.Nullable = true
		col.Len = math.MaxInt64
	case size > 1 && kind != reflect.String:
		// array
		col.Len = int64(size)
	}

	switch etyp.Kind() {
	case reflect.Interface, reflect.Map, reflect.Chan, reflect.Slice, reflect.Array:
		panic(fmt.Errorf("rsqldrv: type %T not supported", reflect.New(etyp).Elem().Interface()))
	case reflect.Int8:
		if unsigned {
			etyp = reflect.TypeOf(uint8(0))
		}
	case reflect.Int16:
		if unsigned {
			etyp = reflect.TypeOf(uint16(0))
		}
	case reflect.Int32:
		if unsigned {
			etyp = reflect.TypeOf(uint32(0))
		}
	case reflect.Int64:
		if unsigned {
			etyp = reflect.TypeOf(uint64(0))
		}
	}

	col.Type = etyp
	return col
}
