// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// disable test on windows because of symlinks
// +build !windows

package rootio

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestFileDirectory(t *testing.T) {
	for _, fname := range []string{
		"testdata/small-flat-tree.root",
		XrdRemote("testdata/small-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			f, err := Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			for _, table := range []struct {
				test  string
				value string
				want  string
			}{
				{"Name", f.Name(), "test-small.root"}, // name when created
				{"Title", f.Title(), "small event file"},
				{"Class", f.Class(), "TFile"},
			} {
				if table.value != table.want {
					t.Fatalf("%v: got=%q, want=%q", table.test, table.value, table.want)
				}
			}

			for _, table := range []struct {
				name string
				want bool
			}{
				{"tree", true},
				{"tree;0", false},
				{"tree;1", true},
				{"tree;9999", true},
				{"tree_nope", false},
				{"tree_nope;0", false},
				{"tree_nope;1", false},
				{"tree_nope;9999", false},
			} {
				_, err := f.Get(table.name)
				if (err == nil) != table.want {
					t.Fatalf("%s: got key (err=%v). want=%v", table.name, err, table.want)
				}
			}

			for _, table := range []struct {
				name string
				want string
			}{
				{"tree", "TTree"},
				{"tree;1", "TTree"},
			} {
				k, err := f.Get(table.name)
				if err != nil {
					t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
				}

				if k.Class() != table.want {
					t.Fatalf("%s: got key with class=%s (want=%s)", table.name, k.Class(), table.want)
				}
			}

			for _, table := range []struct {
				name string
				want string
			}{
				{"tree", "tree"},
				{"tree;1", "tree"},
			} {
				o, err := f.Get(table.name)
				if err != nil {
					t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
				}

				k := o.(Named)
				if k.Name() != table.want {
					t.Fatalf("%s: got key with name=%s (want=%v)", table.name, k.Name(), table.want)
				}
			}

			for _, table := range []struct {
				name string
				want string
			}{
				{"tree", "my tree title"},
				{"tree;1", "my tree title"},
			} {
				o, err := f.Get(table.name)
				if err != nil {
					t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
				}

				k := o.(Named)
				if k.Title() != table.want {
					t.Fatalf("%s: got key with title=%s (want=%v)", table.name, k.Title(), table.want)
				}
			}
		})
	}
}

func TestFileOpenStreamerInfo(t *testing.T) {
	for _, fname := range []string{
		"testdata/small-flat-tree.root",
		"testdata/simple.root",
		XrdRemote("testdata/small-flat-tree.root"),
		XrdRemote("testdata/simple.root"),
	} {
		f, err := Open(fname)
		if err != nil {
			t.Errorf("error opening %q: %v\n", fname, err)
			continue
		}
		defer f.Close()

		_ = f.StreamerInfos()
	}
}

