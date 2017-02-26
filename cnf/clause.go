package cnf

// Clause represents a clause in a formula.
type Clause []Lit

// NewClauseFromInts creates a clause from a slice of integers. Each integer
// uniquely represents a literal. A negative integer represents a negated
// literal.
func NewClauseFromInts(v []int) Clause {
	lits := make([]Lit, len(v))
	for i, raw := range v {
		lits[i] = NewLitInt(raw)
	}

	return Clause(lits)
}

func (c Clause) Int() []int {
	result := make([]int, len(c))
	for i, l := range c {
		result[i] = l.Int()
	}

	return result
}
