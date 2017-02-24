package sat

import (
	"github.com/mitchellh/go-sat/cnf"
)

// AddFormula adds the given formula to the solver.
//
// This can only be called before Solve() is called.
func (s *Solver) AddFormula(f cnf.Formula) {
	for _, c := range f {
		s.AddClause(c)
	}
}

// AddClause adds a Clause to solve to the solver.
//
// This can only be called before Solve() is called.
func (s *Solver) AddClause(c cnf.Clause) {
	ls := make(map[cnf.Literal]struct{})
	for _, l := range c {
		// If this literal is already false in the trail, then don't add
		if s.m.IsLiteralFalse(l) {
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: addClause: not adding literal; literal %d false: %#v",
					l, c)
			}

			continue
		}

		// If the literal is already true, we don't add the clause at all
		if s.m.IsLiteralTrue(l) {
			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: addClause: not adding clause; literal %d already true: %#v",
					l, c)
			}

			return
		}

		// If the clause contains both a positive and negative it is
		// tautological.
		if _, ok := ls[l.Negate()]; ok {
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: addClause: not adding clause; tautology: %#v", c)
			}

			return
		}

		// Add the literal. This will also remove duplicates
		ls[l] = struct{}{}
	}

	if len(ls) == 0 {
		if s.Trace {
			s.Tracer.Printf("[TRACE] sat: addClause: empty clause, forcing unsat")
		}

		s.result = satResultUnsat
		return
	}

	// If this is a single literal clause then we assert it cause it must be
	if len(ls) == 1 {
		for l, _ := range ls {
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: addClause: single literal clause, asserting %d", l)
			}

			s.assertLiteral(l, false)
			s.reasonMap[l] = c

			// Do unit propagation since this may solve already clauses
			s.unitPropagate()
		}

		// We also don't add this clause since we just asserted the value
		return
	}

	// Add it to our formula
	c = make([]cnf.Literal, 0, len(ls))
	for l, _ := range ls {
		c = append(c, cnf.Literal(l))
	}

	s.f = append(s.f, c)
}
