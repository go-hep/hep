// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/shiny/screen"
)

func TestPawgo(t *testing.T) {
	var (
		stdout      = new(bytes.Buffer)
		scr         screen.Screen
		interactive bool
		args        []string
	)
	rc := xmain(stdout, scr, interactive, args)
	if rc != 0 {
		t.Fatalf("invalid exit-code: %d", rc)
	}
}

func TestPawgoScript(t *testing.T) {
	tmp, err := ioutil.TempDir("", "pawgo-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	var (
		stdout      = new(bytes.Buffer)
		scr         screen.Screen
		interactive bool
		fname       = path.Join(tmp, "script.paw")
		args        = []string{fname}
	)

	err = ioutil.WriteFile(fname, []byte(`## comment

## open the rio file hsimple.rio, assign it the variable name 'f'
/file/open f ./testdata/hsimple.rio

## list the content of the rio file 'f'
/file/ls f

## open the hbook.H1D histogram 'h1' from file 'f', assign it the variable name 'h'
/hist/open h /file/id/f/h1

`), 0644)

	rc := xmain(stdout, scr, interactive, args)
	if rc != 0 {
		t.Fatalf("invalid exit-code: %d", rc)
	}

	want := `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /file/open f ./testdata/hsimple.rio
# /file/ls f
/file/id/f name=./testdata/hsimple.rio
 	- h1	(type="*go-hep.org/x/hep/hbook.H1D")
 	- h2	(type="*go-hep.org/x/hep/hbook.H2D")
 	- p1	(type="*go-hep.org/x/hep/hbook.P1D")
 	- s2	(type="*go-hep.org/x/hep/hbook.S2D")

# /hist/open h /file/id/f/h1
bye.
`

	if got, want := stdout.String(), want; got != want {
		t.Fatalf("stdout differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(got),
			),
		)
	}
}

func TestPawgoShellCommand(t *testing.T) {
	tmp, err := ioutil.TempDir("", "pawgo-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	var (
		stdout      = new(bytes.Buffer)
		scr         screen.Screen
		interactive bool
		fname       = path.Join(tmp, "script.paw")
		args        = []string{fname}
	)

	script := "/! ls testdata\n"
	if runtime.GOOS == "windows" {
		script = "/! dir testdata\n"
	}

	err = ioutil.WriteFile(fname, []byte(script), 0644)
	if err != nil {
		t.Fatal(err)
	}

	rc := xmain(stdout, scr, interactive, args)
	if rc != 0 {
		t.Fatalf("invalid exit-code: %d", rc)
	}

	want := ""
	switch runtime.GOOS {
	case "windows":
		want = `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /! dir testdata
hsimple.rio  issue-120.paw  issue-120.rio  script.paw
bye.
`

	default:
		want = `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /! ls testdata
hsimple.rio
issue-120.paw
issue-120.rio
script.paw
bye.
`
	}

	if got, want := stdout.String(), want; got != want {
		t.Fatalf("stdout differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(got),
			),
		)
	}
}

func TestIssue120(t *testing.T) {
	var (
		stdout      = new(bytes.Buffer)
		scr         screen.Screen
		interactive bool
		args        = []string{"./testdata/issue-120.paw"}
	)

	rc := xmain(stdout, scr, interactive, args)
	if rc != 0 {
		t.Fatalf("invalid exit-code: %d", rc)
	}

	want := `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /file/open f ./testdata/issue-120.rio
# /hist/open h /file/id/f/MonoH_Truth/jets
bye.
`

	if got, want := stdout.String(), want; got != want {
		t.Fatalf("stdout differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(got),
			),
		)
	}
}
