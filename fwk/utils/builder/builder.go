// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package builder builds a fwk-app binary from a list of go files.
//
// builder's architecture and sources are heavily inspired from golint:
//
//	https://github.com/golang/lint
package builder // import "go-hep.org/x/hep/fwk/utils/builder"

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type file struct {
	app  *Builder
	f    *ast.File
	fset *token.FileSet
	src  []byte
	name string
}

// Builder generates and builds fwk-based applications.
type Builder struct {
	fset  *token.FileSet
	files map[string]*file

	pkg  *types.Package
	info *types.Info

	funcs []string

	Name  string // name of resulting compiled binary
	Usage string // usage string displayed by compiled binary (with -help)
}

// NewBuilder creates a Builder from a list of file names or directories
func NewBuilder(fnames ...string) (*Builder, error) {
	var err error

	b := &Builder{
		fset:  token.NewFileSet(),
		files: make(map[string]*file, len(fnames)),
		funcs: make([]string, 0),
		Usage: `Usage: %[1]s [options] <input> <output>

ex:
 $ %[1]s -l=INFO -evtmax=-1 input.dat output.dat

options:
`,
	}

	for _, fname := range fnames {
		fi, err := os.Stat(fname)
		if err != nil {
			return nil, fmt.Errorf("builder: could not stat %q: %w", fname, err)
		}
		fm := fi.Mode()
		if fm.IsRegular() {
			src, err := os.ReadFile(fname)
			if err != nil {
				return nil, fmt.Errorf("builder: could not read %q: %w", fname, err)
			}
			f, err := parser.ParseFile(b.fset, fname, src, parser.ParseComments)
			if err != nil {
				return nil, fmt.Errorf("builder: could not parse file %q: %w", fname, err)
			}
			b.files[fname] = &file{
				app:  b,
				f:    f,
				fset: b.fset,
				src:  src,
				name: fname,
			}
		}
		if fm.IsDir() {
			return nil, fmt.Errorf("directories not (yet) handled (got=%q)", fname)
		}
	}
	return b, err
}

// Build applies some type-checking, collects setup functions and generates the sources of the fwk-based application.
func (b *Builder) Build() error {
	var err error

	if b.Name == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("builder: could not fetch current work directory: %w", err)
		}
		b.Name = filepath.Base(pwd)
	}

	err = b.doTypeCheck()
	if err != nil {
		return err
	}

	// check we build a 'main' package
	if !b.isMain() {
		return fmt.Errorf("not a 'main' package")
	}

	err = b.scanSetupFuncs()
	if err != nil {
		return err
	}

	if len(b.funcs) <= 0 {
		return fmt.Errorf("no setup function found")
	}

	err = b.genSources()
	if err != nil {
		return err
	}

	return err
}

func (b *Builder) doTypeCheck() error {
	var err error
	config := &types.Config{
		// By setting a no-op error reporter, the type checker does
		// as much work as possible.
		Error: func(error) {},
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	var anyFile *file
	var astFiles []*ast.File
	for _, f := range b.files {
		anyFile = f
		astFiles = append(astFiles, f.f)
	}
	pkg, err := config.Check(anyFile.f.Name.Name, b.fset, astFiles, info)
	// Remember the typechecking info, even if config.Check failed,
	// since we will get partial information.
	b.pkg = pkg
	b.info = info
	return err
}

func (b *Builder) typeOf(expr ast.Expr) types.Type {
	if b.info == nil {
		return nil
	}
	return b.info.TypeOf(expr)
}

func (b *Builder) scanSetupFuncs() error {
	var err error

	// setupfunc := types.New("func (*job.Job)")
	// fmt.Fprintf(os.Stderr, "looking for type: %#v...\n", setupfunc)

	for _, f := range b.files {
		// fmt.Fprintf(os.Stderr, ">>> [%s]...\n", f.name)
		f.walk(func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// fmt.Fprintf(os.Stderr, "file: %q - func=%s\n", f.name, fn.Name.Name)

			if fn.Recv != nil {
				return true
			}

			if fn.Type.Results != nil && fn.Type.Results.NumFields() != 0 {
				// fmt.Fprintf(os.Stderr,
				// 	"file: %q - func=%s [results=%d]\n",
				// 	f.name, fn.Name.Name,
				// 	fn.Type.Results.NumFields(),
				// )
				return true
			}

			if fn.Type.Params == nil {
				// fmt.Fprintf(os.Stderr,
				// 	"file: %q - func=%s [params=nil]\n",
				// 	f.name, fn.Name.Name,
				// )
				return true
			}

			if fn.Type.Params.NumFields() != 1 {
				// fmt.Fprintf(os.Stderr,
				// 	"file: %q - func=%s [params=%d]\n",
				// 	f.name, fn.Name.Name,
				// 	fn.Type.Params.NumFields(),
				// )
				return true
			}

			param := b.typeOf(fn.Type.Params.List[0].Type)
			// FIXME(sbinet)
			//  - go the extra mile and create a proper type.Type from type.New("func(*job.Job)")
			//  - compare the types
			if param.String() != "*go-hep.org/x/hep/fwk/job.Job" {
				// fmt.Fprintf(os.Stderr,
				// 	"file: %q - func=%s [invalid type=%s]\n",
				// 	f.name, fn.Name.Name,
				// 	param.String(),
				// )
				return true
			}

			// fmt.Fprintf(os.Stderr, "file: %q - func=%s [ok]\n", f.name, fn.Name.Name)

			b.funcs = append(b.funcs, fn.Name.Name)
			return false
		})
	}
	return err
}

