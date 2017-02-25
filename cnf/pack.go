package cnf

import (
	"github.com/mitchellh/go-sat/packed"
)

// Pack converts this formula to a packed Formula. The packed formula
// is a more solver-friendly representation of a formula at the expense
// of a more complicated representation.
func (f Formula) Pack() packed.Formula {
	cs := make([]packed.Clause, len(f))
	for i, rawC := range f {
		// Create the lits
		lits := make([]packed.Lit, len(rawC))
		for j, rawL := range rawC {
			lits[j] = rawL.Pack()
		}

		// Create the packed clause
		cs[i] = packed.Clause(lits)
	}

	return packed.Formula(cs)
}

func (l Literal) Pack() packed.Lit {
	return packed.NewLitInt(int(l))
}
