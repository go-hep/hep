// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type document struct {
	Options struct {
		Syntax           bool `yaml:"getSyntax"`
		ExposePODMembers bool `yaml:"exposePODMembers"`
	} `yaml:"options"`

	Components yaml.Node `yaml:"components"`
	DataTypes  yaml.Node `yaml:"datatypes"`
}

type dtype struct {
	descr            string
	author           string `yaml:"Author"`
	members          []member
	vecmbrs          []member
	one2oneRels      []member
	one2manyRels     []member
	TransientMembers []string  `yaml:"TransientMembers"`
	Typedefs         []string  `yaml:"Typedefs"`
	extraCode        extraCode `yaml:"ExtraCode"`
	constCode        extraCode `yaml:"ConstExtraCode"`
}

type member struct {
	Name string
	Type string
	Doc  string
}

type extraCode struct {
	Includes  string `yaml:"includes"`
	ConstDecl string `yaml:"const_declaration"`
	Decl      string `yaml:"declaration"`
	Impl      string `yaml:"implementation"`
}

var (
	builtins = map[string]string{
		"bool":               "bool",
		"short":              "int16",
		"int":                "int32",
		"long":               "int64",
		"long long":          "int64",
		"unsigned int":       "uint32",
		"unsigned":           "uint32",
		"unsigned long":      "uint64",
		"unsigned long long": "uint64",
		"float":              "float32",
		"double":             "float64",
		"std::string":        "string",
	}

	cxxMangle = strings.NewReplacer(
		":", "_",
		"<", "_",
		">", "_",
	)

	stdArrayRe    = regexp.MustCompile(" *std::array *<([a-zA-Z0-9:]+) *, *([0-9]+)> *")
	stdArrayDocRe = regexp.MustCompile(` *std::array *<([a-zA-Z0-9:]+) *, *([0-9]+)> *(\S+) *// *(.+)`)
)

type generator struct {
	w io.Writer

	buf *bytes.Buffer
	pkg string
	doc document

	comps map[string]string
	types map[string]string
	rules map[string]string
}

func newGenerator(w io.Writer, pkg, fname, rules string) (*generator, error) {
	var (
		err error
		gen = &generator{
			w:   w,
			buf: new(bytes.Buffer),
			pkg: pkg,

			comps: make(map[string]string),
			types: make(map[string]string),
			rules: make(map[string]string),
		}
	)

	raw, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("could not read input YAML file %q: %w", fname, err)
	}

	gen.doc.Options.Syntax = false
	gen.doc.Options.ExposePODMembers = true

	err = yaml.Unmarshal(raw, &gen.doc)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal YAML file %q: %w", fname, err)
	}

	for _, rule := range strings.Split(rules, ",") {
		toks := strings.Split(rule, "->")
		if len(toks) != 2 {
			continue
		}
		gen.rules[toks[0]] = toks[1]
	}

	gen.printf(`// Automatically generated. DO NOT EDIT.

package %s

`, pkg)

	return gen, nil
}

func (g *generator) generate() error {
	var err error

	for i := 0; i < len(g.doc.Components.Content); i += 2 {
		key := g.doc.Components.Content[i]
		val := g.doc.Components.Content[i+1]
		err = g.genComponent(key, val)
		if err != nil {
			return fmt.Errorf("could not handle component %q: %w", key.Value, err)
		}

	}

	for i := 0; i < len(g.doc.DataTypes.Content); i += 2 {
		key := g.doc.DataTypes.Content[i]
		val := g.doc.DataTypes.Content[i+1]
		err = g.genDataType(key, val)
		if err != nil {
			return fmt.Errorf("could not handle datatype %q: %w", key.Value, err)
		}
	}

	out, err := format.Source(g.buf.Bytes())
	if err != nil {
		return fmt.Errorf("could not go/format generated code: %w", err)
	}

	_, err = g.w.Write(out)
	if err != nil {
		return fmt.Errorf("could not write generated code: %w", err)
	}

	return nil
}

func (g *generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *generator) genTypeName(typ string) string {
	for k, v := range g.rules {
		typ = strings.Replace(typ, k, v, 1)
	}

	return g.cxx2go(typ)
}

func (g *generator) cxx2go(typ string) string {
	if v, ok := builtins[typ]; ok {
		return v
	}

	if grps := stdArrayRe.FindStringSubmatch(typ); grps != nil {
		return fmt.Sprintf("[%s]%s", grps[2], g.genTypeName(grps[1]))
	}

	return cxxMangle.Replace(typ)
}

func (g *generator) genComponent(knode, vnode *yaml.Node) error {
	var (
		err  error
		name = knode.Value
		typ  = g.genTypeName(name)
		doc  = strings.TrimSpace(knode.HeadComment)
	)

	handleDoc := func(doc string) string {
		return strings.TrimSpace(strings.Replace(doc, "#", "", 1))
	}

	doc = handleDoc(doc)
	if doc != "" {
		doc = "\n// " + doc
	}

	g.comps[name] = typ
	g.printf(`// %s%s
type %s struct {
`, name, doc, typ,
	)
	for i := 0; i < len(vnode.Content); i += 2 {
		key := vnode.Content[i]
		val := vnode.Content[i+1]
		if key.Value == "ExtraCode" {
			continue
		}
		doc := handleDoc(val.LineComment)
		if doc != "" {
			doc = "// " + doc
		}
		g.printf("\t%s %s%s\n", strings.Title(key.Value), g.genTypeName(val.Value), doc)
	}

	g.printf("}\n\n")

	return err
}

