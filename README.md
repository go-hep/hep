pawgo
=====

`pawgo` is a nod to the old `PAW` physics analysis workstation.

## Installation

```sh
$ go get -u github.com/go-hep/pawgo
```

## Example

```
$ pawgo

:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D to quit.

paw> /?
/? 		-- print help
/file/close 	-- close a file
/file/create 	-- create file for write access
/file/list 	-- list a file's content
/file/open 	-- open file for read access
/hist/open 	-- open a histogram
/hist/plot 	-- plot a histogram
paw> /file/open f hsimple.rio
paw> /file/ls f
/file/id/f name=hsimple.rio
 	- h1	(type="*github.com/go-hep/hbook.H1D")
 	- h2	(type="*github.com/go-hep/hbook.H1D")

paw> /hist/open h /file/id/f/h1
paw> /hist/plot h
== h1d: name="h1"
entries=1000
mean=  -0.059
RMS=   +1.009
```
