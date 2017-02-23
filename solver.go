package sat

import (
	"github.com/mitchellh/go-sat/cnf"
)

// Solver is a SAT solver. This should be created manually with the
// exported fields set as documented.
type Solver struct {
	// Formula is the formula to be solved. Once solving has begun,
	// this shouldn't be changed. If you want to change the formula,
	// a new Solver should be allocated.
	Formula cnf.Formula

	// Trace, if set to true, will output trace debugging information
	// via the standard library `log` package. If true, Tracer must also
	// be set to a non-nil value. The easiest implmentation is a logger
	// created with log.NewLogger.
	Trace  bool
	Tracer Tracer

	m *trail
}

// Solve finds a solution for the formula, returning true on satisfiability.
func (s *Solver) Solve() bool {
	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: starting solver")
	}

	// Get the full list of vars
	varsF := s.Formula.Vars()

	// Create a new empty trail
	s.m = newTrail(len(varsF))

	for {
		// Perform unit propagation
		s.unitPropagate()

		conflictC := s.m.IsFormulaFalse(s.Formula)
		if !conflictC.IsZero() {
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: current trail contains negated formula: %s", s.m)
			}

			// If we have no more decisions within the trail, then we've
			// failed finding a satisfying value.
			if s.m.DecisionsLen() == 0 {
				return false
			}

			// Backtrack since we introduced an invalid literal
			l := s.m.TrimToLastDecision().Negate()
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: backtracking to %s, asserting %d", s.m, l)
			}
			s.m.Assert(l, false)
		} else {
			// If the trail contains the same number of elements as
			// the variables in the formula, then we've found a satisfaction.
			if s.m.Len() == len(varsF) {
				if s.Trace {
					s.Tracer.Printf("[TRACE] sat: solver found solution: %s", s.m)
				}

				return true
			}

			// Choose a literal to assert. For now we naively just select
			// the next literal.
			lit := selectLiteral(varsF, s.m)

			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: asserting: %d", lit)
			}

			s.m.Assert(lit, true)
		}
	}

	return false
}

func (s *Solver) unitPropagate() {
	for {
		for _, c := range s.Formula {
			for _, l := range c {
				if s.m.IsUnit(c, l) {
					if s.Trace {
						s.Tracer.Printf(
							"[TRACE] sat: found unit clause %v with literal %d in trail %s",
							c, l, s.m)
					}

					s.m.Assert(l, false)
					goto UNIT_REPEAT
				}
			}
		}

		// We didn't find a unit clause, close it out
		return

	UNIT_REPEAT:
		// We found a unit clause but we have to check if we violated
		// constraints in the trail.
		if !s.m.IsFormulaFalse(s.Formula).IsZero() {
			return
		}
	}
}

func selectLiteral(vars map[cnf.Literal]struct{}, t *trail) cnf.Literal {
	tMap := map[cnf.Literal]struct{}{}
	for _, e := range t.elems {
		lit := e.Lit
		if lit < 0 {
			lit = cnf.Literal(-int(lit))
		}

		tMap[lit] = struct{}{}
	}

	for k, _ := range vars {
		if _, ok := tMap[k]; !ok {
			return k
		}
	}

	return cnf.Literal(0)
}