func (g *generator) genDataType(knode, vnode *yaml.Node) error {
	var (
		err  error
		name = knode.Value
		typ  = g.genTypeName(name)
	)

	g.types[name] = typ

	var dt dtype

	for i := 0; i < len(vnode.Content); i += 2 {
		key := vnode.Content[i]
		val := vnode.Content[i+1]
		switch key.Value {
		case "Description":
			dt.descr = val.Value
		case "Author":
			dt.author = val.Value
		case "Members":
			dt.members, err = g.genMembers(val)
			if err != nil {
				return fmt.Errorf("could not handle data type members of %q: %w",
					name, err,
				)
			}
		case "OneToOneRelations":
			dt.one2oneRels, err = g.genMembers(val)
			if err != nil {
				return fmt.Errorf("could not handle data type 1-to-1 relations of %q: %w",
					name, err,
				)
			}
		case "OneToManyRelations":
			dt.one2manyRels, err = g.genMembers(val)
			if err != nil {
				return fmt.Errorf("could not handle data type 1-to-n relations of %q: %w",
					name, err,
				)
			}
		case "ExtraCode":
			err = val.Decode(&dt.extraCode)
			if err != nil {
				return fmt.Errorf("could not handle data type extracode of %q: %w",
					name, err,
				)
			}
		case "ConstExtraCode":
			err = val.Decode(&dt.constCode)
			if err != nil {
				return fmt.Errorf("could not handle data type const-extracode of %q: %w",
					name, err,
				)
			}
		case "VectorMembers":
			dt.vecmbrs, err = g.genMembers(val)
			if err != nil {
				return fmt.Errorf("could not handle data type vec-members of %q: %w",
					name, err,
				)
			}
		default:
			return fmt.Errorf("unknown data type field %q", key.Value)
		}
	}

	doc := dt.descr
	if doc != "" {
		doc = "\n// " + dt.descr
	}
	g.printf(`// %s%s
type %s struct {
`, name, doc, typ,
	)
	for _, m := range dt.members {
		doc := ""
		if m.Doc != "" {
			doc = " // " + m.Doc
		}
		g.printf("\t%s %s%s\n", m.Name, m.Type, doc)
	}

	for _, m := range dt.one2oneRels {
		doc := ""
		if m.Doc != "" {
			doc = " // " + m.Doc
		}
		g.printf("\t%s *%s%s\n", m.Name, m.Type, doc)
	}

	for _, m := range dt.one2manyRels {
		doc := ""
		if m.Doc != "" {
			doc = " // " + m.Doc
		}
		g.printf("\t%s []*%s%s\n", m.Name, m.Type, doc)
	}

	for _, m := range dt.vecmbrs {
		doc := ""
		if m.Doc != "" {
			doc = " // " + m.Doc
		}
		g.printf("\t%s []%s%s\n", m.Name, m.Type, doc)
	}
	g.printf("}\n\n")

	return err
}

func (g *generator) genMembers(node *yaml.Node) ([]member, error) {
	mbrs := make([]member, 0, len(node.Content))
	for _, elem := range node.Content {
		mbr := g.genMember(strings.TrimSpace(elem.Value))
		mbr.Name = strings.Title(mbr.Name)
		mbrs = append(mbrs, mbr)
	}
	return mbrs, nil
}

func (g *generator) genMember(v string) member {
	var mbr member
	handleDoc := func(doc string) string {
		return strings.TrimSpace(strings.Replace(doc, "//", "", 1))
	}

	// remove spaces coming from builtins...
	for _, cxx := range []struct {
		k, v string
	}{
		{"unsigned long long", "uint64"},
		{"unsigned long", "uint64"},
		{"long long", "int64"},
		{"unsigned int", "uint32"},
		{"unsigned", "uint32"},
	} {
		if !strings.HasPrefix(v, cxx.k) {
			continue
		}
		v = strings.Replace(v, cxx.k, cxx.v, 1)
	}

	for k, cxx := range builtins {
		if !strings.HasPrefix(v, k) {
			continue
		}
		mbr.Type = cxx
		v = strings.TrimSpace(strings.Replace(v, k, "", 1))
		idx := strings.Index(v, " ")
		mbr.Name = strings.TrimSpace(v[:idx])
		mbr.Doc = handleDoc(strings.TrimSpace(v[idx:]))
		return mbr
	}

	if grps := stdArrayDocRe.FindStringSubmatch(v); grps != nil {
		mbr.Type = fmt.Sprintf("[%s]%s", grps[2], g.genTypeName(grps[1]))
		mbr.Name = grps[3]
		mbr.Doc = handleDoc(grps[4])
		return mbr
	}

	idx := strings.Index(v, " ")
	mbr.Type = g.genTypeName(strings.TrimSpace(v[:idx]))
	v = strings.TrimSpace(v[idx:])
	idx = strings.Index(v, "//")
	mbr.Name = strings.TrimSpace(v[:idx])
	mbr.Doc = handleDoc(strings.TrimSpace(v[idx:]))

	return mbr
}