func TestOpenEmptyFile(t *testing.T) {
	f, err := Open("testdata/uproot/issue70.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	si := f.StreamerInfos()
	if si != nil {
		t.Fatalf("expected no StreamerInfos in empty file")
	}
}

func TestCreate(t *testing.T) {

	rootls := "rootls"
	if runtime.GOOS == "windows" {
		rootls = "rootls.exe"
	}

	rootls, err := exec.LookPath(rootls)
	withROOTCxx := err == nil

	dir, err := ioutil.TempDir("", "rootio-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	type rooter interface {
		Object
		ROOTMarshaler
		ROOTUnmarshaler
	}

	for i, tc := range []struct {
		name string
		skip bool
		want []rooter
	}{
		{name: "", want: nil},
		{
			name: "TObjString",
			want: []rooter{NewObjString("hello")},
		},
		{
			name: "TObjString",
			want: []rooter{NewObjString("hello"), NewObjString("world")},
		},
		{
			name: "TObjString",
			want: func() []rooter {
				var out []rooter
				for _, i := range []int{0, 1, 253, 254, 255, 256, 512, 1024} {
					str := strings.Repeat("=", i)
					out = append(out, NewObjString(str))
				}
				return out
			}(),
		},
		{
			name: "TObject",
			want: []rooter{newObject()},
		},
		{
			name: "TNamed",
			want: []rooter{
				&tnamed{rvers: 1, obj: *newObject(), name: "n0", title: "t0"},
				&tnamed{rvers: 1, obj: *newObject(), name: "n1", title: "t1"},
				&tnamed{rvers: 1, obj: *newObject(), name: "n2", title: "t2"},
			},
		},
		{
			name: "TList",
			want: []rooter{&tlist{
				rvers: 5,
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  "list-name",
				objs: []Object{
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n0", title: "t0"},
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n1", title: "t1"},
					&tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "n2", title: "t2"},
				},
			}},
		},
		{
			name: "TArrayF",
			want: []rooter{
				&ArrayF{Data: []float32{1, 2, 3, 4, 5, 6}},
			},
		},
		{
			name: "TArrayD",
			want: []rooter{
				&ArrayD{Data: []float64{1, 2, 3, 4, 5, 6}},
			},
		},
		{
			name: "TArrays",
			want: []rooter{
				&ArrayF{Data: []float32{1, 2, 3, 4, 5, 6}},
				&ArrayD{Data: []float64{1, 2, 3, 4, 5, 6}},
			},
		},
		{
			name: "TAxis",
			want: []rooter{&taxis{
				rvers:  10,
				tnamed: tnamed{rvers: 1, obj: tobject{id: 0x0, bits: 0x3000000}, name: "xaxis", title: ""},
				attaxis: attaxis{
					rvers: 4,
					ndivs: 510, acolor: 1, lcolor: 1, lfont: 42, loffset: 0.005, lsize: 0.035, ticks: 0.03, toffset: 1, tsize: 0.035, tcolor: 1, tfont: 42,
				},
				nbins: 100, xmin: 0, xmax: 100,
				xbins: ArrayD{Data: nil},
				first: 0, last: 0, bits2: 0x0, time: false, tfmt: "",
				labels:  nil,
				modlabs: nil,
			}},
		},
	} {
		fname := filepath.Join(dir, fmt.Sprintf("out-%d.root", i))
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip()
			}

			w, err := Create(fname)
			if err != nil {
				t.Fatal(err)
			}

			for i := range tc.want {
				var (
					kname = fmt.Sprintf("key-%s-%02d", tc.name, i)
					want  = tc.want[i]
				)

				err = w.Put(kname, want)
				if err != nil {
					t.Fatal(err)
				}
			}

			if got, want := len(w.Keys()), len(tc.want); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			err = w.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			r, err := Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			if got, want := len(r.Keys()), len(tc.want); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			for i := range tc.want {
				var (
					kname = fmt.Sprintf("key-%s-%02d", tc.name, i)
					want  = tc.want[i]
				)

				rgot, err := r.Get(kname)
				if err != nil {
					t.Fatal(err)
				}

				if got := rgot.(rooter); !reflect.DeepEqual(got, want) {
					t.Fatalf("error reading back value[%d].\ngot = %#v\nwant = %#v", i, got, want)
				}
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			if !withROOTCxx {
				t.Logf("skip test with ROOT/C++")
				return
			}

			cmd := exec.Command(rootls, "-l", fname)
			err = cmd.Run()
			if err != nil {
				t.Fatalf("ROOT/C++ could not open file %q", fname)
			}
		})
	}
}

