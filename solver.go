package sat

import (
	"fmt"

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

	// decideLiterals is to be set by tests to force a certain decision
	// literal ordering. This can be used to exercise specific solver
	// behavior being tested.
	decideLiterals []int

	// Internal fields, do not set
	m         *trail
	c, cNeg   cnf.Clause
	reasonMap map[cnf.Literal]cnf.Clause
}

// Solve finds a solution for the formula, returning true on satisfiability.
func (s *Solver) Solve() bool {
	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: starting solver")
	}

	// Get the full list of vars
	varsF := s.Formula.Vars()

	// Create a new empty trail
	s.reasonMap = make(map[cnf.Literal]cnf.Clause)
	s.m = newTrail(len(varsF))

	// Copy f that will hold our learned clauses
	fCopy := make([]cnf.Clause, len(s.Formula))
	copy(fCopy, s.Formula)
	f := cnf.Formula(fCopy)

	for {
		// Perform unit propagation
		s.unitPropagate(f)

		conflictC := s.m.IsFormulaFalse(f)
		if !conflictC.IsZero() {
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: current trail contains negated formula: %s", s.m)
				s.Tracer.Printf("[TRACE] sat: conflict clause: %#v", conflictC)
			}

			// Set our conflict clause
			s.c = conflictC
			s.cNeg = s.c.Negate()

			// If we have no more decisions within the trail, then we've
			// failed finding a satisfying value.
			if s.m.DecisionsLen() == 0 {
				return false
			}

			// Explain to learn our conflict clause
			s.applyExplainUIP()
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: learned clause: %#v", s.c)
			}
			f = append(f, s.c)
			s.applyBackjump()

			/*
				// Backtrack since we introduced an invalid literal
				l := s.m.TrimToLastDecision().Negate()
				if s.Trace {
					s.Tracer.Printf("[TRACE] sat: backtracking to %s, asserting %d", s.m, l)
				}
				s.m.Assert(l, false)
				s.reasonMap[l] = s.c
			*/
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
			lit := s.selectLiteral(varsF)

			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: assert: %d (decision)", lit)
			}

			s.m.Assert(lit, true)
		}
	}

	return false
}

func (s *Solver) selectLiteral(vars map[cnf.Literal]struct{}) cnf.Literal {
	tMap := map[cnf.Literal]struct{}{}
	for _, e := range s.m.elems {
		lit := e.Lit
		if lit < 0 {
			lit = cnf.Literal(-int(lit))
		}

		tMap[lit] = struct{}{}
	}

	if len(s.decideLiterals) > 0 {
		result := cnf.Literal(s.decideLiterals[0])
		s.decideLiterals = s.decideLiterals[1:]

		if _, ok := tMap[result]; ok {
			panic(fmt.Sprintf("decideLiteral taken: %d", result))
		}

		return result
	}

	for k, _ := range vars {
		if _, ok := tMap[k]; !ok {
			return k
		}
	}

	return cnf.Literal(0)
}

func (s *Solver) unitPropagate(f cnf.Formula) {
	for {
		for _, c := range f {
			for _, l := range c {
				if s.m.IsUnit(c, l) {
					if s.Trace {
						s.Tracer.Printf(
							"[TRACE] sat: found unit clause %v with literal %d in trail %s",
							c, l, s.m)
					}

					s.m.Assert(l, false)
					s.reasonMap[l] = c
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

func (s *Solver) applyExplain(lit cnf.Literal) {
	litNeg := lit.Negate()
	reason := s.reasonMap[lit]

	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: applyExplain: lit = %d, reason = %#v", lit, reason)
	}

	resultMap := make(map[cnf.Literal]struct{})
	for _, l := range s.c {
		if l != litNeg {
			resultMap[l] = struct{}{}
		}
	}
	for _, l := range reason {
		if l != lit {
			resultMap[l] = struct{}{}
		}
	}

	result := make([]cnf.Literal, 0, len(resultMap))
	for k, _ := range resultMap {
		result = append(result, k)
	}

	s.c = cnf.Clause(result)
	s.cNeg = s.c.Negate()

	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: applyExplain: new C = %#v", s.c)
	}
}

func (s *Solver) applyExplainUIP() {
	for {
		lit, isUIP := s.isUIP()
		if isUIP {
			return
		}

		s.applyExplain(lit)
	}
}

func (s *Solver) isUIP() (cnf.Literal, bool) {
	lit := s.m.LastAssertedLiteral(s.cNeg)
	litLevel := s.m.Level(lit)
	for _, l := range s.cNeg {
		// Literal must not equal the last asserted lit
		if l == lit {
			continue
		}

		// If these two literals at the same level, then it isn't a UIP
		if s.m.Level(l) == litLevel {
			return lit, false
		}
	}

	return lit, true
}

func (s *Solver) applyBackjump() {
	lit := s.m.LastAssertedLiteral(s.cNeg)
	c := make([]cnf.Literal, 0, len(s.cNeg))
	for _, l := range s.cNeg {
		if l != lit {
			c = append(c, l)
		}
	}

	level := s.m.MaxLevel(cnf.Clause(c))
	if s.Trace {
		s.Tracer.Printf(
			"[TRACE] sat: backjump. C = %#v; l = %d; level = %d",
			s.c, lit, level)
	}

	s.m.TrimToLevel(level)

	lit = lit.Negate()
	s.m.Assert(lit, false)
	s.reasonMap[lit] = s.c

	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: backjump. M = %s", s.m)
	}
}
