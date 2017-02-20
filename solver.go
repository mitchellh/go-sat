package sat

import (
	"log"
)

var Trace = false

// Solve solves the given formula, returning ture on satisfiability and
// false on unsatisfiability. This is just temporary. We'll return the
// actual values for solving eventually.
func Solve(f Formula) bool {
	if Trace {
		log.Printf("[TRACE] sat: starting solver")
	}

	var m trail

	varsF := f.Vars()
	for {
		if m.IsFormulaFalse(f) {
			if Trace {
				log.Printf("[TRACE] sat: current trail contains negated formula: %s", m)
			}

			// If we have no more decisions within the trail, then we've
			// failed finding a satisfying value.
			if m.DecisionsLen() == 0 {
				return false
			}

			// Backtrack since we introduced an invalid literal
			l := m.TrimToLastDecision().Negate()
			if Trace {
				log.Printf("[TRACE] sat: backtracking to %s, asserting %d", m, l)
			}
			m.Assert(l, false)
		} else {
			// If the trail contains the same number of elements as
			// the variables in the formula, then we've found a satisfaction.
			if len(m) == len(varsF) {
				if Trace {
					log.Printf("[TRACE] sat: solver found solution: %s", m)
				}

				return true
			}

			// Choose a literal to assert. For now we naively just select
			// the next literal.
			lit := selectLiteral(varsF, m)

			if Trace {
				log.Printf("[TRACE] sat: asserting: %d", lit)
			}

			m.Assert(lit, true)
		}
	}

	return false
}

func selectLiteral(vars map[Literal]struct{}, t trail) Literal {
	tMap := map[Literal]struct{}{}
	for _, e := range t {
		lit := e.Lit
		if lit < 0 {
			lit = Literal(-int(lit))
		}

		tMap[lit] = struct{}{}
	}

	for k, _ := range vars {
		if _, ok := tMap[k]; !ok {
			return k
		}
	}

	return Literal(0)
}
