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
)

func TestPawgo(t *testing.T) {
	var (
		stdout      = new(bytes.Buffer)
		interactive bool
		args        []string
	)

	rc := xmain(stdout, interactive, args)
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

	for _, tc := range []struct {
		name        string
		script      string
		want        string
		interactive bool
	}{
		{
			name: "basic",
			script: `## comment

## open the rio file hsimple.rio, assign it the variable name 'f'
/file/open f ./testdata/hsimple.rio

## list the content of the rio file 'f'
/file/ls f

## open the hbook.H1D histogram 'h1' from file 'f', assign it the variable name 'h'
/hist/open h /file/id/f/h1

/file/close f
`,
			want: `
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
# /file/close f
bye.
`,
			interactive: false,
		},
		{
			name:   "help",
			script: "/?\n/? /file/open",
			want: `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /?
/! 		-- run a shell command
/? 		-- print help
/file/close 	-- close a file
/file/create 	-- create file for write access
/file/ls 	-- list a file's content
/file/open 	-- open file for read access
/hist/open 	-- open a histogram
/hist/plot 	-- plot a histogram
/quit 		-- quit PAW-Go
# /? /file/open
/file/open 	-- open file for read access
bye.
`,
			interactive: false,
		},
		{
			name:   "quit",
			script: "/quit\n",
			want: `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /quit
bye.
`,
			interactive: false,
		},
		{
			name: "hplot-cmd",
			script: `## comment

/file/open f ./testdata/hsimple.rio
/hist/open h1 /file/id/f/h1
/hist/plot h1
/hist/open h2 /file/id/f/h2
/hist/plot h2
/quit
`,
			want: `
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

# /file/open f ./testdata/hsimple.rio
# /hist/open h1 /file/id/f/h1
# /hist/plot h1
== h1d: name=""
entries=10000
mean=  +0.004
RMS=   +1.005
# /hist/open h2 /file/id/f/h2
# /hist/plot h2
== h2d: name=""
entries=10000
xmean=  +0.027
xRMS=   +2.003
ymean=  +0.992
yRMS=   +1.723
# /quit
bye.
`,
			interactive: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				stdout = new(bytes.Buffer)
				fname  = path.Join(tmp, tc.name+".paw")
				args   = []string{fname}
			)

			err = ioutil.WriteFile(fname, []byte(tc.script), 0644)

			rc := xmain(stdout, tc.interactive, args)
			if rc != 0 {
				t.Fatalf("invalid exit-code: %d", rc)
			}

			if got, want := stdout.String(), tc.want; got != want {
				t.Fatalf("stdout differ:\n%s\n",
					cmp.Diff(
						string(want),
						string(got),
					),
				)
			}
		})
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

	rc := xmain(stdout, interactive, args)
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
		interactive bool
		args        = []string{"./testdata/issue-120.paw"}
	)

	rc := xmain(stdout, interactive, args)
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
