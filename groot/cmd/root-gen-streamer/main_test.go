// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rdatatest"
	"go-hep.org/x/hep/internal/diff"
)

var (
	_ rdatatest.Event // make sure rdatatest is compiled
)

func TestGenerate(t *testing.T) {

	for _, tc := range []struct {
		pkg   string
		types []string
		want  string
	}{
		{
			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
			types: []string{"Event", "HLV", "Particle"},
			want:  "testdata/rdatatest.txt",
		},
	} {
		t.Run(tc.pkg, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := generate(buf, tc.pkg, tc.types)
			if err != nil {
				t.Fatalf("could not generate streamer: %v", err)
			}
			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference streamer: %v", err)
			}

			if got, want := buf.String(), string(want); got != want {
				t.Fatalf("error:\n%s\n", diff.Format(got, want))
			}
		})
	}
}

func TestGenerateCompileRun(t *testing.T) {
	for _, tc := range []struct {
		name  string
		pkg   string
		types []string
		out   string
		tmpl  string
		want  string
	}{
		{
			name:  "builtins",
			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
			types: []string{"Builtins"},
			out:   "../../internal/rdatatest/pkg_gen.go",
			tmpl:  "NewBuiltins",
			want: `>>> file[testdata/out.root]
key[000]: data;1 "" (go_hep_org::x::hep::groot::internal::rdatatest::Builtins) => &{true 8 16 32 64 -8 -16 -32 -64 32.32 64.64 builtins}
`,
		},
		{
			name:  "arr-builtins",
			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
			types: []string{"ArrBuiltins"},
			out:   "../../internal/rdatatest/pkg_gen.go",
			tmpl:  "NewArrBuiltins",
			want: `>>> file[testdata/out.root]
key[000]: data;1 "" (go_hep_org::x::hep::groot::internal::rdatatest::ArrBuiltins) => &{[true false] [8 88] [16 1616] [32 3232] [64 6464] [-8 -88] [-16 -1616] [-32 -3232] [-64 -6464] [32.32 -32.32] [64.64 64.64] [builtins arrays]}
`,
		},
		{
			name:  "struct-t1",
			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
			types: []string{"HLV", "T1"}, // FIXME(sbinet): only select T1 and let root-gen-streamer pick-up HLV
			out:   "../../internal/rdatatest/pkg_gen.go",
			tmpl:  "NewT1",
			want: `>>> file[testdata/out.root]
key[000]: data;1 "" (go_hep_org::x::hep::groot::internal::rdatatest::T1) => &{hello {1 2 3 4}}
`,
		},
		//		{
		//			name:  "struct-t2",
		//			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
		//			types: []string{"HLV", "T2"}, // FIXME(sbinet): only select T2 and let root-gen-streamer pick-up HLV
		//			out:   "../../internal/rdatatest/pkg_gen.go",
		//			tmpl:  "NewT2",
		//			want: `>>> file[testdata/out.root]
		//key[000]: data;1 "" (go_hep_org::x::hep::groot::internal::rdatatest::T1) => &{hello {1 2 3 4}}
		//`,
		//		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			os.Remove(tc.out)
			os.Remove("testdata/run.go")
			defer os.Remove(tc.out)
			defer os.Remove("testdata/run.go")
			defer os.Remove("testdata/out.root")

			var (
				out = new(bytes.Buffer)
				err error
				cmd *exec.Cmd
			)

			out.Reset()
			cmd = exec.Command("go", "get", "-v", tc.pkg)
			cmd.Stdout = out
			cmd.Stderr = out
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not compile package with streamer data:\n%v\nerr: %v", out.String(), err)
			}

			out.Reset()
			cmd = exec.Command(
				"root-gen-streamer",
				"-p", tc.pkg, "-t", strings.Join(tc.types, ","),
				"-o", tc.out,
			)
			cmd.Stdout = out
			cmd.Stderr = out
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not generate streamer data:\n%v\nerr: %v", out.String(), err)
			}

			out.Reset()
			cmd = exec.Command("go", "get", "-v", tc.pkg)
			cmd.Stdout = out
			cmd.Stderr = out
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not recompile package with streamer data:\n%v\nerr: %v", out.String(), err)
			}

			err = os.WriteFile("testdata/run.go", []byte(fmt.Sprintf(`//go:build ignore

package main

import (
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/internal/rdatatest"
)

func main() {
	f, err := groot.Create("testdata/out.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v := rdatatest.%s()
	err = f.Put("data", v)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}
`, tc.tmpl,
			)), 0644)

			if err != nil {
				t.Fatalf("could not generate test-write program: %v", err)
			}

			out.Reset()
			cmd = exec.Command("go", "run", "testdata/run.go")
			cmd.Stdout = out
			cmd.Stderr = out
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not run test-write program:\n%v\nerr: %v\n", out.String(), err)
			}

			out.Reset()
			cmd = exec.Command("root-dump", "testdata/out.root")
			cmd.Stdout = out
			cmd.Stderr = out
			err = cmd.Run()
			if err != nil {
				t.Fatalf("could not run root-dump:\n%v\nerr: %v", out.String(), err)
			}

			if got, want := out.String(), tc.want; got != want {
				t.Fatalf("error:\n%v", diff.Format(got, want))
			}
		})
	}
}
