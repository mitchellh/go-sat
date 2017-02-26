package sat

import (
	"fmt"

	"github.com/mitchellh/go-sat/cnf"
)

// watchClause should be called for any new clause added to the formula.
// This registers watches for the clause.
func (s *Solver) watchClause(c cnf.Clause) {
	c0 := c[0].Neg()
	c1 := c[1].Neg()

	if s.Trace {
		s.Tracer.Printf("[TRACE] sat: registering watchers for clause %s", c)
		s.Tracer.Printf("[TRACE] sat: when %s, check %s", c0, c[1])
		s.Tracer.Printf("[TRACE] sat: when %s, check %s", c1, c[0])
	}

	s.watches[c0] = append(s.watches[c0], &watcher{
		Clause: c,
		Lit:    c[1],
	})
	s.watches[c1] = append(s.watches[c1], &watcher{
		Clause: c,
		Lit:    c[0],
	})
}

// propagate performs unit propagation. This is made extremely efficient
// due to the watched literal algorithm. The core idea of watched literals
// is that a clause only needs to be checked for unit propagation if a
// watched literal is modified.
func (s *Solver) propagate() cnf.Clause {
	// qhead points to the first literal in the trail that we haven't
	// yet checked. This allows literal assertions to occur and only the
	// newly asserted literals (additions to the trail) need to be checked
	// for their affect on clauses.
	for s.qhead < len(s.trail) {
		// Get the next literal assigned in the trail
		p := s.trail[s.qhead]
		s.qhead++

		if s.Trace {
			s.Tracer.Printf("[TRACE] sat: looking for watches for: %s", p)
		}

		// Get the list of watches associated with this literal. We
		// maintain two indexes (i, j) because we'll be removing or
		// modifying watches.
		watches := s.watches[p]
		i := 0
		j := 0

		// This loop goes over each watch and checks if the clause has been
		// affected.
	PROP_LOOP:
		for i < len(watches) {
			iW := watches[i]
			if s.Trace {
				s.Tracer.Printf("[TRACE] sat: watcher: %s", iW)
			}

			// If the value of the literal being watched is true, then
			// this clause is already satisfied. We maintain the watch
			// in the watch list and continue.
			if s.ValueLit(iW.Lit) == True {
				if s.Trace {
					s.Tracer.Printf(
						"[TRACE] sat: watched lit %s became true; clause is true: %s",
						iW.Lit, iW.Clause)
				}

				watches[j] = watches[i]
				j++
				i++
				continue
			}

			// Increment
			i++

			// We're going to need ~p for the remainder
			pNeg := p.Neg()

			// We keep the false literal in c[1]
			iLits := iW.Clause
			first := iLits[0]
			if first == pNeg {
				if s.Trace {
					s.Tracer.Printf(
						"[TRACE] sat: moving false literal %s to position 1",
						iLits[0])
				}

				iLits[0], iLits[1] = iLits[1], pNeg
				first = iLits[0]
			}

			// newW is the new watcher that will replace the current
			// watcher no matter what. We always want to watch the
			// first literal.
			newW := &watcher{Clause: iW.Clause, Lit: first}

			// c[1] always contains the negated literal (above), so we
			// only have to check the first value to see if it is true.
			// We have the first != lit check since that is significantly
			// faster to fail than the ValueLit call (~1 or 2% faster) for
			// a common case.
			if first != iW.Lit && s.ValueLit(first) == True {
				watches[j] = newW
				j++
				continue
			}

			// At this point we know that no literal we're watching
			// is true. We have to check the remainder of the literals
			// in the clause: if they're all false then we have a unit
			// clause. If any are not false (true OR unassigned), then
			// we watch that literal and move on.
			for k := 2; k < len(iLits); k++ {
				if s.ValueLit(iLits[k]) != False {
					iLits[1], iLits[k] = iLits[k], pNeg
					i1 := iLits[1].Neg()
					s.watches[i1] = append(s.watches[i1], newW)
					continue PROP_LOOP
				}
			}

			// Every other value is false! This clause is either unit
			// or a conflict (unit if unassigned, conflict otherwise).
			watches[j] = newW
			j++

			// If it is false, then this is a conflict. We return immediately.
			if s.ValueLit(first) == False {
				// If i != j then we pruned some watches. We need to copy
				// the rest down. copy() is expensive so we avoid it if possible.
				if i != j {
					j += copy(watches[j:], watches[i:])
					s.watches[p] = watches[:j]
				}

				return iW.Clause
			}

			if s.Trace {
				s.Tracer.Printf(
					"[TRACE] sat: asserting unit literal %s in clause %s",
					first, iW.Clause)
			}

			// Assert the unit literal, caused by the clause it is part of
			s.assertLiteral(first, iW.Clause)
		}

		s.watches[p] = watches[:j]
	}

	// If we reached this point, we found no conflicts
	return nil
}

// watcher watches a single literal within a clause for the watched literal
// algorithm.
type watcher struct {
	Clause cnf.Clause // Clause that this literal is part of
	Lit    cnf.Lit    // Lit being watched
}

func (w *watcher) String() string {
	return fmt.Sprintf("watching lit %q in clause %s", w.Lit, w.Clause)
}
