package sat

import (
	"github.com/mitchellh/go-sat/cnf"
	"github.com/mitchellh/go-sat/packed"
)

// AddFormula adds the given formula to the solver.
//
// This can only be called before Solve() is called.
func (s *Solver) AddFormula(f packed.Formula) {
	for _, c := range f {
		s.AddClause(c)
	}
}

// AddClause adds a Clause to solve to the solver.
//
// This can only be called before Solve() is called.
func (s *Solver) AddClause(c *packed.Clause) {
	// Get the actual slice since we'll be modifying this directly.
	// The API docs say not to but its part of our package and we know
	// what we're doing. :)
	lits := c.Lits()

	// Keep track of an index since we'll be slicing as we go. We also
	// keep track of the last value so that we can find tautologies (X | !X)
	idx := 0
	last := packed.LitUndef
	for _, current := range lits {
		// Due to the sorting X and !X will always be next to each other.
		// A cheap way to check for tautologies is to just check the last
		// value.
		if current == last.Neg() {
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: addClause: not adding clause; tautology with var %s",
					current)
			}

			return
		}

		// Check if there is currently an assigned value of the literal.
		// If it is false then we already can skip this literal. If it is
		// true we can avoid adding the entire clause.
		switch s.ValueLit(current) {
		case False:
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: addClause: not adding literal; literal %s false: %s",
					current, c)
			}

			continue

		case True:
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: addClause: not adding clause; literal %s already true: %s",
					current, c)
			}

			return
		}

		// Due to sorting, we can quickly eliminate duplicates by only copying
		// down when the values aren't the same.
		if current != last {
			lits[idx] = current
			last = current
			idx++
		}
	}

	// Reset the size of literals to account for removed duplicates
	lits = lits[:idx]

	// If the clause is empty, then this formula can already not be satisfied
	if len(lits) == 0 {
		if s.Trace {
			s.Tracer.Printf("[TRACE] sat: addClause: empty clause, forcing unsat")
		}

		s.result = satResultUnsat
		return
	}

	// Update literals since we're for sure using this clause
	c.SetLits(lits)

	// Track the available decision variables
	for _, l := range lits {
		s.vars[l.Var()] = struct{}{}
	}

	// If this is a single literal clause then we assert it cause it must be
	if len(lits) == 1 {
		l := cnf.Literal(lits[0].Int())

		if s.Trace {
			s.Tracer.Printf("[TRACE] sat: addClause: single literal clause, asserting %d", l)
		}

		s.assertLiteral(lits[0])
		s.reasonMap[lits[0]] = c

		// Do unit propagation since this may solve already clauses
		s.unitPropagate()

		// We also don't add this clause since we just asserted the value
		return
	}

	// Add it to our formula
	s.clauses = append(s.clauses, *c)
}