func TestFreeList(t *testing.T) {
	var list freeList

	checkSize := func(size int) {
		t.Helper()
		if got, want := len(list), size; got != want {
			t.Fatalf("got=%d, want=%d", got, want)
		}
	}
	checkSegment := func(free *freeSegment, want freeSegment) {
		t.Helper()
		if free == nil {
			t.Fatalf("expected a valid free segment")
		}
		if *free != want {
			t.Fatalf("got=%#v, want=%#v", *free, want)
		}
	}

	list.add(0, 1)
	checkSize(1)

	list.add(3, 10)
	checkSize(2)

	free := list.add(13, 20)
	checkSize(3)
	checkSegment(free, freeSegment{13, 20})

	free = list.add(12, 22)
	checkSize(3)
	checkSegment(free, freeSegment{12, 22})

	if got, want := list, (freeList{
		{0, 1},
		{3, 10},
		{12, 22},
	}); !reflect.DeepEqual(got, want) {
		t.Fatalf("error\ngot = %v\nwant= %v", got, want)
	}

	free = list.add(15, 20)
	checkSize(3)
	checkSegment(free, freeSegment{12, 22})

	free = list.add(40, 50)
	checkSize(4)
	checkSegment(free, freeSegment{40, 50})

	free = list.add(39, 40)
	checkSize(4)
	checkSegment(free, freeSegment{39, 50})

	free = list.add(37, 38)
	checkSize(4)
	checkSegment(free, freeSegment{37, 50})

	list.add(55, 60)
	list.add(65, 70)
	free = list.add(56, 66)
	checkSize(5)
	checkSegment(free, freeSegment{55, 70})

	free = list.add(54, 71)
	checkSize(5)
	checkSegment(free, freeSegment{54, 71})

	for _, tc := range []struct {
		list []freeSegment
		want []freeSegment
		free freeList
	}{
		{
			list: nil,
			want: nil,
			free: nil,
		},
		{
			list: []freeSegment{{0, 1}, {1, 2}},
			want: []freeSegment{{0, 1}, {0, 2}},
			free: freeList{{0, 2}},
		},
		{
			list: []freeSegment{{10, 12}, {10, 13}},
			want: []freeSegment{{10, 12}, {10, 13}},
			free: freeList{{10, 13}},
		},
	} {
		t.Run("", func(t *testing.T) {
			var list freeList
			for i, v := range tc.list {
				free := list.add(v.first, v.last)
				if !reflect.DeepEqual(*free, tc.want[i]) {
					t.Fatalf("error:\ngot[%d] = %#v\nwant[%d]= %#v\n",
						i, *free, i, tc.want[i],
					)
				}
			}
			if !reflect.DeepEqual(list, tc.free) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v\n", list, tc.free)
			}
		})
	}
}

func TestFreeListBest(t *testing.T) {
	for _, tc := range []struct {
		name   string
		nbytes int64
		list   freeList
		want   *freeSegment
	}{
		{
			name:   "empty",
			nbytes: 0,
			list:   nil,
			want:   nil,
		},
		{
			name:   "empty-list",
			nbytes: 10,
			list:   nil,
			want:   nil,
		},
		{
			name:   "exact-match",
			nbytes: 10,
			list:   freeList{{0, 1}, {10, 20 - 1}},
			want:   &freeSegment{10, 20 - 1},
		},
		{
			name:   "match",
			nbytes: 1,
			list:   freeList{{0, 10}},
			want:   &freeSegment{0, 10},
		},
		{
			name:   "match",
			nbytes: 10,
			list:   freeList{{0, 1}, {10, 20 + 4 + 1}},
			want:   &freeSegment{10, 20 + 4 + 1},
		},
		{
			name:   "big-file",
			nbytes: 10,
			list:   freeList{{0, 1}},
			want:   &freeSegment{0, 1000000001},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.list.best(tc.nbytes)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error\ngot = %#v\nwant= %#v\n", got, tc.want)
			}
		})
	}
}

func TestFreeListLast(t *testing.T) {
	for _, tc := range []struct {
		list freeList
		want *freeSegment
	}{
		{
			list: nil,
			want: nil,
		},
		{
			list: freeList{},
			want: nil,
		},
		{
			list: freeList{{0, kStartBigFile}},
			want: &freeSegment{0, kStartBigFile},
		},
		{
			list: freeList{{0, 10}, {12, kStartBigFile}},
			want: &freeSegment{12, kStartBigFile},
		},
	} {
		t.Run("", func(t *testing.T) {
			got := tc.list.last()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error\ngot = %#v\nwant= %#v\n", got, tc.want)
			}
		})
	}
}
