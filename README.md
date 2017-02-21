# SAT Solver written in Go

go-sat is a pure Go library for solving
[Boolean satisfiability problems (SAT)](https://en.wikipedia.org/wiki/Boolean_satisfiability_problem).

Solving SAT problems is at the core of a number of more difficult higher
level problems. From hardware/software verification, scheduling constraints,
version resolving, etc. many problems are reduced fully or partially to a basic
SAT problem.

## Features

The root package is `sat` which contains the SAT solver itself. The SAT
solver is implemented with a standard iterative backtracking algorithm at its
core. To improve efficiency, the following features are supported by this
solver:

  * [Unit propogation](https://en.wikipedia.org/wiki/Unit_propagation)

  * ... more

In addition to the solver, this library contains a number of sub-packages
for working with SAT problems and formulas:

  * `cnf` - Data structure to represent and perform operations on a boolean
    formula in [conjunctive normal form](https://en.wikipedia.org/wiki/Conjunctive_normal_form).

  * `dimacs` - A parser for the [DIMACS CNF format](http://www.domagoj-babic.com/uploads/ResearchProjects/Spear/dimacs-cnf.pdf),
    a widely accepted format for boolean formulas in [CNF](https://en.wikipedia.org/wiki/Conjunctive_normal_form).

## Installation

go-sat is a standard Go library that can be installed and referenced using
standard `go get` and import paths. For example, the root `sat` package:

    go get github.com/mitchellh/go-sat

## Issues and Contributing

If you find an issue with this library, please report an issue. If you'd like,
we welcome any contributions. Fork this library and submit a pull request.

SAT solving is an intensely competitive and fast-moving area of research. As
advancements are made, I welcome any contributions or recommendations to
improve this solver.
