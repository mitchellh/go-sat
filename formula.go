package sat

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

func (f Formula) Vars() []Literal {
	set := make(map[Literal]struct{})
	for _, c := range f {
		for _, l := range c {
			set[l] = struct{}{}
		}
	}

	result := make([]Literal, 0, len(set))
	for k, _ := range set {
		result = append(result, k)
	}

	return result
}
