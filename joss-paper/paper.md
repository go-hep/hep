---
title: 'Go-HEP: libraries for High Energy Physics analyses in Go'
tags:
  - Go
  - ROOT
  - CERN
  - Gonum
authors:
  - name: Sebastien Binet
    orcid: 0000-0003-4913-6104
    affiliation: 1
  - name: Bastian Wieck
  - name: David Blyth
    affiliation: 2
  - name: Emmanuel Busato
    affiliation: 1
  - name: Michaël Ughetto
  - name: Peter Waller
affiliations:
  - name: IN2P3
    index: 1
  - name: Argonne National Laboratory
    index: 2
date: 15 August 2017
bibliography: paper.bib
---

# Summary

Go-HEP provides tools to interface with CERN's ROOT [@ROOT] software
framework and carry analyses or data acquisition from within the Go [@Go]
programming language.

Go-HEP exposes libraries to read and write common High Energy Physics (HEP)
file formats (HepMC [@HepMC], LHEF [@LHEF], SLHA [@SLHA]) but, at the
moment, only read interoperability is achieved with ROOT file format.
Go-HEP also provides tools to carry statistical analyses in the form of
1-dim and 2-dim histograms, 1-dim and 2-dim scatters and n-tuples.
Go-HEP can also graphically represent these results, leveraging the
Gonum [@Gonum] plotting library.

# References

