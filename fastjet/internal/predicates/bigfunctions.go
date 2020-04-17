// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package predicates

import (
	"math/big"
)

var (
	one = new(big.Rat).SetInt64(1)
)

func setBig(f float64) *big.Rat {
	return new(big.Rat).SetFloat64(f)
}

func bigAdd(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Add(a, b)
}

func bigSub(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Sub(a, b)
}

func bigMul(a, b *big.Rat) *big.Rat {
	return new(big.Rat).Mul(a, b)
}
