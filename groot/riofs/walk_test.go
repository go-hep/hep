// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"os"
	stdpath "path"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/root"
)

func TestGet(t *testing.T) {
	f, err := Open("../testdata/dirs-6.14.00.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	h1, err := Get[rhist.H1](f, "dir1/dir11/h1")
	if err != nil {
		t.Fatalf("could not get histo: %+v", err)
	}
	if h1 == nil {
		t.Fatalf("invalid H1 value")
	}

	h1f, err := Get[*rhist.H1F](f, "dir1/dir11/h1")
	if err != nil {
		t.Fatalf("could not get histo: %+v", err)
	}
	if h1f == nil {
		t.Fatalf("invalid H1F value")
	}

	_, err = Get[*rhist.H1D](f, "dir1/dir11/h1")
	if err == nil {
		t.Fatalf("expected an error")
	}
	if want := fmt.Errorf(`riofs: could not convert "dir1/dir11/h1" (*rhist.H1F) to *rhist.H1D`); err.Error() != want.Error() {
		t.Fatalf("invalid error:\ngot= %+v\nwant=%+v", err, want)
	}

	_, err = Get[any](f, "dir1/dir11/h1")
	if err != nil {
		t.Fatalf("could not get histo: %+v", err)
	}

	_, err = Get[any](f, "dir1/dir11/h1_XXX")
	if err == nil {
		t.Fatalf("expected an error")
	}
	if want := fmt.Errorf(`riofs: dir11: could not find key "h1_XXX;9999"`); err.Error() != want.Error() {
		t.Fatalf("invalid error:\ngot= %+v\nwant=%+v", err, want)
	}

}
func TestDir(t *testing.T) {
	f, err := Open("../testdata/dirs-6.14.00.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	rd := Dir(f)
	for _, tc := range []struct {
		path  string
		class string
	}{
		{"dir1/dir11/h1", "TH1F"},
		{"dir1/dir11/h1;1", "TH1F"},
		{"dir1/dir11/h1;9999", "TH1F"},
		{"/dir1/dir11/h1", "TH1F"},
		{"/dir1/dir11/h1;1", "TH1F"},
		{"/dir1/dir11/h1;9999", "TH1F"},
		{"dir1/dir11", "TDirectoryFile"},
		{"dir1/dir11;1", "TDirectoryFile"},
		{"dir1/dir11;9999", "TDirectoryFile"},
		{"dir1", "TDirectoryFile"},
		{"dir2", "TDirectoryFile"},
		{"dir3", "TDirectoryFile"},
		{"", "TFile"},
		{"/", "TFile"},
		{"/dir1", "TDirectoryFile"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			o, err := rd.Get(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := o.Class(), tc.class; got != want {
				t.Fatalf("got=%q, want=%q", got, want)
			}
		})
	}

	keys := make([]string, len(rd.Keys()))
	for i, k := range rd.Keys() {
		keys[i] = k.Name()
	}

	if got, want := keys, []string{"dir1", "dir2", "dir3"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid keys:\ngot = %v\nwant=%v\n", got, want)
	}

	for _, tc := range []struct {
		path   string
		parent string
	}{
		{"dir1/dir11", "dir1"},
		{"/dir1/dir11", "dir1"},
		{"dir1", f.Name()},
		{"/dir1", f.Name()},
		{"", ""},
		{"/", ""},
	} {
		t.Run("parent:"+tc.path, func(t *testing.T) {
			o, err := rd.Get(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			p := o.(Directory).Parent()
			switch p {
			case nil:
				if got, want := "", tc.parent; got != want {
					t.Fatalf("invalid parent: got=%q, want=%q", got, want)
				}
			default:
				if got, want := p.(root.Named).Name(), tc.parent; got != want {
					t.Fatalf("invalid parent: got=%q, want=%q", got, want)
				}
			}
		})
	}

}

func TestRecDirMkdir(t *testing.T) {
	tmp, err := os.CreateTemp("", "groot-riofs-")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	os.Remove(tmp.Name())

	f, err := Create(tmp.Name() + ".root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	rd := Dir(f)

	display := func() string {
		o := new(strings.Builder)
		err := Walk(f, func(path string, obj root.Object, err error) error {
			fmt.Fprintf(o, "%s (%s)\n", path, obj.Class())
			return nil
		})
		if err != nil {
			return fmt.Errorf("could not display file content: %w", err).Error()
		}
		return o.String()
	}

	for _, tc := range []struct {
		path string
		err  error
	}{
		{path: "dir1"},
		{path: "dir2/dir21/dir211"},
		{path: "dir2/dir22"},
		{path: "dir2/dir22/dir222"},
		{path: "/dir3"},
		{path: "/dir3"}, // recursive mkdir does not fail.
		{path: "/dir4/dir44"},
		{path: "/", err: fmt.Errorf("riofs: invalid path \"/\" to Mkdir")},
		{path: "", err: fmt.Errorf("riofs: invalid path \"\" to Mkdir")},
	} {
		t.Run(tc.path, func(t *testing.T) {
			_, err := rd.Mkdir(tc.path)
			switch err {
			case nil:
				if tc.err != nil {
					t.Fatalf("got no error, want=%v\ncontent:\n%v", tc.err, display())
				}
			default:
				if tc.err == nil {

					t.Fatalf("could not create %q: %v\ncontent:\n%v", tc.path, err, display())
				}
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error.\ngot= %v\nwant=%v\ncontent:\n%v", got, want, display())
				}
			}
		})
	}

	// test recursive mkdir does not work on f.
	_, err = f.Mkdir("xdir/xsubdir")
	if err == nil {
		t.Fatalf("expected an error, got=%v\ncontent:\n%v", err, display())
	}
	if got, want := err.Error(), fmt.Errorf("riofs: invalid directory name %q (contains a '/')", "xdir/xsubdir").Error(); got != want {
		t.Fatalf("invalid error. got=%q, want=%q", got, want)
	}

	// test regular mkdir fails when directory already exists
	_, err = f.Mkdir("dir1")
	if err == nil {
		t.Fatalf("expected an error, got=%v\ncontent:\n%v", err, display())
	}
	if got, want := err.Error(), fmt.Errorf("riofs: %q already exists", "dir1").Error(); got != want {
		t.Fatalf("invalid error. got=%q, want=%q", got, want)
	}
}

func TestRecDirPut(t *testing.T) {
	tmp, err := os.CreateTemp("", "groot-riofs-")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	os.Remove(tmp.Name())

	f, err := Create(tmp.Name() + ".root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	rd := Dir(f)

	display := func() string {
		o := new(strings.Builder)
		err := Walk(f, func(path string, obj root.Object, err error) error {
			if err != nil {
				return err
			}
			name := path[len(f.Name()):]
			if name == "" {
				fmt.Fprintf(o, "%s (%s)\n", path, obj.Class())
				return nil
			}
			dir, err := Dir(f).Get(stdpath.Dir(name))
			if err != nil {
				return err
			}
			pdir := dir.(Directory)
			cycle := -1
			for _, k := range pdir.Keys() {
				if k.Name() == stdpath.Base(path) {
					cycle = k.Cycle()
					break
				}
			}
			fmt.Fprintf(o, "%s;%d (%s)\n", path, cycle, obj.Class())
			return nil
		})
		if err != nil {
			return fmt.Errorf("could not display file content: %w", err).Error()
		}
		return o.String()
	}

	for _, tc := range []struct {
		path  string
		obj   string
		cycle int
		err   error
	}{
		{path: "dir1"},
		{path: "dir2/dir21/dir211"},
		{path: "dir2/dir22"},
		{path: "dir2/dir22/dir222"},
		{path: "/dir3"},
		{path: "/dir4/dir44"},

		{path: "/dir5/dir55"},
		{path: "/dir5", obj: "dir55", err: keyTypeError{key: "dir55", class: "TDirectory"}},

		{path: "/dir5/dir55", cycle: 2}, // recreating the same object is ok
	} {
		t.Run(tc.path, func(t *testing.T) {
			obj := tc.obj
			if obj == "" {
				obj = "obj"
			}
			err := rd.Put(stdpath.Join(tc.path, obj), rbase.NewObjString(obj))
			switch err {
			case nil:
				if tc.err != nil {
					t.Fatalf("got no error, want=%v\ncontent:\n%v", tc.err, display())
				}
				cycle := 1
				if tc.cycle != 0 {
					cycle = tc.cycle
				}
				name := stdpath.Join(tc.path, obj)
				namecycle := fmt.Sprintf("%s;%d", name, cycle)
				_, err := rd.Get(namecycle)
				if err != nil {
					t.Fatalf("could not access %q: %v\ncontent:\n%v", namecycle, err, display())
				}
			default:
				if tc.err == nil {

					t.Fatalf("could not create %q: %v\ncontent:\n%v", tc.path, err, display())
				}
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error.\ngot= %v\nwant=%v\ncontent:\n%v", got, want, display())
				}
			}
		})
	}

	// test recursive put does not work on f.
	err = f.Put("xdir/xsubdir/obj", rbase.NewObjString("obj"))
	if err == nil {
		t.Fatalf("expected an error, got=%v\ncontent:\n%v", err, display())
	}
	if got, want := err.Error(), fmt.Errorf("riofs: invalid path name %q (contains a '/')", "xdir/xsubdir/obj").Error(); got != want {
		t.Fatalf("invalid error. got=%q, want=%q", got, want)
	}

	err = rd.Put("", rbase.NewObjString("obj-empty-key"))
	if err != nil {
		t.Fatalf("could not create key-val with empty name: %v", err)
	}
}

func TestFileOf(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-riofs-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	f, err := Create(stdpath.Join(tmp, "file.root"))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer f.Close()

	dir111, err := Dir(f).Mkdir("dir-1/dir-11/dir-111")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, tc := range []struct {
		name   string
		dir    Directory
		panics string
	}{
		{
			name: "file",
			dir:  f,
		},
		{
			name: "file-rec",
			dir:  Dir(f),
		},
		{
			name: "file-dir",
			dir:  &f.dir,
		},
		{
			name: "dir-111",
			dir:  dir111,
		},
		{
			name:   "panics",
			dir:    &unknownDirImpl{},
			panics: "riofs: unknown Directory type *riofs.unknownDirImpl",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message. got=%q, want=%q", got, want)
					}
				}()
			}
			got := fileOf(tc.dir)
			if got != f {
				t.Fatalf("could not retrieve correct file for %q", tc.name)
			}
		})
	}
}

type unknownDirImpl struct{}

func (dir *unknownDirImpl) Get(namecycle string) (root.Object, error) { panic("not implemented") }
func (dir *unknownDirImpl) Put(name string, v root.Object) error      { panic("not implemented") }
func (dir *unknownDirImpl) Keys() []Key                               { panic("not implemented") }
func (dir *unknownDirImpl) Mkdir(name string) (Directory, error)      { panic("not implemented") }
func (dir *unknownDirImpl) Parent() Directory                         { return nil }

var (
	_ Directory = (*unknownDirImpl)(nil)
)
