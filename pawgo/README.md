pawgo
=====

`pawgo` is a nod to the old `PAW` physics analysis workstation.

## Installation

```sh
$ go get -u go-hep.org/x/hep/pawgo
```

## Example

```
$ pawgo

:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

paw> /?
/!              -- run a shell command
/?              -- print help
/file/close     -- close a file
/file/create    -- create file for write access
/file/list      -- list a file's content
/file/open      -- open file for read access
/hist/open      -- open a histogram
/hist/plot      -- plot a histogram
/quit           -- quit PAW-Go

paw> /file/open f testdata/hsimple.rio
paw> /file/ls f
/file/id/f name=testdata/hsimple.rio
 	- h1	(type="*go-hep.org/x/hep/hbook.H1D")
 	- h2	(type="*go-hep.org/x/hep/hbook.H1D")

paw> /hist/open h /file/id/f/h1
paw> /hist/plot h
== h1d: name="h1"
entries=1000
mean=  -0.059
RMS=   +1.009
```
