// Copyright 2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"reflect"
	"testing"
)

func TestParseStdContainers(t *testing.T) {
	for _, tc := range []struct {
		name   string
		parse  func(string) []string
		want   []string
		panics string
	}{
		// std::pair
		{
			name:  "pair<int,int>",
			parse: parseStdPair,
			want:  []string{"int", "int"},
		},
		{
			name:  "pair<int, int>",
			parse: parseStdPair,
			want:  []string{"int", "int"},
		},
		{
			name:  "pair<unsigned int,int>",
			parse: parseStdPair,
			want:  []string{"unsigned int", "int"},
		},
		{
			name:  "pair<int,unsigned int>",
			parse: parseStdPair,
			want:  []string{"int", "unsigned int"},
		},
		{
			name:  "pair<unsigned int,unsigned int>",
			parse: parseStdPair,
			want:  []string{"unsigned int", "unsigned int"},
		},
		{
			name:  "pair< unsigned int, unsigned int >",
			parse: parseStdPair,
			want:  []string{"unsigned int", "unsigned int"},
		},
		{
			name:  "pair<pair<unsigned int,int>, pair<float,unsigned int> >",
			parse: parseStdPair,
			want:  []string{"pair<unsigned int,int>", "pair<float,unsigned int>"},
		},
		{
			name:  "pair<int, pair<int,int>>",
			parse: parseStdPair,
			want:  []string{"int", "pair<int,int>"},
		},
		{
			name:   "pair<int>",
			parse:  parseStdPair,
			panics: `rdict: invalid std::pair template "pair<int>"`,
		},
		{
			name:   "pair<t,>",
			parse:  parseStdPair,
			panics: `rmeta: invalid empty type argument "pair<t,>"`,
		},
		{
			name:   "pair<,u>",
			parse:  parseStdPair,
			panics: `rmeta: invalid empty type argument "pair<,u>"`,
		},
		{
			name:   "pair<,>",
			parse:  parseStdPair,
			panics: `rmeta: invalid empty type argument "pair<,>"`,
		},
		{
			name:   "pair<>",
			parse:  parseStdPair,
			panics: `rdict: invalid std::pair template "pair<>"`,
		},
		{
			name:   "pair< >",
			parse:  parseStdPair,
			panics: `rdict: invalid std::pair template "pair< >"`,
		},
		{
			name:   "pair<t,u",
			parse:  parseStdPair,
			panics: `rmeta: missing '>' in "pair<t,u"`,
		},
		// std::vector
		{
			name:  "vector<int>",
			parse: parseStdVector,
			want:  []string{"int"},
		},
		{
			name:  "std::vector<int>",
			parse: parseStdVector,
			want:  []string{"int"},
		},
		{
			name:  "vector<vector<int>>",
			parse: parseStdVector,
			want:  []string{"vector<int>"},
		},
		{
			name:  "vector<int,allocator<int>>",
			parse: parseStdVector,
			want:  []string{"int", "allocator<int>"},
		},
		{
			name:  "vector<map<int,long int>>",
			parse: parseStdVector,
			want:  []string{"map<int,long int>"},
		},
		{
			name:  "vector<pair<string,string> >",
			parse: parseStdVector,
			want:  []string{"pair<string,string>"},
		},
		{
			name:   "vector<int",
			parse:  parseStdVector,
			panics: `rmeta: missing '>' in "vector<int"`,
		},
		{
			name:   "xvector<int>",
			parse:  parseStdVector,
			panics: `rdict: invalid std::vector template "xvector<int>"`,
		},
		{
			name:   "vector<>",
			parse:  parseStdVector,
			panics: `rdict: invalid empty type argument "vector<>"`,
		},
		{
			name:   "vector<t1,t2,t3>",
			parse:  parseStdVector,
			panics: `rdict: invalid std::vector template "vector<t1,t2,t3>"`,
		},
		// std::map
		{
			name:  "map< int , int >",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "std::map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,string>",
			parse: parseStdMap,
			want:  []string{"int", "string"},
		},
		{
			name:  "map<int,vector<int>>",
			parse: parseStdMap,
			want:  []string{"int", "vector<int>"},
		},
		{
			name:  "map<int,vector<int> >",
			parse: parseStdMap,
			want:  []string{"int", "vector<int>"},
		},
		{
			name:  "map<int,map<string,int> >",
			parse: parseStdMap,
			want:  []string{"int", "map<string,int>"},
		},
		{
			name:  "map<map<string,int>, int>",
			parse: parseStdMap,
			want:  []string{"map<string,int>", "int"},
		},
		{
			name:  "map<map<string,int>, map<int,string>>",
			parse: parseStdMap,
			want:  []string{"map<string,int>", "map<int,string>"},
		},
		{
			name:  "map<long int,long int>",
			parse: parseStdMap,
			want:  []string{"long int", "long int"},
		},
		{
			name:  "map<long int, vector<long int>, allocator<pair<const long int, vector<long int>>>",
			parse: parseStdMap,
			want:  []string{"long int", "vector<long int>", "allocator<pair<const long int, vector<long int>>"},
		},
		{
			name:   "map<k,v",
			parse:  parseStdMap,
			panics: `rmeta: missing '>' in "map<k,v"`,
		},
		{
			name:   "map<k,v,c,a,XXX>",
			parse:  parseStdMap,
			panics: `rdict: invalid std::map template "map<k,v,c,a,XXX>"`,
		},
		{
			name:   "map<>",
			parse:  parseStdMap,
			panics: `rdict: invalid empty type argument "map<>"`,
		},
		{
			name:   "map<k,>",
			parse:  parseStdMap,
			panics: `rmeta: invalid empty type argument "map<k,>"`,
		},
		{
			name:   "map<,v>",
			parse:  parseStdMap,
			panics: `rmeta: invalid empty type argument "map<,v>"`,
		},
		{
			name:   "map<,>",
			parse:  parseStdMap,
			panics: `rmeta: invalid empty type argument "map<,>"`,
		},
		{
			name:   "xmap<k,v>",
			parse:  parseStdMap,
			panics: `rdict: invalid std::map template "xmap<k,v>"`,
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
			got := tc.parse(tc.name)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}
