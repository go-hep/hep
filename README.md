pawgo
=====

`pawgo` is a nod to the old `PAW` physics analysis workstation.

## Installation

```sh
$ go get -u github.com/go-hep/pawgo
```

## Example

```sh
$ pawgo

:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /help for help.
^D to quit.

paw> /file/open 1 hsimple.rio
paw> /file/ls 1
/file/id/1 name=hsimple.rio
  h1
  h2
paw> /hist/open 11 /file/id/1/h1
paw> /hist/plot 11
== h1d: name="h1"
entries=1000
mean=  -0.059
RMS=   +1.009

```
