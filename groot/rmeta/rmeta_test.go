// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rmeta_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rmeta"
)

func TestCxxTemplate(t *testing.T) {
	for _, tt := range []struct {
		n string
		t rmeta.CxxTemplate
	}{
		{
			n: "std::vector<T>",
			t: rmeta.CxxTemplate{
				Name: "std::vector",
				Args: []string{"T"},
			},
		},
		{
			n: "std::vector<std::map<K,V,Cmp>>",
			t: rmeta.CxxTemplate{
				Name: "std::vector",
				Args: []string{"std::map<K,V,Cmp>"},
			},
		},
		{
			n: "std::set<T>",
			t: rmeta.CxxTemplate{
				Name: "std::set",
				Args: []string{"T"},
			},
		},
		{
			n: "std::multiset<T>",
			t: rmeta.CxxTemplate{
				Name: "std::multiset",
				Args: []string{"T"},
			},
		},
		{
			n: "std::unordered_set<T>",
			t: rmeta.CxxTemplate{
				Name: "std::unordered_set",
				Args: []string{"T"},
			},
		},
		{
			n: "std::map<K,V>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "V"},
			},
		},
		{
			n: "std::map<K, V>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "V"},
			},
		},
		{
			n: "std::multimap<K, V>",
			t: rmeta.CxxTemplate{
				Name: "std::multimap",
				Args: []string{"K", "V"},
			},
		},
		{
			n: "std::unordered_map<K, V>",
			t: rmeta.CxxTemplate{
				Name: "std::unordered_map",
				Args: []string{"K", "V"},
			},
		},
		{
			n: "map<K,V>",
			t: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"K", "V"},
			},
		},
		{
			n: "std::map<unsigned long,unsigned int>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"unsigned long", "unsigned int"},
			},
		},
		{
			n: "std::map<K,std::vector<V> >",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			n: "std::map<K,std::vector<V>>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			n: "std::map<K, std::vector<V>>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			n: "map<string,o2::quality_control::core::CheckDefinition>",
			t: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"string", "o2::quality_control::core::CheckDefinition"},
			},
		},
		{
			n: "map<string,o2::quality_control::core::CheckDefinition>",
			t: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"string", "o2::quality_control::core::CheckDefinition"},
			},
		},
		{
			n: "std::map<K, std::map<K2,V2>>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::map<K2,V2>"},
			},
		},
		{
			n: "std::map<std::map<K1,V1>, std::map<K2,V2>>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"std::map<K1,V1>", "std::map<K2,V2>"},
			},
		},
		{
			n: "std::map<std::map<K1,V1>, Foo<T1,T2,T3>>",
			t: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"std::map<K1,V1>", "Foo<T1,T2,T3>"},
			},
		},
		{
			n: "Foo<T1, T2, T3, std::vector<T4, std::allocator<T4>>>",
			t: rmeta.CxxTemplate{
				Name: "Foo",
				Args: []string{"T1", "T2", "T3", "std::vector<T4, std::allocator<T4>>"},
			},
		},
	} {
		t.Run(tt.n, func(t *testing.T) {
			types := rmeta.CxxTemplateFrom(tt.n)
			if got, want := types, tt.t; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid template args.\n got=%q\nwant=%q", got, want)
			}
		})
	}
}
