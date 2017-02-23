// Package cnf contains structures and operations for boolean formulas
// in conjunctive normal form (CNF).
package cnf

// Formula represents a Boolean formula in conjunctive normal form (CNF).
type Formula []Clause

// Clause is a disjunction of literals.
type Clause []Literal

// Literal is a single literal (such as X, Y, etc.).
//
// The value of the literal is a integer to uniquely identify this literal.
// Two literals with the same |value| (absolute value of their value) are
// the same literal. A negative literal indicates a negated literal.
type Literal int

// NewFormulaFromInts is a helper to turn [][]int into a Formula where
// the input is a slice of clauses and each int is a literal.
func NewFormulaFromInts(raw [][]int) Formula {
	cs := make([]Clause, len(raw))
	for i, v := range raw {
		cs[i] = make([]Literal, len(v))
		for j, x := range v {
			cs[i][j] = Literal(x)
		}
	}

	return Formula(cs)
}

func (f Formula) Ints() [][]int {
	result := make([][]int, len(f))
	for i, c := range f {
		result[i] = make([]int, len(c))
		for j, l := range c {
			result[i][j] = int(l)
		}
	}

	return result
}

// Negate the formula by negating all literals in all clauses.
func (f Formula) Negate() Formula {
	result := make([]Clause, len(f))
	for i, c := range f {
		resultC := make([]Literal, len(c))
		for j, l := range c {
			resultC[j] = Literal(-int(l))
		}

		result[i] = Clause(resultC)
	}

	return Formula(result)
}

func (f Formula) Vars() map[Literal]struct{} {
	set := make(map[Literal]struct{})
	for _, c := range f {
		for _, l := range c {
			if l < 0 {
				l = -l
			}

			set[l] = struct{}{}
		}
	}

	return set
}

// IsZero returns true of c represents the zero value.
func (c Clause) IsZero() bool {
	return c == nil
}

// Negate negates all the literals in C. The result is a clause that
// is NOT equivalent to the original Clause.
func (c Clause) Negate() Clause {
	resultC := make([]Literal, len(c))
	for j, l := range c {
		resultC[j] = Literal(-int(l))
	}

	return Clause(resultC)
}

func (l Literal) Negate() Literal {
	return Literal(int(l) * -1)
}
