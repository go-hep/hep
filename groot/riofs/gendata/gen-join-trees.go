// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

const N = 10 // number of events to generate

func main() {
	dir := flag.String("d", "../testdata", "path to directory where to store output ROOT files")

	flag.Parse()

	err := gen(*dir)
	if err != nil {
		log.Fatalf("could not generate 'join' trees: %+v", err)
	}
}

func gen(dir string) error {
	for _, f := range []func(string) error{gen1, gen2, gen3, gen4} {
		err := f(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func gen1(dir string) error {
	fname := filepath.Join(dir, "join1.root")
	f, err := riofs.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var evt struct {
		f1 float64
		f2 int64
		f3 string
	}
	wvars := []rtree.WriteVar{
		{Name: "b10", Value: &evt.f1},
		{Name: "b11", Value: &evt.f2},
		{Name: "b12", Value: &evt.f3},
	}
	tree, err := rtree.NewWriter(f, "j1", wvars, rtree.WithTitle("j1-tree"))
	if err != nil {
		return err
	}
	defer tree.Close()

	for i := 0; i < N; i++ {
		evt.f1 = 100 + float64(i+1)
		evt.f2 = 100 + int64(i+1)
		evt.f3 = fmt.Sprintf("j1-%03d", 100+(i+1))
		_, err = tree.Write()
		if err != nil {
			return fmt.Errorf("could not write evt %d: %w", i, err)
		}
	}

	err = tree.Close()
	if err != nil {
		return fmt.Errorf("could not close tree: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}

	return nil
}

func gen2(dir string) error {
	fname := filepath.Join(dir, "join2.root")
	f, err := riofs.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var evt struct {
		f1 float64
		f2 int64
		f3 string
	}
	wvars := []rtree.WriteVar{
		{Name: "b20", Value: &evt.f1},
		{Name: "b21", Value: &evt.f2},
		{Name: "b22", Value: &evt.f3},
	}
	tree, err := rtree.NewWriter(f, "j2", wvars, rtree.WithTitle("j2-tree"))
	if err != nil {
		return err
	}
	defer tree.Close()

	for i := 0; i < N; i++ {
		evt.f1 = 200 + float64(i+1)
		evt.f2 = 200 + int64(i+1)
		evt.f3 = fmt.Sprintf("j2-%03d", 200+(i+1))
		_, err = tree.Write()
		if err != nil {
			return fmt.Errorf("could not write evt %d: %w", i, err)
		}
	}

	err = tree.Close()
	if err != nil {
		return fmt.Errorf("could not close tree: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}

	return nil
}

func gen3(dir string) error {
	fname := filepath.Join(dir, "join3.root")
	f, err := riofs.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	var evt struct {
		f1 float64
		f2 int64
		f3 string
	}
	wvars := []rtree.WriteVar{
		{Name: "b30", Value: &evt.f1},
		{Name: "b31", Value: &evt.f2},
		{Name: "b32", Value: &evt.f3},
	}
	tree, err := rtree.NewWriter(f, "j3", wvars, rtree.WithTitle("j3-tree"))
	if err != nil {
		return err
	}
	defer tree.Close()

	for i := 0; i < N; i++ {
		evt.f1 = 300 + float64(i+1)
		evt.f2 = 300 + int64(i+1)
		evt.f3 = fmt.Sprintf("j3-%03d", 300+(i+1))
		_, err = tree.Write()
		if err != nil {
			return fmt.Errorf("could not write evt %d: %w", i, err)
		}
	}

	err = tree.Close()
	if err != nil {
		return fmt.Errorf("could not close tree: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}

	return nil
}

func gen4(dir string) error {
	fname := filepath.Join(dir, "join4.root")
	f, err := riofs.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = func() error {
		var evt struct {
			f1 float64
			f2 int64
			f3 string
		}
		wvars := []rtree.WriteVar{
			{Name: "b40", Value: &evt.f1},
			{Name: "b41", Value: &evt.f2},
			{Name: "b42", Value: &evt.f3},
		}
		tree, err := rtree.NewWriter(f, "j41", wvars, rtree.WithTitle("j4-1-tree (evtmax differ)"))
		if err != nil {
			return err
		}
		defer tree.Close()

		for i := 0; i < N+1; i++ {
			evt.f1 = 400 + float64(i+1)
			evt.f2 = 400 + int64(i+1)
			evt.f3 = fmt.Sprintf("j4-1-%03d", 400+(i+1))
			_, err = tree.Write()
			if err != nil {
				return fmt.Errorf("could not write evt %d: %w", i, err)
			}
		}

		err = tree.Close()
		if err != nil {
			return fmt.Errorf("could not close tree: %w", err)
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("could not generate j41 tree: %w", err)
	}

	err = func() error {
		var evt struct {
			f1 float64
			f2 int32 // different type than other evt structs
			f3 string
		}
		wvars := []rtree.WriteVar{
			{Name: "b40", Value: &evt.f1},
			{Name: "b11", Value: &evt.f2},
			{Name: "b22", Value: &evt.f3},
		}
		tree, err := rtree.NewWriter(f, "j42", wvars, rtree.WithTitle("j4-2-tree (branch collision w/ j1, j2))"))
		if err != nil {
			return err
		}
		defer tree.Close()

		for i := 0; i < N; i++ {
			evt.f1 = 400 + float64(i+1)
			evt.f2 = 400 + int32(i+1)
			evt.f3 = fmt.Sprintf("j4-2-%03d", 400+(i+1))
			_, err = tree.Write()
			if err != nil {
				return fmt.Errorf("could not write evt %d: %w", i, err)
			}
		}

		err = tree.Close()
		if err != nil {
			return fmt.Errorf("could not close tree: %w", err)
		}
		return nil
	}()
	if err != nil {
		return fmt.Errorf("could not generate j4 tree: %w", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file: %w", err)
	}

	return nil
}
