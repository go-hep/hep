// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv // import "go-hep.org/x/hep/groot/rsql/rsqldrv"

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/xwb1989/sqlparser"
)

type vctxType map[interface{}]interface{}

func TestExpr(t *testing.T) {
	for _, tc := range []struct {
		expr string
		vctx vctxType
		want interface{}
		err  error
	}{
		{
			expr: `select (x, y) from tbl`,
			vctx: vctxType{"x": int32(1), "y": int32(2)},
			want: []interface{}{int32(1), int32(2)},
		},
		{
			expr: `select (x + "foo") from tbl`,
			vctx: vctxType{"x": "bar"},
			want: "barfoo",
		},
		{
			expr: `select ("bar" + "foo") from tbl`,
			want: "barfoo",
		},
		{
			expr: `select (1 + 1) from tbl`,
			want: idealInt(2),
		},
		{
			expr: `select (10 - 1) from tbl`,
			want: idealInt(9),
		},
		{
			expr: `select (10 - 0xb) from tbl`,
			want: idealInt(-1),
		},
		{
			expr: `select (10 - 0xB) from tbl`,
			want: idealInt(-1),
		},
		{
			expr: `select (10 - 0XB) from tbl`,
			want: idealInt(-1),
		},
		{
			expr: `select (2*2) from tbl`,
			want: idealInt(4),
		},
		{
			expr: `select (30 / 2) from tbl`,
			want: idealInt(15),
		},
		{
			expr: `select (31 / 2.0) from tbl`,
			want: idealFloat(15.5),
		},
		{
			expr: `select (31 / 2e1) from tbl`,
			want: idealFloat(1.55),
		},
		{
			expr: `select (31 / 2E1) from tbl`,
			want: idealFloat(1.55),
		},
		{
			expr: `select (31 / 2) from tbl`,
			want: idealInt(15),
		},
		{
			expr: `select (2.0 + 3.5) from tbl`,
			want: idealFloat(5.5),
		},
		{
			expr: `select (2.0 - 3.5) from tbl`,
			want: idealFloat(-1.5),
		},
		{
			expr: `select (1.0 + 1) from tbl`,
			want: idealFloat(2),
		},
		{
			expr: `select (2.0 * 2) from tbl`,
			want: idealFloat(4),
		},
		{
			expr: `select (2.0 / 4.0) from tbl`,
			want: idealFloat(0.5),
		},
		{
			expr: `select (2.0 < 4.0) from tbl`,
			want: true,
		},
		{
			expr: `select (2.0 <= 4.0) from tbl`,
			want: true,
		},
		{
			expr: `select (2.0 > 4.0) from tbl`,
			want: false,
		},
		{
			expr: `select (2.0 >= 4.0) from tbl`,
			want: false,
		},
		{
			expr: `select (2.0 = 4.0) from tbl`,
			want: false,
		},
		{
			expr: `select (2.0 != 4.0) from tbl`,
			want: true,
		},
		{
			expr: `select (2 < 4) from tbl`,
			want: true,
		},
		{
			expr: `select (2 <= 4) from tbl`,
			want: true,
		},
		{
			expr: `select (2 > 4) from tbl`,
			want: false,
		},
		{
			expr: `select (2 >= 4) from tbl`,
			want: false,
		},
		{
			expr: `select (2 = 4) from tbl`,
			want: false,
		},
		{
			expr: `select (2 != 4) from tbl`,
			want: true,
		},
		{
			expr: `select ("ab" < "bb") from tbl`,
			want: true,
		},
		{
			expr: `select ("ab" <= "bb") from tbl`,
			want: true,
		},
		{
			expr: `select ("ab" > "bb") from tbl`,
			want: false,
		},
		{
			expr: `select ("ab" >= "bb") from tbl`,
			want: false,
		},
		{
			expr: `select ("ab" = "bb") from tbl`,
			want: false,
		},
		{
			expr: `select ("ab" != "bb") from tbl`,
			want: true,
		},
		{
			expr: `select (TRUE || true) from tbl`,
			want: true,
		},
		{
			expr: `select (false || FALSE) from tbl`,
			want: false,
		},
		{
			expr: `select (TRUE != FALSE) from tbl`,
			want: true,
		},
		{
			expr: `select (TRUE = FALSE) from tbl`,
			want: false,
		},
		{
			expr: `select (TRUE = true) from tbl`,
			want: true,
		},
		{
			expr: `select (FALSE = false) from tbl`,
			want: true,
		},
		{
			expr: `select (false && true) from tbl`,
			want: false,
		},
		{
			expr: `select (false || true) from tbl`,
			want: true,
		},
		// idealUint
		{
			expr: `select (x) from tbl`,
			vctx: vctxType{"x": idealUint(5)},
			want: idealUint(5),
		},
		{
			expr: `select (x + y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: idealUint(11),
		},
		{
			expr: `select (x - y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: idealUint(1),
		},
		{
			expr: `select (x / y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: idealUint(1),
		},
		{
			expr: `select (x * y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: idealUint(30),
		},
		{
			expr: `select (x * 0x5) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: idealInt(30),
		},
		{
			expr: `select (x < y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: false,
		},
		{
			expr: `select (x <= y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: false,
		},
		{
			expr: `select (x > y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: true,
		},
		{
			expr: `select (x >= y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: true,
		},
		{
			expr: `select (x = y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: false,
		},
		{
			expr: `select (x != y) from tbl`,
			vctx: vctxType{"x": idealUint(6), "y": idealUint(5)},
			want: true,
		},
		// uint8
		{
			expr: "select (x + y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: uint8(7),
		},
		{
			expr: "select (x - y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: uint8(3),
		},
		{
			expr: "select (x * y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: uint8(10),
		},
		{
			expr: "select (x / y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: uint8(2),
		},
		{
			expr: "select (x < y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // uint8",
			vctx: vctxType{"x": uint8(5), "y": uint8(2)},
			want: true,
		},
		// uint16
		{
			expr: "select (x + y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: uint16(7),
		},
		{
			expr: "select (x - y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: uint16(3),
		},
		{
			expr: "select (x * y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: uint16(10),
		},
		{
			expr: "select (x / y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: uint16(2),
		},
		{
			expr: "select (x < y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // uint16",
			vctx: vctxType{"x": uint16(5), "y": uint16(2)},
			want: true,
		},
		// int32
		{
			expr: "select (x + y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: uint32(7),
		},
		{
			expr: "select (x - y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: uint32(3),
		},
		{
			expr: "select (x * y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: uint32(10),
		},
		{
			expr: "select (x / y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: uint32(2),
		},
		{
			expr: "select (x < y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // uint32",
			vctx: vctxType{"x": uint32(5), "y": uint32(2)},
			want: true,
		},
		// uint64
		{
			expr: "select (x + y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: uint64(7),
		},
		{
			expr: "select (x - y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: uint64(3),
		},
		{
			expr: "select (x * y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: uint64(10),
		},
		{
			expr: "select (x / y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: uint64(2),
		},
		{
			expr: "select (x < y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // uint64",
			vctx: vctxType{"x": uint64(5), "y": uint64(2)},
			want: true,
		},
		// int8
		{
			expr: "select (x + y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: int8(7),
		},
		{
			expr: "select (x - y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: int8(3),
		},
		{
			expr: "select (x * y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: int8(10),
		},
		{
			expr: "select (x / y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: int8(2),
		},
		{
			expr: "select (x < y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // int8",
			vctx: vctxType{"x": int8(5), "y": int8(2)},
			want: true,
		},
		// int16
		{
			expr: "select (x + y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: int16(7),
		},
		{
			expr: "select (x - y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: int16(3),
		},
		{
			expr: "select (x * y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: int16(10),
		},
		{
			expr: "select (x / y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: int16(2),
		},
		{
			expr: "select (x < y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // int16",
			vctx: vctxType{"x": int16(5), "y": int16(2)},
			want: true,
		},
		// int32
		{
			expr: "select (x + y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: int32(7),
		},
		{
			expr: "select (x - y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: int32(3),
		},
		{
			expr: "select (x * y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: int32(10),
		},
		{
			expr: "select (x / y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: int32(2),
		},
		{
			expr: "select (x < y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // int32",
			vctx: vctxType{"x": int32(5), "y": int32(2)},
			want: true,
		},
		// int64
		{
			expr: "select (x + y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: int64(7),
		},
		{
			expr: "select (x - y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: int64(3),
		},
		{
			expr: "select (x * y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: int64(10),
		},
		{
			expr: "select (x / y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: int64(2),
		},
		{
			expr: "select (x < y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // int64",
			vctx: vctxType{"x": int64(5), "y": int64(2)},
			want: true,
		},
		// float32
		{
			expr: "select (x + y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: float32(7),
		},
		{
			expr: "select (x - y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: float32(3),
		},
		{
			expr: "select (x * y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: float32(10),
		},
		{
			expr: "select (x / y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: float32(2.5),
		},
		{
			expr: "select (x < y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // float32",
			vctx: vctxType{"x": float32(5), "y": float32(2)},
			want: true,
		},
		// float64
		{
			expr: "select (x + y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: float64(7),
		},
		{
			expr: "select (x - y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: float64(3),
		},
		{
			expr: "select (x * y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: float64(10),
		},
		{
			expr: "select (x / y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: float64(2.5),
		},
		{
			expr: "select (x < y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: false,
		},
		{
			expr: "select (x <= y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: false,
		},
		{
			expr: "select (x > y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: true,
		},
		{
			expr: "select (x >= y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: true,
		},
		{
			expr: "select (x = y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: false,
		},
		{
			expr: "select (x != y) from tbl // float64",
			vctx: vctxType{"x": float64(5), "y": float64(2)},
			want: true,
		},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			stmt, err := sqlparser.Parse(tc.expr)
			if err != nil {
				t.Fatalf("could not parse %q: %v", tc.expr, err)
			}
			expr, err := newExprFrom(stmt.(*sqlparser.Select).SelectExprs[0].(*sqlparser.AliasedExpr).Expr, nil)
			if err != nil {
				t.Fatalf("could not generate expression: %v", err)
			}
			ectx := newExecCtx(nil, nil)
			v, err := expr.eval(ectx, tc.vctx)
			switch {
			case err == nil && tc.err == nil:
				// ok
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error (got=nil): %v", tc.err)
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %v", err)
			case err.Error() != tc.err.Error():
				t.Fatalf("invalid error.\ngot= %q\nwant=%q", err, tc.err)
			}
			if !reflect.DeepEqual(v, tc.want) {
				t.Fatalf("invalid result.\ngot= %v (%T)\nwant=%v (%T)", v, v, tc.want, tc.want)
			}
		})
	}
}

func TestSelectColumns(t *testing.T) {
	db, err := sql.Open("root", "../../testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type eface = interface{}

	for _, tc := range []struct {
		query string
		cols  []string
		types []interface{}
		args  []interface{}
		vals  [][]eface
	}{
		{
			query: `select one from tree`,
			cols:  []string{"one"},
			types: []interface{}{int32(0)},
			vals: [][]eface{
				[]eface{int32(1)},
				[]eface{int32(2)},
				[]eface{int32(3)},
				[]eface{int32(4)},
			},
		},
		{
			query: `select (one) from tree`,
			cols:  []string{"one"},
			types: []interface{}{int32(0)},
			vals: [][]eface{
				[]eface{int32(1)},
				[]eface{int32(2)},
				[]eface{int32(3)},
				[]eface{int32(4)},
			},
		},
		{
			query: `select (one, two) from tree`,
			cols:  []string{"one", "two"},
			types: []interface{}{int32(0), 0.0},
			vals: [][]eface{
				[]eface{int32(1), 1.1},
				[]eface{int32(2), 2.2},
				[]eface{int32(3), 3.3},
				[]eface{int32(4), 4.4},
			},
		},
		{
			query: `select (one, (two)) from tree`,
			cols:  []string{"one", "two"},
			types: []interface{}{int32(0), 0.0},
			vals: [][]eface{
				[]eface{int32(1), 1.1},
				[]eface{int32(2), 2.2},
				[]eface{int32(3), 3.3},
				[]eface{int32(4), 4.4},
			},
		},
		{
			query: `select (one, ((two))) from tree`,
			cols:  []string{"one", "two"},
			types: []interface{}{int32(0), 0.0},
			vals: [][]eface{
				[]eface{int32(1), 1.1},
				[]eface{int32(2), 2.2},
				[]eface{int32(3), 3.3},
				[]eface{int32(4), 4.4},
			},
		},
		{
			query: `select (((one), ((two)))) from tree`,
			cols:  []string{"one", "two"},
			types: []interface{}{int32(0), 0.0},
			vals: [][]eface{
				[]eface{int32(1), 1.1},
				[]eface{int32(2), 2.2},
				[]eface{int32(3), 3.3},
				[]eface{int32(4), 4.4},
			},
		},
		{
			query: `select three from tree`,
			cols:  []string{"three"},
			types: []interface{}{""},
			vals: [][]eface{
				[]eface{"uno"},
				[]eface{"dos"},
				[]eface{"tres"},
				[]eface{"quatro"},
			},
		},
		{
			query: `select (one, two, three) from tree`,
			cols:  []string{"one", "two", "three"},
			types: []interface{}{int32(0), 0.0, ""},
			vals: [][]eface{
				[]eface{int32(1), 1.1, "uno"},
				[]eface{int32(2), 2.2, "dos"},
				[]eface{int32(3), 3.3, "tres"},
				[]eface{int32(4), 4.4, "quatro"},
			},
		},
		{
			query: `select (?, two, ?) from tree`,
			cols:  []string{"", "two", ""},
			types: []interface{}{"", 0.0, ""},
			args:  []interface{}{"one", "three"},
			vals: [][]eface{
				[]eface{"one", 1.1, "three"},
				[]eface{"one", 2.2, "three"},
				[]eface{"one", 3.3, "three"},
				[]eface{"one", 4.4, "three"},
			},
		},
		{
			query: `select (:v1, two, :v2) from tree`,
			cols:  []string{"", "two", ""},
			types: []interface{}{"", 0.0, ""},
			args:  []interface{}{"one", "three"},
			vals: [][]eface{
				[]eface{"one", 1.1, "three"},
				[]eface{"one", 2.2, "three"},
				[]eface{"one", 3.3, "three"},
				[]eface{"one", 4.4, "three"},
			},
		},
		{
			query: `select (:v2, two, :v1) from tree`,
			cols:  []string{"", "two", ""},
			types: []interface{}{"", 0.0, ""},
			args:  []interface{}{"three", "one"},
			vals: [][]eface{
				[]eface{"one", 1.1, "three"},
				[]eface{"one", 2.2, "three"},
				[]eface{"one", 3.3, "three"},
				[]eface{"one", 4.4, "three"},
			},
		},
		{
			query: `select (:v2, two+:v3, :v1) from tree`,
			cols:  []string{"", "", ""},
			types: []interface{}{"", 0.0, ""},
			args:  []interface{}{"three", "one", 10},
			vals: [][]eface{
				[]eface{"one", 11.1, "three"},
				[]eface{"one", 12.2, "three"},
				[]eface{"one", 13.3, "three"},
				[]eface{"one", 14.4, "three"},
			},
		},
		{
			query: `select (one) from tree where (two > 3)`,
			cols:  []string{"one"},
			types: []interface{}{int32(0)},
			vals: [][]eface{
				[]eface{int32(3)},
				[]eface{int32(4)},
			},
		},
		{
			query: `select (one) from tree where (3 <= two && two < 4)`,
			cols:  []string{"one"},
			types: []interface{}{int32(0)},
			vals: [][]eface{
				[]eface{int32(3)},
			},
		},
		{
			query: `select (one, two) from tree where (three="quatro")`,
			cols:  []string{"one", "two"},
			types: []interface{}{int32(0), 0.0},
			vals: [][]eface{
				[]eface{int32(4), 4.4},
			},
		},
		{
			query: `select (one, two, ?+:v2) from tree where (three="quatro")`,
			cols:  []string{"one", "two", ""},
			types: []interface{}{int32(0), 0.0, uint64(0)},
			args:  []interface{}{idealUint(5), idealUint(10)},
			vals: [][]eface{
				[]eface{int32(4), 4.4, uint64(15)},
			},
		},
	} {
		t.Run(tc.query, func(t *testing.T) {
			rows, err := db.Query(tc.query, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				t.Fatal(err)
			}

			if got, want := cols, tc.cols; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid columns.\ngot= %q\nwant=%q", got, want)
			}

			var got [][]eface
			for rows.Next() {
				vars := make([]interface{}, len(tc.types))
				for i, v := range tc.types {
					vars[i] = reflect.New(reflect.TypeOf(v)).Interface()
				}
				err = rows.Scan(vars...)
				if err != nil {
					t.Fatal(err)
				}
				row := make([]eface, len(vars))
				for i, v := range vars {
					row[i] = reflect.Indirect(reflect.ValueOf(v)).Interface()
				}
				got = append(got, row)
			}

			if got, want := got, tc.vals; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid values.\ngot= %v\nwant=%v\n", got, want)
			}
		})
	}
}
