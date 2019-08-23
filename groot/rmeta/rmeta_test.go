// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rmeta_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rmeta"
)

func TestCxxTemplateArgsOf(t *testing.T) {
	for _, tt := range []struct {
		n string
		t []string
	}{
		{
			n: "std::vector<T>",
			t: []string{"T"},
		},
		{
			n: "std::set<T>",
			t: []string{"T"},
		},
		{
			n: "std::multiset<T>",
			t: []string{"T"},
		},
		{
			n: "std::unordered_set<T>",
			t: []string{"T"},
		},
		{
			n: "std::map<K,V>",
			t: []string{"K", "V"},
		},
		{
			n: "std::map<K, V>",
			t: []string{"K", "V"},
		},
		{
			n: "std::multimap<K, V>",
			t: []string{"K", "V"},
		},
		{
			n: "std::unordered_map<K, V>",
			t: []string{"K", "V"},
		},
		{
			n: "map<K,V>",
			t: []string{"K", "V"},
		},
		{
			n: "std::map<unsigned long,unsigned int>",
			t: []string{"unsigned long", "unsigned int"},
		},
		{
			n: "std::map<K,std::vector<V> >",
			t: []string{"K", "std::vector<V>"},
		},
		{
			n: "std::map<K,std::vector<V>>",
			t: []string{"K", "std::vector<V>"},
		},
		{
			n: "std::map<K, std::vector<V>>",
			t: []string{"K", "std::vector<V>"},
		},
		{
			n: "map<string,o2::quality_control::core::CheckDefinition>",
			t: []string{"string", "o2::quality_control::core::CheckDefinition"},
		},
		{
			n: "map<string,o2::quality_control::core::CheckDefinition>",
			t: []string{"string", "o2::quality_control::core::CheckDefinition"},
		},
		{
			n: "std::map<K, std::map<K2,V2>>",
			t: []string{"K", "std::map<K2,V2>"},
		},
		{
			n: "std::map<std::map<K1,V1>, std::map<K2,V2>>",
			t: []string{"std::map<K1,V1>", "std::map<K2,V2>"},
		},
		{
			n: "std::map<std::map<K1,V1>, Foo<T1,T2,T3>>",
			t: []string{"std::map<K1,V1>", "Foo<T1,T2,T3>"},
		},
		{
			n: "Foo<T1, T2, T3, std::vector<T4, std::allocator<T4>>>",
			t: []string{"T1", "T2", "T3", "std::vector<T4, std::allocator<T4>>"},
		},
	} {
		t.Run(tt.n, func(t *testing.T) {
			types := rmeta.CxxTemplateArgsOf(tt.n)
			if got, want := types, tt.t; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid template args.\n got=%q\nwant=%q", got, want)
			}
		})
	}
}
