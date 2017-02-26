package cnf

// Formaul for now is just a list of clauses. This may change if data is
// required here for the solver.
type Formula []Clause

// NewFormulaFromInts is a helper to construct a Formula from a slice of
// int slices where each int slice represents a clause. See NewClauseFromInts
// for the clause structure information.
func NewFormulaFromInts(v [][]int) Formula {
	cs := make([]Clause, len(v))
	for i, raw := range v {
		cs[i] = NewClauseFromInts(raw)
	}

	return Formula(cs)
}
