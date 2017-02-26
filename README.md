# SAT Solver written in Go

go-sat is a pure Go library for solving
[Boolean satisfiability problems (SAT)](https://en.wikipedia.org/wiki/Boolean_satisfiability_problem).

Solving SAT problems is at the core of a number of more difficult higher
level problems. From hardware/software verification, scheduling constraints,
version resolving, etc. many problems are reduced fully or partially to a basic
SAT problem.

Many existing SAT libraries exist that are efficient and easy to bind to
Unfortunately, cgo has a considerable runtime cost in addition to increasing
the complexity of cross-compilation. This library aims to be a pure Go SAT
solver that requires no cgo.

## Features

The root package is `sat` which contains the SAT solver itself. The SAT
solver is implemented with a standard iterative backtracking algorithm at its
core. To improve efficiency, the following features/heuristics are used by
this solver. Some of these obviously overlap:

  * [Unit propogation](https://en.wikipedia.org/wiki/Unit_propagation)
  * [Backjumping](https://en.wikipedia.org/wiki/Backjumping)
  * [Clause Learning](https://en.wikipedia.org/wiki/Conflict-Driven_Clause_Learning)
  * [Watched Literals](http://constraintmodelling.org/files/2015/07/GentJeffersonMiguelCP06.pdf)

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

## Example

Below is a basic example of using `go-sat`.

Note that the solver itself requires the formula already be in
CNF (conjunctive normal form). The solver expects higher level packages
to convert high level boolean expressions to this form prior to using
the solver.

```
// Consider the example formula already in CNF.
//
// ( ¬x1 ∨ ¬x3 ∨ ¬x4 ) ∧ ( x2 ∨ x3 ∨ ¬x4 ) ∧
// ( x1 ∨ ¬x2 ∨ x4 ) ∧ ( x1 ∨ x3 ∨ x4 ) ∧ ( ¬x1 ∨ x2 ∨ ¬x3 )
// (¬x4)

// Represent each variable as an int where a negative value means negated
formula := cnf.NewFormulaFromInts([][]int{
	[]int{-1, -3, -4},
	[]int{2, 3, -4},
	[]int{1, -2, 4},
	[]int{1, 3, 4},
	[]int{-1, 2, -3},
	[]int{-4},
})

// Create a solver and add the formulas to solve
s := sat.New()
s.AddFormula(formula)

// For low level details on how a solution came to be:
// s.Trace = true
// s.Tracer = log.New(os.Stderr, "", log.LstdFlags)

// Solve it!
sat := s.Solve()

// Solution can be read from Assignments. The key is the variable
// (always positive) and the value is the value.
solution := s.Assignments()

fmt.Printf("Solved: %v\n", sat)
fmt.Printf("Solution:\n")
fmt.Printf("  x1 = %v\n", solution[1])
fmt.Printf("  x2 = %v\n", solution[2])
fmt.Printf("  x3 = %v\n", solution[3])
fmt.Printf("  x4 = %v\n", solution[4])
// Output:
// Solved: true
// Solution:
//   x1 = true
//   x2 = true
//   x3 = true
//   x4 = false
```

## Issues and Contributing

If you find an issue with this library, please report an issue. If you'd like,
we welcome any contributions. Fork this library and submit a pull request.

SAT solving is an intensely competitive and fast-moving area of research. As
advancements are made, I welcome any contributions or recommendations to
improve this solver.

## Thanks

Thanks to [Minisat](http://minisat.se/) for providing understandable
and efficient implementations of SAT solving concepts. go-sat translates many
of their data representations and algorithms directly due to the clarity
of their implementation.
