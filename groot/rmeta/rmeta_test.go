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
	for _, tc := range []struct {
		name   string
		want   rmeta.CxxTemplate
		panics string
	}{
		// std::pair
		{
			name: "pair<int,int>",
			want: rmeta.CxxTemplate{
				Name: "pair",
				Args: []string{"int", "int"},
			},
		},
		{
			name: "pair<int, int>",
			want: rmeta.CxxTemplate{
				Name: "pair",
				Args: []string{"int", "int"},
			},
		},
		{
			name: "pair<unsigned int, unsigned int>",
			want: rmeta.CxxTemplate{
				Name: "pair",
				Args: []string{"unsigned int", "unsigned int"},
			},
		},
		{
			name: "pair<pair<u T,u T>, pair<v T, v T>>",
			want: rmeta.CxxTemplate{
				Name: "pair",
				Args: []string{"pair<u T,u T>", "pair<v T, v T>"},
			},
		},
		{
			name:   "pair<T,>",
			panics: `rmeta: invalid empty type argument "pair<T,>"`,
		},
		{
			name:   "pair<,U>",
			panics: `rmeta: invalid empty type argument "pair<,U>"`,
		},
		{
			name:   "pair<,>",
			panics: `rmeta: invalid empty type argument "pair<,>"`,
		},
		{
			name:   "pair<T,U",
			panics: `rmeta: missing '>' in "pair<T,U"`,
		},
		// std::vector
		{
			name: "std::vector<T>",
			want: rmeta.CxxTemplate{
				Name: "std::vector",
				Args: []string{"T"},
			},
		},
		{
			name: "std::vector<T, alloc<T>>",
			want: rmeta.CxxTemplate{
				Name: "std::vector",
				Args: []string{"T", "alloc<T>"},
			},
		},
		{
			name: "std::vector<std::map<K,V,Cmp>>",
			want: rmeta.CxxTemplate{
				Name: "std::vector",
				Args: []string{"std::map<K,V,Cmp>"},
			},
		},
		{
			name: "vector<pair<int, int>>",
			want: rmeta.CxxTemplate{
				Name: "vector",
				Args: []string{"pair<int, int>"},
			},
		},
		{
			name: "vector<pair<string,string> >",
			want: rmeta.CxxTemplate{
				Name: "vector",
				Args: []string{"pair<string,string>"},
			},
		},
		{
			name: "std::set<T>",
			want: rmeta.CxxTemplate{
				Name: "std::set",
				Args: []string{"T"},
			},
		},
		{
			name: "std::multiset<T>",
			want: rmeta.CxxTemplate{
				Name: "std::multiset",
				Args: []string{"T"},
			},
		},
		{
			name: "std::unordered_set<T>",
			want: rmeta.CxxTemplate{
				Name: "std::unordered_set",
				Args: []string{"T"},
			},
		},
		{
			name: "std::map<K,V>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "V"},
			},
		},
		{
			name: "std::map<K, V>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "V"},
			},
		},
		{
			name: "std::multimap<K, V>",
			want: rmeta.CxxTemplate{
				Name: "std::multimap",
				Args: []string{"K", "V"},
			},
		},
		{
			name: "std::unordered_map<K, V>",
			want: rmeta.CxxTemplate{
				Name: "std::unordered_map",
				Args: []string{"K", "V"},
			},
		},
		{
			name: "map<K,V>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"K", "V"},
			},
		},
		{
			name: "std::map<unsigned long,unsigned int>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"unsigned long", "unsigned int"},
			},
		},
		{
			name: "std::map<K,std::vector<V> >",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			name: "std::map<K,std::vector<V>>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			name: "std::map<K, std::vector<V>>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::vector<V>"},
			},
		},
		{
			name: "map<string,o2::quality_control::core::CheckDefinition>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"string", "o2::quality_control::core::CheckDefinition"},
			},
		},
		{
			name: "map<string,o2::quality_control::core::CheckDefinition>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"string", "o2::quality_control::core::CheckDefinition"},
			},
		},
		{
			name: "std::map<K, std::map<K2,V2>>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"K", "std::map<K2,V2>"},
			},
		},
		{
			name: "std::map<std::map<K1,V1>, std::map<K2,V2>>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"std::map<K1,V1>", "std::map<K2,V2>"},
			},
		},
		{
			name: "std::map<std::map<K1,V1>, Foo<T1,T2,T3>>",
			want: rmeta.CxxTemplate{
				Name: "std::map",
				Args: []string{"std::map<K1,V1>", "Foo<T1,T2,T3>"},
			},
		},
		{
			name: "Foo<T1, T2, T3, std::vector<T4, std::allocator<T4>>>",
			want: rmeta.CxxTemplate{
				Name: "Foo",
				Args: []string{"T1", "T2", "T3", "std::vector<T4, std::allocator<T4>>"},
			},
		},
		{
			name: "Class<>",
			want: rmeta.CxxTemplate{
				Name: "Class",
			},
		},
		{
			name: "Class<T,1<2>",
			want: rmeta.CxxTemplate{
				Name: "Class",
				Args: []string{"T", "1<2"},
			},
		},
		{
			name: "Class<T,1>2>",
			want: rmeta.CxxTemplate{
				Name: "Class",
				Args: []string{"T", "1>2"},
			},
		},
		{
			name: "map<Cls<1>,Cls<2>>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"Cls<1>", "Cls<2>"},
			},
		},
		{
			name: "map<Cls<1>2>,Cls<2>3>>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"Cls<1>2>", "Cls<2>3>"},
			},
		},
		{
			name: "map<Cls<1>2>,Cls<2<3>>",
			want: rmeta.CxxTemplate{
				Name: "map",
				Args: []string{"Cls<1>2>", "Cls<2<3>"},
			},
		},
		// FIXME(sbinet) ?
		// for unknown classes, the ambiguity of C++ templates is
		// undecidable...
		// {
		// 	name: "map<Cls<1<2>,Cls<2>3>>",
		// 	want: rmeta.CxxTemplate{
		// 		Name: "map",
		// 		Args: []string{"Cls<1<2>", "Cls<2>3>"},
		// 	},
		// },
		// {
		// 	name: "map<Cls<1<2>,Cls<2<3>>",
		// 	want: rmeta.CxxTemplate{
		// 		Name: "map",
		// 		Args: []string{"Cls<1<2>", "Cls<2<3>"},
		// 	},
		// },
		// {
		// 	name: "Class<T,map<Class<1>2>,V>",
		// 	want: rmeta.CxxTemplate{
		// 		Name: "Class",
		// 		Args: []string{"T", "map<Class<1>2>,V>"},
		// 	},
		// },
		{
			name: "bitset<42>",
			want: rmeta.CxxTemplate{
				Name: "bitset",
				Args: []string{"42"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic (%s)", tc.panics)
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message: got=%s, want=%s", got, want)
					}
				}()
			}

			cxx := rmeta.CxxTemplateFrom(tc.name)
			if got, want := cxx, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid template args.\n got=%q\nwant=%q", got, want)
			}
		})
	}
}