func (b *Builder) isMain() bool {
	for _, f := range b.files {
		if f.isMain() {
			return true
		}
	}
	return false
}

func (b *Builder) genSources() error {
	var err error
	tmpdir, err := os.MkdirTemp("", "fwk-builder-")
	if err != nil {
		return fmt.Errorf("builder: could not create tmpdir: %w", err)
	}
	defer os.RemoveAll(tmpdir)
	// fmt.Fprintf(os.Stderr, "tmpdir=[%s]...\n", tmpdir)

	// copy sources to dst
	for _, f := range b.files {
		// FIXME(sbinet)
		// only take base. watch out for duplicates!
		fname := filepath.Base(f.name)
		dstname := filepath.Join(tmpdir, fname)
		dst, err := os.Create(dstname)
		if err != nil {
			return fmt.Errorf("builder: could not create dst: %w", err)
		}
		defer dst.Close()

		_, err = dst.Write(f.src)
		if err != nil {
			return fmt.Errorf("builder: could not write dst: %w", err)
		}

		err = dst.Close()
		if err != nil {
			return fmt.Errorf("builder: could not close dst: %w", err)
		}
	}

	// add main.
	f, err := os.Create(filepath.Join(tmpdir, "main.go"))
	if err != nil {
		return fmt.Errorf("builder: could not create main: %w", err)
	}
	defer f.Close()

	data := struct {
		Usage      string
		Name       string
		SetupFuncs []string
	}{
		Usage:      b.Usage,
		Name:       b.Name,
		SetupFuncs: b.funcs,
	}

	err = render(f, tmpl, data)
	if err != nil {
		return fmt.Errorf("builder: could not render: %w", err)
	}

	build := exec.Command(
		"go", "build", "-o", b.Name, ".",
	)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	build.Dir = tmpdir
	err = build.Run()
	if err != nil {
		return fmt.Errorf("builder: could not build %q: %w", b.Name, err)
	}

	// copy final binary.
	{
		src, err := os.Open(filepath.Join(tmpdir, b.Name))
		if err != nil {
			return fmt.Errorf("builder: could not open src %q: %w", filepath.Join(tmpdir, b.Name), err)
		}
		defer src.Close()
		fi, err := src.Stat()
		if err != nil {
			return fmt.Errorf("builder: could not stat src %q: %w", src.Name(), err)
		}

		dst, err := os.OpenFile(b.Name, os.O_CREATE|os.O_WRONLY, fi.Mode())
		if err != nil {
			return fmt.Errorf("builder: could not open dst %q: %w", b.Name, err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return fmt.Errorf("builder: could not copy src to dst: %w", err)
		}

		err = dst.Close()
		if err != nil {
			return fmt.Errorf("builder: could not close %q: %w", dst.Name(), err)
		}
	}

	return err
}

func (f *file) isMain() bool {
	return f.f.Name.Name == "main"
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.f)
}

// walker adapts a function to satisfy the ast.Visitor interface.
// The function return whether the walk should proceed into the node's children.
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}
