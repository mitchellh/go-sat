package sat

import (
	"fmt"

	"github.com/mitchellh/go-sat/cnf"
)

func ExampleSolve() {
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
	s := New()
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
}
